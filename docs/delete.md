# Delete

## Overview
The delete command is used to remove network impairments from a specific interface on a ContainerLab node. This action resets the network settings for the specified interface to their default state, which typically means no delay, jitter, packet loss, or rate limitation.

## Command Syntax
```
clab-telemetry-linker delete -n <clab-node> -i <interface-name>
```
- `--node <clab-node>` or `-n <clab-node>`: Specifies the ContainerLab node name.
- `--interface <interface-name>` or `-i <interface-name>`: Designates the specific interface on the node from which to delete impairments.

## Examples
To delete impairments from interface Gi0-0-0-0 on node XR-1:
```
clab-telemetry-linker delete -n XR-1 -i Gi0-0-0-0
+-----------+-------+--------+-------------+-------------+
| Interface | Delay | Jitter | Packet Loss | Rate (kbit) |
+-----------+-------+--------+-------------+-------------+
| Gi0-0-0-0 | 0ms   | 0s     | 0.00%       |           0 |
+-----------+-------+--------+-------------+-------------+
```