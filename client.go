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
	socket               net.Conn
	optionalCapabilities string
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
			log.Println(fields)
		}
	}
}

type parser struct {
	fields []string
}

func (s *parser) readInt() int {
	result := s.fields[0]
	s.fields = s.fields[1:]
	num, err := strconv.Atoi(result)
	if err != nil {
		panic(err)
	}
	return num
}

func (s *parser) readString() string {
	result := s.fields[0]
	s.fields = s.fields[1:]
	return result
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

// whatToShow - TRADES, MIDPOINT, BID, ASK
// useRth - use regular trading hours
func (c *IbClient) RealTimeBars(ctx context.Context, contract Contract, whatToShow string, useRth bool) (chan Bar, error) {
	if c.ServerVersion < MinServerVer_TRADING_CLASS {
		if contract.TradingClass != "" {
			return nil, stacktrace.NewError("server version %d does not support TradingClass or ContractId fields", c.ServerVersion)
		}
	}

	version := 3
	requestId := 4
	message := messageBuilder{}

	message.addInt(REQ_REAL_TIME_BARS)
	message.addInt(version)
	message.addInt(requestId)

	if c.ServerVersion >= MinServerVer_TRADING_CLASS {
		message.addInt(contract.ContractId)
	}

	message.addString(contract.Symbol)
	message.addString(contract.SecurityType)
	message.addString(contract.LastTradeDateOrContractMonth)
	message.addFloat64(contract.Strike)
	message.addString(contract.Right)
	message.addString(contract.Multiplier)
	message.addString(contract.Exchange)
	message.addString(contract.PrimaryExchange)
	message.addString(contract.Currency)
	message.addString(contract.LocalSymbol)

	if c.ServerVersion >= MinServerVer_TRADING_CLASS {
		message.addString(contract.TradingClass)
	}

	message.addInt(5) // required bar size
	message.addString(whatToShow)
	message.addBool(useRth)

	if c.ServerVersion >= MinServerVer_LINKING {
		// realtime bar options
		message.addString("")
	}

	err := c.writeMessage(&message)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error sending request market data message")
	}

	// add listener of client by request id

	// stream -> bars to caller

	return nil, nil
}
