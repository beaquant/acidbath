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

package orderstatus

import (
	"fmt"
	"math/big"

	"github.com/marklaczynski/acidbath/lib/financial"
	"github.com/marklaczynski/acidbath/lib/orderconst"
)

//OrderStatus represents a order status received from brokerage firm
type OrderStatus struct {
	//activatePrice       financial.Money                     //opt (stop price)
	//exDay               orderconst.OrderInt8                //opt Two digit expiration day, only specified if expire is set to gtc otherwise null.
	//exMonth             orderconst.OrderInt8                //opt Two digit expiration month, only specified if expire is set to gtc otherwise null.
	//exYear              orderconst.OrderInt8                //opt Two digit expiration year, only specified if expire is set to gtc otherwise null. Doc says 2 digits, but i don't believe it
	//specialInstructions orderconst.OrderSpecialInstructions //opt enum: none, fok, aon

	status    string
	orderID   string                   // KEY ... blank for new order, required for cancel and edits
	action    orderconst.OrderAction   //req enum: buytoopen, buytoclose, selltoopen, selltoclose
	orderType orderconst.OrderType     //req enum: market, limit, stop_market, stop_limit
	quantity  int                      //req
	price     financial.Money          //opt limit price
	symbol    string                   //req
	expire    orderconst.OrderExpiry   //req enum: day, gtc
	routing   orderconst.OrderExchange //opt enum: auto, isex, cboe, amex, phlx, pacx, bosx

	//orderEvent orderconst.OrderEvent
}

//New returns a pointer to a new OrderStatus. Default routing is "Auto"
func New() *OrderStatus {
	o := &OrderStatus{}

	o.routing = orderconst.Auto
	//o.activatePrice.Value = big.NewRat(0, 1)
	o.price.Value = big.NewRat(0, 1)

	return o
}

func (o *OrderStatus) Status() string {
	return o.status
}

func (o *OrderStatus) SetStatus(newStatus string) {
	o.status = newStatus
}

func (o *OrderStatus) String() string {
	//tmpString1 := fmt.Sprintf("\nOrderID: %s\n\tSymbol: %s\n\tAction: %s\n\tActivatePrice: %s\n\tOrderStatus: %s\n\tOrderEvent: %s\n\t", o.orderID, o.symbol, o.action, o.activatePrice, o.status, o.orderEvent)
	tmpString1 := fmt.Sprintf("\nOrderID: %s\n\tSymbol: %s\n\tAction: %s\n\tOrderStatus: %s\n\t", o.orderID, o.symbol, o.action, o.status)
	tmpString2 := fmt.Sprintf("Price: %s\n\tQuantity: %d\n\tOrderType: %s\n\t", o.price, o.quantity, o.orderType)
	//tmpString3 := fmt.Sprintf("Exipration: %s\n\tExDay: %s\n\tExMonth: %s\n\tExYear: %s\n\tRouting: %s\n\t", o.expire, o.exDay, o.exMonth, o.exYear, o.routing)
	tmpString3 := fmt.Sprintf("Exipration: %s\n\tRouting: %s\n\t", o.expire, o.routing)
	return tmpString1 + tmpString2 + tmpString3
}

//Copy will return a copy of the order status
func (o *OrderStatus) Copy() *OrderStatus {
	return &OrderStatus{
		orderID: o.orderID,
		action:  o.action,
		//activatePrice:       o.activatePrice,
		expire: o.expire,
		//exDay:               o.exDay,
		//exMonth:             o.exMonth,
		//exYear:              o.exYear,
		orderType: o.orderType,
		price:     o.price,
		quantity:  o.quantity,
		routing:   o.routing,
		//specialInstructions: o.specialInstructions,
		symbol: o.symbol,
		status: o.status,
		//orderEvent: o.orderEvent,
	}
}

//OrderID returns the brokerage's order id
func (o *OrderStatus) OrderID() string {
	return o.orderID
}

//Action returns the order action
func (o *OrderStatus) Action() orderconst.OrderAction {
	return o.action
}

