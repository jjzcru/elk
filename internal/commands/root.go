package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var version = "1"

var rootCmd = &cobra.Command{
	Use:   "elk",
	Short: "Task runner",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(installCmd)
}
