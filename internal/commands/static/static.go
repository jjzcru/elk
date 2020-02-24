package static

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/jjzcru/elk/internal/commands/config"
	"github.com/spf13/cobra"
)


// Cmd Command that runs a task
func Cmd() *cobra.Command {
	var command = &cobra.Command{
		Use:   "static",
		Short: "Load a static file website",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			path := "."

			if len(args) > 0 {
				path = args[0]
			}

			path, err = getWorkingDirectoryPath(path)
			if err != nil {
				config.PrintError(err.Error())
				return
			}

			isDetached, err := cmd.Flags().GetBool("detached")
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

				fmt.Printf("%d", pid)
				return
			}


			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				config.PrintError(err.Error())
				return
			}

			fs := http.FileServer(http.Dir(path))
			http.Handle("/", fs)



			fmt.Printf("Static files: %s\n", aurora.Cyan(path))
			fmt.Printf("Server listening on port 🚀: %s\n", aurora.Cyan(strconv.Itoa(port)))
			err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
			if err != nil {
				config.PrintError(err.Error())
				return
			}
		},
	}

	command.Flags().BoolP("detached", "d", false, "Run the command in detached mode")
	command.Flags().IntP("port", "p", 3000, "Set server port")

	return command
}

func getWorkingDirectoryPath(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if !fileInfo.IsDir() {
		path = filepath.Dir(path)
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absolutePath, nil
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