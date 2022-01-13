package main

import (
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

	time.Sleep(5 * time.Second)
}
