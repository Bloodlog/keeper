package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"keeper/internal/dto"
	"keeper/internal/logger"
	pb "keeper/internal/proto/v1"
	pbModel "keeper/internal/proto/v1/model"
	"keeper/internal/service"
	"keeper/internal/util"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const errorInvalidToken = "invalid token: %w"

type VaultServerHandler struct {
	pb.UnimplementedDataServiceServer
	vaultService service.VaultService
	jwtService   service.JwtService
	logger       *logger.ZapLogger
}

func NewVaultHandler(l *logger.ZapLogger, svc service.VaultService, jwtService service.JwtService) *VaultServerHandler {
	return &VaultServerHandler{
		vaultService: svc,
		jwtService:   jwtService,
		logger:       l,
	}
}

func (s *VaultServerHandler) GetSecret(
	ctx context.Context,
	req *pbModel.GetSecretRequest,
) (*pbModel.SecretResponse, error) {
	userID, err := s.jwtService.GetUserID(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	secret, err := s.vaultService.GetSecret(ctx, userID, req.GetPath())
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	data, err := json.Marshal(secret.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal secret data: %w", err)
	}

	var deletedAt *timestamppb.Timestamp
	if secret.DeletedAt != nil {
		deletedAt = timestamppb.New(*secret.DeletedAt)
	}

	return &pbModel.SecretResponse{
		Path:        util.Ptr(secret.Path),
		Description: util.Ptr(secret.Description),
		Value:       data,
		ExpiredAt:   timestamppb.New(secret.ExpiredAt),
		Version:     util.Ptr(secret.Version),
		DeletedAt:   deletedAt,
	}, nil
}

func (s *VaultServerHandler) ListSecrets(
	ctx context.Context,
	req *pbModel.ListSecretPathsRequest,
) (*pbModel.ListSecretPathsResponse, error) {
	userID, err := s.jwtService.GetUserID(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	paths, err := s.vaultService.ListSecretsPaths(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	return &pbModel.ListSecretPathsResponse{Paths: paths}, nil
}

func (s *VaultServerHandler) SaveSecret(
	ctx context.Context,
	req *pbModel.WriteSecret,
) (*pbModel.SaveSecretResponse, error) {
	userID, err := s.jwtService.GetUserID(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	serverCreateSecretDTO := &dto.ServerCreateSecret{
		UserID:      userID,
		Path:        req.GetPath(),
		Description: req.GetDescription(),
		Payload:     req.GetValue(),
		ExpiredAt:   req.GetExpiredAt().AsTime(),
	}

	err = s.vaultService.SaveSecret(ctx, serverCreateSecretDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to save secret: %w", err)
	}

	return &pbModel.SaveSecretResponse{
		Success: util.Ptr(true),
		Message: util.Ptr("Save: success"),
	}, nil
}

func (s *VaultServerHandler) DeleteSecret(
	ctx context.Context,
	req *pbModel.DeleteSecretRequest,
) (*pbModel.DeleteSecretResponse, error) {
	userID, err := s.jwtService.GetUserID(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	err = s.vaultService.DeleteSecret(ctx, userID, req.GetPath())
	if err != nil {
		return nil, fmt.Errorf("failed to delete secret: %w", err)
	}

	return &pbModel.DeleteSecretResponse{Message: util.Ptr("Soft delete: success")}, nil
}

func (s *VaultServerHandler) DestroySecret(
	ctx context.Context,
	req *pbModel.DeleteSecretRequest,
) (*pbModel.DeleteSecretResponse, error) {
	userID, err := s.jwtService.GetUserID(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	err = s.vaultService.DestroySecret(ctx, userID, req.GetPath())
	if err != nil {
		return nil, fmt.Errorf("failed to destroy secret: %w", err)
	}

	return &pbModel.DeleteSecretResponse{Message: util.Ptr("Destroy: success")}, nil
}

func (s *VaultServerHandler) DeleteMetadata(
	ctx context.Context,
	req *pbModel.DeleteSecretRequest,
) (*pbModel.DeleteSecretResponse, error) {
	userID, err := s.jwtService.GetUserID(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	err = s.vaultService.DeleteMetadata(ctx, userID, req.GetPath())
	if err != nil {
		return nil, fmt.Errorf("failed to delete metadata secret: %w", err)
	}

	return &pbModel.DeleteSecretResponse{Message: util.Ptr("Delete key with metadata: success")}, nil
}

func (s *VaultServerHandler) UndeleteSecret(
	ctx context.Context,
	req *pbModel.UndeleteSecretRequest,
) (*pbModel.DeleteSecretResponse, error) {
	userID, err := s.jwtService.GetUserID(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	err = s.vaultService.UndeleteSecret(ctx, userID, req.GetPath(), req.GetVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to undelete secret: %w", err)
	}

	return &pbModel.DeleteSecretResponse{Message: util.Ptr("Undelete: success")}, nil
}
