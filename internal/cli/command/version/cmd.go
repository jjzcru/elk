package version

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type data struct {
	version   string
	os        string
	arch      string
	commit    string
	date      string
	goVersion string
}

var version data

// Command returns a cobra command for `version` sub command
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version number",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Elk 🦌\n")
			if len(version.version) > 0 {
				fmt.Printf("  Version:    \t %s\n", version.version)
			}

			if len(version.commit) > 0 {
				fmt.Printf("  Git Commit: \t %s\n", version.commit)
			}

			if len(version.date) > 0 {
				fmt.Printf("  Built:      \t %s\n", strings.Replace(version.date, "_", " ", -1))
			}

			if (len(version.os) + len(version.arch)) > 0 {
				fmt.Printf("  OS/Arch:    \t %s/%s\n", version.os, version.arch)
			}

			if len(version.goVersion) > 0 {
				fmt.Printf("  Go Version: \t %s\n", version.goVersion)
			}
		},
	}

	return cmd
}

// SetVersion is a function that prints what is the current version of the cli
func SetVersion(v string, os string, arch string, commit string, date string, goVersion string) {
	version = data{
		version:   v,
		os:        os,
		arch:      arch,
		commit:    commit,
		date:      date,
		goVersion: strings.ReplaceAll(goVersion, "go", ""),
	}
}
