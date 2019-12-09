package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/image_mounter"
)

func init() {
	rootCmd.AddCommand(mountCmd)
}

var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "Manage mounts",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			mountListCmd.Run(cmd, args)
			return
		}
	},
}

var mountListCmd = &cobra.Command{
	Use:   "list",
	Short: "List mounts",
	Run: func(cmd *cobra.Command, args []string) {
		imc, err := image_mounter.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer imc.Close()
		sigs, err := imc.LookupImage(image_mounter.ImageTypeDeveloper)
		if err != nil {
			log.Fatal(err)
		}
		if globalFlags.json {
			json.NewEncoder(os.Stdout).Encode(sigs)
			return
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 1024, 1, ' ', 0)
		fmt.Fprintln(writer, "INDEX\tSIGNATURE")
		for i, sig := range sigs.ImageSignature {
			fmt.Fprintf(writer, "%d\t%s\n", i, hex.EncodeToString(sig))
		}
		writer.Flush()
	},
}

var mountMountCmd = &cobra.Command{
	Use:   "mount",
	Args:  cobra.ExactArgs(2),
	Short: "Mount image",
	Run: func(cmd *cobra.Command, args []string) {
		imc, err := image_mounter.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer imc.Close()
	},
}
