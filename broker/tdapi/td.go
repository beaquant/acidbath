/*
   AcidBath - framework for your trading
   Copyright (C) 2016 Mark Laczynski

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

//Package tdapi supports an interface to the TDAmeritrade API.
package tdapi

// All functions/methods in this file are intended to be reusable, even if by internal private functions/methods

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"log"

	"github.com/marklaczynski/acidbath/broker/tdapi/internal/amtd"
	"github.com/marklaczynski/acidbath/broker/tdapi/tdstream"
	"github.com/marklaczynski/acidbath/broker/tdapi/tdstream/acctactivityfield"
	"github.com/marklaczynski/acidbath/broker/tdapi/tdstream/optrequestfield"
	"github.com/marklaczynski/acidbath/broker/tdapi/tdstream/quoterequestfield"
	"github.com/marklaczynski/acidbath/dm/asset"
	"github.com/marklaczynski/acidbath/dm/optionchain"
	"github.com/marklaczynski/acidbath/dm/optionchain/option"
	"github.com/marklaczynski/acidbath/dm/order"
	"github.com/marklaczynski/acidbath/dm/orderbook"
	"github.com/marklaczynski/acidbath/dm/ordermessage"
	"github.com/marklaczynski/acidbath/dm/orderstatus"
	"github.com/marklaczynski/acidbath/dm/portfolio"
	"github.com/marklaczynski/acidbath/dm/watchlists"
	eventFactory "github.com/marklaczynski/acidbath/eventproc/factory"
	genericEvent "github.com/marklaczynski/acidbath/eventproc/generic"
	"github.com/marklaczynski/acidbath/lib/date"
	"github.com/marklaczynski/acidbath/lib/financial"
	"github.com/marklaczynski/acidbath/lib/mjlog"
	"github.com/marklaczynski/acidbath/lib/orderconst"
)

type operation int

const httpContentType = "application/x-www-form-urlencoded"
const (
	opLogin operation = iota
	opLogout
	opPortfolio
	opSnapshot
	opOptionChain
	opStreamerInfo
	opOptionTrade
	opCancelOrder
	opMessageKey
	opOrderStatus
	opGetWatchlists
	opImpVolHistory
	opPriceHistory
)

type tdSession struct {
	amtdLogin        *amtd.Login
	amtdLogout       *amtd.Logout
	amtdPortfolio    *amtd.Portfolio
	amtdSnapshot     *amtd.SnapshotQuotes
	amtdOptionChain  *amtd.OptionChain
	amtdStreamerInfo *amtd.StreamerInfo
	amtdOrder        *amtd.Order
	amtdMessageKey   *amtd.MessageKey
	amtdCancelOrder  *amtd.CancelOrderMessage
	amtdOrderStatus  *amtd.OrderStatus
	amtdWatchlists   *amtd.Watchlists
}

//Session represents a TD Ameritrade session, which keeps track of various internal state information and caches.
//Current Spec: User can only stream 1 optionchain at a time, and stream the *tracked* options (from various underlyings)
type Session struct {
	sync.RWMutex //RW mutex on the whole structure for next 3 sets of data. Any API call will need to lock and unlock at this level

	// td ameritrade internal structures
	tdSession

	//session info set externally
	sourceID string
	version  string

	//streaming components
	streamingInProgress bool
	streamingBody       io.ReadCloser
	streamingCookies    []*http.Cookie

	strats [eventFactory.Count]genericEvent.Strategy

	// channels
	orderStatusDone chan bool
	endSession      chan bool

	optChanMutex         sync.RWMutex
	optionUpdateChans    map[string]chan *option.Option
	balChanMutex         sync.RWMutex
	portfolioUpdateChans map[string]chan *portfolio.Portfolio
	ordChanMutex         sync.RWMutex
	orderUpdateChans     map[string]chan *ordermessage.Message
}

var (
	logInfo  = log.New(mjlog.CreateInfoFile(), "INFO  [tdapi]: ", log.LstdFlags|log.Lshortfile)
	logDebug = log.New(mjlog.CreateDebugFile(), "DEBUG [tdapi]: ", log.LstdFlags|log.Lshortfile)
	logError = log.New(mjlog.CreateErrorFile(), "ERROR [tdapi]: ", log.LstdFlags|log.Lshortfile)
)

//New returns a pointer to the a new broker session
func New() *Session {
	s := &Session{
		orderStatusDone:      make(chan bool),
		endSession:           make(chan bool),
		optionUpdateChans:    make(map[string]chan *option.Option),
		portfolioUpdateChans: make(map[string]chan *portfolio.Portfolio),
		orderUpdateChans:     make(map[string]chan *ordermessage.Message),
	}

	// initialize all the strategies
	for i := 0; i < int(eventFactory.Count); i++ {
		s.strats[i] = eventFactory.CreateStrategy(eventFactory.StrategyType(i))
	}

	return s
}

func opURL(op operation, sourceid string, version string, param ...string) string {
	switch op {
	case opLogin:
		return "https://apis.tdameritrade.com/apps/300/LogIn?source=" + sourceid + "&version=" + version
	case opLogout:
		return "https://apis.tdameritrade.com/apps/100/LogOut?source=" + sourceid
	case opPortfolio:
		return "https://apis.tdameritrade.com/apps/100/BalancesAndPositions?source=" + sourceid
	case opSnapshot:
		syms := strings.Join(param, ",")
		return "https://apis.tdameritrade.com/apps/100/Quote?source=" + sourceid + "&symbol=" + syms
	case opOptionChain:
		return "https://apis.tdameritrade.com/apps/200/OptionChain?source=" + sourceid + "&symbol=" + param[0]
	case opStreamerInfo:
		return "https://apis.tdameritrade.com/apps/100/StreamerInfo?source=" + sourceid
	case opOptionTrade:
		return "https://apis.tdameritrade.com/apps/100/OptionTrade?source=" + sourceid + "&orderstring=" + param[0]
	case opCancelOrder:
		return "https://apis.tdameritrade.com/apps/100/OrderCancel?source=" + sourceid + "&orderid=" + strings.Join(param, "&orderid=") // + "<#order-id#>&orderid=<#order-id#>"
	case opMessageKey:
		return "https://apis.tdameritrade.com/apps/100/MessageKey?source=" + sourceid + "&accountid=" + param[0]
	case opOrderStatus:
		return "https://apis.tdameritrade.com/apps/100/OrderStatus?source=" + sourceid
	case opGetWatchlists:
		return "https://apis.tdameritrade.com/apps/100/GetWatchlists?source=" + sourceid
	case opImpVolHistory:
		return "https://apis.tdameritrade.com/apps/100/VolatilityHistory?source=" + sourceid +
			"&requestidentifiertype=" + param[0] +
			"&requestvalue=" + param[1] +
			"&volatilityhistorytype=" + param[2] +
			"&intervaltype=" + param[3] +
			"&intervalduration=" + param[4] +
			"&periodtype=" + param[5] +
			"&period=" + param[6] +
			"&startdate=" + param[7] +
			"&enddate=" + param[8] +
			"&daystoexpiration=" + param[9] +
			"&surfacetypeidentifier=" + param[10] +
			"&surfacetypevalue=" + param[11]
	case opPriceHistory:
		return "https://apis.tdameritrade.com/apps/100/PriceHistory?source=" + sourceid +
			"&requestidentifiertype=" + param[0] +
			"&requestvalue=" + param[1] +
			"&intervaltype=" + param[2] +
			"&intervalduration=" + param[3] +
			"&periodtype=" + param[4] +
			"&period=" + param[5] +
			"&startdate=" + param[6] +
			"&enddate=" + param[7] +
			"&extended=" + param[8]
	}
	return ""
}

func (s *Session) RegisterOptionUpdateChan(id string) chan *option.Option {
	s.optChanMutex.Lock()
	s.optionUpdateChans[id] = make(chan *option.Option)
	s.optChanMutex.Unlock()
	return s.optionUpdateChans[id]
}

func (s *Session) DeregisterOptionUpdateChan(id string) {
	s.optChanMutex.Lock()
	close(s.optionUpdateChans[id])
	delete(s.optionUpdateChans, id)
	s.optChanMutex.Unlock()
}

func (s *Session) notifyOptionUpdate(o *option.Option) {
	s.optChanMutex.RLock()
	for _, v := range s.optionUpdateChans {
		v <- o.Copy()
	}
	s.optChanMutex.RUnlock()
}

func (s *Session) RegisterPortfolioUpdateChan(id string) chan *portfolio.Portfolio {
	s.balChanMutex.Lock()
	s.portfolioUpdateChans[id] = make(chan *portfolio.Portfolio)
	s.balChanMutex.Unlock()
	return s.portfolioUpdateChans[id]
}

func (s *Session) DeregisterPortfolioUpdateChan(id string) {
	s.balChanMutex.Lock()
	close(s.portfolioUpdateChans[id])
	delete(s.portfolioUpdateChans, id)
	s.balChanMutex.Unlock()
}

func (s *Session) notifyPortfolioUpdate(portfolioParam *portfolio.Portfolio) {
	s.balChanMutex.RLock()
	for k, v := range s.portfolioUpdateChans {
		logDebug.Printf("sending portfolio data to %s\n", k)
		logDebug.Printf("balances: %#v\n", portfolioParam.Balance())
		v <- portfolioParam.Copy()
	}
	s.balChanMutex.RUnlock()
}

func (s *Session) RegisterOrderUpdateChan(id string) chan *ordermessage.Message {
	s.ordChanMutex.Lock()
	s.orderUpdateChans[id] = make(chan *ordermessage.Message)
	s.ordChanMutex.Unlock()
	return s.orderUpdateChans[id]
}

func (s *Session) DeregisterOrderUpdateChan(id string) {
	s.ordChanMutex.Lock()
	close(s.orderUpdateChans[id])
	delete(s.orderUpdateChans, id)
	s.ordChanMutex.Unlock()
}

func (s *Session) notifyOrderUpdate(message *ordermessage.Message) {
	s.ordChanMutex.RLock()
	for _, v := range s.orderUpdateChans {
		v <- message.Copy()
	}
	s.ordChanMutex.RUnlock()
}

//Login logs the user into the broker
// Post condition: Sets source id for the session, which is used as a param for other method calls
func (s *Session) Login(loginid string, pass string) error {
	logInfo.Printf("Login\n")

	if s.amtdLogin != nil && s.amtdLogin.Result == "OK" {
		return fmt.Errorf("Already logged in")
	}

	s.Lock()
	defer s.Unlock()

	if s.sourceID == "" {
		file, err := os.Open(os.Getenv("GOPATH") + "/src/github.com/marklaczynski/acidbath/broker/tdapi/config/tdconfig.json")
		if err != nil {
			logError.Printf("internal error, unable to open config file\n")
			return errors.New("internal error, unable to open config file\n")
		}

		var config struct {
			SourceID string
			Version  string
		}

		if err = json.NewDecoder(file).Decode(&config); err != nil {
			logError.Printf("internal error, unable to decode config file\n")
			return errors.New("internal error, unable to decode config file\n")
		}

		s.sourceID = config.SourceID
		s.version = config.Version
	}

	loginParams := url.Values{"userid": {loginid}, "password": {pass}, "sourceID": {s.sourceID}, "version": {s.version}}

	err := postRequest(opURL(opLogin, s.sourceID, s.version), loginParams, &s.amtdLogin, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error logging in: %s\n", err)
		return fmt.Errorf("Error logging in: %s", err)
	}

	// We can now check the result
	if s.amtdLogin.Result == "FAIL" {
		logError.Printf("Login service returned failure. Error: %s Result: %s\n", s.amtdLogin.Error, s.amtdLogin.Result)
		return fmt.Errorf("Login service returned failure. Result: %s", s.amtdLogin.Error)
	}

	if err = s.streamAccountActivity(); err != nil {
		logError.Printf("Login service could not start account streaming Error: %s", err)
		return fmt.Errorf("Login service could not start account streaming Error: %s", err)
	}

	return nil
}

//Logout logs the user out from the broker
func (s *Session) Logout() error {
	logInfo.Printf("Logout\n")

	s.Lock()
	defer s.Unlock()

	if s.amtdLogin != nil && s.amtdLogin.Result != "OK" {
		return fmt.Errorf("Not logged in")
	}

	err := postRequest(opURL(opLogout, s.sourceID, ""), nil, &s.amtdLogout, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error logging out: %s\n", err)
		return fmt.Errorf("Error logging out: %s", err)
	}

	if s.amtdLogout.Result != "LoggedOut" {
		return errors.New("Logout service returned failure")
	}

	// at logout invalidate the amtdLogin struct
	s.amtdLogin = nil
	s.streamingInProgress = false
	// end the session in the background, as some clean up needs to occur.
	go func() { s.endSession <- true }()

	return nil
}

//RetrievePortfolio requests the broker for the Portfolio information
func (s *Session) RetrievePortfolio(portfolioParam *portfolio.Portfolio) error {
	logInfo.Printf("RetrievePortfolio\n")

	s.Lock()
	defer s.Unlock()
	//TODO: HIGH: nil out all the amtd* structures... and only keep the info i need elsewhere (ie keep account at ui level, and have it pass in that info in future)
	s.amtdPortfolio = nil

	portfolioParams := url.Values{"accountid": {}, "type": {}, "suppressquotes": {}, "altbalanceformat": {}}

	err := postRequest(opURL(opPortfolio, s.sourceID, ""), portfolioParams, &s.amtdPortfolio, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error calling balances and positions: %s\n", err)
		return fmt.Errorf("Error calling balances and positions: %s", err)
	}

	if s.amtdPortfolio.Result != "OK" {
		return errors.New("BnP service returned failure")
	}

	portfolioParam.Balance().SetNetLiquidity(s.netLiquidity())
	portfolioParam.Balance().SetOptionBuyingPower(s.optionBuyingPower())

	for _, stockPosition := range s.amtdPortfolio.Positions.Stocks.Position {
		tmpPos := portfolio.NewPosition()
		tmpPos.SetSymbol(stockPosition.Security.Symbol)
		tmpPos.SetQuantity(float64(stockPosition.Quantity))
		tmpPos.SetAssetType(mapAsset(stockPosition.Security.AssetType))
		tmpPos.SetCusip(stockPosition.Security.Cusip)
		tmpPos.SetAccountType(int64(stockPosition.AccountType))
		tmpPos.SetClosePrice(stockPosition.ClosePrice)
		tmpPos.SetPositionType(stockPosition.PositionType)
		tmpPos.SetAveragePrice(stockPosition.AveragePrice)
		tmpPos.SetCurrentValue(stockPosition.CurrentValue)
		tmpPos.SetUnderlyingSymbol(stockPosition.UnderlyingSymbol)
		tmpPos.SetPutCallIndicator(stockPosition.PutCall)
		tmpPos.SetUnderlyingStock(asset.NewStock(stockPosition.Security.Symbol))
		tmpPos.SetUnderlyingOption(option.NewNilOption())

		portfolioParam.AddPosition(asset.EquityType, tmpPos)
	}

	for _, optionPosition := range s.amtdPortfolio.Positions.Options.Position {
		tmpPos := portfolio.NewPosition()
		tmpPos.SetQuantity(float64(optionPosition.Quantity))
		tmpPos.SetSymbol(optionPosition.Security.Symbol)
		tmpPos.SetAssetType(mapAsset(optionPosition.Security.AssetType))
		tmpPos.SetCusip(optionPosition.Security.Cusip)
		tmpPos.SetAccountType(int64(optionPosition.AccountType))
		tmpPos.SetClosePrice(optionPosition.ClosePrice)
		tmpPos.SetPositionType(optionPosition.PositionType)
		tmpPos.SetAveragePrice(optionPosition.AveragePrice)
		tmpPos.SetCurrentValue(optionPosition.CurrentValue)
		tmpPos.SetUnderlyingSymbol(optionPosition.UnderlyingSymbol)
		tmpPos.SetPutCallIndicator(optionPosition.PutCall)
		tmpPos.SetUnderlyingStock(asset.NewStock(optionPosition.UnderlyingSymbol))
		tmpPos.SetUnderlyingOption(optionPosition.Quote.NewOption())

		portfolioParam.AddPosition(asset.OptionType, tmpPos)
	}

	for _, fundsPosition := range s.amtdPortfolio.Positions.Funds.Position {
		tmpPos := portfolio.NewPosition()
		tmpPos.SetQuantity(float64(fundsPosition.Quantity))
		tmpPos.SetSymbol(fundsPosition.Security.Symbol)
		tmpPos.SetAssetType(mapAsset(fundsPosition.Security.AssetType))
		tmpPos.SetCusip(fundsPosition.Security.Cusip)
		tmpPos.SetAccountType(int64(fundsPosition.AccountType))
		tmpPos.SetClosePrice(fundsPosition.ClosePrice)
		tmpPos.SetPositionType(fundsPosition.PositionType)
		tmpPos.SetAveragePrice(fundsPosition.AveragePrice)
		tmpPos.SetCurrentValue(fundsPosition.CurrentValue)
		tmpPos.SetUnderlyingSymbol(fundsPosition.UnderlyingSymbol)
		tmpPos.SetPutCallIndicator(fundsPosition.PutCall)
		tmpPos.SetUnderlyingStock(nil)
		tmpPos.SetUnderlyingOption(option.NewNilOption())

		portfolioParam.AddPosition(asset.MutualFundType, tmpPos)
	}

	for _, bondsPosition := range s.amtdPortfolio.Positions.Bonds.Position {
		tmpPos := portfolio.NewPosition()
		tmpPos.SetQuantity(float64(bondsPosition.Quantity))
		tmpPos.SetSymbol(bondsPosition.Security.Symbol)
		tmpPos.SetAssetType(mapAsset(bondsPosition.Security.AssetType))
		tmpPos.SetCusip(bondsPosition.Security.Cusip)
		tmpPos.SetAccountType(int64(bondsPosition.AccountType))
		tmpPos.SetClosePrice(bondsPosition.ClosePrice)
		tmpPos.SetPositionType(bondsPosition.PositionType)
		tmpPos.SetAveragePrice(bondsPosition.AveragePrice)
		tmpPos.SetCurrentValue(bondsPosition.CurrentValue)
		tmpPos.SetUnderlyingSymbol(bondsPosition.UnderlyingSymbol)
		tmpPos.SetPutCallIndicator(bondsPosition.PutCall)
		tmpPos.SetUnderlyingStock(nil)
		tmpPos.SetUnderlyingOption(option.NewNilOption())

		portfolioParam.AddPosition(asset.BondType, tmpPos)
	}

	for _, mmPosition := range s.amtdPortfolio.Positions.MoneyMarket.Position {
		tmpPos := portfolio.NewPosition()
		tmpPos.SetQuantity(float64(mmPosition.Quantity))
		tmpPos.SetSymbol(mmPosition.Security.Symbol)
		tmpPos.SetAssetType(mapAsset(mmPosition.Security.AssetType))
		tmpPos.SetCusip(mmPosition.Security.Cusip)
		tmpPos.SetAccountType(int64(mmPosition.AccountType))
		tmpPos.SetClosePrice(mmPosition.ClosePrice)
		tmpPos.SetPositionType(mmPosition.PositionType)
		tmpPos.SetAveragePrice(mmPosition.AveragePrice)
		tmpPos.SetCurrentValue(mmPosition.CurrentValue)
		tmpPos.SetUnderlyingSymbol(mmPosition.UnderlyingSymbol)
		tmpPos.SetPutCallIndicator(mmPosition.PutCall)
		tmpPos.SetUnderlyingStock(nil)
		tmpPos.SetUnderlyingOption(option.NewNilOption())

		portfolioParam.AddPosition(asset.MoneyMarketType, tmpPos)
	}

	for _, savingsPosition := range s.amtdPortfolio.Positions.Savings.Position {
		tmpPos := portfolio.NewPosition()
		tmpPos.SetQuantity(float64(savingsPosition.Quantity))
		tmpPos.SetSymbol(savingsPosition.Security.Symbol)
		tmpPos.SetAssetType(mapAsset(savingsPosition.Security.AssetType))
		tmpPos.SetCusip(savingsPosition.Security.Cusip)
		tmpPos.SetAccountType(int64(savingsPosition.AccountType))
		tmpPos.SetClosePrice(savingsPosition.ClosePrice)
		tmpPos.SetPositionType(savingsPosition.PositionType)
		tmpPos.SetAveragePrice(savingsPosition.AveragePrice)
		tmpPos.SetCurrentValue(savingsPosition.CurrentValue)
		tmpPos.SetUnderlyingSymbol(savingsPosition.UnderlyingSymbol)
		tmpPos.SetPutCallIndicator(savingsPosition.PutCall)
		tmpPos.SetUnderlyingStock(nil)
		tmpPos.SetUnderlyingOption(option.NewNilOption())

		portfolioParam.AddPosition(asset.SavingsType, tmpPos)
	}

	//debugMarshal(s.amtdPortfolio)
	go func() {
		logDebug.Printf("sending portfolio updates on chan: %#v\n", portfolioParam)
		s.notifyPortfolioUpdate(portfolioParam)
	}()

	return nil
}

// RetrieveSnapshot requests the broker for symbol, and is only a single value (for now... in future if needed, expand to mulitple)
func (s *Session) RetrieveSnapshot(symbol string, assetType asset.AssetType, security interface{}) error {
	logInfo.Printf("RetrieveSnapshot: %s\n", symbol)

	s.Lock()
	defer s.Unlock()

	const FirstResult = 0
	s.amtdSnapshot = nil

	snapshotParams := url.Values{}

	symbols := make([]string, 1, 1)
	symbols[0] = symbol

	err := postRequest(opURL(opSnapshot, s.sourceID, "", symbols...), snapshotParams, &s.amtdSnapshot, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error calling quote snapshot service: %s\n", err)
		return fmt.Errorf("Error calling quote snapshot service: %s", err)
	}

	//debugMarshal(s.amtdSnapshot)

	if s.amtdSnapshot.Result != "OK" {
		logError.Printf("Snapshot service returned failure")
		return errors.New("Snapshot service returned failure")
	}

	if s.amtdSnapshot.QuoteList.Quote[FirstResult].Error != "" {
		logError.Printf("%s returned error: %s", symbol, s.amtdSnapshot.QuoteList.Quote[FirstResult].Error)
		return fmt.Errorf("%s returned error: %s", symbol, s.amtdSnapshot.QuoteList.Quote[FirstResult].Error)
	}

	switch mapAsset(s.amtdSnapshot.QuoteList.Quote[FirstResult].AssetType) {
	case asset.EquityType:
		a := security.(*asset.Stock)

		a.SetDescription(s.amtdSnapshot.QuoteList.Quote[FirstResult].Description)
		a.SetBidPrice(s.amtdSnapshot.QuoteList.Quote[FirstResult].Bid)
		a.SetAskPrice(s.amtdSnapshot.QuoteList.Quote[FirstResult].Ask)
		a.SetLastTradePrice(s.amtdSnapshot.QuoteList.Quote[FirstResult].Last)
		//a.FiftyTwoWeekLow = s.amtdSnapshot.QuoteList.Quote[FirstResult].YearLow
		//a.FiftyTwoWeekHigh = s.amtdSnapshot.QuoteList.Quote[FirstResult].YearHigh
		//a.BidSize, _ = strconv.ParseInt(strings.Split(s.amtdSnapshot.QuoteList.Quote[FirstResult].BidAskSize, "X")[0], 10, 64)
		//a.AskSize, _ = strconv.ParseInt(strings.Split(s.amtdSnapshot.QuoteList.Quote[FirstResult].BidAskSize, "X")[1], 10, 64)

	case asset.OptionType:
		/*
			Note for future:
				function() returns some new instance
				funcation(a) works on some existing function

				func() {
					var a xyz
					funcation(a)
					return a
				}
		*/
		a := security.(*option.Option)
		s.amtdSnapshot.QuoteList.Quote[FirstResult].ProcessOption(a)

	default:
		panic("unsupported type sent to RetrieveSnapshot")
	}

	return nil
}

