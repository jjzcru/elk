package engine

import (
	"bufio"
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

	if len(task.DetachedDeps) > 0 {
		err = e.runTaskDependencies(taskName, true)
		if err != nil {
			return err
		}
	}

	if len(task.Deps) > 0 {
		err = e.runTaskDependencies(taskName, false)
		if err != nil {
			return err
		}
	}

	if len(task.Dir) == 0 {
		task.Dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	var envs []string
	// Load Env variables from system
	for _, env := range os.Environ() {
		envs = append(envs, env)
	}

	// Load Env variables from a file
	if len(e.elk.EnvFile) > 0 {
		envsInFile, err := getEnvFromFile(e.elk.EnvFile)
		if err != nil {
			return err
		}

		for _, env := range envsInFile {
			envs = append(envs, env)
		}
	}

	// Load Env variables from global vars in Elkfile
	for k, v := range e.elk.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	if len(task.EnvFile) > 0 {
		envsInFile, err := getEnvFromFile(task.EnvFile)
		if err != nil {
			return err
		}

		for _, env := range envsInFile {
			envs = append(envs, env)
		}
	}

	// Load Env variables from global vars in Elkfile
	for k, v := range task.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	for _, command := range task.Cmds {
		err = e.runCommand(task, envs, command)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) runTaskDependencies(taskName string, detached bool) error {
	task, err := e.GetTask(taskName)
	if err != nil {
		return err
	}

	if len(append(task.Deps, task.DetachedDeps...)) == 0 {
		return nil
	}

	visitedNodes := make(map[string]bool)

	err = e.HasCircularDependency(taskName, visitedNodes)
	if err != nil {
		return err
	}

	deps := task.Deps
	if detached {
		deps = task.DetachedDeps
	}

	for _, dep := range deps {
		if detached {
			go func() {
				_ = e.Run(dep)
			}()
		} else {
			err = e.Run(dep)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// HasCircularDependency checks if a task has a circular dependency
func (e *Engine) HasCircularDependency(taskName string, visitedNodes map[string]bool) error {
	task, err := e.GetTask(taskName)
	if err != nil {
		return err
	}

	if len(append(task.Deps, task.DetachedDeps...)) == 0 {
		return nil
	}

	dependencyGraph, err := e.getDependencyGraph(task)
	if err != nil {
		return err
	}

	_, hasVisitNode := visitedNodes[taskName]
	if hasVisitNode {
		return fmt.Errorf("The task '%s' has a circular dependency", taskName)
	}

	visitedNodes[taskName] = true

	for _, dep := range dependencyGraph {
		for _, d := range dep {
			err = e.HasCircularDependency(d, visitedNodes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getEnvFromFile(filePath string) ([]string, error) {
	var envs []string

	file, err := os.Open(filePath)
	if err != nil {
		return envs, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		envs = append(envs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return envs, err
	}

	return envs, nil
}

func (e *Engine) getDependencyGraph(task *Task) (map[string][]string, error) {
	dependencyGraph := make(map[string][]string)
	deps := append(task.Deps, task.DetachedDeps...)
	for _, dep := range deps {
		// Validate that the dependency is a valid task
		t, exists := e.elk.Tasks[dep]
		if exists == false {
			return dependencyGraph, fmt.Errorf("The dependency '%s' do not exist as a task", dep)
		}

		dependencyGraph[dep] = append(t.Deps, t.DetachedDeps...)
	}
	return dependencyGraph, nil
}

func (e *Engine) runCommand(task *Task, envs []string, command string) error {
	p, err := syntax.NewParser().Parse(strings.NewReader(command), "")
	if err != nil {
		return err
	}

	stdinReader := e.logger.StdinReader
	stdoutWriter := e.logger.StdoutWriter
	stderrWriter := e.logger.StderrWriter

	if len(task.Log) > 0 {
		logFile, err := os.OpenFile(task.Log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		stdoutWriter = logFile
		stderrWriter = logFile
	}

	r, err := interp.New(
		interp.Dir(task.Dir),
		interp.Env(expand.ListEnviron(envs...)),

		interp.Module(interp.DefaultExec),
		interp.Module(interp.OpenDevImpls(interp.DefaultOpen)),

		interp.StdIO(stdinReader, stdoutWriter, stderrWriter),
	)

	if err != nil {
		return err
	}

	ctx := context.Background()

	return r.Run(ctx, p)
}
