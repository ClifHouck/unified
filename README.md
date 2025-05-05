# Unified

[![CI](https://github.com/ClifHouck/unified/actions/workflows/ci.yaml/badge.svg)](https://github.com/ClifHouck/unified/actions/workflows/ci.yaml)

## Features

Command line utility `unified`: makes requests to UniFi APIs and outputs their responses.

Fully-featured Golang client for type-safe, programmatic access to UniFi APIs.

## Quickstart

Install `unified` from `go`:
```bash
$ go install ClifHouck/unified@latest
```

Set your UniFi API key:
```bash
export UNIFI_API_KEY=$(cat unifi_api.key)
```

Try it out:
```bash
$ unified network info
```

If all goes well, you should see something like:
```json
$ unified network info
{
  "applicationVersion": "9.1.120"
}
```

Note that the output is valid JSON, just like the UniFi applications produce.


## `unified` Command Line Usage

`unified`

`unified` has a full help system accessible through the `--help` flag.

To access Network APIs you will use:

```bash
$ unified network
```
gtkj
and for Protect APIs
```bash
$ unified protect
```

## Golang Client Usage

## Protect Websocket Event Streams

## Project Roadmap

Full Protect API support is planned in short order. After the initial versioned
release of this project, work will proceed towards full Protect V1 API support.
As well as a full v1.0.0 release.

Reference documentation through godoc is also a priority before a v1.0.0 release.

### API Support Status

UniFi Network API V1 is fully supported as of Network application version "9.1.120".

UniFi Protect API is only partially supported, with the following endpoints supported:

- `/v1/meta/info`
- `/v1/subscribe/devices`: only partial type support.
- `/v1/subscribe/events`
- `/v1/cameras/`
- `/v1/cameras/{id}`

## Contributing

## Thanks

Copyright 2025 - Clifton Houck
