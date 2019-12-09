package main

import (
	"github.com/spf13/cobra"
	"github.com/steeve/itool/diagnostics_relay"
)

func init() {
	diagnosticsCmd.AddCommand(diagnosticsQueryCmd)
	diagnosticsCmd.AddCommand(diagnosticsShutdownCmd)
	diagnosticsCmd.AddCommand(diagnosticsRestartCmd)
	diagnosticsCmd.AddCommand(diagnosticsSleepCmd)
	rootCmd.AddCommand(diagnosticsCmd)
}

// com.apple.springboardservices
var diagnosticsCmd = &cobra.Command{
	Use:   "diagnostics",
	Short: "Manage diagnostics",
}

var diagnosticsRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart device",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := diagnostics_relay.NewClient(getUDID())
		if err != nil {
			return err
		}
		defer client.Close()
		return client.Restart()
	},
}

var diagnosticsSleepCmd = &cobra.Command{
	Use:   "sleep",
	Short: "Disconnect USB and put device to sleep",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := diagnostics_relay.NewClient(getUDID())
		if err != nil {
			return err
		}
		defer client.Close()
		return client.Sleep()
	},
}

var diagnosticsShutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Shutdown device",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := diagnostics_relay.NewClient(getUDID())
		if err != nil {
			return err
		}
		defer client.Close()
		return client.Shutdown()
	},
}

var diagnosticsQueryCmd = &cobra.Command{
	Use:   "query [!]KEY ...",
	Short: "Query MobileGestalt clear and encrypted keys",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := diagnostics_relay.NewClient(getUDID())
		if err != nil {
			return err
		}
		defer client.Close()
		return client.MobileGestalt(args...)
	},
}
