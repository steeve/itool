package main

import (
	"fmt"
	"io"
	"log"
	"os"
	pathpkg "path"
	"path/filepath"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/afc"
)

var afcFlags = struct {
}{}

func init() {
	afcCmd.AddCommand(afcLsCmd)
	afcCmd.AddCommand(afcLnCmd)
	afcCmd.AddCommand(afcMvCmd)
	afcCmd.AddCommand(afcRmCmd)
	afcCmd.AddCommand(afcMkdirCmd)
	afcCmd.AddCommand(afcSendCmd)
	afcCmd.AddCommand(afcRecvCmd)
	afcCmd.AddCommand(afcCatCmd)

	rootCmd.AddCommand(afcCmd)
}

var afcCmd = &cobra.Command{
	Use:   "afc",
	Short: "Manage Apple File Conduit (AFC)",
}

var afcLsCmd = &cobra.Command{
	Use:   "ls",
	Args:  cobra.MinimumNArgs(1),
	Short: "list directory contents",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := afc.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		writer := tabwriter.NewWriter(os.Stdout, 0, 128, 1, ' ', tabwriter.AlignRight)
		defer writer.Flush()

		sort.Strings(args)
		for i, root := range args {
			fmt.Fprintf(writer, "%s:\n", root)
			err := client.Walk(root, func(path string, info os.FileInfo, err error) error {
				// Don't print the root
				if path == root && info.IsDir() {
					return nil
				}

				name := pathpkg.Base(path)
				if info.IsDir() {
					name += "/"
				}
				fmt.Fprintf(writer, "%d\t%s\t %s\n", info.Size(), info.ModTime().Format("Jan _2 2006 15:04"), name)

				// Don't recurse in directories
				if info.IsDir() {
					return filepath.SkipDir
				}

				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
			if i < len(args)-1 {
				fmt.Fprintln(writer)
			}
			writer.Flush()
		}
	},
}

var afcSendCmd = &cobra.Command{
	Use:   "send [SOURCE] ... [TARGET]",
	Args:  cobra.MinimumNArgs(2),
	Short: "send files to device",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := afc.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		srcs := args[:len(args)-1]
		target := args[len(args)-1] // last argument is always the destination

		for _, src := range srcs {
			err := client.CopyToDevice(target, src, func(dst, src string, info os.FileInfo) {
				fmt.Println(src, "->", dst)
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var afcRecvCmd = &cobra.Command{
	Use:   "fetch [FROM] [TO]",
	Args:  cobra.ExactArgs(2),
	Short: "fetch files from device",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := afc.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		src := args[0]
		dst := args[1]
		if err := client.CopyFromDevice(dst, src, func(dst, src string, info os.FileInfo) {
			fmt.Println(src, "->", dst)
		}); err != nil {
			log.Fatal(err)
		}
	},
}

var afcLnCmd = &cobra.Command{
	Use:   "ln [FROM] [TO]",
	Args:  cobra.MinimumNArgs(2),
	Short: "make links",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := afc.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
	},
}

var afcMvCmd = &cobra.Command{
	Use:   "mv [FROM] [TO]",
	Args:  cobra.ExactArgs(2),
	Short: "move files",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := afc.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		if err := client.RenamePath(args[0], args[1]); err != nil {
			log.Fatal(err)
		}
	},
}

var afcRmCmd = &cobra.Command{
	Use:   "rm [FILE] ...",
	Args:  cobra.MinimumNArgs(1),
	Short: "remove directory entries",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := afc.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		for _, arg := range args {
			if err := client.RemoveAll(arg); err != nil {
				log.Fatal(fmt.Errorf("can't remove %v: %v", arg, err))
			}
		}
	},
}

var afcMkdirCmd = &cobra.Command{
	Use:   "mkdir",
	Args:  cobra.MinimumNArgs(1),
	Short: "make directories",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := afc.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		for _, arg := range args {
			if err := client.MakeDir(arg); err != nil {
				log.Fatal(fmt.Errorf("can't create %v: %v", arg, err))
			}
		}
	},
}

var afcCatCmd = &cobra.Command{
	Use:   "cat",
	Args:  cobra.MinimumNArgs(1),
	Short: "print files",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := afc.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		for _, arg := range args {
			f, err := client.FileRefOpen(arg, os.O_RDONLY)
			if err != nil {
				log.Fatal(err)
			}
			io.Copy(os.Stdout, f)
			f.Close()
		}
	},
}
