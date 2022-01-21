package ibapi

import (
	"fmt"
	"strconv"
	"strings"
)

// Codes for incoming messages.
const (
	endConn                              = 0
	tickPrice                            = 1
	tickSize                             = 2
	orderStatus                          = 3
	errMsg                               = 4
	openOrder                            = 5
	accountValue                         = 6
	portfolioValue                       = 7
	accountUpdateTime                    = 8
	nextValidId                          = 9
	contractData                         = 10
	executionData                        = 11
	marketDepth                          = 12
	marketDepthL2                        = 13
	newsBulletins                        = 14
	managedAccounts                      = 15
	receiveFa                            = 16
	historicalData                       = 17
	bondContractData                     = 18
	scannerParameters                    = 19
	scannerData                          = 20
	tickOptionComputation                = 21
	tickGeneric                          = 45
	tickString                           = 46
	tickEfp                              = 47
	currentTime                          = 49
	realTimeBars                         = 50
	fundamentalData                      = 51
	contractDataEnd                      = 52
	openOrderEnd                         = 53
	accountDownloadEnd                   = 54
	executionDataEnd                     = 55
	deltaNeutralValidation               = 56
	tickSnapshotEnd                      = 57
	marketDataType                       = 58
	commissionReport                     = 59
	positionData                         = 61
	positionEnd                          = 62
	accountSummary                       = 63
	accountSummaryEnd                    = 64
	verifyMessageApi                     = 65
	verifyCompleted                      = 66
	displayGroupList                     = 67
	displayGroupUpdated                  = 68
	verifyAndAuthMessageApi              = 69
	verifyAndAuthCompleted               = 70
	positionMulti                        = 71
	positionMultiEnd                     = 72
	accountUpdateMulti                   = 73
	accountUpdateMultiEnd                = 74
	securityDefinitionOptionParameter    = 75
	securityDefinitionOptionParameterEnd = 76
	softDollarTiers                      = 77
	familyCode                           = 78
	symbolSample                         = 79
	marketDepthExchanges                 = 80
	tickRequestParameters                = 81
	smartComponents                      = 82
	newsArticls                          = 83
	tickNews                             = 84
	newsProviders                        = 85
	historicalNews                       = 86
	historicalNewsEnd                    = 87
	headTimestamp                        = 88
	histogramData                        = 89
	historicalDataUpdate                 = 90
	rerouteMarketDataRequest             = 91
	rerouteMarketDepthRequest            = 92
	markeRule                            = 93
	pnl                                  = 94
	pnlSingle                            = 95
	historicalTicks                      = 96
	historicalTicksBidAsk                = 97
	historicalTicksLast                  = 98
	TickByTick                           = 99
	orderBound                           = 100
	completedOrder                       = 101
	completedOrdersEnd                   = 102
	replaceFaEnd                         = 103
	wshMetaData                          = 104
	wshEventData                         = 105
	historicalSchedule                   = 106
)

