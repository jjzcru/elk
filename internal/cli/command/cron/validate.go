package cron

import (
	"fmt"
	"github.com/jjzcru/elk/internal/cli/command/config"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/spf13/cobra"
)

func validate(cmd *cobra.Command) error {
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

	isGlobal, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}

	elkFilePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	// Check if the file path is set
	_, err = config.GetElk(elkFilePath, isGlobal)
	if err != nil {
		return err
	}

	return nil
}
