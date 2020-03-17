package cron

import (
	"context"
	"fmt"
	"github.com/jjzcru/elk/internal/cli/command/run"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/engine"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var usageTemplate = `Usage:
  elk cron [crontab] [tasks] [flags]

Examples:
elk cron "*/1 * * * *" foo
elk cron "*/1 * * * *" foo bar
elk cron "*/1 * * * *" foo -d
elk cron "*/2 * * * *" foo -t 1s
elk cron "*/2 * * * *" foo --delay 1s
elk cron "*/2 * * * *" foo -e FOO=BAR --env HELLO=WORLD
elk cron "*/6 * * * *" foo -l ./foo.log -d
elk cron "*/1 * * * *" foo --ignore-log
elk cron "*/2 * * * *" foo --ignore-error
elk cron "*/5 * * * *" foo --deadline 09:41AM
elk cron "*/1 * * * *" foo --start 09:41PM

Flags:
  -d, --detached      Run the task in detached mode and returns the PGID
  -e, --env strings   Overwrite env variable in task   
  -f, --file string   Run elk in a specific file
  -g, --global        Run from the path set in config
  -h, --help          help for run
      --ignore-log    Force task to output to stdout
      --ignore-error  Ignore errors that happened during a task
      --delay         Set a delay to a task
  -l, --log string    File that log output from a task
  -w, --watch         Enable watch mode
  -t, --timeout       Set a timeout to a task
      --deadline      Set a deadline to a task
      --start      	  Set a date/datetime to a task to run
`

// NewRunCommand returns a cobra command for `run` sub command
func NewCronCommand() *cobra.Command {
	var envs []string
	var cmd = &cobra.Command{
		Use:   "cron",
		Short: "Run one or more task as a cron job",
		Args:  cobra.MinimumNArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return run.Validate(cmd, args[1:])
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := Run(cmd, args, envs)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().BoolP("global", "g", false, "Run from the path set in config")
	cmd.Flags().StringSliceVarP(&envs, "env", "e", []string{}, "")
	cmd.Flags().Bool("ignore-log", false, "Force task to output to stdout")
	cmd.Flags().Bool("ignore-error", false, "Ignore errors that happened during a task")
	cmd.Flags().BoolP("detached", "d", false, "Run the command in detached mode and returns the PGID")
	cmd.Flags().StringP("file", "f", "", "Run elk in a specific file")
	cmd.Flags().StringP("log", "l", "", "File that log output from a task")
	cmd.Flags().DurationP("timeout", "t", 0, "Set a timeout for a task in milliseconds")
	cmd.Flags().Duration("delay", 0, "Set a delay for a task in milliseconds")
	cmd.Flags().String("deadline", "", "Set a deadline to a task")
	cmd.Flags().String("start", "", "Set a date/datetime for a task to run")

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func Run(cmd *cobra.Command, args []string, envs []string) error {
	isDetached, err := cmd.Flags().GetBool("detached")
	if err != nil {
		return err
	}

	elkFilePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	isGlobal, err := cmd.Flags().GetBool("global")
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

	// Check if the file path is set
	e, err := utils.GetElk(elkFilePath, isGlobal)
	if err != nil {
		return err
	}

	clientEngine := &engine.Engine{
		Elk: e,
		Executer: engine.DefaultExecuter{
			Logger: &engine.DefaultLogger,
		},
		Build: func() error {
			return run.Build(cmd, e)
		},
	}

	err = clientEngine.Build()
	if err != nil {
		return err
	}

	for name, task := range e.Tasks {
		for _, en := range envs {
			parts := strings.SplitAfterN(en, "=", 2)
			env := strings.ReplaceAll(parts[0], "=", "")
			value := parts[1]
			task.Env[env] = value
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

	cronTab := args[0]
	tasks := args[1:]

	c := cron.New()

	run.DelayStart(delay, start)

	for _, task := range tasks {
		go run.Task(ctx, clientEngine, task)
	}

	_, err = c.AddFunc(cronTab, func() {
		for _, task := range tasks {
			go run.Task(ctx, clientEngine, task)
		}
	})
	if err != nil {
		return err
	}

	c.Start()
	select {
	case <-ctx.Done():
		c.Stop()
		return nil
	}
}
