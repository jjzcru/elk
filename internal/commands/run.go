package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"sync"

	"github.com/jjzcru/elk/pkg/engine"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run one or more task in a terminal",
	Run: func(cmd *cobra.Command, args []string) {
		isGlobal, err := cmd.Flags().GetBool("global")
		if err != nil {
			printError(err.Error())
			return
		}

		elkFilePath, err := cmd.Flags().GetString("file")
		if err != nil {
			printError(err.Error())
			return
		}

		if len(args) == 0 {
			printError("A task name is required")
			return
		}

		elk, err := getElk(elkFilePath, isGlobal)

		if err != nil {
			printError(err.Error())
			return
		}

		logger := &engine.Logger{
			StdoutWriter: os.Stdout,
			StderrWriter: os.Stderr,
			StdinReader:  os.Stdin,
		}

		clientEngine := engine.New(elk, logger)

		var wg sync.WaitGroup

		for _, task := range args {
			wg.Add(1)
			go func(task string, wg *sync.WaitGroup) {
				defer wg.Done()

				if !clientEngine.HasTask(task) {
					printError(fmt.Sprintf("task '%s' do not exist", task))
					return
				}

				err = clientEngine.Run(task)
				if err != nil {
					printError(err.Error())
					return
				}
			}(task, &wg)
		}

		wg.Wait()
	},
}

func getElk(filePath string, isGlobal bool) (*engine.Elk, error) {
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
