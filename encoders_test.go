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
			serverVersion: MinServerVer_TRADING_CLASS,
			version:       3,
			requestId:     1,
			contract:      contract,
			whatToShow:    "TRADES",
			useRth:        true,
		}

		assert.Equal(t, "50\x003\x001\x000\x00ES\x00FUT\x00201803\x000.000000\x00\x00\x00GLOBEX\x00\x00USD\x00ESU6\x00FGBL\x005\x00TRADES\x001\x00", request.encode())
	})

	t.Run("without trading class available", func(t *testing.T) {
		request := realTimeBarsEncoder{
			serverVersion: MinServerVer_TRADING_CLASS - 5,
			version:       3,
			requestId:     1,
			contract:      contract,
			whatToShow:    "TRADES",
			useRth:        true,
		}

		assert.Equal(t, "50\x003\x001\x00ES\x00FUT\x00201803\x000.000000\x00\x00\x00GLOBEX\x00\x00USD\x00ESU6\x005\x00TRADES\x001\x00", request.encode())
	})
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
			serverVersion: MinServerVer_TICK_BY_TICK_IGNORE_SIZE,
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
			serverVersion: MinServerVer_TICK_BY_TICK,
			requestId:     2,
			contract:      contract,
			tickType:      "BidAsk",
			numberOfTicks: 0,
			ignoreSize:    true,
		}

		assert.Equal(t, "97\x002\x000\x00ES\x00FUT\x00201803\x000.000000\x00\x00\x00GLOBEX\x00\x00USD\x00ESU6\x00FGBL\x00BidAsk\x00", request.encode())
	})
}
