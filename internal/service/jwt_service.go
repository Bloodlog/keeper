package service

import (
	"fmt"
	"keeper/internal/config"
	"keeper/internal/dto"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type JwtService interface {
	CreateJwt(userID int64) (string, error)
	GetUserID(tokenString string) (int64, error)
}

type jwtService struct {
	Cfg config.SecurityConfig
}

func NewJwtService(cfg config.SecurityConfig) JwtService {
	return &jwtService{
		Cfg: cfg,
	}
}

func (o *jwtService) CreateJwt(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, dto.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(o.Cfg.TokenTTL)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(o.Cfg.EncryptionKey))
	if err != nil {
		return "", fmt.Errorf("failed to SignedString: %w", err)
	}

	return tokenString, nil
}

func (o *jwtService) GetUserID(tokenString string) (int64, error) {
	claims := &dto.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(o.Cfg.EncryptionKey), nil
		})
	if err != nil {
		return 0, fmt.Errorf("failed to ParseWithClaims: %w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("failed auth token not valid: %w", err)
	}

	return claims.UserID, nil
}
