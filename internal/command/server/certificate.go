package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	serialNumber       = 1658
	ipv4Address        = 127
	ipv6Address        = "0:0:0:0:0:0:0:1"
	certValidity       = 10
	keySize            = 4096
	certDir            = "cert"
	certFile           = "cert/public.cert"
	privKeyFile        = "cert/private.cert"
	dirPerm            = 0o750
	certFilePerm       = 0o600
	privateKeyFilePerm = 0o600
)

func generateCert() error {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(serialNumber),
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(ipv4Address, 0, 0, 1), net.ParseIP(ipv6Address)},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(certValidity, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	var certPEM bytes.Buffer
	if err := pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return fmt.Errorf("failed to encode certificate: %w", err)
	}

	var privateKeyPEM bytes.Buffer
	if err := pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}

	err = os.MkdirAll(certDir, dirPerm)
	if err != nil {
		return fmt.Errorf("failed to create cert directory: %w", err)
	}

	err = os.WriteFile(certFile, certPEM.Bytes(), certFilePerm)
	if err != nil {
		return fmt.Errorf("failed to write cert file: %w", err)
	}

	err = os.WriteFile(privKeyFile, privateKeyPEM.Bytes(), privateKeyFilePerm)
	if err != nil {
		return fmt.Errorf("failed to write private key file: %w", err)
	}

	log.Printf("TLS certificate and private key generated:\n - cert: %s\n - key: %s\n", certFile, privKeyFile)
	return nil
}

func genCertCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gen-cert",
		Short: "Generate self-signed TLS certificate and key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateCert()
		},
	}
}
