package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/screenshotr"
)

var screenshotFlags = struct {
	outfile string
}{}

func init() {
	rootCmd.AddCommand(screenshotCmd)
	screenshotCmd.Flags().StringVarP(&screenshotFlags.outfile, "out", "o", "", "Output file (PNG)")
}

var screenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "Saves a screenshot as a PNG file",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := screenshotr.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		data, err := client.Screenshot()
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(screenshotFlags.outfile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		f.Write(data)
	},
}
