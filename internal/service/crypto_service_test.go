package service

import (
	"encoding/hex"
	"keeper/internal/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCryptoService_EncodeDecode(t *testing.T) {
	key := "6368616e676520746869732070617373"
	cfg := config.SecurityConfig{
		DataEncryptionKey: key,
	}

	svc, err := NewCryptoService(cfg)
	require.NoError(t, err)

	original := []byte("my super secret data")

	encrypted, err := svc.Encode(original)
	require.NoError(t, err)
	require.NotEqual(t, original, encrypted)

	decrypted, err := svc.Decode(encrypted)
	require.NoError(t, err)
	require.Equal(t, original, decrypted)
}

func TestNewCryptoService_InvalidKey(t *testing.T) {
	cfg := config.SecurityConfig{
		DataEncryptionKey: "invalid-hex",
	}

	_, err := NewCryptoService(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decode data key")
}

func TestCryptoService_Encode_Error(t *testing.T) {
	key := hex.EncodeToString([]byte("bad"))

	cfg := config.SecurityConfig{
		DataEncryptionKey: key,
	}

	svc, err := NewCryptoService(cfg)
	require.NoError(t, err)

	_, err = svc.Encode([]byte("test payload"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to encrypt secret")
}

func TestCryptoService_Decode_Error(t *testing.T) {
	key := hex.EncodeToString([]byte("1234567890123456"))

	cfg := config.SecurityConfig{
		DataEncryptionKey: key,
	}

	svc, err := NewCryptoService(cfg)
	require.NoError(t, err)

	_, err = svc.Decode([]byte("not encrypted"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decrypt secret")
}
