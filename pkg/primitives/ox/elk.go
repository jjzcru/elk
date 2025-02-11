package ox

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jjzcru/elk/pkg/file"
	"github.com/jjzcru/elk/pkg/maps"
	"gopkg.in/yaml.v2"
)

// Elk is the structure of the application
type Elk struct {
	filePath string
	Version  string
	Env      map[string]string `yaml:"env"`
	Vars     map[string]string `yaml:"vars"`
	EnvFile  string            `yaml:"env_file"`
	Tasks    map[string]Task
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

// GetFilePath get path used to create the object
func (e *Elk) GetFilePath() string {
	return e.filePath
}

// SetFilePath set the path used to create the object
func (e *Elk) SetFilePath(filepath string) {
	e.filePath = filepath
}

// Build compiles the ox structure and validates its integrity
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

		task.Env = maps.MergeMaps(maps.CopyMap(e.Env), maps.CopyMap(task.Env))
		task.Vars = maps.MergeMaps(maps.CopyMap(e.Vars), maps.CopyMap(task.Vars))

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
		if !exists {
			return dependencyGraph, ErrTaskNotFound
		}

		var depsNames []string
		for _, d := range t.Deps {
			depsNames = append(depsNames, d.Name)
		}
		dependencyGraph[dep.Name] = depsNames
	}
	return dependencyGraph, nil
}

// FromFile loads an elk object from a file
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

// ToFile saves an elk object to a file
func ToFile(elk *Elk, filePath string) error {
	dataBytes, err := yaml.Marshal(elk)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(
		filePath,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(dataBytes)
	if err != nil {
		return err
	}

	return nil
}
