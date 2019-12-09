package notification_proxy

import (
	"github.com/steeve/itool/client"
	"github.com/steeve/itool/lockdownd"
)

const (
	serviceName = "com.apple.mobile.notification_proxy"
)

type RequestBase struct {
	Command string `plist:"Command"`
}

type ObserveNotificationRequest struct {
	RequestBase
	Name string `plist:"Name"`
}

type ObserveNotificationEvent struct {
	RequestBase
	Name string `plist:"Name"`
}

type Client struct {
	c *client.Client
}

func NewClient(udid string) (*Client, error) {
	c, err := lockdownd.NewClientForService(udid, serviceName, false)
	if err != nil {
		return nil, err
	}
	return &Client{
		c: c,
	}, nil
}

func (c *Client) ObserveNotification(notification string) error {
	req := ObserveNotificationRequest{
		RequestBase: RequestBase{"ObserveNotification"},
		Name:        notification,
	}
	_ = req
	return nil
}

func (c *Client) Close() error {
	return c.c.Close()
}