func (s *Session) RetrieveImpliedVolatilityHistory(stockSymbol string, stock *asset.Stock) error {
	logInfo.Printf("RetrieveImpliedVolatilityHistory %s\n", stockSymbol)

	//for some reason VIX freaks this call out, so i'm just returning successfully without doing anytihng for now
	if stockSymbol == "VIX" {
		return nil
	}

	s.Lock()
	defer s.Unlock()

	volHistoryParams := url.Values{}

	volHistoryUrlParamValues := make([]string, 12, 12)
	volHistoryUrlParamValues[0] = "SYMBOL"                // requestidentifiertype : must be SYMBOL
	volHistoryUrlParamValues[1] = stockSymbol             // requestvalue
	volHistoryUrlParamValues[2] = "I"                     // volatilityhistorytype : I or H - I=Implied (calculated)  H=Historical (actual)
	volHistoryUrlParamValues[3] = "DAILY"                 // intervaltype :  DAY  - intervaltype can be DAILY ONLY.  MONTH - The intervaltype can be DAILY, WEEKLY, MONTHLY YEAR - The intervaltype can be WEEKLY, MONTHLY YTD - The intervaltype can be WEEKLY, MONTHLY
	volHistoryUrlParamValues[4] = "1"                     // intervalduration :  Always set to 1
	volHistoryUrlParamValues[5] = "YEAR"                  // periodtype : DAY, WEEK, MONTH,YEAR,YTD
	volHistoryUrlParamValues[6] = "1"                     // period :
	volHistoryUrlParamValues[7] = ""                      // startdate
	volHistoryUrlParamValues[8] = ""                      // enddate
	volHistoryUrlParamValues[9] = ""                      // daystoexpiration
	volHistoryUrlParamValues[10] = "DELTA_WITH_COMPOSITE" // surfacetypeidentifier : DELTA, DELTA_WITH_COMPOSITE , SKEW
	volHistoryUrlParamValues[11] = "50,-50"               // surfacetypevalue :  1 integer if DELTA, 2 integers for composite or skew

	var tmpHistoricalImpVol asset.ImpliedVolatilityTypeSlice
	err := postRequest(opURL(opImpVolHistory, s.sourceID, "", volHistoryUrlParamValues...), volHistoryParams, &tmpHistoricalImpVol, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error calling iv history service: %s\n", err)
		return fmt.Errorf("Error calling iv history service: %s", err)
	}
	stock.SetHistoricalImpliedVol(&tmpHistoricalImpVol)

	return nil
}

