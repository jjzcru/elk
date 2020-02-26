package init

import (
	"github.com/jjzcru/elk/internal/cli/templates"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/primitives/elk"
	"os"
	"path"
	"text/template"

	"github.com/spf13/cobra"
)

// NewInitializeCommand returns a cobra command for `init` sub command
func NewInitializeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a elk in the current directory",
		Run: func(cmd *cobra.Command, args []string) {
			elkFilePath, err := getElkfilePath()
			if err != nil {
				utils.PrintError(err)
				return
			}

			err = CreateElkFile(elkFilePath)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	return cmd
}

// CreateElkFile create an elk file in path
func CreateElkFile(elkFilePath string) error {
	response, err := template.New("installation").Parse(templates.Installation)
	if err != nil {
		return err
	}

	err = response.Execute(os.Stdout, "")
	if err != nil {
		return err
	}

	elkFile, _ := os.Create(elkFilePath)
	defer elkFile.Close()

	response, err = template.New("elk").Parse(templates.Elk)
	if err != nil {
		return err
	}

	err = response.Execute(elkFile, elk.Elk{
		Version: "1",
		Tasks: map[string]elk.Task{
			"shutdown": {
				Cmds: []string{
					"shutdown",
				},
				Description: "Command to shutdown the machine",
			},
			"restart": {
				Cmds: []string{
					"reboot",
				},
				Description: "Command that should restart the machine",
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func getElkfilePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(dir, "elk.yml"), nil
}
