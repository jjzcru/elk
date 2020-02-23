package model

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Task is the data structure for the task to run
type Task struct {
	Cmds         []string
	Env          map[string]string
	detached     bool
	Description  string
	Dir          string
	Log          string
	Deps         []string
	DetachedDeps []string `yaml:"detached_deps"`
	EnvFile      string   `yaml:"env_file"`
}

func (t *Task) SetDetached(detached bool) {
	t.detached = detached
}

// IsDetached Check if the task is detached
func (t *Task) IsDetached() bool {
	return t.detached
}

// LoadEnvFile Log to the variable env the values
func (t *Task) LoadEnvFile() error {
	if len(t.EnvFile) > 0 {
		info, err := os.Stat(t.EnvFile)
		if os.IsNotExist(err) {
			return err
		}

		if info.IsDir() {
			return fmt.Errorf("log path '%s' is a directory", t.EnvFile)
		}

		file, err := os.Open(t.EnvFile)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		envs := t.Env
		for scanner.Scan() {
			parts := strings.SplitAfter(scanner.Text(), "=")
			t.Env[parts[0]] = parts[1]
		}

		for env, value := range envs {
			t.Env[env] = value
		}
	}

	return nil
}

// OverwriteEnvs Overwrites the env variable in the task
func (t *Task) OverwriteEnvs(envs map[string]string) {
	for env, value := range envs {
		t.Env[env] = value
	}
}

// Get Envs return env variables as string
func (t *Task) GetEnvs() []string {
	var envs []string
	for env, value := range t.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", env, value))
	}

	return envs
}