## unified network vouchers generate

Generate one or more hotspot vouchers for a site

```
unified network vouchers generate [site ID] [flags]
```

### Options

```
      --count int         Number of vouchers
      --data-limit int    Data limit in megabytes
      --guest-limit int   Authorized guest limit
  -h, --help              help for generate
      --name string       Name of vouchers
      --rx-limit int      Recieve rate limit in kilobytes
      --time-limit int    Time limit in minutes
      --tx-limit int      Transmit rate limit in kilobytes
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

* [unified network vouchers](unified_network_vouchers.md)	 - Make UniFi Network `vouchers` calls

