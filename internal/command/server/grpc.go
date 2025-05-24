package server

import (
	"context"
	"fmt"
	"keeper/internal/config"
	"keeper/internal/handler"
	"keeper/internal/logger"
	pb "keeper/internal/proto/v1"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc/credentials"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Init GRPC server.
func initGRPCServer(
	ctx context.Context,
	g *errgroup.Group,
	cfg *config.MainServerConfig,
	l *logger.ZapLogger,
	authHandler *handler.AuthServerHandler,
	vaultHandler *handler.VaultServerHandler,
) {
	var grpcServer *grpc.Server

	if cfg.GrpcServerConfig.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(cfg.GrpcServerConfig.CertFile, cfg.GrpcServerConfig.KeyFile)
		if err != nil {
			log.Fatalf("failed to load TLS credentials: %v", err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds))
		l.InfoCtx(ctx, "TLS is enabled for gRPC server")
	} else {
		grpcServer = grpc.NewServer()
		l.InfoCtx(ctx, "TLS is disabled for gRPC server")
	}

	g.Go(func() (err error) {
		lis, err := net.Listen("tcp", cfg.GrpcServerConfig.Address+":"+strconv.Itoa(cfg.GrpcServerConfig.Port))
		if err != nil {
			return fmt.Errorf("failed to run gRPC server: %w", err)
		}
		l.InfoCtx(ctx, "gRPC server listening on: "+cfg.GrpcServerConfig.Address+":"+strconv.Itoa(cfg.GrpcServerConfig.Port))

		pb.RegisterAuthServiceServer(grpcServer, authHandler)
		pb.RegisterDataServiceServer(grpcServer, vaultHandler)

		reflection.Register(grpcServer)
		err = grpcServer.Serve(lis)
		if err != nil {
			return fmt.Errorf("failed to run gRPC server: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		defer log.Print("grpc server has been shutdown")
		<-ctx.Done()

		if grpcServer != nil {
			grpcServer.GracefulStop()
		}

		log.Print("gRPC server has been shutdown")
		return nil
	})
}
