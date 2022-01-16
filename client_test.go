package ibapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Connect(t *testing.T) {
	client, err := Connect("localhost", 4002, 100)

	assert.Nil(t, err, "error connecting: %v", err)
	assert.Equal(t, 10, client.ServerVersion)

	client.Close()
}
