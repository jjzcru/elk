package utils

import (
	"strings"
	"testing"
)

func TestRemoveDetachedFlag(t *testing.T) {
	args := []string{"elk", "run", "test", "-d"}
	args = RemoveDetachedFlag(args)

	expectedCmd := "elk run test"
	cmd := strings.Join(args, " ")
	if cmd != expectedCmd {
		t.Errorf("The command should be '%s' but it is '%s' instead", expectedCmd, cmd)
	}
}