package ibapi

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/palantir/stacktrace"
)

const (
	IbDateLayout  = "20060102 15:04:05 MST"
	ClientVersion = 2
)

type IbClient struct {
	ConnectOptions  string
	ServerVersion   int
	ServerTime      time.Time
	NextValidId     int
	ManagedAccounts string

	clientId             int
	requestId            int
	socket               net.Conn
	optionalCapabilities string
	channels             map[int]chan []string
}

func NewClient() *IbClient {
	return &IbClient{}
}

type Message interface {
	Encode() []byte
}

func (c *IbClient) Connect(host string, port int, clientId int) error {
	c.clientId = clientId

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	c.socket = conn

	c.channels = make(map[int]chan []string)

	prefix := "API\x00"
	version := fmt.Sprintf("v%d..%d", MinClientVer, MaxClientVer)
	if c.ConnectOptions != "" {
		version = version + " " + c.ConnectOptions
	}

	_, err = c.socket.Write([]byte(prefix))
	if err != nil {
		return err
	}

	err = c.writePacket([]byte(version))
	if err != nil {
		return err
	}

	fields, err := c.readFields()
	if err != nil {
		return stacktrace.Propagate(err, "error reading fields")
	}

	if len(fields) != 2 {
		for _, field := range fields {
			fmt.Println("-" + field)
		}
		return stacktrace.NewError("expected 2 fields, got %d: %v", len(fields), fields)
	}

	c.ServerVersion, err = strconv.Atoi(fields[0])
	if err != nil {
		return stacktrace.Propagate(err, "error parsing server version: %v", fields[0])
	}

	c.ServerTime, err = time.Parse(IbDateLayout, fields[1])
	if err != nil {
		return stacktrace.Propagate(err, "error parsing server time: %v", fields[1])
	}

	return c.startApi()
}

func (c *IbClient) Close() error {
	if c.socket != nil {
		return c.socket.Close()
	}

	return nil
}

func (c *IbClient) nextRequestId() int {
	tmp := c.requestId
	c.requestId++
	return tmp + 9000
}

func (c *IbClient) writePacket(data []byte) error {
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))

	_, err := c.socket.Write(header)
	if err != nil {
		return err
	}

	_, err = c.socket.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *IbClient) writeMessage(msg Message) error {
	return c.writePacket(msg.Encode())
}

func (c *IbClient) readPacket() ([]byte, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(c.socket, header)
	if err != nil {
		return nil, err
	}

	size := binary.BigEndian.Uint32(header)

	data := make([]byte, size)
	_, err = io.ReadFull(c.socket, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *IbClient) readFields() ([]string, error) {
	data, err := c.readPacket()
	if err != nil {
		return nil, stacktrace.Propagate(err, "error reading packet")
	}
	return strings.Split(string(data[:len(data)-1]), "\x00"), nil
}

func (c *IbClient) startApi() error {
	msg := fmt.Sprintf("%d\x00%d\x00%d\x00", StartApi, ClientVersion, c.clientId)
	if c.ServerVersion > MinServerVerOptionalCapabilities {
		msg = msg + fmt.Sprintf("%s\x00", c.optionalCapabilities)
	}
	return c.writePacket([]byte(msg))
}

func (c *IbClient) ProcessMessages() {
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
		case EndConn:
			log.Println("connection ended")
			return
		case NextValidId:
			c.handleNextValidId(scanner)
		case ManagedAccounts:
			c.handleManagedAccounts(scanner)
		case ErrMsg:
			c.handleErrorMessage(scanner)
		default:
			requestId := getRequestId(msgId, fields)

			channel, ok := c.channels[requestId]
			if ok {
				channel <- fields
			} else {
				log.Printf("no receiver found for request id %d: %v", requestId, fields)
			}
		}
	}
}

