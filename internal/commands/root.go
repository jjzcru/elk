package commands

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"

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

	runCmd.Flags().BoolP("global", "g", false, "Run from the path set in config")
	runCmd.Flags().StringP("file", "f", "", "Specify an alternate elk file \n(default: elk.yml)")
}

func printError(err string) {
	fmt.Print(aurora.Bold(aurora.Red("ERROR: ")))
	_, _ = fmt.Fprintf(os.Stderr, err)
	fmt.Println()
}
