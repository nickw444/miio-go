# miio-go

[![Coverage Status](https://coveralls.io/repos/github/nickw444/miio-go/badge.svg?branch=master)](https://coveralls.io/github/nickw444/miio-go?branch=master)

An implementation of the miIO home protocol by Xiaomi written in Golang. Heavily inspired by:

 - [The Javascript Implementation](https://github.com/aholstenson/miio)
 - [Protocol Specification](https://github.com/OpenMiHome/mihome-binary-protocol)
 - [API design in this Lifx client implementation (pdx/lifx)](https://github.com/pdf/golifx)

This implementation has been design with the following concerns:
 - Testability
 - Development without a miIO device handy (or performing any real network operations)
 - A simple event-based API.

## Supported Devices
At the moment, only the following devices are officially supported by this library. Feel free to
[submit a pull request](), I'd be more than happy to have more devices supported by this library.

 - Xiaomi Mi Smart WiFi Socket (v1 - no USB)

[API Documentation can be found here]()

## Simulator

A device simulator/emulator exists in the [simulator](simulator/) package. It takes
advantage of the low level network used to communicate with real devices to emulate
hardware devices.

[Give it a try!](simulator/)

## Tokens
Documentation coming soon...

## Examples
Documentation coming soon...

## CLI
Documentation coming soon...
