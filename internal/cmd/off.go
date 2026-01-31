package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/model"
)

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn pod off",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))

		ctx := context.Background()
		sideStr := viper.GetString("off_side")

		if sideStr != "" {
			side, err := model.ParseSide(sideStr)
			if err != nil {
				return err
			}
			userID, err := getUserIDForSide(ctx, cl, side)
			if err != nil {
				return err
			}
			if err := cl.TurnOffUser(ctx, userID); err != nil {
				return err
			}
			fmt.Printf("pod turned off for %s side\n", side)
			return nil
		}

		if err := cl.TurnOff(ctx); err != nil {
			return err
		}
		fmt.Println("pod turned off")
		return nil
	},
}

func init() {
	offCmd.Flags().String("side", "", "Turn off specific side (left or right)")
	viper.BindPFlag("off_side", offCmd.Flags().Lookup("side"))
}
