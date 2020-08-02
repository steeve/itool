package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		fc := fetchsymbols.NewClient(getUDID())
		files, err := fc.List()
		if err != nil {
			return fmt.Errorf("unable to list fetchsymbols: %w", err)
		}
		if globalFlags.json {
			json.NewEncoder(os.Stdout).Encode(&struct {
				Files []string
			}{files})
			return nil
		}
		for _, f := range files {
			fmt.Println(f)
		}
		return nil
	},
}

var fetchsymbolsCopyCmd = &cobra.Command{
	Use:   "copy SRC DST",
	Short: "Copy symbols file to host",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		src := args[0]
		dst := args[1]
		fc := fetchsymbols.NewClient(getUDID())
		files, err := fc.List()
		if err != nil {
			return fmt.Errorf("unable to list fetchsymbols: %w", err)
		}
		idx := -1
		for i, f := range files {
			if f == src {
				idx = i
			}
		}
		if idx < 0 {
			return fmt.Errorf("unable to find source file %v", src)
		}
		srcReader, err := fc.GetFile(uint32(idx))
		if err != nil {
			return fmt.Errorf("unable to get file %d: %w", idx, err)
		}
		dstFile, err := os.Create(dst)
		if err != nil {
			fmt.Errorf("unable to create file %v: %w", dst, err)
		}
		defer dstFile.Close()
		io.Copy(dstFile, srcReader)
		return nil
	},
}
