package config

import (
	"fmt"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/primitives"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
	"path"
)

// NewGetCommand returns a cobra command that get elk file
func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get global elk file path",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			configPath, err := getConfigPath()
			if err != nil {
				utils.PrintError(err)
				return
			}

			fmt.Println(configPath)
		},
	}

	return cmd
}

func getConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	configPath := path.Join(usr.HomeDir, ".elk", "config.yml")
	if !utils.IsPathExist(configPath) {
		return "", fmt.Errorf("installation path do not exist. %s\nPlease run \"elk install\" to create it", configPath)
	}

	config := primitives.Config{}

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
