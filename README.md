# itool

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat)](https://opensource.org/licenses/Apache-2.0)
[![Twitter Follow](https://img.shields.io/twitter/follow/steeve?label=Follow&style=social)](https://twitter.com/steeve)

`itool` is an easy iOS and composable device management command line interface.
It was made to simplify and automate common development and provisioning tasks.

**This is very much a side project and a work in progress.**

## Goals

- Single binary distribution on macOS, Linux and Windows
- Simple, reusable, composable and embeddable
- Feature parity with Xcode

## Install

Head over to the [Releases](https://github.com/steeve/itool/releases) page and
download the `itool` binary for you OS/Arch.

On macOS and Windows, having iTunes installed is sufficient, since `itool` talks
directly to `usbmuxd`. The device should already be paired, as `itool` will
fetch the keys from `usbmuxd`.

On Linux, you'll need to install `usbmuxd` first, but since the pairing command
isn't done yet, there is no way to interact with the device fully.

## Usage
```
$ itool
Easy iOS management

Usage:
  itool [command]

Available Commands:
  afc          Manage Apple File Conduit (AFC)
  apps         Manage apps
  debugserver  Debugserver proxy
  devices      Manage devices
  diagnostics  Manage diagnostics
  fetchsymbols Manage fetchsymbols
  help         Help about any command
  info         Queries device information keys
  location     Simulate location
  mount        Manage mounts
  notification Manage device notifications
  provision    Manage provisioning profiles
  proxy        Proxy TCP connection to device
  screenshot   Saves a screenshot as a PNG file
  syslog       Relays Syslog to stdout

Flags:
  -h, --help          help for itool
      --json          JSON output (not all commands)
  -u, --udid string   UDID

Use "itool [command] --help" for more information about a command.
$
```

## Examples

#### Run an app bundle
Run an app bundle directly from the command line and get the standard output
there provided the following requirements are met:
- bundle has debug entitlements
- device has developer image mounted
  - Although it can list and umount, itool doesn't implement mounting developer
    images yet. Xcode does it automatically until then, and the image stays
    mounted until reboot.

When `itool` exists (either `Ctrl-C` or killed), the app will be killed on the
device also.
```
$ itool apps run my.app.bundle
Hello world!
^C
$
```

#### Install/uninstall apps
```
$ itool apps install myapp.ipa
```

#### Simulate locations from a `.gpx` file
```
$ itool location play route.gpx
```

#### Take a screenshot
```
$ itool screenshot -o shot.png
```

#### Proxy a TCP connection to a local port on the device
```
$ itool proxy :5000 :3215
```

#### Compose device information with `jq`
```
$ itool info --json | jq '.ProductName + " " + .ProductVersion'
```

#### Manage files
```
$ itool afc ls /
```

## Plans / Work In Progress

Some of those commands are sort of working, some are pure plans.

#### `itool devices pair`

Nescessary to make Linux work properly.

#### `itool apps attach`

Attach to an already running process. While this is not that hard, there is
no way yet to get the standard output of an already running process with
`gdbserver`.

#### `itool pcap`

Capture network packets and dump a `.pcap` file for later analysis.

#### `itool ps`

List running processes on a device.

## Building

```
$ go build github.com/steeve/itool/cmd/itool
```

## Packages

One important aspect of `itool` is that all commands are implemented using the
public APIs of Go packages. This means that it should be simple to use those
packages directly, for instance in a gRPC server.
