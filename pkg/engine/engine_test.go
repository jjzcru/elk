package engine

import (
	"os"
	"testing"
)

func getTestEngine() *Engine {
	elk := &Elk{
		Version: "1",
		Tasks: map[string]Task{
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

	for taskName := range engine.elk.Tasks {
		err := engine.Run(taskName)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

/*func TestGetEnvFromFile(t *testing.T) {
	filePath := "/tmp/example"
	_, err := getEnvFromFile(filePath)
	if err != nil {
		t.Error((err.Error()))
	}
}*/
