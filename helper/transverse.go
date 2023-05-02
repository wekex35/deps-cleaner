package helper

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"runtime"
)

func Transverse(startPath string, filter *regexp.Regexp, callback func(string)) {
	files, err := ioutil.ReadDir(startPath)
	if err != nil {
		Log(fmt.Sprintln("Directory Not Found:", startPath))
		return
	}

	done := make(chan struct{})
	defer close(done)

	// Use a worker pool to traverse the directory tree in parallel
	workers := runtime.NumCPU()
	paths := make(chan string, workers)

	go func() {
		for _, file := range files {
			filename := filepath.Join(startPath, file.Name())
			if file.IsDir() {
				if filter.MatchString(filename) {
					callback(filename)
				} else {
					paths <- filename
				}
			}
		}
		close(paths)
	}()

	for i := 0; i < workers; i++ {
		go func() {
			for {
				select {
				case path, ok := <-paths:
					if !ok {
						return
					}
					Transverse(path, filter, callback)
				case <-done:
					return
				}
			}
		}()
	}
}
