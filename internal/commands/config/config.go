package config

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/jjzcru/elk/pkg/engine"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// ConfigTemplate Template for the config.yml
var ConfigTemplate = `# DO NOT modify this file directly
# for alter this file
path: "{{.}}"
`

// Cmd Command that works with configuration
func Cmd() *cobra.Command {
	var command = &cobra.Command{
		Use:   "config",
		Short: "Work with the global installation",
	}

	command.AddCommand(GetCmd)
	command.AddCommand(SetCmd)

	return command
}

// SetCmd Command that set the global file
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the path for the global configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			PrintError("A file path is required")
			return
		}
		targetPath := args[0]

		if !fileExists(targetPath) {
			PrintError(fmt.Sprintf("The file path \"%s\" do not exist", targetPath))
			return
		}

		elkFilePath, err := filepath.Abs(targetPath)
		if err != nil {
			PrintError(err.Error())
			return
		}

		_, err = GetElk(elkFilePath, false)
		if err != nil {
			PrintError(err.Error())
			return
		}

		response, err := template.New("config").Parse(ConfigTemplate)
		if err != nil {
			PrintError(err.Error())
			return
		}

		usr, err := user.Current()
		if err != nil {
			PrintError(err.Error())
			return
		}

		configPath := path.Join(usr.HomeDir, ".elk", "config.yml")
		if err != nil {
			PrintError(err.Error())
			return
		}

		configFile, err := os.Create(configPath)
		if err != nil {
			PrintError(err.Error())
			return
		}

		err = response.Execute(configFile, elkFilePath)
		if err != nil {
			PrintError(err.Error())
			return
		}

		fmt.Println(elkFilePath)
	},
}

// GetCmd Command that get the global file
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the path for the global configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := getConfigPath()
		if err != nil {
			PrintError(err.Error())
			return
		}

		fmt.Println(configPath)
	},
}

func getConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	configPath := path.Join(usr.HomeDir, ".elk", "config.yml")

	if !fileExists(configPath) {
		return "", fmt.Errorf("The installation path \"%s\" do not exist. \nPlease run \"elk install\" to create it", configPath)
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetElk get an elk pointer from a file path
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
		return "", errors.New("Elk is not installed, \nPlease run \"elk install\"")
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
