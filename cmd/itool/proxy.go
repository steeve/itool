package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steeve/itool/usbmuxd"
)

func init() {
	rootCmd.AddCommand(proxyCmd)
}

var proxyCmd = &cobra.Command{
	Use:   "proxy LOCALADDR REMOTEPORT",
	Short: "Proxy TCP connection to device",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		localAddr := args[0]
		remotePort := args[1]
		listener, err := net.Listen("tcp", localAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()
		log.Println("Listening on", localAddr)
		for {
			localConn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			remoteConn, err := usbmuxd.Dial(getUDID() + ":" + remotePort)
			if err != nil {
				localConn.Close()
				log.Println(err)
				continue
			}
			log.Printf("new connection from %s to %s:%s", localConn.RemoteAddr(), getUDID(), remotePort)
			startProxy(localConn, remoteConn)
		}
	},
}

func copyy(dst io.Writer, src io.ReadWriter, prefix string) error {
	for {
		data := make([]byte, 1024)
		n, err := src.Read(data)
		if err != nil {
			return err
		}

		if strings.Contains(string(data[:n]), "QEnableCompression") {
			fmt.Println("DISABLING COMPRESSION")
			// panic(string(data[:n]))
			src.Write([]byte("$OK#00"))
			continue
		}

		fmt.Println(prefix, string(data[:n]))
		dst.Write(data[:n])
	}
}

func startProxy(conn1, conn2 io.ReadWriteCloser) {
	go func() {
		defer conn1.Close()
		defer conn2.Close()
		io.Copy(conn2, conn1)
	}()
	go func() {
		defer conn1.Close()
		defer conn2.Close()
		io.Copy(conn1, conn2)
	}()
}
