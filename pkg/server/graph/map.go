package graph

import (
	"github.com/jjzcru/elk/pkg/primitives/ox"
	"github.com/jjzcru/elk/pkg/server/graph/model"
)

func mapElk(elk *ox.Elk) (*model.Elk, error) {
	elkModel := model.Elk{
		Version: elk.Version,
		Env:     map[string]interface{}{},
		Vars:    map[string]interface{}{},
		Tasks:   []*model.Task{},
	}

	for k, v := range elk.Env {
		elkModel.Env[k] = v
	}

	for k, v := range elk.Vars {
		elkModel.Vars[k] = v
	}

	for k, v := range elk.Tasks {
		task, err := mapTask(v)
		if err != nil {
			return nil, err
		}
		task.Name = k
		elkModel.Tasks = append(elkModel.Tasks, task)
	}

	return &elkModel, nil
}

func mapTask(task ox.Task) (*model.Task, error) {
	taskModel := model.Task{
		Cmds:        []*string{},
		Env:         map[string]interface{}{},
		Vars:        map[string]interface{}{},
		EnvFile:     task.EnvFile,
		Description: task.Description,
		Dir:         task.Dir,
		Log: &(model.Log{
			Out:    task.Log.Out,
			Format: task.Log.Format,
			Error:  task.Log.Err,
		}),
		Sources:     &task.Sources,
		Deps:        []*model.Dep{},
		IgnoreError: task.IgnoreError,
	}

	for _, cmd := range task.Cmds {
		taskModel.Cmds = append(taskModel.Cmds, &cmd)
	}

	for k, v := range task.Env {
		taskModel.Env[k] = v
	}

	for k, v := range task.Vars {
		taskModel.Vars[k] = v
	}

	for _, dep := range task.Deps {
		taskModel.Deps = append(taskModel.Deps, mapDep(dep))
	}

	return &taskModel, nil
}

func mapDep(dep ox.Dep) *model.Dep {
	depModel := model.Dep{
		Name:     dep.Name,
		Detached: dep.Detached,
	}

	return &depModel
}
