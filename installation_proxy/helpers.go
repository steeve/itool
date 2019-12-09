package installation_proxy

import (
	"fmt"
	"os"
	"path"

	"github.com/steeve/itool/afc"
)

func (c *Client) CopyAndInstall(pkgPath string, progressCb ProgressFunc) error {
	afcClient, err := afc.NewClient(c.c.UDID())
	if err != nil {
		return err
	}
	defer afcClient.Close()
	if err := afcClient.CopyToDevice("/", pkgPath, func(dst, src string, info os.FileInfo) {
		fmt.Println(src, "->", dst)
	}); err != nil {
		return err
	}
	return c.Install(path.Join("/", path.Base(pkgPath)), progressCb)
}
