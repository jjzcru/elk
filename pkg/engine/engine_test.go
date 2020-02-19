package engine

import (
	"os"
	"testing"
)

func getTestEngine() *Engine {
	elk := &Elk{
		Version: "1",
		Tasks: map[string]Task{
			"hello": Task{
				Description: "Empty Task",
				Cmds: []string{
					"echo Hello",
				},
			},
			"world": Task{
				Deps: []string{
					"hello",
					"world",
				},
				Env: map[string]string{
					"FOO": "BAR",
				},
				Cmds: []string{
					"echo World $FOO",
				},
			},
		},
	}

	logger := &Logger{
		StdoutWriter: os.Stdout,
		StderrWriter: os.Stderr,
		StdinReader:  os.Stdin,
	}

	return &Engine{
		elk:    elk,
		logger: logger,
	}
}

func TestHasCircularDependency(t *testing.T) {
	engine := getTestEngine()

	visitedNodes := make(map[string]bool)

	for taskName := range engine.elk.Tasks {
		err := engine.HasCircularDependency(taskName, visitedNodes)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestRun(t *testing.T) {
	engine := getTestEngine()

	err := engine.Run("world")
	if err != nil {
		t.Error(err.Error())
	}
}
