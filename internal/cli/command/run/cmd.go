package run

import (
	"context"
	"sync"

	"github.com/jjzcru/elk/internal/cli/command/config"
	"github.com/jjzcru/elk/internal/cli/utils"

	"github.com/jjzcru/elk/pkg/engine"

	"github.com/spf13/cobra"
)

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
	for _, task := range args {
		wg.Add(1)
		go runTask(ctx, clientEngine, task, &wg, isWatch)
	}

	wg.Wait()

	return nil
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
		cancel()
		return
	}
}
