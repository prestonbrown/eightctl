package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/output"
)

var bedtimeCmd = &cobra.Command{
	Use:   "bedtime",
	Short: "Manage bedtime temperature schedule",
}

var bedtimeShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current bedtime schedule",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		schedule, err := cl.Bedtime().Get(context.Background())
		if err != nil {
			return err
		}
		rows := []map[string]any{{
			"enabled": schedule.Enabled,
			"start":   schedule.Time.Start,
			"end":     schedule.Time.End,
			"stages":  len(schedule.Stages),
		}}
		rows = output.FilterFields(rows, viper.GetStringSlice("fields"))
		return output.Print(output.Format(viper.GetString("output")), []string{"enabled", "start", "end", "stages"}, rows)
	},
}

var bedtimeEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable bedtime schedule",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if err := cl.Bedtime().Enable(context.Background()); err != nil {
			return err
		}
		fmt.Println("bedtime schedule enabled")
		return nil
	},
}

var bedtimeDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable bedtime schedule",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if err := cl.Bedtime().Disable(context.Background()); err != nil {
			return err
		}
		fmt.Println("bedtime schedule disabled")
		return nil
	},
}

func init() {
	bedtimeCmd.AddCommand(bedtimeShowCmd, bedtimeEnableCmd, bedtimeDisableCmd)
	rootCmd.AddCommand(bedtimeCmd)
}
