package simulatelocation

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/steeve/itool/client"
	"github.com/steeve/itool/lockdownd"
)

const (
	serviceName = "com.apple.dt.simulatelocation"
)

const (
	commandSetLocation   uint32 = 0
	commandResetLocation uint32 = 1
)

type Client struct {
	udid string
}

func encodeArgs(args ...interface{}) []byte {
	buf := &bytes.Buffer{}
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			binary.Write(buf, binary.BigEndian, uint32(len(v)))
			buf.WriteString(v)
		default:
			binary.Write(buf, binary.BigEndian, v)
		}
	}
	return buf.Bytes()
}

func NewClient(udid string) (*Client, error) {
	return &Client{
		udid: udid,
	}, nil
}

func (c *Client) newClient() (*client.Client, error) {
	return lockdownd.NewClientForService(c.udid, serviceName, false)
}

func (c *Client) SetLocation(latitude, longitude float64) error {
	data := encodeArgs(
		uint32(commandSetLocation),
		strconv.FormatFloat(latitude, 'f', -1, 64),
		strconv.FormatFloat(longitude, 'f', -1, 64),
	)
	lc, err := c.newClient()
	if err != nil {
		return err
	}
	defer lc.Close()
	_, err = lc.Conn().Write(data)
	return err
}

func (c *Client) ResetLocation() error {
	lc, err := c.newClient()
	if err != nil {
		return err
	}
	defer lc.Close()
	return binary.Write(lc.Conn(), binary.BigEndian, uint32(commandResetLocation))
}

func (c *Client) Close() error {
	return nil // no-op
}
