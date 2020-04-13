package server

import (
	"github.com/jjzcru/elk/pkg/server"
	"github.com/jjzcru/elk/pkg/utils"
	"github.com/spf13/cobra"
)

var usageTemplate = `Usage:
  elk server [flags]

Flags:
  -d, --detached      Run the server in detached mode and returns the PGID
  -p, --port          Port where the server is going to run
  -q, --query         Enables graphql playground endpoint 🎮
  -f, --file string   Specify the file to used
  -g, --global        Use global file path
  -h, --help          help for logs
`

// NewServerCommand returns a cobra command for `server` sub command
func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start a graphql server ⚛️",
		Run: func(cmd *cobra.Command, args []string) {
			err := run(cmd, args)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}
	cmd.Flags().IntP("port", "p", 8080, "")
	cmd.Flags().BoolP("query", "q", false, "")
	cmd.Flags().StringP("file", "f", "", "")
	cmd.Flags().BoolP("detached", "d", false, "")
	cmd.Flags().BoolP("global", "g", false, "")

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	isDetached, err := cmd.Flags().GetBool("detached")
	if err != nil {
		return err
	}

	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		return err
	}

	isQueryEnabled, err := cmd.Flags().GetBool("query")
	if err != nil {
		return err
	}

	isGlobal, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}

	elkFilePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	e, err := utils.GetElk(elkFilePath, isGlobal)
	if err != nil {
		return err
	}

	if isDetached {
		return detached()
	}

	return server.Start(port, e.GetFilePath(), isQueryEnabled)
}
