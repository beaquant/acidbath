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

//Order represents the TD Order structure returned after sending an order
type Order struct {
	XMLName      xml.Name `xml:"amtd"`
	Error                 //inline struct
	OrderWrapper orderWrapperXML
}

type orderWrapperXML struct {
	XMLName     xml.Name `xml:"order-wrapper"`
	OrderString string   `xml:"order-string"`
	Error       string   `xml:"error"`
	Order       orderXML `xml:"order"`
}

type orderXML struct {
	XMLName            xml.Name              `xml:"order"`
	AccountID          string                `xml:"account-id"`
	Security           orderSecurityXML      `xml:"security"`
	Quantity           int                   `xml:"quantity"`
	OrderID            string                `xml:"order-id"`
	Action             string                `xml:"action"`     // enum
	TradeType          int                   `xml:"trade-type"` //enum
	RequestDestination requestDestinationXML `xml:"requested-destination"`
	RoutingDisplaySize int                   `xml:"routing-display-size"`
	OrderType          string                `xml:"order-type"` //enum
	LimitPrice         float64               `xml:"limit-price"`
	StopPrice          types.XMLFloat64      `xml:"stop-price"`
	TimeInForce        timeInForceXML        `xml:"time-in-force"`
	PutCall            string                `xml:"put-call"`
	OpenClose          string                `xml:"open-close"`
}

type orderSecurityXML struct {
	XMLName              xml.Name `xml:"security"`
	Symbol               string   `xml:"symbol"`
	SymbolWithTypePrefix string   `xml:"symbol-with-type-prefix"`
	Description          string   `xml:"description"`
	AssetType            string   `xml:"asset-type"`
	Exchange             string   `xml:"exchange"`
}

type requestDestinationXML struct {
	XMLName             xml.Name `xml:"requested-destination"`
	RoutingMode         string   `xml:"routing-mode"`
	MarketMakerID       string   `xml:"market-maker-id"`
	ResponseDescription string   `xml:"response-description"`
}

type timeInForceXML struct {
	XMLName xml.Name `xml:"time-in-force"`
	Session string   `xml:"session"`
}
