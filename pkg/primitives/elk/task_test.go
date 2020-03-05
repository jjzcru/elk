package elk

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func TestTaskLoadEnvFile(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d.env", randomNumber)
	err := ioutil.WriteFile(path, []byte("FOO=BAR"), 0644)
	if err != nil {
		t.Error(err)
	}

	task := Task{
		EnvFile: path,
		Env: map[string]string{
			"HELLO": "World",
		},
	}

	totalOfInitialEnvs := len(task.Env)

	err = task.LoadEnvFile()
	if err != nil {
		t.Error(err)
	}

	if len(task.Env) <= totalOfInitialEnvs {
		t.Error("Expected that the keys from file load to env")
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err)
	}
}

func TestTaskLoadEnvFileNotExist(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d.env", randomNumber)

	task := Task{
		EnvFile: path,
		Env: map[string]string{
			"HELLO": "World",
		},
	}

	err := task.LoadEnvFile()
	if err == nil {
		t.Error("Should throw an error because the file do not exist")
	}
}

func TestTaskLoadEnvFileWithNoEnv(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d.env", randomNumber)
	err := ioutil.WriteFile(path, []byte("FOO=BAR"), 0644)
	if err != nil {
		t.Error(err)
	}

	task := Task{
		EnvFile: path,
		Env:     nil,
	}

	err = task.LoadEnvFile()
	if err != nil {
		t.Error(err)
	}

	if task.Env["FOO"] != "BAR" {
		t.Errorf("The value should be '%s' but is '%s' instead", "BAR", task.Env["FOO"])
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err)
	}
}

func TestTaskGetEnvs(t *testing.T) {
	task := Task{
		Env: map[string]string{
			"HELLO": "World",
		},
	}

	envs := task.GetEnvs()
	if len(envs) == 0 {
		t.Error("Should return 1 env variable")
		return
	}

	if envs[0] != "HELLO=World" {
		t.Errorf("The result is returning '%s' and should be '%s'", envs[0], "HELLO=World")
	}
}
