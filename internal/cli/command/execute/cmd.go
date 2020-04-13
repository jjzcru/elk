package execute

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jjzcru/elk/internal/cli/command/run"
	"github.com/jjzcru/elk/pkg/engine"
	"github.com/jjzcru/elk/pkg/primitives/ox"
	"github.com/jjzcru/elk/pkg/utils"
	"github.com/spf13/cobra"
)

var usageTemplate = `Usage:
  elk exec [commands] [flags]

Flags:
  -d, --detached           Run the commands in detached mode and returns the PGID
  -e, --env strings        Overwrite env variable in commands
      --env-file string    Set an env file
  -v, --var strings        Overwrite var variable in commands
  -h, --help               Help for run
      --delay              Set a delay to a task
      --dir                Set a directory to the command
  -l, --log string         File that log output from the commands
      --ignore-error       Ignore errors that happened during a task
  -t, --timeout            Set a timeout to the commands
      --deadline           Set a deadline to the commands
      --start              Set a date/datetime to the commands to run
  -i, --interval           Set a duration for an interval 
`

// Command returns a cobra command for `exec` sub command
func Command() *cobra.Command {
	var envs []string
	var vars []string
	var cmd = &cobra.Command{
		Use:   "exec",
		Short: "Execute ad-hoc commands ⚡",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := Run(cmd, args, envs, vars)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().Bool("ignore-log-format", false, "")
	cmd.Flags().BoolP("detached", "d", false, "")
	cmd.Flags().StringSliceVarP(&envs, "env", "e", []string{}, "")
	cmd.Flags().String("env-file", "", "")
	cmd.Flags().StringSliceVarP(&vars, "var", "v", []string{}, "")
	cmd.Flags().Duration("delay", 0, "")
	cmd.Flags().String("dir", "", "")
	cmd.Flags().StringP("log", "l", "", "")
	cmd.Flags().Bool("ignore-error", false, "")
	cmd.Flags().DurationP("timeout", "t", 0, "")
	cmd.Flags().String("deadline", "", "")
	cmd.Flags().String("start", "", "")
	cmd.Flags().DurationP("interval", "i", 0, "")

	cmd.Flags().Bool("ignore-log-file", false, "")

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

// Run the command
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

	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}

	if len(dir) > 0 {
		isDir, err := utils.IsPathADir(dir)
		if err != nil {
			return err
		}

		if !isDir {
			return fmt.Errorf("path '%s' is not a directory", dir)
		}
	}

	envFile, err := cmd.Flags().GetString("env-file")
	if err != nil {
		return err
	}

	if len(envFile) > 0 {
		isFile, err := utils.IsPathAFile(envFile)
		if err != nil {
			return err
		}

		if !isFile {
			return fmt.Errorf("path '%s' is not a file", envFile)
		}
	}

	interval, err := cmd.Flags().GetDuration("interval")
	if err != nil {
		return err
	}

	ignoreError, err := cmd.Flags().GetBool("ignore-error")
	if err != nil {
		return err
	}

	elk := ox.Elk{
		Tasks: map[string]ox.Task{
			"elk": {
				Cmds:        args,
				Dir:         dir,
				EnvFile:     envFile,
				Env:         make(map[string]string),
				Vars:        make(map[string]string),
				IgnoreError: ignoreError,
			},
		},
	}

	logger, err := run.Build(cmd, &elk, []string{"elk"})
	if err != nil {
		return err
	}

	clientEngine := &engine.Engine{
		Elk: &elk,
		Executer: engine.DefaultExecuter{
			Logger: logger,
		},
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

	ctx, cancel := context.WithCancel(context.Background())

	if len(start) > 0 {
		startTime, err := run.GetTimeFromString(start)
		if err != nil {
			cancel()
			return err
		}

		now := time.Now()
		if startTime.Before(now) {
			cancel()
			return fmt.Errorf("start can't be before of current time")
		}
	}

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	if len(deadline) > 0 {
		deadlineTime, err := run.GetTimeFromString(deadline)
		if err != nil {
			cancel()
			return err
		}

		ctx, cancel = context.WithDeadline(ctx, deadlineTime)
	}

	run.DelayStart(delay, start)

	if interval > 0 {
		executeTasks := func() {
			go run.TaskWG(ctx, clientEngine, "elk", nil, false)
		}

		go executeTasks()
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				go executeTasks()
			case <-ctx.Done():
				ticker.Stop()
				cancel()
				return nil
			}
		}
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go run.TaskWG(ctx, clientEngine, "elk", &wg, false)

	wg.Wait()
	cancel()
	return nil
}
