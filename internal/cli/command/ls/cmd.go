package ls

import (
	"fmt"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/primitives/ox"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var usageTemplate = `Usage:
  elk ls [flags]

Flags:
  -a, --all           Display task details
  -f, --file string   Specify the file to used
  -g, --global        Search the task in the global path
  -h, --help          help for logs
`

// NewListCommand returns a cobra command for `ls` sub command
func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List tasks",
		Run: func(cmd *cobra.Command, args []string) {
			err := run(cmd, args)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().BoolP("all", "a", false, "")
	cmd.Flags().StringP("file", "f", "", "")
	cmd.Flags().BoolP("global", "g", false, "")

	cmd.SetUsageTemplate(usageTemplate)

	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	isGlobal, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}

	elkFilePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	shouldPrintAll, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}

	e, err := utils.GetElk(elkFilePath, isGlobal)

	if err != nil {
		return err
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 24, 8, 0, '\t', 0)
	defer w.Flush()

	if shouldPrintAll {
		return printAll(w, e)
	}

	return printPlain(w, e)
}

func printAll(w *tabwriter.Writer, e *ox.Elk) error {
	_, err := fmt.Fprintf(w, "\n%s\t%s\t%s\t\n", "TASK NAME", "DESCRIPTION", "DEPENDENCIES")
	if err != nil {
		return err
	}

	for taskName, task := range e.Tasks {
		var deps []string

		for _, dep := range task.Deps {
			deps = append(deps, dep.Name)
		}

		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t\n", taskName, task.Description, strings.Join(deps, ", "))
		if err != nil {
			return err
		}
	}

	return nil
}

func printPlain(w *tabwriter.Writer, elk *ox.Elk) error {
	_, err := fmt.Fprintf(w, "\n%s\t%s\t\n", "TASK NAME", "DESCRIPTION")
	if err != nil {
		return err
	}

	for taskName, task := range elk.Tasks {
		_, err = fmt.Fprintf(w, "%s\t%s\t\n", taskName, task.Description)
		if err != nil {
			return err
		}
	}

	return nil
}