// Codes for outgoing messages.
const (
	REQ_MKT_DATA                  = 1
	CANCEL_MKT_DATA               = 2
	PLACE_ORDER                   = 3
	CANCEL_ORDER                  = 4
	REQ_OPEN_ORDERS               = 5
	REQ_ACCT_DATA                 = 6
	REQ_EXECUTIONS                = 7
	REQ_IDS                       = 8
	REQ_CONTRACT_DATA             = 9
	REQ_MKT_DEPTH                 = 10
	CANCEL_MKT_DEPTH              = 11
	REQ_NEWS_BULLETINS            = 12
	CANCEL_NEWS_BULLETINS         = 13
	SET_SERVER_LOGLEVEL           = 14
	REQ_AUTO_OPEN_ORDERS          = 15
	REQ_ALL_OPEN_ORDERS           = 16
	REQ_MANAGED_ACCTS             = 17
	REQ_FA                        = 18
	REPLACE_FA                    = 19
	REQ_HISTORICAL_DATA           = 20
	EXERCISE_OPTIONS              = 21
	REQ_SCANNER_SUBSCRIPTION      = 22
	CANCEL_SCANNER_SUBSCRIPTION   = 23
	REQ_SCANNER_PARAMETERS        = 24
	CANCEL_HISTORICAL_DATA        = 25
	REQ_CURRENT_TIME              = 49
	REQ_REAL_TIME_BARS            = 50
	CancelRealTimeBars            = 51
	REQ_FUNDAMENTAL_DATA          = 52
	CANCEL_FUNDAMENTAL_DATA       = 53
	REQ_CALC_IMPLIED_VOLAT        = 54
	REQ_CALC_OPTION_PRICE         = 55
	CANCEL_CALC_IMPLIED_VOLAT     = 56
	CANCEL_CALC_OPTION_PRICE      = 57
	REQ_GLOBAL_CANCEL             = 58
	REQ_MARKET_DATA_TYPE          = 59
	REQ_POSITIONS                 = 61
	REQ_ACCOUNT_SUMMARY           = 62
	CANCEL_ACCOUNT_SUMMARY        = 63
	CANCEL_POSITIONS              = 64
	VERIFY_REQUEST                = 65
	VERIFY_MESSAGE                = 66
	QUERY_DISPLAY_GROUPS          = 67
	SUBSCRIBE_TO_GROUP_EVENTS     = 68
	UPDATE_DISPLAY_GROUP          = 69
	UNSUBSCRIBE_FROM_GROUP_EVENTS = 70
	StartApi                      = 71
	VERIFY_AND_AUTH_REQUEST       = 72
	VERIFY_AND_AUTH_MESSAGE       = 73
	REQ_POSITIONS_MULTI           = 74
	CANCEL_POSITIONS_MULTI        = 75
	REQ_ACCOUNT_UPDATES_MULTI     = 76
	CANCEL_ACCOUNT_UPDATES_MULTI  = 77
	REQ_SEC_DEF_OPT_PARAMS        = 78
	REQ_SOFT_DOLLAR_TIERS         = 79
	REQ_FAMILY_CODES              = 80
	REQ_MATCHING_SYMBOLS          = 81
	REQ_MKT_DEPTH_EXCHANGES       = 82
	REQ_SMART_COMPONENTS          = 83
	REQ_NEWS_ARTICLE              = 84
	REQ_NEWS_PROVIDERS            = 85
	REQ_HISTORICAL_NEWS           = 86
	REQ_HEAD_TIMESTAMP            = 87
	REQ_HISTOGRAM_DATA            = 88
	CANCEL_HISTOGRAM_DATA         = 89
	CANCEL_HEAD_TIMESTAMP         = 90
	REQ_MARKET_RULE               = 91
	REQ_PNL                       = 92
	CANCEL_PNL                    = 93
	REQ_PNL_SINGLE                = 94
	CANCEL_PNL_SINGLE             = 95
	REQ_HISTORICAL_TICKS          = 96
	REQ_TICK_BY_TICK_DATA         = 97
	CANCEL_TICK_BY_TICK_DATA      = 98
	REQ_COMPLETED_ORDERS          = 99
	REQ_WSH_META_DATA             = 100
	CANCEL_WSH_META_DATA          = 101
	REQ_WSH_EVENT_DATA            = 102
	CANCEL_WSH_EVENT_DATA         = 103
)

type messageBuilder struct {
	builder strings.Builder
}

func (b *messageBuilder) addInt(i int) {
	fmt.Fprintf(&b.builder, "%d\x00", i)
}

func (b *messageBuilder) addString(s string) {
	fmt.Fprintf(&b.builder, "%s\x00", s)
}

func (b *messageBuilder) addFloat32(num float32) {
	fmt.Fprintf(&b.builder, "%f\x00", num)
}

func (b *messageBuilder) addFloat64(num float64) {
	fmt.Fprintf(&b.builder, "%f\x00", num)
}

func (b *messageBuilder) addBool(flag bool) {
	if flag {
		fmt.Fprintf(&b.builder, "1\x00")
	} else {
		fmt.Fprintf(&b.builder, "0\x00")
	}
}

func (b *messageBuilder) Encode() string {
	return b.builder.String()
}

type parser struct {
	fields []string
}

func (s *parser) readInt() int {
	result := s.fields[0]
	s.fields = s.fields[1:]

	if result == "" {
		return 0
	}

	num, err := strconv.Atoi(result)
	if err != nil {
		panic(err)
	}
	return num
}

func (s *parser) readInt64() int64 {
	return int64(s.readInt())
}

func (s *parser) readFloat64() float64 {
	result := s.fields[0]
	s.fields = s.fields[1:]

	if result == "" {
		return 0
	}

	num, err := strconv.ParseFloat(result, 64)
	if err != nil {
		panic(err)
	}
	return num
}

func (s *parser) readString() string {
	result := s.fields[0]
	s.fields = s.fields[1:]
	return result
}
