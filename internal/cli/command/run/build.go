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

	ignoreLogFile, err := cmd.Flags().GetBool("ignore-log-file")
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
		isFile, _ := utils.IsPathAFile(logFilePath)

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
		 taskLogger := engine.Logger{
			StdinReader:  os.Stdin,
			StdoutWriter: os.Stdout,
			StderrWriter: os.Stderr,
		}

		if len(logFilePath) > 0 {
			task.Log = ox.Log{
				Out: logFilePath,
				Err: logFilePath,
			}
		}

		if len(task.Log.Out) > 0 {
			logFile, err := os.OpenFile(task.Log.Out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return logger, err
			}

			taskLogger.StdoutWriter = logFile
		}

		if len(task.Log.Err) > 0 {
			logFile, err := os.OpenFile(task.Log.Err, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return logger, err
			}

			taskLogger.StderrWriter = logFile
		} else {
			taskLogger.StderrWriter = taskLogger.StdoutWriter
		}

		if ignoreLogFile {
			taskLogger.StdoutWriter = os.Stdout
			taskLogger.StderrWriter = os.Stderr
		}

		logger[name] = taskLogger

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
