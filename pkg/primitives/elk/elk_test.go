package elk

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
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
				Deps: []Dep{
					{
						Name: "hello",
					},
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
		Deps: []Dep{
			{
				Name: "world",
			},
		},
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

func TestElkBuildEnvFileDoNotExist(t *testing.T) {
	err := os.Setenv("BAR", "1")
	if err != nil {
		t.Error(err)
	}

	randomNumber := rand.Intn(100)

	elkEnvPath := fmt.Sprintf("./elk_%d.env", randomNumber)

	taskEnvPath := fmt.Sprintf("./task_%d.env", randomNumber)
	err = ioutil.WriteFile(taskEnvPath, []byte("FOO=BAR"), 0644)
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
	if err == nil {
		t.Error("It should throw an error because the env file do not exist")
	}

	err = os.Remove(taskEnvPath)
	if err != nil {
		t.Error(err)
	}
}

func TestElkBuildEnvFileDoNotExistInTask(t *testing.T) {
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
	if err == nil {
		t.Error("It should throw an error because the env file do not exist")
	}

	err = os.Remove(elkEnvPath)
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
				Deps: []Dep{
					{
						Name: "world",
					},
				},
			},
			"world": {
				Deps: []Dep{
					{
						Name: "hello",
					},
				},
			},
		},
	}

	_, err := e.GetTask("hello")
	if err == nil {
		t.Error("Should throw an error because the task has a circular dependency")
	}
}

func TestFromFile(t *testing.T) {
	e := Elk{
		Env: make(map[string]string),
		Tasks: map[string]Task{
			"hello": {
				Env: make(map[string]string),
				Deps: []Dep{
					{
						Name: "world",
					},
				},
				Cmds: []string{
					"echo Hello",
				},
			},
			"world": {
				Env: make(map[string]string),
				Cmds: []string{
					"echo Hello",
				},
				Deps: []Dep{
					{
						Name: "hello",
					},
				},
			},
		},
	}

	content, err := yaml.Marshal(&e)
	if err != nil {
		t.Error(err)
	}

	path := fmt.Sprint("./elk.yml")
	err = ioutil.WriteFile(path, content, 0644)
	if err != nil {
		t.Error(err)
	}

	elk, err := FromFile(path)
	if err != nil {
		t.Error(err)
	}

	for task := range elk.Tasks {
		if !reflect.DeepEqual(e.Tasks[task], elk.Tasks[task]) {
			t.Errorf("The data is different for the task '%s'", task)
		}
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err)
	}
}

func TestFromFileWithoutTasks(t *testing.T) {
	e := Elk{}

	content, err := yaml.Marshal(&e)
	if err != nil {
		t.Error(err)
	}

	path := fmt.Sprint("./elk.yml")
	err = ioutil.WriteFile(path, content, 0644)
	if err != nil {
		t.Error(err)
	}

	elk, err := FromFile(path)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(elk.Env, make(map[string]string)) {
		t.Error("The env should be an empty map")
	}

	if !reflect.DeepEqual(elk.Tasks, make(map[string]Task)) {
		t.Error("The tasks should be an empty map")
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err)
	}
}

func TestFromFileNotExist(t *testing.T) {
	path := fmt.Sprint("./elk.yml")

	_, err := FromFile(path)
	if err == nil {
		t.Error("it should throw an error because the file do not exist")
	}
}

func TestFromFileInvalidFileContent(t *testing.T) {
	path := fmt.Sprint("./elk.yml")
	err := ioutil.WriteFile(path, []byte("FOO=BAR"), 0644)
	if err != nil {
		t.Error(err)
	}

	_, err = FromFile(path)
	if err == nil {
		t.Error("it should throw an error because the file do not exist")
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err)
	}
}
