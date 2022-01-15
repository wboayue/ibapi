package main

import (
	"context"
	"fmt"
	"log"

	"github.com/wboayue/ibapi"
)

func main() {
	client := ibapi.IbClient{}
	defer client.Close()

	err := client.Connect("localhost", 4002, 100)
	if err != nil {
		log.Printf("error connecting: %v", err)
		return
	}

	fmt.Printf("server version: %v\n", client.ServerVersion)
	fmt.Printf("server time: %v\n", client.ServerTime)

	go client.ProcessMessages()

	ctx := context.Background()

	// realTimeBars(ctx, &client)
	contractDetails(ctx, &client)

	// time.Sleep(10 * time.Minute)
}

func realTimeBars(ctx context.Context, client *ibapi.IbClient) {
	contract := ibapi.Contract{
		LocalSymbol:  "ESH2",
		SecurityType: "FUT",
		Currency:     "USD",
		Exchange:     "GLOBEX",
	}

	_, err := client.RealTimeBars(ctx, contract, "TRADES", false)
	if err != nil {
		log.Printf("error connecting: %v", err)
		return
	}
}

func tickByTickTrades(ctx context.Context, client *ibapi.IbClient) {
	contract := ibapi.Contract{
		// LocalSymbol:  "ESH2",
		// LocalSymbol:  "CLG2",
		LocalSymbol:  "6EF2",
		SecurityType: "FUT",
		Currency:     "USD",
		Exchange:     "GLOBEX",
		// Exchange: "NYMEX",
	}

	bars, err := client.TickByTickTrades(ctx, contract)
	if err != nil {
		log.Printf("error connecting: %v", err)
		return
	}

	for bar := range bars {
		fmt.Println(bar)
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

	fmt.Printf("contracts: %+v\n", contracts)
}
