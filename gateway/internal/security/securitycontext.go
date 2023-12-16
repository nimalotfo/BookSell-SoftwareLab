package security

import (
	"context"
	"errors"
)

type contextKey struct{}

var userContextKey contextKey

var (
	ErrInvalidUserContext = errors.New("user context not found")
)

func NewUserContext(ctx context.Context, userId int64) context.Context {
	return context.WithValue(ctx, userContextKey, userId)
}

// func GetCurrentUser(ctx context.Context) (interface{}, error){}

func GetCurrentUserId(ctx context.Context) (int64, error) {
	userId, ok := ctx.Value(userContextKey).(int64)
	if !ok {
		return 0, ErrInvalidUserContext
	}

	return userId, nil
}
