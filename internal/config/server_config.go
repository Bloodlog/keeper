package config

import "time"

type MainServerConfig struct {
	GrpcServerConfig  GrpcServerConfig
	BuildAgentsConfig BuildAgentsConfig
	Database          DatabaseConfig
	Server            HTTPServerConfig
	Security          SecurityConfig
}

type BuildAgentsConfig struct {
	DownloadDir string
	URLPrefix   string
}

type HTTPServerConfig struct {
	Address string
	Port    int
}

type GrpcServerConfig struct {
	Address   string
	CertFile  string
	KeyFile   string
	Port      int
	EnableTLS bool
}

type DatabaseConfig struct {
	DSN string
}

type SecurityConfig struct {
	EncryptionKey     string
	DataEncryptionKey string
	EnableTLS         bool
	EnableCompression bool
	TokenTTL          time.Duration
}

func NewServerConfig() *MainServerConfig {
	const (
		httpAddress       = "127.0.0.1"
		httpPORT          = 8080
		grpcAddress       = "0.0.0.0"
		grpcPORT          = 8081
		dsn               = "postgres://keeper:password@localhost:5432/keeper?sslmode=disable"
		downloadDir       = "./build/clients"
		URLPrefix         = "/downloads/"
		secret            = "secret"
		dataEncryptionKey = "2fd36a2c3bcd3426f0fc92c84f8c56c1e91b40e372e3f1b739b1c1b0fa6fc457"
		maxTokenTTL       = 24 * time.Hour
	)

	return &MainServerConfig{
		Server: HTTPServerConfig{
			Address: httpAddress,
			Port:    httpPORT,
		},
		GrpcServerConfig: GrpcServerConfig{
			Address:   grpcAddress,
			Port:      grpcPORT,
			EnableTLS: false,
			CertFile:  "",
			KeyFile:   "",
		},
		Database: DatabaseConfig{
			DSN: dsn,
		},
		Security: SecurityConfig{
			EnableTLS:         false,
			EnableCompression: false,
			EncryptionKey:     secret,
			TokenTTL:          maxTokenTTL,
			DataEncryptionKey: dataEncryptionKey,
		},
		BuildAgentsConfig: BuildAgentsConfig{
			DownloadDir: downloadDir,
			URLPrefix:   URLPrefix,
		},
	}
}
