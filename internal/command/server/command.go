package server

import (
	"fmt"
	"keeper/internal/config"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	timeoutServerShutdown = time.Second * 5
	timeoutShutdown       = time.Second * 10
)

func Execute() error {
	cfg := config.NewServerConfig()

	cmd := &cobra.Command{
		Use:   "keeper-server",
		Short: "Run the GophKeeper server",
		RunE: func(cmd *cobra.Command, args []string) error {
			bindFlags(cfg, cmd)
			return run(cfg)
		},
	}

	// HTTP flags
	cmd.Flags().StringVar(&cfg.Server.Address, "address", cfg.Server.Address, "Address to bind the server")
	cmd.Flags().IntVar(&cfg.Server.Port, "port", cfg.Server.Port, "Port to bind the server")

	// gRPC flags
	cmd.Flags().StringVar(
		&cfg.GrpcServerConfig.Address,
		"grpc-address",
		cfg.GrpcServerConfig.Address,
		"Address to bind the gRPC server")
	cmd.Flags().IntVar(
		&cfg.GrpcServerConfig.Port,
		"grpc-port",
		cfg.GrpcServerConfig.Port,
		"Port to bind the gRPC server")

	// Database
	cmd.Flags().StringVar(&cfg.Database.DSN, "dsn", cfg.Database.DSN, "Database DSN")
	// TLS for gRPC
	cmd.Flags().BoolVar(
		&cfg.GrpcServerConfig.EnableTLS,
		"enable-tls", cfg.GrpcServerConfig.EnableTLS,
		"Enable TLS for gRPC server")
	cmd.Flags().StringVar(
		&cfg.GrpcServerConfig.CertFile,
		"cert-file", cfg.GrpcServerConfig.CertFile,
		"Path to TLS certificate file")
	cmd.Flags().StringVar(
		&cfg.GrpcServerConfig.KeyFile,
		"key-file", cfg.GrpcServerConfig.KeyFile,
		"Path to TLS private key file")

	viper.SetEnvPrefix("KEEPER")
	viper.AutomaticEnv()

	cmd.AddCommand(genCertCmd())

	err := cmd.Execute()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
