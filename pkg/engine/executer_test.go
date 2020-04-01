package engine

import (
	"context"
	"testing"

	"github.com/jjzcru/elk/pkg/primitives/ox"
)

func TestDefaultExecuterExecute(t *testing.T) {
	e := ox.Elk{
		Version: "1",
		Tasks: map[string]ox.Task{
			"world": {
				Deps: []ox.Dep{
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
		Logger: make(map[string]Logger),
	}

	_, err := executer.Execute(context.Background(), &e, "world")
	if err != nil {
		t.Error(err)
	}

}

func TestDefaultExecuterExecuteTaskNotExist(t *testing.T) {
	e := ox.Elk{
		Version: "1",
		Tasks: map[string]ox.Task{
			"world": {
				Deps: []ox.Dep{
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
		Logger: make(map[string]Logger),
	}

	_, err := executer.Execute(context.Background(), &e, "bar")
	if err == nil {
		t.Error("task do not exist should throw an error")
	}

}
