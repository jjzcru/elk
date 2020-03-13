package commands

import (
	"os"

	initialize "github.com/jjzcru/elk/internal/cli/command/initialize"
	"github.com/jjzcru/elk/internal/cli/command/logs"
	"github.com/jjzcru/elk/internal/cli/command/ls"
	"github.com/jjzcru/elk/internal/cli/command/run"
	"github.com/jjzcru/elk/internal/cli/command/version"
	"github.com/spf13/cobra"
)

// Execute starts the CLI application
func Execute() error {
	var rootCmd = &cobra.Command{
		Use:   "elk",
		Short: "Minimalist yaml based task runner 🦌",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
		},
	}
	rootCmd.AddCommand(
		version.NewVersionCommand(),
		initialize.NewInitializeCommand(),
		ls.NewListCommand(),
		run.NewRunCommand(),
		logs.NewLogsCommand(),
	)

	return rootCmd.Execute()
}
