## unified network devices list

List all adopted UniFi Network devices by a specific site

### Synopsis

Calls the devices UniFi Network API endpoint for a specific site ID
and prints the results to stdout.

```
unified network devices list [site ID] [flags]
```

### Options

```
  -h, --help                 help for list
      --hide-page            Hides the returned current page information
      --id-only              List only the ID of listed entities, one per line.
      --page-limit uint32    Limit of items per page
      --page-offset uint32   Offset of page to request
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

