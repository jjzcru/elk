package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jjzcru/elk/pkg/engine"
	"github.com/jjzcru/elk/pkg/server/graph/generated"
	"github.com/jjzcru/elk/pkg/server/graph/model"
	"github.com/jjzcru/elk/pkg/utils"
)

func (r *mutationResolver) Run(ctx context.Context, tasks []string, properties *model.TaskProperties) ([]*model.Output, error) {
	elk, err := utils.GetElk(os.Getenv("ELK_FILE"), true)
	if err != nil {
		return nil, err
	}

	err = elk.Build()
	if err != nil {
		return nil, err
	}

	outputs := make(map[string]model.Output)
	for _, task := range tasks {
		outputs[task] = model.Output{
			Task:  task,
			Out:   []string{},
			Error: []string{},
		}
	}

	logger, outChan, errTaskChan, err := GraphQLLogger(elk.Tasks)
	if err != nil {
		return nil, err
	}

	for name, task := range elk.Tasks {
		for k, v := range properties.Vars {
			switch v.(type) {
			case string:
				if task.Vars == nil {
					task.Vars = make(map[string]string)
				}
				task.Vars[k] = fmt.Sprintf("%v", v)
			}
		}

		for k, v := range properties.Env {
			switch v.(type) {
			case string:
				if task.Env == nil {
					task.Env = make(map[string]string)
				}
				task.Env[k] = fmt.Sprintf("%v", v)
			}
		}

		task.IgnoreError = *properties.IgnoreError

		elk.Tasks[name] = task
	}

	for name, task := range elk.Tasks {
		for k, v := range properties.Vars {
			switch v.(type) {
			case string:
				if task.Vars == nil {
					task.Vars = make(map[string]string)
				}
				task.Vars[k] = fmt.Sprintf("%v", v)
			}
		}

		for k, v := range properties.Env {
			switch v.(type) {
			case string:
				if task.Env == nil {
					task.Env = make(map[string]string)
				}
				task.Env[k] = fmt.Sprintf("%v", v)
			}
		}

		task.IgnoreError = *properties.IgnoreError

		elk.Tasks[name] = task
	}

	errChan := make(chan error)

	clientEngine := &engine.Engine{
		Elk: elk,
		Executer: engine.DefaultExecuter{
			Logger: logger,
		},
	}

	go func() {
		var wg sync.WaitGroup
		for _, task := range tasks {
			wg.Add(1)
			go TaskWG(ctx, clientEngine, task, &wg, errChan)
		}

		wg.Wait()
		close(outChan)
		close(errTaskChan)
		close(errChan)
	}()

	for {
		select {
		case out, ok := <-outChan:
			if !ok {
				outChan = nil
			} else {
				for taskName, value := range out {
					if len(value) > 1 {
						output := outputs[taskName]
						output.Out = append(output.Out, value)
						outputs[taskName] = output
					}
				}
			}
		case err, ok := <-errTaskChan:
			if !ok {
				errTaskChan = nil
			} else {
				for taskName, value := range err {
					if len(value) > 1 {
						output := outputs[taskName]
						output.Error = append(output.Error, value)
						outputs[taskName] = output
					}
				}
			}
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
			} else {
				return nil, err
			}
		}

		if outChan == nil && errTaskChan == nil {
			break
		}
	}

	var response []*model.Output

	for task := range outputs {
		resp := outputs[task]
		response = append(response, &resp)
	}

	return response, nil
}

func (r *queryResolver) Elk(_ context.Context) (*model.Elk, error) {
	elk, err := utils.GetElk(os.Getenv("ELK_FILE"), true)
	if err != nil {
		return nil, err
	}

	elkModel, err := mapElk(elk)
	if err != nil {
		return nil, err
	}

	return elkModel, nil
}

func (r *queryResolver) Tasks(_ context.Context) ([]*model.Task, error) {
	elk, err := utils.GetElk(os.Getenv("ELK_FILE"), true)
	if err != nil {
		return nil, err
	}

	elkModel, err := mapElk(elk)
	if err != nil {
		return nil, err
	}

	return elkModel.Tasks, nil
}

func (r *queryResolver) Task(_ context.Context, name string) (*model.Task, error) {
	return getTask(name)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