func (s *Session) RetrievePriceHistory(stockSymbol string, stock *asset.Stock) error {
	logInfo.Printf("RetrievePriceHistory: %s\n", stockSymbol)

	s.Lock()
	defer s.Unlock()

	priceHistoryParams := url.Values{}

	priceHistoryUrlParamValues := make([]string, 9, 9)
	priceHistoryUrlParamValues[0] = "SYMBOL"    // requestidentifiertype : must be SYMBOL
	priceHistoryUrlParamValues[1] = stockSymbol // requestvalue
	priceHistoryUrlParamValues[2] = "DAILY"     // intervaltype :  DAY  - intervaltype can be DAILY ONLY.  MONTH - The intervaltype can be DAILY, WEEKLY, MONTHLY YEAR - The intervaltype can be WEEKLY, MONTHLY YTD - The intervaltype can be WEEKLY, MONTHLY
	priceHistoryUrlParamValues[3] = "1"         // intervalduration :  Always set to 1
	priceHistoryUrlParamValues[4] = "MONTH"     // periodtype : DAY, WEEK, MONTH,YEAR,YTD
	priceHistoryUrlParamValues[5] = "3"         // period : The number of periods for which the data is returned. For example, if periodtype=DAY and period=10, then the request is for 10 days of data
	priceHistoryUrlParamValues[6] = ""          // startdate
	priceHistoryUrlParamValues[7] = ""          // enddate
	priceHistoryUrlParamValues[8] = ""          // extended

	var tmpHistoricalPrices []asset.PriceHistoryType
	err := postRequest(opURL(opPriceHistory, s.sourceID, "", priceHistoryUrlParamValues...), priceHistoryParams, &tmpHistoricalPrices, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error calling price history service: %s\n", err)
		return fmt.Errorf("Error calling price history service: %s", err)
	}
	stock.SetHistoricalPrice(&tmpHistoricalPrices)

	return nil
}

