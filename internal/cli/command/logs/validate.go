package logs

import (
	"fmt"
	"os"

	"github.com/jjzcru/elk/internal/cli/command/config"
	"github.com/jjzcru/elk/pkg/primitives/elk"
	"github.com/spf13/cobra"
)

func validate(cmd *cobra.Command, args []string) error {
	e, err := getElk(cmd)
	if err != nil {
		return err
	}

	for _, name := range args {
		task, err := e.GetTask(name)
		if err != nil {
			return err
		}

		if len(task.Log) == 0 {
			return fmt.Errorf("task '%s' do not have a log file", name)
		}

		info, err := os.Stat(task.Log)
		if os.IsNotExist(err) {
			return err
		}

		if info.IsDir() {
			return fmt.Errorf("log path '%s' is a directory", task.EnvFile)
		}
	}

	return nil
}

func getElk(cmd *cobra.Command) (*elk.Elk, error) {
	isGlobal, err := cmd.Flags().GetBool("global")
	if err != nil {
		return nil, err
	}

	elkFilePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil, err
	}

	// Check if the file path is set
	e, err := config.GetElk(elkFilePath, isGlobal)
	if err != nil {
		return nil, err
	}

	return e, nil
}
