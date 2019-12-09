package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/fetchsymbols"
)

func init() {
	fetchsymbolsCmd.AddCommand(fetchsymbolsListCmd)
	fetchsymbolsCmd.AddCommand(fetchsymbolsCopyCmd)
	rootCmd.AddCommand(fetchsymbolsCmd)
}

var fetchsymbolsCmd = &cobra.Command{
	Use:   "fetchsymbols",
	Short: "Manage fetchsymbols",
}

var fetchsymbolsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List symbols files",
	Run: func(cmd *cobra.Command, args []string) {
		fc := fetchsymbols.NewClient(getUDID())
		files, err := fc.List()
		if err != nil {
			log.Fatal(err)
		}
		if globalFlags.json {
			json.NewEncoder(os.Stdout).Encode(&struct {
				Files []string
			}{files})
			return
		}
		for _, f := range files {
			fmt.Println(f)
		}
	},
}

var fetchsymbolsCopyCmd = &cobra.Command{
	Use:   "copy SRC DST",
	Short: "Copy symbols file to host",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		dst := args[1]
		fc := fetchsymbols.NewClient(getUDID())
		files, err := fc.List()
		if err != nil {
			log.Fatal(err)
		}
		idx := -1
		for i, f := range files {
			if f == src {
				idx = i
			}
		}
		if idx < 0 {
			log.Fatalf("unable to find source file %v", src)
		}
		srcReader, err := fc.GetFile(uint32(idx))
		if err != nil {
			log.Fatal(err)
		}
		dstFile, err := os.Create(dst)
		if err != nil {
			log.Fatal(err)
		}
		defer dstFile.Close()
		io.Copy(dstFile, srcReader)
	},
}
