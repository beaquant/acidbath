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
	"github.com/marklaczynski/acidbath/eventproc/generic"
	"github.com/marklaczynski/acidbath/eventproc/reference"
)

//StrategyType is an enumeration of available brokers
type StrategyType int

//This is a list of all the available strategies
const (
	Reference StrategyType = iota
	Count                  // not a strategy, used to initiate an array of strategies
)

//CreateStrategy returns a concrete instance of generic.Broker interface, based on the StrategyType
func CreateStrategy(s StrategyType) generic.Strategy {
	switch s {

	case Reference:
		return reference.New()

	}

	return nil
}
