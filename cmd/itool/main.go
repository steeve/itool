package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/usbmuxd"
)

var globalFlags = struct {
	udid string
	json bool
}{}

var udidOnce sync.Once

var (
	defaultPairRecord *usbmuxd.PairRecord
)

var rootCmd = &cobra.Command{
	Use:   "itool",
	Short: "Easy iOS management",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&usbmuxd.UsbmuxdURL, "usbmuxd", "m", usbmuxd.UsbmuxdURL, "usbmuxd URL")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.udid, "udid", "u", "", "UDID")
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.json, "json", "", false, "JSON output (not all commands)")
}

func getUDID() string {
	udidOnce.Do(func() {
		if globalFlags.udid != "" {
			return
		}
		conn, err := usbmuxd.Open(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		devices, err := conn.ListDevices()
		if err != nil {
			log.Fatal(err)
		}
		for _, device := range devices {
			globalFlags.udid = device.UDID
			return
		}

		log.Fatal(fmt.Errorf("no devices are connected"))
	})
	return globalFlags.udid
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
