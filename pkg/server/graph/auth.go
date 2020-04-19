package graph

import (
	"context"
	"errors"
)

func auth(ctx context.Context) error {
	token := ctx.Value("token").(string)
	authorization := ctx.Value("authorization").(string)

	if len(token) > 0 {
		if authorization != token {
			return errors.New("authorization error")
		}
	}

	return nil
}
