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

package orderbook

import (
	"fmt"

	"github.com/marklaczynski/acidbath/dm/orderstatus"
)

type OrderBook struct {
	orderStatus map[string]*orderstatus.OrderStatus
}

func New() *OrderBook {
	return &OrderBook{
		orderStatus: make(map[string]*orderstatus.OrderStatus),
	}
}

func (ob *OrderBook) OrderStatuses() map[string]*orderstatus.OrderStatus {
	return ob.orderStatus
}

func (ob *OrderBook) OrderStatus(orderid string) *orderstatus.OrderStatus {
	return ob.orderStatus[orderid]
}

func (ob *OrderBook) AddUpdateOrderStatus(newOrderStatus *orderstatus.OrderStatus) {
	ob.orderStatus[newOrderStatus.OrderID()] = newOrderStatus
}

func (ob *OrderBook) DeleteOrderStatus(existingOrder *orderstatus.OrderStatus) {
	delete(ob.orderStatus, existingOrder.OrderID())
}

func (ob *OrderBook) String() string {
	var finalString string
	finalString = "================== Order Book Start ================== \n"
	for _, currOrderStatus := range ob.orderStatus {
		finalString = finalString + fmt.Sprintf("%s\n", currOrderStatus)
	}
	finalString = finalString + "================== Order Book End ================== \n"
	return finalString
}
