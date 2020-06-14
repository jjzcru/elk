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
	var startDuration time.Duration
	var delayDuration time.Duration
	var sleepDuration time.Duration

	if start != nil {
		now := time.Now()
		startTime := *start

		startDuration = startTime.Sub(now)
	}

	if delay != nil {
		delayDuration = *delay
	}

	if startDuration > 0 && delayDuration > 0 {
		if startDuration > delayDuration {
			sleepDuration = startDuration
		} else {
			sleepDuration = delayDuration
		}
	} else if startDuration > 0 {
		sleepDuration = startDuration
	} else if delayDuration > 0 {
		sleepDuration = delayDuration
	}

	if sleepDuration > 0 {
		time.Sleep(sleepDuration)
	}
}
