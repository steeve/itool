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
	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.AddCommand(deviceKeyCmd)
	rootCmd.AddCommand(deviceCmd)
}

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Manage devices",
}

var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List connected devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := usbmuxd.NewConn()
		if err != nil {
			return err
		}
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
			lc, err := lockdownd.NewClient(device.UDID)
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

var deviceKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Dump TLS key for a device pairing",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := usbmuxd.NewConn()
		if err != nil {
			return err
		}
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
