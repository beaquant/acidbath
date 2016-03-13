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

package order

import (
	"math/big"

	"github.com/marklaczynski/acidbath/lib/financial"
	"github.com/marklaczynski/acidbath/lib/orderconst"
)

//Order represents an order to be sent to brokerage firm
type Order struct {
	//accountId req, but specific to td
	clientOrderID       string                              //opt
	orderID             string                              //blank for new order, required for cancel and edits
	action              orderconst.OrderAction              //req enum: buytoopen, buytoclose, selltoopen, selltoclose
	activatePrice       financial.Money                     //opt (stop price)
	expire              orderconst.OrderExpiry              //req enum: day, gtc
	exDay               orderconst.OrderInt8                //opt Two digit expiration day, only specified if expire is set to gtc otherwise null.
	exMonth             orderconst.OrderInt8                //opt Two digit expiration month, only specified if expire is set to gtc otherwise null.
	exYear              orderconst.OrderInt8                //opt Two digit expiration year, only specified if expire is set to gtc otherwise null. Doc says 2 digits, but i don't believe it
	orderType           orderconst.OrderType                //req enum: market, limit, stop_market, stop_limit
	price               financial.Money                     //opt limit price
	quantity            int                                 //req
	routing             orderconst.OrderExchange            //opt enum: auto, isex, cboe, amex, phlx, pacx, bosx
	specialInstructions orderconst.OrderSpecialInstructions //opt enum: none, fok, aon
	symbol              string                              //req

}

//New returns a pointer to a new Order. Default routing is "Auto"
func New() *Order {
	o := &Order{}

	o.routing = orderconst.Auto
	o.activatePrice.Value = big.NewRat(0, 1)
	o.price.Value = big.NewRat(0, 1)

	return o
}

//Copy will return a copy of the Order structure
func (o *Order) Copy() *Order {
	return &Order{
		clientOrderID:       o.clientOrderID,
		orderID:             o.orderID,
		action:              o.action,
		activatePrice:       o.activatePrice,
		expire:              o.expire,
		exDay:               o.exDay,
		exMonth:             o.exMonth,
		exYear:              o.exYear,
		orderType:           o.orderType,
		price:               o.price,
		quantity:            o.quantity,
		routing:             o.routing,
		specialInstructions: o.specialInstructions,
		symbol:              o.symbol,
	}
}

//ClientOrderID returns the client order id. This is a client created id
func (o *Order) ClientOrderID() string {
	return o.clientOrderID
}

//OrderID returns the brokerage's order id
func (o *Order) OrderID() string {
	return o.orderID
}

//Action returns the order action
func (o *Order) Action() orderconst.OrderAction {
	return o.action
}

//ActivatePrice retuns the activation price
func (o *Order) ActivatePrice() financial.Money {
	return o.activatePrice
}

//Expire returns the expiration type
func (o *Order) Expire() orderconst.OrderExpiry {
	return o.expire
}

//ExDay returns the day the order expires if order is GTC
func (o *Order) ExDay() orderconst.OrderInt8 {
	return o.exDay
}

//ExMonth returns the month the order expires if order is GTC
func (o *Order) ExMonth() orderconst.OrderInt8 {
	return o.exMonth
}

//ExYear returns the year the order expires if order is GTC
func (o *Order) ExYear() orderconst.OrderInt8 {
	return o.exYear
}

//OrderType returns the order type
func (o *Order) OrderType() orderconst.OrderType {
	return o.orderType
}

//Price returns the order price
func (o *Order) Price() financial.Money {
	return o.price
}

//Quantity returns the number of contracts to trade
func (o *Order) Quantity() int {
	return o.quantity
}

//Routing returns the excahge to which orde is set to route to
func (o *Order) Routing() orderconst.OrderExchange {
	return o.routing
}

//SpecialInstructions returns any special instructions associated with the order (fill or kill, all or none)
func (o *Order) SpecialInstructions() orderconst.OrderSpecialInstructions {
	return o.specialInstructions
}

//Symbol retuns the symbol to be traded
func (o *Order) Symbol() string {
	return o.symbol
}

//SetClientOrderID sets the client order id
func (o *Order) SetClientOrderID(id string) {
	o.clientOrderID = id
}

//SetOrderID sets the broker's order id
func (o *Order) SetOrderID(id string) {
	o.orderID = id
}

//SetAction sets the action type on the order
func (o *Order) SetAction(action orderconst.OrderAction) {
	o.action = action
}

//SetActivatePrice sets the activation price on the order
func (o *Order) SetActivatePrice(price financial.Money) {
	o.activatePrice.Value.Set(price.Value)
}

//SetExpire sets the expiration type on the order (GTC/Day)
func (o *Order) SetExpire(expire orderconst.OrderExpiry) {
	o.expire = expire
}

//SetExDay sets the day the order expires if order is GTC
func (o *Order) SetExDay(exDay orderconst.OrderInt8) {
	o.exDay = exDay
}

//SetExMonth sets the month the order expires if order is GTC
func (o *Order) SetExMonth(exMonth orderconst.OrderInt8) {
	o.exMonth = exMonth
}

//SetExYear sets the year the order expires if order is GTC
func (o *Order) SetExYear(exYear orderconst.OrderInt8) {
	o.exYear = exYear
}

//SetOrderType sets the type of order (limit,market,etc)
func (o *Order) SetOrderType(orderType orderconst.OrderType) {
	o.orderType = orderType
}

//SetPrice set the price to execute the order at
func (o *Order) SetPrice(price financial.Money) {
	o.price.Value.Set(price.Value)
}

//SetQuantity sets the number of contracts to trade
func (o *Order) SetQuantity(quantity int) {
	o.quantity = quantity
}

//SetRouting sets the exchange to route the order to. Default is Auto
func (o *Order) SetRouting(routing orderconst.OrderExchange) {
	o.routing = routing
}

//SetSpecialInstructions set any special instructions on the order (ie fill or kill, all or none)
func (o *Order) SetSpecialInstructions(spInstructions orderconst.OrderSpecialInstructions) {
	o.specialInstructions = spInstructions
}

//SetSymbol sets the underlying symbol to trade
func (o *Order) SetSymbol(symbol string) {
	o.symbol = symbol
}
