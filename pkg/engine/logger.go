package engine

import "io"

type Logger struct {
	StderrWriter io.Writer
	StdoutWriter io.Writer
	StdinReader  io.Reader
}
