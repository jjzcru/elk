package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCommand returns a cobra command for `version` sub command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Elk version " + version)
		},
	}

	return cmd
}

// SetVersion is a function that prints what is the current version of the cli
func SetVersion(v string) {
	version = v
}
