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
	"testing"
	"time"
)

func TestNewOptionChain(t *testing.T) {
	testOptionChain := NewOptionChain("XYZ")
	if testOptionChain == nil {
		t.Errorf("testOptionChain was not created")
	}
}

func todayDate(t *testing.T) time.Time {
	d, err := date.NewDate(time.Now())
	if err != nil {
		t.Errorf("Error setting up test")
		t.FailNow()
	}
	return d
}

func old(t *testing.T) {
	/*
		cases := []struct {
			ul      string
			strike  float32
			exp     time.Time
			optType asset.TypeOfOption
		}{
			{"XYZ", 1, todayDate(t), asset.PUT},
			{"XYZ", 1, todayDate(t), asset.PUT},
		}

		oc := NewOptionChain("XYZ")

		for _, v := range cases {
			o, err := asset.NewOption(v.ul, v.strike, v.exp, v.optType)
			if err != nil {
				t.Errorf("Error setting up test")
				t.FailNow()
			}
		}
	*/
}

func TestAddOptions(t *testing.T) {
	ul := "XYZ"
	oc := NewOptionChain(ul)

	var strike float32
	strike = 1.0
	exp := todayDate(t)
	oType := asset.CALL

	o, err := asset.NewOption(ul, strike, exp, oType)
	if err != nil {
		t.Errorf("Error setting up test")
		t.FailNow()
	}

	if err := oc.AddOption(o); err != nil {
		t.Errorf("Error adding first option")
	}

	//t.Logf("OC:\n%#v\n", oc)
	//t.Logf("OC:\n%#v\n", oc.expirations)

	if err := oc.AddOption(o); err == nil {
		t.Errorf("Incorrectly added the same option twice")
	}

	// should not be able to add an option for ABC to option chain for XYZ
	o.SetUnderlying("ABC")
	if err := oc.AddOption(o); err == nil {
		t.Errorf("Incorrectly added option with different underlying than the chain")
	}

}
