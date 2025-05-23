## unified network clients list

List connected clients of a site

### Synopsis

List connected clients of a site (paginated). Clients are either
physical devices (computers, smartphones, connected by wire or wirelessly),
or active VPN connections.

```
unified network clients list [site ID] [flags]
```

### Options

```
      --filter string        Filter results based on expression
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

* [unified network clients](unified_network_clients.md)	 - Make UniFi Network `clients` calls

