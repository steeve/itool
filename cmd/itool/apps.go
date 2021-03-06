package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/debugserver"
	"github.com/steeve/itool/installation_proxy"
)

func init() {
	rootCmd.AddCommand(appsRootCmd)

	appsListCmd.Flags().StringVarP(&appListFlags.bundleID, "bundleid", "b", "", "Detailed information about bundle id (JSON)")
	appsListCmd.Flags().BoolVarP(&appListFlags.path, "path", "p", false, "Print full path to binary for bundle id")
	appsRootCmd.AddCommand(appsListCmd)

	appsRootCmd.AddCommand(appsInstallCmd)
	appsRootCmd.AddCommand(appsUninstallCmd)
	appsRootCmd.AddCommand(appsRunCmd)

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
		defer client.Close()
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
			fmt.Fprintf(
				writer,
				"%s\t%s\t%s\t%s\n",
				strings.TrimSpace(app.CFBundleIdentifier),
				strings.TrimSpace(app.CFBundleDisplayName),
				strings.TrimSpace(app.CFBundleShortVersionString),
				strings.TrimSpace(app.ApplicationType),
			)
		}
		writer.Flush()
	},
}

var appsRunCmd = &cobra.Command{
	Use:   "run BUNDLEID",
	Short: "Run app",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bundleID := args[0]
		client, err := installation_proxy.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		path, err := client.LookupPath(bundleID)
		if err != nil {
			log.Fatal(err)
		}

		appArgs := []string{path}
		appEnv := []string{}
		if os.Getenv("IDE_DISABLED_OS_ACTIVITY_DT_MODE") == "" {
			appEnv = append(appEnv, "OS_ACTIVITY_DT_MODE=enable")
		}
		proc, err := debugserver.NewProcess(getUDID(), appArgs, appEnv)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			io.Copy(os.Stdout, proc.Stdout())
		}()
		if err := proc.Start(); err != nil {
			log.Fatal(err)
		}
		defer proc.Kill()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGPIPE, syscall.SIGTERM)
		<-c
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
