package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/model"
	"github.com/steipete/eightctl/internal/output"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show device status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))

		ctx := context.Background()
		var st *client.TempStatus
		var err error

		sideStr := viper.GetString("status_side")
		if sideStr != "" {
			side, err := model.ParseSide(sideStr)
			if err != nil {
				return err
			}
			userID, err := getUserIDForSide(ctx, cl, side)
			if err != nil {
				return err
			}
			st, err = cl.GetUserTemperature(ctx, userID)
			if err != nil {
				return err
			}
		} else {
			st, err = cl.GetStatus(ctx)
			if err != nil {
				return err
			}
		}

		row := map[string]any{"mode": st.CurrentState.Type, "level": st.CurrentLevel}
		fields := viper.GetStringSlice("fields")
		rows := output.FilterFields([]map[string]any{row}, fields)
		headers := fields
		if len(headers) == 0 {
			headers = []string{"mode", "level"}
		}
		return output.Print(output.Format(viper.GetString("output")), headers, rows)
	},
}

func init() {
	statusCmd.Flags().String("side", "", "Show status for specific side (left or right)")
	viper.BindPFlag("status_side", statusCmd.Flags().Lookup("side"))
}
