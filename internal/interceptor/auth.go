package interceptor

import (
	"context"
	"keeper/internal/service"
	utils "keeper/internal/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(jwtService service.JwtService) grpc.UnaryServerInterceptor {
	skipAuth := map[string]bool{
		"/keeper.go.grpc.v1.AuthService/Login":    true,
		"/keeper.go.grpc.v1.AuthService/Register": true,
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if skipAuth[info.FullMethod] {
			return handler(ctx, req)
		}
		msg, ok := req.(interface {
			GetToken() string
		})
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "request does not contain a token")
		}

		userID, err := jwtService.GetUserID(msg.GetToken())
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		ctx = utils.SetUserID(ctx, userID)

		return handler(ctx, req)
	}
}
