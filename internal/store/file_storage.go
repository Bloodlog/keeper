package store

import (
	"context"
	"fmt"
	"keeper/internal/config"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileStorage struct {
	FileClient *minio.Client
}

func NewFileStorage(ctx context.Context, cfg *config.FileStorageConfig) (*FileStorage, error) {
	var fileAddress = cfg.Address + ":" + strconv.Itoa(cfg.Port)
	client, err := minio.New(
		fileAddress, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.Username, cfg.Password, ""),
			Secure: false,
		})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	return &FileStorage{
		FileClient: client,
	}, nil
}
