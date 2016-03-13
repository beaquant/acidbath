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

	"github.com/marklaczynski/acidbath/lib/types"
)

//Login represents the TD Login structure
type Login struct {
	XMLName xml.Name `xml:"amtd"`
	Error            //inline struct
	Login   loginXML
}

type loginXML struct {
	XMLName             xml.Name        `xml:"xml-log-in"`
	SessionID           string          `xml:"session-id"`
	UserID              string          `xml:"user-id"`
	Cdi                 string          `xml:"cdi"`
	Timeout             types.XMLByte   `xml:"timeout"`
	LoginTime           string          `xml:"login-time"`
	AssociatedAccountID types.XMLUInt64 `xml:"associated-account-id"`
	NyseQuotes          string          `xml:"nyse-quotes"`
	NasdaqQuotes        string          `xml:"nasdaq-quotes"`
	OpraQuotes          string          `xml:"opra-quotes"`
	AmexQuotes          string          `xml:"amex-quotes"`
	CMEQuotes           string          `xml:"cme-quotes"`
	ICEQuotes           string          `xml:"ice-quotes"`
	ForexQuotes         string          `xml:"forex-quotes"`
	ExchangeStatus      string          `xml:"exchange-status"`
	Options360          bool            `xml:"authorizations>options360"` //doesn't exist in TD doc, but it's returned in API
	Accounts            []account       `xml:"accounts>account"`
}

type account struct {
	XMLName           xml.Name        `xml:"account"`
	AccountID         types.XMLUInt64 `xml:"account-id"`
	DisplayName       string          `xml:"display-name"`
	Cdi               string          `xml:"cdi"`
	Description       string          `xml:"description"`
	AssociatedAccount bool            `xml:"associated-account"`
	Company           string          `xml:"company"`
	Segment           string          `xml:"segment"`
	Unified           bool            `xml:"unified"`
	Preferences       preferencesXML
	Authorizations    authorizationsXML
}

type preferencesXML struct {
	XMLName                         xml.Name       `xml:"preferences"`
	ExpressTrading                  bool           `xml:"express-trading"`
	OptionDirectRouting             bool           `xml:"option-direct-routing"`
	StockDirectRouting              bool           `xml:"stock-direct-routing"`
	DefaultStockAction              string         `xml:"default-stock-action"`
	DefaultStockOrder               string         `xml:"default-stock-order-type"`
	DefaultStockQuantity            types.XMLInt64 `xml:"default-stock-quantity"` //this is really an int, but if it's empty there's a ParseInt error
	DefaultStockExpiration          string         `xml:"default-stock-expiration"`
	DefaultStockSpecailInstructions string         `xml:"default-stock-special-instructions"`
	DefaultStockRouting             string         `xml:"default-stock-routing"`
	DefaultStockDisplaySize         types.XMLInt64 `xml:"default-stock-display-size"` //this is really an int, but if it's empty there's a ParseInt error
	StockTaxLotMethod               string         `xml:"stock-tax-lot-method"`
	OptionTaxLotMethod              string         `xml:"option-tax-lot-method"`
	MutualFundTaxLotMethod          string         `xml:"mutual-fund-tax-lot-method"`
	DefaultAdvancedToolLaunch       string         `xml:"default-advanced-tool-launch"`
}

type authorizationsXML struct {
	XMLName        xml.Name `xml:"authorizations"`
	Apex           bool     `xml:"apex"`
	Level2         bool     `xml:"level2"`
	StockTrading   bool     `xml:"stock-trading"`
	MarginTrading  bool     `xml:"margin-trading"`
	StreamingNews  bool     `xml:"streaming-news"`
	OptionTrading  string   `xml:"option-trading"`
	Streamer       bool     `xml:"streamer"`
	AdvancedMargin bool     `xml:"advanced-margin"`
}
