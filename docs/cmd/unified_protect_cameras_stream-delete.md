## unified protect cameras stream-delete

Delete RTSPS stream(s), based on qualities specified, for a camera

```
unified protect cameras stream-delete [camera ID] [flags]
```

### Options

```
  -h, --help      help for stream-delete
      --high      high stream quality
      --low       low stream quality
      --medium    medium stream quality
      --package   package stream quality
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

