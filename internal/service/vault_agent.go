package service

import (
	"context"
	"fmt"
	"keeper/internal/dto"
	pb "keeper/internal/proto/v1"
	pbModel "keeper/internal/proto/v1/model"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type RemoteVaultService interface {
	GetSecret(ctx context.Context, token string, path string) (*dto.AgentGetSecret, error)
	ListSecretPaths(ctx context.Context, token string) ([]string, error)
	SaveSecret(ctx context.Context, req *dto.AgentCreateSecret) error
	DeleteSecret(ctx context.Context, token, path string) error
	DestroySecret(ctx context.Context, token, path string) error
	DeleteMetadata(ctx context.Context, token, path string) error
	UndeleteSecret(ctx context.Context, token, path string, version int64) error
}

type remoteVaultService struct {
	client pb.DataServiceClient
}

func NewRemoteVaultService(client pb.DataServiceClient) RemoteVaultService {
	return &remoteVaultService{client: client}
}

func (s *remoteVaultService) GetSecret(ctx context.Context, token string, path string) (*dto.AgentGetSecret, error) {
	pbReq := &pbModel.GetSecretRequest{}
	pbReq.SetToken(token)
	pbReq.SetPath(path)
	resp, err := s.client.GetSecret(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	var deletedAt *time.Time
	if ts := resp.GetDeletedAt(); ts != nil {
		t := ts.AsTime()
		deletedAt = &t
	}

	return &dto.AgentGetSecret{
		Path:        resp.GetPath(),
		Description: resp.GetDescription(),
		Payload:     resp.GetValue(),
		DeletedAt:   deletedAt,
		Version:     resp.GetVersion(),
		ExpiredAt:   resp.GetExpiredAt().AsTime(),
		CreatedAt:   resp.GetCreatedAt().AsTime(),
	}, nil
}

func (s *remoteVaultService) ListSecretPaths(ctx context.Context, token string) ([]string, error) {
	req := &pbModel.ListSecretPathsRequest{}
	req.SetToken(token)
	resp, err := s.client.ListSecrets(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}
	paths := resp.GetPaths()

	return paths, nil
}

func (s *remoteVaultService) SaveSecret(ctx context.Context, req *dto.AgentCreateSecret) error {
	pbReq := &pbModel.WriteSecret{}
	pbReq.SetToken(req.Token)
	pbReq.SetPath(req.Path)
	pbReq.SetDescription(req.Description)
	pbReq.SetValue(req.Payload)
	pbReq.SetExpiredAt(timestamppb.New(req.ExpiredAt))

	_, err := s.client.SaveSecret(ctx, pbReq)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}
	return nil
}

func (s *remoteVaultService) DeleteSecret(ctx context.Context, token string, path string) error {
	pbReq := &pbModel.DeleteSecretRequest{}
	pbReq.SetToken(token)
	pbReq.SetPath(path)
	_, err := s.client.DeleteSecret(ctx, pbReq)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}
	return nil
}

func (s *remoteVaultService) DestroySecret(ctx context.Context, token string, path string) error {
	pbReq := &pbModel.DeleteSecretRequest{}
	pbReq.SetToken(token)
	pbReq.SetPath(path)
	_, err := s.client.DestroySecret(ctx, pbReq)
	if err != nil {
		return fmt.Errorf("failed to destroy secret: %w", err)
	}
	return nil
}

func (s *remoteVaultService) DeleteMetadata(ctx context.Context, token string, path string) error {
	pbReq := &pbModel.DeleteSecretRequest{}
	pbReq.SetToken(token)
	pbReq.SetPath(path)
	_, err := s.client.DeleteMetadata(ctx, pbReq)
	if err != nil {
		return fmt.Errorf("failed to delete metadata secret: %w", err)
	}
	return nil
}

func (s *remoteVaultService) UndeleteSecret(ctx context.Context, token string, path string, version int64) error {
	pbReq := &pbModel.UndeleteSecretRequest{}
	pbReq.SetToken(token)
	pbReq.SetPath(path)
	pbReq.SetVersion(version)
	_, err := s.client.UndeleteSecret(ctx, pbReq)
	if err != nil {
		return fmt.Errorf("failed to undelete secret: %w", err)
	}
	return nil
}
