package utils

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/jjzcru/elk/pkg/primitives/ox"
)

// GetElk get an ox pointer from a file path
func GetElk(filePath string, isGlobal bool) (*ox.Elk, error) {
	var elkConfigPath string
	var err error

	if len(filePath) > 0 {
		elkConfigPath = filePath
	} else {
		elkConfigPath, err = getElkFilePath(isGlobal)
		if err != nil {
			return nil, err
		}
	}

	response, err := ox.FromFile(elkConfigPath)
	if err != nil {
		return nil, err
	}

	response.SetFilePath(elkConfigPath)

	return response, nil
}

// SetElk saves elk object in a file
func SetElk(elk *ox.Elk, filePath string) error {
	return ox.ToFile(elk, filePath)
}

func getElkFilePath(isGlobal bool) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var elkFilePath string

	if isGlobal {
		elkFilePath, err = getGlobalElkFile()
		if err != nil {
			return "", err
		}
		return elkFilePath, nil
	}

	isLocal := isLocalElkFile(path.Join(dir, "ox.yml"))
	if isLocal {
		elkFilePath = path.Join(dir, "ox.yml")
	} else {
		elkFilePath, err = getGlobalElkFile()
		if err != nil {
			return "", err
		}
	}

	return elkFilePath, nil
}

func isLocalElkFile(localDirectory string) bool {
	if _, err := os.Stat(localDirectory); os.IsNotExist(err) {
		return false
	}
	return true
}

func getGlobalElkFile() (string, error) {
	globalElkFilePath := os.Getenv("ELK_FILE")
	if len(globalElkFilePath) > 0 {
		isAFile, err := IsPathAFile(globalElkFilePath)
		if err != nil {
			return "", err
		}

		if !isAFile {
			return "", errors.New("ELK_FILE path must be a file")
		}

		return globalElkFilePath, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	globalElkFilePath = path.Join(usr.HomeDir, "ox.yml")
	if _, err := os.Stat(globalElkFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("default global path %s do not exist, please create it or set the env variable ELK_FILE", globalElkFilePath)
	}

	return globalElkFilePath, nil
}
