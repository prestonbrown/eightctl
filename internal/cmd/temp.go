package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/daemon"
	"github.com/steipete/eightctl/internal/model"
)

var tempCmd = &cobra.Command{
	Use:   "temp <value>",
	Short: "Set pod temperature (e.g., 68F, 20C, or heating level -100..100)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		lvl, err := daemon.ParseTemp(args[0])
		if err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))

		ctx := context.Background()
		sideStr := viper.GetString("temp_side")
		durationStr := viper.GetString("duration")

		if sideStr != "" {
			if durationStr != "" {
				return fmt.Errorf("--duration is not supported with --side")
			}
			side, err := model.ParseSide(sideStr)
			if err != nil {
				return err
			}
			userID, err := getUserIDForSide(ctx, cl, side)
			if err != nil {
				return err
			}
			if err := cl.SetUserTemperature(ctx, userID, lvl); err != nil {
				return err
			}
			fmt.Printf("temperature set (level %d) for %s side\n", lvl, side)
			return nil
		}

		if durationStr != "" {
			minutes, err := parseDuration(durationStr)
			if err != nil {
				return err
			}
			if err := cl.SetTemperatureWithDuration(ctx, lvl, minutes); err != nil {
				return err
			}
			fmt.Printf("temperature set (level %d) for %d minutes\n", lvl, minutes)
			return nil
		}

		if err := cl.SetTemperature(ctx, lvl); err != nil {
			return err
		}
		fmt.Printf("temperature set (level %d)\n", lvl)
		return nil
	},
}

// parseDuration parses duration strings like "30m", "1h", "90" (minutes).
func parseDuration(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty duration")
	}

	// Try standard Go duration parsing first
	if d, err := time.ParseDuration(s); err == nil {
		return int(d.Minutes()), nil
	}

	// Try plain number (treated as minutes)
	if minutes, err := strconv.Atoi(s); err == nil {
		return minutes, nil
	}

	return 0, fmt.Errorf("invalid duration format: %s (use e.g., 30m, 1h, or plain minutes)", s)
}

func init() {
	tempCmd.Flags().String("duration", "", "Duration for temperature setting (e.g., 30m, 1h)")
	viper.BindPFlag("duration", tempCmd.Flags().Lookup("duration"))
	tempCmd.Flags().String("side", "", "Set temperature for specific side (left or right)")
	viper.BindPFlag("temp_side", tempCmd.Flags().Lookup("side"))
}
