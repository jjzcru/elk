package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"sync"
	"time"

	"github.com/jjzcru/elk/pkg/engine"
	"github.com/jjzcru/elk/pkg/server/graph/generated"
	"github.com/jjzcru/elk/pkg/server/graph/model"
	"github.com/jjzcru/elk/pkg/utils"
)

func (r *mutationResolver) Run(ctx context.Context, tasks []string, properties *model.TaskProperties) ([]*model.Output, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}

	elkFilePath := ctx.Value(ElkFileKey).(string)
	elk, err := utils.GetElk(elkFilePath, true)
	if err != nil {
		return nil, err
	}

	if properties != nil && properties.EnvFile != nil {
		if len(*properties.EnvFile) > 0 {
			elk.EnvFile = *properties.EnvFile
		}
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

	logger, outChan, errTaskChan, err := gqlLogger(elk.Tasks, tasks)
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

func (r *mutationResolver) Detached(ctx context.Context, tasks []string, properties *model.TaskProperties, config *model.RunConfig) (*model.DetachedTask, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}

	elkFilePath := ctx.Value(ElkFileKey).(string)

	id := getDetachedTaskID()
	elk, err := utils.GetElk(elkFilePath, true)
	if err != nil {
		return nil, err
	}

	var start *time.Time
	var delay *time.Duration

	loadTaskProperties(elk, properties)

	err = elk.Build()
	if err != nil {
		return nil, err
	}

	isInFuture := func(start *time.Time) bool {
		now := time.Now()
		return start.After(now)
	}

	if config != nil {
		delay = config.Delay

		if isInFuture(config.Start) {
			start = config.Start
		}
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

	logger, outChan, errTaskChan, err := gqlLogger(elk.Tasks, tasks)
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

	ctx, cancel := getConfigContext(ServerCtx, config)

	go func(id string) {
		contextMap := detachedContext{
			ctx:    ctx,
			cancel: cancel,
		}

		DetachedCtxMap[id] = &contextMap

		defer closeChannels()
		var wg sync.WaitGroup
		delayStart(delay, start)

		resp := getResponseFromDetached(id)
		resp.Status = "running"
		updateDetachedTask(id, resp)

		for _, task := range tasks {
			wg.Add(1)
			go TaskWG(ctx, clientEngine, task, &wg, errChan)
		}

		wg.Wait()
	}(id)

	detachedTasks, err := func() ([]*model.Task, error) {
		var result []*model.Task
		for _, task := range tasks {
			taskModel, err := mapTask(elk.Tasks[task], task)
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

	delayDuration := getDelayDuration(delay, start)

	status := "running"
	if delayDuration > 0 {
		status = "waiting"
	}

	response := model.DetachedTask{
		ID:       id,
		Tasks:    detachedTasks,
		Outputs:  outputs,
		Status:   status,
		Duration: 0,
		StartAt:  time.Now().Add(delayDuration),
	}

	DetachedTasksMap[id] = &response

	go func(id string) {
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
				resp := getResponseFromDetached(id)
				resp.Status = "killed"
				updateDetachedTask(id, resp)
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
						resp := getResponseFromDetached(id)
						resp.Status = "error"
						updateDetachedTask(id, resp)

						output := outputMap[taskName]
						output.Error = append(output.Error, message)
						outputMap[taskName] = output
					}
				}
			}

			if outChan == nil && errTaskChan == nil {
				resp := getResponseFromDetached(id)
				if resp.Status == "running" {
					resp.Status = "success"
				}
				endAt := time.Now()
				resp.EndAt = &endAt
				duration := endAt.Sub(response.StartAt)
				resp.Duration = duration
				updateDetachedTask(id, resp)
				break
			}
		}
	}(id)

	result := response
	return &result, nil
}

func (r *mutationResolver) Kill(ctx context.Context, id string) (*model.DetachedTask, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}

	if detachedTask, ok := DetachedTasksMap[id]; ok {
		contextMap := DetachedCtxMap[id]
		if contextMap.ctx.Err() != nil {
			return detachedTask, nil
		}

		endAt := time.Now()
		detachedTask.EndAt = &endAt
		detachedTask.Status = "killed"

		duration := endAt.Sub(detachedTask.StartAt)
		detachedTask.Duration = duration

		contextMap.cancel()
		DetachedTasksMap[id] = detachedTask
		return detachedTask, nil
	}
	return nil, nil
}

func (r *mutationResolver) Remove(ctx context.Context, name string) (*model.Task, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}

	elkFilePath := ctx.Value(ElkFileKey).(string)

	elk, err := utils.GetElk(elkFilePath, true)
	if err != nil {
		return nil, err
	}

	task, err := elk.GetTask(name)
	if err != nil {
		return nil, err
	}

	task.Title = name
	taskModel, err := mapTask(*task, name)
	if err != nil {
		return nil, err
	}

	delete(elk.Tasks, name)

	return taskModel, utils.SetElk(elk, elkFilePath)
}

func (r *mutationResolver) Put(ctx context.Context, task model.TaskInput) (*model.Task, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}

	elkFilePath := ctx.Value(ElkFileKey).(string)

	elk, err := utils.GetElk(elkFilePath, true)
	if err != nil {
		return nil, err
	}

	t := mapTaskInput(task)

	if _, exist := elk.Tasks[task.Name]; exist {
		t = mergeTaskInput(task, t)
	}

	elk.Tasks[task.Name] = t
	taskModel, err := mapTask(t, task.Name)
	if err != nil {
		return nil, err
	}

	return taskModel, utils.SetElk(elk, elkFilePath)
}

func (r *queryResolver) Health(ctx context.Context) (bool, error) {
	return true, nil
}

func (r *queryResolver) Elk(ctx context.Context) (*model.Elk, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}

	elkFilePath := ctx.Value(ElkFileKey).(string)
	elk, err := utils.GetElk(elkFilePath, true)
	if err != nil {
		return nil, err
	}

	elkModel, err := mapElk(elk)
	if err != nil {
		return nil, err
	}

	return elkModel, nil
}

func (r *queryResolver) Tasks(ctx context.Context, name *string) ([]*model.Task, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}

	elkFilePath := ctx.Value(ElkFileKey).(string)
	elk, err := utils.GetElk(elkFilePath, true)
	if err != nil {
		return nil, err
	}

	elkModel, err := mapElk(elk)
	if err != nil {
		return nil, err
	}

	var tasks []*model.Task
	if name != nil {
		for _, task := range elkModel.Tasks {
			if task != nil && task.Name == *name {
				tasks = append(tasks, task)
			}
		}
	} else {
		tasks = elkModel.Tasks
	}

	return tasks, nil
}

func (r *queryResolver) Detached(ctx context.Context, ids []string, status []model.DetachedTaskStatus) ([]*model.DetachedTask, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}

	detachedTaskIDs := getDetachedTaskIDs()

	// First filter by status
	if len(status) > 0 {
		detachedTaskIDs = getDetachedTasksByStatus(status)
	}

	// Filter by id
	if ids != nil {
		detachedTaskIDs = getDetachedTasksByID(ids, detachedTaskIDs)
	}

	return getDetachedTaskFromIDs(detachedTaskIDs), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
