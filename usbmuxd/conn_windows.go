// +build windows

package usbmuxd

import "net"

func usbmuxdDial() (net.Conn, error) {
	return net.Dial("tcp", "localhost:27015")
}
