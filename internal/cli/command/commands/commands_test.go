package commands

import (
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	os.Args = []string{"elk"}
	err := Execute()
	if err != nil {
		t.Error(err)
	}
}
