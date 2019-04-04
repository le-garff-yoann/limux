package in

import (
	"limux/file/availabilityness"
	"limux/helper"
	"limux/processor"
	"fmt"
	"os"

	logger "github.com/apsdehal/go-logger"
	"github.com/fsnotify/fsnotify"
	"github.com/mholt/archiver"
)

// In check if there are any files in Src.
// If so, he detar and send them to Dst.
type In struct {
	watcher *fsnotify.Watcher // TODO: replace this with a time.Ticker (refactor parts of Out/In codebase).
	Src     string            `json:"src" yaml:"src"`
	Dst     string            `json:"dst" yaml:"dst"`
}

// Configure implements `processor.Processor` interface.
func (s *In) Configure() error {
	for _, path := range []string{s.Src, s.Dst} {
		val, err := helper.IsDir(path)
		if err != nil {
			return processor.EnvironmentError(fmt.Sprintf("%s Cannot stat %s: %s", s, path, err))
		}
		if !val {
			return processor.EnvironmentError(fmt.Sprintf("%s %s is not a directory", s, path))
		}
	}

	return nil
}

// Start implements `processor.Processor` interface.
func (s *In) Start(events chan (processor.Event)) error {
	var err error
	s.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	s.watcher.Add(s.Src)

	go func() {
		for {
			e, ok := <-s.watcher.Events
			if !ok {
				return
			}

			if e.Op&fsnotify.Create == fsnotify.Create {
				events <- processor.Event{Processor: s, Type: processor.Bgn, Message: fmt.Sprintf("%s created", e.Name), Level: logger.DebugLevel}

				go func() {
					fanessFiles, fanessLogs := availabilityness.Watch(e.Name)

					for {
						select {
						case f, ok := <-fanessFiles:
							if !ok {
								return
							}

							if f.File == nil {
								events <- processor.Event{Processor: s, Type: processor.Log, Message: fmt.Sprintf("The file %s no longer exists", f.Fname), Level: logger.InfoLevel}

								return
							}

							if err := f.File.Close(); err != nil {
								events <- processor.Event{Processor: s, Type: processor.Log, Message: err.Error(), Level: logger.WarningLevel}
							}

							events <- processor.Event{Processor: s, Type: processor.Log, Message: fmt.Sprintf("Extracting the file %s", f.Fname), Level: logger.NoticeLevel}

							tarer := archiver.Tar{
								OverwriteExisting: true,
							}

							if err := tarer.Unarchive(f.Fname, s.Dst); err == nil {
								if err = os.RemoveAll(f.Fname); err == nil {
									events <- processor.Event{Processor: s, Type: processor.Fin, Message: fmt.Sprintf("The processing for %s is finished", f.Fname), Level: logger.NoticeLevel}
								} else {
									events <- processor.Event{Processor: s, Type: processor.Fin, Message: err.Error(), Level: logger.ErrorLevel}
								}
							} else {
								events <- processor.Event{Processor: s, Type: processor.Fin, Message: err.Error(), Level: logger.ErrorLevel}
							}

							if err := tarer.Close(); err != nil {
								events <- processor.Event{Processor: s, Type: processor.Log, Message: err.Error(), Level: logger.InfoLevel}
							}

							return
						case e := <-fanessLogs:
							events <- processor.Event{Processor: s, Type: processor.Log, Message: e, Level: logger.InfoLevel}
						}
					}
				}()
			}
		}
	}()

	return nil
}

// Stop implements `processor.Processor` interface.
func (s *In) Stop() error {
	if err := s.watcher.Remove(s.Src); err != nil {
		return err
	}

	s.watcher.Close()

	return nil
}

func (s In) String() string {
	return fmt.Sprintf("[in: %s -> %s]", s.Src, s.Dst)
}
