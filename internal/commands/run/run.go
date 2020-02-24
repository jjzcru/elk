package run

import (
	"fmt"
	"github.com/jjzcru/elk/pkg/primitives"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/jjzcru/elk/internal/commands/config"
	"github.com/jjzcru/elk/pkg/engine"

	"github.com/spf13/cobra"
)

// Cmd Command that runs a task
func Cmd() *cobra.Command {
	var envs []string
	var command = &cobra.Command{
		Use:   "run",
		Short: "Run one or more task in a terminal",
		Run: func(cmd *cobra.Command, args []string) {
			isDetached, err := cmd.Flags().GetBool("detached")
			if err != nil {
				config.PrintError(err.Error())
				return
			}

			logFilePath, err := cmd.Flags().GetString("log")
			if err != nil {
				config.PrintError(err.Error())
				return
			}

			if isDetached {
				cwd, err := os.Getwd()
				if err != nil {
					config.PrintError(err.Error())
					return
				}

				command := removeDetachedFlag(os.Args)
				cmd := exec.Command(command[0], command[1:]...)
				pid := os.Getpid()
				cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: pid}
				cmd.Dir = cwd

				err = cmd.Start()
				if err != nil {
					config.PrintError(err.Error())
					return
				}

				// _ = cmd.Process.Release()

				fmt.Printf("%d", pid)
				return
			}

			isGlobal, err := cmd.Flags().GetBool("global")
			if err != nil {
				config.PrintError(err.Error())
				return
			}

			elkFilePath, err := cmd.Flags().GetString("file")
			if err != nil {
				config.PrintError(err.Error())
				return
			}

			if len(args) == 0 {
				config.PrintError("A task name is required")
				return
			}

			elk, err := config.GetElk(elkFilePath, isGlobal)

			if err != nil {
				config.PrintError(err.Error())
				return
			}

			executer := engine.DefaultExecuter{
				Logger: &engine.DefaultLogger,
			}

			elk.OverwriteEnvs(engine.MapEnvs(envs))
			// clientEngine := engine.New(elk, executer)
			clientEngine := &engine.Engine{
				Elk:      elk,
				Executer: executer,
				Build: func(elk *primitives.Elk) error {
					// Validate if there is a circular dependency
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
							elk.Tasks[name] = task
						}
						/*err := elk.HasCircularDependency(task)
						if err != nil {
							return err
						}*/
					}

					return nil
				},
			}

			var wg sync.WaitGroup

			for _, task := range args {
				wg.Add(1)
				go func(task string, wg *sync.WaitGroup) {
					defer wg.Done()

					if !elk.HasTask(task) {
						config.PrintError(fmt.Sprintf("task '%s' do not exist", task))
						return
					}

					err = clientEngine.Run(task)
					if err != nil {
						config.PrintError(err.Error())
						return
					}
				}(task, &wg)
			}

			wg.Wait()
		},
	}

	command.Flags().StringSliceVarP(&envs, "env", "e", []string{}, "")
	command.Flags().BoolP("detached", "d", false, "Run the command in detached mode and returns the PGID")
	command.Flags().StringP("file", "f", "", "Run elk in a specific file")
	command.Flags().StringP("log", "l", "", "File that log output from a task")

	return command
}

func removeDetachedFlag(args []string) []string {
	var cmd []string

	for _, arg := range args {
		if len(arg) > 0 && arg != "-d" && arg != "--detached" {
			cmd = append(cmd, strings.TrimSpace(arg))
		}
	}

	return cmd
}

type detachedLogger struct{}

func (d detachedLogger) Write(p []byte) (n int, err error) {
	f, err := os.OpenFile("/home/jjzcru/Desktop/test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return len(p), err
	}

	defer f.Close()

	_, err = f.Write(p)
	if err != nil {
		return len(p), err
	}
	return len(p), nil
}
