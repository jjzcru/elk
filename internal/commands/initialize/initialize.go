package init

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/jjzcru/elk/pkg/engine"
	"github.com/spf13/cobra"
)

// Cmd Command that initialize elk in current directory
func Cmd() *cobra.Command {
	var command = &cobra.Command{
		Use:   "init",
		Short: "Create a 'elk.yml' file in current directory",
		Run: func(cmd *cobra.Command, args []string) {
			elkFilePath, err := getElkfilePath()
			if err != nil {
				_ = fmt.Errorf(err.Error())
				return
			}

			err = CreateElkFile(elkFilePath)
			if err != nil {
				_ = fmt.Errorf(err.Error())
			}
		},
	}

	return command
}

// CreateElkFile create an elk file in path
func CreateElkFile(elkFilePath string) error {
	response, err := template.New("installation").Parse(installationTemplate)
	if err != nil {
		panic(err)
	}

	err = response.Execute(os.Stdout, "")

	elk := engine.Elk{
		Version: "1",
		Tasks: map[string]engine.Task{
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
	}

	// Create elk file
	elkFile, _ := os.Create(elkFilePath)
	response, err = template.New("elk").Parse(elkTemplate)
	if err != nil {
		return err
	}

	err = response.Execute(elkFile, elk)
	if err != nil {
		return err
	}

	err = elkFile.Close()
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

var elkTemplate = `version: '1'
tasks: {{ range $name, $task := .Tasks }}
  {{ $name }}:
    description: '{{$task.Description}}'
    cmds: {{range $cmd := $task.Cmds}}
      - {{$cmd}}{{end}}
{{end}}
`

var installationTemplate = `
This will create a default Elkfile

It only covers just a few tasks. 

The installation will include some default events like 'shutdown' 
or 'restart' just to get started but you will be able to add more 
events in the configuration file.

`
