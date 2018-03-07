# MiiO Simulator

This simulator has been created to allow easier development when away from
hardware devices and to test the entire integration of the stack, rather
than the unit testing currently available in this library.

```
usage: simulator [<flags>] [<device>]

Flags:
  --help                Show context-sensitive help (also try --help-long and --help-man).
  --device-id=12341234  Device ID for the simulated device
  --device-token=00ff00ff00ff00ff00ff00ff00ff00ff
                        The device token to use for encrypted payloads
  --(no-)reveal-token   Whether or not to reveal the device token

Args:
  [<device>]  Device to simulate

```

### Available Devices

- `powerplug`
- `yeelight`

Devices are a work in progress. All devices are a built from a collection of capabilities.
Available capabilities include:

- power
- info

## Building

```
go build
```
