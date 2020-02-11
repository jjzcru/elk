package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/jjzcru/elk/pkg/engine"
	"gopkg.in/yaml.v2"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run task declared in the 'Elkfile.yml'",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			printError("A task name is required")
			return
		}

		/*elkFilePath, err := getElkFilePath()
		if err != nil {
			printError(err.Error())
			return
		}
		fmt.Printf("Elkfile path: '%s'\n", elkFilePath)*/

		task := args[0]

		elk, err := getElk()

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

		if !clientEngine.HasTask(task) {
			printError(fmt.Sprintf("event '%s' should exist", task))
			return
		}

		err = clientEngine.Run(task)
		if err != nil {
			printError(err.Error())
			return
		}
	},
}

func printError(err string) {
	fmt.Print(aurora.Bold(aurora.Red("ERROR: ")))

	fmt.Fprintf(os.Stderr, err)
}

func getElk() (*engine.Elk, error) {
	usr, _ := user.Current()
	elkConfigPath := path.Join(usr.HomeDir, ".elk", "Elkfile.yml")
	elk := engine.Elk{}
	if _, err := os.Stat(elkConfigPath); os.IsNotExist(err) {
		return nil, errors.New("Elk is not installed, please run \n elk init")
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

/*func getConfig() (*engine.Config, error) {
	usr, _ := user.Current()
	elkConfigPath := path.Join(usr.HomeDir, ".elk", "Elkfile.yml")
	config := engine.Config{}
	if _, err := os.Stat(elkConfigPath); os.IsNotExist(err) {
		return nil, errors.New("Elk is not installed, please run \n elk install")
	}

	data, err := ioutil.ReadFile(elkConfigPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}*/
