package graph

import (
	"context"
	"errors"
)

func auth(ctx context.Context) error {
	token := ctx.Value(TokenKey).(string)
	authorization := ctx.Value(AuthorizationKey).(string)

	if len(token) > 0 {
		if authorization != token {
			return errors.New("authorization error")
		}
	}

	return nil
}
