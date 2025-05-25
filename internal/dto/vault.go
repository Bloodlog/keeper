package dto

import "time"

type AgentCreateSecret struct {
	ExpiredAt   time.Time
	Token       string
	Path        string
	Description string
	Payload     []byte
}

type AgentGetSecret struct {
	ExpiredAt   time.Time
	CreatedAt   time.Time
	DeletedAt   *time.Time
	Path        string
	Description string
	Payload     []byte
	Version     int64
}

type ServerCreateSecret struct {
	ExpiredAt   time.Time
	Path        string
	Description string
	Payload     []byte
	UserID      int64
}

type DecryptedSecretResponse struct {
	ExpiredAt   time.Time
	CreatedAt   time.Time
	DeletedAt   *time.Time
	Path        string
	Description string
	Data        []byte
	Version     int64
}
