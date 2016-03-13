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

package optionstrike

import (
	"testing"
	"time"

	"github.com/marklaczynski/acidbath/dm/optionchain/option"
)

// This test the creation of an Option Strike with valid values
func TestNewOptionStrike(t *testing.T) {
	tos, err := NewOptionStrike(100)
	if tos == nil && err != nil {
		t.Errorf("Option Strike was not created")
	}

	if _, err := NewOptionStrike(0); err == nil {
		t.Errorf("Option Strike did not return error for a 0 strike")
	}

}

func TestStrike(t *testing.T) {
	tos, _ := NewOptionStrike(1)

	cases := []struct {
		strike float32
	}{
		{-100.5},
		{-10},
		{0},
		{10},
		{100.5},
	}

	for _, v := range cases {

		err := tos.SetStrike(v.strike)
		if v.strike <= 0 && err == nil {
			t.Errorf("Did not receive error for non-positive strike %v\n", v.strike)
		}

		if v.strike > 0 && tos.Strike() != v.strike {
			t.Errorf("Did not get correct strike back.\nIP: %s\nOP: %s\n", v.strike, tos.Strike())
		}

	}

}

func newOpt(t *testing.T) *option.Option {
	today, err := date.NewDate(time.Now())
	if err != nil {
		t.Errorf("Error setting up test option's date")
		t.FailNow()
	}

	op, err := option.NewOption("XYZ", 1, today, option.CALL)
	if err != nil {
		t.Errorf("Error setting up test option")
		t.FailNow()
	}
	return op
}

func TestAddOption(t *testing.T) {
	tos, _ := NewOptionStrike(1)

	cases := []struct {
		optType option.TypeOfOption
		ticker  string
	}{
		{option.PUT, "XYZ150220P175"},
		{option.CALL, "ABC150220C175"},
	}

	for _, v := range cases {
		testOption := newOpt(t)
		testOption.SetOptionType(v.optType)
		testOption.SetOptionTickerSymbol(v.ticker)

		tos.AddOption(testOption)

		if tos.Option(testOption.OptionType()).OptionTickerSymbol() != v.ticker {
			t.Errorf("Did not get correct option back.\nIP: %s\nOP: %s\n", v.ticker, tos.Option(testOption.OptionType()).OptionTickerSymbol())
		}

	}

}
