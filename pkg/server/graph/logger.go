package graph

import (
	"github.com/jjzcru/elk/pkg/engine"
	"github.com/jjzcru/elk/pkg/primitives/ox"
	"io"
	"os"
)

func GraphQLLogger(tasks map[string]ox.Task) (map[string]engine.Logger, chan map[string]string, chan map[string]string, error) {
	loggerMapper := make(map[string]engine.Logger)
	outChan := make(chan map[string]string)
	errChan := make(chan map[string]string)

	for name, task := range tasks {
		var stdOutWriter io.Writer = GraphQLWriter{
			task:   name,
			output: outChan,
		}

		var stdErrWriter io.Writer = GraphQLWriter{
			task:   name,
			output: errChan,
		}

		if len(task.Log.Out) > 0 {
			logFile, err := os.OpenFile(task.Log.Out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, nil, nil, err
			}

			stdOutWriter = io.MultiWriter(stdOutWriter, logFile)
		}

		if len(task.Log.Err) > 0 {
			logFile, err := os.OpenFile(task.Log.Err, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, nil, nil, err
			}

			stdErrWriter = io.MultiWriter(stdErrWriter, logFile)
		}

		loggerMapper[name] = engine.Logger{
			StdoutWriter: stdOutWriter,
			StderrWriter: stdErrWriter,
			StdinReader:  nil,
		}
	}

	return loggerMapper, outChan, errChan, nil
}

type GraphQLWriter struct {
	task   string
	output chan map[string]string
}

func (w GraphQLWriter) Write(p []byte) (int, error) {
	w.output <- map[string]string{w.task: string(p)}
	return len(p), nil
}
