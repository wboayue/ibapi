package ibapi

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ibDateLayout  = "20060102 15:04:05 MST"
	clientVersion = 2
	noRequest     = -1
)

type IbClient struct {
	ServerVersion    int        // IB server version
	ServerTime       time.Time  // IB server time
	NextValidOrderId int        // next valid order id
	ManagedAccounts  string     // Ids of managed accounts
	MessageBus       MessageBus // bus used to communicate with server

	currentRequestId int                   // used to generate sequence of request Ids
	channels         map[int]chan []string // message exchange
	ready            chan struct{}

	mu                   sync.Mutex
	requestIdMutex       sync.Mutex
	contractDetailsMutex sync.Mutex
}

type MessageBus interface {
	ReadPacket() (string, error)
	Write(string) error
	WritePacket(string) error
	Close() error
}

// Connect creates a socket connection to TWS/IBG.
//
// Parameters:
// 	host 	- host to connect to
// 	port 	- port to connect to
// 	client 	- client id. can connect up to 32 clients
func Connect(host string, port int, clientId int) (*IbClient, error) {
	bus := TcpMessageBus{}
	if err := bus.Connect(host, port, clientId); err != nil {
		return nil, err
	}

	client := IbClient{
		MessageBus: &bus,
		channels:   make(map[int]chan []string),
	}

	if err := client.handshake(); err != nil {
		return nil, err
	}

	log.Println("sent handshake")

	if err := client.startApi(clientId); err != nil {
		return nil, err
	}

	client.ready = make(chan struct{})

	go client.processMessages()

	<-client.ready

	return &client, nil
}

func (c *IbClient) handshake() error {
	prefix := "API\x00"
	version := fmt.Sprintf("v%d..%d", minClientVer, maxClientVer)

	if err := c.MessageBus.Write(prefix); err != nil {
		return fmt.Errorf("error sending prefix: %w", err)
	}

	if err := c.MessageBus.WritePacket(version); err != nil {
		return fmt.Errorf("error sending version: %w", err)
	}

	fields, err := c.readFirstPacket()
	if err != nil {
		return fmt.Errorf("error reading first packet: %w", err)
	}

	c.ServerVersion, err = strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("error parsing server version %v: %w", fields[0], err)
	}
	log.Printf("server version: %d", c.ServerVersion)

	c.ServerTime, err = time.Parse(ibDateLayout, fields[1])
	if err != nil {
		return fmt.Errorf("error parsing server time %v: %w", fields[1], err)
	}
	log.Printf("server time: %s", c.ServerTime)

	return nil
}

func (c *IbClient) Close() error {
	if c.MessageBus != nil {
		return c.MessageBus.Close()
	}

	return nil
}

func (c *IbClient) nextRequestId() int {
	c.requestIdMutex.Lock()
	defer c.requestIdMutex.Unlock()

	tmp := c.currentRequestId
	c.currentRequestId++

	return tmp + 9000
}

func (c *IbClient) readFields() ([]string, error) {
	data, err := c.MessageBus.ReadPacket()
	if err != nil {
		return nil, fmt.Errorf("error reading packet: %w", err)
	}
	return strings.Split(string(data[:len(data)-1]), "\x00"), nil
}

func (c *IbClient) readFirstPacket() ([]string, error) {
	fields, err := c.readFields()
	if err != nil {
		return nil, fmt.Errorf("error reading fields: %w", err)
	}

	if len(fields) != 2 {
		for _, field := range fields {
			fmt.Println("-" + field)
		}
		return nil, fmt.Errorf("expected 2 fields, got %d: %v", len(fields), fields)
	}

	return fields, nil
}

func (c *IbClient) startApi(clientId int) error {
	msg := fmt.Sprintf("%d\x00%d\x00%d\x00", startApi, clientVersion, clientId)
	if c.ServerVersion > minServerVerOptionalCapabilities {
		msg = msg + "\x00"
	}
	return c.MessageBus.WritePacket(msg)
}

