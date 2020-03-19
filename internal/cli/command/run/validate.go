package run

import (
	"fmt"

	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/primitives/elk"
	"github.com/spf13/cobra"
)

// Validate if the arguments are valid
func Validate(cmd *cobra.Command, tasks []string) error {
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
		isWatch = false
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
	e, err := utils.GetElk(elkFilePath, isGlobal)
	if err != nil {
		return err
	}

	for _, name := range tasks {
		task, err := e.GetTask(name)
		if err != nil {
			if err == elk.ErrTaskNotFound {
				return fmt.Errorf("task \"%s\" not found", name)
			}
			return err
		}

		if isWatch {
			if len(task.Sources) == 0 {
				return fmt.Errorf("task '%s' do now have a watch property", name)
			}
		}
	}

	return nil
}
