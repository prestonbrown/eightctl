// Package mqtt provides Home Assistant MQTT integration for Eight Sleep Pods.
package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/steipete/eightctl/internal/adapter"
	"github.com/steipete/eightctl/internal/model"
	"github.com/steipete/eightctl/internal/state"
)

// Config holds MQTT adapter configuration.
type Config struct {
	BrokerURL    string        // e.g., "tcp://localhost:1883"
	TopicPrefix  string        // e.g., "homeassistant" for HA discovery
	DeviceID     string        // Eight Sleep device ID
	DeviceName   string        // Human-readable name like "Bedroom Pod"
	PollInterval time.Duration // How often to poll state
	ClientID     string        // MQTT client ID
	Username     string        // Optional MQTT username
	Password     string        // Optional MQTT password
}

// Adapter implements the adapter.Adapter interface for MQTT/Home Assistant.
type Adapter struct {
	cfg          Config
	stateManager *state.Manager
	client       mqtt.Client
	stopCh       chan struct{}
	wg           sync.WaitGroup
}

// Compile-time check that Adapter implements adapter.Adapter.
var _ adapter.Adapter = (*Adapter)(nil)

// New creates a new MQTT adapter.
func New(cfg Config, stateManager *state.Manager) *Adapter {
	return &Adapter{
		cfg:          cfg,
		stateManager: stateManager,
		stopCh:       make(chan struct{}),
	}
}

// Start connects to the MQTT broker, publishes discovery configs, and starts polling.
func (a *Adapter) Start(ctx context.Context) error {
	// Configure MQTT client options
	opts := mqtt.NewClientOptions().
		AddBroker(a.cfg.BrokerURL).
		SetClientID(a.cfg.ClientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetOnConnectHandler(a.onConnect).
		SetConnectionLostHandler(a.onConnectionLost)

	if a.cfg.Username != "" {
		opts.SetUsername(a.cfg.Username)
	}
	if a.cfg.Password != "" {
		opts.SetPassword(a.cfg.Password)
	}

	// Create and connect client
	a.client = mqtt.NewClient(opts)
	token := a.client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", err)
	}

	// Publish discovery configs
	if err := a.publishDiscovery(); err != nil {
		return fmt.Errorf("failed to publish discovery configs: %w", err)
	}

	// Publish initial state
	if err := a.publishState(ctx); err != nil {
		return fmt.Errorf("failed to publish initial state: %w", err)
	}

	// Subscribe to command topics
	if err := a.subscribeCommands(); err != nil {
		return fmt.Errorf("failed to subscribe to command topics: %w", err)
	}

	// Publish online status
	a.publishAvailability("online")

	// Start polling goroutine
	a.wg.Add(1)
	go a.pollLoop(ctx)

	return nil
}

// HandleCommand processes a command from the smart home platform.
func (a *Adapter) HandleCommand(ctx context.Context, cmd adapter.Command) error {
	switch cmd.Action {
	case adapter.ActionOn:
		return a.stateManager.TurnOn(ctx, cmd.Side)
	case adapter.ActionOff:
		return a.stateManager.TurnOff(ctx, cmd.Side)
	case adapter.ActionSetTemp:
		if cmd.Temperature == nil {
			return fmt.Errorf("temperature required for set_temperature action")
		}
		return a.stateManager.SetTemperature(ctx, cmd.Side, *cmd.Temperature)
	default:
		return fmt.Errorf("unknown action: %s", cmd.Action)
	}
}

// Stop gracefully shuts down the adapter.
func (a *Adapter) Stop() error {
	// Signal polling goroutine to stop
	close(a.stopCh)
	a.wg.Wait()

	if a.client != nil && a.client.IsConnected() {
		// Publish offline status
		a.publishAvailability("offline")

		// Unsubscribe from command topics
		a.unsubscribeCommands()

		// Disconnect with 1 second timeout
		a.client.Disconnect(1000)
	}

	return nil
}

// onConnect is called when the MQTT connection is established.
func (a *Adapter) onConnect(_ mqtt.Client) {
	// Re-publish discovery and re-subscribe on reconnect
	_ = a.publishDiscovery()
	_ = a.subscribeCommands()
	a.publishAvailability("online")
}

// onConnectionLost is called when the MQTT connection is lost.
func (a *Adapter) onConnectionLost(_ mqtt.Client, _ error) {
	// Auto-reconnect is enabled, so we just wait for reconnection
}

// publishDiscovery publishes Home Assistant MQTT discovery configs for both sides.
func (a *Adapter) publishDiscovery() error {
	configs := GenerateDiscoveryConfigs(a.cfg.TopicPrefix, a.cfg.DeviceID, a.cfg.DeviceName)

	for side, config := range configs {
		topic := DiscoveryTopic(a.cfg.TopicPrefix, a.cfg.DeviceID, side)
		payload, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal discovery config for %s: %w", side, err)
		}

		token := a.client.Publish(topic, 1, true, payload) // QoS 1, retained
		token.Wait()
		if err := token.Error(); err != nil {
			return fmt.Errorf("failed to publish discovery config for %s: %w", side, err)
		}
	}

	return nil
}

