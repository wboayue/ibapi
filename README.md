# Interactive Brokers API - Go

[![test](https://github.com/wboayue/ibapi/workflows/ci/badge.svg)](https://github.com/wboayue/ibapi/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/wboayue/ibapi/branch/main/graph/badge.svg)](https://codecov.io/gh/wboayue/ibapi)

This is a partial implementation of the Interactive Brokers API in Go. This implemention does not provide a one to one match of the official Interactive Brokers API, it attempts to provide an idiomatic Go API to the TWS functionality.

This implementation is a work in progress and focuses on parts of the API used by my application. The initial set of APIs implemented are around realtime market data. Next in line are APIs to support order execution and account management.

This implementation uses the [TWS API 10.12](https://interactivebrokers.github.io/#) source as a reference.

The official API has some required complexity to be backwards compatible with older versions of the IB Gateway. To simplify things the API was implemented and test only against IB Gateway 10.12.

# Installation

```bash
go get -u github.com/wboayue/ibapi
```

# Usage

## Market Data

Connect to gateway

```go
port := 4002
clientId := 100
client, err := ibapi.Connect("localhost", port, client)
if err != nil {
    log.Printf("error connecting: %v", err)
    return
}

defer client.Close()

fmt.Printf("server version: %v\n", client.ServerVersion)
fmt.Printf("server time: %v\n", client.ServerTime)
```

Request real time bars

```go
contract := ibapi.Contract{
    LocalSymbol: "ESH2",
    SecurityType: "FUT",
    Currency:     "USD",
    Exchange:     "GLOBEX",
}

ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, 60*time.Second) // stop streaming after 60 seconds
defer cancel()

bars, err := client.RealTimeBars(ctx, contract, "TRADES", false)
if err != nil {
    log.Printf("error connecting: %v", err)
    return
}

for bar := range bars {
    fmt.Println(bar)
}
```

# Reference

* [API Documentation](https://interactivebrokers.github.io/tws-api/)

# Development

* [Publish Module](https://go.dev/doc/modules/publishing)

```bash
git tag                 # list tags
git tag vx.y.z          # create tag
git push origin vx.y.z  # push tags to origin
GOPROXY=proxy.golang.org go list -m github.com/wboayue/ibapi@vx.y.z     # publish to package repo
```

