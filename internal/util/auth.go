package util

import (
	"context"
	"errors"
)

type contextKey string

const userIDKey contextKey = "userID"

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(userIDKey).(int64)
	if !ok {
		return 0, errors.New("unauthorized")
	}
	return userID, nil
}
