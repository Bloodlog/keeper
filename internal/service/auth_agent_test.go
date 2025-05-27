package service

import (
	"errors"
	"keeper/internal/dto"
	"keeper/internal/proto/v1/mock"
	pbModel "keeper/internal/proto/v1/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRemoteAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockAuthServiceClient(ctrl)
	svc := NewRemoteAuthService(mockClient)

	ctx := t.Context()
	req := dto.RegisterUser{
		Login:    "user",
		Password: "pass",
	}

	t.Run("success", func(t *testing.T) {
		resp := &pbModel.RegisterResponse{}
		resp.SetMessage("")
		resp.SetToken("token123")
		resp.SetSuccess(true)
		mockClient.EXPECT().
			Register(ctx, gomock.Any()).
			Return(resp, nil)

		token, err := svc.Register(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "token123", token)
	})

	t.Run("rpc error", func(t *testing.T) {
		mockClient.EXPECT().
			Register(ctx, gomock.Any()).
			Return(nil, errors.New("connection lost"))

		token, err := svc.Register(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "register error")
		assert.Empty(t, token)
	})

	t.Run("unsuccessful response", func(t *testing.T) {
		resp := &pbModel.RegisterResponse{}
		resp.SetMessage("error")
		resp.SetToken("")
		resp.SetSuccess(false)

		mockClient.EXPECT().
			Register(ctx, gomock.Any()).
			Return(resp, nil)

		token, err := svc.Register(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "register failed")
		assert.Empty(t, token)
	})
}

func TestRemoteAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockAuthServiceClient(ctrl)
	svc := NewRemoteAuthService(mockClient)

	ctx := t.Context()
	req := dto.LoginUser{
		Login:    "user",
		Password: "pass",
	}

	t.Run("success", func(t *testing.T) {
		resp := &pbModel.LoginResponse{}
		resp.SetMessage("")
		resp.SetToken("token456")
		resp.SetSuccess(true)
		mockClient.EXPECT().
			Login(ctx, gomock.Any()).
			Return(resp, nil)

		token, err := svc.Login(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, "token456", token)
	})

	t.Run("rpc error", func(t *testing.T) {
		mockClient.EXPECT().
			Login(ctx, gomock.Any()).
			Return(nil, errors.New("server unavailable"))

		token, err := svc.Login(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "login error")
		assert.Empty(t, token)
	})

	t.Run("unsuccessful response", func(t *testing.T) {
		resp := &pbModel.LoginResponse{}
		resp.SetMessage("error")
		resp.SetToken("")
		resp.SetSuccess(false)
		mockClient.EXPECT().
			Login(ctx, gomock.Any()).
			Return(resp, nil)

		token, err := svc.Login(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "login failed")
		assert.Empty(t, token)
	})
}
