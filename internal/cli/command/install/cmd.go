package install

import (
	"fmt"
	"github.com/jjzcru/elk/internal/cli/templates"
	"github.com/jjzcru/elk/internal/cli/utils"
	"os"
	"os/user"
	"path"
	"text/template"

	in "github.com/jjzcru/elk/internal/cli/command/initialize"
	"github.com/spf13/cobra"
)

// NewInstallCommand returns a cobra command for 'install' sub command
func NewInstallCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "install",
		Short: "Install elk in the system",
		Run: func(cmd *cobra.Command, args []string) {
			err := install()
			if err != nil {
				utils.PrintError(err)
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

func isInstallationDirExist(dirPath string) (bool, error) {
	if !utils.IsPathExist(dirPath) || !utils.IsPathExist(path.Join(dirPath, "config.yml")) {
		return false, nil
	}

	return true, nil
}

func createInstallationDir(installationDirPath string) error {
	if !utils.IsPathExist(installationDirPath) {
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

	response, err := template.New("config").Parse(templates.Config)
	if err != nil {
		return err
	}

	err = response.Execute(configFile, path.Join(usr.HomeDir, "elk.yml"))
	if err != nil {
		return err
	}

	return createGlobalIfNotExist()
}

func createGlobalIfNotExist() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	elkFilePath := path.Join(usr.HomeDir, "elk.yml")

	if !utils.IsPathExist(elkFilePath) {
		fmt.Printf("Elk file: %s\n", elkFilePath)
		return in.CreateElkFile(elkFilePath)
	}

	return nil
}
