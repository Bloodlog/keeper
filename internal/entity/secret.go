package entity

import "time"

type Secret struct {
	ExpiredAt   time.Time
	CreatedAt   time.Time
	DeletedAt   *time.Time
	Path        string
	Description string
	Value       []byte
	UserID      int64
	Version     int64
}