//AddOptionToStrategy adds a strategy to an option
func (s *Session) AddOptionToStrategy(opt *option.Option, strategy eventFactory.StrategyType) ([]string, error) {
	logInfo.Printf("AddOptionToStrategy\n")
	s.Lock()
	defer s.Unlock()

	s.strats[strategy].AddOption(opt) //start tracking option
	logInfo.Printf("attached option to strategy, and tracking %s\n", opt)
	listOfTrackedInstruments := s.getTrackedOptions() //mutex is obtained by calling updateOption()

	return listOfTrackedInstruments, nil
}

//RemoveOptionFromStrategy sets the strategy back to null
func (s *Session) RemoveOptionFromStrategy(opt *option.Option, strategy eventFactory.StrategyType) ([]string, error) {
	logInfo.Printf("RemoveOptionFromStrategy\n")
	s.Lock()
	defer s.Unlock()

	s.strats[strategy].RemoveOption(opt)
	logInfo.Printf("removed option from strategy, and untracking %s\n", opt)
	listOfTrackedInstruments := s.getTrackedOptions() //mutex is obtained by calling updateOption()

	return listOfTrackedInstruments, nil
}

//RetrieveWatchlists retrieves a watchlist from td
func (s *Session) RetrieveWatchlists(wls *watchlists.Watchlists) error {
	logInfo.Printf("RetrieveWatchlists\n")

	s.Lock()
	defer s.Unlock()

	watchlistsParams := url.Values{"accountid": {strconv.Itoa(int(s.amtdLogin.Login.Accounts[0].AccountID))}, "listid": {}}

	//wip
	err := postRequest(opURL(opGetWatchlists, s.sourceID, ""), watchlistsParams, &s.amtdWatchlists, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error calling watchlists: %s\n", err)
		return fmt.Errorf("Error calling watchlists: %s", err)
	}

	if s.amtdWatchlists.Result != "OK" {
		return errors.New("Watchlists service returned failure")
	}

	//map amtd data to local data model
	for _, amtdWl := range s.amtdWatchlists.WatchlistResults.Watchlist {
		wls.AddWatchlist(amtdWl.Name, int64(amtdWl.ID))
		for _, amtdWatchedSymbol := range amtdWl.SymbolList.WatchedSymbols {
			if amtdWatchedSymbol.Security.Symbol != "VIX" {
				wls.AddWatchedSymbol(int64(amtdWl.ID), amtdWatchedSymbol.Security.Symbol, mapAsset(amtdWatchedSymbol.Security.AssetType))
			}
		}
	}

	debugMarshal(s.amtdWatchlists)

	return nil
}

func mapAsset(assetType string) asset.AssetType {
	switch assetType {
	case "E":
		return asset.EquityType
	case "F":
		return asset.MutualFundType
	case "I":
		return asset.IndexType
	case "O":
		return asset.OptionType
	case "B":
		return asset.BondType
	case "M":
		return asset.MoneyMarketType
	case "V":
		return asset.SavingsType
	default:
		panic(fmt.Sprintf("Invalid asset type: %s\n", assetType))
	}

}

//retrieveOptionChain will call TD to retrieve the Option Chain for a single symbol.
func (s *Session) retrieveOptionChain(stock *asset.Stock) error {
	logInfo.Printf("retrieveOptionChain\n")

	oc := optionchain.NewOptionChain(stock.Symbol())

	// clear out any existing data if it already exists
	if s.amtdOptionChain != nil {
		s.amtdOptionChain = nil
	}

	// I need to use param range:O because I've seen all SPY options request cause a failure response
	optionChainParams := url.Values{"type": {}, "interval": {}, "strike": {},
		"expire": {"a"}, "range": {"O"}, "neardate": {}, "fardate": {}, "quotes": {"true"}}

	err := postRequest(opURL(opOptionChain, s.sourceID, "", []string{oc.Underlying()}...), optionChainParams, &s.amtdOptionChain, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error calling Option Chain service: %s\n", err)
		return fmt.Errorf("Error calling Option Chain service: %s", err)
	}

	//debugMarshal(s.amtdOptionChain)

	if s.amtdOptionChain.Result != "OK" {
		logError.Printf("Option Chain service returned failure. Error: %s Result: %s\n", s.amtdOptionChain.Error, s.amtdOptionChain.Result)
		return fmt.Errorf("Option Chain service returned failure. Result: %s", s.amtdOptionChain.Error)
	}

	for _, optDate := range s.amtdOptionChain.OptionChainResults.OptionDate {

		exp, _ := time.Parse(date.ParseOptionExpDate, optDate.Date)

		if err != nil {
			logError.Printf("Option Chain service had internal failure creating an option. Error: %s\n", err)
			return fmt.Errorf("Option Chain service had internal failure creating an option. Error: %s", err)
		}

		for _, optStrike := range optDate.OptionStrike {
			if optStrike.StandardOption == false {
				continue
			}

			if optStrike.Call != nil {

				o, err := oc.NewOption(oc.Underlying(), float64(optStrike.StrikePrice), exp, option.CALL, float64(optStrike.Call.Multiplier))
				if err != nil {
					logError.Printf("Error creating option\n")
					return errors.New("Error creating option\n")
				}

				o.SetOptionTickerSymbol(optStrike.Call.OptionSymbol)
				o.SetStrike(float64(optStrike.StrikePrice))
				o.SetExpirationDate(exp)
				o.SetDaysToExpiration(int64(optDate.DaysToExpiration))
				o.SetMultiplier(float64(optStrike.Call.Multiplier))
				o.SetLast(optStrike.Call.Last)
				o.SetBid(optStrike.Call.Bid)
				o.SetAsk(optStrike.Call.Ask)
				o.SetDelta(float64(optStrike.Call.Delta))
				o.SetGamma(float64(optStrike.Call.Gamma))
				o.SetTheta(float64(optStrike.Call.Theta))
				o.SetVega(float64(optStrike.Call.Vega))

				if o.Error() != nil {
					logError.Printf("Error constructing option: %s\n", o.Error())
					return fmt.Errorf("Error constructing option: %s\n", o.Error())

				}

				err = oc.AddOption(o)
				if err != nil {
					logError.Printf("Option Chain service had internal failure adding option to option chain. Error: %s\n", err)
					return fmt.Errorf("Option Chain service had internal failure adding option to option chain. Error: %s", err)
				}
			}

			if optStrike.Put != nil {
				o, err := oc.NewOption(oc.Underlying(), float64(optStrike.StrikePrice), exp, option.PUT, float64(optStrike.Put.Multiplier))
				if err != nil {
					logError.Printf("Error creating option\n")
					return errors.New("Error creating option\n")
				}

				o.SetOptionTickerSymbol(optStrike.Put.OptionSymbol)
				o.SetStrike(float64(optStrike.StrikePrice))
				o.SetExpirationDate(exp)
				o.SetDaysToExpiration(int64(optDate.DaysToExpiration))
				o.SetMultiplier(float64(optStrike.Put.Multiplier))
				o.SetLast(optStrike.Put.Last)
				o.SetBid(optStrike.Put.Bid)
				o.SetAsk(optStrike.Put.Ask)
				o.SetDelta(float64(optStrike.Put.Delta))
				o.SetGamma(float64(optStrike.Put.Gamma))
				o.SetTheta(float64(optStrike.Put.Theta))
				o.SetVega(float64(optStrike.Put.Vega))

				if o.Error() != nil {

					logError.Printf("Error constructing option: %s\n", o.Error())
					return fmt.Errorf("Error constructing option: %s\n", o.Error())
				}

				err = oc.AddOption(o)
				if err != nil {
					logError.Printf("Option Chain service had internal failure adding option to option chain. Error: %s\n", err)
					return fmt.Errorf("Option Chain service had internal failure adding option to option chain. Error: %s", err)
				}
			}
		}
	}

	stock.SetOptionChain(oc)

	logDebug.Printf("Exit retrieveOptionChain")
	return nil
}

