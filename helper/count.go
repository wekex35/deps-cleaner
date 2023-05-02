package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func Count(path string, pattern *regexp.Regexp) (int, error) {
	// Check if the provided path exists and is a directory
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if !info.IsDir() {
		return 0, fmt.Errorf("%s is not a directory", path)
	}

	// Check if the provided regular expression is valid
	if pattern == nil {
		return 0, fmt.Errorf("invalid regular expression")
	}

	var count int
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && pattern.MatchString(info.Name()) {
			count++
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}
