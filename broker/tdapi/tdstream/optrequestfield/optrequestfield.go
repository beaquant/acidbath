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

package optrequestfield

//OptionColumnNumber is an index that represent a column of information returned from TD. 0 represents "SYMBOL", 1 represents "Bid", etc
type OptionColumnNumber int

func (colNum OptionColumnNumber) String() string {
	switch colNum {
	case Symbol:
		return "Symbol"
	case Contract:
		return "Contract"
	case Bid:
		return "Bid"
	case Ask:
		return "Ask"
	case Last:
		return "Last"
	case High:
		return "High"
	case Low:
		return "Low"
	case Close:
		return "Close"
	case Volume:
		return "Volume"
	case OpenInterest:
		return "Open Interest"
	case Volatility:
		return "Volatility"
	case QuoteTime:
		return "Quote Time"
	case TradeTime:
		return "Trade Time"
	case InTheMoney:
		return "In The Money"
	case QuoteDate:
		return "Quote Date"
	case TradeDate:
		return "Trade Date"
	case Year:
		return "Year"
	case Multiplier:
		return "Multiplier"
	case Open:
		return "Open"
	case BidSize:
		return "BidSize"
	case AskSize:
		return "AskSize"
	case LastSize:
		return "LastSize"
	case Change:
		return "Change"
	case Strike:
		return "Strike"
	case ContractType:
		return "ContractType"
	case Underlying:
		return "Underlying"
	case Month:
		return "Month"
	case Note:
		return "Note"
	case TimeValue:
		return "TimeValue"
	case DaysToExp:
		return "DaysToExp"
	case DeltaIndex:
		return "DeltaIndex"
	case GammaIndex:
		return "GammaIndex"
	case ThetaIndex:
		return "ThetaIndex"
	case VegaIndex:
		return "VegaIndex"
	case RhoIndex:
		return "RhoIndex"
	}

	return ""
}

//Following constants represent a dev friendly name to an index
const (
	Symbol       OptionColumnNumber = 0
	Contract                        = 1
	Bid                             = 2
	Ask                             = 3
	Last                            = 4
	High                            = 5
	Low                             = 6
	Close                           = 7
	Volume                          = 8
	OpenInterest                    = 9
	Volatility                      = 10
	QuoteTime                       = 11
	TradeTime                       = 12
	InTheMoney                      = 13
	QuoteDate                       = 14
	TradeDate                       = 15
	Year                            = 16
	Multiplier                      = 17
	Open                            = 19
	BidSize                         = 20
	AskSize                         = 21
	LastSize                        = 22
	Change                          = 23
	Strike                          = 24
	ContractType                    = 25
	Underlying                      = 26
	Month                           = 27
	Note                            = 28
	TimeValue                       = 29
	DaysToExp                       = 31
	DeltaIndex                      = 32
	GammaIndex                      = 33
	ThetaIndex                      = 34
	VegaIndex                       = 35
	RhoIndex                        = 36
)
