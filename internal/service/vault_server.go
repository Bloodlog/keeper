package service

import (
	"context"
	"fmt"
	"keeper/internal/dto"
	"keeper/internal/entity"
	"keeper/internal/repository"
)

const errorDeleteSecret = "failed to delete secret: %w"

type VaultService interface {
	GetSecret(ctx context.Context, userID int64, path string) (dto.DecryptedSecretResponse, error)
	ListSecretsPaths(ctx context.Context, userID int64) ([]string, error)
	SaveSecret(ctx context.Context, request *dto.ServerCreateSecret) error
	DeleteSecret(ctx context.Context, userID int64, path string) error
	DestroySecret(ctx context.Context, userID int64, path string) error
	DeleteMetadata(ctx context.Context, userID int64, path string) error
	UndeleteSecret(ctx context.Context, userID int64, path string, version int64) error
}

type vaultService struct {
	repo          repository.VaultRepositoryInterface
	cryptoService CryptoService
	fileRepo      repository.FileRepository
}

func NewVaultService(
	repo repository.VaultRepositoryInterface,
	cryptoService CryptoService,
	fileRepo repository.FileRepository,
) VaultService {
	return &vaultService{
		repo:          repo,
		cryptoService: cryptoService,
		fileRepo:      fileRepo,
	}
}

func (s *vaultService) GetSecret(
	ctx context.Context,
	userID int64,
	path string,
) (dto.DecryptedSecretResponse, error) {
	secret, err := s.repo.GetByUserAndPath(ctx, userID, path)
	if err != nil {
		return dto.DecryptedSecretResponse{}, fmt.Errorf("failed to get secret: %w", err)
	}

	var decrypted []byte
	if secret.FilePath != nil {
		var file []byte
		file, err = s.fileRepo.Load(ctx, *secret.FilePath)
		if err != nil {
			return dto.DecryptedSecretResponse{}, fmt.Errorf("failed to load secret: %w", err)
		}
		decrypted, err = s.cryptoService.Decode(file)
		if err != nil {
			return dto.DecryptedSecretResponse{}, fmt.Errorf("failed to decrypt secret data: %w", err)
		}
	} else {
		decrypted, err = s.cryptoService.Decode(secret.Value)
		if err != nil {
			return dto.DecryptedSecretResponse{}, fmt.Errorf("failed to decrypt secret: %w", err)
		}
	}

	return dto.DecryptedSecretResponse{
		Path:        secret.Path,
		ExpiredAt:   secret.ExpiredAt,
		Description: secret.Description,
		CreatedAt:   secret.CreatedAt,
		Data:        decrypted,
		Version:     secret.Version,
		DeletedAt:   secret.DeletedAt,
		FilePath:    secret.FilePath,
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

func (s *vaultService) SaveSecret(ctx context.Context, request *dto.ServerCreateSecret) error {
	var secretVersion *entity.SecretVersion

	encrypted, err := s.cryptoService.Encode(request.Payload)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %w", err)
	}

	secretMetadata := &entity.SecretMetadata{
		UserID:      request.UserID,
		Path:        request.Path,
		ExpiredAt:   request.ExpiredAt,
		Description: request.Description,
	}
	if request.FilePath != nil && s.fileRepo != nil {
		err = s.fileRepo.Save(ctx, *request.FilePath, encrypted)
		if err != nil {
			return fmt.Errorf("failed to store file: %w", err)
		}
		placeholder, err := s.cryptoService.Encode([]byte(`{"status": "FILE-UPLOADED"}`))
		if err != nil {
			return fmt.Errorf("failed to encrypt file placeholder: %w", err)
		}

		secretVersion = &entity.SecretVersion{
			Value:    placeholder,
			FilePath: request.FilePath,
		}
	} else {
		secretVersion = &entity.SecretVersion{
			Value: encrypted,
		}
	}

	_, err = s.repo.SaveOrUpdate(ctx, secretMetadata, secretVersion)
	if err != nil {
		return fmt.Errorf("failed to save secret: %w", err)
	}
	return nil
}

func (s *vaultService) DeleteSecret(ctx context.Context, userID int64, path string) error {
	err := s.repo.Delete(ctx, userID, path)
	if err != nil {
		return fmt.Errorf(errorDeleteSecret, err)
	}
	return nil
}

func (s *vaultService) DestroySecret(ctx context.Context, userID int64, path string) error {
	err := s.repo.DestroySecret(ctx, userID, path)
	if err != nil {
		return fmt.Errorf(errorDeleteSecret, err)
	}
	return nil
}

func (s *vaultService) DeleteMetadata(ctx context.Context, userID int64, path string) error {
	err := s.repo.DeleteMetadata(ctx, userID, path)
	if err != nil {
		return fmt.Errorf(errorDeleteSecret, err)
	}
	return nil
}

func (s *vaultService) UndeleteSecret(ctx context.Context, userID int64, path string, version int64) error {
	err := s.repo.UndeleteSecret(ctx, userID, path, version)
	if err != nil {
		return fmt.Errorf(errorDeleteSecret, err)
	}
	return nil
}
