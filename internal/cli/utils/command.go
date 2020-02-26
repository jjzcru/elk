package utils

import "strings"

// RemoveDetachedFlag returns a command without the detached flag
func RemoveDetachedFlag(args []string) []string {
	var cmd []string

	for _, arg := range args {
		if len(arg) > 0 && arg != "-d" && arg != "--detached" {
			cmd = append(cmd, strings.TrimSpace(arg))
		}
	}

	return cmd
}
