package commands

import (
	"fmt"
	"os"

	"github.com/jjzcru/elk/internal/commands/ls"

	"github.com/spf13/cobra"
)

var version = "1"

var rootCmd = &cobra.Command{
	Use:   "elk",
	Short: "A simple yml based task runner",
}

// Execute starts the CLI application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(ls.Cmd)

	runCmd.Flags().BoolP("global", "g", false, "Run from the path set in config")
	runCmd.Flags().StringP("file", "f", "", "Specify an alternate elk file \n(default: elk.yml)")

	ls.Cmd.Flags().BoolP("global", "g", false, "Run from the path set in config")
	ls.Cmd.Flags().StringP("file", "f", "", "Specify an alternate elk file \n(default: elk.yml)")
	ls.Cmd.Flags().BoolP("all", "a", false, "Print tasks details")
}
