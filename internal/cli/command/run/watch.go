package run

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/jjzcru/elk/internal/cli/utils"
	"github.com/jjzcru/elk/pkg/engine"
	"github.com/jjzcru/elk/pkg/primitives/elk"
	"os"
	"path/filepath"
	"regexp"
)

func runWatch(ctx context.Context, cliEngine *engine.Engine, task string, t elk.Task) {
	taskCtx, cancel := context.WithCancel(ctx)

	go func() {
		err := cliEngine.Run(taskCtx, task)
		if err != nil {
			utils.PrintError(err)
			return
		}
	}()

	files, err := getWatcherFiles(t.Watch, t.Dir)
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
		fmt.Println(file)
		if err != nil {

			fmt.Println(err.Error())
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
				cancel()
			}
		case <-taskCtx.Done():
			taskCtx, cancel = context.WithCancel(ctx)
			go func() {
				err = cliEngine.Run(taskCtx, task)
				if err != nil {
					utils.PrintError(err)
				}
			}()
		case err := <-watcher.Errors:
			utils.PrintError(err)
			return
		}
	}
}

func getWatcherFiles(reg string, dir string) ([]string, error) {
	if len(dir) == 0 {
		d, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		dir = d
	}

	re := regexp.MustCompile(reg)
	var files []string
	walk := func(fn string, fi os.FileInfo, err error) error {
		if re.MatchString(fn) == false {
			return nil
		}
		if !fi.IsDir() {
			files = append(files, fn)
		}
		return nil
	}

	err := filepath.Walk(dir, walk)
	if err != nil {
		return files, err
	}

	return files, nil
}
