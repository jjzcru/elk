package graph

import "github.com/jjzcru/elk/pkg/engine"

func GraphQLLogger() (engine.Logger, chan string, chan string) {
	outChan := make(chan string)
	errChan := make(chan string)

	logger := engine.Logger{
		StdoutWriter: GraphQLWriter{
			output: outChan,
		},
		StderrWriter: GraphQLWriter{
			output: errChan,
		},
		StdinReader: nil,
	}

	return logger, outChan, errChan
}

type GraphQLWriter struct {
	output chan string
}

func (w GraphQLWriter) Write(p []byte) (int, error) {
	w.output <- string(p)
	return len(p), nil
}
