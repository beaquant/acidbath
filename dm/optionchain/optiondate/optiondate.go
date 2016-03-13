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

//Package optiondate implements an option expirationDate
package optiondate

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/marklaczynski/acidbath/dm/optionchain/optionstrike"
	"github.com/marklaczynski/acidbath/lib/date"
)

//As per http://golang.org/pkg/time/ time struct should be passed by value

//OptionDate representa an expiration date within an option chain
type OptionDate struct {
	expirationDate   time.Time //key
	strikes          map[float64]*optionstrike.OptionStrike
	daysToExpiration int64
	//expirationType TBD enum
}

//NewOptionDate instantiate an OptionDate with a date
func NewOptionDate(date time.Time) *OptionDate {
	sk := make(map[float64]*optionstrike.OptionStrike)
	od := &OptionDate{strikes: sk}
	if err := od.SetDate(date); err != nil {
		return nil
	}
	return od
}

//Date returns the current expiration
func (od *OptionDate) Date() time.Time {
	return od.expirationDate
}

//SetDate will set the date and ignore any timestamp.
func (od *OptionDate) SetDate(expirationDate time.Time) error {
	expDate, err := date.New(expirationDate)
	if err != nil {
		return err
	}

	od.expirationDate = expDate
	return nil
}

//AddStrikes add a number of strikes to the current option expiration date
func (od *OptionDate) AddStrikes(strike ...*optionstrike.OptionStrike) error {
	// http://blog.golang.org/go-maps-in-action
	for _, v := range strike {
		if od.strikes[v.Strike()] == nil {
			od.strikes[v.Strike()] = v
		} else {

			return fmt.Errorf("Strike %f already exists", v.Strike())
		}
	}
	return nil
}

//GetStrike returns a strike object that represents the strike parameter value
func (od *OptionDate) GetStrike(strike float64) *optionstrike.OptionStrike {
	if od.strikes[strike] != nil {
		return od.strikes[strike]
	}

	return nil
}

//Strikes returns all strikes as map
func (od *OptionDate) Strikes() map[float64]*optionstrike.OptionStrike {
	return od.strikes
}

func (od *OptionDate) SortedStrikes() []*optionstrike.OptionStrike {
	var keys sort.Float64Slice

	for k := range od.strikes {
		keys = append(keys, k)
	}

	sort.Sort(keys)

	var sortedStrikes []*optionstrike.OptionStrike
	for _, k := range keys {
		sortedStrikes = append(sortedStrikes, od.GetStrike(k))
	}

	return sortedStrikes
}

//DaysToExpiration returns the number of days to expiration for current expiration date
func (od *OptionDate) DaysToExpiration() int64 {
	return od.daysToExpiration
}

//SetDaysToExpiration sets the number of days to expiration for current expiration date
func (od *OptionDate) SetDaysToExpiration(daysToExpiration int64) error {
	if daysToExpiration < 0 {
		return errors.New("daysToExpiration cannot be < 0")
	}
	od.daysToExpiration = daysToExpiration
	return nil
}
