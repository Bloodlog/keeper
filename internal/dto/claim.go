package dto

import (
	jwt "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}
