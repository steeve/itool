package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/usbmuxd"
)

var globalFlags = struct {
	udid       string
	json       bool
	usbmuxdUrl string
}{}

var udidOnce sync.Once

var rootCmd = &cobra.Command{
	Use:   "itool",
	Short: "Easy iOS management",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&globalFlags.udid, "usbmuxd", "m", usbmuxd.UsbmuxdURL, "usbmuxd URL")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.udid, "udid", "u", "", "UDID")
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.json, "json", "", false, "JSON output (not all commands)")
}

func getUDID() string {
	udidOnce.Do(func() {
		if globalFlags.udid != "" {
			return
		}
		conn, err := usbmuxd.Dial(cmd.Context(), globalFlags.usbmuxdUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		devices, err := conn.ListDevices()
		if err != nil {
			log.Fatal(err)
		}
		if len(devices) < 1 {
			log.Fatal(fmt.Errorf("no devices are connected"))
		}

		globalFlags.udid = devices[0].UDID
	})
	return globalFlags.udid
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
