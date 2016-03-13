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

package amtd

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/marklaczynski/acidbath/dm/optionchain/option"
	"github.com/marklaczynski/acidbath/lib/date"
	"github.com/marklaczynski/acidbath/lib/financial"
	"github.com/marklaczynski/acidbath/lib/types"
)

//SnapshotQuotes represents the TD SnapshotQuotes structure
type SnapshotQuotes struct {
	XMLName   xml.Name `xml:"amtd"`
	Error              //inline struct
	QuoteList quoteListXML
}

type quoteListXML struct {
	XMLName xml.Name   `xml:"quote-list"`
	Error   string     `xml:"error"`
	Quote   []QuoteXML `xml:"quote"`
}

type QuoteXML struct {
	XMLName           xml.Name         `xml:"quote"`
	Error             string           `xml:"error"`
	Symbol            string           `xml:"symbol"`
	Description       string           `xml:"description"`
	Bid               financial.Money  `xml:"bid"`
	Ask               financial.Money  `xml:"ask"`
	BidAskSize        string           `xml:"bid-ask-size"`
	Last              financial.Money  `xml:"last"`
	LastTradeSize     types.XMLInt64   `xml:"last-trade-size"`
	LastTradeDate     string           `xml:"last-trade-date"`
	Open              financial.Money  `xml:"open"`
	High              financial.Money  `xml:"high"`
	Low               financial.Money  `xml:"low"`
	Close             financial.Money  `xml:"close"`
	Volume            types.XMLInt64   `xml:"volume"`
	StrikePrice       types.XMLFloat64 `xml:"strike-price"` //options start
	OpenInterest      types.XMLInt64   `xml:"open-interest"`
	ExpirationMonth   types.XMLInt64   `xml:"expiration-month"`
	ExpirationDay     types.XMLInt64   `xml:"expiration-day"`
	ExpirationYear    types.XMLInt64   `xml:"expiration-year"`
	UnderlyingSymbol  string           `xml:"underlying-symbol"`
	PutCall           string           `xml:"put-call"`
	Delta             types.XMLFloat64 `xml:"delta"`
	Gamma             types.XMLFloat64 `xml:"gamma"`
	Theta             types.XMLFloat64 `xml:"theta"`
	Vega              types.XMLFloat64 `xml:"vega"`
	Rho               types.XMLFloat64 `xml:"rho"`
	ImpliedVolatility types.XMLFloat64 `xml:"implied-volatility"`
	DTE               types.XMLInt64   `xml:"days-to-expiration"`
	TimeValueIndex    financial.Money  `xml:"time-value-index"`
	Multiplier        types.XMLFloat64 `xml:"multiplier"` //options end (incorrectly documented)
	YearHigh          financial.Money  `xml:"year-high"`
	YearLow           financial.Money  `xml:"year-low"`
	RealTime          bool             `xml:"real-time"`
	Exchange          string           `xml:"exchange"`
	AssetType         string           `xml:"asset-type"`
	Change            types.XMLFloat64 `xml:"change"`
	ChangePercent     string           `xml:"change-percent"`
	Nav               types.XMLFloat64 `xml:"nav"` //funds start
	Offer             types.XMLFloat64 `xml:"offer"`
}

func (amtdQuote QuoteXML) NewOption() *option.Option {
	o := option.NewNilOption()
	amtdQuote.ProcessOption(o)
	return o
}

func (amtdQuote QuoteXML) ProcessOption(o *option.Option) {
	location, _ := time.LoadLocation(date.LocalLocation)
	exp, _ := date.New(time.Date(int(amtdQuote.ExpirationYear), time.Month(amtdQuote.ExpirationMonth), int(amtdQuote.ExpirationDay), 0, 0, 0, 0, location))

	if strings.Contains(strings.Split(amtdQuote.Symbol, "_")[1], "C") {
		o.SetOptionType(option.CALL)
	} else if strings.Contains(strings.Split(amtdQuote.Symbol, "_")[1], "P") {
		o.SetOptionType(option.PUT)
	} else {
		panic(fmt.Sprintf("Didn't parse any option type from symbol: %s", amtdQuote.Symbol))
	}

	o.SetSymbol(amtdQuote.Symbol)
	o.SetStrike(float64(amtdQuote.StrikePrice))
	o.SetExpirationDate(exp)
	o.SetMultiplier(float64(amtdQuote.Multiplier))

	o.SetLast(amtdQuote.Last)

	o.SetBid(amtdQuote.Bid)
	o.SetAsk(amtdQuote.Ask)
	//o.SetVolume(int64(amtdQuote.Volume))
	o.SetDelta(float64(amtdQuote.Delta))
	o.SetTheta(float64(amtdQuote.Theta))
	o.SetGamma(float64(amtdQuote.Gamma))
	o.SetVega(float64(amtdQuote.Vega))
}
