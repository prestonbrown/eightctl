package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/model"
)

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn pod on",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))

		ctx := context.Background()
		sideStr := viper.GetString("on_side")

		if sideStr != "" {
			side, err := model.ParseSide(sideStr)
			if err != nil {
				return err
			}
			userID, err := getUserIDForSide(ctx, cl, side)
			if err != nil {
				return err
			}
			if err := cl.TurnOnUser(ctx, userID); err != nil {
				return err
			}
			fmt.Printf("pod turned on for %s side\n", side)
			return nil
		}

		if err := cl.TurnOn(ctx); err != nil {
			return err
		}
		fmt.Println("pod turned on")
		return nil
	},
}

func init() {
	onCmd.Flags().String("side", "", "Turn on specific side (left or right)")
	viper.BindPFlag("on_side", onCmd.Flags().Lookup("side"))
}
