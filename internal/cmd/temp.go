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

		durationStr := viper.GetString("duration")
		if durationStr != "" {
			minutes, err := parseDuration(durationStr)
			if err != nil {
				return err
			}
			if err := cl.SetTemperatureWithDuration(context.Background(), lvl, minutes); err != nil {
				return err
			}
			fmt.Printf("temperature set (level %d) for %d minutes\n", lvl, minutes)
			return nil
		}

		if err := cl.SetTemperature(context.Background(), lvl); err != nil {
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
}
