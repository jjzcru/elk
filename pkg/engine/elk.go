package engine

// Elk is the structure of the application
type Elk struct {
	Version string
	Env     map[string]string
	EnvFile string `yaml:"env_file"`
	Tasks   map[string]Task
}

// Task is the data structure for the task to run
type Task struct {
	Cmds         []string
	Env          map[string]string
	Description  string
	Dir          string
	Log          string
	Deps         []string
	DetachedDeps []string `yaml:"detached_deps"`
	EnvFile      string   `yaml:"env_file"`
}

// Config is the structure of the global configurations
type Config struct {
	Path string
}
