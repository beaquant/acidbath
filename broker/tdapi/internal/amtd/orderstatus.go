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
)

//OrderStatus represents the TD Order structure returned after sending an order
type OrderStatus struct {
	XMLName         xml.Name `xml:"amtd"`
	Error                    //inline struct
	OrderStatusList orderStatusListXML
}

type orderStatusListXML struct {
	XMLName     xml.Name         `xml:"orderstatus-list"`
	AccountID   string           `xml:"account-id"`
	OrderStatus []orderStatusXML `xml:"orderstatus"`
}

type orderStatusXML struct {
	XMLName                 xml.Name        `xml:"orderstatus"`
	OrderNumber             string          `xml:order-number`
	Cancelable              bool            `xml:cancelable`
	Editable                bool            `xml:editable`
	ComplexOption           bool            `xml:complex-option`
	EnhancedOrder           bool            `xml:"enhanced-order"`
	EnhancedType            string          `xml:"enhanced-type"`  //enum
	DisplayStatus           string          `xml:"display-status"` //enum
	OrderRoutingStatus      string          `xml:"order-routing-status"`
	OrderReceivedDateTime   string          `xml:"order-received-date-time"`
	ReportedTime            string          `xml:"reported-time"`
	ExpiredDateTime         string          `xml:"exp-date-time"`
	CanceledDateTime        string          `xml:"cancel-date-time"`
	RemainingQuantity       float64         `xml:"remaining-quantity"`
	TrailingActivationPrice financial.Money `xml:"trailing-activation-price"`
	UnderlyingSymbol        string          `xml:"underlying-symbol"`
	Order                   OrderXML
	Strategy                string `xml:"strategy"`
	Fills                   fillsXML
	RelatedOrders           relatedOrdersXML
}

//OrderXML is
type OrderXML struct {
	XMLName              xml.Name `xml:"order"`
	Security             orderStatusSecurityXML
	Quantity             float64         `xml:"quantity"`
	OrderID              string          `xml:"order-id"`
	Action               string          `xml:"action"`     //enum
	TradeType            int             `xml:"trade-type"` //enum
	RequestedDestination destinationXML  `xml:"requested-destination"`
	ActualDestination    destinationXML  `xml:"actual-destination"`
	RoutingDisplaySize   int             `xml:"routing-display-size"`
	OrderType            string          `xml:"order-type"` //enum
	LimitPrice           financial.Money `xml:"limit-price"`
	StopPrice            financial.Money `xml:"stop-price"`
	TrailingStopMethod   string          `xml:"trailing-stop-method"` //enum
	SpecialCondition     string          `xml:"special-conditions"`
	TimeInForce          tifXML
	OpenClose            string `xml:"open-close"` //enum
	PutCall              string `xml:"put-call"`   //enum
}

type orderStatusSecurityXML struct {
	XMLName              xml.Name `xml:"security"`
	Symbol               string   `xml:"symbol"`
	SymbolWithTypePrefix string   `xml:"symbol-with-type-prefix"`
	Description          string   `xml:"description"`
	AssetType            string   `xml:"asset-type"` //enum
}

type destinationXML struct {
	//XMLName             xml.Name
	RoutingMode         string `xml:"routing-mode"`         //enum ??
	OptionExchange      string `xml:"option-exchange"`      //enum ??
	ResponseDescription string `xml:"response-description"` //enum ??
}

type tifXML struct {
	XMLName    xml.Name `xml:"time-in-force"`
	Session    string   `xml:"session"` //enum
	Expiration string   `xml:"expiration"`
}

type fillsXML struct {
	XMLName           xml.Name        `xml:"fills"`
	FillID            string          `xml:"fill-id"`
	FillQuantity      float64         `xml:"fill-quantity"`
	FillPrice         financial.Money `xml:"fill-price"`
	ExecutionDateTime string          `xml:"execution-reported-date-time"`
}

type relatedOrdersXML struct {
	XMLName          xml.Name         `xml:"related-orders"`
	RelationshipType string           `xml:"relationship-type"`
	OrderStatus      []orderStatusXML `xml:"orderstatus"`
}
