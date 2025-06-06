## unified protect cameras snapshot

Get a live snapshot image from a specified camera and save it to a file

```
unified protect cameras snapshot [camera ID] [filename] [flags]
```

### Options

```
  -h, --help               help for snapshot
      --jpeg-quality int   JPEG Quality from 1 to 100 (default 100)
      --low-quality        snapshot low quality
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

* [unified protect cameras](unified_protect_cameras.md)	 - Make UniFi Protect `cameras` calls

