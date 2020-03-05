package main

import (
	"github.com/jjzcru/elk/internal/cli/command/commands"
	"github.com/jjzcru/elk/internal/cli/command/version"
	"github.com/jjzcru/elk/internal/cli/utils"
	"os"
)

var v = "0.1.0"

func main() {
	version.SetVersion(v)
	err := commands.Execute()
	if err != nil {
		utils.PrintError(err)
		os.Exit(1)
	}
}
