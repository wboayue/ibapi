package ibapi

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecodeRealTimeBars(t *testing.T) {
	// Assemble

	messageId := RealTimeBars
	version := 3
	requestId := 9000
	timestamp := int64(1642465785)
	open := 4658.00
	high := 4658.25
	low := 4658.00
	close := 4658.00
	volume := int64(5)
	wap := 4658.05
	count := 3

	packet := []string{
		fmt.Sprintf("%d", messageId),
		fmt.Sprintf("%d", version),
		fmt.Sprintf("%d", requestId),
		fmt.Sprintf("%d", timestamp),
		fmt.Sprintf("%f", open),
		fmt.Sprintf("%f", high),
		fmt.Sprintf("%f", low),
		fmt.Sprintf("%f", close),
		fmt.Sprintf("%d", volume),
		fmt.Sprintf("%f", wap),
		fmt.Sprintf("%d", count),
	}

	// Activate

	bar := decodeRealTimeBars(packet)

	// Assert

	assert.Equal(t, time.Unix(timestamp, 0), bar.Time)
	assert.Equal(t, open, bar.Open)
	assert.Equal(t, high, bar.High)
	assert.Equal(t, low, bar.Low)
	assert.Equal(t, close, bar.Close)
	assert.Equal(t, volume, bar.Volume)
	assert.Equal(t, wap, bar.WAP)
	assert.Equal(t, count, bar.Count)
}
