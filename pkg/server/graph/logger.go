package graph

import (
	"io"
	"os"
	"regexp"

	"github.com/jjzcru/elk/pkg/engine"
	"github.com/jjzcru/elk/pkg/primitives/ox"
)

func gqlLogger(elkTasks map[string]ox.Task, tasks []string) (map[string]engine.Logger, chan map[string]string, chan map[string]string, error) {
	var err error
	loggerMapper := make(map[string]engine.Logger)
	outChan := make(chan map[string]string)
	errChan := make(chan map[string]string)

	taskMaps := make(map[string]bool)
	for _, task := range tasks {
		taskMaps[task] = true
	}

	for name, task := range elkTasks {
		if _, ok := taskMaps[name]; !ok {
			continue
		}

		logger := engine.DefaultLogger()

		var stdOutWriter io.Writer = QLWriter{
			task:   name,
			output: outChan,
		}

		var stdErrWriter io.Writer = QLWriter{
			task:   name,
			output: errChan,
		}

		if len(task.Log.Out) > 0 {
			logger.StdoutWriter, err = os.OpenFile(task.Log.Out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, nil, nil, err
			}

			if len(task.Log.Format) > 0 {
				format, err := engine.TimeStampWriter{}.GetDateFormat(task.Log.Format)
				if err != nil {
					return nil, nil, nil, err
				}

				timeStampLogger, err := engine.TimeStampLogger(logger, format)
				if err != nil {
					return nil, nil, nil, err
				}

				logger.StdoutWriter = timeStampLogger.StdoutWriter
			}

			stdOutWriter = io.MultiWriter(stdOutWriter, logger.StdoutWriter)
		}

		if len(task.Log.Err) > 0 {
			logger.StderrWriter, err = os.OpenFile(task.Log.Err, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, nil, nil, err
			}

			if len(task.Log.Format) > 0 {
				format, err := engine.TimeStampWriter{}.GetDateFormat(task.Log.Format)
				if err != nil {
					return nil, nil, nil, err
				}

				timeStampLogger, err := engine.TimeStampLogger(logger, format)
				if err != nil {
					return nil, nil, nil, err
				}

				logger.StderrWriter = timeStampLogger.StderrWriter
			}

			stdErrWriter = io.MultiWriter(stdErrWriter, logger.StderrWriter)
		}

		loggerMapper[name] = engine.Logger{
			StdoutWriter: stdOutWriter,
			StderrWriter: stdErrWriter,
			StdinReader:  nil,
		}
	}

	return loggerMapper, outChan, errChan, nil
}

// QLWriter writes the logs from a task to an specific output
type QLWriter struct {
	task   string
	output chan map[string]string
}

func (w QLWriter) Write(p []byte) (int, error) {
	re := regexp.MustCompile(`\r?\n`)
	message := re.ReplaceAllString(string(p), " ")
	w.output <- map[string]string{w.task: message}
	return len(p), nil
}
