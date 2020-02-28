package run

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/engine"
	"github.com/jjzcru/elk/pkg/primitives/elk"
)

func runWatch(cliEngine *engine.Engine, taskCtx context.Context, task string, t *elk.Task, cancel context.CancelFunc, ctx context.Context) {
	go func() {
		err := cliEngine.Run(taskCtx, task)
		if err != nil {
			utils.PrintError(err)
			return
		}
	}()

	files, err := t.GetWatcherFiles(t.Watch)
	if err != nil {
		utils.PrintError(err)
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.PrintError(err)
		return
	}
	defer watcher.Close()

	for _, file := range files {
		err = watcher.Add(file)
		if err != nil {
			utils.PrintError(err)
			return
		}
	}

	for {
		select {
		case event := <-watcher.Events:
			switch {
			case event.Op&fsnotify.Write == fsnotify.Write:
				fallthrough
			case event.Op&fsnotify.Create == fsnotify.Create:
				fallthrough
			case event.Op&fsnotify.Remove == fsnotify.Remove:
				fallthrough
			case event.Op&fsnotify.Rename == fsnotify.Rename:
				go func() {
					cancel()
					taskCtx, cancel = context.WithCancel(ctx)

					err = cliEngine.Run(taskCtx, task)
					if err != nil && err != context.Canceled {
						utils.PrintError(err)
						return
					}
				}()
			}
		case err := <-watcher.Errors:
			utils.PrintError(err)
			return
		}
	}
}
