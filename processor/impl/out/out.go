package out

import (
	"bytes"
	"filemux/file/availabilityness"
	"filemux/helper"
	"filemux/processor"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	logger "github.com/apsdehal/go-logger"
	"github.com/mholt/archiver"
)

// DefaultInterval
const defaultInterval time.Duration = 5

// Out check if there are any files in Src at each Interval.
// If so, he tares them in a single archive and sends them to Dst.
// Exec will be executed if not null.
type Out struct {
	ticker              *time.Ticker
	Interval            *time.Duration `json:"interval" yaml:"interval"`
	Src                 string         `json:"src" yaml:"src"`
	Dst                 string         `json:"dst" yaml:"dst"`
	ArchiveBasename     string         `json:"archive_basename" yaml:"archive_basename"`
	ArchiveInnerDirname string         `json:"archive_inner_dirname" yaml:"archive_inner_dirname"`
	Exec                *[]string      `json:"exec" yaml:"exec"`
}

// Configure implements `processor.Processor` interface.
func (s *Out) Configure() error {
	if s.Exec == nil {
		var emptySlice []string

		s.Exec = &emptySlice
	}

	if s.Interval == nil {
		*s.Interval = defaultInterval
	}

	if len(s.ArchiveBasename) == 0 {
		return fmt.Errorf("%s base_dir_tpl len() should be superior to 0", s)
	}

	var tpl bytes.Buffer
	if err := helper.ParseTemplate(s.ArchiveBasename, &tpl, map[string]interface{}{
		"Now": time.Now(),
	}); err != nil {
		return fmt.Errorf("%s Cannot render %s: %s", s, s.ArchiveBasename, err)
	}

	for i, arg := range *s.Exec {
		if err := helper.ParseTemplate(arg, &tpl, map[string]interface{}{
			"ArchiveFullPath": s.Dst,
		}); err != nil {
			return fmt.Errorf("%s Cannot render exec[%d] %s: %s", s, i, arg, err)
		}
	}

	val, err := helper.IsDir(s.Dst)
	if err != nil {
		return processor.EnvironmentError(fmt.Sprintf("%s Cannot stat %s: %s", s, s.Dst, err))
	}
	if !val {
		return processor.EnvironmentError(fmt.Sprintf("%s %s is not a directory", s, s.Dst))
	}

	return nil
}

