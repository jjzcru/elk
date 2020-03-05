package elk

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func TestElkLoadEnvFile(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d.env", randomNumber)
	err := ioutil.WriteFile(path, []byte("FOO=BAR"), 0644)
	if err != nil {
		t.Error(err)
	}

	e := Elk{
		EnvFile: path,
		Env: map[string]string{
			"HELLO": "World",
		},
	}

	totalOfInitialEnvs := len(e.Env)

	err = e.LoadEnvFile()
	if err != nil {
		t.Error(err)
	}

	if len(e.Env) <= totalOfInitialEnvs {
		t.Error("Expected that the keys from file load to env")
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err)
	}
}

func TestElkLoadEnvFileNotExist(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d.env", randomNumber)

	e := Elk{
		EnvFile: path,
		Env: map[string]string{
			"HELLO": "World",
		},
	}

	err := e.LoadEnvFile()
	if err == nil {
		t.Error("Should throw an error because the file do not exist")
	}
}

func TestElkLoadEnvFileWithNoEnv(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d.env", randomNumber)
	err := ioutil.WriteFile(path, []byte("FOO=BAR"), 0644)
	if err != nil {
		t.Error(err)
	}

	e := Elk{
		EnvFile: path,
		Env:     nil,
	}

	err = e.LoadEnvFile()
	if err != nil {
		t.Error(err)
	}

	if e.Env["FOO"] != "BAR" {
		t.Errorf("The value should be '%s' but is '%s' instead", "BAR", e.Env["FOO"])
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err)
	}
}

func TestHasCircularDependency(t *testing.T) {
	e := Elk{
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

	for taskName := range e.Tasks {
		err := e.HasCircularDependency(taskName)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestElkBuild(t *testing.T) {
	err := os.Setenv("BAR", "1")
	if err != nil {
		t.Error(err)
	}

	randomNumber := rand.Intn(100)

	elkEnvPath := fmt.Sprintf("./elk_%d.env", randomNumber)
	err = ioutil.WriteFile(elkEnvPath, []byte("FOO=BAR"), 0644)
	if err != nil {
		t.Error(err)
	}

	taskEnvPath := fmt.Sprintf("./task_%d.env", randomNumber)
	err = ioutil.WriteFile(taskEnvPath, []byte("FOO=FOO"), 0644)
	if err != nil {
		t.Error(err)
	}

	e := Elk{
		EnvFile: elkEnvPath,
		Env: map[string]string{
			"HELLO": "World",
		},
		Tasks: map[string]Task{
			"hello": {
				EnvFile: taskEnvPath,
			},
		},
	}

	err = e.Build()
	if err != nil {
		t.Error(err)
	}

	hello, err := e.GetTask("hello")
	if err != nil {
		t.Error(err)
	}

	if hello.Env["BAR"] != "1" {
		t.Errorf("The env variable should be '%s' but was '%s' instead", "1", e.Tasks["hello"].Env["BAR"])
	}

	if hello.Env["FOO"] != "FOO" {
		t.Errorf("The env variable should be '%s' but was '%s' instead", "FOO", e.Tasks["hello"].Env["FOO"])
	}

	e.Tasks["world"] = Task{
		Deps: []string{"world"},
	}

	err = e.Build()
	if err == nil {
		t.Error("Should throw an error because it has circular dependency")
	}

	err = os.Remove(elkEnvPath)
	if err != nil {
		t.Error(err)
	}

	err = os.Remove(taskEnvPath)
	if err != nil {
		t.Error(err)
	}
}

func TestHasTask(t *testing.T) {
	e := Elk{
		Tasks: map[string]Task{
			"hello": {},
		},
	}

	hasTask := e.HasTask("hello")
	if !hasTask {
		t.Error("It should have a task")
	}
}

func TestNotHasTask(t *testing.T) {
	e := Elk{
		Tasks: map[string]Task{
			"hello": {},
		},
	}

	hasTask := e.HasTask("world")
	if hasTask {
		t.Error("It should not have a task")
	}
}

func TestGetTask(t *testing.T) {
	e := Elk{
		Tasks: map[string]Task{
			"hello": {},
		},
	}

	_, err := e.GetTask("hello")
	if err != nil {
		t.Error(err)
	}
}

func TestGetTaskNotExist(t *testing.T) {
	e := Elk{
		Tasks: map[string]Task{
			"hello": {},
		},
	}

	_, err := e.GetTask("world")
	if err == nil {
		t.Error("Should throw an error because the task do not exist")
	}
}

func TestGetTaskCircularDependency(t *testing.T) {
	e := Elk{
		Tasks: map[string]Task{
			"hello": {
				Deps: []string{"world"},
			},
			"world": {
				Deps: []string{"hello"},
			},
		},
	}

	_, err := e.GetTask("hello")
	if err == nil {
		t.Error("Should throw an error because the task has a circular dependency")
	}
}
