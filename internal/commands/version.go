package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// SetVersion is a function that prints what is the current version of the cli
func SetVersion(v string) {
	version = v
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Elk version " + version)
	},
}
