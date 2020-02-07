package main

import "github.com/jjzcru/elk/internal/commands"

/*import (
	"context"
	"fmt"
	"os"
	"strings"

	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)*/

var version = "0.1.0"

func main() {
	commands.SetVersion(version)
	commands.Execute()
	// run("docker-compose -v")
}

/*
func run(command string) error {
	p, err := syntax.NewParser().Parse(strings.NewReader(command), "")
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	r, err := interp.New(
		interp.Dir(dir),
		// interp.Env(expand.ListEnviron(envs...)),

		interp.Module(interp.DefaultExec),
		interp.Module(interp.OpenDevImpls(interp.DefaultOpen)),

		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
	)

	if err != nil {
		return err
	}

	ctx := context.Background()

	return r.Run(ctx, p)
}
*/
