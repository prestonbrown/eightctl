package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/output"
)

var alarmCmd = &cobra.Command{
	Use:   "alarm",
	Short: "Manage alarms",
}

var alarmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List alarms",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		alarms, err := cl.ListAlarms(context.Background())
		if err != nil {
			return err
		}
		rows := make([]map[string]any, 0, len(alarms))
		for _, a := range alarms {
			rows = append(rows, map[string]any{
				"id":        a.ID,
				"time":      a.Time,
				"enabled":   a.Enabled,
				"repeat":    a.Repeat.Enabled,
				"vibration": a.Vibration.Enabled,
				"thermal":   a.Thermal.Enabled,
				"snoozing":  a.Snoozing,
			})
		}
		rows = output.FilterFields(rows, viper.GetStringSlice("fields"))
		return output.Print(output.Format(viper.GetString("output")), []string{"id", "time", "enabled", "repeat", "vibration", "thermal"}, rows)
	},
}

var alarmCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an alarm",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		timeStr := viper.GetString("time")
		if timeStr == "" {
			return fmt.Errorf("--time required")
		}
		// Ensure time is in HH:MM:SS format
		if len(timeStr) == 5 {
			timeStr += ":00"
		}
		days := viper.GetIntSlice("days")
		weekDays := make(map[string]bool)
		dayNames := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
		for _, d := range days {
			if d >= 0 && d < 7 {
				weekDays[dayNames[d]] = true
			}
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		alarm := client.Alarm{
			Enabled: !viper.GetBool("disabled"),
			Time:    timeStr,
			Repeat: client.AlarmRepeat{
				Enabled:  len(days) > 0,
				WeekDays: weekDays,
			},
			Vibration: client.AlarmVibration{
				Enabled: !viper.GetBool("no-vibration"),
			},
			Thermal: client.AlarmThermal{
				Enabled: viper.GetBool("thermal"),
			},
		}
		res, err := cl.CreateAlarm(context.Background(), alarm)
		if err != nil {
			return err
		}
		fmt.Printf("created alarm %s\n", res.ID)
		return nil
	},
}

var alarmUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an alarm",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		patch := map[string]any{}
		if f := viper.GetString("time"); f != "" {
			timeStr := f
			if len(timeStr) == 5 {
				timeStr += ":00"
			}
			patch["time"] = timeStr
		}
		if days := viper.GetIntSlice("days"); len(days) > 0 {
			weekDays := make(map[string]bool)
			dayNames := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
			for _, d := range days {
				if d >= 0 && d < 7 {
					weekDays[dayNames[d]] = true
				}
			}
			patch["repeat"] = map[string]any{
				"enabled":  true,
				"weekDays": weekDays,
			}
		}
		if cmd.Flags().Changed("enabled") {
			patch["enabled"] = viper.GetBool("enabled")
		}
		if cmd.Flags().Changed("no-vibration") {
			patch["vibration"] = map[string]any{
				"enabled": !viper.GetBool("no-vibration"),
			}
		}
		if cmd.Flags().Changed("thermal") {
			patch["thermal"] = map[string]any{
				"enabled": viper.GetBool("thermal"),
			}
		}
		if len(patch) == 0 {
			return fmt.Errorf("no fields to update")
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if _, err := cl.UpdateAlarm(context.Background(), args[0], patch); err != nil {
			return err
		}
		fmt.Println("updated")
		return nil
	},
}

var alarmDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an alarm",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if err := cl.DeleteAlarm(context.Background(), args[0]); err != nil {
			return err
		}
		fmt.Println("deleted")
		return nil
	},
}

func init() {
	alarmCreateCmd.Flags().String("time", "", "HH:MM or HH:MM:SS time")
	alarmCreateCmd.Flags().IntSlice("days", nil, "Comma-separated days 0=Sun..6=Sat (for repeating)")
	alarmCreateCmd.Flags().Bool("disabled", false, "Create disabled")
	alarmCreateCmd.Flags().Bool("no-vibration", false, "Disable vibration")
	alarmCreateCmd.Flags().Bool("thermal", false, "Enable thermal wake")
	viper.BindPFlag("time", alarmCreateCmd.Flags().Lookup("time"))
	viper.BindPFlag("days", alarmCreateCmd.Flags().Lookup("days"))
	viper.BindPFlag("disabled", alarmCreateCmd.Flags().Lookup("disabled"))
	viper.BindPFlag("no-vibration", alarmCreateCmd.Flags().Lookup("no-vibration"))
	viper.BindPFlag("thermal", alarmCreateCmd.Flags().Lookup("thermal"))

	alarmUpdateCmd.Flags().String("time", "", "HH:MM or HH:MM:SS time")
	alarmUpdateCmd.Flags().IntSlice("days", nil, "Comma-separated days 0=Sun..6=Sat")
	alarmUpdateCmd.Flags().Bool("enabled", true, "Set enabled true/false")
	alarmUpdateCmd.Flags().Bool("no-vibration", false, "Disable vibration")
	alarmUpdateCmd.Flags().Bool("thermal", false, "Enable thermal wake")
	viper.BindPFlag("time", alarmUpdateCmd.Flags().Lookup("time"))
	viper.BindPFlag("days", alarmUpdateCmd.Flags().Lookup("days"))
	viper.BindPFlag("enabled", alarmUpdateCmd.Flags().Lookup("enabled"))
	viper.BindPFlag("no-vibration", alarmUpdateCmd.Flags().Lookup("no-vibration"))
	viper.BindPFlag("thermal", alarmUpdateCmd.Flags().Lookup("thermal"))

	// add subcommands
	alarmCmd.AddCommand(alarmListCmd, alarmCreateCmd, alarmUpdateCmd, alarmDeleteCmd, alarmSnoozeCmd, alarmDismissCmd, alarmDismissAllCmd, alarmVibeCmd)
}

// snooze
var alarmSnoozeCmd = &cobra.Command{Use: "snooze <id>", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Alarms().Snooze(context.Background(), args[0])
}}

var alarmDismissCmd = &cobra.Command{Use: "dismiss <id>", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Alarms().Dismiss(context.Background(), args[0])
}}

var alarmDismissAllCmd = &cobra.Command{Use: "dismiss-all", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Alarms().DismissAll(context.Background())
}}

var alarmVibeCmd = &cobra.Command{Use: "vibration-test", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	return cl.Alarms().VibrationTest(context.Background())
}}

// parseDays convenience to support comma inputs (unused, kept for future).
func parseDays(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	res := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		var v int
		if _, err := fmt.Sscanf(p, "%d", &v); err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}
