package kill

import (
	"github.com/jjzcru/elk/internal/commands/config"
	"github.com/spf13/cobra"
	"strconv"
	"syscall"
)

func Cmd() *cobra.Command {
	var command = &cobra.Command{
		Use:   "kill",
		Short: "Kill a process by PID",
		Run: func(cmd *cobra.Command, args []string) {
			isPgid, err := cmd.Flags().GetBool("pgid")
			if err != nil {
				config.PrintError(err.Error())
				return
			}

			for _, arg := range args {
				id, err := strconv.Atoi(arg)
				if err != nil {
					config.PrintError(err.Error())
					return
				}

				if !isPgid {
					id = id * -1
				}

				err = syscall.Kill(id * -1, syscall.SIGKILL)
				if err != nil {
					config.PrintError(err.Error())
					return
				}
			}
		},
	}

	command.Flags().BoolP("pgid", "p", false, "Kill process by PGID")
	return command
}