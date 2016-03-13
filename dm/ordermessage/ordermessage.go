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

package ordermessage

import "github.com/marklaczynski/acidbath/lib/orderconst"

type Message struct {
	orderID    string
	orderEvent orderconst.OrderEvent
}

func New(orderid string, orderevent orderconst.OrderEvent) *Message {
	return &Message{
		orderID:    orderid,
		orderEvent: orderevent,
	}
}

//OrderEvent returns the last event to occur for this order
func (m *Message) OrderEvent() orderconst.OrderEvent {
	return m.orderEvent
}

//SetOrderEvent sets the last event to occur for this order
func (m *Message) SetOrderEvent(event orderconst.OrderEvent) {
	m.orderEvent = event
}

//SetOrderID sets the broker's order id
func (m *Message) SetOrderID(id string) {
	m.orderID = id
}

//OrderID returns the brokerage's order id
func (m *Message) OrderID() string {
	return m.orderID
}

func (m *Message) Copy() *Message {
	return &Message{
		orderID:    m.orderID,
		orderEvent: m.orderEvent,
	}
}
