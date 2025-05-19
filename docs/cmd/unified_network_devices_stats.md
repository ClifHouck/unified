## unified network devices stats

Get latest (live) statistics of a specific adopted device.

### Synopsis

Get latest (live) statistics of a specific adopted device.
Response contains latest readings from a single device, such as CPU and
memory utilization, uptime, uplink tx/rx rates etc

```
unified network devices stats [site ID] [device ID] [flags]
```

### Options

```
  -h, --help   help for stats
```

### Options inherited from parent commands

```
      --config string                  config file (default is $HOME/.unified.yaml)
      --debug                          Enable debug logging
      --host string                    Hostname of UniFi API (default "unifi")
      --insecure                       Skip verification of UniFi TLS certificate. (default true)
      --keep-alive-interval duration   Interval between keep-alive pings sent for websocket streams (default 30s)
      --trace                          Enable trace logging
```

### SEE ALSO

* [unified network devices](unified_network_devices.md)	 - Make UniFi Network `devices` calls

