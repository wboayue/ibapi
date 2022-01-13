package ibapi

import (
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
	ConnectOptions string
	ServerVersion  int
	ServerTime     time.Time

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
		data, err := c.readPacket()
		if err != nil {
			log.Printf("error reading: %v\n", err)
			break
		}
		log.Println(string(data))
	}
}
