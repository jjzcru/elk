package graph

import (
	"github.com/jjzcru/elk/pkg/server/graph/model"
	"github.com/jjzcru/elk/pkg/utils"
	"os"
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
