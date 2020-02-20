package install

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"text/template"

	"github.com/jjzcru/elk/internal/commands/config"
	in "github.com/jjzcru/elk/internal/commands/initialize"
	"github.com/spf13/cobra"
)

// Cmd Command that install elk in the system
func Cmd() *cobra.Command {
	var command = &cobra.Command{
		Use:   "install",
		Short: "Install elk in the system",
		Run: func(cmd *cobra.Command, args []string) {
			err := install()
			if err != nil {
				config.PrintError("A task name is required")
				return
			}
			fmt.Println("Elk was installed successfully")
		},
	}

	return command
}

func install() error {
	installationPath, err := getInstallationPath(".elk")
	if err != nil {
		return err
	}

	installationExist, err := isInstallationDirExist(installationPath)
	if err != nil {
		return err
	}

	if !installationExist {
		fmt.Printf("Installation directory: %s\n", installationPath)
		err = createInstallationDir(installationPath)
		if err != nil {
			return err
		}
	}

	return createGlobalIfNotExist()
}

func getInstallationPath(dir string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(usr.HomeDir, dir), nil
}

func isInstallationDirExist(installationDirPath string) (bool, error) {
	if _, err := os.Stat(installationDirPath); os.IsNotExist(err) {
		return false, nil
	}

	if _, err := os.Stat(path.Join(installationDirPath, "config.yml")); os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

func createInstallationDir(installationDirPath string) error {
	if _, err := os.Stat(installationDirPath); os.IsNotExist(err) {
		err := os.Mkdir(installationDirPath, 0777)
		if err != nil {

			return err
		}
	}

	return createInstallConfigFile(installationDirPath)
}

func createInstallConfigFile(installationDirPath string) error {

	configPath := path.Join(installationDirPath, "config.yml")
	configFile, err := os.Create(configPath)

	usr, err := user.Current()
	if err != nil {
		return err
	}

	response, err := template.New("config").Parse(config.ConfigTemplate)
	if err != nil {
		return err
	}

	err = response.Execute(configFile, path.Join(usr.HomeDir, "elk.yml"))

	return createGlobalIfNotExist()
}

func createGlobalIfNotExist() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	existGlobal, err := isExistGlobalElkFile(path.Join(usr.HomeDir, "elk.yml"))
	if err != nil {
		return err
	}

	if !existGlobal {
		elkFilePath := path.Join(usr.HomeDir, "elk.yml")
		fmt.Printf("Elkfile: %s\n", elkFilePath)
		return in.CreateElkFile(elkFilePath)
	}

	return nil
}

func isExistGlobalElkFile(elkFilePath string) (bool, error) {
	if _, err := os.Stat(elkFilePath); os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}
