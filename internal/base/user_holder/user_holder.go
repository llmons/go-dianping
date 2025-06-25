package user_holder

import (
	"context"
	"go-dianping/api/v1"
)

type ctxKey struct{}

var userKey ctxKey

func WithUser(ctx context.Context, user *v1.SimpleUser) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func GetUser(ctx context.Context) *v1.SimpleUser {
	if user, ok := ctx.Value(userKey).(*v1.SimpleUser); ok {
		return user
	}
	return nil
}
