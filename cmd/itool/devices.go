package main

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/steeve/itool/diagnostics_relay"
	"github.com/steeve/itool/lockdownd"
	"github.com/steeve/itool/usbmuxd"
)

func init() {
	devicesCmd.AddCommand(devicesListCmd)
	devicesCmd.AddCommand(devicesKeyCmd)
	devicesCmd.AddCommand(devicesQueryCmd)
	devicesCmd.AddCommand(devicesShutdownCmd)
	devicesCmd.AddCommand(devicesRestartCmd)
	devicesCmd.AddCommand(devicesSleepCmd)
	devicesCmd.AddCommand(devicesInfoCmd)
	devicesCmd.AddCommand(devicesRecoveryCmd)
	rootCmd.AddCommand(devicesCmd)
}

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Manage devices",
}

var devicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List connected devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := usbmuxd.Dial(context.Background(), globalFlags.usbmuxdUrl)
		if err != nil {
			return err
		}
		defer conn.Close()
		devices, err := conn.ListDevices()
		if err != nil {
			return err
		}
		if globalFlags.json {
			return json.NewEncoder(os.Stdout).Encode(&struct {
				Devices []*usbmuxd.DeviceAttachment
			}{devices})
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 32, 2, ' ', 0)
		fmt.Fprintln(writer, "UUID\tNAME\tCONNECTION")
		for _, device := range devices {
			lc, err := lockdownd.NewClient(device.SerialNumber)
			if err != nil {
				return err
			}
			defer lc.Close()
			name, err := lc.GetValue("DeviceName")
			if err != nil {
				return err
			}
			fmt.Fprintf(writer, "%s\t%s\t%s\n", device.SerialNumber, name, device.ConnectionType)
		}
		writer.Flush()
		return nil
	},
}

var devicesKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Dump TLS key for a device pairing",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := usbmuxd.Dial(context.Background(), globalFlags.usbmuxdUrl)
		if err != nil {
			return err
		}
		defer conn.Close()
		pairRecord, err := conn.ReadPairRecord(getUDID())
		if err != nil {
			return err
		}
		pem.Encode(
			os.Stdout,
			&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: pairRecord.HostCertificate,
			},
		)
		pem.Encode(
			os.Stdout,
			&pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: pairRecord.HostPrivateKey,
			},
		)
		return nil
	},
}

var devicesRestartCmd = &cobra.Command{
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

var devicesSleepCmd = &cobra.Command{
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

var devicesShutdownCmd = &cobra.Command{
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

var devicesQueryCmd = &cobra.Command{
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

var devicesInfoCmd = &cobra.Command{
	Use:   "info [KEY]",
	Short: "Queries device information keys",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := ""
		if len(args) > 0 {
			key = args[0]
		}
		client, err := lockdownd.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		v, err := client.GetValue(key)
		if err != nil {
			log.Fatal(err)
		}
		if globalFlags.json {
			json.NewEncoder(os.Stdout).Encode(v)
			return
		}
		switch v := v.(type) {
		case map[string]interface{}:
			keys := make([]string, 0)
			for k := range v {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Print(k, ": ", v[k])
				fmt.Println("")
			}
		}
		if key != "" {
			fmt.Println(v)
		}
	},
}

var devicesRecoveryCmd = &cobra.Command{
	Use:   "recovery",
	Short: "Enter recovery",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := lockdownd.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		if err := client.EnterRecovery(); err != nil {
			return fmt.Errorf("unable to enter recovery: %w", err)
		}
		return nil
	},
}
