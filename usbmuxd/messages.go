package usbmuxd

type ResultValue int

const (
	ResultValueOK ResultValue = iota
	ResultValueBadCommand
	ResultValueBadDevice
	ResultValueConnectionRefused
	ResultValueConnectionUnknown1
	ResultValueConnectionUnknown2
	ResultValueBadVersion
)

type RequestBase struct {
	MessageType string
}

type ListenRequest struct {
	RequestBase
	ClientVersionString string
	ProgName            string
}

type ResultResponse struct {
	Number ResultValue
}

type DeviceAttachment struct {
	ConnectionSpeed int
	ConnectionType  string
	DeviceID        int
	LocationID      int
	ProductID       int
	SerialNumber    string
	UDID            string
	USBSerialNumber string
}

type DeviceAttached struct {
	RequestBase
	DeviceID   int
	Properties *DeviceAttachment
}

type DeviceDetached struct {
	RequestBase
	DeviceID int
}

type ConnectRequest struct {
	RequestBase
	ClientVersionString string `plist:"ClientVersionString,omitempty"`
	ProgName            string `plist:"ProgName,omitempty"`
	DeviceID            int    `plist:"DeviceID"`
	PortNumber          uint16 `plist:"PortNumber"`
}

type ListDevicesRequest struct {
	RequestBase
}

type ListDevicesResponse struct {
	DeviceList []*DeviceAttached
}

type ReadPairRecordRequest struct {
	RequestBase
	PairRecordID   string `plist:"PairRecordID"`
	PairRecordData []byte `plist:"PairRecordData,omitempty"`
	DeviceID       int    `plist:"DeviceID,omitempty"`
}

type ReadPairRecordResponse struct {
	PairRecordData []byte
}

type PairRecord struct {
	DeviceCertificate []byte
	EscrowBag         []byte
	HostCertificate   []byte
	HostID            string
	HostPrivateKey    []byte
	RootCertificate   []byte
	RootPrivateKey    []byte
	SystemBUID        string
	WiFiMACAddress    string
}
