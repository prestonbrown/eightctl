package cmd

import (
	"context"
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/output"
)

// ErrNoAdjustableBase indicates the user doesn't have an adjustable base.
var ErrNoAdjustableBase = errors.New("no adjustable base found for this account (base commands require an Eight Sleep adjustable base)")

// isNoBaseError checks if an error indicates no adjustable base is available.
func isNoBaseError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Device not found") ||
		strings.Contains(msg, "PodOffline") ||
		strings.Contains(msg, "not found")
}

var baseCmd = &cobra.Command{Use: "base", Short: "Adjustable base controls"}

var baseInfoCmd = &cobra.Command{Use: "info", Short: "Show base status", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	res, err := cl.Base().Info(context.Background())
	if err != nil {
		if isNoBaseError(err) {
			return ErrNoAdjustableBase
		}
		return err
	}
	return output.Print(output.Format(viper.GetString("output")), []string{"info"}, []map[string]any{{"info": res}})
}}

var baseAngleCmd = &cobra.Command{Use: "angle", Short: "Set head/foot angle", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	head := viper.GetInt("head")
	foot := viper.GetInt("foot")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	if err := cl.Base().SetAngle(context.Background(), head, foot); err != nil {
		if isNoBaseError(err) {
			return ErrNoAdjustableBase
		}
		return err
	}
	return nil
}}

var basePresetsCmd = &cobra.Command{Use: "presets", Short: "List available presets", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	res, err := cl.Base().Presets(context.Background())
	if err != nil {
		if isNoBaseError(err) {
			return ErrNoAdjustableBase
		}
		return err
	}
	return output.Print(output.Format(viper.GetString("output")), []string{"presets"}, []map[string]any{{"presets": res}})
}}

var basePresetRunCmd = &cobra.Command{Use: "preset-run", Short: "Run a preset", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	name := viper.GetString("name")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	if err := cl.Base().RunPreset(context.Background(), name); err != nil {
		if isNoBaseError(err) {
			return ErrNoAdjustableBase
		}
		return err
	}
	return nil
}}

var baseTestCmd = &cobra.Command{Use: "test", Short: "Run vibration test", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	if err := cl.Base().VibrationTest(context.Background()); err != nil {
		if isNoBaseError(err) {
			return ErrNoAdjustableBase
		}
		return err
	}
	return nil
}}

var baseStopCmd = &cobra.Command{Use: "stop", Short: "Stop base movement", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	if err := cl.Base().StopMovement(context.Background()); err != nil {
		if isNoBaseError(err) {
			return ErrNoAdjustableBase
		}
		return err
	}
	return nil
}}

func init() {
	baseAngleCmd.Flags().Int("head", 0, "head angle")
	baseAngleCmd.Flags().Int("foot", 0, "foot angle")
	viper.BindPFlag("head", baseAngleCmd.Flags().Lookup("head"))
	viper.BindPFlag("foot", baseAngleCmd.Flags().Lookup("foot"))
	basePresetRunCmd.Flags().String("name", "", "preset name")
	viper.BindPFlag("name", basePresetRunCmd.Flags().Lookup("name"))

	baseCmd.AddCommand(baseInfoCmd, baseAngleCmd, basePresetsCmd, basePresetRunCmd, baseTestCmd, baseStopCmd)
}
