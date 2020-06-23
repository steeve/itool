package mobileconfig

import (
	"fmt"

	"howett.net/plist"

	"github.com/steeve/itool/client"
	"github.com/steeve/itool/lockdownd"
)

const (
	serviceName = "com.apple.mobile.MCInstall"
)

type Client struct {
	c *client.Client
}

func NewClient(udid string) (*Client, error) {
	c, err := lockdownd.NewClientForService(udid, serviceName, false)
	if err != nil {
		return nil, err
	}
	return &Client{
		c: c,
	}, nil
}

func (c *Client) validateError(resp ResponseBase) error {
	if resp.Status != "Acknowledged" {
		var err error
		for _, e := range resp.ErrorChain {
			if err != nil {
				err = fmt.Errorf("%s: %w", e.LocalizedDescription, err)
			} else {
				err = fmt.Errorf("%s", e.LocalizedDescription)
			}
		}
		return err
	}
	return nil
}

func (c *Client) InstallProfile(profileData []byte) error {
	req := &InstallProfileRequest{
		RequestBase: RequestBase{"InstallProfile"},
		Payload:     profileData,
	}
	resp := &InstallProfileResponse{}
	if err := c.c.Request(req, resp); err != nil {
		return err
	}
	if err := c.validateError(resp.ResponseBase); err != nil {
		return err
	}
	return nil
}

func (c *Client) ListProfiles() ([]*Profile, error) {
	req := &GetProfileListRequest{
		RequestBase: RequestBase{"GetProfileList"},
	}
	resp := &GetProfileListResponse{}
	if err := c.c.Request(req, resp); err != nil {
		return nil, err
	}
	if err := c.validateError(resp.ResponseBase); err != nil {
		return nil, err
	}
	profiles := make([]*Profile, 0)
	for _, identifier := range resp.OrderedIdentifiers {
		manifest := resp.ProfileManifest[identifier]
		metadata := resp.ProfileMetadata[identifier]
		profiles = append(profiles, &Profile{
			Identifier: identifier,
			Manifest:   manifest,
			Metadata:   metadata,
		})
	}
	return profiles, nil
}

func (c *Client) RemoveProfile(identifier string) error {
	profiles, err := c.ListProfiles()
	if err != nil {
		return fmt.Errorf("unable to get profiles: %w", err)
	}
	var profile *Profile
	for _, p := range profiles {
		if p.Identifier == identifier {
			profile = p
			break
		}
	}
	if profile == nil {
		return fmt.Errorf("profile %s wasn't found on device", identifier)
	}
	pid := &RemoveProfilePayload{
		PayloadType:       "Configuration",
		PayloadIdentifier: profile.Identifier,
		PayloadUUID:       profile.Metadata.PayloadUUID,
		PayloadVersion:    profile.Metadata.PayloadVersion,
	}
	payloadData, err := plist.Marshal(pid, plist.XMLFormat)
	if err != nil {
		return fmt.Errorf("unable to marshal payload paylist: %w", err)
	}
	req := &RemoveProfileRequest{
		RequestBase:       RequestBase{"RemoveProfile"},
		ProfileIdentifier: payloadData,
	}
	resp := &RemoveProfileResponse{}
	if err := c.c.Request(req, resp); err != nil {
		return err
	}
	if err := c.validateError(resp.ResponseBase); err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	return c.c.Close()
}
