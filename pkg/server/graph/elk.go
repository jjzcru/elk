package graph

import (
	"context"
	"fmt"
	"github.com/jjzcru/elk/pkg/engine"
	"github.com/jjzcru/elk/pkg/primitives/ox"
	"github.com/jjzcru/elk/pkg/server/graph/model"
	"github.com/jjzcru/elk/pkg/utils"
	"os"
	"sync"
)

func getTask(name string) (*model.Task, error) {
	elk, err := utils.GetElk(os.Getenv("ELK_FILE"), true)
	if err != nil {
		return nil, err
	}

	elkModel, err := mapElk(elk)
	if err != nil {
		return nil, err
	}

	var task *model.Task

	for _, t := range elkModel.Tasks {
		if t.Name == name {
			task = t
		}
	}

	return task, nil
}

func TaskWG(ctx context.Context, cliEngine *engine.Engine, task string, wg *sync.WaitGroup, errChan chan map[string]error) {
	if wg != nil {
		defer wg.Done()
	}

	err := cliEngine.Run(ctx, task)
	if err != nil {
		errChan <- map[string]error{task: err}
	}
}

func loadTaskProperties(elk *ox.Elk, properties *model.TaskProperties) {
	if properties != nil {
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
	}
}
