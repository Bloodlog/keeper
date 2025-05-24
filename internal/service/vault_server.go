package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"keeper/internal/dto"
	"keeper/internal/security"

	"keeper/internal/entity"
	"keeper/internal/repository"
)

type VaultService interface {
	SaveSecret(ctx context.Context, request *dto.ServerCreateSecret) error
	GetSecret(ctx context.Context, userID int64, path string) (dto.DecryptedSecretResponse, error)
	ListSecretsPaths(ctx context.Context, userID int64) ([]string, error)
	DeleteSecret(ctx context.Context, userID int64, path string) error
}

type vaultService struct {
	repo              repository.VaultRepositoryInterface
	dataEncryptionKey string
}

func NewVaultService(repo repository.VaultRepositoryInterface, dataEncryptionKey string) VaultService {
	return &vaultService{
		repo:              repo,
		dataEncryptionKey: dataEncryptionKey,
	}
}

func (s *vaultService) SaveSecret(ctx context.Context, request *dto.ServerCreateSecret) error {
	key, err := hex.DecodeString(s.dataEncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decode data key: %w", err)
	}

	encrypted, err := security.EncryptAESGCM(request.Payload, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %w", err)
	}

	secret := &entity.Secret{
		UserID:      request.UserID,
		Path:        request.Path,
		ExpiredAt:   request.ExpiredAt,
		Description: request.Description,
		Value:       encrypted,
	}

	_, err = s.repo.SaveOrUpdate(ctx, secret)
	if err != nil {
		return fmt.Errorf("failed to save secret: %w", err)
	}
	return nil
}

func (s *vaultService) GetSecret(ctx context.Context, userID int64, path string) (dto.DecryptedSecretResponse, error) {
	secret, err := s.repo.GetByUserAndPath(ctx, userID, path)
	if err != nil {
		return dto.DecryptedSecretResponse{}, fmt.Errorf("failed to get secret: %w", err)
	}

	key, err := hex.DecodeString(s.dataEncryptionKey)
	if err != nil {
		return dto.DecryptedSecretResponse{}, fmt.Errorf("failed to decode data key: %w", err)
	}

	decrypted, err := security.DecryptAESGCM(key, secret.Value)
	if err != nil {
		return dto.DecryptedSecretResponse{}, fmt.Errorf("failed to decrypt secret: %w", err)
	}

	return dto.DecryptedSecretResponse{
		UserID:      secret.UserID,
		Path:        secret.Path,
		ExpiredAt:   secret.ExpiredAt,
		Description: secret.Description,
		CreatedAt:   secret.CreatedAt,
		Data:        decrypted,
		Version:     secret.Version,
		DeletedAt:   secret.DeletedAt,
	}, nil
}

func (s *vaultService) ListSecretsPaths(ctx context.Context, userID int64) ([]string, error) {
	secrets, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	paths := make([]string, 0, len(secrets))
	for i := range secrets {
		paths = append(paths, secrets[i].Path)
	}
	return paths, nil
}

func (s *vaultService) DeleteSecret(ctx context.Context, userID int64, path string) error {
	err := s.repo.Delete(ctx, userID, path)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}
	return nil
}
