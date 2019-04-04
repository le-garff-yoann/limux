package availabilityness

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	// CheckInterval is the verification interval of each file of a watcher.
	CheckInterval = 500 * time.Millisecond

	// CheckAttempts is the number of times that a file should not have "changed".
	CheckAttempts = 10
)

type sizeCounter struct {
	size    int64
	attemps int
}

// TrulyFullInfo is basically os.File with the full path not far.
type TrulyFullInfo struct {
	Fname string
	File  *os.File
}

// Watch monitors a file list and warns if they are ready for reading.
func Watch(files ...string) (<-chan TrulyFullInfo, <-chan string) {
	fileQueue := make(chan TrulyFullInfo)
	logs := make(chan string)

	var wg sync.WaitGroup

	wg.Add(len(files))

	for _, fname := range files {
		go func(fname string) {
			sizeCounter := &sizeCounter{0, 0}

			tk := time.NewTicker(CheckInterval)
			defer func() {
				tk.Stop()

				wg.Done()
			}()

			for range tk.C {
				file, err := os.Open(fname)
				if err != nil {
					if os.IsNotExist(err) || os.IsPermission(err) {
						fileQueue <- TrulyFullInfo{fname, nil}

						return
					}

					logs <- fmt.Sprint(err)
				}

				if fi, err := file.Stat(); err == nil {
					if fi.IsDir() {
						fileQueue <- TrulyFullInfo{fname, nil}

						return
					}

					csize := fi.Size()

					logs <- fmt.Sprintf("Current size (attempt %d) of %s is %d", sizeCounter.attemps, fname, csize)

					if csize == sizeCounter.size {
						sizeCounter.attemps++
					} else {
						sizeCounter.attemps = 0
					}

					if sizeCounter.attemps > CheckAttempts {
						fileQueue <- TrulyFullInfo{fname, file}

						return
					}

					sizeCounter.size = csize
				} else {
					logs <- fmt.Sprint(err)
				}

				if err := file.Close(); err != nil {
					logs <- fmt.Sprint(err)
				}
			}
		}(fname)
	}

	go func() {
		wg.Wait()

		close(fileQueue)
		close(logs)
	}()

	return fileQueue, logs
}
