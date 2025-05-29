package repository

import (
	"context"
	"errors"
	"fmt"
	"keeper/internal/entity"

	pgx "github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VaultRepositoryInterface interface {
	GetByUserAndPath(ctx context.Context, userID int64, path string) (entity.OneSecretVersionWithMetadata, error)
	ListByUser(ctx context.Context, userID int64) ([]entity.SecretMetadata, error)
	SaveOrUpdate(ctx context.Context, secretMetadata *entity.SecretMetadata,
		secretVersion *entity.SecretVersion) (entity.SecretMetadata, error)
	Delete(ctx context.Context, userID int64, path string) error
	DestroySecret(ctx context.Context, userID int64, path string) error
	DeleteMetadata(ctx context.Context, userID int64, path string) error
	UndeleteSecret(ctx context.Context, userID int64, path string, version int64) error
}

type vaultRepository struct {
	Pool *pgxpool.Pool
}

func NewVaultRepository(db *pgxpool.Pool) VaultRepositoryInterface {
	return &vaultRepository{Pool: db}
}

func (r *vaultRepository) GetByUserAndPath(
	ctx context.Context,
	userID int64,
	path string,
) (entity.OneSecretVersionWithMetadata, error) {
	var secret entity.OneSecretVersionWithMetadata
	query := `
		SELECT sm.title, sm.expired_at, sm.description,
			sv.content, sv.created_at, sv.version, sv.deleted_at, sv.file_path
		FROM secrets_metadata sm
		JOIN secret_versions sv ON sm.id = sv.metadata_id
		WHERE sm.user_id = $1 AND sm.title = $2 AND sv.deleted_at IS NULL
		ORDER BY sv.version DESC LIMIT 1
	`
	err := r.Pool.QueryRow(ctx, query, userID, path).Scan(
		&secret.Path, &secret.ExpiredAt, &secret.Description,
		&secret.Value, &secret.CreatedAt, &secret.Version, &secret.DeletedAt, &secret.FilePath,
	)
	if err != nil {
		return secret, fmt.Errorf("failed to get secret: %w", err)
	}
	return secret, nil
}

func (r *vaultRepository) ListByUser(ctx context.Context, userID int64) ([]entity.SecretMetadata, error) {
	query := `
		SELECT sm.title
		FROM secrets_metadata sm 
		WHERE sm.user_id = $1
		ORDER BY sm.title
	`
	rows, err := r.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}
	defer rows.Close()

	var secrets []entity.SecretMetadata
	for rows.Next() {
		var s entity.SecretMetadata
		if err := rows.Scan(&s.Path); err != nil {
			return nil, fmt.Errorf("failed to scan secret: %w", err)
		}
		secrets = append(secrets, s)
	}
	return secrets, nil
}

func (r *vaultRepository) SaveOrUpdate(ctx context.Context, secretMetadata *entity.SecretMetadata,
	secretVersion *entity.SecretVersion) (entity.SecretMetadata,
	error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return *secretMetadata, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
	}(tx, ctx)

	var metadataID int64
	metaUpsert := `
		INSERT INTO secrets_metadata (user_id, title, expired_at, description, file_path)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, title) DO UPDATE
		SET expired_at = EXCLUDED.expired_at, description = EXCLUDED.description, deleted_at = NULL
		RETURNING id
	`
	err = tx.QueryRow(ctx, metaUpsert, secretMetadata.UserID, secretMetadata.Path,
		secretMetadata.ExpiredAt, secretMetadata.Description, secretVersion.FilePath).Scan(&metadataID)
	if err != nil {
		return *secretMetadata, fmt.Errorf("failed to upsert metadata: %w", err)
	}

	versionInsert := `
		INSERT INTO secret_versions (metadata_id, version, content)
		SELECT $1, COALESCE(MAX(version), 0)+1, $2 
		FROM secret_versions WHERE metadata_id = $1
		RETURNING version, created_at
	`
	err = tx.QueryRow(ctx, versionInsert, metadataID,
		secretVersion.Value).Scan(&secretVersion.Version, &secretVersion.CreatedAt)
	if err != nil {
		return *secretMetadata, fmt.Errorf("failed to insert version: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return *secretMetadata, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return *secretMetadata, nil
}

func (r *vaultRepository) Delete(ctx context.Context, userID int64, path string) error {
	query := `
		UPDATE secret_versions SET deleted_at = NOW()
		WHERE metadata_id = (SELECT id FROM secrets_metadata WHERE user_id = $1 AND title = $2)
		AND deleted_at IS NULL
	`
	ct, err := r.Pool.Exec(ctx, query, userID, path)
	if err != nil {
		return fmt.Errorf("failed to delete secret versions: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errors.New("no versions deleted")
	}
	return nil
}

func (r *vaultRepository) DestroySecret(ctx context.Context, userID int64, path string) error {
	query := `
		UPDATE secret_versions SET destroyed = TRUE, content = '', deleted_at = NOW(), file_path = ''
		WHERE metadata_id = (SELECT id FROM secrets_metadata WHERE user_id = $1 AND title = $2)
		AND destroyed = FALSE
	`
	ct, err := r.Pool.Exec(ctx, query, userID, path)
	if err != nil {
		return fmt.Errorf("failed to destroy version: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errors.New("no version destroyed")
	}
	return nil
}

func (r *vaultRepository) DeleteMetadata(ctx context.Context, userID int64, path string) error {
	query := `
		UPDATE secrets_metadata SET deleted_at = NOW()
		WHERE user_id = $1 AND title = $2 AND deleted_at IS NULL
	`
	ct, err := r.Pool.Exec(ctx, query, userID, path)
	if err != nil {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errors.New("no metadata deleted")
	}
	return nil
}

func (r *vaultRepository) UndeleteSecret(ctx context.Context, userID int64, path string, version int64) error {
	query := `
		UPDATE secret_versions SET deleted_at = NULL
		WHERE metadata_id = (SELECT id FROM secrets_metadata WHERE user_id = $1 AND title = $2)
		AND version = $3 AND deleted_at IS NOT NULL
	`
	ct, err := r.Pool.Exec(ctx, query, userID, path, version)
	if err != nil {
		return fmt.Errorf("failed to undelete version: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errors.New("no version undeleted")
	}
	return nil
}
