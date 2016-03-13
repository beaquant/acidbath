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

// Package factory is responsible for returning a concrete implemenation of a generic broker interface
package factory

import (
	"github.com/marklaczynski/acidbath/broker/generic"
	"github.com/marklaczynski/acidbath/broker/tdapi"
)

//BrokerType is an enumeration of available brokers
type BrokerType int

const (
	//TD supports the TDAmeritrade broker
	TD BrokerType = iota
)

//CreateBroker returns a concrete instance of generic.Broker interface, based on the BrokerType
func CreateBroker(b BrokerType) generic.Broker {
	switch b {

	case TD:
		return tdapi.New()
	}

	return nil
}
