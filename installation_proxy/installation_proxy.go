package installation_proxy

import (
	"sort"

	"github.com/steeve/itool/client"
	"github.com/steeve/itool/lockdownd"
)

const (
	serviceName = "com.apple.mobile.installation_proxy"
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

func (c *Client) Lookup() (map[string]*AppBundle, error) {
	req := &Lookup{
		Command: NewCommand("Lookup"),
	}
	resp := &LookupResult{}
	if err := c.c.Request(req, resp); err != nil {
		return nil, err
	}
	return resp.LookupResult, nil
}

func (c *Client) LookupRaw(keys ...string) (map[string]interface{}, error) {
	req := &Lookup{
		Command: NewCommand("Lookup", keys...),
	}
	resp := &struct {
		LookupResult map[string]interface{}
	}{}
	if err := c.c.Request(req, resp); err != nil {
		return nil, err
	}
	return resp.LookupResult, nil
}

func (c *Client) LookupPath(bundleId string) (string, error) {
	apps, err := c.LookupRaw("CFBundleExecutable", "Path")
	if err != nil {
		return "", err
	}
	if d, ok := apps[bundleId]; ok {
		values := d.(map[string]interface{})
		path := values["Path"].(string) + "/" + values["CFBundleExecutable"].(string)
		return path, nil
	}
	return "", nil
}

func (c *Client) InstalledApps() ([]string, error) {
	apps, err := c.LookupRaw("CFBundleIdentifier")
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(apps))
	for k, _ := range apps {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret, nil
}

func (c *Client) Browse() (interface{}, error) {
	req := &Browse{
		Command: NewCommand("Browse"),
	}
	resp := map[string]interface{}{}
	if err := c.c.Request(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type ProgressFunc func(*ProgressEvent)

func (c *Client) watchProgress(cb ProgressFunc) error {
	for {
		ev := &ProgressEvent{}
		if err := c.c.Recv(ev); err != nil {
			return err
		}
		// Some iOS versions send a message that is not a status message.
		// Ignore it.
		if ev.Status != "" {
			continue
		}
		if ev.Status == "Complete" {
			ev.PercentComplete = 100
		}
		if cb != nil {
			cb(ev)
		}
		if ev.Status == "Complete" {
			return nil
		}
	}
}

func (c *Client) installOrUpgrade(cmd, packagePath string, progressCb ProgressFunc) error {
	req := &InstallOrUpgradeRequest{
		Command:     NewCommand(cmd),
		PackagePath: packagePath,
	}
	if err := c.c.Send(req); err != nil {
		return err
	}
	return c.watchProgress(progressCb)
}

func (c *Client) Install(packagePath string, progressCb ProgressFunc) error {
	return c.installOrUpgrade("Install", packagePath, progressCb)
}

func (c *Client) Upgrade(packagePath string, progressCb ProgressFunc) error {
	return c.installOrUpgrade("Upgrade", packagePath, progressCb)
}

func (c *Client) commandForBundle(cmd, bundleId string, progressCb ProgressFunc) error {
	req := &ApplicationIdentifierRequest{
		Command:               NewCommand(cmd),
		ApplicationIdentifier: bundleId,
	}
	if err := c.c.Send(req); err != nil {
		return err
	}
	return c.watchProgress(progressCb)
}

func (c *Client) Uninstall(bundleId string, progressCb ProgressFunc) error {
	return c.commandForBundle("Uninstall", bundleId, progressCb)
}

func (c *Client) Archive(bundleId string, progressCb ProgressFunc) error {
	return c.commandForBundle("Archive", bundleId, progressCb)
}

func (c *Client) RestoreArchive(bundleId string, progressCb ProgressFunc) error {
	return c.commandForBundle("RestoreArchive", bundleId, progressCb)
}

func (c *Client) RemoveArchive(bundleId string, progressCb ProgressFunc) error {
	return c.commandForBundle("RemoveArchive", bundleId, progressCb)
}

func (c *Client) LookupArchives() error {
	req := &LookupArchivesRequest{
		Command: NewCommand("LookupArchives"),
	}
	_ = req
	return nil
}

func (c *Client) Close() error {
	return c.c.Close()
}
