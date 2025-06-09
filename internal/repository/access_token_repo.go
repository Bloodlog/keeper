package repository

import (
	"context"
	"errors"
	"fmt"
	"keeper/internal/entity"

	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

type AccessTokenRepository interface {
	Create(ctx context.Context, tx pgx.Tx, token *entity.AccessToken) error
	FindValidByUserID(ctx context.Context, tx pgx.Tx, userID int64) (entity.AccessToken, error)
}

type accessTokenRepository struct {
	Pool *pgxpool.Pool
}

func NewAccessRepository(db *pgxpool.Pool) AccessTokenRepository {
	return &accessTokenRepository{
		Pool: db,
	}
}

func (r *accessTokenRepository) Create(ctx context.Context, tx pgx.Tx, token *entity.AccessToken) error {
	query := `
		INSERT INTO access_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := tx.Exec(ctx, query, token.UserID, token.Token, token.ExpiresAt)

	if err != nil {
		return fmt.Errorf("error create Token: %w", err)
	}

	return nil
}

func (r *accessTokenRepository) FindValidByUserID(
	ctx context.Context,
	tx pgx.Tx,
	userID int64,
) (entity.AccessToken, error) {
	var token entity.AccessToken

	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM access_tokens
		WHERE user_id = $1 AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := tx.QueryRow(ctx, query, userID).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return token, fmt.Errorf("no valid token found: %w", err)
		}
		return token, fmt.Errorf("failed to find valid token: %w", err)
	}

	return token, nil
}
