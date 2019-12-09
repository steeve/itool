package usbmuxd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"syscall"
	"time"

	"howett.net/plist"
)

var (
	ErrNotFound = errors.New("device not found")
)

type Conn struct {
	conn net.Conn
}

func htonl(v uint16) uint16 {
	return (v << 8 & 0xFF00) | (v >> 8 & 0xFF)
}

func NewConn() (*Conn, error) {
	conn, err := usbmuxdDial()
	if err != nil {
		return nil, err
	}
	return &Conn{
		conn: conn,
	}, nil
}

func Dial(address string) (net.Conn, error) {
	conn, err := NewConn()
	if err != nil {
		return nil, err
	}
	if err := conn.Dial(address); err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Conn) ReadPairRecord(udid string) (*PairRecord, error) {
	req := &ReadPairRecordRequest{
		RequestBase:  RequestBase{"ReadPairRecord"},
		PairRecordID: udid,
	}
	resp := &ReadPairRecordResponse{}
	if err := c.Request(req, resp); err != nil {
		return nil, err
	}
	record := &PairRecord{}
	if _, err := plist.Unmarshal(resp.PairRecordData, record); err != nil {
		return nil, err
	}
	return record, nil
}

func (c *Conn) ListDevices() ([]*DeviceAttachment, error) {
	req := &ListDevicesRequest{
		RequestBase: RequestBase{"ListDevices"},
	}
	resp := &ListDevicesResponse{}
	if err := c.Request(req, resp); err != nil {
		return nil, err
	}
	devices := make([]*DeviceAttachment, 0, len(resp.DeviceList))
	for _, dev := range resp.DeviceList {
		if dev.Properties.UDID != "" {
			devices = append(devices, dev.Properties)
		}
	}
	return devices, nil
}

func (c *Conn) Dial(address string) error {
	udid, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	devices, err := c.ListDevices()
	if err != nil {
		return err
	}
	deviceID := -1
	for _, device := range devices {
		if device.UDID == udid {
			deviceID = device.DeviceID
		}
	}
	if deviceID < 0 {
		return fmt.Errorf("unable to find device with udid: %v", udid)
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	req := &ConnectRequest{
		RequestBase: RequestBase{"Connect"},
		DeviceID:    deviceID,
		PortNumber:  htonl(uint16(portNumber)),
	}
	// fmt.Println(req, portNumber)
	resp := &ResultResponse{}
	if err := c.Request(req, resp); err != nil {
		return err
	}
	if resp.Number == ResultValueConnectionRefused {
		return syscall.ECONNREFUSED
	}
	return nil
}

func (c *Conn) Send(msg interface{}) error {
	data, err := plist.Marshal(msg, plist.XMLFormat)
	if err != nil {
		return err
	}
	hdr := NewHeader(len(data))
	if err := binary.Write(c, binary.LittleEndian, hdr); err != nil {
		return err
	}
	if _, err := c.Write(data); err != nil {
		return err
	}
	return nil
}

func (c *Conn) Recv(msg interface{}) error {
	hdr := Header{}
	if err := binary.Read(c, binary.LittleEndian, &hdr); err != nil {
		return err
	}
	data := make([]byte, hdr.Length-HeaderSize)
	if _, err := io.ReadFull(c, data); err != nil {
		return err
	}
	if _, err := plist.Unmarshal(data, msg); err != nil {
		return err
	}
	return nil
}

func (c *Conn) Request(req, resp interface{}) error {
	if err := c.Send(req); err != nil {
		return err
	}
	return c.Recv(resp)
}

// Read reads data from the connection.
// Read can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *Conn) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

// Write writes data to the connection.
// Write can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (c *Conn) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to Read or
// Write. After a deadline has been exceeded, the connection
// can be refreshed by setting a deadline in the future.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
//
// Note that if a TCP connection has keep-alive turned on,
// which is the default unless overridden by Dialer.KeepAlive
// or ListenConfig.KeepAlive, then a keep-alive failure may
// also return a timeout error. On Unix systems a keep-alive
// failure on I/O can be detected using
// errors.Is(err, syscall.ETIMEDOUT).
func (c *Conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
