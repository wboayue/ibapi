package ibapi

type realTimeBarsEncoder struct {
	serverVersion int
	version       int
	requestId     int

	contract   Contract
	whatToShow string
	useRth     bool
}

func (e *realTimeBarsEncoder) encode() string {
	message := messageBuilder{}

	message.addInt(REQ_REAL_TIME_BARS)
	message.addInt(e.version)
	message.addInt(e.requestId)

	if e.serverVersion >= MinServerVer_TRADING_CLASS {
		message.addInt(e.contract.ContractId)
	}

	message.addString(e.contract.Symbol)
	message.addString(e.contract.SecurityType)
	message.addString(e.contract.LastTradeDateOrContractMonth)
	message.addFloat64(e.contract.Strike)
	message.addString(e.contract.Right)
	message.addString(e.contract.Multiplier)
	message.addString(e.contract.Exchange)
	message.addString(e.contract.PrimaryExchange)
	message.addString(e.contract.Currency)
	message.addString(e.contract.LocalSymbol)

	if e.serverVersion >= MinServerVer_TRADING_CLASS {
		message.addString(e.contract.TradingClass)
	}

	message.addInt(5) // required bar size
	message.addString(e.whatToShow)
	message.addBool(e.useRth)

	if e.serverVersion >= MinServerVer_LINKING {
		// realtime bar options
		message.addString("")
	}

	return message.Encode()
}

type tickByTickEncoder struct {
	serverVersion int
	version       int
	requestId     int

	contract      Contract
	tickType      string
	numberOfTicks int
	ignoreSize    bool
}

func (e *tickByTickEncoder) encode() string {
	message := messageBuilder{}

	message.addInt(REQ_TICK_BY_TICK_DATA)
	message.addInt(e.requestId)

	message.addInt(e.contract.ContractId)
	message.addString(e.contract.Symbol)
	message.addString(e.contract.SecurityType)
	message.addString(e.contract.LastTradeDateOrContractMonth)
	message.addFloat64(e.contract.Strike)
	message.addString(e.contract.Right)
	message.addString(e.contract.Multiplier)
	message.addString(e.contract.Exchange)
	message.addString(e.contract.PrimaryExchange)
	message.addString(e.contract.Currency)
	message.addString(e.contract.LocalSymbol)
	message.addString(e.contract.TradingClass)
	message.addString(e.tickType)

	if e.serverVersion > MinServerVer_TICK_BY_TICK_IGNORE_SIZE {
		message.addInt(e.numberOfTicks)
		message.addBool(e.ignoreSize)
	}

	return message.Encode()
}

type contractDetailsEncoder struct {
	serverVersion int
	version       int
	requestId     int

	contract Contract
}

func (e *contractDetailsEncoder) encode() string {
	message := messageBuilder{}

	message.addInt(REQ_CONTRACT_DATA)
	message.addInt(e.version)
	message.addInt(e.requestId)

	message.addInt(e.contract.ContractId)
	message.addString(e.contract.Symbol)
	message.addString(e.contract.SecurityType)
	message.addString(e.contract.LastTradeDateOrContractMonth)
	message.addFloat64(e.contract.Strike)
	message.addString(e.contract.Right)
	message.addString(e.contract.Multiplier)
	message.addString(e.contract.Exchange)
	message.addString(e.contract.PrimaryExchange)
	message.addString(e.contract.Currency)
	message.addString(e.contract.LocalSymbol)
	message.addString(e.contract.TradingClass)
	message.addBool(e.contract.IncludeExpired)
	message.addString(e.contract.SecurityIdType)
	message.addString(e.contract.SecurityId)

	return message.Encode()
}
