package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/lockdownd"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info [KEY]",
	Short: "Queries device information keys",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := ""
		if len(args) > 0 {
			key = args[0]
		}
		client, err := lockdownd.NewClient(getUDID())
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		v, err := client.GetValue(key)
		if err != nil {
			log.Fatal(err)
		}
		if globalFlags.json {
			json.NewEncoder(os.Stdout).Encode(v)
			return
		}
		switch v := v.(type) {
		case map[string]interface{}:
			keys := make([]string, 0)
			for k := range v {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Print(k, ": ", v[k])
				fmt.Println("")
			}
		}
		if key != "" {
			fmt.Println(v)
		}
	},
}