// publishState fetches current state and publishes to state topics.
func (a *Adapter) publishState(ctx context.Context) error {
	deviceState, err := a.stateManager.GetState(ctx)
	if err != nil {
		return fmt.Errorf("failed to get device state: %w", err)
	}

	// Publish state for each side
	sides := []struct {
		name string
		side model.Side
		user *model.UserState
	}{
		{"left", model.Left, deviceState.LeftUser},
		{"right", model.Right, deviceState.RightUser},
	}

	for _, s := range sides {
		if s.user == nil {
			continue
		}

		// Publish temperature level
		tempTopic := fmt.Sprintf("eightsleep/%s/%s/temperature", a.cfg.DeviceID, s.name)
		a.publish(tempTopic, strconv.Itoa(s.user.TargetLevel))

		// Publish mode
		modeTopic := fmt.Sprintf("eightsleep/%s/%s/mode", a.cfg.DeviceID, s.name)
		mode := a.powerStateToMode(s.user.State, s.user.TargetLevel)
		a.publish(modeTopic, mode)

		// Publish current bed temperature
		currentTempTopic := fmt.Sprintf("eightsleep/%s/%s/current_temperature", a.cfg.DeviceID, s.name)
		a.publish(currentTempTopic, fmt.Sprintf("%.1f", s.user.BedTemperature))
	}

	return nil
}

// powerStateToMode converts PowerState and level to MQTT mode string.
func (a *Adapter) powerStateToMode(state model.PowerState, level int) string {
	switch state {
	case model.PowerOff:
		return "off"
	case model.PowerSmart, model.PowerManual:
		if level >= 0 {
			return "heat"
		}
		return "cool"
	default:
		return "off"
	}
}

// subscribeCommands subscribes to command topics for both sides.
func (a *Adapter) subscribeCommands() error {
	sides := []string{"left", "right"}

	for _, side := range sides {
		// Subscribe to temperature commands
		tempTopic := fmt.Sprintf("eightsleep/%s/%s/set_temperature", a.cfg.DeviceID, side)
		if err := a.subscribe(tempTopic, a.handleTemperatureCommand(side)); err != nil {
			return err
		}

		// Subscribe to mode commands
		modeTopic := fmt.Sprintf("eightsleep/%s/%s/set_mode", a.cfg.DeviceID, side)
		if err := a.subscribe(modeTopic, a.handleModeCommand(side)); err != nil {
			return err
		}
	}

	return nil
}

// unsubscribeCommands unsubscribes from all command topics.
func (a *Adapter) unsubscribeCommands() {
	sides := []string{"left", "right"}
	topics := make([]string, 0, 4)

	for _, side := range sides {
		topics = append(topics,
			fmt.Sprintf("eightsleep/%s/%s/set_temperature", a.cfg.DeviceID, side),
			fmt.Sprintf("eightsleep/%s/%s/set_mode", a.cfg.DeviceID, side),
		)
	}

	token := a.client.Unsubscribe(topics...)
	token.Wait()
}

// handleTemperatureCommand returns a handler for temperature set commands.
func (a *Adapter) handleTemperatureCommand(sideName string) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		level, err := strconv.Atoi(strings.TrimSpace(string(msg.Payload())))
		if err != nil {
			return // Invalid payload, ignore
		}

		side, err := model.ParseSide(sideName)
		if err != nil {
			return
		}

		cmd := adapter.Command{
			Action:      adapter.ActionSetTemp,
			Side:        side,
			Temperature: &level,
		}

		// Use background context since MQTT handlers don't have one
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := a.HandleCommand(ctx, cmd); err != nil {
			log.Printf("[mqtt] error handling temperature command for %s: %v", sideName, err)
			return
		}

		// Publish updated state
		if err := a.publishState(ctx); err != nil {
			log.Printf("[mqtt] error publishing state after temperature command: %v", err)
		}
	}
}

// handleModeCommand returns a handler for mode set commands.
func (a *Adapter) handleModeCommand(sideName string) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		mode := strings.TrimSpace(strings.ToLower(string(msg.Payload())))

		side, err := model.ParseSide(sideName)
		if err != nil {
			return
		}

		var cmd adapter.Command
		cmd.Side = side

		switch mode {
		case "off":
			cmd.Action = adapter.ActionOff
		case "heat", "cool":
			cmd.Action = adapter.ActionOn
		default:
			return // Unknown mode, ignore
		}

		// Use background context since MQTT handlers don't have one
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := a.HandleCommand(ctx, cmd); err != nil {
			log.Printf("[mqtt] error handling mode command for %s: %v", sideName, err)
			return
		}

		// Publish updated state
		if err := a.publishState(ctx); err != nil {
			log.Printf("[mqtt] error publishing state after mode command: %v", err)
		}
	}
}

// subscribe subscribes to a topic with the given handler.
func (a *Adapter) subscribe(topic string, handler mqtt.MessageHandler) error {
	token := a.client.Subscribe(topic, 1, handler) // QoS 1
	token.Wait()
	return token.Error()
}

// publish publishes a message to a topic.
func (a *Adapter) publish(topic, payload string) {
	token := a.client.Publish(topic, 1, true, payload) // QoS 1, retained
	token.Wait()
}

// publishAvailability publishes the availability status.
func (a *Adapter) publishAvailability(status string) {
	topic := fmt.Sprintf("eightsleep/%s/availability", a.cfg.DeviceID)
	a.publish(topic, status)
}

// pollLoop polls the state manager and publishes updates.
func (a *Adapter) pollLoop(ctx context.Context) {
	defer a.wg.Done()

	ticker := time.NewTicker(a.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-a.stopCh:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Invalidate cache to get fresh state
			a.stateManager.InvalidateCache()

			pollCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			if err := a.publishState(pollCtx); err != nil {
				log.Printf("[mqtt] error publishing state during poll: %v", err)
			}
			cancel()
		}
	}
}