func (c *IbClient) processMessages() {
	for {
		fields, err := c.readFields()
		if err != nil {
			log.Printf("error reading: %v\n", err)
			break
		}

		msgId, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Printf("error parsing: %v\n", err)
			continue
		}

		scanner := &parser{fields[1:]}

		switch msgId {
		case endConn:
			log.Println("connection ended")
			panic("connection ended")
		case nextValidId:
			c.handleNextValidId(scanner)
		case managedAccounts:
			c.handleManagedAccounts(scanner)
		case errMsg:
			c.handleErrorMessage(scanner, fields)
		default:
			requestId := getRequestId(msgId, fields)

			channel := c.getChannel(requestId)
			if channel == nil {
				log.Printf("no receiver found for request id %d: %v", requestId, fields)
				continue
			}

			channel <- fields
		}
	}
}

func getRequestId(msgId int, fields []string) int {
	text := ""

	switch msgId {
	case contractData, tickByTick:
		text = fields[1]
	case contractDataEnd, realTimeBars:
		text = fields[2]
	default:
		log.Fatalf("could not determine request id for message ID %d: %v\n", msgId, fields)
	}

	requestId, err := strconv.Atoi(text)
	if err != nil {
		panic(err)
	}

	return requestId
}

func (c *IbClient) handleNextValidId(scanner *parser) {
	scanner.readInt() // skip version
	c.NextValidOrderId = scanner.readInt()

	close(c.ready)

	log.Printf("next valid id: %v", c.NextValidOrderId)
}

func (c *IbClient) handleManagedAccounts(scanner *parser) {
	scanner.readInt() // skip version
	c.ManagedAccounts = scanner.readString()

	log.Printf("managed accounts: %v", c.ManagedAccounts)
}

func (c *IbClient) handleErrorMessage(scanner *parser, fields []string) {
	version := scanner.readInt()
	if version < 2 {
		log.Println(scanner.readString())
	} else {
		requestId := scanner.readInt()
		code := scanner.readInt()
		msg := scanner.readString()

		if requestId == noRequest {
			log.Printf("error message[%d]: %s", code, msg)
		} else {
			channel, ok := c.channels[requestId]
			if ok {
				channel <- fields
			} else {
				log.Printf("no receiver found for request id %d:%d: %v", requestId, code, msg)
			}
		}
	}
}

// RealTimeBars requests real time bars.
// Currently, only 5 seconds bars are provided. This request is subject to the same pacing as any historical data request: no more than 60 API queries in more than 600 seconds.
// Real time bars subscriptions are also included in the calculation of the number of Level 1 market data subscriptions allowed in an account.
//
// Parameters:
// 	contract 	- the Contract for which the depth is being requested
// 	whatToShow 	- TRADES, MIDPOINT, BID, ASK
// 	useRth 		- use regular trading hours
func (c *IbClient) RealTimeBars(ctx context.Context, contract Contract, whatToShow string, useRth bool) (<-chan Bar, error) {
	if c.ServerVersion < minServerVersionRealTimeBars {
		return nil, fmt.Errorf("server version %d does not support real time bars", c.ServerVersion)
	}

	if c.ServerVersion < minServerVersionTradingClass {
		return nil, fmt.Errorf("server version %d does not support TradingClass or ContractId fields", c.ServerVersion)
	}

	encoder := realTimeBarsEncoder{
		serverVersion: c.ServerVersion,
		version:       3,
		requestId:     c.nextRequestId(),
		contract:      contract,
		whatToShow:    whatToShow,
		useRth:        useRth,
	}

	messages := c.addChannel(encoder.requestId)

	err := c.MessageBus.WritePacket(encoder.encode())
	if err != nil {
		return nil, fmt.Errorf("error sending request market data message: %w", err)
	}

	// process response

	bars := make(chan Bar)

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.cancelRealTimeBars(ctx, encoder.requestId)
				c.removeChannel(encoder.requestId)
				close(bars)
				return

			case message := <-messages:
				if message == nil {
					close(bars)
					return
				}

				messageId, err := strconv.Atoi(message[0])
				if err != nil {
					log.Printf("error parsing messageId [%s]: %v", message[0], err)
				}

				if messageId == realTimeBars {
					bar := decodeRealTimeBars(message)
					bars <- bar
				} else {
					log.Printf("unexpected message: %v", message)
				}
			}
		}
	}()

	return bars, nil
}

