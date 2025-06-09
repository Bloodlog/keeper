package service

import (
	"errors"
	"keeper/internal/dto"
	"keeper/internal/proto/v1/mock"
	"keeper/internal/proto/v1/model"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRemoteVaultService_GetSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockDataServiceClient(ctrl)
	svc := NewRemoteVaultService(mockClient)

	now := time.Now()
	deletedAt := now.Add(-1 * time.Hour)
	createdAt := now.Add(-2 * time.Hour)
	expiredAt := now.Add(1 * time.Hour)

	mockResp := &model.SecretResponse{}
	mockResp.SetPath("secret/foo")
	mockResp.SetDescription("my secret")
	mockResp.SetValue([]byte("supersecret"))
	mockResp.SetDeletedAt(timestamppb.New(deletedAt))
	mockResp.SetVersion(2)
	mockResp.SetCreatedAt(timestamppb.New(createdAt))
	mockResp.SetExpiredAt(timestamppb.New(expiredAt))

	mockClient.EXPECT().
		GetSecret(gomock.Any(), gomock.Any()).
		Return(mockResp, nil)

	result, err := svc.GetSecret(t.Context(), "token123", "secret/foo")
	require.NoError(t, err)
	require.Equal(t, "secret/foo", result.Path)
	require.Equal(t, "my secret", result.Description)
	require.Equal(t, []byte("supersecret"), result.Payload)
	require.Equal(t, int64(2), result.Version)
	require.WithinDuration(t, deletedAt, *result.DeletedAt, time.Second)
	require.WithinDuration(t, createdAt, result.CreatedAt, time.Second)
	require.WithinDuration(t, expiredAt, result.ExpiredAt, time.Second)
}

func TestRemoteVaultService_GetSecret_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockDataServiceClient(ctrl)
	svc := NewRemoteVaultService(mockClient)

	mockClient.EXPECT().
		GetSecret(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("rpc error"))

	result, err := svc.GetSecret(t.Context(), "bad-token", "bad/path")
	require.Nil(t, result)
	require.ErrorContains(t, err, "failed to get secret")
}

func TestRemoteVaultService_ListSecretPaths(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockDataServiceClient(ctrl)
	svc := NewRemoteVaultService(mockClient)

	mockResp := &model.ListSecretPathsResponse{}
	mockResp.SetPaths([]string{"secret/foo", "secret/bar"})

	mockClient.EXPECT().
		ListSecrets(gomock.Any(), gomock.Any()).
		Return(mockResp, nil)

	paths, err := svc.ListSecretPaths(t.Context(), "token123")
	require.NoError(t, err)
	require.Equal(t, []string{"secret/foo", "secret/bar"}, paths)
}

func TestRemoteVaultService_SaveSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockDataServiceClient(ctrl)
	svc := NewRemoteVaultService(mockClient)

	req := &dto.AgentCreateSecret{
		Token:       "token123",
		Path:        "secret/foo",
		Description: "desc",
		Payload:     []byte("data"),
		ExpiredAt:   time.Now().Add(time.Hour),
	}

	mockClient.EXPECT().
		SaveSecret(gomock.Any(), gomock.Any()).
		Return(&model.SaveSecretResponse{}, nil)

	err := svc.SaveSecret(t.Context(), req)
	require.NoError(t, err)
}

func TestRemoteVaultService_DeleteSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockDataServiceClient(ctrl)
	svc := NewRemoteVaultService(mockClient)

	mockClient.EXPECT().
		DeleteSecret(gomock.Any(), gomock.Any()).
		Return(&model.DeleteSecretResponse{}, nil)

	err := svc.DeleteSecret(t.Context(), "token", "secret/foo")
	require.NoError(t, err)
}

func TestRemoteVaultService_DeleteSecret_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockDataServiceClient(ctrl)
	svc := NewRemoteVaultService(mockClient)

	mockClient.EXPECT().
		DeleteSecret(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("fail"))

	err := svc.DeleteSecret(t.Context(), "token", "secret/foo")
	require.ErrorContains(t, err, "failed to delete secret")
}
