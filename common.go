package ibapi

import "time"

type (
	// Describes an instrument's definition
	Contract struct {
		ContractId int    // The unique IB contract identifier.
		Symbol     string // The underlying's asset symbol.

		// The security's type:
		// STK - stock (or ETF)
		// OPT - option
		// FUT - future
		// IND - index
		// FOP - futures option
		// CASH - forex pair
		// BAG - combo
		// WAR - warrant
		// BOND - bond
		// CMDTY - commodity
		// NEWS - news
		// FUND - mutual fund.
		SecurityType string

		// The contract's last trading day or contract month (for Options and Futures).
		// Strings with format YYYYMM will be interpreted as the Contract Month whereas YYYYMMDD will be interpreted as Last Trading Day.
		LastTradeDateOrContractMonth string

		Strike      float64 // The option's strike price.
		Right       string  // Either Put or Call (i.e. Options). Valid values are P, PUT, C, CALL.
		Multiplier  string  // The instrument's multiplier (i.e. options, futures).
		Exchange    string  // The destination exchange.
		Currency    string  // The underlying's currency.
		LocalSymbol string  // The contract's symbol within its primary exchange For options, this will be the OCC symbol.

		// The contract's primary exchange.
		// For smart routed contracts, used to define contract in case of ambiguity.
		// Should be defined as native exchange of contract, e.g. ISLAND for MSFT For exchanges which contain a period in name, will only be part of exchange name prior to period, i.e. ENEXT for ENEXT.BE.
		PrimaryExchange string

		TradingClass   string // The trading class name for this contract. Available in TWS contract description window as well. For example, GBL Dec '13 future's trading class is "FGBL".
		IncludeExpired bool   // If set to true, contract details requests and historical data queries can be performed pertaining to expired futures contracts. Expired options or other instrument types are not available.
		SecurityIdType string // Security's identifier when querying contract's details or placing orders ISIN - Example: Apple: US0378331005 CUSIP - Example: Apple: 037833100.
		SecurityId     string // Identifier of the security type.

		ComboLegsDescription string // Description of the combo legs.
		ComboLegs            []ComboLeg
		DeltaNeutralContract DeltaNeutralContract // elta and underlying price for Delta-Neutral combo orders. Underlying (STK or FUT), delta and underlying price goes into this attribute.
	}

	ComboLeg struct {
		ContractId int    // The Contract's IB's unique id.
		Ratio      int    // Select the relative number of contracts for the leg you are constructing. To help determine the ratio for a specific combination order, refer to the Interactive Analytics section of the User's Guide.
		Action     string //The side (buy or sell) of the leg:
		Exchange   string // The destination exchange to which the order will be routed.

		// Specifies whether an order is an open or closing order. For instituational customers to determine if this order is to open or close a position. 0 - Same as the parent security. This is the only option for retail customers.
		// 1 - Open. This value is only valid for institutional customers.
		// 2 - Close. This value is only valid for institutional customers.
		// 3 - Unknown.
		OpenClose int

		ShortSaleSlot      int    // For stock legs when doing short selling. Set to 1 = clearing broker, 2 = third party.
		DesignatedLocation string // When ShortSaleSlot is 2, this field shall contain the designated location.
		ExemptCode         int    // DOC_TODO.
	}

	// elta and underlying price for Delta-Neutral combo orders. Underlying (STK or FUT), delta and underlying price goes into this attribute.
	DeltaNeutralContract struct {
		ContractId string  // The unique contract identifier specifying the security. Used for Delta-Neutral Combo contracts.
		Delta      float64 // The underlying stock or future delta. Used for Delta-Neutral Combo contracts.
		Price      float64 // The price of the underlying. Used for Delta-Neutral Combo contracts.
	}

	Bar struct {
		Time   time.Time // The bar's date and time (either as a yyyymmss hh:mm:ss formatted string or as system time according to the request). Time zone is the TWS time zone chosen on login.
		Open   float64   // The bar's open price.
		High   float64   // The bar's high price.
		Low    float64   // The bar's low price.
		Close  float64   // The bar's close price.
		Volume int64     // The bar's traded volume if available (only available for TRADES)
		WAP    float64   // The bar's Weighted Average Price (only available for TRADES)
		Count  int       // The number of trades during the bar's timespan (only available for TRADES)
	}
)
