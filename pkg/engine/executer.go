package engine

import (
	"context"
	"fmt"
	"github.com/jjzcru/elk/pkg/primitives/ox"
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
	Logger *Logger
}

// Execute task and returns a PID
func (e DefaultExecuter) Execute(ctx context.Context, elk *ox.Elk, name string) (int, error) {
	ctx, _ = context.WithCancel(ctx)
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
			depCtx, _ := context.WithCancel(ctx)
			go e.Execute(depCtx, elk, dep)
		}
	}

	if len(deps) > 0 {
		for _, dep := range deps {
			depCtx, _ := context.WithCancel(ctx)
			_, err := e.Execute(depCtx, elk, dep)
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

	stdinReader := e.Logger.StdinReader
	stdoutWriter := e.Logger.StdoutWriter
	stderrWriter := e.Logger.StderrWriter

	if len(task.Log) > 0 {
		logFile, err := os.OpenFile(task.Log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}

		stdoutWriter = logFile
		stderrWriter = logFile
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

func getEnvs(envMap map[string]string) []string {
	var envs []string
	for env, value := range envMap {
		envs = append(envs, fmt.Sprintf("%s=%s", env, value))
	}
	return envs
}