//Expire returns the expiration type
func (o *OrderStatus) Expire() orderconst.OrderExpiry {
	return o.expire
}

//OrderType returns the order type
func (o *OrderStatus) OrderType() orderconst.OrderType {
	return o.orderType
}

//Price returns the order price
func (o *OrderStatus) Price() financial.Money {
	return o.price
}

//Quantity returns the number of contracts to trade
func (o *OrderStatus) Quantity() int {
	return o.quantity
}

//Routing returns the excahge to which orde is set to route to
func (o *OrderStatus) Routing() orderconst.OrderExchange {
	return o.routing
}

//Symbol retuns the symbol to be traded
func (o *OrderStatus) Symbol() string {
	return o.symbol
}

//SetOrderID sets the broker's order id
func (o *OrderStatus) SetOrderID(id string) {
	o.orderID = id
}

//SetAction sets the action type on the order
func (o *OrderStatus) SetAction(action orderconst.OrderAction) {
	o.action = action
}

//SetExpire sets the expiration type on the order (GTC/Day)
func (o *OrderStatus) SetExpire(expire orderconst.OrderExpiry) {
	o.expire = expire
}

//SetOrderType sets the type of order (limit,market,etc)
func (o *OrderStatus) SetOrderType(orderType orderconst.OrderType) {
	o.orderType = orderType
}

//SetPrice set the price to execute the order at
func (o *OrderStatus) SetPrice(price financial.Money) {
	o.price.Value.Set(price.Value)
}

//SetQuantity sets the number of contracts to trade
func (o *OrderStatus) SetQuantity(quantity int) {
	o.quantity = quantity
}

//SetRouting sets the exchange to route the order to. Default is Auto
func (o *OrderStatus) SetRouting(routing orderconst.OrderExchange) {
	o.routing = routing
}

//SetSymbol sets the underlying symbol to trade
func (o *OrderStatus) SetSymbol(symbol string) {
	o.symbol = symbol
}

/*
//ActivatePrice retuns the activation price
func (o *OrderStatus) ActivatePrice() financial.Money {
	return o.activatePrice
}

//SetActivatePrice sets the activation price on the order
func (o *OrderStatus) SetActivatePrice(price financial.Money) {
	o.activatePrice.Value.Set(price.Value)
}

//ExDay returns the day the order expires if order is GTC
func (o *OrderStatus) ExDay() orderconst.OrderInt8 {
	return o.exDay
}

//ExMonth returns the month the order expires if order is GTC
func (o *OrderStatus) ExMonth() orderconst.OrderInt8 {
	return o.exMonth
}

//ExYear returns the year the order expires if order is GTC
func (o *OrderStatus) ExYear() orderconst.OrderInt8 {
	return o.exYear
}

//SetExDay sets the day the order expires if order is GTC
func (o *OrderStatus) SetExDay(exDay orderconst.OrderInt8) {
	o.exDay = exDay
}

//SetExMonth sets the month the order expires if order is GTC
func (o *OrderStatus) SetExMonth(exMonth orderconst.OrderInt8) {
	o.exMonth = exMonth
}

//SetExYear sets the year the order expires if order is GTC
func (o *OrderStatus) SetExYear(exYear orderconst.OrderInt8) {
	o.exYear = exYear
}

//SpecialInstructions returns any special instructions associated with the order (fill or kill, all or none)
func (o *OrderStatus) SpecialInstructions() orderconst.OrderSpecialInstructions {
	return o.specialInstructions
}

//SetSpecialInstructions set any special instructions on the order (ie fill or kill, all or none)
func (o *OrderStatus) SetSpecialInstructions(spInstructions orderconst.OrderSpecialInstructions) {
	o.specialInstructions = spInstructions
}

//OrderEvent returns the last event to occur for this order
func (o *OrderStatus) OrderEvent() orderconst.OrderEvent {
	return o.orderEvent
}

//SetOrderEvent sets the last event to occur for this order
func (o *OrderStatus) SetOrderEvent(event orderconst.OrderEvent) {
	o.orderEvent = event
}



*/
