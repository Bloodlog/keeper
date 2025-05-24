package agent

import (
	"fmt"
	"keeper/internal/client"
	"keeper/internal/config"
	"keeper/internal/service"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	flagPath            = flagKeyName
	flagToken           = "token"
	flagTokenFile       = "token-file"
	flagKeyName         = "path"
	envAuthToken        = "TOKEN"
	permissionTokenFile = 0o600
	defaultTokenFile    = ".keeper-token"

	flagTokenDescription     = "User token (can also be set via TOKEN)"
	flagTokenFileDescription = "Path to file with token (used if --token and TOKEN are unset)"

	errorReadTokenFile = "failed to read token from file (%s): %w"
	errorTokenRequired = "token is required (--token, TOKEN, or --token-file)"

	flagLogin        = "login"
	flagPassword     = "password"
	errorConnectGrpc = "failed to connect to gRPC server: %w"

	flagGrpcAddress = "grpc-address"
	flagGrpcPort    = "grpc-port"
)

func runWithAuthService(action func(service.RemoteAuthService, time.Duration) error) error {
	grpcClient, cfg, err := initGrpcAuthClient()
	if err != nil {
		return fmt.Errorf(errorConnectGrpc, err)
	}
	defer func(grpcClient *client.GrpcAuthClient) {
		err := grpcClient.Close()
		if err != nil {
			fmt.Printf("failed to close gRPC client connection: %v", err)
		}
	}(grpcClient)

	auth := service.NewRemoteAuthService(grpcClient)
	return action(auth, cfg.RemoteServer.Timeout)
}

func runWithVaultService(action func(service.RemoteVaultService, time.Duration) error) error {
	grpcClient, cfg, err := initGrpcVaultClient()
	if err != nil {
		return fmt.Errorf(errorConnectGrpc, err)
	}
	defer func(grpcClient *client.GrpcVaultClient) {
		err := grpcClient.Close()
		if err != nil {
			fmt.Printf("failed to close gRPC client connection: %v", err)
		}
	}(grpcClient)

	auth := service.NewRemoteVaultService(grpcClient)
	return action(auth, cfg.RemoteServer.Timeout)
}

func initGrpcAuthClient() (*client.GrpcAuthClient, *config.MainAgentConfig, error) {
	cfg := config.NewAgentConfig()
	cfg.RemoteServer.Address = viper.GetString(flagGrpcAddress)
	cfg.RemoteServer.Port = viper.GetInt(flagGrpcPort)

	grpcClient, err := client.NewGrpcAuthClient(cfg)
	if err != nil {
		return nil, cfg, fmt.Errorf(errorConnectGrpc, err)
	}

	return grpcClient, cfg, nil
}

func initGrpcVaultClient() (*client.GrpcVaultClient, *config.MainAgentConfig, error) {
	cfg := config.NewAgentConfig()
	cfg.RemoteServer.Address = viper.GetString(flagGrpcAddress)
	cfg.RemoteServer.Port = viper.GetInt(flagGrpcPort)

	grpcClient, err := client.NewGrpcVaultClient(cfg)
	if err != nil {
		return nil, cfg, fmt.Errorf(errorConnectGrpc, err)
	}

	return grpcClient, cfg, nil
}

func saveTokenAndPrintInfo(token string, tokenFilePath string) error {
	if err := os.WriteFile(tokenFilePath, []byte(token), permissionTokenFile); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

func printSuccessInfoAboutLogin(tokenFilePath string, successMessage string) {
	fmt.Println(successMessage)
	fmt.Println("You do NOT need to run vault login again.")
	fmt.Println("Future Keeper requests will automatically use this token.")
	fmt.Printf("Token information will be saved to: %s\n", tokenFilePath)
}
