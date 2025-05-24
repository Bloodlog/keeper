package agent

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "keeper-agent",
	Short: "Run the GophKeeper agent",
}

const (
	flagGRPCPort    = "grpc-port"
	flagGRPCAddress = "grpc-address"
	flagEnableTLS   = "enable-tls"
	flagCACert      = "ca-cert"
	grpcPort        = 8081
)

func init() {
	rootCmd.PersistentFlags().String(flagGRPCAddress, "app-keeper", "gRPC server address")
	rootCmd.PersistentFlags().Int(flagGRPCPort, grpcPort, "gRPC server port")
	rootCmd.PersistentFlags().Bool(flagEnableTLS, false, "Enable TLS when connecting to the server")
	rootCmd.PersistentFlags().String(flagCACert, "cert/public.cert", "Path to CA certificate file")

	viper.SetEnvPrefix("KEEPER")
	viper.AutomaticEnv()

	_ = viper.BindPFlag(flagGRPCAddress, rootCmd.PersistentFlags().Lookup(flagGRPCAddress))
	_ = viper.BindPFlag(flagGRPCPort, rootCmd.PersistentFlags().Lookup(flagGRPCPort))
	_ = viper.BindPFlag(flagEnableTLS, rootCmd.PersistentFlags().Lookup(flagEnableTLS))
	_ = viper.BindPFlag(flagCACert, rootCmd.PersistentFlags().Lookup(flagCACert))

	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(writeCmd)
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
}

func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return fmt.Errorf("error agent: %w", err)
	}
	return nil
}
