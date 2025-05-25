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
		destroy, _ := cmd.Flags().GetBool("destroy")
		metadata, _ := cmd.Flags().GetBool("metadata")
		undelete, _ := cmd.Flags().GetBool("undelete")
		version, _ := cmd.Flags().GetInt("version")

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

		if undelete && version <= 0 {
			return errors.New("--undelete requires --version to be set to a positive integer")
		}

		return runWithVaultService(func(vault service.RemoteVaultService, timeout time.Duration) error {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			switch {
			case metadata:
				if err := vault.DeleteMetadata(ctx, token, path); err != nil {
					return fmt.Errorf("failed to delete metadata: %w", err)
				}
				fmt.Printf("❌ Metadata deleted for secret: %s\n", path)

			case destroy:
				if err := vault.DestroySecret(ctx, token, path); err != nil {
					return fmt.Errorf("failed to destroy secret: %w", err)
				}
				fmt.Printf("🔥 Secret destroyed: %s\n", path)

			case undelete:
				if err := vault.UndeleteSecret(ctx, token, path, int64(version)); err != nil {
					return fmt.Errorf("failed to undelete version %d: %w", version, err)
				}
				fmt.Printf("♻️  Secret version %d restored at: %s\n", version, path)

			default:
				if err := vault.DeleteSecret(ctx, token, path); err != nil {
					return fmt.Errorf("failed to delete secret: %w", err)
				}
				fmt.Printf("🗑️  Secret deleted from: %s\n", path)
			}
			return nil
		})
	},
}

func init() {
	deleteCmd.Flags().String(flagPath, "", "Path of the secret to delete")
	deleteCmd.Flags().String(flagToken, "", flagTokenDescription)
	deleteCmd.Flags().String(flagTokenFile, defaultTokenFile, flagTokenFileDescription)

	deleteCmd.Flags().Bool("destroy", false, "Permanently destroy the secret")
	deleteCmd.Flags().Bool("metadata", false, "Delete metadata for the secret")
	deleteCmd.Flags().Bool("undelete", false, "Undelete a previously deleted secret version")
	deleteCmd.Flags().Int("version", 0, "Secret version to undelete")

	_ = deleteCmd.MarkFlagRequired(flagKeyName)
}
