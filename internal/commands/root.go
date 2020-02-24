package commands

import (
	"fmt"
	"github.com/jjzcru/elk/internal/commands/kill"
	"os"

	"github.com/jjzcru/elk/internal/commands/config"
	in "github.com/jjzcru/elk/internal/commands/initialize"
	"github.com/jjzcru/elk/internal/commands/install"
	"github.com/jjzcru/elk/internal/commands/ls"
	"github.com/jjzcru/elk/internal/commands/run"
	"github.com/jjzcru/elk/internal/commands/static"
	"github.com/jjzcru/elk/internal/commands/version"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "elk",
	Short: "Simple yml based task runner 🦌",
}

// Execute starts the CLI application
func Execute() {
	start()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func start() {
	rootCmd.AddCommand(version.Cmd())
	rootCmd.AddCommand(install.Cmd())
	rootCmd.AddCommand(in.Cmd())
	rootCmd.AddCommand(run.Cmd())
	rootCmd.AddCommand(kill.Cmd())
	rootCmd.AddCommand(ls.Cmd())
	rootCmd.AddCommand(config.Cmd())
	rootCmd.AddCommand(static.Cmd())

	rootCmd.PersistentFlags().BoolP("global", "g", false, "Run from the path set in config")
}
