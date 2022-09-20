package ibapi

// Decoders convert raw messages into a structured responses

import (
	"log"
	"time"
)

// decodeRealTimeBars converts a RealTimeBars incoming message into a Bar
func decodeRealTimeBars(fields []string) Bar {
	scanner := &parser{fields[3:]}

	return Bar{
		Time:   time.Unix(scanner.readInt64(), 0),
		Open:   scanner.readFloat64(),
		High:   scanner.readFloat64(),
		Low:    scanner.readFloat64(),
		Close:  scanner.readFloat64(),
		Volume: scanner.readInt64(),
		WAP:    scanner.readFloat64(),
		Count:  scanner.readInt(),
	}
}

func decodeTickByTickBidAsk(serverVersion int, fields []string) BidAsk {
	scanner := &parser{fields[3:]}

	timestamp := scanner.readInt64()

	bidPrice := scanner.readFloat64()
	askPrice := scanner.readFloat64()
	bidSize := scanner.readInt64()
	askSize := scanner.readInt64()

	mask := scanner.readInt()
	attribute := BidAskAttribute{
		BidPastLow:  mask&0x1 == 0x1,
		AskPastHigh: mask&0x2 == 0x2,
	}

	return BidAsk{
		Time:            time.Unix(timestamp, 0),
		BidPrice:        bidPrice,
		AskPrice:        askPrice,
		BidSize:         bidSize,
		AskSize:         askSize,
		BidAskAttribute: attribute,
	}
}

func decodeTickByTickTrade(serverVersion int, fields []string) Trade {
	scanner := &parser{fields[2:]}

	tickType := scanner.readInt()
	timestamp := scanner.readInt64()

	if tickType != 2 {
		log.Printf("expected tick type 2, got: %v", tickType)
		return Trade{}
	}

	price := scanner.readFloat64()
	size := scanner.readInt64()
	mask := scanner.readInt()
	attribute := TradeAttribute{
		PastLimit:  mask&0x1 == 0x1,
		Unreported: mask&0x2 == 0x2,
	}
	exchange := scanner.readString()
	specialConditions := scanner.readString()

	return Trade{
		TickType:          "Last",
		Time:              time.Unix(timestamp, 0),
		Price:             price,
		Size:              size,
		TradeAttribute:    attribute,
		Exchange:          exchange,
		SpecialConditions: specialConditions,
	}
}

func decodeContractDetails(serverVersion int, fields []string) ContractDetails {
	details := ContractDetails{}

	scanner := &parser{fields[2:]}

	details.Contract.Symbol = scanner.readString()
	details.Contract.SecurityType = scanner.readString()
	details.Contract.LastTradeDateOrContractMonth = scanner.readString()
	details.Contract.Strike = scanner.readFloat64()
	details.Contract.Right = scanner.readString()
	details.Contract.Exchange = scanner.readString()
	details.Contract.Currency = scanner.readString()
	details.Contract.LocalSymbol = scanner.readString()
	details.MarketName = scanner.readString()
	details.Contract.TradingClass = scanner.readString()
	details.Contract.ContractId = scanner.readInt()
	details.MinTick = scanner.readFloat64()
	details.Contract.Multiplier = scanner.readString()
	details.OrderTypes = scanner.readString()
	details.ValidExchanges = scanner.readString()
	details.PriceMagnifier = scanner.readInt()
	details.UnderContractId = scanner.readInt()
	details.LongName = scanner.readString()
	details.Contract.PrimaryExchange = scanner.readString()
	details.ContractMonth = scanner.readString()
	details.Industry = scanner.readString()
	details.Category = scanner.readString()
	details.Subcategory = scanner.readString()
	details.TimeZoneId = scanner.readString()
	details.TradingHours = scanner.readString()
	details.LiquidHours = scanner.readString()
	details.EvRule = scanner.readString()
	details.EvMultiplier = scanner.readInt()

	secIdListCount := scanner.readInt()
	if secIdListCount > 0 {
		details.SecIdList = make([]TagValue, secIdListCount)
		for i := 0; i < secIdListCount; i++ {
			tag := scanner.readString()
			value := scanner.readString()

			details.SecIdList[i] = TagValue{Tag: tag, Value: value}
		}
	}

	details.AggGroup = scanner.readInt()
	details.UnderSymbol = scanner.readString()
	details.UnderSecType = scanner.readString()
	details.MarketRuleIds = scanner.readString()
	details.RealExpirationDate = scanner.readString()
	details.StockType = scanner.readString()

	details.MinSize = scanner.readFloat64()
	details.SizeIncrement = scanner.readFloat64()
	details.SuggestedSizeIncrement = scanner.readFloat64()

	return details
}
