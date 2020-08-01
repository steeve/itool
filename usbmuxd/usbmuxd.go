package usbmuxd

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"syscall"

	"github.com/steeve/itool/netutils"
	"howett.net/plist"
)

var (
	ErrNotFound = errors.New("device not found")
)

type Conn struct {
	net.Conn
}

func htonl(v uint16) uint16 {
	return (v << 8 & 0xFF00) | (v >> 8 & 0xFF)
}

func Dial(ctx context.Context, usbmuxdURL string) (*Conn, error) {
	c, err := netutils.DialURLContext(ctx, usbmuxdURL)
	if err != nil {
		return nil, fmt.Errorf("unable to dial usbmuxd: %w", err)
	}
	return &Conn{
		Conn: c,
	}, nil
}

func DialDefault(ctx context.Context) (*Conn, error) {
	return Dial(ctx, UsbmuxdURL)
}

func DialUDID(ctx context.Context, usbmuxdURL, udidAddr string) (net.Conn, error) {
	conn, err := Dial(ctx, usbmuxdURL)
	if err != nil {
		return nil, err
	}
	if err := conn.Dial(udidAddr); err != nil {
		return nil, err
	}
	return conn, nil
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

func (c Conn) Dial(address string) error {
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
		if device.SerialNumber == udid {
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
	if err := c.Send(req); err != nil {
		return err
	}
	return c.Recv(resp)
}
