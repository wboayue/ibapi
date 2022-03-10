package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wboayue/ibapi"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := ibapi.Connect("127.0.0.1", 4002, 200)
	if err != nil {
		log.Printf("error connecting: %v", err)
		return
	}

	defer client.Close()

	fmt.Printf("server version: %v\n", client.ServerVersion)
	fmt.Printf("server time: %v\n", client.ServerTime)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// realTimeBars(ctx, client)
	// contractDetails(ctx, client)
	// tickByTickTrades(ctx, client)
	// tickByTickSpreads(ctx, client)
	tickByTick(ctx, client)
	fmt.Println("done")
}

func realTimeBars(ctx context.Context, client *ibapi.IbClient) {
	contract := ibapi.Contract{
		LocalSymbol: "ESJ2",
		// LocalSymbol:  "6EF2",
		SecurityType: "FUT",
		Currency:     "USD",
		Exchange:     "GLOBEX",
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	bars, err := client.RealTimeBars(ctx, contract, "TRADES", false)
	if err != nil {
		log.Printf("error connecting: %v", err)
		return
	}

	for bar := range bars {
		fmt.Printf("bar: %+v\n", bar)
	}
}

func tickByTickTrades(ctx context.Context, client *ibapi.IbClient) {
	contract := ibapi.Contract{
		// LocalSymbol: "ESH2",
		// Exchange:     "GLOBEX",
		SecurityType: "FUT",
		Currency:     "USD",
		// Exchange:     "GLOBEX",
		// Exchange: "NYMEX",
		LocalSymbol: "NMH2",
		Exchange:    "GLOBEX",
	}

	trades, err := client.TickByTickTrades(ctx, contract)
	if err != nil {
		log.Printf("error connecting: %v", err)
		return
	}

	for trade := range trades {
		fmt.Printf("trade: %+v\n", trade)
	}
}

func contractDetails(ctx context.Context, client *ibapi.IbClient) {
	contract := ibapi.Contract{
		Symbol:                       "ES",
		SecurityType:                 "FUT",
		Currency:                     "USD",
		LastTradeDateOrContractMonth: "2022",
	}

	contracts, err := client.ContractDetails(ctx, contract)
	if err != nil {
		log.Printf("error connecting: %v", err)
		return
	}

	for i, contract := range contracts {
		fmt.Printf("%d - %+v\n", i, contract)
	}
}

func tickByTickSpreads(ctx context.Context, client *ibapi.IbClient) {
	contract := ibapi.Contract{
		LocalSymbol: "ESH2",
		// LocalSymbol:  "CLG2",
		SecurityType: "FUT",
		Currency:     "USD",
		Exchange:     "GLOBEX",
		// Exchange: "NYMEX",
	}

	spreads, err := client.TickByTickBidAsk(ctx, contract)
	if err != nil {
		log.Printf("error connecting: %v", err)
		return
	}

	for spread := range spreads {
		fmt.Printf("bid/ask: %+v\n", spread)
	}
}

func tickByTick(ctx context.Context, client *ibapi.IbClient) {
	contract := ibapi.Contract{
		Symbol: "ES",
		// LocalSymbol:  "ESH2",
		SecurityType:                 "FUT",
		Currency:                     "USD",
		Exchange:                     "GLOBEX",
		LastTradeDateOrContractMonth: "202203",
	}

	log.Printf("stream tick")

	// contract.LastTradeDateOrContractMonth = "201803";
	spreads, err := client.TickByTickBidAsk(ctx, contract)
	if err != nil {
		log.Printf("error connecting: %v", err)
		return
	}

	trades, err := client.TickByTickTrades(ctx, contract)
	if err != nil {
		log.Printf("error connecting: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("done")
			return

		case spread := <-spreads:
			fmt.Printf("spread: %+v\n", spread)

		case trade := <-trades:
			fmt.Printf("trade: %+v\n", trade)
		}
	}
}
