package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/mobileconfig"
)

func init() {
	mobileconfigCmd.AddCommand(mobileconfigInstallCmd)
	mobileconfigCmd.AddCommand(mobileconfigListCmd)
	mobileconfigCmd.AddCommand(mobileconfigRemoveCmd)
	rootCmd.AddCommand(mobileconfigCmd)
}

var mobileconfigCmd = &cobra.Command{
	Use:   "mobileconfig",
	Short: "Manage mobileconfigs",
}

var mobileconfigInstallCmd = &cobra.Command{
	Use:   "install MOBILECONFIG",
	Args:  cobra.ExactArgs(1),
	Short: "Install mobileconfig",
	RunE: func(cmd *cobra.Command, args []string) error {
		mobileconfigFile := args[0]
		profileData, err := ioutil.ReadFile(mobileconfigFile)
		if err != nil {
			return fmt.Errorf("unable to open %s: %w", mobileconfigFile, err)
		}
		mc, err := mobileconfig.NewClient(getUDID())
		if err != nil {
			return fmt.Errorf("unable to open connection to mobileconfig service: %w", err)
		}
		defer mc.Close()
		if err := mc.InstallProfile(profileData); err != nil {
			return fmt.Errorf("unable to install mobileconfig: %w", err)
		}
		return nil
	},
}

var mobileconfigListCmd = &cobra.Command{
	Use:   "list",
	Short: "List mobileconfigs",
	RunE: func(cmd *cobra.Command, args []string) error {
		mc, err := mobileconfig.NewClient(getUDID())
		if err != nil {
			return fmt.Errorf("unable to open connection to mobileconfig service: %w", err)
		}
		defer mc.Close()
		profiles, err := mc.ListProfiles()
		if err != nil {
			return fmt.Errorf("unable to list mobileconfigs: %w", err)
		}
		if globalFlags.json {
			return json.NewEncoder(os.Stdout).Encode(profiles)
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 32, 2, ' ', 0)
		defer writer.Flush()
		fmt.Fprintln(writer, "ID\tNAME\tVERSION\tACTIVE")
		for _, profile := range profiles {
			fmt.Fprintf(writer, "%s\t%s\t%d\t%t\n",
				profile.Identifier,
				profile.Metadata.PayloadDisplayName,
				profile.Metadata.PayloadVersion,
				profile.Manifest.IsActive,
			)
		}
		return nil
	},
}

var mobileconfigRemoveCmd = &cobra.Command{
	Use:   "remove ID",
	Short: "remove mobileconfig",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileId := args[0]
		mc, err := mobileconfig.NewClient(getUDID())
		if err != nil {
			return fmt.Errorf("unable to open connection to mobileconfig service: %w", err)
		}
		defer mc.Close()
		if err := mc.RemoveProfile(profileId); err != nil {
			return fmt.Errorf("unable to remove mobileconfig %s: %w", profileId, err)
		}
		return nil
	},
}
