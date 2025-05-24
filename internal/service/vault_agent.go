package service

import (
	"context"
	"fmt"
	"keeper/internal/dto"
	pb "keeper/internal/proto/v1"
	pbModel "keeper/internal/proto/v1/model"
	"keeper/internal/util"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type RemoteVaultService interface {
	SaveSecret(ctx context.Context, req *dto.AgentCreateSecret) error
	GetSecret(ctx context.Context, token string, path string) (*dto.AgentGetSecret, error)
	DeleteSecret(ctx context.Context, token string, path string) error
	ListSecretPaths(ctx context.Context, token string) ([]string, error)
}

type remoteVaultService struct {
	client pb.DataServiceClient
}

func NewRemoteVaultService(client pb.DataServiceClient) RemoteVaultService {
	return &remoteVaultService{client: client}
}

func (s *remoteVaultService) SaveSecret(ctx context.Context, req *dto.AgentCreateSecret) error {
	pbReq := &pbModel.WriteSecret{
		Token:       util.Ptr(req.Token),
		Path:        util.Ptr(req.Path),
		Description: util.Ptr(req.Description),
		Value:       req.Payload,
		ExpiredAt:   timestamppb.New(req.ExpiredAt),
	}

	_, err := s.client.SaveSecret(ctx, pbReq)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}
	return nil
}

func (s *remoteVaultService) GetSecret(ctx context.Context, token string, path string) (*dto.AgentGetSecret, error) {
	pbReq := &pbModel.GetSecretRequest{
		Token: util.Ptr(token),
		Path:  util.Ptr(path),
	}
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

func (s *remoteVaultService) DeleteSecret(ctx context.Context, token string, path string) error {
	pbReq := &pbModel.DeleteSecretRequest{
		Token: util.Ptr(token),
		Path:  util.Ptr(path),
	}
	_, err := s.client.DeleteSecret(ctx, pbReq)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}
	return nil
}

func (s *remoteVaultService) ListSecretPaths(ctx context.Context, token string) ([]string, error) {
	req := &pbModel.ListSecretPathsRequest{
		Token: util.Ptr(token),
	}
	resp, err := s.client.ListSecrets(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}
	paths := resp.GetPaths()

	return paths, nil
}
