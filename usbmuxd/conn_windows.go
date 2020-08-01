// +build windows

package usbmuxd

import "net"

var (
	UsbmuxdURL = "tcp://localhost:27015"
)

func usbmuxdDial() (net.Conn, error) {
	return net.Dial("tcp", "localhost:27015")
}
