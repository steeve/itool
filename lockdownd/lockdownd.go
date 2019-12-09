package lockdownd

import (
	"github.com/steeve/itool/client"
)

const (
	port = 62078
)

type Client struct {
	c *client.Client
}

func NewClientForService(udid, serviceName string, withEscrowBag bool) (*client.Client, error) {
	lc, err := NewClient(udid)
	if err != nil {
		return nil, err
	}
	defer lc.Close()
	svc, err := lc.StartService(serviceName, withEscrowBag)
	if err != nil {
		return nil, err
	}
	c, err := client.NewClient(udid, svc.Port)
	if err != nil {
		return nil, err
	}
	if svc.EnableServiceSSL {
		c.EnableSSL()
	}
	return c, nil
}

func NewClient(udid string) (*Client, error) {
	c, err := client.NewClient(udid, port)
	if err != nil {
		return nil, err
	}
	req := &StartSessionRequest{
		RequestBase: RequestBase{"StartSession"},
		HostID:      c.PairRecord().HostID,
		SystemBUID:  c.PairRecord().SystemBUID,
	}
	resp := &StartSessionResponse{}
	c.Request(req, resp)
	if resp.EnableSessionSSL {
		if err := c.EnableSSL(); err != nil {
			return nil, err
		}
	}
	return &Client{
		c: c,
	}, nil
}

func (c *Client) GetValues() (*DeviceValues, error) {
	req := &GetValueRequest{
		RequestBase: RequestBase{"GetValue"},
	}
	resp := &GetValueResponse{}
	if err := c.c.Request(req, resp); err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (c *Client) GetValue(key string) (interface{}, error) {
	req := &GetValueRequest{
		RequestBase: RequestBase{"GetValue"},
		Key:         key,
	}
	var resp map[string]interface{}
	if err := c.c.Request(req, &resp); err != nil {
		return nil, err
	}
	return resp["Value"], nil
}

func (c *Client) QueryType() (string, error) {
	req := &QueryTypeRequest{
		RequestBase: RequestBase{"QueryType"},
	}
	resp := &QueryTypeResponse{}
	if err := c.c.Request(req, resp); err != nil {
		return "", err
	}
	return resp.Type, nil
}

func (c *Client) StartService(service string, withEscrowBag bool) (*StartServiceResponse, error) {
	req := &StartServiceRequest{
		RequestBase: RequestBase{"StartService"},
		Service:     service,
	}
	if withEscrowBag {
		req.EscrowBag = c.c.PairRecord().EscrowBag
	}
	resp := &StartServiceResponse{}
	if err := c.c.Request(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Close() error {
	return c.c.Close()
}