// Start implements `processor.Processor` interface.
func (s *Out) Start(events chan (processor.Event)) error {
	s.ticker = time.NewTicker(*s.Interval)

	go func() {
		defer s.ticker.Stop()

		for range s.ticker.C { // Should block.
			events <- processor.Event{Processor: s, Type: processor.Log, Message: "Looking for files", Level: logger.InfoLevel}

			if files, err := filepath.Glob(filepath.Join(s.Src, "*")); err == nil && len(files) > 0 {
				events <- processor.Event{Processor: s, Type: processor.Bgn, Message: "Found something", Level: logger.DebugLevel}

				fanessRdy, fanessLogs := availabilityness.Watch(files...)

				var (
					wg sync.WaitGroup

					tar        archiver.Tar
					arFullPath string
					arf        *os.File

					created      = false
					emptyArchive = true
				)

				wg.Add(len(files))

				go func() {
					for {
						select {
						case f, ok := <-fanessRdy:
							if !ok {
								return
							}

							// It block because archiver is NOT safe for concurrent use
							// (https://github.com/mholt/archiver/blob/master/archiver.go#L40)
							func() {
								defer wg.Done()

								if f.File == nil {
									events <- processor.Event{Processor: s, Type: processor.Log, Message: fmt.Sprintf("The file %s no longer exists", f.Fname), Level: logger.InfoLevel}

									return
								}

								defer func() {
									if err := f.File.Close(); err != nil {
										events <- processor.Event{Processor: s, Type: processor.Log, Message: err.Error(), Level: logger.InfoLevel}
									}
								}()

								fi, err := f.File.Stat()
								if err != nil {
									events <- processor.Event{Processor: s, Type: processor.Log, Message: fmt.Sprintf("Cannot stat the file %s: %s", f.Fname, err), Level: logger.ErrorLevel}

									return
								}

								if !created {
									var tpl bytes.Buffer
									helper.ParseTemplate(s.ArchiveBasename, &tpl, map[string]interface{}{
										"Now": time.Now(),
									})

									arFullPath = filepath.Join(s.Dst, tpl.String())
									arFullPath += ".tar"

									events <- processor.Event{Processor: s, Type: processor.Log, Message: fmt.Sprintf("Archiving file %s to %s", s.Src, arFullPath), Level: logger.NoticeLevel}

									tar = archiver.Tar{
										MkdirAll:          true,
										OverwriteExisting: true,
									}

									arf, err = os.Create(arFullPath)
									if err != nil {
										events <- processor.Event{Processor: s, Type: processor.Log, Message: err.Error(), Level: logger.ErrorLevel}

										return
									}

									if err = tar.Create(arf); err != nil {
										events <- processor.Event{Processor: s, Type: processor.Log, Message: err.Error(), Level: logger.ErrorLevel}

										return
									}

									created = true
								}

								internalName, err := archiver.NameInArchive(fi, f.Fname, f.Fname)
								if err != nil {
									events <- processor.Event{Processor: s, Type: processor.Log, Message: err.Error(), Level: logger.ErrorLevel}

									return
								}

								customName := internalName
								if s.ArchiveInnerDirname != "" {
									var tpl bytes.Buffer
									if err = helper.ParseTemplate(s.ArchiveInnerDirname, &tpl, map[string]interface{}{
										"SrcFragments": strings.Split(f.Fname, string(filepath.Separator)),
									}); err == nil {
										customName = filepath.Join(tpl.String(), internalName)
									} else {
										events <- processor.Event{Processor: s, Type: processor.Log, Message: fmt.Sprintf("Cannot render the template %s: %s", s.ArchiveInnerDirname, err), Level: logger.ErrorLevel}

										return
									}
								}

								err = tar.Write(archiver.File{
									FileInfo: archiver.FileInfo{
										FileInfo:   fi,
										CustomName: customName,
									},
									ReadCloser: f.File,
								})
								if err := f.File.Close(); err != nil {
									events <- processor.Event{Processor: s, Type: processor.Log, Message: err.Error(), Level: logger.InfoLevel}
								}

								if err == nil {
									emptyArchive = false

									events <- processor.Event{Processor: s, Type: processor.Log, Message: fmt.Sprintf("File %s written as %s to %s", f.Fname, customName, arFullPath), Level: logger.NoticeLevel}

									if err := os.Remove(f.Fname); err != nil {
										events <- processor.Event{Processor: s, Type: processor.Log, Message: fmt.Sprintf("Cannot remove the file %s: %s", f.Fname, err), Level: logger.ErrorLevel}
									}
								} else {
									events <- processor.Event{Processor: s, Type: processor.Log, Message: fmt.Sprintf("Cannot write %s to the archive: %s", f.Fname, err), Level: logger.ErrorLevel}
								}
							}()
						case e := <-fanessLogs:
							events <- processor.Event{Processor: s, Type: processor.Log, Message: e, Level: logger.InfoLevel}
						}
					}
				}()

				wg.Wait()

				if created {
					if err := arf.Close(); err != nil {
						events <- processor.Event{Processor: s, Type: processor.Log, Message: err.Error(), Level: logger.InfoLevel}
					}

					if emptyArchive {
						if err := os.Remove(arFullPath); err != nil {
							events <- processor.Event{Processor: s, Type: processor.Log, Message: fmt.Sprintf("Cannot remove the file %s: %s", arFullPath, err), Level: logger.ErrorLevel}
						}
					} else {
						execSlice := *s.Exec

						if len(execSlice) > 0 {
							cmd := exec.Command(execSlice[0])
							for _, arg := range execSlice[1:] {
								var tpl bytes.Buffer
								helper.ParseTemplate(arg, &tpl, map[string]interface{}{
									"ArchiveFullPath": arFullPath,
								})

								cmd.Args = append(cmd.Args, tpl.String())
							}

							if b, err := cmd.CombinedOutput(); err == nil {
								events <- processor.Event{Processor: s, Type: processor.Log, Message: fmt.Sprintf("Command %s has succeeded: %s", cmd.Args, string(b)), Level: logger.DebugLevel}
							} else {
								events <- processor.Event{Processor: s, Type: processor.Log, Message: fmt.Sprintf("Command %s has failed: %s", cmd.Args, string(b)), Level: logger.ErrorLevel}
							}
						}
					}
				}

				events <- processor.Event{Processor: s, Type: processor.Fin, Message: "The processing is finished", Level: logger.NoticeLevel}
			}
		}
	}()

	return nil
}

// Stop implements `processor.Processor` interface.
func (s *Out) Stop() error {
	s.ticker.Stop()

	return nil
}

func (s Out) String() string {
	return fmt.Sprintf("[out: %s -> %s]", s.Src, s.Dst)
}
