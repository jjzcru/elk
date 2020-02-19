package run

import (
	"fmt"
	"os"
	"sync"

	"github.com/jjzcru/elk/internal/commands/config"
	"github.com/jjzcru/elk/pkg/engine"

	"github.com/spf13/cobra"
)

// Cmd Command that runs a task
var Cmd = &cobra.Command{
	Use:   "run",
	Short: "Run one or more task in a terminal",
	Run: func(cmd *cobra.Command, args []string) {
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

		logger := &engine.Logger{
			StdoutWriter: os.Stdout,
			StderrWriter: os.Stderr,
			StdinReader:  os.Stdin,
		}

		clientEngine := engine.New(elk, logger)

		var wg sync.WaitGroup

		for _, task := range args {
			wg.Add(1)
			go func(task string, wg *sync.WaitGroup) {
				defer wg.Done()

				if !clientEngine.HasTask(task) {
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
