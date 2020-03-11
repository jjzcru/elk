package engine

import (
	"context"
	"fmt"
	elk2 "github.com/jjzcru/elk/pkg/primitives/elk"
	"testing"
)

func getTestEngine() *Engine {
	elk := &elk2.Elk{
		Version: "1",
		Tasks: map[string]elk2.Task{
			"hello": {
				Description: "Empty Task",
				Cmds: []string{
					"echo Hello",
				},
			},
			"world": {
				Deps: []elk2.Dep{
					{
						Name: "hello",
					},
				},
				Env: map[string]string{
					"FOO": "BAR",
				},
				Cmds: []string{
					"echo World",
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

	ctx := context.Background()

	for taskName := range engine.Elk.Tasks {
		ctx, cancel := context.WithCancel(ctx)
		err := engine.Run(ctx, taskName)
		if err != nil {
			t.Error(err.Error())
		}
		cancel()
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
