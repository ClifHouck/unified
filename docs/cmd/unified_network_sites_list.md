## unified network sites list

List all sites managed by the Network application

### Synopsis

List local sites managed by this Network application (paginated).
Setups using Multi-Site option enabled will return all created sites,
while if option is disabled it will return just the default site.

```
unified network sites list [flags]
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

* [unified network sites](unified_network_sites.md)	 - Make UniFi Network `sites` calls

