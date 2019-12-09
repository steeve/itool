package main

import (
	"log"
	"net"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/debugserver"
)

func init() {
	rootCmd.AddCommand(debugserverCmd)
}

var debugserverCmd = &cobra.Command{
	Use:   "debugserver LOCALADDR",
	Short: "Debugserver proxy",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		listenAddr := args[0]
		listener, err := net.Listen("tcp", listenAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()
		log.Println("listening on", listenAddr)
		for {
			localConn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			dc, err := debugserver.NewClient(getUDID())
			if err != nil {
				localConn.Close()
				log.Println(err)
				continue
			}
			log.Printf("new connection from %s", localConn.RemoteAddr())
			startProxy(dc.Conn(), localConn)
		}
	},
}
