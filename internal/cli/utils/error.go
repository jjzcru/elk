package utils

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"os"
)

// PrintError display the error message in the cli
func PrintError(err error) {
	fmt.Print(aurora.Bold(aurora.Red("ERROR: ")))
	_, _ = fmt.Fprintf(os.Stderr, err.Error())
	fmt.Println()
}
