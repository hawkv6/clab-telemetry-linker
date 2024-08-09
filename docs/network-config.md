# Network Config
Here you can get more information about which network configs are used.

## YANG models
The following YANG models are used to influence the impairments.

## Delay / Jitter
```
sensor-group delay
  sensor-path Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail/delay-measurement-session/last-advertisement-information/advertised-values
```

## Bandwidth
```
 sensor-group bandwidth
  sensor-path Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface/interface-status-and-data/enabled/bandwidth
```

## Packet Loss
```
 sensor-group packet-loss
  sensor-path Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface/interface-status-and-data/enabled/packet-loss-percentage
```

### Utilization
```
 sensor-group utilization
  sensor-path openconfig-interfaces:interfaces/interface/state/counters/in-octets
  sensor-path openconfig-interfaces:interfaces/interface/state/counters/out-octets
```


## Config files
Detailed network configs can be found [here](https://github.com/hawkv6/network)