package run

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jjzcru/elk/pkg/utils"

	"github.com/jjzcru/elk/pkg/engine"

	"github.com/spf13/cobra"
)

var usageTemplate = `Usage:
  elk run [tasks] [flags]

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
  -i, --interval            Set a duration for an interval
`

// Command returns a cobra command for `run` sub command
func Command() *cobra.Command {
	var envs []string
	var vars []string
	var cmd = &cobra.Command{
		Use:   "run",
		Short: "Run one or more tasks 🤖",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := Validate(cmd, args)
			if err != nil {
				utils.PrintError(err)
				return
			}

			err = run(cmd, args, envs, vars)
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
	cmd.Flags().BoolP("watch", "w", false, "")
	cmd.Flags().StringP("file", "f", "", "")
	cmd.Flags().StringP("log", "l", "", "")
	cmd.Flags().DurationP("timeout", "t", 0, "")
	cmd.Flags().Duration("delay", 0, "")
	cmd.Flags().String("deadline", "", "")
	cmd.Flags().String("start", "", "")
	cmd.Flags().DurationP("interval", "i", 0, "")

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func run(cmd *cobra.Command, args []string, envs []string, vars []string) error {
	isDetached, err := cmd.Flags().GetBool("detached")
	if err != nil {
		return err
	}

	isWatch, err := cmd.Flags().GetBool("watch")
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

	interval, err := cmd.Flags().GetDuration("interval")
	if err != nil {
		return err
	}

	// Check if the file path is set
	e, err := utils.GetElk(elkFilePath, isGlobal)
	if err != nil {
		return err
	}

	logger, err := Build(cmd, e, args)
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
			task.Env[env] = parts[1]
		}

		for _, v := range vars {
			parts := strings.SplitAfterN(v, "=", 2)
			k := strings.ReplaceAll(parts[0], "=", "")
			task.Vars[k] = parts[1]
		}

		clientEngine.Elk.Tasks[name] = task
	}

	if isDetached {
		return Detached()
	}

	ctx, cancel := context.WithCancel(context.Background())

	if len(start) > 0 {
		startTime, err := GetTimeFromString(start)
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
		deadlineTime, err := GetTimeFromString(deadline)
		if err != nil {
			cancel()
			return err
		}

		ctx, cancel = context.WithDeadline(ctx, deadlineTime)
	}

	DelayStart(delay, start)

	if interval > 0 {
		executeTasks := func() {
			for _, task := range args {
				go TaskWG(ctx, clientEngine, task, nil, false)
			}
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
	for _, task := range args {
		wg.Add(1)
		go TaskWG(ctx, clientEngine, task, &wg, isWatch)
	}

	wg.Wait()
	cancel()

	return nil
}

// DelayStart sleep the program by an amount of time
func DelayStart(delay time.Duration, start string) {
	var startDuration time.Duration
	var delayDuration time.Duration
	var sleepDuration time.Duration

	if len(start) > 0 {
		startTime, _ := GetTimeFromString(start)
		now := time.Now()
		diff := startTime.Sub(now)

		startDuration = diff
	}

	if delay > 0 {
		delayDuration = delay
	}

	if startDuration > 0 && delayDuration > 0 {
		if startDuration > delayDuration {
			sleepDuration = startDuration
		} else {
			sleepDuration = delayDuration
		}
	} else if startDuration > 0 {
		sleepDuration = startDuration
	} else if delayDuration > 0 {
		sleepDuration = delayDuration
	}

	if sleepDuration > 0 {
		time.Sleep(sleepDuration)
	}
}

// GetTimeFromString transform a string to a duration
func GetTimeFromString(input string) (time.Time, error) {
	validTimeFormats := []string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
	}

	for _, layout := range validTimeFormats {
		deadlineTime, err := time.Parse(layout, input)
		if err == nil {
			if layout == time.Kitchen {
				now := time.Now()
				deadlineTime = time.Date(now.Year(),
					now.Month(),
					now.Day(),
					deadlineTime.Hour(),
					deadlineTime.Minute(),
					0,
					0,
					now.Location())

				// If time is before now i refer to that time but the next day
				if deadlineTime.Before(now) {
					deadlineTime = deadlineTime.Add(24 * time.Hour)
				}
			}
			return deadlineTime, nil
		}
	}

	return time.Now(), errors.New("invalid input")
}

// TaskWG runs task with a wait group
func TaskWG(ctx context.Context, cliEngine *engine.Engine, task string, wg *sync.WaitGroup, isWatch bool) {
	if wg != nil {
		defer wg.Done()
	}

	t, err := cliEngine.Elk.GetTask(task)
	if err != nil {
		utils.PrintError(err)
		return
	}

	if len(t.Sources) > 0 && isWatch {
		Watch(ctx, cliEngine, task, *t)
		return
	}

	Task(ctx, cliEngine, task)
}

// Task runs a task on the engine
func Task(ctx context.Context, cliEngine *engine.Engine, task string) {
	ctx, cancel := context.WithCancel(ctx)

	err := cliEngine.Run(ctx, task)
	cancel()
	if err != nil {
		utils.PrintError(err)
		return
	}
}
