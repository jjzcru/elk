package run

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jjzcru/elk/internal/cli/utils"

	"github.com/jjzcru/elk/pkg/engine"

	"github.com/spf13/cobra"
)

var usageTemplate = `Usage:
  elk run [tasks] [flags]

Flags:
  -d, --detached         Run the task in detached mode and returns the PGID
  -e, --env strings      Overwrite env variable in task   
  -f, --file string      Run elk in a specific file
  -g, --global           Run from the path set in config
  -h, --help             Help for run
      --ignore-log-file  Force task to output to stdout
      --ignore-error     Ignore errors that happened during a task
      --delay            Set a delay to a task
  -l, --log string       File that log output from a task
  -w, --watch            Enable watch mode
  -t, --timeout          Set a timeout to a task
      --deadline         Set a deadline to a task
      --start            Set a date/datetime to a task to run
  -i, --interval         Set a duration for an interval
`

// NewRunCommand returns a cobra command for `run` sub command
func NewRunCommand() *cobra.Command {
	var envs []string
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

			err = run(cmd, args, envs)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().BoolP("global", "g", false, "Run from the path set in config")
	cmd.Flags().StringSliceVarP(&envs, "env", "e", []string{}, "")
	cmd.Flags().Bool("ignore-log-file", false, "Force task to output to stdout")
	cmd.Flags().Bool("ignore-error", false, "Ignore errors that happened during a task")
	cmd.Flags().BoolP("detached", "d", false, "Run the command in detached mode and returns the PGID")
	cmd.Flags().BoolP("watch", "w", false, "Enable watch mode")
	cmd.Flags().StringP("file", "f", "", "Run elk in a specific file")
	cmd.Flags().StringP("log", "l", "", "File that log output from a task")
	cmd.Flags().DurationP("timeout", "t", 0, "Set a timeout for a task in milliseconds")
	cmd.Flags().Duration("delay", 0, "Set a delay for a task in milliseconds")
	cmd.Flags().String("deadline", "", "Set a deadline to a task")
	cmd.Flags().String("start", "", "Set a date/datetime for a task to run")
	cmd.Flags().DurationP("interval", "i", 0, "Set a duration for an interval")

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func run(cmd *cobra.Command, args []string, envs []string) error {
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

	clientEngine := &engine.Engine{
		Elk: e,
		Executer: engine.DefaultExecuter{
			Logger: &engine.DefaultLogger,
		},
		Build: func() error {
			return Build(cmd, e)
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
		return Detached()
	}

	ctx, cancel := context.WithCancel(context.Background())

	if len(start) > 0 {
		startTime, err := GetTimeFromString(start)
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
		deadlineTime, err := GetTimeFromString(deadline)
		if err != nil {
			return err
		}

		ctx, _ = context.WithDeadline(ctx, deadlineTime)
	}

	DelayStart(delay, start)

	if interval > 0 {
		executeTasks := func() {
			for _, task := range args {
				go runTask(ctx, clientEngine, task, nil, false)
			}
		}

		go executeTasks()
		ticker := time.NewTicker(interval)
		quit := make(chan struct{})
		for {
			select {
			case <-ticker.C:
				go executeTasks()
			case <-ctx.Done():
				ticker.Stop()
				return nil
			case <-quit:
				ticker.Stop()
				cancel()
				return nil
			}
		}
	}

	var wg sync.WaitGroup
	for _, task := range args {
		wg.Add(1)
		go runTask(ctx, clientEngine, task, &wg, isWatch)
	}

	wg.Wait()
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

func runTask(ctx context.Context, cliEngine *engine.Engine, task string, wg *sync.WaitGroup, isWatch bool) {
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
	ctx, _ = context.WithCancel(ctx)

	err := cliEngine.Run(ctx, task)
	if err != nil {
		utils.PrintError(err)
		return
	}
}
