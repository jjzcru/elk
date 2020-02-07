package commands

import (
	"bufio"
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
		createConfigFile()
	},
}

func createConfigFile() {
	// Get default hostname
	elkPath := getDefaultElkfilePath()

	scanner := bufio.NewScanner(os.Stdin)
	response, err := template.New("elk").Parse(installationTemplate)
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

	response, err = template.New("success_installation").Parse(elkTemplate)
	if err != nil {
		panic(err)
	}

	err = response.Execute(os.Stdout, elk)

	usr, _ := user.Current()
	configPath := path.Join(usr.HomeDir, ".elk", "config.yml")

	// Create the .elk directory if not exist
	elkInstallationPath := path.Join(usr.HomeDir, ".elk")
	if _, err := os.Stat(elkInstallationPath); os.IsNotExist(err) {
		_ = os.Mkdir(elkInstallationPath, 0777)
	}

	// Create config path
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		_, err = os.Create(configPath)
		if err != nil {
			panic(err)
		}
	}

	// Create elk file
	elkFile, _ := os.Create(elkPath)
	response, err = template.New("elk").Parse(elkTemplate)
	if err != nil {
		panic(err)
	}

	err = response.Execute(elkFile, elk)
	if err != nil {
		panic(err)
	}

	err = elkFile.Close()
	if err != nil {
		panic(err)
	}

	response, err = template.New("config").Parse(configTemplate)
	if err != nil {
		panic(err)
	}

	configFile, err := os.Create(configPath)
	if err != nil {
		panic(err)
	}

	err = response.Execute(configFile, elkPath)
	// err = ioutil.WriteFile(configPath, []byte(fmt.Sprintf(`path: "%s"`, elkPath)), 0777)
	if err != nil {
		panic(err)
	}

	if scanner.Err() != nil {
		// handle error.
	}
}

func getDefaultElkfilePath() string {
	usr, _ := user.Current()
	return path.Join(usr.HomeDir, ".elk", "Elkfile.yml")
}

var elkTemplate = `version: '1'
tasks: {{ range $name, $task := .Tasks }}
  {{ $name }}:
    cmds: {{range $cmd := $task.Cmds}}
      - {{$cmd}}{{end}}
    description: '{{$task.Description}}'
{{end}}
`

var installationTemplate = `
This will create a default Elkfile

It only covers just a few details. 

The installation will include some default events like 'shutdown' 
or 'restart' just to get started but you will be able to add more 
events in the configuration file.

`

var configTemplate = `# DO NOT modify this file directly
# for alter this file
paths: 
  - "{{.}}"
`
