package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd Command that prints current version
var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Elk version " + version)
	},
}

// SetVersion is a function that prints what is the current version of the cli
func SetVersion(v string) {
	version = v
}
