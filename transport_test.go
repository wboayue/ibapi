package ibapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTcpMessageBus_Connect(t *testing.T) {
	bus := TcpMessageBus{}

	host := "localhost"
	port := 4002
	clientId := 200

	err := bus.Connect(host, port, clientId)
	assert.Nil(t, err, "error connecting: %v", err)

	assert.Equal(t, host, bus.host)
	assert.Equal(t, port, bus.port)
	assert.Equal(t, clientId, bus.clientId)

	bus.Close()
}

// type TcpMessageBus struct {
