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

	"github.com/marklaczynski/acidbath/lib/financial"
	"github.com/marklaczynski/acidbath/lib/types"
)

//OptionChain represents the TD OptionChain structure
type OptionChain struct {
	XMLName            xml.Name `xml:"amtd"`
	Error                       //inline struct
	OptionChainResults optionChainResultsXML
}

type optionChainResultsXML struct {
	XMLName          xml.Name         `xml:"option-chain-results"`
	Error            string           `xml:"error"`
	Symbol           string           `xml:"symbol"`
	Description      string           `xml:"description"`
	Bid              financial.Money  `xml:"bid"`
	Ask              financial.Money  `xml:"ask"`
	BidAskSize       string           `xml:"bid-ask-size"`
	Last             financial.Money  `xml:"last"`
	Open             financial.Money  `xml:"open"`
	High             financial.Money  `xml:"high"`
	Low              financial.Money  `xml:"low"`
	Close            financial.Money  `xml:"close"`
	Volume           string           `xml:"volume"`
	Change           types.XMLFloat64 `xml:"change"`
	QuotePunctuality string           `xml:"quote-punctuality"`
	Time             string           `xml:"time"`
	OptionDate       []optionDateXML  `xml:"option-date"`
}

type optionDateXML struct {
	XMLName          xml.Name          `xml:"option-date"`
	Date             string            `xml:"date"`
	ExpirationType   string            `xml:"expiration-type"` //ENUM
	DaysToExpiration types.XMLInt64    `xml:"days-to-expiration"`
	OptionStrike     []optionStrikeXML `xml:"option-strike"`
}

type optionStrikeXML struct {
	XMLName        xml.Name         `xml:"option-strike"`
	StrikePrice    types.XMLFloat64 `xml:"strike-price"`
	StandardOption bool             `xml:"standard-option"`
	Put            *optionXML       `xml:"put,omitempty"`
	Call           *optionXML       `xml:"call,omitempty"`
}

type optionXML struct {
	// not sure cause one is a put and one is a call XMLName        xml.Name `xml:"option-strike"`
	OptionSymbol      string             `xml:"option-symbol"`
	Description       string             `xml:"description"`
	Bid               financial.Money    `xml:"bid"`
	Ask               financial.Money    `xml:"ask"`
	BidAskSize        string             `xml:"bid-ask-size"`
	Last              financial.Money    `xml:"last"`
	LastTradeDate     string             `xml:"last-trade-date"`
	Volume            types.XMLInt64     `xml:"volume"`
	OpenInterest      types.XMLInt64     `xml:"open-interest"`
	RealTime          string             `xml:"real-time"`
	UnderlyingSymbol  string             `xml:"underlying-symbol"`
	Delta             types.XMLFloat64   `xml:"delta"`
	Gamma             types.XMLFloat64   `xml:"gamma"`
	Theta             types.XMLFloat64   `xml:"theta"`
	Vega              types.XMLFloat64   `xml:"vega"`
	Rho               types.XMLFloat64   `xml:"rho"`
	ImpliedVolatility types.XMLFloat64   `xml:"implied-volatility"`
	TimeValueIndex    types.XMLFloat64   `xml:"time-value-index"`
	Multiplier        types.XMLFloat64   `xml:"multiplier"` //options end (incorrectly documented)
	Change            types.XMLFloat64   `xml:"change"`
	ChangePercent     string             `xml:"change-percent"`
	InTheMoney        bool               `xml:"in-the-money"`
	NearTheMoney      bool               `xml:"near-the-money"`
	TheoreticalValue  financial.Money    `xml:"theoretical-value"`
	DeliverableList   deliverableListXML `xml:"deliverable-list"`
}

type deliverableListXML struct {
	XMLName                xml.Name         `xml:"deliverable-list"`
	CashInLieuDollarAmount types.XMLFloat64 `xml:"cash-in-lieu-dollar-amount"`
	CashDollarAmount       types.XMLFloat64 `xml:"cash-dollar-amount"`
	IndexOption            bool             `xml:"index-option"`
	NotesDescription       string           `xml:"notes-description"`
	Row                    []rowXML         `xml:"row"`
}

type rowXML struct {
	Symbol string         `xml:"symbol"`
	Shares types.XMLInt64 `xml:"shares"`
}
