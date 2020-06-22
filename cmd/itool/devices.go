package main

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/lockdownd"
	"github.com/steeve/itool/usbmuxd"
)

func init() {
	devicesCmd.AddCommand(devicesListCmd)
	devicesCmd.AddCommand(devicesKeyCmd)
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
		conn, err := usbmuxd.NewConn()
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
		conn, err := usbmuxd.NewConn()
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
