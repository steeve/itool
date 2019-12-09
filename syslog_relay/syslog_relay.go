package syslog_relay

import (
	"io"

	"github.com/steeve/itool/lockdownd"
)

const (
	serviceName = "com.apple.syslog_relay"
)

func Syslog(udid string) (io.ReadCloser, error) {
	c, err := lockdownd.NewClientForService(udid, serviceName, false)
	if err != nil {
		return nil, err
	}
	if err := c.Send("watch"); err != nil {
		return nil, err
	}
	return c.Conn(), nil
}
