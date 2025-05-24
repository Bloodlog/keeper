package agent

import (
	"context"
	"errors"
	"fmt"
	"keeper/internal/service"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a secret by path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString(flagPath)
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

		return runWithVaultService(func(vault service.RemoteVaultService, timeout time.Duration) error {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			err := vault.DeleteSecret(ctx, token, path)
			if err != nil {
				return fmt.Errorf("failed to delete secret: %w", err)
			}

			fmt.Printf("🗑️  Secret deleted from: %s\n", path)
			return nil
		})
	},
}

func init() {
	deleteCmd.Flags().String(flagPath, "", "Path of the secret to delete")
	deleteCmd.Flags().String(flagToken, "", flagTokenDescription)
	deleteCmd.Flags().String(flagTokenFile, defaultTokenFile, flagTokenFileDescription)

	_ = deleteCmd.MarkFlagRequired(flagKeyName)
}
