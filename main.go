package main

import (
	OS "os"

	"github.com/jjzcru/elk/internal/cli/command"

	v "github.com/jjzcru/elk/internal/cli/command/version"
	"github.com/jjzcru/elk/pkg/utils"
)

var version = ""
var os = ""
var arch = ""
var commit = ""
var date = ""
var goversion = ""

func main() {
	v.SetVersion(version, os, arch, commit, date, goversion)
	err := command.Execute()
	if err != nil {
		utils.PrintError(err)
		OS.Exit(1)
	}
}
