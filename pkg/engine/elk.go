package engine

// Elk is the structure of the application
type Elk struct {
	Version string
	Env     map[string]string
	Tasks   map[string]Task
}

// Task is the data structure for the task to run
type Task struct {
	Cmds        []string
	Env         map[string]string
	Description string
	Dir         string
}

// Config is the structure of the global configurations
type Config struct {
	Path string
}
