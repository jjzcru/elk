package graph

import (
	"context"
	"github.com/jjzcru/elk/pkg/server/graph/model"
)

type ContextKey int

const (
	ElkFileKey ContextKey = iota
	TokenKey ContextKey = iota
	AuthorizationKey ContextKey = iota
)

func getConfigContext(parentContext context.Context, config *model.RunConfig) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parentContext)

	if config != nil {
		if config.Timeout != nil {
			ctx, cancel = context.WithTimeout(ctx, *config.Timeout)
		}

		if config.Deadline != nil {
			ctx, cancel = context.WithDeadline(ctx, *config.Deadline)
		}
	}

	return ctx, cancel
}