//retrieveStreamerInfo should be called before any streaming requests, since it cannot be guarnteeed when a
//streaming session will start because it's a user request. It is re-entrant, which will get streamer info the
//first time it's called, and return true on subsequent calls assuming it succeeded the first time.
//Maybe in the future, I'll redesign in order to make this start up right after login
func (s *Session) retrieveStreamerInfo() error {
	logInfo.Printf("StreamerInfo\n")

	// "Singleton" function, basically we only need to pull the request successfully the first time, all subsequent calls return here
	if s.amtdStreamerInfo != nil && s.amtdStreamerInfo.Result == "OK" {
		logInfo.Printf("Streamer Info already set\n")
		return nil
	}

	streamerInfoParams := url.Values{"accountid": {}}

	err := postRequest(opURL(opStreamerInfo, s.sourceID, ""), streamerInfoParams, &s.amtdStreamerInfo, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error calling StreamerInfo service : %s\n", err)
		return fmt.Errorf("Error calling StreamerInfo service: %s", err)
	}

	//debugMarshal(s.amtdStreamerInfo)

	if s.amtdStreamerInfo.Result != "OK" {
		logError.Printf("StreamerInfo service returned failure. Error: %s Result: %s ErrorMsg: %s\n", s.amtdStreamerInfo.Error, s.amtdStreamerInfo.Result, s.amtdStreamerInfo.StreamerInfo.ErrorMsg)
		return fmt.Errorf("StreamerInfo service returned failure. Result: %s", s.amtdStreamerInfo.Error)
	}

	return nil
}

func (s *Session) retrieveMessageKey() error {
	logInfo.Printf("retrieveMessageKey\n")

	messageKeyParams := url.Values{}

	err := postRequest(opURL(opMessageKey, s.sourceID, "", strconv.Itoa(int(s.amtdLogin.Login.Accounts[0].AccountID))), messageKeyParams, &s.amtdMessageKey, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error calling message key service : %s\n", err)
		return fmt.Errorf("Error calling message key service: %s", err)
	}

	//debugMarshal(s.amtdMessageKey)

	if s.amtdMessageKey.Result != "OK" {
		logError.Printf("Message Key service returned failure. Error: %s Result: %s\n", s.amtdStreamerInfo.Error, s.amtdStreamerInfo.Result)
		return fmt.Errorf("Message Key service returned failure. Result: %s", s.amtdStreamerInfo.Error)
	}

	return nil

}

//AddStockOptionsToStream streams the stockSymbol and it's corresponding options. Currently only supporting OTM option only, FUTURE to make that configurable
func (s *Session) AddStockOptionsToStream(stock *asset.Stock) error {
	logInfo.Printf("AddStockOptionsToStream %s\n", stock.Symbol())

	s.Lock()
	defer s.Unlock()

	// start streaming service
	err := s.retrieveStreamerInfo()
	if err != nil {
		logInfo.Printf("Calling retrieveStreamerInfo: %s\n", err)
		return fmt.Errorf("Calling retrieveStreamerInfo: %s\n", err)
	}

	err = s.streamAllOptionsForStock(stock)
	if err != nil {
		logInfo.Printf("Error streaming option for stock: %s\n", err)
		return fmt.Errorf("Error streaming option for stock: %s\n", err)
	}

	err = s.streamStock(stock.Symbol())
	if err != nil {
		logInfo.Printf("Error streaming stock: %s with error:%s\n", stock.Symbol(), err)
		return fmt.Errorf("Error streaming stock: %s with error: %s\n", stock.Symbol(), err)
	}

	return nil
}

//RemoveStockOptionsFromStream streams the stockSymbol and it's corresponding options. Currently only supporting OTM option only, FUTURE to make that configurable
func (s *Session) RemoveStockOptionsFromStream(stock *asset.Stock) error {
	logInfo.Printf("RemoveStockOptionsFromStream\n")

	s.Lock()
	defer s.Unlock()

	// start streaming service
	err := s.retrieveStreamerInfo()
	if err != nil {
		logInfo.Printf("Calling retrieveStreamerInfo: %s\n", err)
		return fmt.Errorf("Calling retrieveStreamerInfo: %s\n", err)
	}

	err = s.stream(stock.OptionChain().OptionSymbols(), tdstream.Quote, cmdUnsubs)
	if err != nil {
		logError.Printf("Error unsubscribing option from stream for %s", stock.Symbol())
		return fmt.Errorf("Error unsubscribing option from stream for %s", stock.Symbol())
	}

	err = s.stream([]string{stock.Symbol()}, tdstream.Option, cmdUnsubs)
	if err != nil {
		logError.Printf("Error unsubscribing quote from stream for %s", stock.Symbol())
		return fmt.Errorf("Error unsubscribing quote from stream for %s", stock.Symbol())
	}

	return nil
}

/*
streamAllOptionsForStock will start streaming all available OTM (FUTURE: make this user configurable, currently hardcoded) options

PostCondition:
After calling this activity the OptionChain field should be populated & updated assuming no errors
*/
func (s *Session) streamAllOptionsForStock(stock *asset.Stock) error {
	logInfo.Printf("streamAllOptionsForStock\n")

	// start streaming service
	err := s.retrieveStreamerInfo()
	if err != nil {
		logInfo.Printf("Error calling retrieveStreamerInfo: %s\n", err)
		return fmt.Errorf("Error calling retrieveStreamerInfo: %s\n", err)
	}

	err = s.retrieveOptionChain(stock)
	if err != nil {
		logError.Printf("Error retrieving Option chain %s\n", err)
		return fmt.Errorf("Error retrieving Option chain %s\n", err)
	}

	oc := stock.OptionChain()

	if s.streamingInProgress {
		s.stream(oc.OptionSymbols(), tdstream.Option, cmdAdd)
	} else {
		s.stream(oc.OptionSymbols(), tdstream.Option, cmdSubs)
	}

	return nil
}

func (s *Session) streamStock(stockSymbol string) error {
	logInfo.Printf("AddToStreamer\n")

	// start streaming service
	err := s.retrieveStreamerInfo()
	if err != nil {
		logInfo.Printf("Calling retrieveStreamerInfo: %s\n", err)
		return fmt.Errorf("Calling retrieveStreamerInfo: %s\n", err)
	}

	if s.streamingInProgress {
		err = s.stream([]string{stockSymbol}, tdstream.Quote, cmdAdd)
	} else {
		err = s.stream([]string{stockSymbol}, tdstream.Quote, cmdSubs)
	}

	return err
}

func (s *Session) streamAccountActivity() error {

	// start streaming service
	err := s.retrieveStreamerInfo()
	if err != nil {
		logInfo.Printf("Calling retrieveStreamerInfo: %s\n", err)
		return fmt.Errorf("Calling retrieveStreamerInfo: %s\n", err)
	}

	if err = s.stream([]string{}, tdstream.AcctActivity, cmdSubs); err != nil {
		return err
	}

	return nil
}