// cancelRealTimeBar cancels a request for real time bars.
func (c *IbClient) cancelRealTimeBars(ctx context.Context, requestId int) error {
	if c.ServerVersion < minServerVersionRealTimeBars {
		return fmt.Errorf("server version %d does not support real time bars cancellation", c.ServerVersion)
	}

	log.Printf("canceling real time bar request %v.", requestId)

	message := messageBuilder{}

	version := 1
	message.addInt(cancelRealTimeBars)
	message.addInt(version)
	message.addInt(requestId)

	// interface for this
	if err := c.MessageBus.WritePacket(message.Encode()); err != nil {
		return fmt.Errorf("error sending request to cancel market data: %w", err)
	}

	return nil
}

// TickByTickTrades requests tick by tick trades.
func (c *IbClient) TickByTickTrades(ctx context.Context, contract Contract) (chan Trade, error) {
	if c.ServerVersion < minServerVerTickByTick {
		return nil, fmt.Errorf("server version %d does not support tick-by-tick data requests", c.ServerVersion)
	}

	if c.ServerVersion < minServerVerTickByTickIgnoreSize {
		return nil, fmt.Errorf("server version %d does not support ignore_size and number_of_ticks parameters in tick-by-tick data requests", c.ServerVersion)
	}

	encoder := tickByTickEncoder{
		serverVersion: c.ServerVersion,
		requestId:     c.nextRequestId(),
		contract:      contract,
		tickType:      "AllLast",
		numberOfTicks: 0,
		ignoreSize:    false,
	}

	messages := c.addChannel(encoder.requestId)

	err := c.MessageBus.WritePacket(encoder.encode())
	if err != nil {
		return nil, fmt.Errorf("error sending request for tick by tick trades: %w", err)
	}

	// process response

	trades := make(chan Trade)

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.cancelTickByTickData(ctx, encoder.requestId)
				c.removeChannel(encoder.requestId)
				close(trades)
				return

			case message := <-messages:
				if message == nil {
					close(trades)
					return
				}

				messageId, err := strconv.Atoi(message[0])
				if err != nil {
					log.Printf("error parsing messageId [%s]: %v", message[0], err)
				}

				if messageId == tickByTick {
					trade := decodeTickByTickTrade(c.ServerVersion, message)
					trades <- trade
				} else {
					log.Printf("unexpected message: %v", message)
				}
			}
		}
	}()

	return trades, nil
}

// cancelTickByTickData cancels a request for tick by tick data.
func (c *IbClient) cancelTickByTickData(ctx context.Context, requestId int) error {
	if c.ServerVersion < minServerVerTickByTick {
		return fmt.Errorf("server version %d does not support tick by tick cancellation", c.ServerVersion)
	}

	log.Printf("canceling tick by tick data request %v.", requestId)

	message := messageBuilder{}

	message.addInt(cancelTickByTickData)
	message.addInt(requestId)

	if err := c.MessageBus.WritePacket(message.Encode()); err != nil {
		return fmt.Errorf("error sending request to cancel tick by tick data: %w", err)
	}

	return nil
}

