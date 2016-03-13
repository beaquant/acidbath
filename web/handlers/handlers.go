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

package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"text/template"

	"log"

	genericBroker "github.com/marklaczynski/acidbath/broker/generic"
	"github.com/marklaczynski/acidbath/dm/asset"
	"github.com/marklaczynski/acidbath/dm/optionchain/option"
	"github.com/marklaczynski/acidbath/dm/order"
	"github.com/marklaczynski/acidbath/dm/orderbook"
	"github.com/marklaczynski/acidbath/dm/portfolio"
	eventProcFactory "github.com/marklaczynski/acidbath/eventproc/factory"
	"github.com/marklaczynski/acidbath/lib/financial"
	"github.com/marklaczynski/acidbath/lib/mjlog"
	"github.com/marklaczynski/acidbath/lib/orderconst"
)

const (
	baseHTML = "/src/github.com/marklaczynski/acidbath/web/html/"
)

var (
	templates *template.Template
	logInfo   = log.New(mjlog.CreateInfoFile(), "INFO  [handlers]: ", log.LstdFlags|log.Lshortfile)
	logDebug  = log.New(mjlog.CreateDebugFile(), "DEBUG [handlers]: ", log.LstdFlags|log.Lshortfile)
	logError  = log.New(mjlog.CreateErrorFile(), "ERROR [handlers]: ", log.LstdFlags|log.Lshortfile)

	//tmp TODO: clean this up...
	readyToSendOptUpdatesChan chan bool                                = make(chan bool)
	userSelectedStock         *asset.Stock                             = nil
	loginFunc                 func(brokerSession genericBroker.Broker) = nil
)

func init() {
	templates = template.New("login.xhtml").Delims("{{%", "%}}")
	_, err := templates.ParseFiles(os.Getenv("GOPATH") + baseHTML + "login.xhtml")
	if err != nil {
		panic(err)
	}
}

func RegisterLoginActivity(fn func(brokerSession genericBroker.Broker)) {
	loginFunc = fn
}

