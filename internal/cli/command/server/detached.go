package server

import (
	"fmt"
	"github.com/jjzcru/elk/pkg/utils"
	"os"
	"os/exec"
)

func detached(token string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	command := utils.RemoveDetachedFlag(os.Args)
	var commands []string
	commands = append(commands, command[1:]...)

	if len(token) > 0 {
		commands = append(commands, []string{"--token", token}...)
	}

	cmd := exec.Command(command[0], commands...)
	cmd.Dir = cwd

	err = cmd.Start()
	if err != nil {
		return err
	}

	pid := cmd.Process.Pid

	if len(token) > 0 {
		fmt.Printf("%d %s", pid, token)
	} else {
		fmt.Println(pid)
	}

	return nil
}
