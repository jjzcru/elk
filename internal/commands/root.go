package commands

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/spf13/cobra"
)

var version = "1"
var hermesFilePath string

var rootCmd = &cobra.Command{
	Use:   "elk",
	Short: "Task runner",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Placeholder text")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().Bool("author", true, "Author name for copyright attribution")
	// rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	// rootCmd.PersistentFlags().StringVarP(&hermesFilePath, "file", "f", "", "Hermesfile path")
	registerCommands()
}

func getElkFilePath() (string, error) {
	// Hermesfile was not provided by the flag
	if len(hermesFilePath) == 0 {
		currentDir, err := os.Getwd()
		if err != nil {
			return "", err
		}

		hermesFilePath = path.Join(currentDir, "Elkfile.yml")
		if _, err := os.Stat(hermesFilePath); err == nil {
			return hermesFilePath, nil
		}

		usr, _ := user.Current()
		hermesFilePath := path.Join(usr.HomeDir, ".elk", "Elkfile.yml")
		if _, err := os.Stat(hermesFilePath); err == nil {
			return hermesFilePath, nil
		}
	}

	if _, err := os.Stat(hermesFilePath); err == nil {
		return hermesFilePath, nil
	}

	return "", fmt.Errorf("the file 'Elkfile.yml' was not found in the path or system")
}

func registerCommands() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(runCmd)
}
