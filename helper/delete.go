package helper

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

func Delete(path string, wg *sync.WaitGroup) {
	defer wg.Done()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		Log(fmt.Sprintf("Error reading directory %s: %s", path, err))
		return
	}

	for _, file := range files {
		curPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			wg.Add(1)
			Delete(curPath, wg)
		} else {
			if err := os.Remove(curPath); err != nil {
				Log(fmt.Sprintf("Error deleting file %s: %s", curPath, err))
			}
		}
	}

	if err := os.Remove(path); err != nil {
		Log(fmt.Sprintf("Error deleting directory %s: %s", path, err))
	}
}
