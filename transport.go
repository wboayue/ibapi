package ibapi

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/palantir/stacktrace"
)

// TcpMessageBus implements the MessageBus over TCP
type TcpMessageBus struct {
	host     string
	port     int
	clientId int
	socket   net.Conn
}

// Connect establises a connection to the remote host
func (b *TcpMessageBus) Connect(host string, port int, clientId int) error {
	b.host = host
	b.port = port
	b.clientId = clientId

	var err error
	b.socket, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return stacktrace.Propagate(err, "error dialing %s:%d", host, port)
	}

	return nil
}

// Close closes the network connection
func (b *TcpMessageBus) Close() error {
	if b.socket != nil {
		return b.socket.Close()
	}

	return nil
}

// WritePacket writes raw data to message bus
func (b *TcpMessageBus) Write(data string) error {
	_, err := b.socket.Write([]byte(data))
	if err != nil {
		return stacktrace.Propagate(err, "error writing bytes")
	}

	return nil
}

// WritePacket writes a packet of data to message bus adding the length prefix
func (b *TcpMessageBus) WritePacket(data string) error {
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))

	_, err := b.socket.Write(header)
	if err != nil {
		return stacktrace.Propagate(err, "error writing packet")
	}

	_, err = b.socket.Write([]byte(data))
	if err != nil {
		return err
	}

	return nil
}

// ReadPacket reads the next data packet from the message bus
func (b *TcpMessageBus) ReadPacket() (string, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(b.socket, header)
	if err != nil {
		return "", stacktrace.Propagate(err, "error reading packet header")
	}

	size := binary.BigEndian.Uint32(header)

	data := make([]byte, size)
	_, err = io.ReadFull(b.socket, data)
	if err != nil {
		return "", stacktrace.Propagate(err, "error reading packet body")
	}

	return string(data), nil
}

// MessageBusRecorder records the MessageBus interactions
type MessageBusRecorder struct {
	Bus MessageBus
}

func (b *MessageBusRecorder) ReadPacket() (string, error) {
	return b.Bus.ReadPacket()
}

func (b *MessageBusRecorder) Write(data string) error {
	return b.Bus.Write(data)
}

func (b *MessageBusRecorder) WritePacket(packet string) error {
	return b.Bus.WritePacket(packet)
}

func (b *MessageBusRecorder) Close() error {
	return b.Bus.Close()
}
