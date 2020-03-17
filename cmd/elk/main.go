package main

import (
	"github.com/jjzcru/elk/internal/cli/command"
	"os"

	"github.com/jjzcru/elk/internal/cli/command/version"
	"github.com/jjzcru/elk/internal/cli/utils"
)

var v = "0.2.1"

func main() {
	version.SetVersion(v)
	err := command.Execute()
	if err != nil {
		utils.PrintError(err)
		os.Exit(1)
	}
}
