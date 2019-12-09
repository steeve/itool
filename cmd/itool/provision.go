package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/misagent"
)

var provisionRemoveFlags = struct {
	uuid string
}{}

var provisionCopyFlags = struct {
	uuid    string
	outfile string
}{}

func init() {
	rootCmd.AddCommand(provisionCmd)
	provisionCmd.AddCommand(provisionListCmd)
	provisionCmd.AddCommand(provisionRemoveCmd)
	provisionCmd.AddCommand(provisionCopyCmd)
}

var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Manage provisioning profiles",
}

var provisionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed provisioning profiles",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := misagent.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		profiles, err := client.ProvisioningProfiles()
		if err != nil {
			log.Fatal(err)
		}
		if globalFlags.json {
			json.NewEncoder(os.Stdout).Encode(profiles)
			return
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 32, 2, ' ', 0)
		fmt.Fprintln(writer, "UUID\tNAME")
		for _, profile := range profiles {
			fmt.Fprintf(writer, "%s\t%s\n", profile.UUID, profile.Name)
		}
		writer.Flush()
	},
}

var provisionRemoveCmd = &cobra.Command{
	Use:   "remove UUID ...",
	Args:  cobra.MinimumNArgs(1),
	Short: "Remove a provisioning profile",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := misagent.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		for _, uuid := range args {
			if err := client.Remove(uuid); err != nil {
				log.Fatal(err)
			}
		}
	},
}

var provisionCopyCmd = &cobra.Command{
	Use:   "copy",
	Args:  cobra.ExactArgs(2),
	Short: "Copy a provisioning profile from a device to a file",
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		dst := args[1]
		client, err := misagent.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		profilesData, err := client.CopyAll()
		if err != nil {
			log.Fatal(err)
		}
		for _, profileData := range profilesData {
			profile, err := misagent.NewMobileProvisionFromData(profileData)
			if err != nil {
				log.Fatal(err)
			}
			if profile.UUID == src {
				if err := ioutil.WriteFile(dst, profileData, 0666); err != nil {
					log.Fatal(err)
				}
				break
			}
		}
	},
}
