package agent

import (
	"context"
	"fmt"
	"keeper/internal/dto"
	"keeper/internal/service"
	"time"

	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	RunE: func(cmd *cobra.Command, args []string) error {
		login, _ := cmd.Flags().GetString(flagLogin)
		password, _ := cmd.Flags().GetString(flagPassword)
		tokenFilePath, _ := cmd.Flags().GetString("token-file")

		return runWithAuthService(func(auth service.RemoteAuthService, timeout time.Duration) error {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			token, err := auth.Register(ctx, dto.RegisterUser{Login: login, Password: password})
			if err != nil {
				return fmt.Errorf("registration failed: %w", err)
			}
			err = saveTokenAndPrintInfo(token, tokenFilePath)
			if err != nil {
				return fmt.Errorf("save token failed: %w", err)
			}
			printSuccessInfoAboutLogin(tokenFilePath, "✅ Registration successful.")
			return nil
		})
	},
}

func init() {
	registerCmd.Flags().String(flagLogin, "", "Login for new user")
	registerCmd.Flags().String(flagPassword, "", "Password for new user")
	registerCmd.Flags().String(flagTokenFile, defaultTokenFile, "Path to token file")

	_ = registerCmd.MarkFlagRequired(flagLogin)
	_ = registerCmd.MarkFlagRequired(flagPassword)
}