// stream subscribes to a list of symbols/tickers, for a given SID. OPTION & QUOTE are currently the only supported SIDs
// In the near future this function may take on more responsibility if I pass a streamingCommand parameter, but for now I don't need it
// This function is reentrant, because it creates only a ONE streaming go routine
// For now I think this is safe to say it can only be called once, since it's wrapped in a mutex in API calls, but in case the app
// becomes really dynamic, we need to ensure there is only 1 URL request at a time as per documenation:
// streamer_request_data.htm
func (s *Session) stream(tickerSymbols []string, sid tdstream.StreamingID, cmd streamingCommand) error {
	logInfo.Printf("stream\n")

	rawurl := "https://" + s.amtdStreamerInfo.StreamerInfo.StreamerURL + "/"

	//logDebug.Printf("streaming url whole: %s\n", rawurl)

	// Future... I should cache this client somehow instead of having to create it each time... low priority
	client, err := createClient(rawurl, s.streamingCookies, s.isLoggedIn(), s.sessionID())
	if err != nil {
		return err
	}

	postData := bytes.NewBufferString(s.streamRequest(sid, cmd, tickerSymbols))

	//TODO: Add retry logic
	var resp *http.Response
	for idx := 0; idx < 3; idx++ {
		resp, err = client.Post(rawurl, httpContentType, postData)
		if err != nil {
			logError.Printf("Error sending request: %s\n", err)
			//return fmt.Errorf("Error sending request: %s", err)
		} else {
			break
		}
	}

	if s.streamingInProgress == false {

		tmp := resp.Cookies()
		s.streamingCookies = append(s.streamingCookies, &http.Cookie{})
		s.streamingCookies[0].Name = tmp[0].Name
		s.streamingCookies[0].Value = tmp[0].Value

		s.streamingBody = resp.Body
		go streamParser(s)
	}

	s.streamingInProgress = true

	return nil
}

func (s *Session) streamRequest(sid tdstream.StreamingID, c streamingCommand, ulSymbols []string) string {
	var data, auth string

	// This is pretty standard for now...
	//FUTURE... populate some array list with all accounts, and use primary account index instead of "0"
	auth = "!U=" + strconv.Itoa(int(s.amtdLogin.Login.Accounts[0].AccountID)) +
		"&W=" + s.amtdStreamerInfo.StreamerInfo.Token +
		"&A=userid=" + strconv.Itoa(int(s.amtdLogin.Login.Accounts[0].AccountID)) +
		"&token=" + s.amtdStreamerInfo.StreamerInfo.Token +
		"&company=" + s.amtdLogin.Login.Accounts[0].Company +
		"&segment=" + s.amtdLogin.Login.Accounts[0].Segment +
		"&cddomain=" + s.amtdStreamerInfo.StreamerInfo.CDDomainID +
		"&usergroup=" + s.amtdStreamerInfo.StreamerInfo.Usergroup +
		"&accesslevel=" + s.amtdStreamerInfo.StreamerInfo.AccessLevel +
		"&authorized=" + s.amtdStreamerInfo.StreamerInfo.Authorized +
		"&acl=" + s.amtdStreamerInfo.StreamerInfo.Acl +
		"&timestamp=" + strconv.Itoa(int(s.amtdStreamerInfo.StreamerInfo.Timestamp)) +
		"&appid=" + s.amtdStreamerInfo.StreamerInfo.Appid

	logDebug.Printf("control flag: %s\n", s.controlFlag())

	symbolListing := ""
	fieldListing := ""

	switch c {
	case cmdUnsubsAll:
		//no params needed
	case cmdView:
		//cmdView currently unsupported... FUTURE
	case cmdUnsubs:
		symbolListing = "&P=" + strings.Join(ulSymbols, "+")
	case cmdSubs, cmdAdd:
		switch sid {
		case tdstream.Quote, tdstream.Option:

			symbolListing = "&P=" + strings.Join(ulSymbols, "+")
			if sid == tdstream.Option && len(s.getTrackedOptions()) > 0 {
				symbolListing += "+" + strings.Join(s.getTrackedOptions(), "+")
			}

		case tdstream.AcctActivity:
			s.retrieveMessageKey()
			symbolListing = "&P=" + s.amtdMessageKey.MessageKeyData.Token
		}

		fieldListing = "&T=" + fields(sid)
	default:
		panic("Developer error")
	}

	data = auth +
		"|source=" + s.sourceID +
		"|control=" + s.controlFlag() +
		"|S=" + fmt.Sprintf("%s", sid) +
		"&C=" + fmt.Sprintf("%s", c) +
		symbolListing +
		fieldListing

	data += "\n\n"

	logDebug.Printf("req string: :%s:\n", data)

	return data
}

/*
Control Flag (control=true or control=false)
The communication protocol for streaming data was updated to support the ability to update a streaming
subscription via secondary "control" requests to add/remove/change the subscribed data.

control parameter -- if "false", identifies a new streaming connection; if "true", identifies an update command
source parameter -- your application source ID

You start with a normal streaming connection, with "control=false" and "source" parameters. Then you can
make subscription changes by making subsequent control requests using a separate connection (without closing the
original streaming connection) with "control=true" and "source" parameters. This works for all streaming services.
*/
func (s *Session) controlFlag() string {
	if s.streamingInProgress {
		return "true"
	}
	return "false"
}

type streamingCommand int

const (
	cmdSubs streamingCommand = iota
	cmdAdd
	cmdUnsubs
	cmdUnsubsAll
	cmdView
)

func (cf streamingCommand) String() string {
	switch cf {
	case cmdSubs:
		return "SUBS"
	case cmdAdd:
		return "ADD"
	case cmdUnsubs:
		return "UNSUBS"
	case cmdUnsubsAll:
		return "UNSUBS"
	case cmdView:
		return "VIEW"
	}
	return ""
}

func fields(sid tdstream.StreamingID) string {
	switch sid {
	case tdstream.Quote:
		// for right now, i'm only going to stream symbol,bid,ask,(delta)... future is to make this configurable
		return fmt.Sprintf("%d+%d+%d", quoterequestfield.Symbol, quoterequestfield.Bid, quoterequestfield.Ask)
	case tdstream.TimeSale:
	case tdstream.Response:
	case tdstream.Option:
		return fmt.Sprintf("%d+%d+%d+%d+%d+%d+%d+%d", optrequestfield.Symbol,
			optrequestfield.Bid,
			optrequestfield.Ask,
			optrequestfield.Last,
			optrequestfield.DeltaIndex,
			optrequestfield.GammaIndex,
			optrequestfield.ThetaIndex,
			optrequestfield.VegaIndex)
	case tdstream.ActivesNYSE:
	case tdstream.ActivesNASDAQ:
	case tdstream.ActivesOTCBB:
	case tdstream.ActivesOptions:
	case tdstream.News:
	case tdstream.NewsHistory:
	case tdstream.AdapNASDAQ:
	case tdstream.NYSEBook:
	case tdstream.NYSEChart:
	case tdstream.NASDAQChart:
	case tdstream.OpraBook:
	case tdstream.IndexChart:
	case tdstream.TotalView:
	case tdstream.AcctActivity:
		return fmt.Sprintf("%d+%d+%d+%d", acctactivityfield.SubscriptionKey, acctactivityfield.AccountNumber, acctactivityfield.MessageType, acctactivityfield.MessageData)
	case tdstream.Chart:
	case tdstream.StreamerServer:
	}
	return ""
}

func createClient(rawurl string, cookiesParam []*http.Cookie, isLoggedIn bool, sessionID string) (*http.Client, error) {
	client := &http.Client{}

	if isLoggedIn {
		jar, err := cookiejar.New(nil)
		if err != nil {
			logError.Printf("Error setting up cookie jar")
			return nil, errors.New("Error setting up cookie jar")
		}

		c := http.Cookie{Name: "JSESSIONID", Value: sessionID}
		var cookies []*http.Cookie
		if cookiesParam != nil {
			cookies = append(cookies, cookiesParam[0])
		}
		cookies = append(cookies, &c)

		logInfo.Printf("sending cookie: %v\n", cookies)

		u, err := url.Parse(rawurl)
		if err != nil {
			logError.Printf("Error parsing URL: %s\n", rawurl)
			return nil, fmt.Errorf("Error parsing URL: %s", rawurl)
		}

		jar.SetCookies(u, cookies)
		client = &http.Client{Jar: jar}
	}

	return client, nil
}

func (s *Session) processOrderMessage(message *ordermessage.Message) {
	logInfo.Printf("processOrderMessage: %s\n", message.OrderID())

	s.notifyOrderUpdate(message)

	/* TBD
	for k := range s.strats {
		s.strats[k].OrderStatusUpdate(s.orderStatus[orderid].Copy())
	}
	*/
}

