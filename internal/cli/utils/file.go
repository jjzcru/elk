package utils

import (
	"fmt"
	"os"
)

// IsFileExist returns a bool depending if a file exist or not
func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// IsPathADir returns a bool depending if a path is a dir and an error if the file do not exist
func IsPathADir(path string) (bool, error) {
	if !IsPathExist(path) {
		return false, fmt.Errorf("file \"%s\" do not exist", path)
	}

	info, _ := os.Stat(path)
	return info.IsDir(), nil
}

// IsPathAFile returns a bool depending if a path is a file and an error if the file do not exist
func IsPathAFile(path string) (bool, error) {
	isPathADir, err := IsPathADir(path)
	if err != nil {
		return false, err
	}

	return !isPathADir, nil
}
