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

package orderconst

//OrderInt8 is used to return string value "" if the value is 0, otherwise it returns the string value of the int
type OrderInt8 int8

//String returns the string of v. if the value is 0 return "", otherwise it returns the string value of the int
func (v OrderInt8) String() string {
	if v == 0 {
		return ""
	}
	return string(v)
}

//OrderSpecialInstructions define special instruction for order, such as fill or kill, all or none
type OrderSpecialInstructions int

//enumerations for OrderSpecialInstructions
const (
	None OrderSpecialInstructions = iota //none
	Fok                                  //fill or kill
	Aon                                  // all or none
)

func (osi OrderSpecialInstructions) String() string {
	switch osi {
	case None:
		return "none"
	case Fok:
		return "fok"
	case Aon:
		return "aon"
	}
	return ""
}

//OrderExchange defines type to handle exchanges to which to send the order to
type OrderExchange int

//enumerations for OrderExchange
const (
	Auto OrderExchange = iota
	ISEX
	CBOE
	AMEX
	PHLX
	PACX
	BOSX
)

func (oe OrderExchange) String() string {
	switch oe {
	case Auto:
		return "auto"
	case ISEX:
		return "isex"
	case CBOE:
		return "cboe"
	case AMEX:
		return "amex"
	case PHLX:
		return "phlx"
	case PACX:
		return "pacx"
	case BOSX:
		return "bosx"
	}
	return ""
}

//OrderExpiry defines expiration types of an order (day/gtc)
type OrderExpiry int

//enumeration values for OrderExpiry
const (
	InvalidOrderExpiry OrderExpiry = iota
	Day
	GTC
)

func (oe OrderExpiry) String() string {
	switch oe {
	case InvalidOrderExpiry:
		return "invalid"
	case Day:
		return "day"
	case GTC:
		return "gtc"
	}
	return ""
}

//OrderAction defines actions available for an order (buytoopen/selltoopen/etc)
type OrderAction int

//enumeration values for OrderAction
const (
	InvalidOrderAction OrderAction = iota
	BuyToOpen
	BuyToClose
	SellToOpen
	SellToClose
)

func (oa OrderAction) String() string {
	switch oa {
	case InvalidOrderAction:
		return "Invalid or Unsupported Order Action"
	case BuyToOpen:
		return "buytoopen"
	case BuyToClose:
		return "buytoclose"
	case SellToOpen:
		return "selltoopen"
	case SellToClose:
		return "selltoclose"
	}
	return ""
}

//OrderType defines type of order (ie Limit/market/etc)
type OrderType int

//enumerations values for OrderType
const (
	InvalidOrderType OrderType = iota
	Limit
	Market
	StopMarket
	StopLimit
)

func (ot OrderType) String() string {
	switch ot {
	case InvalidOrderType:
		return "Invalid or Unsupported Order Type"
	case Market:
		return "market"
	case Limit:
		return "limit"
	case StopMarket:
		return "stop_market"
	case StopLimit:
		return "stop_limit"
	}

	return ""
}

//OrderEvent represents an Order Event that arrives from brokerage
type OrderEvent int

//Order event values
const (
	OrderEventInvalid OrderEvent = iota
	OrderEventNil
	OrderBroken          // After an order was filled, the trade is reversed or "Broken" and the order is changed to Canceled.
	OrderManualExecution // The order is manually entered (and filled) by the broker.  Usually due to some system issue.
	OrderActivation      // A Stop order has been Activated
	OrderCancelReplace   // A request to modify an order (Cancel/Replace) has been received (You will also get a UROUT for the original order)
	OrderCancel          // A request to cancel an order has been received
	OrderEntry           // A new order has been submitted
	OrderFill            // An order has been completely filled
	OrderPartialFill     // An order has been partial filled
	OrderRejection       // An order was rejected
	OrderTooLateToCancel // A request to cancel an order has been received but the order cannot be canceled either because it was already canceled, filled, or for some other reason
	OrderOut             // Indicates "You Are Out" - that the order has been canceled
)

func (oe OrderEvent) String() string {
	switch oe {
	case OrderEventNil:
		return "nil"
	case OrderEventInvalid:
		return "Invalid or Unsupported Order Event"
	case OrderBroken:
		return "OrderBroken"
	case OrderManualExecution:
		return "OrderManualExecution"
	case OrderActivation:
		return "OrderActivation"
	case OrderCancelReplace:
		return "OrderCancelReplace"
	case OrderCancel:
		return "OrderCancel"
	case OrderEntry:
		return "OrderEntry"
	case OrderFill:
		return "OrderFill"
	case OrderPartialFill:
		return "OrderPartialFill"
	case OrderRejection:
		return "OrderRejection"
	case OrderTooLateToCancel:
		return "OrderTooLateToCancel"
	case OrderOut:
		return "OrderOut"
	}
	return ""
}
