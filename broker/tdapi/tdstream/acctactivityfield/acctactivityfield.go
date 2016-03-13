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

package acctactivityfield

import (
	"encoding/xml"

	"github.com/marklaczynski/acidbath/lib/types"
)

//AcctActivityNumber is an index that represent a column of information returned from TD. 0 represents "SYMBOL", 1 represents "Bid", etc
type AcctActivityNumber int

func (colNum AcctActivityNumber) String() string {
	switch colNum {
	case SubscriptionKey:
		return "SubscriptionKey"
	case AccountNumber:
		return "AccountNumber"
	case MessageType:
		return "MessageType"
	case MessageData:
		return "MessageData"

	}

	return ""
}

//Following constants represent a dev friendly name to an index
const (
	SubscriptionKey AcctActivityNumber = 0
	AccountNumber                      = 1
	MessageType                        = 2
	MessageData                        = 3
)

//AcctActivityMessageType represents a message type for an order
type AcctActivityMessageType string

//Account Activity (order) message types
const (
	Subscribed                AcctActivityMessageType = "SUBSCRIBED"
	Error                                             = "ERROR"
	BrokenTrade                                       = "BrokenTrade"
	ManualExecution                                   = "ManualExecution"
	OrderActivation                                   = "OrderActivation"
	OrderCancelReplaceRequest                         = "OrderCancelReplaceRequest"
	OrderCancelRequest                                = "OrderCancelRequest"
	OrderEntryRequest                                 = "OrderEntryRequest"
	OrderFill                                         = "OrderFill"
	OrderPartialFill                                  = "OrderPartialFill"
	OrderRejection                                    = "OrderRejection"
	TooLateToCancel                                   = "TooLateToCancel"
	UrOut                                             = "UROUT"
)

//UROUTMessage Order cancelation message
type UROUTMessage struct {
	XMLName      xml.Name `xml:"UROUTMessage"`
	OrderGroupID orderGroupIDXML
	Order        orderXML
}

//TooLateToCancelMessage Too late to cancel message
type TooLateToCancelMessage struct {
	XMLName      xml.Name `xml:"TooLateToCancelMessage"`
	OrderGroupID orderGroupIDXML
	Order        orderXML
}

//OrderRejectionMessage Order rejection message
type OrderRejectionMessage struct {
	XMLName      xml.Name `xml:"OrderRejectionMessage"`
	OrderGroupID orderGroupIDXML
	Order        orderXML
}

//OrderPartialFillMessage Order was partially filled
type OrderPartialFillMessage struct {
	XMLName      xml.Name `xml:"OrderPartialFillMessage"`
	OrderGroupID orderGroupIDXML
	Order        orderXML
}

//OrderFillMessage Order filled message
type OrderFillMessage struct {
	XMLName      xml.Name `xml:"OrderFillMessage"`
	OrderGroupID orderGroupIDXML
	Order        orderXML
}

//OrderEntryRequestMessage Order entry request made
type OrderEntryRequestMessage struct {
	XMLName      xml.Name `xml:"OrderEntryRequestMessage"`
	OrderGroupID orderGroupIDXML
	Order        orderXML
}

//OrderCancelRequestMessage Order cancel request message
type OrderCancelRequestMessage struct {
	XMLName      xml.Name `xml:"OrderCancelRequestMessage"`
	OrderGroupID orderGroupIDXML
	Order        orderXML
}

//OrderActivationMessage Order activation message
type OrderActivationMessage struct {
	XMLName      xml.Name `xml:"OrderActivationMessage"`
	OrderGroupID orderGroupIDXML
	Order        orderXML
}

//ManualExecutionMessage Order was manually executed message
type ManualExecutionMessage struct {
	XMLName      xml.Name `xml:"ManualExecutionMessage"`
	OrderGroupID orderGroupIDXML
	Order        orderXML
}

//BrokenTradeMessage Broken trade message
type BrokenTradeMessage struct {
	XMLName      xml.Name `xml:"BrokenTradeMessage"`
	OrderGroupID orderGroupIDXML
	Order        orderXML
}

//OrderCancelReplaceRequestMessage Order cancel & replace request message
type OrderCancelReplaceRequestMessage struct {
	XMLName      xml.Name `xml:"OrderCancelReplaceRequestMessage"`
	OrderGroupID orderGroupIDXML
	Order        orderXML
}

type orderGroupIDXML struct {
	XMLName           xml.Name `xml:"OrderGroupID"`
	Firm              types.XMLInt64
	Branch            int
	ClientKey         string
	AccountKey        string
	SubAccountType    string
	ActivityTimestamp string
}

type orderXML struct {
	XMLName                 xml.Name `xml:"Order"`
	OrderKey                string
	Security                securityXML
	OrderPricing            orderPricingXML
	OrderType               string
	OrderDuration           string
	OrderEnteredDateTime    string
	OrderInstructions       string
	OriginalQuantity        float32
	SpecialInstructions     specialInstructionsXML
	DoNotReduceIncreaseFlag string
	Discretionary           bool
	OrderSource             string
	Solicited               bool
	MarketCode              string
	Capacity                string
	GoodTilDate             string
	ActivationPrice         float32
	LastUpdated             string
	OriginalOrderID         int
	PendingCancelQuantity   float32
	CancelledQuantity       float32
	RejectCode              int
	RejectReason            string
	ReportedBy              string
	RemainingQuantity       float32
	OrderCompletionCode     string
	Charges                 chargesXML
	OrderAssociation        orderAssociationXML
	ComplexOrderType        string
	CreditOrDebit           string
	ExecutionInformation    executionInformationXML
}

type securityXML struct {
	XMLName          xml.Name `xml:"Security"`
	CUSIP            string
	Symbol           string
	SecurityType     string
	SecurityCategory string
	ShortDescription string
	SymbolUnderlying string
}

type orderPricingXML struct {
	XMLName xml.Name `xml:"OrderPricing"`
	Last    float32
	Ask     float32
	Bid     float32
	Limit   float32
	Method  string
	Amount  float32
}

type specialInstructionsXML struct {
	XMLName   xml.Name `xml:"SpecialInstructions"`
	AllOrNone int
}

type chargesXML struct {
	XMLName xml.Name `xml:"Charges"`
	Charge  []chargeXML
}

type chargeXML struct {
	XMLName xml.Name `xml:"Charge"`
	Type    string
	Amount  float32
}

type orderAssociationXML struct {
	XMLName xml.Name `xml:"OrderAssociation"`
	Type    typeXML
}

type typeXML struct {
	XMLName          xml.Name `xml:"Type"`
	AssociatedOrders associatedOrdersXML
}

type associatedOrdersXML struct {
	XMLName      xml.Name `xml:"AssociatedOrders"`
	OrderKey     string
	Relationship string
}

type executionInformationXML struct {
	XMLName               xml.Name `xml:"ExecutionInformation"`
	Type                  string
	Timestamp             string
	Quantity              float32
	ExecutionPrice        float32
	AveragePriceIndicator bool
	LeavesQuantity        float32
	ID                    string
}
