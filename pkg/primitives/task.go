package primitives

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Task is the data structure for the task to run
type Task struct {
	Cmds         []string
	Env          map[string]string
	detached     bool
	Description  string
	Dir          string
	Log          string `yaml:"log"`
	Watch		 string
	Deps         []string
	IgnoreError  bool     `yaml:"ignore_error"`
	DetachedDeps []string `yaml:"detached_deps"`
	EnvFile      string   `yaml:"env_file"`
}

// SetDetached make a task run in attached mode
func (t *Task) SetDetached(detached bool) {
	t.detached = detached
}

// IsDetached Check if the task is detached
func (t *Task) IsDetached() bool {
	return t.detached
}

// LoadEnvFile Log to the variable env the values
func (t *Task) LoadEnvFile() error {
	if t.Env == nil {
		t.Env = make(map[string]string)
	}

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
			parts := strings.SplitAfterN(scanner.Text(), "=", 2)
			env := strings.ReplaceAll(parts[0], "=", "")
			value := parts[1]
			t.Env[env] = value
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
	re := regexp.MustCompile(reg)
	var files []string
	walk := func(fn string, fi os.FileInfo, err error) error {
		if re.MatchString(fn) == false {
			return nil
		}
		if fi.IsDir() {
			// fmt.Println(fn + string(os.PathSeparator))
		} else {
			// fmt.Println(fn)
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
