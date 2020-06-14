package graph

import (
	"context"
	"fmt"
	"sync"

	"github.com/jjzcru/elk/pkg/engine"
	"github.com/jjzcru/elk/pkg/primitives/ox"
	"github.com/jjzcru/elk/pkg/server/graph/model"
)

// TaskWG run a working group of tasks
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

			if properties.IgnoreError != nil {
				task.IgnoreError = *properties.IgnoreError
			}

			elk.Tasks[name] = task
		}
	}
}
