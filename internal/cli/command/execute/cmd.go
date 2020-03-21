package execute

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jjzcru/elk/internal/cli/command/run"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/engine"
	"github.com/jjzcru/elk/pkg/primitives/ox"
	"github.com/spf13/cobra"
)

var usageTemplate = `Usage:
  elk exec [commands] [flags]

Flags:
  -d, --detached         Run the task in detached mode and returns the PGID
  -e, --env strings      Overwrite env variable in task
  -v, --var strings      Overwrite var variable in task
  -h, --help             Help for run
      --delay            Set a delay to a task
  -l, --log string       File that log output from a task
  -w, --watch            Enable watch mode
  -t, --timeout          Set a timeout to a task
      --deadline         Set a deadline to a task
      --start            Set a date/datetime to a task to run
`

// NewRunCommand returns a cobra command for `run` sub command
func NewExecCommand() *cobra.Command {
	var envs []string
	var vars []string
	var cmd = &cobra.Command{
		Use:   "exec",
		Short: "Run ad-hoc commands ⚡",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := Run(cmd, args, envs, vars)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().StringSliceVarP(&envs, "env", "e", []string{}, "")
	cmd.Flags().StringSliceVarP(&vars, "var", "v", []string{}, "")
	cmd.Flags().Bool("ignore-log-file", false, "Force task to output to stdout")
	cmd.Flags().Bool("ignore-error", false, "Ignore errors that happened during a task")
	cmd.Flags().BoolP("detached", "d", false, "Run the command in detached mode and returns the PGID")
	cmd.Flags().StringP("log", "l", "", "File that log output from a task")
	cmd.Flags().DurationP("timeout", "t", 0, "Set a timeout for a task in milliseconds")
	cmd.Flags().Duration("delay", 0, "Set a delay for a task in milliseconds")
	cmd.Flags().String("deadline", "", "Set a deadline to a task")
	cmd.Flags().String("start", "", "Set a date/datetime for a task to run")

	// TODO Support for interval

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func Run(cmd *cobra.Command, args []string, envs []string, vars []string) error {
	isDetached, err := cmd.Flags().GetBool("detached")
	if err != nil {
		return err
	}

	delay, err := cmd.Flags().GetDuration("delay")
	if err != nil {
		return err
	}

	timeout, err := cmd.Flags().GetDuration("timeout")
	if err != nil {
		return err
	}

	deadline, err := cmd.Flags().GetString("deadline")
	if err != nil {
		return err
	}

	start, err := cmd.Flags().GetString("start")
	if err != nil {
		return err
	}

	elk := ox.Elk{
		Env:  make(map[string]string),
		Vars: make(map[string]string),
		Tasks: map[string]ox.Task{
			"elk": {
				Cmds:        args,
				Env:         make(map[string]string),
				Vars:        make(map[string]string),
				IgnoreError: false,
			},
		},
	}

	clientEngine := &engine.Engine{
		Elk: &elk,
		Executer: engine.DefaultExecuter{
			Logger: &engine.DefaultLogger,
		},
		Build: func() error {
			return run.Build(cmd, &elk)
		},
	}

	err = clientEngine.Build()
	if err != nil {
		return err
	}

	for name, task := range elk.Tasks {
		for _, en := range envs {
			parts := strings.SplitAfterN(en, "=", 2)
			env := strings.ReplaceAll(parts[0], "=", "")
			value := parts[1]
			task.Env[env] = value
		}

		for _, v := range vars {
			parts := strings.SplitAfterN(v, "=", 2)
			k := strings.ReplaceAll(parts[0], "=", "")
			task.Vars[k] = parts[1]
		}

		clientEngine.Elk.Tasks[name] = task
	}

	if isDetached {
		return run.Detached()
	}

	ctx := context.Background()

	if len(start) > 0 {
		startTime, err := run.GetTimeFromString(start)
		if err != nil {
			return err
		}

		now := time.Now()
		if startTime.Before(now) {
			return fmt.Errorf("start can't be before of current time")
		}
	}

	if timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, timeout)
	}

	if len(deadline) > 0 {
		deadlineTime, err := run.GetTimeFromString(deadline)
		if err != nil {
			return err
		}

		ctx, _ = context.WithDeadline(ctx, deadlineTime)
	}

	run.DelayStart(delay, start)

	var wg sync.WaitGroup

	wg.Add(1)
	go run.TaskWG(ctx, clientEngine, "elk", &wg, false)

	wg.Wait()
	return nil
}
