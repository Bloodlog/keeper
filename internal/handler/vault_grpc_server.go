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
	utils "keeper/internal/util"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const errorInvalidToken = "invalid token: %w"

type VaultServerHandler struct {
	pb.UnimplementedDataServiceServer
	vaultService service.VaultService
	logger       *logger.ZapLogger
}

func NewVaultHandler(l *logger.ZapLogger, svc service.VaultService) *VaultServerHandler {
	return &VaultServerHandler{
		vaultService: svc,
		logger:       l,
	}
}

func (s *VaultServerHandler) GetSecret(
	ctx context.Context,
	req *pbModel.GetSecretRequest,
) (*pbModel.SecretResponse, error) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	secret, err := s.vaultService.GetSecret(ctx, userID, req.GetPath())
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	var value []byte
	if secret.FilePath != nil {
		value = secret.Data
	} else {
		value, err = json.Marshal(secret.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal secret data: %w", err)
		}
	}

	var deletedAt *timestamppb.Timestamp
	if secret.DeletedAt != nil {
		deletedAt = timestamppb.New(*secret.DeletedAt)
	}
	var filePath string
	if secret.FilePath != nil {
		filePath = *secret.FilePath
	}
	resp := &pbModel.SecretResponse{}
	resp.SetPath(secret.Path)
	resp.SetDescription(secret.Description)
	resp.SetValue(value)
	resp.SetExpiredAt(timestamppb.New(secret.ExpiredAt))
	resp.SetVersion(secret.Version)
	resp.SetDeletedAt(deletedAt)
	resp.SetFilePath(filePath)

	return resp, nil
}

func (s *VaultServerHandler) ListSecrets(
	ctx context.Context,
	req *pbModel.ListSecretPathsRequest,
) (*pbModel.ListSecretPathsResponse, error) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	paths, err := s.vaultService.ListSecretsPaths(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	resp := &pbModel.ListSecretPathsResponse{}
	resp.SetPaths(paths)

	return resp, nil
}

func (s *VaultServerHandler) SaveSecret(
	ctx context.Context,
	req *pbModel.WriteSecret,
) (*pbModel.SaveSecretResponse, error) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	var filePath *string
	if fp := req.GetFilePath(); fp != "" {
		filePath = &fp
	}

	serverCreateSecretDTO := &dto.ServerCreateSecret{
		UserID:      userID,
		Path:        req.GetPath(),
		Description: req.GetDescription(),
		Payload:     req.GetValue(),
		ExpiredAt:   req.GetExpiredAt().AsTime(),
		FilePath:    filePath,
	}

	err = s.vaultService.SaveSecret(ctx, serverCreateSecretDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to save secret: %w", err)
	}
	resp := &pbModel.SaveSecretResponse{}
	resp.SetSuccess(true)
	resp.SetMessage("Save: success")

	return resp, nil
}

func (s *VaultServerHandler) DeleteSecret(
	ctx context.Context,
	req *pbModel.DeleteSecretRequest,
) (*pbModel.DeleteSecretResponse, error) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	err = s.vaultService.DeleteSecret(ctx, userID, req.GetPath())
	if err != nil {
		return nil, fmt.Errorf("failed to delete secret: %w", err)
	}

	resp := &pbModel.DeleteSecretResponse{}
	resp.SetMessage("Soft delete: success")

	return resp, nil
}

func (s *VaultServerHandler) DestroySecret(
	ctx context.Context,
	req *pbModel.DeleteSecretRequest,
) (*pbModel.DeleteSecretResponse, error) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	err = s.vaultService.DestroySecret(ctx, userID, req.GetPath())
	if err != nil {
		return nil, fmt.Errorf("failed to destroy secret: %w", err)
	}

	resp := &pbModel.DeleteSecretResponse{}
	resp.SetMessage("Destroy: success")

	return resp, nil
}

func (s *VaultServerHandler) DeleteMetadata(
	ctx context.Context,
	req *pbModel.DeleteSecretRequest,
) (*pbModel.DeleteSecretResponse, error) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	err = s.vaultService.DeleteMetadata(ctx, userID, req.GetPath())
	if err != nil {
		return nil, fmt.Errorf("failed to delete metadata secret: %w", err)
	}

	resp := &pbModel.DeleteSecretResponse{}
	resp.SetMessage("Delete key with metadata: success")

	return resp, nil
}

func (s *VaultServerHandler) UndeleteSecret(
	ctx context.Context,
	req *pbModel.UndeleteSecretRequest,
) (*pbModel.DeleteSecretResponse, error) {
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf(errorInvalidToken, err)
	}

	err = s.vaultService.UndeleteSecret(ctx, userID, req.GetPath(), req.GetVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to undelete secret: %w", err)
	}

	resp := &pbModel.DeleteSecretResponse{}
	resp.SetMessage("Undelete: success")

	return resp, nil
}
