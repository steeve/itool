package debugserver

import (
	"net"

	"github.com/steeve/itool/client"
	"github.com/steeve/itool/lockdownd"
)

const (
	serviceName = "com.apple.debugserver"
)

type Client struct {
	c *client.Client
}

func NewClient(udid string) (*Client, error) {
	c, err := lockdownd.NewClientForService(udid, serviceName, false)
	if err != nil {
		return nil, err
	}
	// Disable TLS after the handshake
	// See https://github.com/libimobiledevice/libimobiledevice/issues/793
	c.DisableSSL()

	return &Client{
		c: c,
	}, nil
}

func (c *Client) Conn() net.Conn {
	return c.c.Conn()
}

func (c *Client) Close() error {
	return c.c.Close()
}
