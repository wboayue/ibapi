package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	contract := ibapi.Contract{
		// LocalSymbol:  "ESH2",
		// LocalSymbol:  "CLG2",
		LocalSymbol:  "6EF2",
		SecurityType: "FUT",
		Currency:     "USD",
		Exchange:     "GLOBEX",
		// Exchange: "NYMEX",
	}

	_, err = client.RealTimeBars(ctx, contract, "TRADES", false)
	if err != nil {
		log.Printf("error connecting: %v", err)
		return
	}

	_, err = client.TickByTickTrades(ctx, contract)
	if err != nil {
		log.Printf("error connecting: %v", err)
		return
	}

	// for bar := range bars {
	// 	fmt.Println(bar)
	// }

	time.Sleep(10 * time.Minute)
}
