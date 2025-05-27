package service

import (
	"encoding/hex"
	"fmt"
	"keeper/internal/config"
	"keeper/internal/security"
)

type CryptoService interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

type cryptoService struct {
	dataEncryptionKey []byte
}

func NewCryptoService(cfg config.SecurityConfig) (CryptoService, error) {
	key, err := hex.DecodeString(cfg.DataEncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode data key: %w", err)
	}

	return &cryptoService{
		dataEncryptionKey: key,
	}, nil
}

func (svc *cryptoService) Encode(payload []byte) ([]byte, error) {
	encrypted, err := security.EncryptAESGCM(payload, svc.dataEncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt secret: %w", err)
	}

	return encrypted, nil
}

func (svc *cryptoService) Decode(data []byte) ([]byte, error) {
	decrypted, err := security.DecryptAESGCM(svc.dataEncryptionKey, data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	}

	return decrypted, nil
}
