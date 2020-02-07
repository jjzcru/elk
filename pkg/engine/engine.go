package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"mvdan.cc/sh/expand"
	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)

// Engine is the data structure responsible of processing the content
type Engine struct {
	logger *Logger
	elk    *Elk
}

// New creates a new instance of the engine
func New(elk *Elk, logger *Logger) *Engine {
	return &Engine{
		elk:    elk,
		logger: logger,
	}
}

// Set overwrite elk object data
func (e *Engine) Set(elk *Elk) {
	e.elk = elk
}

// HasTask return a boolean if the incoming event exist
func (e *Engine) HasTask(taskName string) bool {
	for task := range e.elk.Tasks {
		if task == taskName {
			return true
		}
	}
	return false
}

// GetTask return the Task for the incoming event
func (e *Engine) GetTask(taskName string) (*Task, error) {
	for task, event := range e.elk.Tasks {
		if task == taskName {
			return &event, nil
		}
	}
	return nil, errors.New("task not found")
}

// Run task declared in Elkfile
func (e *Engine) Run(taskName string) error {
	if !e.HasTask(taskName) {
		return errors.New("task not found")
	}

	task, err := e.GetTask(taskName)
	if err != nil {
		return err
	}

	if len(task.Dir) == 0 {
		task.Dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	var envs []string
	for k, v := range task.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	for _, env := range os.Environ() {
		envs = append(envs, env)
	}

	for _, command := range task.Cmds {
		err = e.runCommand(task, envs, command)
	}

	return nil
}

func (e *Engine) runCommand(task *Task, envs []string, command string) error {
	p, err := syntax.NewParser().Parse(strings.NewReader(command), "")
	if err != nil {
		return err
	}

	r, err := interp.New(
		interp.Dir(task.Dir),
		interp.Env(expand.ListEnviron(envs...)),

		interp.Module(interp.DefaultExec),
		interp.Module(interp.OpenDevImpls(interp.DefaultOpen)),

		interp.StdIO(e.logger.StdinReader, e.logger.StdoutWriter, e.logger.StderrWriter),
	)

	if err != nil {
		return err
	}

	ctx := context.Background()

	return r.Run(ctx, p)
}
