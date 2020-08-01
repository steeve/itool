// +build !windows

package usbmuxd

import "net"

var (
	UsbmuxdURL = "unix:///var/run/usbmuxd"
)

func usbmuxdDial() (net.Conn, error) {
	return net.Dial("unix", "/var/run/usbmuxd")
}
