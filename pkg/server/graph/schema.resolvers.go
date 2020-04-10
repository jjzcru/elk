package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

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

	loadTaskProperties(elk, properties)

	errChan := make(chan map[string]error)

	clientEngine := &engine.Engine{
		Elk: elk,
		Executer: engine.DefaultExecuter{
			Logger: logger,
		},
	}

	closeChannels := func() {
		close(outChan)
		close(errTaskChan)
		close(errChan)
	}

	go func() {
		defer closeChannels()
		var wg sync.WaitGroup
		for _, task := range tasks {
			wg.Add(1)
			go TaskWG(ctx, clientEngine, task, &wg, errChan)
		}

		wg.Wait()
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
				for _, taskError := range err {
					return nil, taskError
				}
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

func (r *mutationResolver) RunDetached(ctx context.Context, tasks []string, properties *model.TaskProperties, config *model.RunConfig) (*model.DetachedTask, error) {
	ctx = ServerCtx
	ctx, cancel := context.WithCancel(ctx)
	id := getDetachedTaskID()
	elk, err := utils.GetElk(os.Getenv("ELK_FILE"), true)
	if err != nil {
		return nil, err
	}

	var start *time.Time
	var delay *time.Duration

	loadTaskProperties(elk, properties)

	if config != nil {
		start = config.Start
		delay = config.Delay

		if config.Timeout != nil {
			ctx, _ = context.WithTimeout(ctx, *config.Timeout)
			fmt.Println(config.Delay)
		}

		if config.Deadline != nil {
			ctx, _ = context.WithDeadline(ctx, *config.Deadline)
			fmt.Println(config.Deadline)
		}

	}

	err = elk.Build()
	if err != nil {
		return nil, err
	}

	outputMap := make(map[string]*model.Output)
	var outputs []*model.Output
	for _, task := range tasks {
		output := model.Output{
			Task:  task,
			Out:   []string{},
			Error: []string{},
		}

		outputMap[task] = &output
		outputs = append(outputs, outputMap[task])
	}

	logger, outChan, errTaskChan, err := GraphQLLogger(elk.Tasks)
	if err != nil {
		return nil, err
	}

	errChan := make(chan map[string]error)

	clientEngine := &engine.Engine{
		Elk: elk,
		Executer: engine.DefaultExecuter{
			Logger: logger,
		},
	}

	closeChannels := func() {
		close(outChan)
		close(errTaskChan)
		close(errChan)
	}

	go func() {
		defer closeChannels()
		var wg sync.WaitGroup
		delayStart(delay, start)
		for _, task := range tasks {
			wg.Add(1)
			go TaskWG(ctx, clientEngine, task, &wg, errChan)
		}

		wg.Wait()
	}()

	detachedTasks, err := func() ([]*model.Task, error) {
		var result []*model.Task
		for _, task := range tasks {
			taskModel, err := mapTask(elk.Tasks[task])
			if err != nil {
				return nil, err
			}
			taskModel.Name = task
			result = append(result, taskModel)
		}

		return result, nil
	}()
	if err != nil {
		return nil, err
	}

	response := model.DetachedTask{
		ID:       id,
		Tasks:    detachedTasks,
		Outputs:  outputs,
		Status:   "running",
		Duration: 0,
		StartAt:  time.Now(),
	}

	DetachedTasksMap[id] = &response

	contextMap := detachedContext{
		ctx:    ctx,
		cancel: cancel,
	}

	DetachedCtxMap[id] = &contextMap

	go func() {
		for {
			select {
			case out, ok := <-outChan:
				if !ok {
					outChan = nil
				} else {
					for taskName, value := range out {
						if len(value) > 1 {
							output := outputMap[taskName]
							if output != nil {
								output.Out = append(output.Out, value)
								outputMap[taskName] = output
							}
						}
					}
				}
			case <-ctx.Done():
				outChan = nil
				errTaskChan = nil
				response.Status = "killed"
				break
			case err, ok := <-errTaskChan:
				if !ok {
					errTaskChan = nil
				} else {
					for taskName, value := range err {
						if len(value) > 1 {
							output := outputMap[taskName]
							output.Error = append(output.Error, value)
							outputMap[taskName] = output
						}
					}
				}
			case err, ok := <-errChan:
				if !ok {
					errChan = nil
				} else {
					for taskName, taskError := range err {
						message := taskError.Error()
						response.Status = "error"

						output := outputMap[taskName]
						output.Error = append(output.Error, message)
						outputMap[taskName] = output
					}
				}
			}

			if outChan == nil && errTaskChan == nil {
				if response.Status == "running" {
					response.Status = "success"
				}
				endAt := time.Now()
				response.EndAt = &endAt
				duration := float64(response.EndAt.Sub(response.StartAt)*time.Millisecond) / math.Pow(10, 12)
				response.Duration = int(duration)
				break
			}
		}
	}()

	result := response
	return &result, nil
}

func (r *mutationResolver) Kill(ctx context.Context, id string) (*model.DetachedTask, error) {
	if detachedTask, ok := DetachedTasksMap[id]; ok {
		contextMap := DetachedCtxMap[id]
		if contextMap.ctx.Err() != nil {
			return detachedTask, nil
		}

		endAt := time.Now()
		detachedTask.EndAt = &endAt
		detachedTask.Status = "killed"

		duration := float64(detachedTask.EndAt.Sub(detachedTask.StartAt)*time.Millisecond) / math.Pow(10, 12)
		detachedTask.Duration = int(duration)

		contextMap.cancel()
		DetachedTasksMap[id] = detachedTask
		return detachedTask, nil
	}
	return nil, nil
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

func (r *queryResolver) DetachedTask(ctx context.Context, id string) (*model.DetachedTask, error) {
	return DetachedTasksMap[id], nil
}

func (r *queryResolver) DetachedTasks(ctx context.Context) ([]*model.DetachedTask, error) {
	var detachedTasks []*model.DetachedTask
	for _, v := range DetachedTasksMap {
		detachedTasks = append(detachedTasks, v)
	}
	return detachedTasks, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
