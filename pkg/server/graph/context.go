package graph

import (
	"context"

	"github.com/jjzcru/elk/pkg/server/graph/model"
)

// ContextKey is a type for the values that are send on each request
type ContextKey int

const (
	// ElkFileKey stores which is the file path used to run the server
	ElkFileKey ContextKey = iota

	// TokenKey token that is sent by the user request
	TokenKey ContextKey = iota

	// AuthorizationKey stores the valid authorization key
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
