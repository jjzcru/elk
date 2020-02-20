package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd Command that prints current version
func Cmd() *cobra.Command {
	var command = &cobra.Command{
		Use:   "version",
		Short: "Print version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Elk version " + version)
		},
	}
	return command
}

// SetVersion is a function that prints what is the current version of the cli
func SetVersion(v string) {
	version = v
}
