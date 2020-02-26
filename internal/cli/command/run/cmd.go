package run

import (
	"context"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/primitives/elk"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/jjzcru/elk/internal/cli/command/config"
	"github.com/jjzcru/elk/pkg/engine"

	"github.com/spf13/cobra"
)

// NewRunCommand returns a cobra command for `run` sub command
func NewRunCommand() *cobra.Command {
	var envs []string
	var cmd = &cobra.Command{
		Use:   "run",
		Short: "Run one or more task in a terminal",
		Run: func(cmd *cobra.Command, args []string) {
			err := run(cmd, args, envs)
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

func run(cmd *cobra.Command, args []string, envs []string) error {
	isDetached, err := cmd.Flags().GetBool("detached")
	if err != nil {
		return err
	}

	isWatch, err := cmd.Flags().GetBool("watch")
	if err != nil {
		return err
	}

	logFilePath, err := cmd.Flags().GetString("log")
	if err != nil {
		return err
	}

	ignoreLog, err := cmd.Flags().GetBool("ignore-log")
	if err != nil {
		return err
	}

	if isDetached {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		command := utils.RemoveDetachedFlag(os.Args)
		cmd := exec.Command(command[0], command[1:]...)
		pid := os.Getpid()
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: pid}
		cmd.Dir = cwd

		err = cmd.Start()
		if err != nil {
			return err
		}

		// _ = cmd.Process.Release()

		fmt.Printf("%d", pid)
		return nil
	}

	isGlobal, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}

	elkFilePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("A task name is required")
	}

	e, err := config.GetElk(elkFilePath, isGlobal)

	if err != nil {
		return err
	}

	executer := engine.DefaultExecuter{
		Logger: &engine.DefaultLogger,
	}

	envMap := engine.MapEnvs(envs)

	e.OverwriteEnvs(envMap)
	// clientEngine := engine.New(elk, executer)

	clientEngine := &engine.Engine{
		Elk:      e,
		Executer: executer,
		Build: func(elk *elk.Elk) error {
			if len(logFilePath) > 0 {
				_, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return err
				}

				absolutePath, err := filepath.Abs(logFilePath)
				if err != nil {
					return err
				}

				logFilePath = absolutePath
			}

			for name, task := range elk.Tasks {
				if len(logFilePath) > 0 {
					task.Log = logFilePath
				}

				if ignoreLog {
					task.Log = ""
				}

				elk.Tasks[name] = task
				/*err := elk.HasCircularDependency(task)
				if err != nil {
					return err
				}*/
			}

			return nil
		},
	}

	var wg sync.WaitGroup

	ctx := context.Background()

	for _, task := range args {
		wg.Add(1)
		go func(task string, wg *sync.WaitGroup) {
			defer wg.Done()

			taskCtx, cancel := context.WithCancel(ctx)

			t, err := e.GetTask(task)
			if err != nil {
				utils.PrintError(err)
				return
			}

			err = e.HasCircularDependency(task)
			if err != nil {
				utils.PrintError(err)
				return
			}

			err = clientEngine.Run(taskCtx, task)
			if err != nil {
				utils.PrintError(err)
				return
			}

			if len(t.Watch) > 0 && isWatch {
				files, err := t.GetWatcherFiles(t.Watch)
				if err != nil {
					utils.PrintError(err)
					return
				}

				watcher, err := fsnotify.NewWatcher()
				if err != nil {
					utils.PrintError(err)
					return
				}
				defer watcher.Close()

				for _, file := range files {
					err = watcher.Add(file)
					if err != nil {
						utils.PrintError(err)
						return
					}
				}

				for {
					select {
					case event := <-watcher.Events:
						switch {
						case event.Op&fsnotify.Write == fsnotify.Write:
							fallthrough
						case event.Op&fsnotify.Create == fsnotify.Create:
							fallthrough
						case event.Op&fsnotify.Remove == fsnotify.Remove:
							fallthrough
						case event.Op&fsnotify.Rename == fsnotify.Rename:
							go func() {
								cancel()
								taskCtx, cancel = context.WithCancel(ctx)

								err = clientEngine.Run(taskCtx, task)
								if err != nil && err != context.Canceled {
									utils.PrintError(err)
									return
								}
							}()
						}
					case err := <-watcher.Errors:
						utils.PrintError(err)
						return
					}
				}
			}
		}(task, &wg)
	}

	wg.Wait()

	return nil
}
