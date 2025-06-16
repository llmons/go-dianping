package user_holder

import (
	"context"
	"fmt"
	"go-dianping/api"
)

type ctxKey struct{}

var userKey ctxKey

func WithUser(ctx context.Context, user *api.SimpleUser) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func GetUser(ctx context.Context) *api.SimpleUser {
	fmt.Println(ctx.Value(userKey))

	if user, ok := ctx.Value(userKey).(*api.SimpleUser); ok {
		return user
	}
	return nil
}
