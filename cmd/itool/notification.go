package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/notification_proxy"
)

func init() {
	notificationCmd.AddCommand(notificationObserveCmd)
	notificationCmd.AddCommand(notificationPostCmd)
	rootCmd.AddCommand(notificationCmd)
}

var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Manage device notifications",
}

var notificationObserveCmd = &cobra.Command{
	Use:   "observe ID",
	Args:  cobra.ExactArgs(1),
	Short: "Observe notification",
	Run: func(cmd *cobra.Command, args []string) {
		notification := args[0]
		nc, err := notification_proxy.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer nc.Close()
		nc.ObserveNotification(notification)
	},
}

var notificationPostCmd = &cobra.Command{
	Use:   "post ID",
	Args:  cobra.ExactArgs(1),
	Short: "Post notification",
	Run: func(cmd *cobra.Command, args []string) {
		notification := args[0]
		nc, err := notification_proxy.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer nc.Close()
		_ = notification
	},
}
