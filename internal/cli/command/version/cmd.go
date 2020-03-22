package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

type data struct {
	version string
	os      string
	arch    string
	commit  string
	date    string
}

var version data

// NewVersionCommand returns a cobra command for `version` sub command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version number",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Println("Elk version " + version)
			fmt.Print("Elk 🦌\n\n")
			fmt.Printf("Version:    \t %s\n", version.version)
			fmt.Printf("Git Commit: \t %s\n", version.commit)
			fmt.Printf("Built:      \t %s\n", strings.Replace(version.date, "_", " ", -1))
			fmt.Printf("OS:         \t %s\n", version.os)
			fmt.Printf("Arch:       \t %s\n", version.arch)
		},
	}

	return cmd
}

// SetVersion is a function that prints what is the current version of the cli
func SetVersion(v string, os string, arch string, commit string, date string) {
	version = data{
		version: v,
		os:      os,
		arch:    arch,
		commit:  commit,
		date:    date,
	}
}
