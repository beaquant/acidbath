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

package optiondate

import (
	"testing"
	"time"

	"github.com/marklaczynski/acidbath/dm/optionchain/optionstrike"
)

func TestNewOptionDate(t *testing.T) {
	initDate, err := date.NewDate(time.Now())
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
		t.FailNow()
	}

	od, err := NewOptionDate(initDate)
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
	}

	if od == nil {
		t.Errorf("Option Date was not created")
	}
}

func utilNewOptionStrike(strike float32) *optionstrike.OptionStrike {
	os, _ := optionstrike.NewOptionStrike(strike)
	return os
}

func TestSetDaysToExpiration(t *testing.T) {
	initDate, err := date.NewDate(time.Now())
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
		t.FailNow()
	}

	od, err := NewOptionDate(initDate)
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
	}

	cases := []struct {
		dte int32
	}{
		{-100},
		{-10},
		{0},
		{10},
		{100},
	}

	for _, v := range cases {
		err := od.SetDaysToExpiration(v.dte)
		if v.dte < 0 && err == nil {
			t.Errorf("DaysToExpiration cannot be less than 0: %v", v.dte)
		}

		if v.dte >= 0 && od.DaysToExpiration() != v.dte {
			t.Errorf("DaysToExpiration is not set to expected value.\n IP:%v\nOP%v\n", v.dte, od.DaysToExpiration())
		}

	}

}

func TestSetDate(t *testing.T) {
	initDate, err := date.NewDate(time.Now())
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
		t.FailNow()
	}

	od, err := NewOptionDate(initDate)
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
	}

	currDate, err := date.NewDate(time.Now())
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
		t.FailNow()
	}

	futDate, err := date.NewDate(currDate.AddDate(1, 0, 0))
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
		t.FailNow()
	}

	pastDate, err := time.Parse("2006-January-02", "2014-January-01")
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
		t.FailNow()
	}

	cases := []struct {
		d time.Time
	}{
		{pastDate}, //past date
		{currDate}, //current date
		{futDate},  // some future date
	}

	for _, v := range cases {
		if err := od.SetDate(v.d); err == nil && od.Date().Before(currDate) {
			t.Errorf("Expiration date is in the past .\nIP Date: %v\nCP Date: %v\n", v.d, currDate)
		}

	}
}

func TestAddStrikes(t *testing.T) {
	initDate, err := date.NewDate(time.Now())
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
		t.FailNow()
	}

	od, err := NewOptionDate(initDate)
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
	}

	/*  interesting read, as optionStrike.NewOptionStrike returns 2 values
	http://blog.vladimirvivien.com/2014/03/hacking-go-filter-values-from-multi.html
	*/

	cases := []struct {
		strike *optionstrike.OptionStrike
	}{
		// cannot do this {optionstrike.NewOptionStrike(-100.5)}, see note above
		{utilNewOptionStrike(10)},
		{utilNewOptionStrike(100.5)},
	}

	for _, v := range cases {
		od.AddStrikes(v.strike)
	}

	if len(od.strikes) != len(cases) {
		t.Errorf("Did not add all strikes. \nActual: %v\nExpected: ", len(od.strikes), len(cases))

	}

	//Test scenario 2, add multiple strikes in single call
	//Reset strikes
	od, err = NewOptionDate(initDate)
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
	}

	od.AddStrikes(cases[0].strike, cases[1].strike)

	if len(od.strikes) != 2 {
		t.Errorf("Did not add all strikes. \nActual: %v\nExpected: %v", len(od.strikes), 2)

	}

	//Test scenario 3, should not be able to add same strike twice
	od, err = NewOptionDate(initDate)
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
	}

	os1 := utilNewOptionStrike(1)
	os2 := utilNewOptionStrike(1)

	if err = od.AddStrikes(os1, os2); err == nil {
		t.Errorf("Added the same strike twice")
	}

}

// test being able to find a strike
func TestStrike(t *testing.T) {
	initDate, err := date.NewDate(time.Now())
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
		t.FailNow()
	}

	od, err := NewOptionDate(initDate)
	if err != nil {
		t.Errorf("Error setting up test. Error: %v", err)
	}

	if s := od.GetStrike(1); s != nil {
		t.Errorf("Strike function not returning nil for non existant strike")
	}

	os1 := utilNewOptionStrike(1)
	os2 := utilNewOptionStrike(2)
	os3 := utilNewOptionStrike(3)

	od.AddStrikes(os1, os2, os3)

	if s := od.GetStrike(1); s == nil {
		t.Errorf("Strike function not able to find existing strike")
	}

}
