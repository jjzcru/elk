package logs

import (
	"bufio"
	"fmt"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/spf13/cobra"
	"io"
	"os"
	"time"
)

var usageTemplate = `Usage:
  elk logs [task] [flags]

Flags:
  -f, --file string   Specify the file to used
      --follow        Run in follow mode
  -g, --global        Search the task in the global path
  -h, --help          Help for logs
`

// NewLogsCommand returns a cobra command for `logs` sub command
func NewLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Attach logs from a task to the terminal 📝",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := validate(cmd, args)
			if err != nil {
				utils.PrintError(err)
				return
			}

			err = run(cmd, args)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().BoolP("global", "g", false, "")
	cmd.Flags().StringP("file", "f", "", "")
	cmd.Flags().Bool("follow", false, "")

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	isFollow, err := cmd.Flags().GetBool("follow")
	if err != nil {
		return err
	}

	e, err := getElk(cmd)
	if err != nil {
		return err
	}

	ch := make(chan string)
	errCh := make(chan error)

	for _, name := range args {
		task, err := e.GetTask(name)
		if err != nil {
			return err
		}

		file, err := os.Open(task.Log.Out)
		if err != nil {
			return err
		}

		go readLogFile(file, ch, errCh, isFollow)
	}

	for {
		select {
		case line := <-ch:
			fmt.Print(line)
		case err := <-errCh:
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func readLogFile(file *os.File, ch chan string, errCh chan error, isFollow bool) {
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if isFollow {
				if err == io.EOF {
					time.Sleep(250 * time.Millisecond)
					reader = bufio.NewReader(file)
				} else {
					break
				}
			} else {
				errCh <- err
			}
		}

		ch <- line
	}
}
