// +build !windows

package usbmuxd

import "net"

func usbmuxdDial() (net.Conn, error) {
	return net.Dial("unix", "/var/run/usbmuxd")
}
