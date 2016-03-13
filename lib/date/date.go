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

package date

import (
	"fmt"
	"time"
)

//Dates is a slice of time.Times
type Dates []time.Time

//gives user options for ways to parse a date
const (
	ParseLongFormat    = "2000-01-01 14:30:00 EDT"
	ParseShortFormat   = "2006-January-02"
	ParseOptionExpDate = "20060102"
	LocalLocation      = "America/New_York"
)

// New returns a new time.Time structure with no timestamp
// Date will always use UTC as the locale.
func New(date time.Time) (time.Time, error) {
	year, month, day := date.Date()

	loc, err := time.LoadLocation(LocalLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("Error loading location %s with error %v\n", LocalLocation, err)
	}

	newDate, err := time.ParseInLocation(ParseShortFormat, fmt.Sprintf("%04d-%s-%02d", year, month, day), loc)
	if err != nil {
		return time.Time{}, err
	}

	return newDate.Local(), nil
}

func (d Dates) Len() int {
	return len(d)
}

func (d Dates) Less(i, j int) bool {
	return d[i].Before(d[j])
}

func (d Dates) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
