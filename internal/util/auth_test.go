package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGetUserID(t *testing.T) {
	ctx := t.Context()

	var userID int64 = 123
	ctx = SetUserID(ctx, userID)

	retrievedUserID, err := GetUserID(ctx)

	assert.NoError(t, err)
	assert.Equal(t, userID, retrievedUserID)
}

func TestGetUserID_NoUserID(t *testing.T) {
	ctx := t.Context()

	_, err := GetUserID(ctx)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized", err.Error())
}
