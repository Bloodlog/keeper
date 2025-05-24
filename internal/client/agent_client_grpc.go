package client

import (
	"fmt"
	"keeper/internal/config"
	pb "keeper/internal/proto/v1"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcVaultClient struct {
	pb.DataServiceClient
	conn *grpc.ClientConn
}

func (dc *GrpcVaultClient) Close() error {
	err := dc.conn.Close()
	if err != nil {
		return fmt.Errorf("close grpc client: %w", err)
	}
	return nil
}

func getGrpcDialOptions(cfg *config.RemoteServer) ([]grpc.DialOption, error) {
	if cfg.EnableTLS {
		creds, err := credentials.NewClientTLSFromFile(cfg.CACert, "")
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
		}
		return []grpc.DialOption{grpc.WithTransportCredentials(creds)}, nil
	}
	return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
}

func NewGrpcVaultClient(cfg *config.MainAgentConfig) (*GrpcVaultClient, error) {
	opts, err := getGrpcDialOptions(&cfg.RemoteServer)
	if err != nil {
		return nil, err
	}

	grpcAddress := fmt.Sprintf("%s:%d", cfg.RemoteServer.Address, cfg.RemoteServer.Port)

	conn, err := grpc.NewClient(grpcAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new client: %w", err)
	}

	client := pb.NewDataServiceClient(conn)

	return &GrpcVaultClient{
		conn:              conn,
		DataServiceClient: client,
	}, nil
}

type GrpcAuthClient struct {
	pb.AuthServiceClient
	conn *grpc.ClientConn
}

func (dc *GrpcAuthClient) Close() error {
	err := dc.conn.Close()
	if err != nil {
		return fmt.Errorf("close grpc client: %w", err)
	}
	return nil
}

func NewGrpcAuthClient(cfg *config.MainAgentConfig) (*GrpcAuthClient, error) {
	opts, err := getGrpcDialOptions(&cfg.RemoteServer)
	if err != nil {
		return nil, err
	}

	grpcAddress := fmt.Sprintf("%s:%d", cfg.RemoteServer.Address, cfg.RemoteServer.Port)

	conn, err := grpc.NewClient(grpcAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new client: %w", err)
	}

	client := pb.NewAuthServiceClient(conn)

	return &GrpcAuthClient{
		conn:              conn,
		AuthServiceClient: client,
	}, nil
}
