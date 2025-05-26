package handler

import (
	"context"
	"fmt"
	"keeper/internal/dto"
	"keeper/internal/logger"
	pb "keeper/internal/proto/v1"
	pbModel "keeper/internal/proto/v1/model"
	"keeper/internal/service"
	"keeper/internal/util"
)

type AuthServerHandler struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
	logger      *logger.ZapLogger
}

func NewAuthHandler(l *logger.ZapLogger, svc service.AuthService) *AuthServerHandler {
	return &AuthServerHandler{
		authService: svc,
		logger:      l,
	}
}

func (s *AuthServerHandler) Register(
	ctx context.Context,
	req *pbModel.RegisterRequest,
) (*pbModel.RegisterResponse, error) {
	requestDto := &dto.RegisterUser{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}
	token, err := s.authService.Register(ctx, requestDto)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return &pbModel.RegisterResponse{
		Success: util.Ptr(true),
		Message: util.Ptr("Register successful."),
		Token:   &token.Token,
	}, nil
}

func (s *AuthServerHandler) Login(
	ctx context.Context,
	req *pbModel.LoginRequest,
) (*pbModel.LoginResponse, error) {
	requestDto := &dto.LoginUser{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}
	token, err := s.authService.Login(ctx, requestDto)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return &pbModel.LoginResponse{
		Success: util.Ptr(true),
		Message: util.Ptr("Login successful."),
		Token:   &token.Token,
	}, nil
}
