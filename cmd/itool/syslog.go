package main

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/syslog_relay"
)

func init() {
	rootCmd.AddCommand(syslogCmd)
}

var syslogCmd = &cobra.Command{
	Use:   "syslog",
	Short: "Relays Syslog to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		rc, err := syslog_relay.Syslog(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer rc.Close()
		io.Copy(os.Stdout, rc)
	},
}
