package graph

import "github.com/jjzcru/elk/pkg/engine"

func GraphQLLogger(tasks []string) (map[string]engine.Logger, chan map[string]string, chan map[string]string) {
	loggerMapper := make(map[string]engine.Logger)
	outChan := make(chan map[string]string)
	errChan := make(chan map[string]string)

	for _, task := range tasks {
		loggerMapper[task] = engine.Logger{
			StdoutWriter: GraphQLWriter{
				task:   task,
				output: outChan,
			},
			StderrWriter: GraphQLWriter{
				task:   task,
				output: errChan,
			},
			StdinReader: nil,
		}
	}

	return loggerMapper, outChan, errChan
}

type GraphQLWriter struct {
	task   string
	output chan map[string]string
}

func (w GraphQLWriter) Write(p []byte) (int, error) {
	w.output <- map[string]string{w.task: string(p)}
	return len(p), nil
}
