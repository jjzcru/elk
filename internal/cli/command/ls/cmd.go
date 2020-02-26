package ls

import (
	"fmt"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/primitives/elk"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/jjzcru/elk/internal/cli/command/config"
	"github.com/spf13/cobra"
)

// NewListCommand returns a cobra command for `ls` sub command
func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List tasks",
		Run: func(cmd *cobra.Command, args []string) {
			err := run(cmd)
			if err != nil {
				utils.PrintError(err)
			}
		},
	}

	cmd.Flags().BoolP("all", "a", false, "Print tasks details")
	cmd.Flags().StringP("file", "f", "", "Run elk in a specific file")
	cmd.Flags().BoolP("global", "g", false, "Run from the path set in config")

	return cmd
}

func run(cmd *cobra.Command) error {
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

	e, err := config.GetElk(elkFilePath, isGlobal)

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

func printAll(w *tabwriter.Writer, e *elk.Elk) error {
	_, err := fmt.Fprintf(w, "\n%s\t%s\t%s\t\n", "TASK", "DESCRIPTION", "DEPENDENCIES")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s\t%s\t%s\t\n", "----", "-----------", "------------")
	if err != nil {
		return err
	}

	for taskName, task := range e.Tasks {
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t\n", taskName, task.Description, strings.Join(append(task.Deps, task.DetachedDeps...), ", "))
		if err != nil {
			return err
		}
	}

	return nil
}

func printPlain(w *tabwriter.Writer, elk *elk.Elk) error {
	_, err := fmt.Fprintf(w, "\n%s\t%s\t\n", "TASK", "DESCRIPTION")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s\t%s\t\n", "----", "-----------")
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
