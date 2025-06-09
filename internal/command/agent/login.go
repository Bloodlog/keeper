package agent

import (
	"context"
	"fmt"
	"keeper/internal/dto"
	"keeper/internal/service"
	"time"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   flagLogin,
	Short: "Login user",
	RunE: func(cmd *cobra.Command, args []string) error {
		login, _ := cmd.Flags().GetString(flagLogin)
		password, _ := cmd.Flags().GetString(flagPassword)
		tokenFilePath, _ := cmd.Flags().GetString(flagTokenFile)

		return runWithAuthService(func(auth service.RemoteAuthService, timeout time.Duration) error {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			token, err := auth.Login(ctx, dto.LoginUser{Login: login, Password: password})
			if err != nil {
				return fmt.Errorf("login failed: %w", err)
			}
			err = saveTokenAndPrintInfo(token, tokenFilePath)
			if err != nil {
				return fmt.Errorf("save token failed: %w", err)
			}
			printSuccessInfoAboutLogin(tokenFilePath, "✅ Success! You are now logged in.")
			return nil
		})
	},
}

func init() {
	loginCmd.Flags().String(flagLogin, "", "User login")
	loginCmd.Flags().String(flagPassword, "", "User password")
	loginCmd.Flags().String(flagTokenFile, defaultTokenFile, "Path to token file")

	_ = loginCmd.MarkFlagRequired(flagLogin)
	_ = loginCmd.MarkFlagRequired(flagPassword)
}
