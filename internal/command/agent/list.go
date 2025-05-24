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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List secret paths",
	RunE: func(cmd *cobra.Command, args []string) error {
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

			paths, err := vault.ListSecretPaths(ctx, token)
			if err != nil {
				return fmt.Errorf("failed to list secrets: %w", err)
			}

			fmt.Println("Keys")
			fmt.Println("----")

			if len(paths) == 0 {
				fmt.Println("No secrets found.")
				return nil
			}

			for _, path := range paths {
				fmt.Println(path)
			}
			return nil
		})
	},
}

func init() {
	listCmd.Flags().String(flagToken, "", flagTokenDescription)
	listCmd.Flags().String(flagTokenFile, defaultTokenFile, flagTokenFileDescription)
}
