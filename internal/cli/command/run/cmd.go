package run

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jjzcru/elk/internal/cli/command/config"
	"github.com/jjzcru/elk/internal/cli/utils"

	"github.com/jjzcru/elk/pkg/engine"

	"github.com/spf13/cobra"
)

var usageTemplate = `Usage:
  elk run [tasks] [flags]

Examples:
elk run foo
elk run foo bar
elk run foo -d
elk run foo -d -w
elk run foo -t 10000
elk run foo -e FOO=BAR -e HELLO=WORLD
elk run foo -l ./foo.log -d
elk run foo --ignore-log
elk run foo --deadline 10:00PM

Flags:
  -d, --detached      Run the command in detached mode and returns the PGID
  -e, --env strings   Overwrite env variable in task   
  -f, --file string   Run elk in a specific file
  -g, --global        Run from the path set in config
  -h, --help          help for run
      --ignore-log    Force task to output to stdout
  -l, --log string    File that log output from a task
  -w, --watch         Enable watch mode
  -t, --timeout       Set a timeout for a task in miliseconds
      --deadline      Set a deadline to a task
`

// NewRunCommand returns a cobra command for `run` sub command
func NewRunCommand() *cobra.Command {
	var envs []string
	var cmd = &cobra.Command{
		Use:   "run",
		Short: "Run one or more task in a terminal",
		Args:  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validate(cmd, args)
			// return validate(cmd, args, &e)
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := run(cmd, args)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().BoolP("global", "g", false, "Run from the path set in config")
	cmd.Flags().StringSliceVarP(&envs, "env", "e", []string{}, "")
	cmd.Flags().Bool("ignore-log", false, "Force task to output to stdout")
	cmd.Flags().BoolP("detached", "d", false, "Run the command in detached mode and returns the PGID")
	cmd.Flags().BoolP("watch", "w", false, "Enable watch mode")
	cmd.Flags().StringP("file", "f", "", "Run elk in a specific file")
	cmd.Flags().StringP("log", "l", "", "File that log output from a task")
	cmd.Flags().Int32P("timeout", "t", 0, "Set a timeout for a task")
	cmd.Flags().String("deadline", "", "Set a deadline to a task")

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
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

	// Check if the file path is set
	e, err := config.GetElk(elkFilePath, isGlobal)
	if err != nil {
		return err
	}

	clientEngine := &engine.Engine{
		Elk: e,
		Executer: engine.DefaultExecuter{
			Logger: &engine.DefaultLogger,
		},
		Build: func() error {
			return build(cmd, e)
		},
	}

	err = clientEngine.Build()
	if err != nil {
		return err
	}

	if isDetached {
		return runDetached()
	}

	var wg sync.WaitGroup
	ctx := context.Background()

	timeout, err := cmd.Flags().GetInt32("timeout")
	if err != nil {
		return err
	}

	if timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	}

	deadline, err := cmd.Flags().GetString("deadline")
	if err != nil {
		return err
	}

	if len(deadline) > 0 {
		deadlineTime, err := getDeadlineTime(deadline)
		if err != nil {
			return err
		}

		ctx, _ = context.WithDeadline(ctx, deadlineTime)
	}

	for _, task := range args {
		wg.Add(1)
		go runTask(ctx, clientEngine, task, &wg, isWatch)
	}

	wg.Wait()

	return nil
}

func getDeadlineTime(deadline string) (time.Time, error) {
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
		deadlineTime, err := time.Parse(layout, deadline)
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

	return time.Now(), errors.New("invalid deadline")
}

func runTask(ctx context.Context, cliEngine *engine.Engine, task string, wg *sync.WaitGroup, isWatch bool) {
	defer wg.Done()

	taskCtx, cancel := context.WithCancel(ctx)

	t, err := cliEngine.Elk.GetTask(task)
	if err != nil {
		utils.PrintError(err)
		cancel()
		return
	}

	if len(t.Watch) > 0 && isWatch {
		runWatch(cliEngine, taskCtx, task, t, cancel, ctx)
		cancel()
		return
	}

	err = cliEngine.Run(taskCtx, task)
	if err != nil {
		utils.PrintError(err)
		return
	}
	cancel()
}