func getRequestId(msgId int, fields []string) int {
	text := ""

	switch msgId {
	case ContractData:
		text = fields[1]
	case ContractDataEnd:
		text = fields[2]
	}

	requestId, err := strconv.Atoi(text)
	if err != nil {
		panic(err)
	}

	return requestId
}

func (c *IbClient) handleNextValidId(scanner *parser) {
	scanner.readInt() // version
	c.NextValidId = scanner.readInt()

	log.Printf("next valid id: %v", c.NextValidId)
}

func (c *IbClient) handleManagedAccounts(scanner *parser) {
	scanner.readInt() // version
	c.ManagedAccounts = scanner.readString()

	log.Printf("managed accounts: %v", c.ManagedAccounts)
}

func (c *IbClient) handleErrorMessage(scanner *parser) {
	version := scanner.readInt()
	if version < 2 {
		log.Println(scanner.readString())
	} else {
		id := scanner.readInt()
		code := scanner.readInt()
		msg := scanner.readString()

		log.Printf("id: %d, code: %d, msg: %s", id, code, msg)
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
	if c.ServerVersion < MinServerVer_REAL_TIME_BARS {
		return nil, stacktrace.NewError("server version %d does not support real time bars", c.ServerVersion)
	}

	if c.ServerVersion < MinServerVersionTradingClass {
		return nil, stacktrace.NewError("server version %d does not support TradingClass or ContractId fields", c.ServerVersion)
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

	err := c.writePacket([]byte(encoder.encode()))
	if err != nil {
		return nil, stacktrace.Propagate(err, "error sending request market data message")
	}

	// process response

	bars := make(chan Bar)

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.cancelRealTimeBars(ctx, encoder.requestId)
				time.Sleep(200 * time.Millisecond)
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

				if messageId == REAL_TIME_BARS {
					bar := decodeRealTimeBars(c.ServerVersion, message)
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
	if c.ServerVersion < MinServerVer_REAL_TIME_BARS {
		return stacktrace.NewError("server version %d does not support real time bars cancellation", c.ServerVersion)
	}

	log.Printf("canceling real time bar request %v.", requestId)

	message := messageBuilder{}

	version := 1
	message.addInt(CANCEL_REAL_TIME_BARS)
	message.addInt(version)
	message.addInt(requestId)

	// interface for this
	if err := c.writePacket([]byte(message.Encode())); err != nil {
		return stacktrace.Propagate(err, "error sending request to cancel market data")
	}

	return nil
}

// TickByTickTrades requests tick-by-tick trades.
func (c *IbClient) TickByTickTrades(ctx context.Context, contract Contract) (chan Trade, error) {
	if c.ServerVersion < MinServerVer_TICK_BY_TICK {
		return nil, stacktrace.NewError("server version %d does not support tick-by-tick data requests.", c.ServerVersion)
	}

	if c.ServerVersion < MinServerVer_TICK_BY_TICK_IGNORE_SIZE {
		return nil, stacktrace.NewError("server version %d does not support ignore_size and number_of_ticks parameters in tick-by-tick data requests.", c.ServerVersion)
	}

	encoder := tickByTickEncoder{
		serverVersion: c.ServerVersion,
		requestId:     c.nextRequestId(),
		contract:      contract,
		tickType:      "AllLast",
		numberOfTicks: 0,
		ignoreSize:    false,
	}

	// messages := c.addChannel(encoder.requestId)

	err := c.writePacket([]byte(encoder.encode()))
	if err != nil {
		return nil, stacktrace.Propagate(err, "error sending request for tick by tick trades")
	}

	// add listener of client by request id

	// process response

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		c.removeChannel(encoder.requestId)
	// 		return contracts, fmt.Errorf("contract details request %d cancelled", encoder.requestId)

	// 	case message := <-messages:
	// 		if message == nil {
	// 			return contracts, nil
	// 		}

	// 		messageId, err := strconv.Atoi(message[0])
	// 		if err != nil {
	// 			log.Printf("error parsing messageId [%s]: %v", message[0], err)
	// 		}

	// 		if messageId == ContractDataEnd {
	// 			c.removeChannel(encoder.requestId)
	// 		} else if messageId == ContractData {
	// 			contract := decodeContractDetails(c.ServerVersion, message)
	// 			contracts = append(contracts, contract)
	// 		} else {
	// 			log.Printf("unexpected message: %v", message)
	// 		}
	// 	}
	// }

	return nil, nil
}

// TickByTickBidAsk requests tick-by-tick bid/ask.
func (c *IbClient) TickByTickBidAsk(ctx context.Context, contract Contract) (chan BidAsk, error) {
	if c.ServerVersion < MinServerVer_REAL_TIME_BARS {
		return nil, stacktrace.NewError("server version %d does not support real time bars", c.ServerVersion)
	}

	if c.ServerVersion < MinServerVersionTradingClass {
		if contract.TradingClass != "" {
			return nil, stacktrace.NewError("server version %d does not support TradingClass or ContractId fields", c.ServerVersion)
		}
	}

	encoder := realTimeBarsEncoder{
		serverVersion: c.ServerVersion,
		version:       3,
		requestId:     c.nextRequestId(),
		contract:      contract,
	}

	// messages := c.addChannel(encoder.requestId)

	err := c.writePacket([]byte(encoder.encode()))
	if err != nil {
		return nil, stacktrace.Propagate(err, "error sending request market data message")
	}

	// // process response

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		c.removeChannel(encoder.requestId)
	// 		return contracts, fmt.Errorf("contract details request %d cancelled", encoder.requestId)

	// 	case message := <-messages:
	// 		if message == nil {
	// 			return contracts, nil
	// 		}

	// 		messageId, err := strconv.Atoi(message[0])
	// 		if err != nil {
	// 			log.Printf("error parsing messageId [%s]: %v", message[0], err)
	// 		}

	// 		if messageId == ContractDataEnd {
	// 			c.removeChannel(encoder.requestId)
	// 		} else if messageId == ContractData {
	// 			contract := decodeContractDetails(c.ServerVersion, message)
	// 			contracts = append(contracts, contract)
	// 		} else {
	// 			log.Printf("unexpected message: %v", message)
	// 		}
	// 	}
	// }

	return nil, nil
}

// ContractDetails requests contract information.
// This method will provide all the contracts matching the contract provided.
// It can also be used to retrieve complete options and futures chains.
func (c *IbClient) ContractDetails(ctx context.Context, contract Contract) ([]ContractDetails, error) {
	if c.ServerVersion < MinServerVersionSecurityIdType {
		return nil, stacktrace.NewError("server version %d does not support SecurityIdType or SecurityId fields", c.ServerVersion)
	}

	if c.ServerVersion < MinServerVersionTradingClass {
		return nil, stacktrace.NewError("server version %d does not support TradingClass field in Contract", c.ServerVersion)
	}

	if c.ServerVersion < MinServerVersionLinking {
		return nil, stacktrace.NewError("server version %d does not support PrimaryExchange field in Contract", c.ServerVersion)
	}

	// create and send request

	encoder := contractDetailsEncoder{
		serverVersion: c.ServerVersion,
		version:       8,
		requestId:     c.nextRequestId(),
		contract:      contract,
	}

	messages := c.addChannel(encoder.requestId)

	err := c.writePacket([]byte(encoder.encode()))
	if err != nil {
		return nil, stacktrace.Propagate(err, "error sending contract details request")
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

			if messageId == ContractDataEnd {
				c.removeChannel(encoder.requestId)
			} else if messageId == ContractData {
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
	channel := make(chan []string)
	c.channels[requestId] = channel
	return channel
}

func (c *IbClient) removeChannel(requestId int) {
	channel := c.channels[requestId]
	if channel != nil {
		delete(c.channels, requestId)
		close(channel)
	}
}
