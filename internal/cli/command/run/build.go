package run

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/primitives/elk"
	"github.com/spf13/cobra"
)

func build(cmd *cobra.Command, e *elk.Elk) error {
	ignoreLog, err := cmd.Flags().GetBool("ignore-log")
	if err != nil {
		return err
	}

	logFilePath, err := cmd.Flags().GetString("log")
	if err != nil {
		return err
	}

	if len(logFilePath) > 0 {
		isFile, err := utils.IsPathAFile(logFilePath)
		if err != nil {
			return err
		}

		if !isFile {
			return fmt.Errorf("path is not a file: %s", logFilePath)
		}
	}

	if e.Env == nil {
		e.Env = make(map[string]string)
	}

	if len(logFilePath) > 0 {
		_, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		absolutePath, err := filepath.Abs(logFilePath)
		if err != nil {
			return err
		}

		logFilePath = absolutePath
	}

	elkEnvCopy := make(map[string]string)
	for k, v := range e.Env {
		elkEnvCopy[k] = v
	}

	loadSystemEnvVariable(e)
	err = loadElkEnvFile(e)
	if err != nil {
		return err
	}

	for env, value := range elkEnvCopy {
		e.Env[env] = value
	}

	for name, task := range e.Tasks {
		if len(logFilePath) > 0 {
			task.Log = logFilePath
		}

		if ignoreLog {
			task.Log = ""
		}

		err := e.HasCircularDependency(name)
		if err != nil {
			return err
		}

		// If not env variable is set create an empty map
		if task.Env == nil {
			task.Env = make(map[string]string)
		}

		// Keep a copy of the existing map
		taskEnvCopy := make(map[string]string)
		for k, v := range task.Env {
			taskEnvCopy[k] = v
		}

		// Load env variables from Elk
		for env, value := range e.Env {
			task.Env[env] = value
		}

		// Load from file
		err = loadTaskEnvFile(&task)
		if err != nil {
			return err
		}

		for env, value := range taskEnvCopy {
			task.Env[env] = value
		}

		e.Tasks[name] = task
	}

	return nil
}

func loadSystemEnvVariable(e *elk.Elk) {
	for _, en := range os.Environ() {
		parts := strings.SplitAfterN(en, "=", 2)
		env := strings.ReplaceAll(parts[0], "=", "")
		value := parts[1]
		e.Env[env] = value
	}
}

func loadElkEnvFile(e *elk.Elk) error {
	if len(e.EnvFile) > 0 {
		info, err := os.Stat(e.EnvFile)
		if os.IsNotExist(err) {
			return err
		}

		if info.IsDir() {
			return fmt.Errorf("log path '%s' is a directory", e.EnvFile)
		}

		file, err := os.Open(e.EnvFile)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			text := scanner.Text()
			if len(text) > 1 {
				parts := strings.SplitAfterN(text, "=", 2)
				env := strings.ReplaceAll(parts[0], "=", "")
				value := parts[1]
				e.Env[env] = value
			}
		}
	}

	return nil
}

func loadTaskEnvFile(task *elk.Task) error {
	if len(task.EnvFile) > 0 {
		info, err := os.Stat(task.EnvFile)
		if os.IsNotExist(err) {
			return err
		}

		if info.IsDir() {
			return fmt.Errorf("log path '%s' is a directory", task.EnvFile)
		}

		file, err := os.Open(task.EnvFile)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			text := scanner.Text()
			if len(text) > 1 {
				parts := strings.SplitAfterN(text, "=", 2)
				env := strings.ReplaceAll(parts[0], "=", "")
				value := parts[1]
				task.Env[env] = value
			}
		}
	}

	return nil
}
