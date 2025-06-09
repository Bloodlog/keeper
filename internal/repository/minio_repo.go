package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
)

type FileRepository interface {
	Save(ctx context.Context, fileName string, data []byte) error
	Load(ctx context.Context, fileName string) ([]byte, error)
	Delete(ctx context.Context, fileName string) error
}

type MinIORepository struct {
	Client        *minio.Client
	BucketName    string
	URLExpiredTTL time.Duration
}

func NewMinIORepository(client *minio.Client, bucketName string, uRLExpiredTTL time.Duration) *MinIORepository {
	return &MinIORepository{
		Client:        client,
		BucketName:    bucketName,
		URLExpiredTTL: uRLExpiredTTL,
	}
}

func (m *MinIORepository) Save(ctx context.Context, fileName string, data []byte) error {
	_, err := m.Client.PutObject(
		ctx, m.BucketName,
		fileName, bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func (m *MinIORepository) Load(ctx context.Context, fileName string) ([]byte, error) {
	obj, err := m.Client.GetObject(ctx, m.BucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	defer func() {
		if err := obj.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("read object: %w", err)
	}
	return data, nil
}

func (m *MinIORepository) Delete(ctx context.Context, fileName string) error {
	err := m.Client.RemoveObject(
		ctx, m.BucketName,
		fileName,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("remove object: %w", err)
	}

	return nil
}
