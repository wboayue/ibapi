package ibapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRealTimeBarsEncoder(t *testing.T) {
	contract := Contract{
		Symbol:                       "ES",
		LocalSymbol:                  "ESU6",
		SecurityType:                 "FUT",
		Currency:                     "USD",
		Exchange:                     "GLOBEX",
		TradingClass:                 "FGBL",
		LastTradeDateOrContractMonth: "201803",
	}

	t.Run("with trading class available", func(t *testing.T) {
		request := realTimeBarsEncoder{
			serverVersion: minServerVersionLinking,
			version:       3,
			requestId:     1,
			contract:      contract,
			whatToShow:    "TRADES",
			useRth:        true,
		}

		assert.Equal(t, "50\x003\x001\x000\x00ES\x00FUT\x00201803\x000.000000\x00\x00\x00GLOBEX\x00\x00USD\x00ESU6\x00FGBL\x005\x00TRADES\x001\x00\x00", request.encode())
	})

	// t.Run("without trading class available", func(t *testing.T) {
	// 	request := realTimeBarsEncoder{
	// 		serverVersion: MinServerVersionTradingClass - 5,
	// 		version:       3,
	// 		requestId:     1,
	// 		contract:      contract,
	// 		whatToShow:    "TRADES",
	// 		useRth:        true,
	// 	}

	// 	assert.Equal(t, "50\x003\x001\x00ES\x00FUT\x00201803\x000.000000\x00\x00\x00GLOBEX\x00\x00USD\x00ESU6\x005\x00TRADES\x001\x00", request.encode())
	// })
}

func TestTickByTickEncoder(t *testing.T) {
	contract := Contract{
		Symbol:                       "ES",
		LocalSymbol:                  "ESU6",
		SecurityType:                 "FUT",
		Currency:                     "USD",
		Exchange:                     "GLOBEX",
		TradingClass:                 "FGBL",
		LastTradeDateOrContractMonth: "201803",
	}

	t.Run("tick by tick trades", func(t *testing.T) {
		request := tickByTickEncoder{
			serverVersion: minServerVerTickByTickIgnoreSize,
			requestId:     2,
			contract:      contract,
			tickType:      "AllLast",
			numberOfTicks: 0,
			ignoreSize:    false,
		}

		assert.Equal(t, "97\x002\x000\x00ES\x00FUT\x00201803\x000.000000\x00\x00\x00GLOBEX\x00\x00USD\x00ESU6\x00FGBL\x00AllLast\x00", request.encode())
	})

	t.Run("tick by tick top of book", func(t *testing.T) {
		request := tickByTickEncoder{
			serverVersion: minServerVerTickByTick,
			requestId:     2,
			contract:      contract,
			tickType:      "BidAsk",
			numberOfTicks: 0,
			ignoreSize:    true,
		}

		assert.Equal(t, "97\x002\x000\x00ES\x00FUT\x00201803\x000.000000\x00\x00\x00GLOBEX\x00\x00USD\x00ESU6\x00FGBL\x00BidAsk\x00", request.encode())
	})
}

func TestContractDetailsEncoder(t *testing.T) {
	contract := Contract{
		Symbol:                       "ES",
		SecurityType:                 "FUT",
		Currency:                     "USD",
		LastTradeDateOrContractMonth: "2021",
	}

	request := contractDetailsEncoder{
		serverVersion: minServerVersionLinking,
		requestId:     2,
		contract:      contract,
	}

	assert.Equal(t, "9\x000\x002\x000\x00ES\x00FUT\x002021\x000.000000\x00\x00\x00\x00\x00USD\x00\x00\x000\x00\x00\x00", request.encode())
}
