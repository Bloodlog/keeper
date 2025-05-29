package entity

import "time"

type SecretMetadata struct {
	ExpiredAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Path        string
	Description string
	UserID      int64
}

type SecretVersion struct {
	ExpiredAt  time.Time
	CreatedAt  time.Time
	DeletedAt  *time.Time
	FilePath   *string
	Value      []byte
	MetadataID int64
	Version    int64
	Destroyed  bool
}

type OneSecretVersionWithMetadata struct {
	ExpiredAt   time.Time
	CreatedAt   time.Time
	DeletedAt   *time.Time
	FilePath    *string
	Path        string
	Description string
	Value       []byte
	MetadataID  int64
	Version     int64
	Destroyed   bool
}
