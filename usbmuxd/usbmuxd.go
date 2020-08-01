package usbmuxd

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"sync"
	"syscall"

	"github.com/steeve/itool/netutils"
	"howett.net/plist"
)

var (
	ErrNotFound = errors.New("device not found")
)

func init() {
	netutils.RegisterURLScheme("usbmux", func(ctx context.Context, u *url.URL) (net.Conn, error) {
		return Dial(ctx, u.Host)
	})
}

type Conn struct {
	net.Conn
	sync.RWMutex
}

func htonl(v uint16) uint16 {
	return (v << 8 & 0xFF00) | (v >> 8 & 0xFF)
}

func OpenWithUrl(ctx context.Context, usbmuxdURL string) (*Conn, error) {
	c, err := netutils.URLDialContext(ctx, usbmuxdURL)
	if err != nil {
		return nil, fmt.Errorf("unable to dial usbmuxd: %w", err)
	}
	return &Conn{
		Conn: c,
	}, nil
}

func Open(ctx context.Context) (*Conn, error) {
	return OpenWithUrl(ctx, UsbmuxdURL)
}

func Connect(ctx context.Context, udid string, port uint16) (net.Conn, error) {
	conn, err := Open(ctx)
	if err != nil {
		return nil, err
	}
	if err := conn.Connect(udid, port); err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

func Dial(ctx context.Context, udidAddr string) (net.Conn, error) {
	udid, portStr, err := net.SplitHostPort(udidAddr)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	return Connect(ctx, udid, uint16(port))
}

func (c Conn) ReadPairRecord(udid string) (*PairRecord, error) {
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

func (c Conn) SavePairRecord(udid string, pairRecord *PairRecord) error {
	recordData, err := plist.Marshal(pairRecord, plist.XMLFormat)
	if err != nil {
		return fmt.Errorf("unable to marshal PairRecord: %w", err)
	}
	req := &SavePairRecordRequest{
		RequestBase:    RequestBase{"SavePairRecord"},
		PairRecordID:   udid,
		PairRecordData: recordData,
	}
	resp := &SavePairRecordResponse{}
	if err := c.Request(req, resp); err != nil {
		return fmt.Errorf("unable to save PairRecord: %w", err)
	}
	return nil
}

func (c Conn) DeletePairRecord(udid string) error {
	req := &DeletePairRecordRequest{
		RequestBase:  RequestBase{"DeletePairRecord"},
		PairRecordID: udid,
	}
	resp := &DeletePairRecordResponse{}
	if err := c.Request(req, resp); err != nil {
		return fmt.Errorf("unable to delete PairRecord: %w", err)
	}
	return nil
}

func (c Conn) ListDevices() ([]*DeviceAttachment, error) {
	req := &ListDevicesRequest{
		RequestBase: RequestBase{"ListDevices"},
	}
	resp := &ListDevicesResponse{}
	if err := c.Request(req, resp); err != nil {
		return nil, err
	}
	devices := make([]*DeviceAttachment, 0, len(resp.DeviceList))
	for _, dev := range resp.DeviceList {
		devices = append(devices, dev.Properties)
	}
	return devices, nil
}

func (c Conn) DeviceIDFromUDID(udid string) (int, error) {
	devices, err := c.ListDevices()
	if err != nil {
		return 0, err
	}
	for _, device := range devices {
		if udid == "" {
			return device.DeviceID, nil
		}
		if device.SerialNumber == udid {
			return device.DeviceID, nil
		}
	}
	return 0, fmt.Errorf("unable to find device with udid: %v", udid)
}

func (c Conn) Connect(udid string, port uint16) error {
	deviceId, err := c.DeviceIDFromUDID(udid)
	if err != nil {
		return fmt.Errorf("unable to connect: %w", err)
	}
	req := &ConnectRequest{
		RequestBase: RequestBase{"Connect"},
		DeviceID:    deviceId,
		PortNumber:  htonl(port),
	}
	resp := &ResultResponse{}
	if err := c.Request(req, resp); err != nil {
		return err
	}
	if resp.Number == ResultValueConnectionRefused {
		return syscall.ECONNREFUSED
	}
	return nil
}

func (c Conn) ReadBUID() (string, error) {
	req := &ReadBUIDRequest{
		RequestBase: RequestBase{"ReadBUID"},
	}
	resp := &ReadBUIDResponse{}
	if err := c.Request(req, resp); err != nil {
		return "", err
	}
	return resp.BUID, nil
}

func (c Conn) Send(msg interface{}) error {
	c.Lock()
	defer c.Unlock()
	return c.send(msg)
}

func (c Conn) send(msg interface{}) error {
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

func (c Conn) Recv(msg interface{}) error {
	c.Lock()
	defer c.Unlock()
	return c.recv(msg)
}

func (c Conn) recv(msg interface{}) error {
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

func (c Conn) Request(req, resp interface{}) error {
	c.Lock()
	defer c.Unlock()
	if err := c.send(req); err != nil {
		return err
	}
	return c.recv(resp)
}
