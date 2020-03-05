package utils

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func TestIsPathExist(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d", randomNumber)
	err := ioutil.WriteFile(path, []byte(""), 0644)
	if err != nil {
		t.Error(err.Error())
	}

	exist := IsPathExist(path)

	if !exist {
		t.Errorf("The path '%s' should exist", path)
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestIsPathNotExist(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d", randomNumber)

	exist := IsPathExist(path)

	if exist {
		t.Errorf("The path '%s' should not exist", path)
	}
}

func TestIsPathADir(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d", randomNumber)

	err := os.Mkdir(path,0777)
	if err != nil {
		t.Error(err.Error())
	}

	isADir, err := IsPathADir(path)
	if err != nil {
		t.Error(err.Error())
	}

	if !isADir {
		t.Errorf("The path '%s' should be a directory", path)
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestIsPathADirNotExist(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d", randomNumber)

	_, err := IsPathADir(path)
	if err == nil {
		t.Errorf("The path '%s' should not exist", path)
	}
}

func TestIsPathIsNotADir(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d", randomNumber)
	err := ioutil.WriteFile(path, []byte(""), 0644)
	if err != nil {
		t.Error(err.Error())
	}

	isADir, err := IsPathADir(path)
	if err != nil {
		t.Error(err.Error())
	}

	if isADir {
		t.Errorf("The path '%s' should be a file", path)
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestIsPathAFile(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d", randomNumber)

	err := ioutil.WriteFile(path, []byte(""), 0644)
	if err != nil {
		t.Error(err.Error())
	}

	isAFile, err := IsPathAFile(path)
	if err != nil {
		t.Error(err.Error())
	}

	if !isAFile {
		t.Errorf("The path '%s' should be a file", path)
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestIsPathNotAFile(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d", randomNumber)

	err := os.Mkdir(path,0777)
	if err != nil {
		t.Error(err.Error())
	}

	isAFile, err := IsPathAFile(path)
	if err != nil {
		t.Error(err.Error())
	}

	if isAFile {
		t.Errorf("The path '%s' should be a dir", path)
	}

	err = os.Remove(path)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestIsPathNotAFileNotExist(t *testing.T) {
	randomNumber := rand.Intn(100)
	path := fmt.Sprintf("./%d", randomNumber)

	_, err := IsPathAFile(path)
	if err == nil {
		t.Errorf("The path '%s' should not exist", path)
	}
}