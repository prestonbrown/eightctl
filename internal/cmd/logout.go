package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/steipete/eightctl/internal/tokencache"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear cached authentication token",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := tokencache.Clear(); err != nil {
			return fmt.Errorf("clear token: %w", err)
		}
		fmt.Println("Logged out (token cache cleared)")
		return nil
	},
}
