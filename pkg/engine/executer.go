package engine

import (
	"context"
	"fmt"
	"github.com/jjzcru/elk/pkg/primitives/ox"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/expand"
	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)

// Executer runs a task and returns a PID and an error
type Executer interface {
	Execute(context.Context, *ox.Elk, string) (int, error)
}

// DefaultExecuter Execute task with a POSIX emulator
type DefaultExecuter struct {
	Logger map[string]Logger
}

// Execute task and returns a PID
func (e DefaultExecuter) Execute(ctx context.Context, elk *ox.Elk, name string) (int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pid := os.Getpid()

	task, err := elk.GetTask(name)
	if err != nil {
		return pid, err
	}

	var detachedDeps []string
	var deps []string

	for _, dep := range task.Deps {
		if dep.Detached {
			detachedDeps = append(detachedDeps, dep.Name)
		} else {
			deps = append(deps, dep.Name)
		}
	}

	if len(detachedDeps) > 0 {
		for _, dep := range detachedDeps {
			go e.ExecuteDetached(ctx, elk, dep)
		}
	}

	if len(deps) > 0 {
		for _, dep := range deps {
			_, err := e.Execute(ctx, elk, dep)
			if err != nil {
				return pid, err
			}
		}
	}

	if len(task.Dir) == 0 {
		task.Dir, err = os.Getwd()
		if err != nil {
			return pid, err
		}
	}

	var stdinReader io.Reader
	var stdoutWriter io.Writer
	var stderrWriter io.Writer

	logger, exists := e.Logger[name]

	if !exists {
		stdinReader = os.Stdin
		stdoutWriter = os.Stdout
		stderrWriter = os.Stderr
	} else {
		stdinReader = logger.StdinReader
		stdoutWriter = logger.StdoutWriter
		stderrWriter = logger.StderrWriter
	}

	for _, command := range task.Cmds {
		command, err := ox.GetCmdFromVars(task.Vars, command)
		if err != nil {
			return pid, err
		}

		p, err := syntax.NewParser().Parse(strings.NewReader(command), "")
		if err != nil {
			return pid, err
		}

		envs := getEnvs(task.Env)

		r, err := interp.New(
			interp.Dir(task.Dir),

			interp.Env(expand.ListEnviron(envs...)),

			interp.Module(interp.DefaultExec),
			interp.Module(interp.OpenDevImpls(interp.DefaultOpen)),

			interp.StdIO(stdinReader, stdoutWriter, stderrWriter),
		)

		if err != nil {
			return pid, err
		}
		err = r.Run(ctx, p)
		if err != nil && !task.IgnoreError {
			return pid, err
		}
	}
	return pid, nil
}

// ExecuteDetached do not keep track of the execution of the task
func (e DefaultExecuter) ExecuteDetached(ctx context.Context, elk *ox.Elk, name string) {
	_, _ = e.Execute(ctx, elk, name)
}

func getEnvs(envMap map[string]string) []string {
	var envs []string
	for env, value := range envMap {
		envs = append(envs, fmt.Sprintf("%s=%s", env, value))
	}
	return envs
}
