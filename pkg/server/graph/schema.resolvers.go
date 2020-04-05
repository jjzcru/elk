package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/jjzcru/elk/pkg/engine"
	"os"

	"github.com/jjzcru/elk/pkg/server/graph/generated"
	"github.com/jjzcru/elk/pkg/server/graph/model"
	"github.com/jjzcru/elk/pkg/utils"
)

func (r *mutationResolver) Run(ctx context.Context, task string, detached *bool) ([]*string, error) {
	elk, err := utils.GetElk(os.Getenv("ELK_FILE"), true)
	if err != nil {
		return nil, err
	}

	logger, outChan, errChan := GraphQLLogger()

	taskModel, err := getTask(task)
	if err != nil {
		return nil, err
	}

	if taskModel == nil {
		return nil, nil
	}

	clientEngine := &engine.Engine{
		Elk: elk,
		Executer: engine.DefaultExecuter{
			Logger: map[string]engine.Logger{
				task: logger,
			},
		},
	}

	go func() {
		err = clientEngine.Run(ctx, task)
		if err != nil {
			errChan <- err.Error()
		}
		close(outChan)
		close(errChan)
	}()

	var response []*string
	for {
		select {
		case out, ok := <-outChan:
			if len(out) > 1 {
				response = append(response, &out)
			}

			if !ok {
				outChan = nil
			}
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
			} else {
				return nil, fmt.Errorf(err)
			}
		}

		if outChan == nil && errChan == nil {
			break
		}
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

func (r *queryResolver) Task(ctx context.Context, name string) (*model.Task, error) {
	return getTask(name)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
