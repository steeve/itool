package misagent

import (
	"time"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/steeve/itool/client"
	"github.com/steeve/itool/lockdownd"
	"howett.net/plist"
)

const (
	serviceName = "com.apple.misagent"
)

type RequestBase struct {
	MessageType string `plist:"MessageType"`
}

type CopyRequest struct {
	RequestBase
	ProfileType string `plist:"ProfileType"`
}

type CopyResponse struct {
	Payload [][]byte `plist:"Payload"`
	Status  int      `plist:"Status"`
}

type CopyAllRequest struct {
	RequestBase
	ProfileType string `plist:"ProfileType"`
}

type CopyAllResponse struct {
	Payload [][]byte `plist:"Payload"`
	Status  int      `plist:"Status"`
}

type InstallRequest struct {
	RequestBase
	Profile     []byte `plist:"Profile"`
	ProfileType string `plist:"ProfileType"`
}

type InstallResponse struct {
	Payload [][]byte `plist:"Payload"`
	Status  int      `plist:"Status"`
}

type RemoveRequest struct {
	RequestBase
	ProfileID   string `plist:"ProfileID"`
	ProfileType string `plist:"ProfileType"`
}

type MobileProvision struct {
	AppIDName                   string                 `plist:"AppIDName"`
	ApplicationIdentifierPrefix []string               `plist:"ApplicationIdentifierPrefix"`
	CreationDate                time.Time              `plist:"CreationDate"`
	Platform                    []string               `plist:"Platform"`
	IsXcodeManaged              bool                   `plist:"IsXcodeManaged"`
	DeveloperCertificates       [][]byte               `plist:"DeveloperCertificates"`
	Entitlements                map[string]interface{} `plist:"Entitlements"`
	ExpirationDate              time.Time              `plist:"ExpirationDate"`
	Name                        string                 `plist:"Name"`
	ProvisionsAllDevices        bool                   `plist:"ProvisionsAllDevices"`
	TeamIdentifier              []string               `plist:"TeamIdentifier"`
	TeamName                    string                 `plist:"TeamName"`
	TimeToLive                  int                    `plist:"TimeToLive"`
	UUID                        string                 `plist:"UUID"`
	Version                     int                    `plist:"Version"`
}

type Client struct {
	c *client.Client
}

func NewMobileProvisionFromData(data []byte) (*MobileProvision, error) {
	packet := ber.DecodePacket(data)
	plistData := packet.Children[1].Children[0].Children[2].Children[1].Children[0]
	profile := &MobileProvision{}
	if _, err := plist.Unmarshal(plistData.ByteValue, profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// func ([]*MobileProvision, error) {
// 	req := &CopyAllRequest{
// 		RequestBase: RequestBase{"CopyAll"},
// 		ProfileType: "Provisioning",
// 	}
// 	resp := &CopyAllResponse{}
// 	if err := c.c.Request(req, resp); err != nil {
// 		return nil, err
// 	}
// 	profiles := make([]*MobileProvision, 0, len(resp.Payload))
// 	for _, payload := range resp.Payload {
// 		packet := ber.DecodePacket(payload)
// 		plistData := packet.Children[1].Children[0].Children[2].Children[1].Children[0]
// 		profile := &MobileProvision{}
// 		if _, err := plist.Unmarshal(plistData.ByteValue, profile); err != nil {
// 			return nil, err
// 		}
// 		profiles = append(profiles, profile)
// 	}
// 	return profiles, nil
// }

func NewClient(udid string) (*Client, error) {
	c, err := lockdownd.NewClientForService(udid, serviceName, false)
	if err != nil {
		return nil, err
	}
	return &Client{
		c: c,
	}, nil
}

func (c *Client) Install(profileData []byte) error {
	req := &InstallRequest{
		RequestBase: RequestBase{"Install"},
		Profile:     profileData,
		ProfileType: "Provisioning",
	}
	resp := &InstallResponse{}
	if err := c.c.Request(req, resp); err != nil {
		return err
	}
	return nil
}

func (c *Client) Copy() ([]*MobileProvision, error) {
	req := &CopyRequest{
		RequestBase: RequestBase{"Copy"},
		ProfileType: "Provisioning",
	}
	resp := &CopyResponse{}
	if err := c.c.Request(req, resp); err != nil {
		return nil, err
	}
	profiles := make([]*MobileProvision, 0, len(resp.Payload))
	for _, payload := range resp.Payload {
		packet := ber.DecodePacket(payload)
		// Found this path by dumping the whole payload as JSON.
		// It may be fragile.
		plistData := packet.Children[1].Children[0].Children[2].Children[1].Children[0]
		profile := &MobileProvision{}
		if _, err := plist.Unmarshal(plistData.ByteValue, profile); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (c *Client) CopyAll() ([][]byte, error) {
	req := &CopyAllRequest{
		RequestBase: RequestBase{"CopyAll"},
		ProfileType: "Provisioning",
	}
	resp := &CopyAllResponse{}
	if err := c.c.Request(req, resp); err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (c *Client) ProvisioningProfiles() ([]*MobileProvision, error) {
	profilesData, err := c.CopyAll()
	if err != nil {
		return nil, err
	}
	profiles := make([]*MobileProvision, 0, len(profilesData))
	for _, profileData := range profilesData {
		profile, err := NewMobileProvisionFromData(profileData)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// func (c *Client) CopyAll() ([]*MobileProvision, error) {
// 	req := &CopyAllRequest{
// 		RequestBase: RequestBase{"CopyAll"},
// 		ProfileType: "Provisioning",
// 	}
// 	resp := &CopyAllResponse{}
// 	if err := c.c.Request(req, resp); err != nil {
// 		return nil, err
// 	}
// 	profiles := make([]*MobileProvision, 0, len(resp.Payload))
// 	for _, payload := range resp.Payload {
// 		packet := ber.DecodePacket(payload)
// 		plistData := packet.Children[1].Children[0].Children[2].Children[1].Children[0]
// 		profile := &MobileProvision{}
// 		if _, err := plist.Unmarshal(plistData.ByteValue, profile); err != nil {
// 			return nil, err
// 		}
// 		profiles = append(profiles, profile)
// 	}
// 	return profiles, nil
// }

func (c *Client) Remove(uuid string) error {
	req := &RemoveRequest{
		RequestBase: RequestBase{"Remove"},
		ProfileID:   uuid,
		ProfileType: "Provisioning",
	}
	resp := &CopyAllResponse{}
	if err := c.c.Request(req, resp); err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	return c.c.Close()
}
