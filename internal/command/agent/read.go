package agent

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"keeper/internal/service"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	permissionOutFile = 0o600
	flagOutFile       = "out-file"
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read a secret by path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString(flagKeyName)
		token, _ := cmd.Flags().GetString(flagToken)
		tokenFile, _ := cmd.Flags().GetString(flagTokenFile)
		outFile, _ := cmd.Flags().GetString(flagOutFile)

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

		return runWithVaultService(func(vault service.RemoteVaultService, timeout time.Duration) error {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			secret, err := vault.GetSecret(ctx, token, path)
			if err != nil {
				return fmt.Errorf("failed to read secret: %w", err)
			}

			var deletionTime string
			var destroyed bool

			if secret.DeletedAt != nil {
				deletionTime = secret.DeletedAt.Format(time.RFC3339)
				destroyed = true
			} else {
				deletionTime = "n/a"
				destroyed = false
			}

			createdAt := "n/a"
			if !secret.CreatedAt.IsZero() {
				createdAt = secret.CreatedAt.Format(time.RFC3339)
			}
			const separator1 = "%-16s %s\n"
			const separator2 = "%-16s %v\n"

			fmt.Println("====== Metadata ======")
			fmt.Printf(separator1, "Key", "Value")
			fmt.Printf(separator1, "---", "-----")
			fmt.Printf(separator1, "created_time", createdAt)
			fmt.Printf(separator1, "deletion_time", deletionTime)
			fmt.Printf(separator2, "destroyed", destroyed)
			fmt.Printf(separator2, "version", secret.Version)

			var encoded string
			if err := json.Unmarshal(secret.Payload, &encoded); err != nil {
				return fmt.Errorf("failed to unmarshal payload as string: %w", err)
			}

			decoded, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				return fmt.Errorf("failed to decode base64 payload: %w", err)
			}

			if outFile != "" {
				if err := os.WriteFile(outFile, decoded, permissionOutFile); err != nil {
					return fmt.Errorf("failed to write to file: %w", err)
				}
				fmt.Printf("\n✅ Secret written to file: %s\n", outFile)
				return nil
			}

			var data map[string]interface{}
			if err := json.Unmarshal(decoded, &data); err != nil {
				return fmt.Errorf("failed to decode JSON string payload: %w", err)
			}

			const separator3 = "%-12s %s\n"

			const separator4 = "%-12s %v\n"
			fmt.Println("\n====== Data ======")
			fmt.Printf(separator3, "Key", "Value")
			fmt.Printf(separator3, "---", "-----")
			for k, v := range data {
				fmt.Printf(separator4, k, v)
			}

			return nil
		})
	},
}

func init() {
	readCmd.Flags().String(flagKeyName, "", "Path of the secret to read")
	readCmd.Flags().String(flagToken, "", flagTokenDescription)
	readCmd.Flags().String(flagTokenFile, defaultTokenFile, flagTokenFileDescription)
	readCmd.Flags().String(flagOutFile, "", "Optional path to write the secret payload to a file")

	_ = readCmd.MarkFlagRequired(flagKeyName)
}
