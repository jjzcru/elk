package ls

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/jjzcru/elk/internal/commands/config"
	"github.com/jjzcru/elk/pkg/engine"
	"github.com/spf13/cobra"
)

// Cmd Command that display the task in the elk file
var Cmd = &cobra.Command{
	Use:   "ls",
	Short: "List tasks",
	Run: func(cmd *cobra.Command, args []string) {
		isGlobal, err := cmd.Flags().GetBool("global")
		if err != nil {
			config.PrintError(err.Error())
		}

		elkFilePath, err := cmd.Flags().GetString("file")
		if err != nil {
			config.PrintError(err.Error())
			return
		}

		shouldPrintAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			config.PrintError(err.Error())
			return
		}

		elk, err := config.GetElk(elkFilePath, isGlobal)

		if err != nil {
			config.PrintError(err.Error())
			return
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 24, 8, 0, '\t', 0)
		defer w.Flush()

		if shouldPrintAll {
			printAll(w, elk)
		} else {
			printPlain(w, elk)
		}
	},
}

func printAll(w *tabwriter.Writer, elk *engine.Elk) {
	fmt.Fprintf(w, "\n%s\t%s\t%s\t\n", "Task", "Description", "Dependencies")
	fmt.Fprintf(w, "%s\t%s\t%s\t\n", "----", "----", "----")

	for taskName, task := range elk.Tasks {
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", taskName, task.Description, strings.Join(append(task.Deps, task.DetachedDeps...), ", "))
	}
}

func printPlain(w *tabwriter.Writer, elk *engine.Elk) {
	fmt.Fprintf(w, "\n%s\t%s\t\n", "Task", "Description")
	fmt.Fprintf(w, "%s\t%s\t\n", "----", "----")

	for taskName, task := range elk.Tasks {
		fmt.Fprintf(w, "%s\t%s\t\n", taskName, task.Description)
	}
}
