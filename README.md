<h1 align="center">containerlab telemetry linker</h1>
<p align="center">
    <br>
    <img alt="GitHub Release" src="https://img.shields.io/github/v/release/hawkv6/clab-telemetry-linker?display_name=release&style=flat-square">
    <img src="https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat-square">
    <img src="https://img.shields.io/github/actions/workflow/status/hawkv6/clab-telemetry-linker/testing.yaml?style=flat-square&label=tests">
    <img src="https://img.shields.io/codecov/c/github/hawkv6/clab-telemetry-linker?style=flat-square">
    <img src="https://img.shields.io/github/actions/workflow/status/hawkv6/clab-telemetry-linker/golangci-lint.yaml?style=flat-square&label=checks">
</p>

<p align="center">
</p>

---

## Overview
The clab-telemetry-linker service was developed to solve the missing performance measurement support issue in virtual XRd routers, leading to streaming telemetry messages with empty or fixed values. This service allows for the simulation of network impairments like delay, jitter, packet loss, and bandwidth limitations within a containerlab network. By utilizing the containerlab interface, the clab-telemetry-linker can adjust these impairments in the virtual network and populate the streaming telemetry with these values, plus or minus a random value. It can be used with [Jalapeno](https://github.com/cisco-open/jalapeno) or another similar tech stack using Telegraf, Kafka, and InfluxDB.

![](docs/images/clab-telemetry-linker-overview.drawio.svg)

### Functionality

The Cisco IOS-XRd devices deployed with containerlab transmit telemetry data (Cisco MDT / YANG PUSH) with empty/static values to Telegraf Ingress. The messages are then converted into JSON format and forwarded to Kafka, where they become available in the receiver topic for the clab-telemetry-linker. The data is then processed with the applied impairment values.
After processing, the data is converted into Influx Line Protocol and sent to Kafka Publisher Topic. From there, each message is taken by Telegraf Egress and added to the InfluxDB.
The following impairments can be linked:
- delay
- jitter (delay variation)
- packet loss
- bandwidth / rate

## Usage
- **Set Impairments** - [`set`](docs/set.md)
- **Show Impairments** - [`show`](docs/show.md)
- **Delete Impairments** - [`delete`](docs/delete.md)
- **Start the Service** - [`start`](docs/start.md)
- **Check Version** - `version`

## Installation 
### Using Package Manager
For Debian-based systems, install the package using `apt`:
```bash
sudo apt install ./clab-telemetry-linker_{version}_amd64.deb
```
### Using Binary
```
git clone https://github.com/hawkv6/clab-telemetry-linker
cd clab-telemetry-linker && make binary
sudo ./bin/clab-telemetry-linker
```

## Getting Started

1. Start the collector pipeline.
   - For more information, visit the [hawkv6 deployment guide](https://github.com/hawkv6/deployment).

2. Install the network.
   - Detailed instructions can be found in the [hawkv6 testnetwork guide](https://github.com/hawkv6/network).

3. Install `clab-telemetry-linker`.

4. Set the initial impairments using the `set` command.

5. Start the service using the `start` command.

## Additional Info
- The default configuration file is located at `$HOME/.clab-telemetry-linker/config.yaml`
- The default containerlab prefix is: `clab-hawkv6` (can be modified in the config file)
- More details about network configurations are available in [network config documentation](docs/network-config.md)
- Example telemetry messages can be found in the `examples` folder
- clab-telemetry-linker forwards impairments to the relevant containerlab command. More information can be found [here](https://containerlab.dev/cmd/tools/netem/set/)
