# Unified

[![CI](https://github.com/ClifHouck/unified/actions/workflows/ci.yaml/badge.svg)](https://github.com/ClifHouck/unified/actions/workflows/ci.yaml)

An Unofficial UniFi Network & Protect API Client & CLI command, written in Golang.

# Features

* Command line utility `unified`: makes it easy to send requests to UniFi APIs
    and get their responses.
* Fully-featured Golang client for type-safe, programmatic access to UniFi APIs.

# Quickstart

Install `unified` from `go`:
```bash
$ go install ClifHouck/unified@latest
```

Set your UniFi [API key](#UniFi-API-Key-Instructions)
```bash
export UNIFI_API_KEY=$(cat unifi_api.key)
```

Try it out:
```bash
$ unified network info
```

If all goes well, you should see something like:
```json
{
  "applicationVersion": "9.1.120"
}
```

Note that the output is valid JSON, just like the UniFi applications produce.


# `unified` Command Line Usage

`unified` has a full help system accessible through the `--help` flag.

To access Network APIs you will use:

```bash
$ unified network
```

and for Protect APIs
```bash
$ unified protect
```

# Golang Client Usage

To instantiate a client you should call `client.NewClient`:

```golang
    // ctx is a context.Context.
    // NewDefaultConfig returns a client.Config struct, apiKey is a string populated with your UniFi API key.
    // log is a *logrus.Logger.
    unifiClient := client.NewClient(ctx, client.NewDefaultConfig(apiKey), log)
```

Access to actual API client calls is mediated through [NetworkV1](https://github.com/ClifHouck/unified/blob/main/types/network.go)
and [ProtectV1](https://github.com/ClifHouck/unified/blob/main/types/protect.go) interfaces. Like so:

```golang
    networkInfo, err := unifiClient.Network.Info()
    if err != nil {
        return err
    }
    fmt.Println(networkInfo.ApplicationVersion) // Prints '9.1.20' or similar.
```

These interfaces strive to closely mirror the actual APIs exposed by the
Network and Protect applications.

# Protect Websocket Event Streams

Protect's API has a couple of interesting endpoints which allow a client to subscribe
to a Websocket event stream. The `ProtectV1` interface exposes a pair of
methods which provide easy access to those streams as golang channels:

```golang
    // Websocket updates
    SubscribeDeviceEvents() (<-chan *ProtectDeviceEvent, error)
    SubscribeProtectEvents() (<-chan *ProtectEvent, error)
```

and while it's certainly do-able to consume those event channels, unified also
provides a handler for each event type that makes consuming and re-acting to
these event's even easier. Here's a brief example of using `ProtectEventStreamHandler`:

```golang
    eventChan, err := unifiClient.Protect.SubscribeProtectEvents()
    if err != nil {
        return err
    }

    streamHandler := client.NewProtectEventStreamHandler(ctx, eventChan)

    // Register handler callback function for type-safe access to the Protect
    // RingEvent.
    streamHandler.SetRingEventHandler(func(eventType string, _ *types.RingEvent) {
        if eventType == "add" {
            fmt.Println("Got add ring event!")
        }
        ...
    })

    // Start processing events from the stream channel.
    go streamHandler.Process()

    <-ctx.Done()
```

A nearly-identical struct exists to handle `ProtectDeviceEvent`s: `ProtectDeviceEventStreamHandler`.

[doorbell.go](https://github.com/ClifHouck/unified/blob/main/examples/doorbell/doorbell.go) is a full
example of using a stream handler. Example programs can be built via:
```bash
$ mage buildExamples
```

# UniFi API Key Instructions
Learn how to generate an API key from [UniFi's offcial documentation](https://help.ui.com/hc/en-us/articles/30076656117655-Getting-Started-with-the-Official-UniFi-API).
Network and Protect are "Local Applications".

>[!WARNING]
>Your API key is a sensitive secret! Please keep it securely stored.
>It gives *FULL* API access to your UniFi applications. If you need to revoke
>an API key navigate to your UniFi application and go to your Admin user:
>Settings -> Admins & User -> Select Admin account associated with the API key -> Click on API Key -> Remove.

You can always generate a new API key if necessary.

# Project Roadmap

Full Protect API support is planned in short order. After the initial versioned
release of this project, work will proceed towards full Protect V1 API support.
As well as a full v1.0.0 release.

Reference documentation through godoc is also a priority before a v1.0.0 release.

Access API might be supported in a future release. Contributions welcome here.

## TLS Issue

Unifi's provided TLS certificates are self-signed and do not sign for the `unifi`
hostname. They *DO* sign for `unifi.local`, but the default DNS configuration
for my UDM Pro does not seem to add an entry for `unifi.local`. Therefore TLS
verification generally fails when `https` protocol connections are attemped.
`unified` defaults to TLS verification being off.

In short, at least two general issues need to be solved to enable TLS verification
in general:

1. UniFi certificates should be signed by a trusted authority. Or Ubiquiti needs
   to provide an easy way to import their authority chain.
2. UniFi certificates should at least be signed for `unifi` *or* provide a
   `unifi.local` DNS entry by default.

## API Support Status

UniFi Network API V1 is fully supported as of Network application version "9.1.120".

UniFi Protect API is only partially supported, with the following endpoints supported:

- `/v1/meta/info`
- `/v1/subscribe/devices`: only partial type support.
- `/v1/subscribe/events`
- `/v1/cameras/`
- `/v1/cameras/{id}`

UniFi Access API is not supported yet.

# Contributing

Contributions are welcome!

Before submiting any PRs, please ensure your commits build and lint. Also, take
care to test any changes, and ideally submit unit tests or integration tests
which cover your changes. PRs must pass all CI

If you'd like to contribute a feature or API support that doesn't yet exist,
please communicate your intention through a new or existing GitHub issue. Let's
try to avoid duplicating work, and make sure the work aligns with project goals.

If you find a bug, please report it via GitHub issues. The more descriptive you
can be, the better. Please include specific steps to reproduce. Bonus points for
submitting a PR to fix it!

## Building `unified`

Unified is built primarily via [`mage`](https://magefile.org/).

All command examples assume current working directory is the repository root.

To build the `unified` CLI command you would run:

```bash
$ mage buildCmd
```

Which will build `unified` as well as any of its dependencies. If successful,
it should place the binary at `build/unified`.

## Testing

```bash
$ mage test
```
Will run any unit tests available. Which are not many at this point. Any
reasonable unit test contributions are welcome.

If you have a UniFi API host available, you can run integration tests:
```bash
UNIFIED_HAVE_UNIFI_API_HOST=true go test -v ./test/integration/
```
>[!WARNING]
>While designed to be non-destructive to existing application objects and
>configuration, some non-`GET` endpoints are called. Please take a look at
>existing integration tests and verify you're comfortable running them
>against your API host. We are *NOT* resposible for any harm they might
>cause to your network device/control-plane.

## Linting

```bash
$ mage lint
```

# Thanks

Copyright 2025 - Clifton Houck
