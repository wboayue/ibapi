package ibapi

func decodeRequestMarketData() {

}

func decodeTickByTickBidAsk() {

}

func decodeTickByTickBidAllLast() {

}

func decodeContractDetails(serverVersion int, fields []string) ContractDetails {
	details := ContractDetails{}

	scanner := &parser{fields[1:]}

	_ = scanner.readInt()

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
		details.SecIdList = make([]string, secIdListCount)
		for i := 0; i < secIdListCount; i++ {
			tag := scanner.readString()
			value := scanner.readString()

			details.SecIdList[i] = "" + tag + ":" + value
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
