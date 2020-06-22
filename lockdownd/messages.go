package lockdownd

type RequestBase struct {
	Request string `plist:"Request"`
}

type ResponseBase struct {
	Request string
	Result  string
}

type QueryTypeRequest struct {
	RequestBase
}

type QueryTypeResponse struct {
	ResponseBase
	Type string
}

type StartSessionRequest struct {
	RequestBase
	HostID     string
	SystemBUID string
}

type StartSessionResponse struct {
	ResponseBase
	EnableSessionSSL bool
	SessionID        string
}

type DeviceValues struct {
	BasebandCertId             int
	BasebandKeyHashInformation struct {
		AKeyStatus int
		SKeyStatus int
	}
	BasebandSerialNumber    []byte
	BasebandVersion         string
	BoardId                 int
	BuildVersion            string
	CPUArchitecture         string
	ChipID                  int
	DeviceClass             string
	DeviceColor             string
	DeviceName              string
	DevicePublicKey         []byte
	DeviceCertificate       []byte
	DieID                   int
	HardwareModel           string
	HasSiDP                 bool
	PartitionType           string
	ProductName             string
	ProductType             string
	ProductVersion          string
	ProductionSOC           bool
	ProtocolVersion         string
	SupportedDeviceFamilies []int
	TelephonyCapability     bool
	UniqueChipID            int64
	UniqueDeviceID          string
	WiFiAddress             string
}

type GetValueRequest struct {
	RequestBase
	Domain string `plist:"Domain,omitempty"`
	Key    string `plist:"Key,omitempty"`
}

type GetValueResponse struct {
	ResponseBase
	Value *DeviceValues
}

type StartServiceRequest struct {
	RequestBase
	Service   string
	EscrowBag []byte `plist:",omitempty"`
}

type StartServiceResponse struct {
	ResponseBase
	Service          string
	Port             int
	EnableServiceSSL bool
}

type EnterRecoveryRequest struct {
	RequestBase
}

type EnterRecoveryResponse struct {
	ResponseBase
}
