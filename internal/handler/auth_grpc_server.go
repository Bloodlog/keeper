package handler

import (
	"context"
	"fmt"
	"keeper/internal/dto"
	"keeper/internal/logger"
	pb "keeper/internal/proto/v1"
	pbModel "keeper/internal/proto/v1/model"
	"keeper/internal/service"
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
	registerDto := &dto.RegisterUser{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}
	token, err := s.authService.Register(ctx, registerDto)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return fillAuthResponse(&pbModel.RegisterResponse{}, "Register successful.", token.Token), nil
}

func (s *AuthServerHandler) Login(
	ctx context.Context,
	req *pbModel.LoginRequest,
) (*pbModel.LoginResponse, error) {
	loginDto := &dto.LoginUser{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}
	token, err := s.authService.Login(ctx, loginDto)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return fillAuthResponse(&pbModel.LoginResponse{}, "Login successful.", token.Token), nil
}

func fillAuthResponse[T interface {
	SetSuccess(bool)
	SetMessage(string)
	SetToken(string)
}](resp T, message, token string) T {
	resp.SetSuccess(true)
	resp.SetMessage(message)
	resp.SetToken(token)
	return resp
}
