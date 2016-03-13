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

package reference

import (
	"log"

	"github.com/marklaczynski/acidbath/dm/optionchain/option"
	"github.com/marklaczynski/acidbath/dm/orderstatus"
	"github.com/marklaczynski/acidbath/lib/mjlog"
)

var (
	logInfo  = log.New(mjlog.CreateInfoFile(), "INFO  [reference processor]: ", log.LstdFlags|log.Lshortfile)
	logDebug = log.New(mjlog.CreateDebugFile(), "DEBUG [reference processor]: ", log.LstdFlags|log.Lshortfile)
	logError = log.New(mjlog.CreateErrorFile(), "ERROR [reference processor]: ", log.LstdFlags|log.Lshortfile)
)

//Reference implementation of a strategy. Do not use in PROD, as it's only for educational purposes
type Reference struct {
}

//New creates and returns a pointer to a new Reference implementation
func New() *Reference {
	return &Reference{}
}

//OrderStatusUpdate handler for order status updates
func (r *Reference) OrderStatusUpdate(orderStatus *orderstatus.OrderStatus) {
	logInfo.Printf("Reference OrderStatusUpdate\n")
}

//Execute runs the strategy on the instrument it's attached to
func (r *Reference) Execute(optionTicker string) {
	logDebug.Printf("Executing reference strategy\n")
}

//AddOption adds the option to be monitored by this strategy
func (r *Reference) AddOption(o *option.Option) {
	logInfo.Printf("AddOption\n")
}

//RemoveOption removes the option from monitoring by this strategy
func (r *Reference) RemoveOption(o *option.Option) {
	logInfo.Printf("RemoveOption\n")
}

//TrackedOptions retuns an array of options ticker symbols being tracked by this strategy
func (r *Reference) TrackedOptions() []string {
	return []string{}
}

//Option retuns the option (if it exists) based on the optionTicker for quick access
func (r *Reference) Option(optionTicker string) (o *option.Option, ok bool) {
	return nil, false
}
