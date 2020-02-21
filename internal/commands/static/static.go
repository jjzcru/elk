package static

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"net/http"

	"github.com/jjzcru/elk/internal/commands/config"
	"github.com/spf13/cobra"
)


// Cmd Command that runs a task
func Cmd() *cobra.Command {
	var command = &cobra.Command{
		Use:   "static",
		Short: "Load a static file website",
		Run: func(cmd *cobra.Command, args []string) {
			path := "."
			if len(args) > 0 {
				path = args[0]
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
				cmd.Dir = cwd

				err = cmd.Start()
				if err != nil {
					config.PrintError(err.Error())
					return
				}

				cmd.Process.Release()
				return
			}


			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				config.PrintError(err.Error())
				return
			}

			fs := http.FileServer(http.Dir(path))
			http.Handle("/", fs)

			fmt.Printf("Server listening on port: %d\n", port)
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

func removeDetachedFlag(args []string) []string {
	cmd := []string{}

	for _, arg := range args {
		if len(arg) > 0 && arg != "-d" && arg != "--detached" {
			cmd = append(cmd, strings.TrimSpace(arg))
		}
	}

	return cmd
}