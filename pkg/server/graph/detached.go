package graph

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/jjzcru/elk/pkg/server/graph/model"
	"time"
)

var detachedTasksMap = make(map[string]*model.DetachedTask)

func getDetachedTaskID() string {
	hash := md5.New()
	hash.Write([]byte(time.Now().Format(time.RFC3339)))
	var id string
	for {
		id = hex.EncodeToString(hash.Sum(nil))
		if _, ok := detachedTasksMap[id]; !ok {
			break
		}
	}

	return id
}
