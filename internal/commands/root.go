package commands

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"os"

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
}

func printError(err string) {
	fmt.Print(aurora.Bold(aurora.Red("ERROR: ")))
	_, _ = fmt.Fprintf(os.Stderr, err)
	fmt.Println()
}
