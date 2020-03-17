package config

import (
	"errors"
	"fmt"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/primitives/elk"
	"os"
	"os/user"
	"path"
)

// GetElk get an elk pointer from a file path
func GetElk(filePath string, isGlobal bool) (*elk.Elk, error) {
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

	return elk.FromFile(elkConfigPath)
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

	isLocal := isLocalElkFile(path.Join(dir, "elk.yml"))
	if isLocal {
		elkFilePath = path.Join(dir, "elk.yml")
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
		isAFile, err := utils.IsPathAFile(globalElkFilePath)
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

	globalElkFilePath = path.Join(usr.HomeDir, "elk.yml")
	if _, err := os.Stat(globalElkFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("default global path %s do not exist, please create it or set the env variable ELK_FILE", globalElkFilePath)
	}

	return globalElkFilePath, nil
}
