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

package generic

import (
	"github.com/marklaczynski/acidbath/dm/optionchain/option"
	"github.com/marklaczynski/acidbath/dm/orderstatus"
)

/*
Notes:

FUTURE: add more params such as stock, so that you can stream stock data as input for options (if need be)
AddOption(*option.Option, stock, etc)
*/

//Strategy interface
type Strategy interface {
	Execute(optionTicker string)
	AddOption(o *option.Option)
	RemoveOption(o *option.Option)
	TrackedOptions() []string
	Option(optionTicker string) (o *option.Option, ok bool)
	OrderStatusUpdate(orderStatus *orderstatus.OrderStatus)
}
