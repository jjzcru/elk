package engine

import (
	"context"
	"github.com/jjzcru/elk/pkg/primitives/elk"
	"testing"
)

func TestDefaultExecuterExecute(t *testing.T) {
	e := elk.Elk{
		Version: "1",
		Tasks: map[string]elk.Task{
			"world": {
				Deps: []elk.Dep{
					{
						Name: "hello",
					},
					{
						Name:     "foo",
						Detached: true,
					},
				},
				Env: map[string]string{
					"FOO": "BAR",
				},
				Cmds: []string{
					"echo $FOO",
				},
			},
			"hello": {
				Description: "Empty Task",
				Env: map[string]string{
					"FOO": "Bar",
				},
				Cmds: []string{
					"echo $FOO",
				},
			},
			"foo": {
				Description: "Empty Task",
				Env: map[string]string{
					"FOO": "Bar",
				},
				Cmds: []string{
					"echo $FOO",
				},
			},
		},
	}

	executer := DefaultExecuter{
		Logger: &DefaultLogger,
	}

	_, err := executer.Execute(context.Background(), &e, "world")
	if err != nil {
		t.Error(err)
	}

}

func TestDefaultExecuterExecuteTaskNotExist(t *testing.T) {
	e := elk.Elk{
		Version: "1",
		Tasks: map[string]elk.Task{
			"world": {
				Deps: []elk.Dep{
					{
						Name: "hello",
					},
					{
						Name:     "foo",
						Detached: true,
					},
				},
				Env: map[string]string{
					"FOO": "BAR",
				},
				Cmds: []string{
					"echo $FOO",
				},
			},
			"hello": {
				Description: "Empty Task",
				Env: map[string]string{
					"FOO": "Bar",
				},
				Cmds: []string{
					"echo $FOO",
				},
			},
			"foo": {
				Description: "Empty Task",
				Env: map[string]string{
					"FOO": "Bar",
				},
				Cmds: []string{
					"echo $FOO",
				},
			},
		},
	}

	executer := DefaultExecuter{
		Logger: &DefaultLogger,
	}

	_, err := executer.Execute(context.Background(), &e, "bar")
	if err == nil {
		t.Error("task do not exist should throw an error")
	}

}
