package ox

import (
	"fmt"
	"github.com/jjzcru/elk/pkg/file"
)

// Task is the data structure for the task to run
type Task struct {
	Cmds        []string          `yaml:"cmds"`
	Env         map[string]string `yaml:"env,omitempty"`
	Vars        map[string]string `yaml:"vars,omitempty"`
	EnvFile     string            `yaml:"env_file,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Dir         string            `yaml:"dir,omitempty"`
	Log         Log               `yaml:"log,omitempty"`
	Sources     string            `yaml:"sources,omitempty"`
	Deps        []Dep             `yaml:"deps,omitempty"`
	IgnoreError bool              `yaml:"ignore_error,omitempty"`
}

type Dep struct {
	Name     string `yaml:"name"`
	Detached bool   `yaml:"detached"`
}

type Log struct {
	Out string `yaml:"out"`
	Err string `yaml:"error"`
}

// LoadEnvFile Log to the variable env the values
func (t *Task) LoadEnvFile() error {
	if t.Env == nil {
		t.Env = make(map[string]string)
	}

	if len(t.EnvFile) > 0 {
		envFromFile, err := file.GetEnvFromFile(t.EnvFile)
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
