# Set impairments

## Overview
The `set` command allows you to configure or overwrite network impairments on specific interfaces of a ContainerLab node. These impairments include delay, jitter, packet loss, and bandwidth rate limitation.

## Command Syntax
```
sudo clab-telemetry-linker set -n <clab-node> -i <interface-name> --delay <value in ms> --jitter <value in ms>  --loss <value in %> --rate <value in kbit/s>
```
- `--node <clab-node>` or `-n <clab-node>`: Specify the ContainerLab node name.
- `--interface <interface-name>` or`-i <interface-name>`: Designate the interface on the node to set impairments.
- `--delay <value in ms>`or `-d <value in ms>`: Set the delay time in milliseconds.
- `--jitter <value in ms>` or `-j <value in ms>`: Set the jitter value in milliseconds.
- `--loss <value in %>` or `-l <value in %>`: Define the packet loss percentage.
- `--rate <value in kbit/s>` or `-r <value in kbit/s>`: Limit the bandwidth rate in kilobits per second.


## Example
To set a delay of 1ms, jitter of 1ms, packet loss of 5%, and a rate limit of 100000 kbit/s on interface Gi0-0-0-0 of node XR-1:
```
sudo clab-telemetry-linker set -n XR-1 -i Gi0-0-0-0 --delay 1ms --jitter 1ms --loss 5 --rate 100000 
-----------+-------+--------+-------------+-------------+
| Interface | Delay | Jitter | Packet Loss | Rate (kbit) |
+-----------+-------+--------+-------------+-------------+
| Gi0-0-0-0 | 1ms   | 1ms    | 5.00%       |      100000 |
+-----------+-------+--------+-------------+-------------+
```