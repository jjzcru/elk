package main

import (
	"os"

	"github.com/jjzcru/elk/internal/cli/command/commands"
	"github.com/jjzcru/elk/internal/cli/command/version"
	"github.com/jjzcru/elk/internal/cli/utils"
)

var v = "0.2.0"

func main() {
	version.SetVersion(v)
	err := commands.Execute()
	if err != nil {
		utils.PrintError(err)
		os.Exit(1)
	}
}
