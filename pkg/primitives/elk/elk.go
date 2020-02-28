package elk

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Elk is the structure of the application
type Elk struct {
	Version     string
	Env         map[string]string
	EnvFile     string `yaml:"env_file"`
	IgnoreError bool   `yaml:"ignore_error"`
	Tasks       map[string]Task
}

// GetTask Get a task object by its name
func (e *Elk) GetTask(name string) (*Task, error) {
	err := e.HasCircularDependency(name)
	if err != nil {
		return nil, err
	}

	task := e.Tasks[name]
	return &task, nil
}

// HasTask return a boolean if the incoming event exist
func (e *Elk) HasTask(name string) bool {
	if _, ok := e.Tasks[name]; ok {
		return true
	}
	return false
}

// LoadEnvFile Log to the variable env the values
func (e *Elk) LoadEnvFile() error {
	if e.Env == nil {
		e.Env = make(map[string]string)
	}

	envCopy := make(map[string]string)
	for k, v := range e.Env {
		envCopy[k] = v
	}

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

		envs := e.Env
		for scanner.Scan() {
			parts := strings.SplitAfterN(scanner.Text(), "=", 2)
			env := strings.ReplaceAll(parts[0], "=", "")
			value := parts[1]
			e.Env[env] = value
		}

		for env, value := range envs {
			e.Env[env] = value
		}

		for env, value := range envCopy {
			e.Env[env] = value
		}

		for _, task := range e.Tasks {
			if task.Env == nil {
				task.Env = make(map[string]string)
			}

			taskEnvs := task.Env
			for env, value := range e.Env {
				task.Env[env] = value
			}

			for env, value := range taskEnvs {
				task.Env[env] = value
			}

			err = task.LoadEnvFile()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// LoadEnvsInTasks Load env variable from elk to its tasks
func (e *Elk) LoadEnvsInTasks() {
	for _, task := range e.Tasks {
		envs := make(map[string]string)
		for k, v := range e.Env {
			envs[k] = v
		}

		for k, v := range task.Env {
			envs[k] = v
		}

		task.OverwriteEnvs(envs)
	}

}

// OverwriteEnvs Overwrites the env variable in the elk file and all the tasks
func (e *Elk) OverwriteEnvs(envs map[string]string) {
	if e.Env == nil {
		e.Env = make(map[string]string)
	}

	for env, value := range envs {
		e.Env[env] = value
	}

	for name, task := range e.Tasks {
		task.OverwriteEnvs(envs)
		e.Tasks[name] = task
	}
}

// HasCircularDependency checks if a task has a circular dependency
func (e *Elk) HasCircularDependency(name string, visitedNodes ...string) error {
	if !e.HasTask(name) {
		return ErrTaskNotFound
	}

	task := e.Tasks[name]

	if len(append(task.Deps, task.DetachedDeps...)) == 0 {
		return nil
	}

	dependencyGraph, err := e.getDependencyGraph(&task)
	if err != nil {
		return err
	}

	for _, node := range visitedNodes {
		if node == name {
			// return fmt.Errorf("the task '%s' has a circular dependency", name)
			return ErrCircularDependency
		}
	}

	visitedNodes = append(visitedNodes, name)

	for _, dep := range dependencyGraph {
		for _, d := range dep {
			err = e.HasCircularDependency(d, visitedNodes...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Elk) getDependencyGraph(task *Task) (map[string][]string, error) {
	dependencyGraph := make(map[string][]string)
	deps := append(task.Deps, task.DetachedDeps...)
	for _, dep := range deps {
		// Validate that the dependency is a valid task
		t, exists := e.Tasks[dep]
		if exists == false {
			return dependencyGraph, fmt.Errorf("The dependency '%s' do not exist as a task", dep)
		}

		dependencyGraph[dep] = append(t.Deps, t.DetachedDeps...)
	}
	return dependencyGraph, nil
}

func FromFile(filePath string) (*Elk, error) {
	elk := Elk{}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("path do not exist: '%s'", filePath)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &elk)
	if err != nil {
		return nil, err
	}

	if elk.Tasks == nil {
		elk.Tasks = make(map[string]Task)
	}

	if elk.Env == nil {
		elk.Env = make(map[string]string)
	}

	for name := range elk.Tasks {
		task := elk.Tasks[name]

		if task.Env == nil {
			task.Env = make(map[string]string)
		}

		elk.Tasks[name] = task
	}

	return &elk, nil
}
