package main

import (
	"os"

	"github.com/jjzcru/elk/internal/cli/command"

	"github.com/jjzcru/elk/internal/cli/command/version"
	"github.com/jjzcru/elk/internal/cli/utils"
)

var v = ""
var o = ""
var arch = ""
var commit = ""
var date = ""

func main() {
	version.SetVersion(v, o, arch, commit, date)
	err := command.Execute()
	if err != nil {
		utils.PrintError(err)
		os.Exit(1)
	}
}
