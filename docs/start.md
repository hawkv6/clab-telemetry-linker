# Start the service
## Overview
The `start` command is used to initiate the clab-telemetry-linker tool. This tool listens for telemetry data, enriches it with container lab impairments, and forwards it to a Kafka broker. Upon starting, the tool initiates several threads responsible for different threads:
1. Consumer - consumes the messages from the receiver topic
2. Processor - processes the messages according to the config (applied impairments)
3. Publisher - publishes the messages to the publisher-topic

## Command Syntax
```
sudo clab-telemetry-linker start -b <kafka-host>:<port> -r <receiver-topic> -p <publisher-topic> 
```
- `--broker <kafka-host:port>` or `-b <kafka-host>:<port>`: Specifies the Kafka broker host and port.
- `--receiver-topic <receiver-topic>` or `-r <receiver-topic>`: Designates the Kafka topic to receive unprocessed telemetry data.
- `--publisher-topic <publisher-topic>` or `-p <publisher-topic>`: Indicates the Kafka topic for publishing processed telemetry data.

## Example
To start the service with Kafka broker at 172.16.19.77:9094, receiving data from hawkv6.telemetry.unprocessed, and publishing to hawkv6.telemetry.processed:
```
sudo clab-telemetry-linker start -b 172.16.19.77:9094 -r hawkv6.telemetry.unprocessed -p hawkv6.telemetry.processed
INFO[2024-01-21T11:31:19Z] Read config file:  /home/ins/.clab-telemetry-linker/config.yaml  subsystem=config
INFO[2024-01-21T11:31:19Z] Start all services                            subsystem=service
INFO[2024-01-21T11:31:19Z] Start consuming messages from broker 172.16.19.77:9094 and topic hawkv6.telemetry.unprocessed  subsystem=consumer
INFO[2024-01-21T11:31:19Z] Starting processing messages                  subsystem=processor
INFO[2024-01-21T11:31:19Z] Starting publishing messages to broker 172.16.19.77:9094 and topic hawkv6.telemetry.processed  subsystem=publisher

^CINFO[2024-01-21T11:31:24Z] Received interrupt signal, shutting down      subsystem=cmd
INFO[2024-01-21T11:31:24Z] Stopping all services                         subsystem=service
INFO[2024-01-21T11:31:24Z] Stop consumer with values:  172.16.19.77:9094 hawkv6.telemetry.unprocessed  subsystem=consumer
INFO[2024-01-21T11:31:24Z] Stopping processor                            subsystem=processor
INFO[2024-01-21T11:31:24Z] Stopping publisher with broker  172.16.19.77:9094  and topic  hawkv6.telemetry.processed  subsystem=publisher
```

## Additional Info
Network impairments can be adjusted even after the service has started. The service automatically detects configuration changes and adapts accordingly.