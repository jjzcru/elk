package engine

import (
	"io"
	"os"
)

// Logger is used by the engine to store the output
type Logger struct {
	StderrWriter io.Writer
	StdoutWriter io.Writer
	StdinReader  io.Reader
}

// DefaultLogger is the standard output for a logger
var DefaultLogger = Logger{
	StdoutWriter: os.Stdout,
	StderrWriter: os.Stderr,
	StdinReader:  os.Stdin,
}