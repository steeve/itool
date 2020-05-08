package fetchsymbols

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/steeve/itool/lockdownd"
)

const (
	serviceName = "com.apple.dt.fetchsymbols"
)

var (
	ListFilesPlistRequest uint32 = 0x30303030
	GetFileRequest        uint32 = 0x00000001
)

type ListFilesResponse struct {
	Files   []string `plist:"files"`
	Version int      `plist:"version"`
}

type Client struct {
	udid string
}

func NewClient(udid string) *Client {
	return &Client{
		udid: udid,
	}
}

func (c *Client) List() ([]string, error) {
	fc, err := lockdownd.NewClientForService(c.udid, serviceName, false)
	if err != nil {
		return nil, err
	}
	if err := c.sendCommand(fc.Conn(), ListFilesPlistRequest); err != nil {
		return nil, err
	}
	resp := &ListFilesResponse{}
	if err := fc.Recv(resp); err != nil {
		return nil, err
	}
	return resp.Files, nil
}

func (c *Client) GetFile(idx uint32) (io.Reader, error) {
	fc, err := lockdownd.NewClientForService(c.udid, serviceName, false)
	if err != nil {
		return nil, err
	}
	if err := c.sendCommand(fc.Conn(), GetFileRequest); err != nil {
		return nil, err
	}
	size := uint64(0)
	if err := c.sendRecv(fc.Conn(), idx, &size); err != nil {
		return nil, err
	}
	return io.LimitReader(fc.Conn(), int64(size)), nil
}

func (c *Client) sendRecv(rw io.ReadWriter, req, resp interface{}) error {
	if err := binary.Write(rw, binary.BigEndian, req); err != nil {
		return err
	}
	return binary.Read(rw, binary.BigEndian, resp)
}

func (c *Client) sendCommand(rw io.ReadWriter, cmd uint32) error {
	respCmd := uint32(0)
	c.sendRecv(rw, cmd, &respCmd)
	if respCmd != cmd {
		return fmt.Errorf("invalid response: wanted %v, got %v", cmd, respCmd)
	}
	return nil
}
