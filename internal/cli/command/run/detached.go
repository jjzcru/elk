// +build !windows

package run

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/jjzcru/elk/pkg/utils"
)

// Detached runs ox in detached mode
func Detached() error {
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

	fmt.Println(pid)
	return nil
}
