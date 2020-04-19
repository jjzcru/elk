package ox

import (
	"errors"
	"fmt"
	"testing"
)

func TestVarsWrite(t *testing.T) {
	vars := Vars{
		Map: map[string]string{
			"foo": "bar",
		},
	}

	cmd := "echo {{.foo}}"

	bytes, err := vars.Write([]byte(cmd))
	if err != nil {
		t.Error(err)
	}

	if bytes != len([]byte(cmd)) {
		t.Error(fmt.Errorf("the total of bytes should be %d but it was %d instead", len([]byte(cmd)), bytes))
	}

	if cmd != vars.Cmd {
		t.Error(fmt.Errorf("the command should be '%s' but it was '%s' instead", cmd, vars.Cmd))
	}
}

func TestVarsProcess(t *testing.T) {
	vars := Vars{
		Map: map[string]string{
			"foo": "bar",
		},
	}

	inputCmd := "echo {{.foo}}"
	expectedCmd := "echo bar"

	cmd, err := vars.Process(inputCmd)
	if err != nil {
		t.Error(err)
	}

	if cmd != expectedCmd {
		t.Error(fmt.Errorf("the command should be '%s' but it was '%s' instead", expectedCmd, cmd))
	}
}

func TestVarsProcessErrorParsing(t *testing.T) {
	vars := Vars{
		Map: map[string]string{
			"foo": "bar",
		},
	}

	_, err := vars.Process("{{.foo{}")
	if err == nil {
		t.Error(errors.New("it should throw an error of invalid syntax"))
	}
}
