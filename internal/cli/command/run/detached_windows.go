// +build windows

package run

import (
	"fmt"
	"github.com/jjzcru/elk/pkg/utils"
	"os"
	"os/exec"
)

func Detached() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	command := utils.RemoveDetachedFlag(os.Args)
	cmd := exec.Command(command[0], command[1:]...)
	pid := os.Getpid()
	cmd.Dir = cwd

	err = cmd.Start()
	if err != nil {
		return err
	}

	fmt.Println(pid)
	return nil
}
