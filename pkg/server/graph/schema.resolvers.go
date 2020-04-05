package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"os"
	"sync"

	"github.com/jjzcru/elk/pkg/engine"
	"github.com/jjzcru/elk/pkg/server/graph/generated"
	"github.com/jjzcru/elk/pkg/server/graph/model"
	"github.com/jjzcru/elk/pkg/utils"
)

func (r *mutationResolver) Run(ctx context.Context, tasks []string, detached *bool) ([]*model.Output, error) {
	elk, err := utils.GetElk(os.Getenv("ELK_FILE"), true)
	if err != nil {
		return nil, err
	}

	ouputs := make(map[string]model.Output)

	for _, task := range tasks {
		ouputs[task] = model.Output{
			Task:  task,
			Out:   []string{},
			Error: []string{},
		}
	}

	logger, outChan, errTaskChan := GraphQLLogger(tasks)

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
						output := ouputs[taskName]
						output.Out = append(output.Out, value)
						ouputs[taskName] = output
					}
				}
			}
		case err, ok := <-errTaskChan:
			if !ok {
				errTaskChan = nil
			} else {
				for taskName, value := range err {
					if len(value) > 1 {
						output := ouputs[taskName]
						output.Error = append(output.Error, value)
						ouputs[taskName] = output
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

	for task := range ouputs {
		resp := ouputs[task]
		response = append(response, &resp)
	}

	return response, nil
}

func (r *queryResolver) Elk(ctx context.Context) (*model.Elk, error) {
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

func (r *queryResolver) Tasks(ctx context.Context) ([]*model.Task, error) {
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

func (r *queryResolver) Task(ctx context.Context, name string) (*model.Task, error) {
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
func TaskWG(ctx context.Context, cliEngine *engine.Engine, task string, wg *sync.WaitGroup, errChan chan error) {
	if wg != nil {
		defer wg.Done()
	}

	err := cliEngine.Run(ctx, task)
	if err != nil {
		errChan <- err
	}
}
