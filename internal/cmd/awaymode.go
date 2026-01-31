package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/output"
)

var awayModeCmd = &cobra.Command{
	Use:   "away-mode",
	Short: "Manage away mode",
}

var awayModeShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current away mode status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		status, err := cl.AwayMode().Get(context.Background())
		if err != nil {
			return err
		}
		rows := []map[string]any{{
			"enabled": status.Enabled,
		}}
		rows = output.FilterFields(rows, viper.GetStringSlice("fields"))
		return output.Print(output.Format(viper.GetString("output")), []string{"enabled"}, rows)
	},
}

var awayModeOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Enable away mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if err := cl.AwayMode().Enable(context.Background()); err != nil {
			return err
		}
		fmt.Println("away mode enabled")
		return nil
	},
}

var awayModeOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Disable away mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if err := cl.AwayMode().Disable(context.Background()); err != nil {
			return err
		}
		fmt.Println("away mode disabled")
		return nil
	},
}

func init() {
	awayModeCmd.AddCommand(awayModeShowCmd, awayModeOnCmd, awayModeOffCmd)
	rootCmd.AddCommand(awayModeCmd)
}
