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

	"github.com/steipete/eightctl/internal/adapter/mqtt"
	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/state"
)

var mqttCmd = &cobra.Command{
	Use:   "mqtt",
	Short: "Run MQTT bridge for Home Assistant",
	Long: `Starts an MQTT bridge that publishes Eight Sleep Pod state to Home Assistant
via MQTT Discovery and subscribes to command topics for control.

The bridge:
  - Publishes MQTT Discovery configs for climate entities
  - Polls device state and publishes to state topics
  - Subscribes to command topics for temperature and mode changes

Requires an MQTT broker (e.g., Mosquitto) accessible from both this machine
and Home Assistant.`,
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

		pollInterval := viper.GetDuration("mqtt.poll-interval")
		mgr := state.NewManager(cl, deviceID, state.WithCacheTTL(pollInterval))

		cfg := mqtt.Config{
			BrokerURL:    viper.GetString("mqtt.broker"),
			TopicPrefix:  viper.GetString("mqtt.topic-prefix"),
			DeviceID:     deviceID,
			DeviceName:   viper.GetString("mqtt.device-name"),
			PollInterval: pollInterval,
			ClientID:     viper.GetString("mqtt.client-id"),
			Username:     viper.GetString("mqtt.mqtt-username"),
			Password:     viper.GetString("mqtt.mqtt-password"),
		}

		adapter := mqtt.New(cfg, mgr)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		if err := adapter.Start(ctx); err != nil {
			return fmt.Errorf("failed to start MQTT bridge: %w", err)
		}

		fmt.Printf("MQTT bridge connected to %s\n", cfg.BrokerURL)
		fmt.Printf("Publishing to %s discovery prefix\n", cfg.TopicPrefix)

		<-sigChan
		fmt.Println("\nShutting down...")

		if err := adapter.Stop(); err != nil {
			return fmt.Errorf("failed to stop MQTT bridge: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mqttCmd)

	mqttCmd.Flags().String("broker", "tcp://localhost:1883", "MQTT broker URL")
	mqttCmd.Flags().String("topic-prefix", "homeassistant", "Topic prefix for discovery")
	mqttCmd.Flags().String("device-name", "Eight Sleep Pod", "Device name in Home Assistant")
	mqttCmd.Flags().String("client-id", "eightctl", "MQTT client ID")
	mqttCmd.Flags().String("mqtt-username", "", "MQTT username (optional)")
	mqttCmd.Flags().String("mqtt-password", "", "MQTT password (optional)")
	mqttCmd.Flags().Duration("poll-interval", 30*time.Second, "State polling interval")

	viper.BindPFlag("mqtt.broker", mqttCmd.Flags().Lookup("broker"))
	viper.BindPFlag("mqtt.topic-prefix", mqttCmd.Flags().Lookup("topic-prefix"))
	viper.BindPFlag("mqtt.device-name", mqttCmd.Flags().Lookup("device-name"))
	viper.BindPFlag("mqtt.client-id", mqttCmd.Flags().Lookup("client-id"))
	viper.BindPFlag("mqtt.mqtt-username", mqttCmd.Flags().Lookup("mqtt-username"))
	viper.BindPFlag("mqtt.mqtt-password", mqttCmd.Flags().Lookup("mqtt-password"))
	viper.BindPFlag("mqtt.poll-interval", mqttCmd.Flags().Lookup("poll-interval"))
}
