package handler

import (
	"errors"
	"keeper/internal/dto"
	"keeper/internal/entity"
	"keeper/internal/logger"
	"keeper/internal/proto/v1/model"
	mocks "keeper/internal/service/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func mockAuthHandler(mockAuthService *mocks.MockAuthService) *AuthServerHandler {
	zapLog, _ := logger.NewZapLogger(zap.InfoLevel)
	return NewAuthHandler(zapLog, mockAuthService)
}

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)
	h := mockAuthHandler(mockAuthService)
	ctx := t.Context()
	expectedToken := entity.AccessToken{Token: "mocked-token"}
	mockAuthService.
		EXPECT().
		Register(ctx, getRegisterUserDto()).
		Return(expectedToken, nil)
	resp, err := h.Register(ctx, getRegisterDto())

	require.NoError(t, err)
	require.Equal(t, true, resp.GetSuccess())
	require.Equal(t, "Register successful.", resp.GetMessage())
	require.Equal(t, "mocked-token", resp.GetToken())
}

func TestRegister_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)
	h := mockAuthHandler(mockAuthService)
	ctx := t.Context()
	mockAuthService.
		EXPECT().
		Register(ctx, getRegisterUserDto()).
		Return(entity.AccessToken{}, errors.New("something went wrong"))
	resp, err := h.Register(ctx, getRegisterDto())

	require.Nil(t, resp)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to register user")
}

func getRegisterUserDto() *dto.RegisterUser {
	return &dto.RegisterUser{Login: "user", Password: "pass"}
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	h := mockAuthHandler(mockAuthService)
	ctx := t.Context()
	expectedToken := entity.AccessToken{Token: "login-token"}
	mockAuthService.
		EXPECT().
		Login(ctx, &dto.LoginUser{Login: "user", Password: "pass"}).
		Return(expectedToken, nil)
	resp, err := h.Login(ctx, getLoginDto())

	require.NoError(t, err)
	require.Equal(t, true, resp.GetSuccess())
	require.Equal(t, "Login successful.", resp.GetMessage())
	require.Equal(t, "login-token", resp.GetToken())
}

func getRegisterDto() *model.RegisterRequest {
	reqRegister := &model.RegisterRequest{}
	reqRegister.SetLogin("user")
	reqRegister.SetPassword("pass")
	return reqRegister
}

func getLoginDto() *model.LoginRequest {
	reqLogin := &model.LoginRequest{}
	reqLogin.SetLogin("user")
	reqLogin.SetPassword("pass")
	return reqLogin
}
