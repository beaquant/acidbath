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

import "encoding/xml"

//CancelOrderMessage represents the cancled TD Order structure returned after sending an order
type CancelOrderMessage struct {
	XMLName             xml.Name `xml:"amtd"`
	Error                        //inline struct
	CancelOrderMessages cancelOrderMessagesXML
}

type cancelOrderMessagesXML struct {
	XMLName       xml.Name           `xml:"cancel-order-messages"`
	AccountID     string             `xml:"account-id"`
	CanceledOrder []canceledOrderXML `xml:"order"`
}

type canceledOrderXML struct {
	XMLName xml.Name `xml:"order"`
	OrderID string   `xml:"order-id"`
	Message string   `xml:"message"`
	Error   string   `xml:"error"`
}
