package usbmuxd

type Result int

const (
	ResultOK Result = iota
	ResultBadCommand
	ResultBadDevice
	ResultConnectionRefused
	ResultConnectionUnknown1
	ResultConnectionUnknown2
	ResultBadVersion
)

// type MessageType string

// const (
// 	MessageTypeResult      MessageType = "Result"
// 	MessageTypeConnect     MessageType = "Connect"
// 	MessageTypeListen      MessageType = "Listen"
// 	MessageTypeAttached    MessageType = "Attached"
// 	MessageTypeDetached    MessageType = "Detached"
// 	MessageTypeListDevices MessageType = "ListDevices"
// )

// Total size is always 16 bytes
type Header struct {
	Length      uint32
	Version     uint32
	MessageType uint32
	Tag         uint32
}

const HeaderSize = 16

func NewHeader(length int) Header {
	return Header{
		Length:      uint32(length) + HeaderSize,
		Version:     1,
		MessageType: 8, // plist
		Tag:         1,
	}
}
