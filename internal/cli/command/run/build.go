package run

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jjzcru/elk/pkg/engine"

	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/primitives/ox"
	"github.com/spf13/cobra"
)

// Build loads ox object with the values
func Build(cmd *cobra.Command, e *ox.Elk) (map[string]engine.Logger, error) {
	logger := make(map[string]engine.Logger)

	ignoreLog, err := cmd.Flags().GetBool("ignore-log-file")
	if err != nil {
		return logger, err
	}

	ignoreError, err := cmd.Flags().GetBool("ignore-error")
	if err != nil {
		return logger, err
	}

	ignoreDep, err := cmd.Flags().GetBool("ignore-deps")
	if err != nil {
		ignoreDep = false
	}

	logFilePath, err := cmd.Flags().GetString("log")
	if err != nil {
		return logger, err
	}

	if len(logFilePath) > 0 {
		isFile, err := utils.IsPathAFile(logFilePath)
		if err != nil {
			return logger, err
		}

		if !isFile {
			return logger, fmt.Errorf("path is not a file: %s", logFilePath)
		}
	}

	if e.Env == nil {
		e.Env = make(map[string]string)
	}

	if len(logFilePath) > 0 {
		_, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return logger, err
		}

		absolutePath, err := filepath.Abs(logFilePath)
		if err != nil {
			return logger, err
		}

		logFilePath = absolutePath
	}

	for name, task := range e.Tasks {
		if len(logFilePath) > 0 {
			task.Log = logFilePath
		}

		if len(task.Log) > 0 {
			logFile, err := os.OpenFile(task.Log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return logger, err
			}

			logger[name] = engine.Logger{
				StderrWriter: logFile,
				StdoutWriter: logFile,
				StdinReader:  os.Stdin,
			}
		}

		if ignoreLog {
			task.Log = ""
		}

		if ignoreError {
			task.IgnoreError = true
		}

		if ignoreDep {
			task.Deps = []ox.Dep{}
		}

		e.Tasks[name] = task
	}

	err = e.Build()
	if err != nil {
		return logger, err
	}

	return logger, nil
}
