package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/adapter/hubitat"
	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/state"
)

var hubitatCmd = &cobra.Command{
	Use:   "hubitat",
	Short: "Run HTTP server for Hubitat integration",
	Long: `Starts an HTTP server that exposes Eight Sleep Pod status and controls
for integration with Hubitat Elevation via custom Groovy drivers.

The server provides endpoints for:
  - GET /status - Full device status
  - GET /{side}/status - Status for left or right side
  - PUT /{side}/on - Turn on a side
  - PUT /{side}/off - Turn off a side
  - PUT /{side}/temperature?level=N - Set temperature level (-100 to 100)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}

		cl := client.New(
			viper.GetString("email"),
			viper.GetString("password"),
			viper.GetString("user_id"),
			viper.GetString("client_id"),
			viper.GetString("client_secret"),
		)

		ctx := context.Background()

		deviceID, err := cl.EnsureDeviceID(ctx)
		if err != nil {
			return fmt.Errorf("failed to get device ID: %w", err)
		}

		pollInterval := viper.GetDuration("hubitat.poll-interval")
		mgr := state.NewManager(cl, deviceID, state.WithCacheTTL(pollInterval))

		port := viper.GetInt("hubitat.port")
		adapter := hubitat.New(mgr, port, pollInterval)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		if err := adapter.Start(ctx); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}

		fmt.Printf("Hubitat server listening on port %d\n", port)

		<-sigChan
		fmt.Println("\nShutting down...")

		if err := adapter.Stop(); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(hubitatCmd)

	hubitatCmd.Flags().Int("port", 8080, "HTTP server port")
	hubitatCmd.Flags().Duration("poll-interval", 30*time.Second, "State polling interval")

	viper.BindPFlag("hubitat.port", hubitatCmd.Flags().Lookup("port"))
	viper.BindPFlag("hubitat.poll-interval", hubitatCmd.Flags().Lookup("poll-interval"))
}
