package client

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"io"
	"net"

	"github.com/steeve/itool/usbmuxd"

	"howett.net/plist"
)

type Client struct {
	tlsConn    *tls.Conn
	conn       net.Conn
	udid       string
	pairRecord *usbmuxd.PairRecord
}

func NewClient(udid string, port int) (*Client, error) {
	usbmuxConn, err := usbmuxd.Open(context.TODO())
	if err != nil {
		return nil, err
	}
	pairRecord, err := usbmuxConn.ReadPairRecord(udid)
	if err != nil {
		return nil, err
	}
	if err := usbmuxConn.Connect(udid, uint16(port)); err != nil {
		return nil, err
	}
	c := &Client{
		conn:       usbmuxConn,
		pairRecord: pairRecord,
		udid:       udid,
	}
	return c, nil
}

func Dial(ctx context.Context, conn *usbmuxd.Conn, udid string, port int) (*Client, error) {
	usbmuxConn, err := usbmuxd.Open(ctx)
	if err != nil {
		return nil, err
	}
	pairRecord, err := usbmuxConn.ReadPairRecord(udid)
	if err != nil {
		return nil, err
	}
	if err := usbmuxConn.Connect(udid, uint16(port)); err != nil {
		return nil, err
	}
	c := &Client{
		conn:       usbmuxConn,
		pairRecord: pairRecord,
		udid:       udid,
	}
	return c, nil
}

func NewClient2(ctx context.Context, conn net.Conn) (*Client, error) {
	c := &Client{
		conn: conn,
	}
	return c, nil
}

func (c *Client) EnableSSL() error {
	crt, err := tls.X509KeyPair(c.pairRecord.HostCertificate, c.pairRecord.HostPrivateKey)
	if err != nil {
		return err
	}
	config := &tls.Config{
		Certificates:       []tls.Certificate{crt},
		InsecureSkipVerify: true,
	}
	c.tlsConn = tls.Client(c.conn, config)
	if err := c.tlsConn.Handshake(); err != nil {
		return err
	}
	return nil
}

func (c *Client) EnableSSL2(pairRecord *usbmuxd.PairRecord) error {
	crt, err := tls.X509KeyPair(pairRecord.HostCertificate, pairRecord.HostPrivateKey)
	if err != nil {
		return err
	}
	config := &tls.Config{
		Certificates:       []tls.Certificate{crt},
		InsecureSkipVerify: true,
	}
	c.tlsConn = tls.Client(c.conn, config)
	if err := c.tlsConn.Handshake(); err != nil {
		return err
	}
	return nil
}

func (c *Client) DisableSSL() {
	c.tlsConn = nil
}

func (c *Client) PairRecord() *usbmuxd.PairRecord {
	return c.pairRecord
}

func (c *Client) UDID() string {
	return c.udid
}

func (c *Client) Close() error {
	return c.Conn().Close()
}

func (c *Client) Request(req, resp interface{}) error {
	if err := c.Send(req); err != nil {
		return err
	}
	return c.Recv(resp)
}

func (c *Client) Send(req interface{}) error {
	data, err := plist.Marshal(req, plist.XMLFormat)
	if err != nil {
		return err
	}
	// fmt.Println(">>>", string(data))
	if err := binary.Write(c.Conn(), binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}
	if _, err := c.Conn().Write(data); err != nil {
		return err
	}
	return nil
}

func (c *Client) Recv(resp interface{}) error {
	data, err := c.RecvBytes()
	if err != nil {
		return err
	}
	// fmt.Println("<<<", string(data))
	if _, err := plist.Unmarshal(data, resp); err != nil {
		return err
	}
	return nil
}

func (c *Client) RecvBytes() ([]byte, error) {
	respLen := uint32(0)
	if err := binary.Read(c.Conn(), binary.BigEndian, &respLen); err != nil {
		return nil, err
	}
	data := make([]byte, respLen)
	if _, err := io.ReadFull(c.Conn(), data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) Conn() net.Conn {
	if c.tlsConn != nil {
		return c.tlsConn
	}
	return c.conn
}

func (c *Client) DeviceLinkHandshake() error {
	versionExchange := []interface{}{}
	if err := c.Recv(&versionExchange); err != nil {
		return err
	}
	reply := []interface{}{"DLMessageVersionExchange", "DLVersionsOk", versionExchange[1]}
	if err := c.Send(reply); err != nil {
		return err
	}
	ready := []interface{}{}
	return c.Recv(&ready)
}

func (c *Client) DeviceLinkSend(msg interface{}) error {
	return c.Send([]interface{}{"DLMessageProcessMessage", msg})
}

func (c *Client) DeviceLinkRecv() (interface{}, error) {
	dlMsg := []interface{}{}
	if err := c.Recv(&dlMsg); err != nil {
		return nil, err
	}
	return dlMsg[1], nil
}
