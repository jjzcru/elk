package cron

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jjzcru/elk/internal/cli/command/run"
	"github.com/jjzcru/elk/pkg/engine"
	"github.com/jjzcru/elk/pkg/utils"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var usageTemplate = `Usage:
  elk cron [crontab] [tasks] [flags]

Flags:
  -d, --detached            Run the task in detached mode and returns the PGID
  -e, --env strings         Overwrite env variable in task
  -v, --var strings         Overwrite var variable in task   
  -f, --file string         Run elk in a specific file
  -g, --global              Run from the path set in config
  -h, --help                Help for run
      --ignore-log-file     Ignores task log property
      --ignore-log-format   Ignores format value in log
      --ignore-error        Ignore errors that happened during a task
      --ignore-deps         Ignore task dependencies
      --delay               Set a delay to a task
  -l, --log string          File that log output from a task
  -w, --watch               Enable watch mode
  -t, --timeout             Set a timeout to a task
      --deadline            Set a deadline to a task
      --start               Set a date/datetime to a task to run
`

// Command returns a cobra command for `run` sub command
func Command() *cobra.Command {
	var envs []string
	var vars []string
	var cmd = &cobra.Command{
		Use:   "cron",
		Short: "Run one or more task as a cron job ⏱",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			err := run.Validate(cmd, args[1:])
			if err != nil {
				utils.PrintError(err)
				return
			}

			err = Run(cmd, args, envs, vars)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().BoolP("global", "g", false, "")
	cmd.Flags().StringSliceVarP(&envs, "env", "e", []string{}, "")
	cmd.Flags().StringSliceVarP(&vars, "var", "v", []string{}, "")
	cmd.Flags().Bool("ignore-log-file", false, "")
	cmd.Flags().Bool("ignore-log-format", false, "")
	cmd.Flags().Bool("ignore-error", false, "")
	cmd.Flags().Bool("ignore-deps", false, "")
	cmd.Flags().BoolP("detached", "d", false, "")
	cmd.Flags().StringP("file", "f", "", "")
	cmd.Flags().StringP("log", "l", "", "")
	cmd.Flags().DurationP("timeout", "t", 0, "")
	cmd.Flags().Duration("delay", 0, "")
	cmd.Flags().String("deadline", "", "")
	cmd.Flags().String("start", "", "")

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func Run(cmd *cobra.Command, args []string, envs []string, vars []string) error {
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

	logger, err := run.Build(cmd, e, args)
	if err != nil {
		return err
	}

	clientEngine := &engine.Engine{
		Elk: e,
		Executer: engine.DefaultExecuter{
			Logger: logger,
		},
	}

	for name, task := range e.Tasks {
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
	var cancel context.CancelFunc

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
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	if len(deadline) > 0 {
		deadlineTime, err := run.GetTimeFromString(deadline)
		if err != nil {
			return err
		}

		ctx, cancel = context.WithDeadline(ctx, deadlineTime)
	}

	if cancel != nil {
		defer cancel()
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
	/*select {
	case <-ctx.Done():
		c.Stop()
		return nil
	}*/

	<-ctx.Done()
	c.Stop()
	return nil
}
