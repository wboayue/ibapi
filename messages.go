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
	tickByTick                           = 99
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
	requestMarketData                  = 1
	cancelMarketData                   = 2
	placeOrder                         = 3
	cancelOrder                        = 4
	requestOpenOrders                  = 5
	requestAccountData                 = 6
	requestExecutions                  = 7
	requestIds                         = 8
	requestContractData                = 9
	requestMarketDepth                 = 10
	cancelMarketDepth                  = 11
	requestNewsBulletins               = 12
	cancelNewsBulletins                = 13
	setServerLoglevel                  = 14
	requestAutoOpenOrders              = 15
	requestAllOpenOrders               = 16
	requestManagedAccounts             = 17
	requestFa                          = 18
	replaceFa                          = 19
	requestHistoricalData              = 20
	exerciseOptions                    = 21
	requestScannerSubscription         = 22
	cancelScannerSubscription          = 23
	requestScannerParameters           = 24
	cancelHistoricalData               = 25
	requestCurrentTime                 = 49
	requestRealTimeBars                = 50
	cancelRealTimeBars                 = 51
	requestFundamentalData             = 52
	cancelFundamentalData              = 53
	requestCalculateImpliedVolatility  = 54
	requestCalculateOptionPrice        = 55
	cancelCalculateImpliedVolatility   = 56
	cancelCalculateOptionPrice         = 57
	requestGlobalCancel                = 58
	requestMarketDataType              = 59
	requestPositions                   = 61
	requestAccountSummary              = 62
	cancelAccountSummary               = 63
	cancelPositions                    = 64
	verifyRequest                      = 65
	veridyMessage                      = 66
	queryDisplayGroups                 = 67
	subscribeToGroupEvents             = 68
	updateDisplayGroup                 = 69
	unsubscribeFromGroupEvents         = 70
	startApi                           = 71
	verifyAndAuthRequest               = 72
	verifyAndAuthMessage               = 73
	requestPositionsMulti              = 74
	cancelPositionsMulti               = 75
	requestAccountUpdatesMulti         = 76
	cancelAccountUpdatesMulti          = 77
	requestSecurityDefinitionOptParams = 78
	requestSoftDollarTiers             = 79
	requestFamilyCodes                 = 80
	requestMatchingSymbols             = 81
	requestMarketDepthExchanges        = 82
	requestSmartComponents             = 83
	requestNewsArticle                 = 84
	requestNewsProviders               = 85
	requestHistoricalNews              = 86
	requestHeadTimestamp               = 87
	requestHistogramData               = 88
	cancelHistogramData                = 89
	cancelHeadTimestamp                = 90
	requestMarketRule                  = 91
	requestPnl                         = 92
	cancelPnl                          = 93
	requestPnlSingle                   = 94
	cancelPnlSingle                    = 95
	requestHistoricalTicks             = 96
	requestTickByTickData              = 97
	cancelTickByTickData               = 98
	requestCompleteOrders              = 99
	requestWshMetaData                 = 100
	cancelWshMetaData                  = 101
	requestWshEventData                = 102
	cancelWashEventData                = 103
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
