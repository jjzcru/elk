package main

import (
	"github.com/jjzcru/elk/internal/cli/command/commands"
	"github.com/jjzcru/elk/internal/cli/command/version"
)

var v = "0.1.0"

func main() {
	version.SetVersion(v)
	commands.Execute()
}
