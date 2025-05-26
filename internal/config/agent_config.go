package config

import "time"

type MainAgentConfig struct {
	RemoteServer RemoteServer
}

type RemoteServer struct {
	Address   string
	CACert    string
	Port      int
	Timeout   time.Duration
	EnableTLS bool
}

func NewAgentConfig() *MainAgentConfig {
	const (
		remoteGrpcAddress = "127.0.0.1"
		remoteGrpcPORT    = 8081
		defaultTimeout    = 5 * time.Second
	)
	return &MainAgentConfig{
		RemoteServer{
			Address: remoteGrpcAddress,
			Port:    remoteGrpcPORT,
			Timeout: defaultTimeout,
		},
	}
}
