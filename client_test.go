package ibapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Connect(t *testing.T) {

	client := IbClient{}

	err := client.Connect("localhost", 4002, 100)

	assert.Nil(t, err, "error connecting: %v", err)
	assert.Equal(t, 10, client.ServerVersion)
	assert.Equal(t, "", client.ConnectOptions)

	client.Close()
}
