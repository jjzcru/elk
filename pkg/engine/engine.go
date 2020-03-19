package engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/jjzcru/elk/pkg/primitives/elk"
)

// Engine is the data structure responsible of processing the content
type Engine struct {
	Elk      *elk.Elk
	Executer Executer
	Build    func() error
}

// Run task declared in ox.yml file
func (e *Engine) Run(ctx context.Context, task string) error {
	if !e.Elk.HasTask(task) {
		return fmt.Errorf("task '%s' not found", task)
	}

	_, err := e.Executer.Execute(ctx, e.Elk, task)
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
