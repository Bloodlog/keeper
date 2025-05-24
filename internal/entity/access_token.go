package entity

import "time"

type AccessToken struct {
	ExpiresAt time.Time
	CreatedAt time.Time
	Token     string
	ID        int64
	UserID    int64
}
