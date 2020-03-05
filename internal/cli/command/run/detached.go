// +build !windows

package run

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/jjzcru/elk/internal/cli/utils"
)

func runDetached() error {
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

	fmt.Println(pid)
	return nil
}
