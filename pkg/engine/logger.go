package engine

import "io"

// Logger is used by the engine to store the output
type Logger struct {
	StderrWriter io.Writer
	StdoutWriter io.Writer
	StdinReader  io.Reader
}
