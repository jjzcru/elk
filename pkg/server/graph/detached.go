package graph

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"sync"
	"time"

	"github.com/jjzcru/elk/pkg/server/graph/model"
)

type detachedContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// ServerCtx stores the context on which the server is running
var ServerCtx context.Context

// DetachedTasksMap links detached task with an id
var DetachedTasksMap = make(map[string]*model.DetachedTask)

// DetachedCtxMap stores the context of each task using an id
var DetachedCtxMap = make(map[string]*detachedContext)

func getDetachedTaskID() string {
	hash := md5.New()
	_, _ = hash.Write([]byte(time.Now().Format(time.RFC3339)))
	var id string
	for {
		id = hex.EncodeToString(hash.Sum(nil))
		if _, ok := DetachedTasksMap[id]; !ok {
			break
		}
	}

	return id
}

func getResponseFromDetached(id string) *model.DetachedTask {
	return DetachedTasksMap[id]
}

func updateDetachedTask(id string, task *model.DetachedTask) {
	DetachedTasksMap[id] = task
}

// CancelDetachedTasks call cancel on all the context
func CancelDetachedTasks() {
	var wg sync.WaitGroup
	for id := range DetachedCtxMap {
		detachedTask := DetachedTasksMap[id]
		if detachedTask == nil {
			continue
		}

		switch detachedTask.Status {
		case "running":
			break
		default:
			continue
		}

		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			detachedCtx := DetachedCtxMap[id]
			if detachedCtx != nil {
				detachedCtx.cancel()
			}
		}(id)
	}
	wg.Wait()
}

func delayStart(delay *time.Duration, start *time.Time) {
	sleepDuration := getDelayDuration(delay, start)
	if sleepDuration > 0 {
		time.Sleep(sleepDuration)
	}
}

func getDelayDuration(delay *time.Duration, start *time.Time) time.Duration {
	var startDuration time.Duration
	var delayDuration time.Duration

	if start != nil {
		now := time.Now()
		var startTime time.Time

		if start.Before(now) {
			startTime = time.Date(now.Year(), now.Month(), now.Day(),
				start.Hour(), start.Minute(), start.Second(),
				start.Nanosecond(), start.Location())

			startTime.Add(24 * time.Hour)
		} else {
			startTime = *start
		}
		startDuration = startTime.Sub(now)
	}

	if delay != nil {
		delayDuration = *delay
	}

	if startDuration > 0 && delayDuration > 0 {
		if startDuration > delayDuration {
			return startDuration
		}
		return delayDuration
	} else if startDuration > 0 {
		return startDuration
	} else if delayDuration > 0 {
		return delayDuration
	}

	return 0
}
