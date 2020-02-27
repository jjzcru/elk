package kill

import (
	"errors"
	"fmt"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/spf13/cobra"
	"strconv"
	"syscall"
)

// NewKillCommand returns a cobra command for `kill` sub command
func NewKillCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kill",
		Short: "Kill a process by PID",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires an ID argument")
			}

			for _, arg := range args {
				_, err := strconv.Atoi(arg)
				if err != nil {
					return fmt.Errorf("only integers, value not valid: %s", arg)
				}
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := run(cmd, args)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().BoolP("pgid", "g", false, "Kill process by PGID")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	isPgid, err := cmd.Flags().GetBool("pgid")
	if err != nil {
		return err
	}

	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return err
		}

		if !isPgid {
			id = id * -1
		}

		err = syscall.Kill(id*-1, syscall.SIGKILL)
		if err != nil {
			return err
		}
	}

	return nil
}
