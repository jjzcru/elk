package logs

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/file"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var usageTemplate = `Usage:
  elk logs [task] [flags]

Flags:
  -f, --file string   Specify ox.yml file to be used
      --follow        Run in follow mode
  -g, --global        Search the task in the global path
  -h, --help          Help for logs
`

// NewLogsCommand returns a cobra command for `logs` sub command
func NewLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Attach logs from a task to the terminal 📝",
		Args:  cobra.MinimumNArgs(1),
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

	// Will use this to get the task with the larger name
	taskNameLength := 0
	for _, name := range args {
		if len(name) > taskNameLength {
			taskNameLength = len(name)
		}
	}

	for _, name := range args {
		task, err := e.GetTask(name)
		if err != nil {
			return err
		}

		f, err := os.Open(task.Log.Out)
		if err != nil {
			return err
		}

		f.Close()
		var prefix string
		if len(args) > 1 || (len(task.Log.Out) > 0 && len(task.Log.Err) > 0) {
			prefix = getColorPrefix(name, taskNameLength)
		}

		go readLogFile(task.Log.Out, ch, errCh, isFollow, prefix)

		if len(task.Log.Err) > 0 {
			f, err = os.Open(task.Log.Err)
			if err != nil {
				return err
			}

			f.Close()

			prefix = getErrorColorPrefix(name, taskNameLength)
			go readLogFile(task.Log.Err, ch, errCh, isFollow, prefix)
		}
	}

	for {
		select {
		case line := <-ch:
			if len(line) > 0 {
				fmt.Println(line)
			}
		case err := <-errCh:
			if err == io.EOF {
				fmt.Println()
				return nil
			}
			return err
		}
	}
}

func readLogFile(filename string, ch chan string, errCh chan error, isFollow bool, prefix string) {
	f, _ := os.Open(filename)
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()
	defer f.Close()
	_ = watcher.Add(filename)
	delimiter := byte('\n')

	r := bufio.NewReader(f)
	if isFollow {
		for {
			by, err := r.ReadBytes(delimiter)
			if err != nil && err != io.EOF {
				errCh <- err
			}
			ch <- getStringFromBytes(by, prefix)
			if err != io.EOF {
				continue
			}
			if err = waitForChange(watcher); err != nil {
				errCh <- err
			}
		}
	} else {
		for {
			by, err := r.ReadBytes(delimiter)
			if err != nil && err != io.EOF {
				errCh <- err
			}
			ch <- getStringFromBytes(by, prefix)
			if err == io.EOF {
				errCh <- err
				close(ch)
				break
			}

		}
	}
}

func getStringFromBytes(by []byte, prefix string) string {
	content := string(by)
	content = strings.ReplaceAll(content, file.BreakLine, "")
	clearScreenSequece := "[2J"

	if len(content) == 0 {
		return ""
	}

	if strings.Contains(content, clearScreenSequece) {
		return ""
	}

	// Do not add prefix if we are clearing the screen
	if len(prefix) > 0 {
		return prefix + content
	}

	return content
}

func waitForChange(w *fsnotify.Watcher) error {
	for {
		select {
		case event := <-w.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				return nil
			}
		case err := <-w.Errors:
			return err
		}
	}
}

func getColorPrefix(name string, taskNameLength int) string {
	name = getPrefixName(name, taskNameLength)
	rand.Seed(time.Now().UnixNano())
	switch n := rand.Intn(10); n {
	case 0:
		return aurora.Bold(aurora.Green(fmt.Sprintf("%s | ", name))).String()
	case 1:
		return aurora.Bold(aurora.Yellow(fmt.Sprintf("%s | ", name))).String()
	case 3:
		return aurora.Bold(aurora.BrightMagenta(fmt.Sprintf("%s | ", name))).String()
	case 4:
		return aurora.Bold(aurora.Blue(fmt.Sprintf("%s | ", name))).String()
	case 5:
		return aurora.Bold(aurora.Magenta(fmt.Sprintf("%s | ", name))).String()
	case 6:
		return aurora.Bold(aurora.Cyan(fmt.Sprintf("%s | ", name))).String()
	case 7:
		return aurora.Bold(aurora.BrightGreen(fmt.Sprintf("%s | ", name))).String()
	case 8:
		return aurora.Bold(aurora.BrightYellow(fmt.Sprintf("%s | ", name))).String()
	case 9:
		return aurora.Bold(aurora.BrightCyan(fmt.Sprintf("%s | ", name))).String()
	default:
		return aurora.Bold(aurora.BrightBlue(fmt.Sprintf("%s | ", name))).String()
	}
}

func getErrorColorPrefix(name string, taskNameLength int) string {
	return aurora.Bold(aurora.Red(fmt.Sprintf("%s | ", getPrefixName(name, taskNameLength)))).String()
}

func getPrefixName(name string, taskNameLength int) string {
	difference := taskNameLength - len(name)
	for i := 0; i < difference; i++ {
		name += " "
	}

	return name
}
