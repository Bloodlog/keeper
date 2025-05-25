package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"keeper/internal/dto"
	"keeper/internal/service"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const flagKeyValue = "value"
const flagKeyDescription = "description"
const flagKeyMaxTTL = "max-ttl"

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Store a new secret",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString(flagKeyName)
		description, _ := cmd.Flags().GetString(flagKeyDescription)
		value, _ := cmd.Flags().GetString(flagKeyValue)
		expired, _ := cmd.Flags().GetInt(flagKeyMaxTTL)
		token, _ := cmd.Flags().GetString(flagToken)
		tokenFile, _ := cmd.Flags().GetString(flagTokenFile)

		if token == "" {
			token = os.Getenv(envAuthToken)
		}
		if token == "" && tokenFile != "" {
			data, err := os.ReadFile(tokenFile)
			if err != nil {
				return fmt.Errorf(errorReadTokenFile, tokenFile, err)
			}
			token = strings.TrimSpace(string(data))
		}

		if token == "" {
			return errors.New(errorTokenRequired)
		}

		if !json.Valid([]byte(value)) {
			return errors.New("value is not valid JSON")
		}

		ttl := expired
		if ttl <= 0 {
			ttl = 86400
		}

		expiredAt := time.Now().Add(time.Duration(ttl) * time.Second)

		return runWithVaultService(func(vault service.RemoteVaultService, timeout time.Duration) error {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			err := vault.SaveSecret(ctx, &dto.AgentCreateSecret{
				Token:       token,
				Path:        path,
				Description: description,
				Payload:     []byte(value),
				ExpiredAt:   expiredAt,
			})
			if err != nil {
				return fmt.Errorf("failed to store secret: %w", err)
			}

			fmt.Printf("✅ Success! Data written to: %s\n", path)
			return nil
		})
	},
}

func init() {
	writeCmd.Flags().String(flagKeyName, "", "Path to store secret under")
	writeCmd.Flags().String(flagKeyDescription, "", "Description of the secret")
	writeCmd.Flags().String(flagKeyValue, "", "Secret value (must be valid JSON)")
	writeCmd.Flags().Int(flagKeyMaxTTL, 0, "Optional TTL in seconds")
	writeCmd.Flags().String(flagToken, "", flagTokenDescription)
	writeCmd.Flags().String(
		flagTokenFile,
		defaultTokenFile,
		flagTokenFileDescription)

	_ = writeCmd.MarkFlagRequired(flagKeyName)
	_ = writeCmd.MarkFlagRequired(flagKeyValue)
}
