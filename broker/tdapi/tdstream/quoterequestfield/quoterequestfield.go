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

package quoterequestfield

//QuoteColumnNumber is an index that represent a column of information returned from TD. 0 represents "Symbol", 1 represents "Bid", etc
type QuoteColumnNumber int

func (colNum QuoteColumnNumber) String() string {

	switch colNum {
	case Symbol:
		return "Symbol"
	case Bid:
		return "Bid"
	case Ask:
		return "Ask"
	case Last:
		return "Last"
	case BidSize:
		return "BidSize"
	case AskSize:
		return "AskSize"
	case BidID:
		return "BidID"
	case AskID:
		return "AskID"
	case Volume:
		return "Volume"
	case LastSize:
		return "LastSize"
	case TradeTime:
		return "Trade Time"
	case QuoteTime:
		return "Quote Time"
	case High:
		return "High"
	case Low:
		return "Low"
	case Tick:
		return "Tick"
	case Close:
		return "Close"
	case EXChange:
		return "EXChange"
	case Marginable:
		return "Marginable"
	case Shortable:
		return "Shortable"
	case QuoteDate:
		return "Quote Date"
	case TradeDate:
		return "Trade Date"
	case Volatility:
		return "Volatility"
	case Description:
		return "Description"
	case TradeID:
		return "TradeID"
	case Digits:
		return "Digits"
	case Open:
		return "Open"
	case Change:
		return "Change"
	case WeekHigh52:
		return "WeekHigh52"
	case WeekLow52:
		return "WeekLow52"
	case PERatio:
		return "PERatio"
	case DividendAmt:
		return "DividendAmt"
	case DividendYield:
		return "DividendYield"
	case Nav:
		return "Nav"
	case Fund:
		return "Fund"
	case ExchangeName:
		return "ExchangeName"
	case DividendDate:
		return "DividendDate"
	case LastMarketHours:
		return "LastMarketHours"
	case LastSizeMarketHours:
		return "LastSizeMarketHours"
	case TradeDateMarketHours:
		return "TradeDateMarketHours"
	case TradeTimeMarketHours:
		return "TradeTimeMarketHours"
	case ChangeMarketHours:
		return "ChangeMarketHours"
	case IsRegularMarketQuote:
		return "IsRegularMarketQuote"
	case IsRegularMarketTrade:
		return "IsRegularMarketTrade"

	}
	return ""
}

//Following constants represent a dev friendly name to an index
const (
	Symbol               QuoteColumnNumber = 0
	Bid                                    = 1
	Ask                                    = 2
	Last                                   = 3
	BidSize                                = 4
	AskSize                                = 5
	BidID                                  = 6
	AskID                                  = 7
	Volume                                 = 8
	LastSize                               = 9
	TradeTime                              = 10
	QuoteTime                              = 11
	High                                   = 12
	Low                                    = 13
	Tick                                   = 14
	Close                                  = 15
	EXChange                               = 16
	Marginable                             = 17
	Shortable                              = 18 //skip
	QuoteDate                              = 22
	TradeDate                              = 23
	Volatility                             = 24
	Description                            = 25
	TradeID                                = 26
	Digits                                 = 27
	Open                                   = 28
	Change                                 = 29
	WeekHigh52                             = 30
	WeekLow52                              = 31
	PERatio                                = 32
	DividendAmt                            = 33
	DividendYield                          = 34 //skip
	Nav                                    = 37
	Fund                                   = 38
	ExchangeName                           = 39
	DividendDate                           = 40
	LastMarketHours                        = 41
	LastSizeMarketHours                    = 42
	TradeDateMarketHours                   = 43
	TradeTimeMarketHours                   = 44
	ChangeMarketHours                      = 45
	IsRegularMarketQuote                   = 46
	IsRegularMarketTrade                   = 47
)
