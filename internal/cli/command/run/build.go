package run

import (
	"fmt"
	"github.com/jjzcru/elk/pkg/engine"
	"os"
	"path/filepath"
	"time"

	"github.com/jjzcru/elk/pkg/primitives/ox"
	"github.com/jjzcru/elk/pkg/utils"
	"github.com/spf13/cobra"
)

// Build loads ox object with the values
func Build(cmd *cobra.Command, e *ox.Elk, tasks []string) (map[string]engine.Logger, error) {
	logger := make(map[string]engine.Logger)

	ignoreLogFile, err := cmd.Flags().GetBool("ignore-log-file")
	if err != nil {
		return logger, err
	}

	ignoreError, err := cmd.Flags().GetBool("ignore-error")
	if err != nil {
		return logger, err
	}

	ignoreLogFormat, err := cmd.Flags().GetBool("ignore-log-format")
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

	taskMaps := make(map[string]bool)
	for _, task := range tasks {
		taskMaps[task] = true
	}

	for name, task := range e.Tasks {
		if _, ok := taskMaps[name]; !ok {
			continue
		}

		taskLogger := engine.DefaultLogger()

		if ignoreLogFile {
			logger[name] = taskLogger
			continue
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

		if len(task.Log.Format) > 0 && !ignoreLogFormat {
			format, err := getDateFormat(task.Log.Format)
			if err != nil {
				return nil, err
			}

			taskLogger, err = engine.TimeStampLogger(taskLogger, format)
			if err != nil {
				return nil, err
			}
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

func getDateFormat(format string) (string, error) {
	switch format {
	case "ANSIC":
		fallthrough
	case "ansic":
		return time.ANSIC, nil
	case "UnixDate":
		fallthrough
	case "unixdate":
		return time.UnixDate, nil
	case "rubydate":
		fallthrough
	case "RubyDate":
		return time.RubyDate, nil
	case "RFC822":
		return time.RFC822, nil
	case "RFC822Z":
		return time.RFC822Z, nil
	case "RFC850":
		return time.RFC850, nil
	case "RFC1123":
		return time.RFC1123, nil
	case "RFC1123Z":
		return time.RFC1123Z, nil
	case "RFC3339":
		return time.RFC3339, nil
	case "RFC3339Nano":
		return time.RFC3339Nano, nil
	case "kitchen":
		fallthrough
	case "Kitchen":
		return time.Kitchen, nil
	default:
		return "", fmt.Errorf("%s is an invalid timestamp format", format)
	}
}
