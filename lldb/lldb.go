package lldb

import (
	"fmt"

	"github.com/steeve/itool/installation_proxy"
	"github.com/steeve/itool/lockdownd"
)

func SDKPathForUDID(udid string, deviceType string) (string, error) {
	client, err := lockdownd.NewClient(udid)
	if err != nil {
		return "", err
	}
	defer client.Close()
	values, err := client.GetValues()
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("~/Library/Developer/Xcode/%s DeviceSupport/%s (%s) %s/Symbols", deviceType, values.ProductVersion, values.BuildVersion, values.CPUArchitecture)
	return path, nil
}

func LLDBTargetCreate(udid, addr, bundleId, appPath string) (string, error) {
	sdkPath, err := SDKPathForUDID(udid, "iOS")
	if err != nil {
		return "", err
	}
	ipc, err := installation_proxy.NewClient(udid)
	if err != nil {
		return "", err
	}
	containerPath, err := ipc.LookupPath(bundleId)
	if err != nil {
		return "", err
	}
	ret := fmt.Sprintf("platform select remote-ios --sysroot \"%s\"\n", sdkPath)
	ret += fmt.Sprintf("target create \"%s\"\n", appPath)
	ret += fmt.Sprintf("script lldb.target.modules[0].SetPlatformFileSpec(lldb.SBFileSpec(\"%s\"))\n", containerPath)
	ret += fmt.Sprintf("process connect connect://%s\n", addr)
	return ret, nil
}
