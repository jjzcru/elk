package init

import (
	"os"
	"path"
	"runtime"
	"text/template"

	"github.com/jjzcru/elk/internal/cli/templates"
	"github.com/jjzcru/elk/pkg/primitives/ox"
	"github.com/jjzcru/elk/pkg/utils"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

// NewInitializeCommand returns a cobra command for `init` sub command
func NewInitializeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates an ox.yml file in the current directory",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			elkFilePath, err := getElkfilePath()
			if err != nil {
				utils.PrintError(err)
				return
			}

			_, err = os.Stat(elkFilePath)
			if os.IsNotExist(err) {
				err = CreateElkFile(elkFilePath)
				if err != nil {
					utils.PrintError(err)
				}
			}
		},
	}

	return cmd
}

// CreateElkFile create an ox file in path
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

	response, err = template.New("ox").Parse(templates.Elk)
	if err != nil {
		return err
	}

	restart := "reboot"
	shutdown := "shutdown"

	if runtime.GOOS == "windows" {
		restart = "shutdown /r"
		shutdown = "shutdown /s"
	}

	e := ox.Elk{
		Version: "1",
		Env: map[string]string{
			"HELLO": "World",
		},
		Tasks: map[string]ox.Task{
			"hello": {
				Description: "Print hello world",
				Env: map[string]string{
					"HELLO": "Hello",
				},
				Cmds: []string{
					"echo $HELLO",
				},
			},
			"test-log": {
				Description: "Print World",
				Log: ox.Log{
					Out: "./test.log",
				},
				Cmds: []string{
					"echo $HELLO",
				},
			},
			"ts-run": {
				Description: "Run a typescript app",
				Cmds: []string{
					"npm start",
				},
				Deps: []ox.Dep{
					{
						Name:     "ts-build",
						Detached: false,
					},
				},
			},
			"ts-build": {
				Description: "Watch files and re-run to compile typescript",
				Sources:     "[a-zA-Z]*.ts$",
				Cmds: []string{
					"npm run build",
				},
			},
			"shutdown": {
				Description: "Command to shutdown the machine",
				Cmds: []string{
					shutdown,
				},
			},
			"restart": {
				Description: "Command that should restart the machine",
				Cmds: []string{
					restart,
				},
			},
		},
	}

	b, err := yaml.Marshal(e)
	if err != nil {
		return err
	}

	_, err = elkFile.Write(b)

	//err = response.Execute(elkFile, e)

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

	return path.Join(dir, "ox.yml"), nil
}