func (s *Session) updateOption(newOptionData *option.Option) {
	logInfo.Printf("updateOption")

	go func() {

		// 1 execute strategy on option (default is null, so nothing will happen)
		for k := range s.strats {
			if newOptionData != nil {
				s.strats[k].Execute(newOptionData.OptionTickerSymbol())
			}
		}

		// 2 fwd the option to channel for UI (and Execute() results will be there)
		if newOptionData != nil {
			/*
				The reason for Copy() is that newOptionData is a option.Option (interface), which holds a pointer
				to an option. If I had done something like
					oTmp := newOptionData
				then oTmp would have a copy of the pointer to the actual data, which still could be modified
				by some other go routine. This way I make a temp copy of the option, and send it to the ui
			*/
			s.notifyOptionUpdate(newOptionData)
		} else {
			for k := range s.strats {
				if opt, ok := s.strats[k].Option(newOptionData.OptionTickerSymbol()); ok {
					s.notifyOptionUpdate(opt)
				}
			}
		}

	}()

	if newOptionData == nil {

		logDebug.Printf("newOptionData is nil, looking for option in tracked options %#v\n", s.getTrackedOptions())
		for k := range s.strats {
			if opt, ok := s.strats[k].Option(newOptionData.OptionTickerSymbol()); ok {
				logDebug.Printf("tracked option %#v", opt)
				return
			}
		}

		logDebug.Printf("Did not find option in Option Chain or in trakedOptions, return nil ")
		return
	}

	return
}

func streamParser(s *Session) {
	defer s.streamingBody.Close()
	streamReader := tdstream.NewDecoder(s.streamingBody)

	sidHandler := &tdstream.SidHandlers{
		OptionCallback:          s.updateOption,
		AccountActivityCallback: s.processOrderMessage,
	}

	for {
		select {
		case <-s.endSession:
			logDebug.Printf("Ending the streamParser go routine\n")
			return

		default:
			switch h := streamReader.DecodeHeader(); h {
			case 'H':
				streamReader.DecodeHeartbeat()
			case 'N':
				streamReader.DecodeSnapshotResponse(sidHandler)
			case 'S':
				streamReader.DecodeCommonStreamingHeader(sidHandler)
			case 'X':
				logInfo.Printf("EOF reached")
				return
			case 'Y':
				// if streaming body is still running, then some error occured with the stream
				if s.streamingInProgress {
					logError.Printf("Invalid Header %x", h)
				} else {
					logInfo.Printf("Closing streamming connection\n")
				}
				return
			default:
				panic("Something unexpected occured")
			}
			// If I do not include this, the default case will starve resources
			runtime.Gosched()
		}
	}
}

func (s *Session) getTrackedOptions() []string {
	var result []string

	for k := range s.strats {
		result = append(result, s.strats[k].TrackedOptions()...)
	}

	return result
}

func (s *Session) sessionID() string {
	if s.amtdLogin != nil {
		return s.amtdLogin.Login.SessionID
	}
	return ""
}

func (s *Session) isLoggedIn() bool {
	if s.sessionID() == "" {
		return false
	}
	return true
}

func (s *Session) isBnPValid() bool {
	if s.amtdPortfolio != nil && s.amtdPortfolio.Result != "FAIL" {
		return true
	}
	return false
}

func (s *Session) netLiquidity() float64 {
	if s.isLoggedIn() && s.isBnPValid() {
		return float64(s.amtdPortfolio.Balance.AccountValue.Current)
	}
	return 0
}

func (s *Session) optionBuyingPower() float64 {
	if s.isLoggedIn() && s.isBnPValid() {
		return float64(s.amtdPortfolio.Balance.OptionBuyingPower)
	}
	return 0
}

func postRequest(rawurl string, postParams url.Values, v interface{}, isLoggedIn bool, sessionID string) error {
	logDebug.Printf("postRequest")
	client, err := createClient(rawurl, nil, isLoggedIn, sessionID)
	if err != nil {
		return err
	}

	//TODO: Add retry logic
	logDebug.Printf("posting request to broker\n")

	var resp *http.Response
	for idx := 0; idx < 3; idx++ {
		resp, err = client.PostForm(rawurl, postParams)
		if err != nil {
			logError.Printf("Error sending request: %s\n", err)
			return fmt.Errorf("Error sending request: %s", err)
		} else {
			break
		}
	}
	logDebug.Printf("received response from broker\n")
	defer resp.Body.Close()

	//logDebug.Printf("response header: %s\n", resp.Header)

	switch t := v.(type) {
	case *[]asset.PriceHistoryType:
		logDebug.Printf("pricehistory type parsing a %t", t)

		r := bufio.NewReader(resp.Body)

		symbolCount := tdstream.ReadInt32(r)
		logDebug.Printf("symbol count %d\n", symbolCount)

		var idx int32
		for idx = 0; idx < symbolCount; idx++ {
			symbol := tdstream.ReadString(r, int(tdstream.ReadInt16(r)))
			logDebug.Printf("parsed symbol: %s\n", symbol)
			errCode := tdstream.ReadInt8(r)
			logDebug.Printf("Error Code is %d", errCode)
			if errCode == 1 {
				errorMessage := tdstream.ReadString(r, int(tdstream.ReadInt16(r)))
				logError.Printf("Received an error parsing response: %s", errorMessage)
				return fmt.Errorf("Received an error parsing response")
			}
			numValues := tdstream.ReadInt32(r)

			var currValIdx int32
			var localHistoricalPrice *[]asset.PriceHistoryType
			localHistoricalPrice, ok := v.(*[]asset.PriceHistoryType)
			if !ok {
				logDebug.Printf("Panic caused by %#v\n", localHistoricalPrice)
				panic("Someone fucked up passing parameter v into postRequest")
			}

			for currValIdx = 0; currValIdx < numValues; currValIdx++ {
				_ = financial.Money{(&big.Rat{}).SetFloat64(float64(tdstream.ReadFloat32(r)))} // open
				_ = financial.Money{(&big.Rat{}).SetFloat64(float64(tdstream.ReadFloat32(r)))} // high
				_ = financial.Money{(&big.Rat{}).SetFloat64(float64(tdstream.ReadFloat32(r)))} //low
				closePrice := financial.Money{(&big.Rat{}).SetFloat64(float64(tdstream.ReadFloat32(r)))}
				_ = tdstream.ReadFloat32(r) // volume
				timeStamp := time.Unix(0, tdstream.ReadInt64(r)*int64(time.Millisecond))

				currHistoricalPrice := asset.NewPriceHistoryPoint(closePrice, timeStamp)

				logDebug.Printf("price data added: %s\n", currHistoricalPrice)
				*localHistoricalPrice = append(*localHistoricalPrice, currHistoricalPrice)
			}

			termCode := byte(tdstream.ReadInt8(r))
			if termCode != 0xFF {
				logError.Printf("Error with data terminator 1: %x\n", termCode)
			}
			termCode = byte(tdstream.ReadInt8(r))
			if termCode != 0xFF {
				logError.Printf("Error with data terminator 2: %x\n", termCode)
			}

		}

	case *asset.ImpliedVolatilityTypeSlice:
		logDebug.Printf("voldata type parsing a %t", t)

		r := bufio.NewReader(resp.Body)

		symbolCount := tdstream.ReadInt32(r)
		logDebug.Printf("symbol count %d\n", symbolCount)
		if symbolCount > 1 {
			logError.Printf("Currently only 1 symbol at a time is supported, and for some reason. %d is the number of symbols claimed by API call", symbolCount)
			return fmt.Errorf("Currently only 1 symbol at a time is supported, and for some reason. %d is the number of symbols claimed by API call", symbolCount)
		}

		var idx int32
		for idx = 0; idx < symbolCount; idx++ {
			symbol := tdstream.ReadString(r, int(tdstream.ReadInt16(r)))
			logDebug.Printf("parsed symbol: %s\n", symbol)
			errCode := tdstream.ReadInt8(r)
			logDebug.Printf("Error Code is %d", errCode)
			if errCode == 1 {
				errorMessage := tdstream.ReadString(r, int(tdstream.ReadInt16(r)))
				logError.Printf("Received an error parsing response: %s", errorMessage)
				return fmt.Errorf("Received an error parsing response")
			}
			numValues := tdstream.ReadInt32(r)

			var currValIdx int32
			var localVolData *asset.ImpliedVolatilityTypeSlice
			localVolData, ok := v.(*asset.ImpliedVolatilityTypeSlice)
			if !ok {
				logDebug.Printf("Panic caused by %#v\n", localVolData)
				panic("Someone fucked up passing parameter v into postRequest")
			}

			for currValIdx = 0; currValIdx < numValues; currValIdx++ {
				currVolData := asset.NewImpliedVolInstance(tdstream.ReadFloat32(r), time.Unix(0, tdstream.ReadInt64(r)*int64(time.Millisecond)))
				logDebug.Printf("vol data added: %s\n", currVolData)
				*localVolData = append(*localVolData, currVolData)
			}

			termCode := byte(tdstream.ReadInt8(r))
			if termCode != 0xFF {
				logError.Printf("Error with data terminator 1: %x\n", termCode)
			}
			termCode = byte(tdstream.ReadInt8(r))
			if termCode != 0xFF {
				logError.Printf("Error with data terminator 2: %x\n", termCode)
			}

		}

	default:
		logDebug.Printf("default parsing a %t", t)
		//logic switch
		if false {
			//remove ReadAll https://www.datadoghq.com/2014/07/crossing-streams-love-letter-gos-io-reader/
			//I'm actually keeping this code for reference, because I do believe I have a legit use for logging the
			//response body for debugging purposes. But once debugging is complete, I totally agree with the
			//piped approach in the article

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logError.Printf("Error reading response body: %s\n", err)
				return fmt.Errorf("Error reading response body: %s", err)
			}

			logDebug.Printf("xml: %s\n", body)
			err = xml.Unmarshal([]byte(body), v)
			if err != nil {
				logError.Printf("Error unmarshaling response: %s\n", err)
				return fmt.Errorf("Error unmarshaling response: %s", err)
			}
		} else {

			err = xml.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				logError.Printf("Error unmarshaling response: %s\n", err)
				return fmt.Errorf("Error unmarshaling response: %s", err)
			}
		}
	}

	return nil
}

