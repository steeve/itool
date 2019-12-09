package screenshotr

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/steeve/itool/client"
	"github.com/steeve/itool/lockdownd"
)

const (
	serviceName = "com.apple.mobile.screenshotr"
)

type Client struct {
	c *client.Client
}

type ScreenShotRequest struct {
	MessageType string `plist:"MessageType"`
}

type ScreenShotResponse struct {
	ScreenShotData []byte `plist:"ScreenShotData"`
}

func NewClient(udid string) (*Client, error) {
	c, err := lockdownd.NewClientForService(udid, serviceName, false)
	if err != nil {
		return nil, err
	}
	c.DeviceLinkHandshake()
	return &Client{
		c: c,
	}, nil
}

func (c *Client) Screenshot() ([]byte, error) {
	req := ScreenShotRequest{
		MessageType: "ScreenShotRequest",
	}
	c.c.DeviceLinkSend(req)
	resp, err := c.c.DeviceLinkRecv()
	if err != nil {
		return nil, err
	}
	respMap := resp.(map[string]interface{})
	return respMap["ScreenShotData"].([]byte), nil
}

func (c *Client) ScreenshotImage() (image.Image, error) {
	data, err := c.Screenshot()
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewBuffer(data))
	return img, err
}

func (c *Client) Close() error {
	return c.c.Close()
}
