package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/installation_proxy"
)

func init() {
	rootCmd.AddCommand(appsRootCmd)

	appsListCmd.Flags().StringVarP(&appListFlags.bundleID, "bundleid", "b", "", "Detailed information about bundle id (JSON)")
	appsListCmd.Flags().BoolVarP(&appListFlags.path, "path", "p", false, "Print full path to binary for bundle id")
	appsRootCmd.AddCommand(appsListCmd)

	appsRootCmd.AddCommand(appsInstallCmd)
	appsRootCmd.AddCommand(appsUninstallCmd)

	appsArchiveCmd.AddCommand(appsArchiveCreateCmd)
	appsArchiveCmd.AddCommand(appsRestoreArchiveCmd)
	appsArchiveCmd.AddCommand(appsRemoveArchiveCmd)
	appsRootCmd.AddCommand(appsArchiveCmd)
}

var appsRootCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage apps",
}

var appListFlags = &struct {
	bundleID string
	path     bool
}{}

var appsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List apps",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := installation_proxy.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		if appListFlags.path {
			path, err := client.LookupPath(appListFlags.bundleID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(path)
			return
		}

		if globalFlags.json {
			data, err := client.LookupRaw()
			if err != nil {
				log.Fatal(err)
			}
			if appListFlags.bundleID != "" {
				json.NewEncoder(os.Stdout).Encode(data[appListFlags.bundleID])
			} else {
				json.NewEncoder(os.Stdout).Encode(data)
			}
			return
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 128, 2, ' ', 0)
		fmt.Fprintln(writer, "BUNDLE\tNAME\tVERSION\tTYPE")
		apps, err := client.Lookup()
		if err != nil {
			log.Fatal(err)
		}
		appids := make([]string, 0, len(apps))
		for appid := range apps {
			appids = append(appids, appid)
		}
		sort.Strings(appids)
		for _, appid := range appids {
			app := apps[appid]
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", app.CFBundleIdentifier, app.CFBundleDisplayName, app.CFBundleVersion, app.ApplicationType)
		}
		writer.Flush()
	},
}

var appsInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install .ipa or .app",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := installation_proxy.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		for _, apppkg := range args {
			log.Println("Installing", apppkg)
			if err := client.CopyAndInstall(apppkg, func(ev *installation_proxy.ProgressEvent) {
				log.Printf("%s (%d%%)\n", ev.Status, ev.PercentComplete)
			}); err != nil {
				log.Fatal(err)
			}
		}
	},
}

var appsUninstallCmd = &cobra.Command{
	Use:   "uninstall [BUNDLEID] ...",
	Short: "unininstall apps",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := installation_proxy.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		for _, bundleId := range args {
			if err := client.Uninstall(bundleId, func(ev *installation_proxy.ProgressEvent) {
				log.Printf("%s (%d%%)\n", ev.Status, ev.PercentComplete)
			}); err != nil {
				log.Fatal(err)
			}
		}
	},
}

var appsArchiveCmd = &cobra.Command{
	Use:   "archive [BUNDLEID] ...",
	Short: "app archives management",
}

var appsArchiveCreateCmd = &cobra.Command{
	Use:   "create [BUNDLEID] ...",
	Short: "archive apps",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := installation_proxy.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		for _, bundleId := range args {
			if err := client.Archive(bundleId, func(ev *installation_proxy.ProgressEvent) {
				log.Printf("%s (%d%%)\n", ev.Status, ev.PercentComplete)
			}); err != nil {
				log.Fatal(err)
			}
		}
	},
}

var appsArchiveLookupCmd = &cobra.Command{
	Use:   "list",
	Short: "list app archives",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := installation_proxy.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		client.LookupArchives()
	},
}

var appsRestoreArchiveCmd = &cobra.Command{
	Use:   "restore [BUNDLEID] ...",
	Short: "restore apps archive",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := installation_proxy.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		for _, bundleId := range args {
			if err := client.RestoreArchive(bundleId, func(ev *installation_proxy.ProgressEvent) {
				log.Printf("%s (%d%%)\n", ev.Status, ev.PercentComplete)
			}); err != nil {
				log.Fatal(err)
			}
		}
	},
}

var appsRemoveArchiveCmd = &cobra.Command{
	Use:   "remove [BUNDLEID] ...",
	Short: "remove apps archive",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := installation_proxy.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		for _, bundleId := range args {
			if err := client.RemoveArchive(bundleId, func(ev *installation_proxy.ProgressEvent) {
				log.Printf("%s (%d%%)\n", ev.Status, ev.PercentComplete)
			}); err != nil {
				log.Fatal(err)
			}
		}
	},
}
