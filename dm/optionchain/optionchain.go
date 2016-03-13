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

package optionchain

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/marklaczynski/acidbath/dm/optionchain/option"
	"github.com/marklaczynski/acidbath/dm/optionchain/optiondate"
	"github.com/marklaczynski/acidbath/dm/optionchain/optionstrike"
	"github.com/marklaczynski/acidbath/lib/date"
	"github.com/marklaczynski/acidbath/lib/mjlog"
)

var (
	logInfo  = log.New(mjlog.CreateInfoFile(), "INFO  [optionchain]: ", log.LstdFlags|log.Lshortfile)
	logDebug = log.New(mjlog.CreateDebugFile(), "DEBUG [optionchain]: ", log.LstdFlags|log.Lshortfile)
	logError = log.New(mjlog.CreateErrorFile(), "ERROR [optionchain]: ", log.LstdFlags|log.Lshortfile)
)

//OptionChain is an implemenation of optionchain interface
type OptionChain struct {
	sync.RWMutex
	expirations map[time.Time]*optiondate.OptionDate
	underlying  string //key
}

//SortedExpirations retuns a slice of OptionDates that is sorted from most recent to furthest out
func (oc *OptionChain) SortedExpirations() []*optiondate.OptionDate {
	oc.RLock()
	defer oc.RUnlock()
	var keys date.Dates

	for k := range oc.expirations {
		keys = append(keys, k)
	}
	sort.Sort(keys)

	var sortedExpirations []*optiondate.OptionDate
	for _, k := range keys {
		sortedExpirations = append(sortedExpirations, oc.expirations[k])
	}

	return sortedExpirations
}

//OptionSymbols returns a slice of string which contains all the symbols for the current options chain's options
func (oc *OptionChain) OptionSymbols() []string {
	oc.RLock()
	defer oc.RUnlock()

	options := make([]string, 0, 0)
	for dateKey := range oc.expirations {
		for _, strike := range oc.expirations[dateKey].Strikes() {
			if strike.Option(option.CALL) != nil {
				options = append(options, strike.Option(option.CALL).OptionTickerSymbol())
			}
			if strike.Option(option.PUT) != nil {
				options = append(options, strike.Option(option.PUT).OptionTickerSymbol())
			}
		}
	}

	return options
}

//NewOptionChain returns a pointer to a new OptionChain
func NewOptionChain(ul string) *OptionChain {
	x := make(map[time.Time]*optiondate.OptionDate)

	return &OptionChain{expirations: x, underlying: ul}
}

func (oc *OptionChain) addOptionDates(expirationsData ...*optiondate.OptionDate) error {
	oc.Lock()
	defer oc.Unlock()

	for _, v := range expirationsData {
		if oc.expirations[v.Date()] == nil {
			oc.expirations[v.Date()] = v
		} else {
			return fmt.Errorf("Date %v already exists", v.Date())
		}
	}
	return nil
}

//NewOption returns a new option structure initialized to the values provided in params
func (oc *OptionChain) NewOption(ul string, strike float64, expDate time.Time, optType option.TypeOfOption, multiplier float64) (*option.Option, error) {
	return option.NewOption(ul, strike, expDate, optType, multiplier)
}

//Option return a pointer to the option in the OptionChain based on key paramters
func (oc *OptionChain) Option(optTickerSymbol string) *option.Option {
	oc.RLock()
	defer oc.RUnlock()

	for dateKey := range oc.expirations {
		for _, strike := range oc.expirations[dateKey].Strikes() {
			if strike.Option(option.CALL) != nil {
				if strike.Option(option.CALL).OptionTickerSymbol() == optTickerSymbol {
					return strike.Option(option.CALL)
				}
			}
			if strike.Option(option.PUT) != nil {
				if strike.Option(option.PUT).OptionTickerSymbol() == optTickerSymbol {
					return strike.Option(option.PUT)
				}
			}
		}
	}
	return nil
}

//AddOption adds an option to the option chain. If the strike or exp date do not exist, then it will create them. It will return an error if the option already exists
func (oc *OptionChain) AddOption(o *option.Option) error {
	//Sanity checking
	if o.Underlying() != oc.Underlying() {
		return errors.New("Incompatable underlyings. Option U/L: " + o.Underlying() + "Option Chain U/L: " + oc.Underlying())
	}

	// retrieve the od (OptionDate) for the option's exp date
	od := oc.getOptionDate(o.ExpirationDate())
	if od == nil {
		// if it doesn't exist, add it
		tmpod := optiondate.NewOptionDate(o.ExpirationDate())
		if tmpod == nil {
			return errors.New("Creating a new Option Date failed")
		}

		tmpod.SetDaysToExpiration(o.DaysToExpiration())
		oc.addOptionDates(tmpod)
		od = oc.getOptionDate(o.ExpirationDate())
	}

	// retrieve the os (OptionStrike) for the option's strike
	os := od.GetStrike(o.Strike())
	if os == nil {
		// if it doesn't exist, add it
		tmpos, err := optionstrike.NewOptionStrike(o.Strike())
		if err != nil {
			return errors.New("Creating a new strike failed because of error: " + err.Error())
		}
		od.AddStrikes(tmpos)
		os = od.GetStrike(o.Strike())
	}

	// check to make sure there's no option already there, you should not overwrite anything present
	if os.Option(o.OptionType()) != nil && os.Option(o.OptionType()).Multiplier() == o.Multiplier() {
		return fmt.Errorf("Option already exists %#v", o)
	}

	//finally add the option
	os.AddOption(o)
	return nil
}

//Underlying returns the underlying stock/equity that the option chain represents
func (oc *OptionChain) Underlying() string {
	return oc.underlying
}

//SetUnderlying sets the underlying stock/equity that the option chain represents
func (oc *OptionChain) SetUnderlying(underlying string) {
	oc.underlying = underlying
}

//getOptionDate returns the pointer to expiration date otherwise nil
func (oc *OptionChain) getOptionDate(date time.Time) *optiondate.OptionDate {
	oc.RLock()
	defer oc.RUnlock()

	if oc.expirations[date] != nil {
		return oc.expirations[date]
	}
	return nil
}

//FindExpirationClosestTo returns an OptionDate closest to days param
func (oc *OptionChain) FindExpirationBetween(minDTE, maxDTE int64) *optiondate.OptionDate {
	for _, currExpiration := range oc.SortedExpirations() {
		logDebug.Printf("Days to exp %d", currExpiration.DaysToExpiration())
		if currExpiration.DaysToExpiration() > minDTE && currExpiration.DaysToExpiration() < maxDTE {
			return currExpiration
		}
	}

	return nil
}
