package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/output"
)

var sidesCmd = &cobra.Command{
	Use:   "sides",
	Short: "Show user assignments for each side of the bed",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))

		device, err := cl.Device().GetWithUsers(context.Background())
		if err != nil {
			return err
		}

		rows := []map[string]any{
			{"side": "left", "user_id": device.LeftUserID},
			{"side": "right", "user_id": device.RightUserID},
		}

		fields := viper.GetStringSlice("fields")
		rows = output.FilterFields(rows, fields)
		headers := fields
		if len(headers) == 0 {
			headers = []string{"side", "user_id"}
		}
		return output.Print(output.Format(viper.GetString("output")), headers, rows)
	},
}

func init() {
	rootCmd.AddCommand(sidesCmd)
}
