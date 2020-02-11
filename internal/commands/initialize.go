package commands

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"text/template"

	"github.com/jjzcru/elk/pkg/engine"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create elk configuration file",
	Long:  fmt.Sprintf(`Create the configuration file where the daemon is going to read the properties (default: %s)`, getDefaultElkfilePath()),
	Run: func(cmd *cobra.Command, args []string) {
		err := createElkFile(getDefaultElkfilePath())
		if err != nil {
			_ = fmt.Errorf(err.Error())
		}
	},
}

func createElkFile(elkFilePath string) error {
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

func getDefaultElkfilePath() string {
	usr, _ := user.Current()
	return path.Join(usr.HomeDir, ".elk", "Elkfile.yml")
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
