package repository

import (
	"context"
	"errors"
	"fmt"
	"keeper/internal/entity"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VaultRepositoryInterface interface {
	GetByUserAndPath(ctx context.Context, userID int64, path string) (entity.Secret, error)
	SaveOrUpdate(ctx context.Context, secret *entity.Secret) (entity.Secret, error)
	Delete(ctx context.Context, userID int64, path string) error
	ListByUser(ctx context.Context, userID int64) ([]entity.Secret, error)
}

type vaultRepository struct {
	Pool *pgxpool.Pool
}

func NewVaultRepository(db *pgxpool.Pool) VaultRepositoryInterface {
	return &vaultRepository{Pool: db}
}

func (r *vaultRepository) SaveOrUpdate(ctx context.Context, secret *entity.Secret) (entity.Secret, error) {
	query := `
		INSERT INTO secrets (user_id, title, expired_at, description, version, content, deleted_at)
		VALUES ($1, $2, $3, $4, 1, $5, NULL)
		ON CONFLICT (user_id, title) DO UPDATE
		SET expired_at = EXCLUDED.expired_at,
		    description = EXCLUDED.description,
		    content = EXCLUDED.content,
		    version = secrets.version + 1,
			deleted_at = NULL
		RETURNING created_at, version
	`

	var createdAt time.Time
	err := r.Pool.QueryRow(ctx, query,
		secret.UserID, secret.Path,
		secret.ExpiredAt, secret.Description,
		secret.Value,
	).Scan(&createdAt, &secret.Version)
	if err != nil {
		return *secret, fmt.Errorf("failed to save or update secret: %w", err)
	}

	secret.CreatedAt = createdAt
	return *secret, nil
}

func (r *vaultRepository) GetByUserAndPath(ctx context.Context, userID int64, path string) (entity.Secret, error) {
	var secret entity.Secret
	query := `
		SELECT user_id, title, expired_at, description, content, created_at, version, deleted_at
		FROM secrets
		WHERE user_id = $1 AND title = $2
	`
	err := r.Pool.QueryRow(ctx, query, userID, path).Scan(
		&secret.UserID, &secret.Path,
		&secret.ExpiredAt, &secret.Description, &secret.Value, &secret.CreatedAt, &secret.Version, &secret.DeletedAt,
	)
	if err != nil {
		return secret, fmt.Errorf("failed to get secret: %w", err)
	}
	return secret, nil
}

func (r *vaultRepository) Delete(ctx context.Context, userID int64, path string) error {
	query := `
		UPDATE secrets
		SET deleted_at = NOW()
		WHERE user_id = $1 AND title = $2 AND deleted_at IS NULL
	`
	ct, err := r.Pool.Exec(ctx, query, userID, path)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}

func (r *vaultRepository) ListByUser(ctx context.Context, userID int64) ([]entity.Secret, error) {
	query := `
		SELECT user_id, title, expired_at, description, content, created_at, version, deleted_at
		FROM secrets
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}
	defer rows.Close()

	var secrets []entity.Secret
	for rows.Next() {
		var s entity.Secret
		if err := rows.Scan(&s.UserID, &s.Path, &s.ExpiredAt, &s.Description, &s.Value, &s.CreatedAt,
			&s.Version, &s.DeletedAt); err != nil {
			return nil, fmt.Errorf("failed to scan secret: %w", err)
		}
		secrets = append(secrets, s)
	}
	return secrets, nil
}
