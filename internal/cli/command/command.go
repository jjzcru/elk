package command

import (
	"github.com/jjzcru/elk/internal/cli/command/cron"
	"github.com/jjzcru/elk/internal/cli/command/execute"
	"github.com/jjzcru/elk/internal/cli/command/server"
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
		version.Command(),
		initialize.Command(),
		ls.Command(),
		run.Command(),
		execute.Command(),
		cron.Command(),
		logs.Command(),
		server.NewServerCommand(),
	)

	return rootCmd.Execute()
}
