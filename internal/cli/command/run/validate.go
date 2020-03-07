package run

import (
	"fmt"
	"github.com/jjzcru/elk/internal/cli/command/config"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/spf13/cobra"
)

func validate(cmd *cobra.Command, args []string) error {
	logFilePath, err := cmd.Flags().GetString("log")
	if err != nil {
		return err
	}

	if len(logFilePath) > 0 {
		isFile, err := utils.IsPathAFile(logFilePath)
		if err != nil {
			return err
		}

		if !isFile {
			return fmt.Errorf("path is not a file: %s", logFilePath)
		}
	}

	isWatch, err := cmd.Flags().GetBool("watch")
	if err != nil {
		return err
	}

	isGlobal, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}

	elkFilePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	// Check if the file path is set
	e, err := config.GetElk(elkFilePath, isGlobal)
	if err != nil {
		return err
	}

	for _, name := range args {
		task, err := e.GetTask(name)
		if err != nil {
			return err
		}

		if isWatch {
			if len(task.Watch) == 0 {
				return fmt.Errorf("task '%s' do now have a watch property", name)
			}
		}
	}

	return nil
}
