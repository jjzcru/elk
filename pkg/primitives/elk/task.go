package elk

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// Task is the data structure for the task to run
type Task struct {
	Cmds           []string          `yaml:"cmds"`
	Env            map[string]string `yaml:"env,omitempty"`
	EnvFile        string            `yaml:"env_file,omitempty"`
	Description    string            `yaml:"description,omitempty"`
	Dir            string            `yaml:"dir,omitempty"`
	Log            string            `yaml:"log,omitempty"`
	Watch          string            `yaml:"watch,omitempty"`
	Deps           []string          `yaml:"deps,omitempty"`
	IgnoreError    bool              `yaml:"ignore_error,omitempty"`
	BackgroundDeps []string          `yaml:"background_deps,omitempty"`
}

// LoadEnvFile Log to the variable env the values
func (t *Task) LoadEnvFile() error {
	if t.Env == nil {
		t.Env = make(map[string]string)
	}

	if len(t.EnvFile) > 0 {
		envFromFile, err := GetEnvFromFile(t.EnvFile)
		if err != nil {
			return err
		}

		envs := make(map[string]string)
		for env, value := range envFromFile {
			envs[env] = value
		}

		for env, value := range t.Env {
			envs[env] = value
		}

		t.Env = envs
	}

	return nil
}

// GetEnvs return env variables as string
func (t *Task) GetEnvs() []string {
	var envs []string
	for env, value := range t.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", env, value))
	}

	return envs
}

// GetWatcherFiles return a list of the files that are going to be watched
func (t *Task) GetWatcherFiles(reg string) ([]string, error) {
	dir := t.Dir
	if len(dir) == 0 {
		d, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		dir = d
	}

	re := regexp.MustCompile(reg)
	var files []string
	walk := func(fn string, fi os.FileInfo, err error) error {
		if re.MatchString(fn) == false {
			return nil
		}
		if fi.IsDir() {
			files = append(files, fn+string(os.PathSeparator))
		} else {
			files = append(files, fn)
		}
		return nil
	}

	err := filepath.Walk(dir, walk)
	if err != nil {
		return files, err
	}

	return files, nil
}
