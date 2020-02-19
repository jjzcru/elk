package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/jjzcru/elk/pkg/engine"
	"github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"
)

func GetElk(filePath string, isGlobal bool) (*engine.Elk, error) {
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

	elk := engine.Elk{}
	if _, err := os.Stat(elkConfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("the path for elk.yml do not exist '%s'", elkConfigPath)
	}

	data, err := ioutil.ReadFile(elkConfigPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &elk)
	if err != nil {
		return nil, err
	}

	return &elk, nil
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
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	configPath := path.Join(usr.HomeDir, ".elk", "config.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", errors.New("Elk is not installed, please run \nelk install")
	}

	config := engine.Config{}

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return "", err
	}

	return config.Path, nil
}

// PrintError displays a error message
func PrintError(err string) {
	fmt.Print(aurora.Bold(aurora.Red("ERROR: ")))
	_, _ = fmt.Fprintf(os.Stderr, err)
	fmt.Println()
}
