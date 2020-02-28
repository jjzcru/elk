package run

import (
	"fmt"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/primitives/elk"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func build(cmd *cobra.Command, e *elk.Elk) error {
	ignoreLog, err := cmd.Flags().GetBool("ignore-log")
	if err != nil {
		return err
	}

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

	if len(logFilePath) > 0 {
		_, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		absolutePath, err := filepath.Abs(logFilePath)
		if err != nil {
			return err
		}

		logFilePath = absolutePath
	}

	for name, task := range e.Tasks {
		if len(logFilePath) > 0 {
			task.Log = logFilePath
		}

		if ignoreLog {
			task.Log = ""
		}

		e.Tasks[name] = task
		err := e.HasCircularDependency(name)
		if err != nil {
			return err
		}
	}

	// TODO Should process env variables

	return nil
}
