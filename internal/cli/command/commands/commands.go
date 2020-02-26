package commands

import (
	"os"

	"github.com/jjzcru/elk/internal/cli/command/config"
	initialize "github.com/jjzcru/elk/internal/cli/command/initialize"
	"github.com/jjzcru/elk/internal/cli/command/install"
	"github.com/jjzcru/elk/internal/cli/command/kill"
	"github.com/jjzcru/elk/internal/cli/command/ls"
	"github.com/jjzcru/elk/internal/cli/command/run"
	"github.com/jjzcru/elk/internal/cli/command/static"
	"github.com/jjzcru/elk/internal/cli/command/version"

	"github.com/jjzcru/elk/internal/cli/utils"

	"github.com/spf13/cobra"
)

// Execute starts the CLI application
func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "elk",
		Short: "Simple yml based task runner 🦌",
	}
	rootCmd.AddCommand(
		config.NewConfigCommand(),
		version.NewVersionCommand(),
		install.NewInstallCommand(),
		initialize.NewInitializeCommand(),
		kill.NewKillCommand(),
		static.NewStaticCommand(),
		ls.NewListCommand(),
		run.NewRunCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		utils.PrintError(err)
		os.Exit(1)
	}
}