// this is a general util can move to "common" place
func debugMarshal(data interface{}) {
	// these steps marshal indent the structure received, and prints it to console
	// This way I can compare this to what was received
	// start debug

	op, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		logDebug.Printf("Error: %v\n", err)
	}
	logDebug.Printf("\n%s\n", string(op))
}

//ErrOrderValidation represents an error when validation of order structure fails brokerage rules.
var ErrOrderValidation = errors.New("Validating order failed")

//SendSingleLegOptionTrade sends the order to TD. It always validates the order before sending request.
func (s *Session) SendSingleLegOptionTrade(order *order.Order) error {
	logInfo.Printf("SendOptionTrade\n")

	s.Lock()
	defer s.Unlock()

	tdo := &tdOrder{
		accountID: strconv.Itoa(int(s.amtdLogin.Login.Accounts[0].AccountID)),
		order:     order,
	}

	if err := tdo.validate(); err != nil {
		logError.Printf("Validating order failed: %s", err)
		return ErrOrderValidation
	}

	orderString := tdo.orderString()
	logDebug.Printf("orderString: %s", orderString)

	err := postRequest(opURL(opOptionTrade, s.sourceID, "", orderString), nil, &s.amtdOrder, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error calling option trade: %s\n", err)
		return fmt.Errorf("Error calling option trade: %s", err)
	}

	if s.amtdOrder.Result == "FAIL" {
		logError.Printf("Error from TD: %s\n", s.amtdOrder.Error)
		return fmt.Errorf("Error from TD: %s\n", s.amtdOrder.Error)
	} else if s.amtdOrder.Result == "OK" && s.amtdOrder.OrderWrapper.Error != "" {
		logError.Printf("Error from TD: %s\n", s.amtdOrder.OrderWrapper.Error)
		return fmt.Errorf("Error from TD: %s\n", s.amtdOrder.OrderWrapper.Error)
	}

	return nil
}

//CancelOrder cancels an order that has been accepted by the broker. Note currently it only cancels the first order in orderids
func (s *Session) CancelOrder(orderids []string) error {
	logInfo.Printf("CancelOrder\n")

	if len(orderids) > 1 {
		logError.Printf("Currently unsupporting multiple cancel orders")
		return fmt.Errorf("Currently unsupporting multiple cancel orders")
	}

	s.Lock()
	defer s.Unlock()

	err := postRequest(opURL(opCancelOrder, s.sourceID, "", orderids...), nil, &s.amtdCancelOrder, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error calling option trade: %s\n", err)
		return fmt.Errorf("Error calling option trade: %s", err)
	}

	if s.amtdCancelOrder.Result == "FAIL" {
		logError.Printf("Error from TD: %s\n", s.amtdCancelOrder.Error)
		return fmt.Errorf("Error from TD: %s\n", s.amtdCancelOrder.Error)
	} else if s.amtdCancelOrder.Result == "OK" && s.amtdCancelOrder.CancelOrderMessages.CanceledOrder[0].Error != "" {
		logError.Printf("Error from TD: %s\n", s.amtdCancelOrder.CancelOrderMessages.CanceledOrder[0].Error)
		return fmt.Errorf("Error from TD: %s\n", s.amtdCancelOrder.CancelOrderMessages.CanceledOrder[0].Error)
	}

	return nil
}

//RetrieveOrderBook retrieves the orders from the brokerage firm in an asynchronous method (even though the function call is syncronous)
func (s *Session) RetrieveOrderBook(accountid string, ob *orderbook.OrderBook) error {
	logInfo.Printf("RetrieveOrderBook\n")

	if s.amtdLogin != nil && s.amtdLogin.Result != "OK" {
		return fmt.Errorf("Not logged in")
	}

	s.Lock()
	defer s.Unlock()

	orderStatusParams := url.Values{"accountid": {strconv.Itoa(int(s.amtdLogin.Login.Accounts[0].AccountID))}, "time": {}, "orderid": {}, "type": {}, "fromdate": {}, "todate": {}, "days": {}, "numrec": {}, "underlying": {}}

	err := postRequest(opURL(opOrderStatus, s.sourceID, ""), orderStatusParams, &s.amtdOrderStatus, s.isLoggedIn(), s.sessionID())
	if err != nil {
		logError.Printf("Error getting order status: %s\n", err)
		return fmt.Errorf("Error getting order status: %s", err)
	}

	//debugMarshal(s.amtdOrderStatus)

	if s.amtdOrderStatus.Result != "OK" {
		return errors.New("Error received from OrderStatus service")
	}

	for _, currOrderStatus := range s.amtdOrderStatus.OrderStatusList.OrderStatus {
		os := orderstatus.New()

		os.SetStatus(currOrderStatus.DisplayStatus)
		os.SetOrderID(currOrderStatus.Order.OrderID)
		os.SetAction(orderAction(currOrderStatus.Order))
		os.SetOrderType(orderTypeValue(currOrderStatus.Order))
		os.SetQuantity(int(currOrderStatus.Order.Quantity))
		//FUTURE : assume that 0 price is "invalid" may need to adjuts this type
		//of logic in future to something like "Exists() or Null()" or whatever boolean function
		if currOrderStatus.Order.LimitPrice.Value.Cmp(big.NewRat(0, 1)) > 0 {
			os.SetPrice(currOrderStatus.Order.LimitPrice)
		} else {
			os.SetPrice(currOrderStatus.Order.StopPrice)
		}
		os.SetSymbol(currOrderStatus.Order.Security.Symbol)
		os.SetExpire(orderExpiry(currOrderStatus.Order))
		os.SetRouting(mapRouting(currOrderStatus.Order.ActualDestination.OptionExchange))

		ob.AddUpdateOrderStatus(os)
	}

	return nil
}

func mapRouting(route string) orderconst.OrderExchange {
	switch route {
	case "Auto":
		return orderconst.Auto
	case "ISE":
		return orderconst.ISEX
	case "CBOE":
		return orderconst.CBOE
	case "AMEX":
		return orderconst.AMEX
	case "PHLX":
		return orderconst.PHLX
	case "PACX":
		return orderconst.PACX
	case "BOSX":
		return orderconst.BOSX
	default:
		// TODO: FOLLOWUP because documentation doesn't say what values can come back
		//i'm using this as debug
		logError.Printf("FIX THIS ROUTE:: %s", route)
		return orderconst.Auto
	}
}

//orderAction is just a start for now to get the info I want, later it'll need a more formal expansion in future
func orderAction(o amtd.OrderXML) orderconst.OrderAction {
	//Option order
	switch o.Security.AssetType {
	case "O":
		switch {
		case o.OpenClose == "O" && o.Action == "B":
			return orderconst.BuyToOpen

		case o.OpenClose == "O" && o.Action == "S":
			return orderconst.SellToOpen

		case o.OpenClose == "C" && o.Action == "B":
			return orderconst.BuyToClose

		case o.OpenClose == "C" && o.Action == "S":
			return orderconst.SellToClose
		}

		//FUTURE: finish
	}

	return orderconst.InvalidOrderAction
}

//orderExpiry returns the mapping between TD expiration type and generic expiration type
func orderExpiry(o amtd.OrderXML) orderconst.OrderExpiry {
	switch o.TimeInForce.Session {
	case "G":
		return orderconst.GTC
	case "D":
		return orderconst.Day

		// more cases here in future

	}

	return orderconst.InvalidOrderExpiry
}

//orderTypeValue returns the mapping between TD order type type and generic order type
func orderTypeValue(o amtd.OrderXML) orderconst.OrderType {
	switch o.OrderType {
	case "M":
		return orderconst.Market
	case "L":
		return orderconst.Limit
	case "S":
		return orderconst.StopMarket
	case "X":
		return orderconst.StopLimit
	default:
		if o.LimitPrice.Value.Cmp(big.NewRat(0, 1)) > 0 {
			return orderconst.Limit
		}
	}

	return orderconst.InvalidOrderType

}
