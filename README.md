# Unified

[![CI](https://github.com/ClifHouck/unified/actions/workflows/ci.yaml/badge.svg)](https://github.com/ClifHouck/unified/actions/workflows/ci.yaml)

## Features

- `unified`: command makes requests to UniFi APIs and prints their responses
- `import "github.com/ClifHouck/unified/client": Fully-featured Golang client 
for programmatic access to UniFi APIs.

## Quickstart

Install `unified`:
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

## Golang Client Usage

## Protect Websocket Event Streams

## Thanks

Copyright 2025 - Clifton Houck