// TickByTickBidAsk requests tick-by-tick bid/ask.
func (c *IbClient) TickByTickBidAsk(ctx context.Context, contract Contract) (chan BidAsk, error) {
	if c.ServerVersion < minServerVerTickByTick {
		return nil, fmt.Errorf("server version %d does not support tick-by-tick data requests", c.ServerVersion)
	}

	if c.ServerVersion < minServerVerTickByTickIgnoreSize {
		return nil, fmt.Errorf("server version %d does not support ignore_size and number_of_ticks parameters in tick-by-tick data requests", c.ServerVersion)
	}

	encoder := tickByTickEncoder{
		serverVersion: c.ServerVersion,
		requestId:     c.nextRequestId(),
		contract:      contract,
		tickType:      "BidAsk",
		numberOfTicks: 0,
		ignoreSize:    false,
	}

	messages := c.addChannel(encoder.requestId)

	err := c.MessageBus.WritePacket(encoder.encode())
	if err != nil {
		return nil, fmt.Errorf("error sending request for tick by tick bid/ask: %w", err)
	}

	// process response

	spreads := make(chan BidAsk)

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.cancelTickByTickData(ctx, encoder.requestId)
				c.removeChannel(encoder.requestId)
				close(spreads)
				return

			case message := <-messages:
				if message == nil {
					close(spreads)
					return
				}

				messageId, err := strconv.Atoi(message[0])
				if err != nil {
					log.Printf("error parsing messageId [%s]: %v", message[0], err)
				}

				if messageId == tickByTick {
					spread := decodeTickByTickBidAsk(c.ServerVersion, message)
					spreads <- spread
				} else if messageId == errMsg {
					log.Printf("error: %v", message)
					close(spreads)
					return
				} else {
					log.Printf("unexpected message: %v", message)
				}
			}
		}
	}()

	return spreads, nil
}

// ContractDetails requests contract information.
// This method will provide all the contracts matching the contract provided.
// It can also be used to retrieve complete options and futures chains.
func (c *IbClient) ContractDetails(ctx context.Context, contract Contract) ([]ContractDetails, error) {
	if c.ServerVersion < minServerVersionSecurityIdType {
		return nil, fmt.Errorf("server version %d does not support SecurityIdType or SecurityId fields", c.ServerVersion)
	}

	if c.ServerVersion < minServerVersionTradingClass {
		return nil, fmt.Errorf("server version %d does not support TradingClass field in Contract", c.ServerVersion)
	}

	if c.ServerVersion < minServerVersionLinking {
		return nil, fmt.Errorf("server version %d does not support PrimaryExchange field in Contract", c.ServerVersion)
	}

	c.contractDetailsMutex.Lock()
	defer c.contractDetailsMutex.Unlock()

	// create and send request

	encoder := contractDetailsEncoder{
		serverVersion: c.ServerVersion,
		version:       8,
		requestId:     c.nextRequestId(),
		contract:      contract,
	}

	messages := c.addChannel(encoder.requestId)

	err := c.MessageBus.WritePacket(encoder.encode())
	if err != nil {
		return nil, fmt.Errorf("error sending contract details request: %w", err)
	}

	// process response

	contracts := []ContractDetails{}

	for {
		select {
		case <-ctx.Done():
			c.removeChannel(encoder.requestId)
			return contracts, fmt.Errorf("contract details request %d cancelled", encoder.requestId)

		case message := <-messages:
			if message == nil {
				return contracts, nil
			}

			messageId, err := strconv.Atoi(message[0])
			if err != nil {
				log.Printf("error parsing messageId [%s]: %v", message[0], err)
			}

			if messageId == contractDataEnd {
				c.removeChannel(encoder.requestId)
			} else if messageId == contractData {
				contract := decodeContractDetails(c.ServerVersion, message)
				contracts = append(contracts, contract)
			} else {
				log.Printf("unexpected message: %v", message)
			}
		}
	}
}

// Utility Methods

func (c *IbClient) addChannel(requestId int) chan []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel := make(chan []string)
	c.channels[requestId] = channel

	return channel
}

func (c *IbClient) removeChannel(requestId int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel := c.channels[requestId]
	if channel != nil {
		delete(c.channels, requestId)
		close(channel)
	}
}

func (c *IbClient) getChannel(requestId int) chan []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.channels[requestId]
}
