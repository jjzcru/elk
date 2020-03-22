package command

import (
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	os.Args = []string{"ox"}
	err := Execute()
	if err != nil {
		t.Error(err)
	}
}
