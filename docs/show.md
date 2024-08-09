# Show

## Overview
The `show` command displays the current network impairments applied to a ContainerLab node. This command provides a detailed view of the network settings, including each interface's delay, jitter, packet loss, and bandwidth rate.

## Command Syntax
```
sudo clab-telemetry-linker show -n <clab-node>
```
- `--node <clab-node>` or  `n <clab-node>`: Specifies the ContainerLab node name for which you want to view the impairments.

## Example
To display the network impairments for the node XR-1:
```
sudo clab-telemetry-linker show -n XR-1
+-----------+-------+--------+-------------+-------------+
| Interface | Delay | Jitter | Packet Loss | Rate (kbit) |
+-----------+-------+--------+-------------+-------------+
| lo        | N/A   | N/A    | N/A         | N/A         |
| eth0      | N/A   | N/A    | N/A         | N/A         |
| Gi0-0-0-0 | 0s    | 0s     | 0.00%       |           0 |
| Gi0-0-0-1 | 4ms   | 0s     | 0.00%       |           0 |
| Gi0-0-0-2 | 6ms   | 0s     | 5.00%       |           0 |
+-----------+-------+--------+-------------+-------------+
```