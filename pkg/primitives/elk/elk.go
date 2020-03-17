package elk

import (
	"fmt"
	"github.com/jjzcru/elk/pkg/file"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

// Elk is the structure of the application
type Elk struct {
	Version string
	Env     map[string]string
	EnvFile string `yaml:"env_file"`
	Tasks   map[string]Task
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

// Build compiles the elk structure and validates its integrity
func (e *Elk) Build() error {
	osEnvs := make(map[string]string)
	for _, en := range os.Environ() {
		parts := strings.SplitAfterN(en, "=", 2)
		env := strings.ReplaceAll(parts[0], "=", "")
		value := parts[1]
		osEnvs[env] = value
	}

	err := e.LoadEnvFile()
	if err != nil {
		return err
	}

	for env, value := range e.Env {
		osEnvs[env] = value
	}

	e.Env = osEnvs

	for name, task := range e.Tasks {
		err = e.HasCircularDependency(name)
		if err != nil {
			return err
		}

		err = task.LoadEnvFile()
		if err != nil {
			return err
		}

		envs := e.Env
		for env, value := range task.Env {
			envs[env] = value
		}
		task.Env = envs

		e.Tasks[name] = task
	}

	return nil
}

// LoadEnvFile Log to the variable env the values
func (e *Elk) LoadEnvFile() error {
	if e.Env == nil {
		e.Env = make(map[string]string)
	}

	if len(e.EnvFile) > 0 {
		envFromFile, err := file.GetEnvFromFile(e.EnvFile)
		if err != nil {
			return err
		}

		envs := make(map[string]string)
		for env, value := range envFromFile {
			envs[env] = value
		}

		for env, value := range e.Env {
			envs[env] = value
		}

		e.Env = envs
	}

	return nil
}

// HasCircularDependency checks if a task has a circular dependency
func (e *Elk) HasCircularDependency(name string, visitedNodes ...string) error {
	if !e.HasTask(name) {
		return ErrTaskNotFound
	}

	task := e.Tasks[name]

	if len(task.Deps) == 0 {
		return nil
	}

	dependencyGraph, err := e.getDependencyGraph(&task)
	if err != nil {
		return err
	}

	for _, node := range visitedNodes {
		if node == name {
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
	deps := task.Deps
	for _, dep := range deps {
		// Validate that the dependency is a valid task
		t, exists := e.Tasks[dep.Name]
		if exists == false {
			return dependencyGraph, fmt.Errorf("the dependency '%s' do not exist as a task", dep.Name)
		}

		var depsNames []string
		for _, d := range t.Deps {
			depsNames = append(depsNames, d.Name)
		}
		dependencyGraph[dep.Name] = depsNames
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
