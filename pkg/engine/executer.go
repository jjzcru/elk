package engine

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jjzcru/elk/pkg/primitives"
	"mvdan.cc/sh/expand"
	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)

// Executer runs a task and returns a PID and an error
type Executer interface {
	Execute(context.Context, *primitives.Elk, string) (int, error)
}

// DefaultExecuter Execute task with a POSIX emulator
type DefaultExecuter struct {
	Logger *Logger
}

// Execute task and returns a PID
func (e DefaultExecuter) Execute(ctx context.Context, elk *primitives.Elk, name string) (int, error) {
	pid := os.Getpid()

	task, err := elk.GetTask(name)
	if err != nil {
		return pid, err
	}

	if len(task.Deps) > 0 {
		for _, dep := range task.Deps {
			_, err := e.Execute(ctx, elk, dep)
			if err != nil {
				return pid, err
			}
		}
	}

	if len(task.DetachedDeps) > 0 {
		for _, dep := range task.DetachedDeps {
			go func() {
				_, _ = e.Execute(ctx, elk, dep)
			}()
		}
	}

	if len(task.Dir) == 0 {
		task.Dir, err = os.Getwd()
		if err != nil {
			return pid, err
		}
	}

	err = task.LoadEnvFile()
	if err != nil {
		return 0, err
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
		p, err := syntax.NewParser().Parse(strings.NewReader(command), "")
		if err != nil {
			return pid, err
		}

		envs := uniqueEnvs(append(os.Environ(), getEnvs(task.Env)...))

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

		ctx, _:=context.WithCancel(ctx)
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

func uniqueEnvs(envs []string) []string {
	var response []string

	for k, v := range MapEnvs(envs) {
		response = append(response, fmt.Sprintf("%s=%s", k, v))
	}

	return envs
}
