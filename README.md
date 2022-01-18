# Interactive Brokers API - Go

This is a partial, unofficial and idomatic Go implementation of the Interactive Brokers API.
This implementation focuses on parts of the API used by my application. The focus is currently realtime market data. Roadmap is to implement APIs to support order execution and account management next.

The official API has some required complexity because it needs to support old version of the IB Gateway. This implementation make some simplifications and was implemented and test only against the latest version of the IB Gateway.

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
