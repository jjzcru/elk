package utils

import (
	"strings"
	"testing"
)

func TestRemoveDetachedFlag(t *testing.T) {
	args := []string{"ox", "run", "test", "-d"}
	args = RemoveDetachedFlag(args)

	expectedCmd := "ox run test"
	cmd := strings.Join(args, " ")
	if cmd != expectedCmd {
		t.Errorf("The command should be '%s' but it is '%s' instead", expectedCmd, cmd)
	}
}
