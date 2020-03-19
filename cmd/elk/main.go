package main

import (
	"os"

	"github.com/jjzcru/elk/internal/cli/command"

	"github.com/jjzcru/elk/internal/cli/command/version"
	"github.com/jjzcru/elk/internal/cli/utils"
)

var v = "0.3.1"

func main() {
	version.SetVersion(v)
	err := command.Execute()
	if err != nil {
		utils.PrintError(err)
		os.Exit(1)
	}
}
