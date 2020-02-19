package main

import (
	"github.com/jjzcru/elk/internal/commands"
	"github.com/jjzcru/elk/internal/commands/version"
)

var v = "0.1.0"

func main() {
	version.SetVersion(v)
	commands.Execute()
}
