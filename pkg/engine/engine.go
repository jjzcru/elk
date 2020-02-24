package engine

import (
	"fmt"
	"strings"

	"github.com/jjzcru/elk/pkg/primitives"
)

// Engine is the data structure responsible of processing the content
type Engine struct {
	Elk      *primitives.Elk
	Executer Executer
	Build func(*primitives.Elk) error
}

// New creates a new instance of the engine
func New(elk *primitives.Elk, executer Executer) *Engine {
	return &Engine{
		Elk:      elk,
		Executer: executer,
	}
}

// Run task declared in Elkfile
func (e *Engine) Run(task string) error {
	err := e.Build(e.Elk)
	if err != nil {
		return err
	}

	if !e.Elk.HasTask(task) {
		return fmt.Errorf("task '%s' not found", task)
	}

	err = e.Elk.LoadEnvFile()

	if err != nil {
		return err
	}

	e.Elk.LoadEnvsInTasks()

	_, err = e.Executer.Execute(e.Elk, task)
	if err != nil {
		return err
	}

	return nil
}

// MapEnvs map an array of string env
func MapEnvs(envs []string) map[string]string {
	envMap := make(map[string]string)
	for _, env := range envs {
		result := strings.SplitAfterN(env, "=", 2)
		env := strings.ReplaceAll(result[0], "=", "")
		value := result[1]
		envMap[env] = value
	}

	return envMap
}
