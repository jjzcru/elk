package file

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func TestGetEnvFromFile(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d.env", randomNumber)
	err := ioutil.WriteFile(path, []byte("FOO=BAR"), 0644)
	if err != nil {
		t.Error(err)
	}

	mapEnv, err := GetEnvFromFile(path)
	if err != nil {
		t.Error(err)
	}

	if mapEnv["FOO"] != "BAR" {
		t.Errorf("Expected to be '%s' but was '%s'", "BAR", mapEnv["FOO"])
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err)
	}
}

func TestGetEnvFromFileNotExist(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d.env", randomNumber)
	err := ioutil.WriteFile(path, []byte("FOO=BAR"), 0644)
	if err != nil {
		t.Error(err)
	}

	randomNumber = rand.Intn(100)
	wrongPath := fmt.Sprintf("./%d.env", randomNumber)

	_, err = GetEnvFromFile(wrongPath)
	if err == nil {
		t.Error("It should throw an error because the file do not exist")
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err)
	}
}

func TestGetEnvFDir(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d", randomNumber)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	_, err = GetEnvFromFile(path)
	if err == nil {
		t.Error("It should throw an error because the file do not exist")
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err)
	}
}
