package config

import (
	"errors"
	"fmt"
	"github.com/jjzcru/elk/internal/cli/templates"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/spf13/cobra"
	"html/template"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

// NewSetCommand returns a cobra command that set elk file
func NewSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set global elk file",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a file path argument")
			}

			targetPath := args[0]
			if !utils.IsPathExist(targetPath) {
				return fmt.Errorf("provided file path do not exist: %s", targetPath)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			elkFilePath, err := filepath.Abs(args[0])
			if err != nil {
				utils.PrintError(err)
				return
			}

			_, err = GetElk(elkFilePath, false)
			if err != nil {
				utils.PrintError(err)
				return
			}

			response, err := template.New("config").Parse(templates.Config)
			if err != nil {
				utils.PrintError(err)
				return
			}

			usr, err := user.Current()
			if err != nil {
				utils.PrintError(err)
				return
			}

			configFile, err := os.Create(path.Join(usr.HomeDir, ".elk", "config.yml"))
			if err != nil {
				utils.PrintError(err)
				return
			}

			err = response.Execute(configFile, elkFilePath)
			if err != nil {
				utils.PrintError(err)
				return
			}

			fmt.Println(elkFilePath)
		},
	}

	return cmd
}
