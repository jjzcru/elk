package static

import (
	"fmt"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/logrusorgru/aurora"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

// NewStaticCommand returns a cobra command for `static` sub command
func NewStaticCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "static",
		Short: "Load a static file website",
		Run: func(cmd *cobra.Command, args []string) {
			err := run(cmd, args)
			if err != nil {
				utils.PrintError(err)
				return
			}
		},
	}

	cmd.Flags().BoolP("detached", "d", false, "Run the command in detached mode")
	cmd.Flags().IntP("port", "p", 3000, "Set server port")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	path := "."

	if len(args) > 0 {
		path = args[0]
	}

	path, err := getWorkingDirectoryPath(path)
	if err != nil {
		return err
	}

	isDetached, err := cmd.Flags().GetBool("detached")
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

		fmt.Printf("%d", pid)
		return nil
	}

	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		return err
	}

	fs := http.FileServer(http.Dir(path))
	http.Handle("/", fs)

	fmt.Printf("Static files: %s\n", aurora.Cyan(path))
	fmt.Printf("Server listening on port 🚀: %s\n", aurora.Cyan(strconv.Itoa(port)))
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		return err
	}

	return nil
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
