package engine

import (
	"fmt"
	"testing"

	"github.com/jjzcru/elk/pkg/primitives"
)

func getTestEngine() *Engine {
	elk := &primitives.Elk{
		Version: "1",
		Tasks: map[string]primitives.Task{
			"hello": {
				Description: "Empty Task",
				Cmds: []string{
					"clear",
				},
			},
			"world": {
				Deps: []string{
					"hello",
				},
				Env: map[string]string{
					"FOO": "BAR",
				},
				Cmds: []string{
					"clear",
				},
			},
		},
	}

	return &Engine{
		Elk: elk,
		Executer: DefaultExecuter{
			Logger: &DefaultLogger,
		},
	}
}

func TestRun(t *testing.T) {
	engine := getTestEngine()

	for taskName := range engine.Elk.Tasks {
		err := engine.Run(taskName)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestMapEnvs(t *testing.T) {
	key := "FOO"
	value := "http://localhost:7777?id=20"

	envs := []string{fmt.Sprintf("%s=%s", key, value)}
	envMap := MapEnvs(envs)
	if envMap[key] != value {
		t.Errorf("The key '%s' should have a value of '%s' but have a value of '%s' instead", key, value, envMap[key])
	}
}
