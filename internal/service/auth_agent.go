package service

import (
	"context"
	"fmt"
	"keeper/internal/dto"
	pb "keeper/internal/proto/v1"
	pbModel "keeper/internal/proto/v1/model"
)

type RemoteAuthService interface {
	Register(ctx context.Context, req dto.RegisterUser) (string, error)
	Login(ctx context.Context, req dto.LoginUser) (string, error)
}

type remoteAuthService struct {
	client pb.AuthServiceClient
}

func NewRemoteAuthService(client pb.AuthServiceClient) RemoteAuthService {
	return &remoteAuthService{client: client}
}

func (s *remoteAuthService) Register(ctx context.Context, req dto.RegisterUser) (string, error) {
	requestDto := &pbModel.RegisterRequest{
		Login:    &req.Login,
		Password: &req.Password,
	}
	resp, err := s.client.Register(ctx, requestDto)
	if err != nil {
		return "", fmt.Errorf("register error: %w", err)
	}

	if !resp.GetSuccess() {
		return "", fmt.Errorf("register failed: %s", resp.GetMessage())
	}

	return resp.GetToken(), nil
}

func (s *remoteAuthService) Login(ctx context.Context, req dto.LoginUser) (string, error) {
	requestDto := &pbModel.LoginRequest{
		Login:    &req.Login,
		Password: &req.Password,
	}
	resp, err := s.client.Login(ctx, requestDto)
	if err != nil {
		return "", fmt.Errorf("login error: %w", err)
	}

	if !resp.GetSuccess() {
		return "", fmt.Errorf("login failed: %s", resp.GetMessage())
	}

	return resp.GetToken(), nil
}
