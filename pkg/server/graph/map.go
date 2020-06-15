package graph

import (
	"fmt"

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
		task, err := mapTask(v, k)
		if err != nil {
			return nil, err
		}
		task.Name = k

		if len(task.Title) == 0 {
			task.Title = task.Name
		}

		elkModel.Tasks = append(elkModel.Tasks, task)
	}

	return &elkModel, nil
}

func mapTask(task ox.Task, name string) (*model.Task, error) {
	taskModel := model.Task{
		Title:       task.Title,
		Name:        name,
		Tags:        uniqueString(task.Tags),
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

	for i := range task.Cmds {
		cmd := task.Cmds[i]
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

func uniqueString(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func mapDep(dep ox.Dep) *model.Dep {
	depModel := model.Dep{
		Name:     dep.Name,
		Detached: dep.Detached,
	}

	return &depModel
}

func mapTaskInput(task model.TaskInput) ox.Task {
	env := make(map[string]string)
	vars := make(map[string]string)

	var deps []ox.Dep
	var log ox.Log

	title := ""
	envFile := ""
	description := ""
	dir := ""
	sources := ""

	ignoreError := false

	if task.Env != nil {
		for k, v := range task.Env {
			env[k] = fmt.Sprintf("%v", v)
		}
	}

	if task.Vars != nil {
		for k, v := range task.Vars {
			vars[k] = fmt.Sprintf("%v", v)
		}
	}

	if task.Title != nil {
		title = *task.Title
	}

	if task.EnvFile != nil {
		envFile = *task.EnvFile
	}

	if task.Description != nil {
		description = *task.Description
	}

	if task.Dir != nil {
		dir = *task.Description
	}

	if task.Sources != nil {
		sources = *task.Sources
	}

	if task.IgnoreError != nil {
		ignoreError = *task.IgnoreError
	}

	if task.Deps != nil {
		for _, dep := range task.Deps {
			deps = append(deps, ox.Dep{
				Name:        dep.Name,
				Detached:    dep.Detached,
				IgnoreError: dep.IgnoreError,
			})
		}
	}

	if task.Log != nil {
		logFormat := ""

		if task.Log.Format != nil {
			logFormat = task.Log.Format.String()
		}

		log = ox.Log{
			Out:    task.Log.Out,
			Err:    task.Log.Error,
			Format: logFormat,
		}
	}

	return ox.Task{
		Title:       title,
		Tags:        task.Tags,
		Cmds:        task.Cmds,
		Env:         env,
		Vars:        vars,
		EnvFile:     envFile,
		Description: description,
		Dir:         dir,
		Sources:     sources,
		IgnoreError: ignoreError,
		Log:         log,
		Deps:        deps,
	}
}
