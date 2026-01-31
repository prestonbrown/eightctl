package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/tokencache"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Eight Sleep and cache credentials",
	Long: `Authenticate with Eight Sleep using email and password.

On success, the access token is cached in your system keychain (macOS Keychain,
Linux SecretService, or Windows Credential Manager) for future use.

You can provide credentials via flags, environment variables (EIGHTCTL_EMAIL,
EIGHTCTL_PASSWORD), or interactively when prompted.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		email := viper.GetString("email")
		password := viper.GetString("password")

		// Prompt for email if not provided
		if email == "" {
			fmt.Print("Email: ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("read email: %w", err)
			}
			email = strings.TrimSpace(input)
			if email == "" {
				return fmt.Errorf("email is required")
			}
		}

		// Prompt for password if not provided
		if password == "" {
			fmt.Print("Password: ")
			pwBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println() // newline after hidden input
			if err != nil {
				return fmt.Errorf("read password: %w", err)
			}
			password = string(pwBytes)
			if password == "" {
				return fmt.Errorf("password is required")
			}
		}

		cl := client.New(
			email,
			password,
			viper.GetString("user_id"),
			viper.GetString("client_id"),
			viper.GetString("client_secret"),
		)

		ctx := context.Background()
		if err := cl.Authenticate(ctx); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		// Load cached token to get expiration time
		cached, err := tokencache.Load(cl.Identity(), cl.UserID)
		if err != nil {
			// Authentication succeeded but couldn't read cache - still report success
			fmt.Printf("Logged in as %s\n", cl.UserID)
			return nil
		}

		remaining := time.Until(cached.ExpiresAt).Round(time.Minute)
		fmt.Printf("Logged in as %s (token expires in %s)\n", cl.UserID, remaining)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
