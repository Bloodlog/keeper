package server

import (
	"keeper/internal/config"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func bindFlags(cfg *config.MainServerConfig, cmd *cobra.Command) {
	for _, name := range []string{"address", "port", "grpc-address", "grpc-port", "dsn"} {
		if err := viper.BindPFlag(name, cmd.Flags().Lookup(name)); err != nil {
			log.Fatalf("failed to bind flag '%s': %v", name, err)
		}
	}

	cfg.Server.Address = viper.GetString("address")
	cfg.Server.Port = viper.GetInt("port")
	cfg.GrpcServerConfig.Address = viper.GetString("grpc-address")
	cfg.GrpcServerConfig.Port = viper.GetInt("grpc-port")
	cfg.Database.DSN = viper.GetString("dsn")
	cfg.GrpcServerConfig.Address = "0.0.0.0" // жёстко задано
}
