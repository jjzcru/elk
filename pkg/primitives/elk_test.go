package model

import (
	"testing"
)

func TestHasCircularDependency(t *testing.T) {
	elk := Elk{
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

	for taskName := range elk.Tasks {
		err := elk.HasCircularDependency(taskName)
		if err != nil {
			t.Error(err.Error())
		}
	}
}