func renderTemplate(pagename string, w http.ResponseWriter, data interface{}) {
	logInfo.Printf("Rendering page: %s\n", pagename)

	err := templates.ExecuteTemplate(w, pagename, data)
	if err != nil {
		logError.Printf("Error rendering page: %s with error: %s\n", pagename, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

type AppHandler func(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error

func MakeHandler(fn AppHandler, brokerSession genericBroker.Broker) http.HandlerFunc {
	logInfo.Printf("Making handler\n")

	return func(w http.ResponseWriter, r *http.Request) {
		// FUTURE: eventually I'll pass in a user token with each request,
		// so that I won't need to use the brokerSession for right now it's not a high priority
		err := fn(w, r, brokerSession)
		if err != nil {
			logDebug.Printf("Received an error, and quietly discarding")
		}
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	logInfo.Printf("RootHandler\n")

	renderTemplate("login.xhtml", w, nil)

	return nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	logInfo.Printf("LogoutHandler\n")

	err := brokerSession.Logout()
	if err != nil {
		logError.Printf("Error logging out with error: %s\n", err)
		renderTemplate("login.xhtml", w, nil)
		return fmt.Errorf("Error logging out with error: %s\n", err)
	}

	renderTemplate("login.xhtml", w, nil)
	return nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	logInfo.Printf("LoginHandler\n")

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	type loginResponse struct {
		Token string `json:"token"`
		Error string `json:"error"`
	}

	var loginParams struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginParams)
	if err != nil {
		logError.Printf("Error parsing parameters: %#v\n", loginParams)
		return fmt.Errorf("Error parsing parameters: %#v\n", loginParams)
	}

	if loginParams.Login == "" || loginParams.Password == "" {
		logError.Printf("Not enough params to login: %#v\n", loginParams)
		//logError.Printf("Not enough params to login: %#v\n", r.Form)

		emptyToken := loginResponse{Error: fmt.Sprintf("Not enough params to login")}
		if err := json.NewEncoder(w).Encode(emptyToken); err != nil {
			logError.Printf("System error: %s", err)
			return err
		}
		return nil
	}

	err = brokerSession.Login(loginParams.Login, loginParams.Password)
	if err != nil {
		emptyToken := loginResponse{Error: fmt.Sprintf("Error logging in: %s", err)}
		if err := json.NewEncoder(w).Encode(emptyToken); err != nil {
			logError.Printf("System error: %s", err)
			return err
		}

		return nil
	}

	if loginFunc != nil {
		loginFunc(brokerSession)
	}

	err = brokerSession.RetrievePortfolio(portfolio.NewPortfolio())
	if err != nil {
		logError.Printf("Failed to get Balance and Positions: %s\n", err)
		return fmt.Errorf("Failed to get Balance and Positions: %s\n", err)
	}

	//look at: https://github.com/dgrijalva/jwt-go to use a real token in future or ?use token from td?
	successToken := loginResponse{Token: "you're in with a fake token"}
	if err := json.NewEncoder(w).Encode(successToken); err != nil {
		logError.Printf("System error: %s", err)
		return err
	}

	return nil
}

func ReqOrderBookHandler(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	logInfo.Printf("RetrieveOrderBookHandler\n")

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	type uiOrderStatusModel struct {
		Status    string
		Action    string
		OrderID   string
		OrderType string
		Quantity  string
		Price     string
		Symbol    string
		Expire    string
		Routing   string
		Event     string
	}

	type uiOrderBookModel struct {
		UiOrderStatuses map[string]uiOrderStatusModel
	}

	ob := orderbook.New()
	if err := brokerSession.RetrieveOrderBook("", ob); err != nil {
		emptyResponse := uiOrderBookModel{}
		if err := json.NewEncoder(w).Encode(emptyResponse); err != nil {
			logError.Printf("System error: %s", err)
			return err
		}

		return nil
	}

	uiOrderBookResponse := uiOrderBookModel{
		UiOrderStatuses: make(map[string]uiOrderStatusModel),
	}

	for _, currOrderStatus := range ob.OrderStatuses() {
		tmpUiOrderStatus := uiOrderStatusModel{
			Status:    currOrderStatus.Status(),
			Action:    currOrderStatus.Action().String(),
			OrderID:   currOrderStatus.OrderID(),
			OrderType: currOrderStatus.OrderType().String(),
			Quantity:  fmt.Sprintf("%d", currOrderStatus.Quantity()),
			Price:     fmt.Sprintf("%s", currOrderStatus.Price().Value.FloatString(2)),
			Symbol:    currOrderStatus.Symbol(),
			Expire:    currOrderStatus.Expire().String(),
			Routing:   currOrderStatus.Routing().String(),
			Event:     "",
		}

		uiOrderBookResponse.UiOrderStatuses[tmpUiOrderStatus.OrderID] = tmpUiOrderStatus
	}

	if err := json.NewEncoder(w).Encode(uiOrderBookResponse); err != nil {
		logError.Printf("System error: %s", err)
		return err
	}

	return nil
}

func ReleaseOptionUpdatesEventsHandler(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	logDebug.Printf("ReleaseOptionUpdatesEventsHandler")
	readyToSendOptUpdatesChan <- true
	return nil
}

func ReqOptChainHandler(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	logInfo.Printf("ReqOptChainHandler\n")

	readyToSendOptUpdatesChan <- false

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var reqOptChainHandlerParams struct {
		Symbol string `json:"symbol"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqOptChainHandlerParams)
	if err != nil {
		logError.Printf("Error parsing parameters: %#v\n", reqOptChainHandlerParams)
		return fmt.Errorf("Error parsing parameters: %#v\n", reqOptChainHandlerParams)
	}

	logDebug.Printf("Requesting data for: %s", reqOptChainHandlerParams.Symbol)

	if reqOptChainHandlerParams.Symbol == "" {
		logError.Printf("Not enough params: %#v\n", reqOptChainHandlerParams.Symbol)
		return fmt.Errorf("Not enough params : %#v\n", reqOptChainHandlerParams.Symbol)
	}

	// snapshot xml
	// TODO: wrap this function in something like "IsValidSymbol"
	previousUserSelectedSTock := userSelectedStock
	userSelectedStock = asset.NewStock(reqOptChainHandlerParams.Symbol)
	err = brokerSession.RetrieveSnapshot(reqOptChainHandlerParams.Symbol, asset.EquityType, userSelectedStock)
	if err != nil {
		logError.Printf("Error Getting Snapshot: %s", err)
		return fmt.Errorf("Error Getting Snapshot: %s", err)
	}

	if previousUserSelectedSTock != nil {
		err = brokerSession.RemoveStockOptionsFromStream(previousUserSelectedSTock)
		if err != nil {
			logInfo.Printf("Error calling RemoveStockOptionsFromStream: %s\n", err)
			return fmt.Errorf("Error calling RemoveStockOptionsFromStream: %s\n", err)
		}
	}

	err = brokerSession.AddStockOptionsToStream(userSelectedStock)
	if err != nil {
		logInfo.Printf("Error calling AddStockOptionsToStream: %s\n", err)
		return fmt.Errorf("Error calling AddStockOptionsToStream: %s\n", err)
	}

	type UiOptionType struct {
		Bid    string
		Ask    string
		Ticker string
	}

	type strike struct {
		Strike string
		Date   string
		Option map[string]UiOptionType
	}

	type expiration struct {
		Strikes map[string]strike
	}

	var optionChain struct {
		Expirations map[string]expiration
	}

	// send data to ui
	//oc, doneChan := brokerSession.OptionChain()
	oc := userSelectedStock.OptionChain()

	optionChain.Expirations = make(map[string]expiration)
	for _, exp := range oc.SortedExpirations() {

		uiExp := expiration{Strikes: make(map[string]strike)}

		for _, stk := range exp.Strikes() {
			uiStk := strike{Option: make(map[string]UiOptionType)}

			uiStk.Strike = fmt.Sprintf("%.2f", stk.Strike())
			uiStk.Date = fmt.Sprintf("%s", exp.Date())
			if stk.Option(option.CALL) != nil {
				uiCallOption := UiOptionType{}
				uiCallOption.Bid = stk.Option(option.CALL).Bid().Value.FloatString(2)
				uiCallOption.Ask = stk.Option(option.CALL).Ask().Value.FloatString(2)
				uiCallOption.Ticker = stk.Option(option.CALL).OptionTickerSymbol()

				uiStk.Option[fmt.Sprintf("%s", option.CALL)] = uiCallOption

			} else if stk.Option(option.PUT) != nil {
				uiPutOption := UiOptionType{}
				uiPutOption.Bid = stk.Option(option.PUT).Bid().Value.FloatString(2)
				uiPutOption.Ask = stk.Option(option.PUT).Ask().Value.FloatString(2)
				uiPutOption.Ticker = stk.Option(option.PUT).OptionTickerSymbol()

				uiStk.Option[fmt.Sprintf("%s", option.PUT)] = uiPutOption
			}

			uiExp.Strikes[fmt.Sprintf("%.2f", stk.Strike())] = uiStk
		}

		optionChain.Expirations[fmt.Sprintf("%s", exp.Date())] = uiExp
	}

	//doneChan <- true

	if err := json.NewEncoder(w).Encode(optionChain); err != nil {
		logError.Printf("System error: %s", err)
		return err
	}

	logDebug.Printf("new request in progress = false")
	return nil
}

func PortfolioUpdateEvent(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	f, ok := w.(http.Flusher)
	if !ok {
		logError.Printf("Error with Serve HTTP")
		http.Error(w, "Streaming unsupported! make better handling in future", http.StatusInternalServerError)
		return nil
	}

	conClosedNotification := w.(http.CloseNotifier).CloseNotify()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	type uiPortfolioModel struct {
		NetLiquidity      string
		OptionBuyingPower string
	}

	portfolioChan := brokerSession.RegisterPortfolioUpdateChan("handler")

	// If I don't include this for loop, then i get an error on the web browser. I think it's because the connection gets lost if this routine ends
	logDebug.Printf("starting for loop in PortfolioUpdateEvent\n")
	for {
		select {
		case recentPortfolio := <-portfolioChan:
			logInfo.Printf("Received an recentPortfolio update %v\n", recentPortfolio)
			logInfo.Printf("Received an recentPortfolio update %v\n", recentPortfolio.Balance())

			uiPortfolio := uiPortfolioModel{
				NetLiquidity:      fmt.Sprintf("%.2f", recentPortfolio.Balance().NetLiquidity()),
				OptionBuyingPower: fmt.Sprintf("%.2f", recentPortfolio.Balance().OptionBuyingPower()),
			}

			data, err := json.Marshal(uiPortfolio)
			if err != nil {
				logError.Printf("Could not marshal recentPortfolio into json\n")
				return errors.New("Could not marshal recentPortfolio into json\n")
			}
			logDebug.Printf("data:%s\n", data)

			//"data:" must the the first thing sent (part of the SSE contract) and ended with 2 newlines \n\n
			fmt.Fprintf(w, "data:%s\n\n", data)
			f.Flush()
		case <-conClosedNotification:
			logDebug.Printf("BnP HTTP Connection closed\n")
			return nil
		}
	}

	return nil
}

func OrderUpdateEvent(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	f, ok := w.(http.Flusher)
	if !ok {
		logError.Printf("Error with Serve HTTP")
		http.Error(w, "Streaming unsupported! make better handling in future", http.StatusInternalServerError)
		return nil
	}

	conClosedNotification := w.(http.CloseNotifier).CloseNotify()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	type uiOrderMessageModel struct {
		OrderID    string
		OrderEvent string
	}

	orderMessageChan := brokerSession.RegisterOrderUpdateChan("handler")

	// If I don't include this for loop, then i get an error on the web browser. I think it's because the connection gets lost if this routine ends
	logDebug.Printf("starting for loop in OrderUpdateEvent\n")
	for {
		select {
		case orderMessage := <-orderMessageChan:
			logInfo.Printf("Received an orderMessage update %v\n", orderMessage)

			uiOrder := uiOrderMessageModel{
				OrderID:    orderMessage.OrderID(),
				OrderEvent: orderMessage.OrderEvent().String(),
			}

			data, err := json.Marshal(uiOrder)
			if err != nil {
				logError.Printf("Could not marshal orderMessage into json\n")
				return errors.New("Could not marshal orderMessage into json\n")
			}
			logDebug.Printf("data:%s\n", data)

			//"data:" must the the first thing sent (part of the SSE contract) and ended with 2 newlines \n\n
			fmt.Fprintf(w, "data:%s\n\n", data)
			f.Flush()
		case <-conClosedNotification:
			logDebug.Printf("Order HTTP Connection closed\n")
			return nil
		}
	}

	return nil
}

func OptionUpdateEvent(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	f, ok := w.(http.Flusher)
	if !ok {
		logError.Printf("Error with Serve HTTP")
		http.Error(w, "Streaming unsupported! make better handling in future", http.StatusInternalServerError)
		return nil
	}

	conClosedNotification := w.(http.CloseNotifier).CloseNotify()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	type uiOptionModel struct {
		Expiration string
		Strike     string
		OptionType string

		TickerSymbol string
		Bid          string
		Ask          string
	}

	optionChan := brokerSession.RegisterOptionUpdateChan("handler")

	// If I don't include this for loop, then i get an error on the web browser. I think it's because the connection gets lost if this routine ends
	logDebug.Printf("starting for loop in OptionUpdateEvent\n")
	ready := false
	for {
		select {
		case ready = <-readyToSendOptUpdatesChan:
		case o := <-optionChan:
			//update local dm
			if userSelectedStock != nil && userSelectedStock.OptionChain() != nil {
				userSelectedStock.OptionChain().Option(o.OptionTickerSymbol()).SetBid(o.Bid())
				userSelectedStock.OptionChain().Option(o.OptionTickerSymbol()).SetAsk(o.Ask())
			}

			//send update to ui dm
			if ready == true {
				logInfo.Printf("Received an option update %v\n", o)

				opt := uiOptionModel{
					Expiration: fmt.Sprintf("%s", o.ExpirationDate()),
					Strike:     fmt.Sprintf("%.2f", o.Strike()),
					OptionType: fmt.Sprintf("%s", o.OptionType()),

					TickerSymbol: o.OptionTickerSymbol(),
					Bid:          fmt.Sprintf("%s", o.Bid().Value.FloatString(2)),
					Ask:          fmt.Sprintf("%s", o.Ask().Value.FloatString(2)),
				}

				data, err := json.Marshal(opt)
				if err != nil {
					logError.Printf("Could not marshal option into json\n")
					return errors.New("Could not marshal option into json\n")
				}
				logDebug.Printf("data:%s\n", data)

				//"data:" must the the first thing sent (part of the SSE contract) and ended with 2 newlines \n\n
				fmt.Fprintf(w, "data:%s\n\n", data)
				f.Flush()
			} else {
				logDebug.Printf("flushing the data, because ui is not ready... FUTURE try to buffer these up and send them, but that'll be a lot of work")
			}
		case <-conClosedNotification:
			logDebug.Printf("Option HTTP Connection closed\n")
			return nil
		}
	}

	return nil
}

func TrackOptionHandler(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	logInfo.Printf("TrackOptionHandler\n")

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var trackOptionParam struct {
		Symbol string `json:"symbol"`
	}

	err := json.NewDecoder(r.Body).Decode(&trackOptionParam)
	if err != nil {
		logError.Printf("Error parsing parameters: %#v\n", trackOptionParam)
		return fmt.Errorf("Error parsing parameters: %#v\n", trackOptionParam)
	}

	if trackOptionParam.Symbol == "" {
		logError.Printf("Not enough params: %#v\n", trackOptionParam.Symbol)
		return fmt.Errorf("Not enough params to track option: %#v\n", trackOptionParam.Symbol)
	}

	logInfo.Printf("tracking: %s \n", trackOptionParam.Symbol)
	list, _ := brokerSession.AddOptionToStrategy(userSelectedStock.OptionChain().Option(trackOptionParam.Symbol), eventProcFactory.Reference)

	var instrumentList struct {
		Instrument []string
	}

	instrumentList.Instrument = list

	if err := json.NewEncoder(w).Encode(instrumentList); err != nil {
		logError.Printf("System error: %s", err)
		return err
	}

	return nil
}

func UntrackOptionHandler(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	logInfo.Printf("UntrackOptionHandler\n")

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var trackOptionParam struct {
		Symbol string `json:"symbol"`
	}

	err := json.NewDecoder(r.Body).Decode(&trackOptionParam)
	if err != nil {
		logError.Printf("Error parsing parameters: %#v\n", trackOptionParam)
		return fmt.Errorf("Error parsing parameters: %#v\n", trackOptionParam)
	}

	if trackOptionParam.Symbol == "" {
		logError.Printf("Not enough params: %#v\n", trackOptionParam.Symbol)
		return fmt.Errorf("Not enough params to track option: %#v\n", trackOptionParam.Symbol)
	}

	logInfo.Printf("tracking: %s \n", trackOptionParam.Symbol)
	list, _ := brokerSession.RemoveOptionFromStrategy(userSelectedStock.OptionChain().Option(trackOptionParam.Symbol), eventProcFactory.Reference)

	var instrumentList struct {
		Instrument []string
	}

	instrumentList.Instrument = list

	if err := json.NewEncoder(w).Encode(instrumentList); err != nil {
		logError.Printf("System error: %s", err)
		return err
	}

	return nil
}

func TestOrderHandler(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	logInfo.Printf("TestOrderHandler\n")

	order := order.New()

	order.SetAction(orderconst.SellToOpen)
	order.SetExpire(orderconst.Day)
	order.SetOrderType(orderconst.Limit)
	order.SetQuantity(1)
	order.SetPrice(financial.Money{big.NewRat(100, 1)})
	//TODO: fix this up this is hardcoded for right now, which is fine, but after 03/17/2017 this will stop working.
	order.SetSymbol("SPY_031717P100")

	err := brokerSession.SendSingleLegOptionTrade(order)
	if err != nil {
		logError.Printf("Error sending order")
		return errors.New("Error sending order")
	}

	return nil
}

func TestCancelOrderHandler(w http.ResponseWriter, r *http.Request, brokerSession genericBroker.Broker) error {
	logInfo.Printf("TestCancelOrderHandler\n")

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var cancelOrderParams struct {
		OrderID string `json:"orderid"`
	}

	err := json.NewDecoder(r.Body).Decode(&cancelOrderParams)
	if err != nil {
		logError.Printf("Error parsing parameters: %#v\n", cancelOrderParams)
		return fmt.Errorf("Error parsing parameters: %#v\n", cancelOrderParams)
	}

	if cancelOrderParams.OrderID == "" {
		logError.Printf("Not enough params: %#v\n", cancelOrderParams.OrderID)
		return fmt.Errorf("Not enough params to login: %#v\n", cancelOrderParams.OrderID)
	}

	err = brokerSession.CancelOrder([]string{cancelOrderParams.OrderID})
	if err != nil {
		logError.Printf("Error canceling order")
		return errors.New("Error canceling order")
	}

	return nil
}
