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

package option

import (
	"testing"
	"time"
)

func newRawOption(t *testing.T) *Option {
	today, err := date.NewDate(time.Now())
	if err != nil {
		t.Errorf("Error setting up test option's date")
		t.FailNow()
	}

	op, err := NewOption("XYZ", 1, today, CALL)
	if err != nil {
		t.Errorf("Error setting up test option")
		t.FailNow()
	}
	return op
}

func TestSetBid(t *testing.T) {
	op := newRawOption(t)

	cases := []struct {
		bid float32
	}{
		{-100.5},
		{-10},
		{0},
		{10},
		{100.5},
	}

	for _, v := range cases {
		err := op.SetBid(v.bid)
		if v.bid < 0 && err == nil {
			t.Errorf("Bid cannot be < 0: %v", v.bid)
		}

		if v.bid >= 0 && op.Bid() != v.bid {
			t.Errorf("Bid is not set to expected value.\n IP:%v\nOP%v\n", v.bid, op.Bid())
		}

	}

}

func TestSetAsk(t *testing.T) {
	op := newRawOption(t)

	cases := []struct {
		ask float32
	}{
		{-100.5},
		{-10},
		{0},
		{10},
		{100.5},
	}

	for _, v := range cases {
		err := op.SetAsk(v.ask)
		if v.ask < 0 && err == nil {
			t.Errorf("Ask cannot be < 0: %v", v.ask)
		}

		if v.ask >= 0 && op.Ask() != v.ask {
			t.Errorf("Ask is not set to expected value.\n IP:%v\nOP%v\n", v.ask, op.Ask())
		}

	}

}

func TestSetBidSize(t *testing.T) {
	op := newRawOption(t)

	cases := []struct {
		bidSize int32
	}{
		{-100},
		{-10},
		{0},
		{10},
		{100},
	}

	for _, v := range cases {
		err := op.SetBidSize(v.bidSize)
		//t.Logf("err: %v, bidSize: %v", err, op.BidSize())
		if v.bidSize <= 0 && err == nil {
			t.Errorf("bidSize cannot be <= 0: %v", v.bidSize)
		}

		if v.bidSize > 0 && op.BidSize() != v.bidSize {
			t.Errorf("bidSize is not set to expected value.\n IP:%v\nOP%v\n", v.bidSize, op.BidSize())
		}

	}

}

func TestSetAskSize(t *testing.T) {
	op := newRawOption(t)

	cases := []struct {
		askSize int32
	}{
		{-100},
		{-10},
		{0},
		{10},
		{100},
	}

	for _, v := range cases {
		err := op.SetAskSize(v.askSize)
		if v.askSize <= 0 && err == nil {
			t.Errorf("AskSize cannot be <= 0: %v", v.askSize)
		}

		if v.askSize > 0 && op.AskSize() != v.askSize {
			t.Errorf("AskSize is not set to expected value.\n IP:%v\nOP%v\n", v.askSize, op.AskSize())
		}

	}

}

func TestSetLast(t *testing.T) {
	op := newRawOption(t)

	cases := []struct {
		last float32
	}{
		{-100.5},
		{-10},
		{0},
		{10},
		{100.5},
	}

	for _, v := range cases {
		err := op.SetLast(v.last)
		if v.last < 0 && err == nil {
			t.Errorf("Last cannot be less than 0: %v", v.last)
		}

		if v.last >= 0 && op.Last() != v.last {
			t.Errorf("Last is not set to expected value.\n IP:%v\nOP%v\n", v.last, op.Last())
		}

	}

}

func TestSetVolume(t *testing.T) {
	op := newRawOption(t)

	cases := []struct {
		volume int64
	}{
		{-100},
		{-10},
		{0},
		{10},
		{100},
	}

	for _, v := range cases {
		err := op.SetVolume(v.volume)
		if v.volume < 0 && err == nil {
			t.Errorf("Volume cannot be less than 0: %v", v.volume)
		}

		if v.volume >= 0 && op.Volume() != v.volume {
			t.Errorf("Volume is not set to expected value.\n IP:%v\nOP%v\n", v.volume, op.Volume())
		}

	}

}

func TestSetOpenInterest(t *testing.T) {
	op := newRawOption(t)

	cases := []struct {
		openInterest int32
	}{
		{-100},
		{-10},
		{0},
		{10},
		{100},
	}

	for _, v := range cases {
		err := op.SetOpenInterest(v.openInterest)
		if v.openInterest < 0 && err == nil {
			t.Errorf("OpenInterest cannot be less than 0: %v", v.openInterest)
		}

		if v.openInterest >= 0 && op.OpenInterest() != v.openInterest {
			t.Errorf("OpenInterest is not set to expected value.\n IP:%v\nOP%v\n", v.openInterest, op.OpenInterest())
		}

	}

}

func TestSetTheoPrice(t *testing.T) {
	op := newRawOption(t)

	cases := []struct {
		theoPrice float32
	}{
		{-100.5},
		{-10},
		{0},
		{10},
		{100.5},
	}

	for _, v := range cases {
		err := op.SetTheoPrice(v.theoPrice)
		if v.theoPrice < 0 && err == nil {
			t.Errorf("TheoPrice cannot be less than 0: %v", v.theoPrice)
		}

		if v.theoPrice >= 0 && op.TheoPrice() != v.theoPrice {
			t.Errorf("TheoPrice is not set to expected value.\n IP:%v\nOP%v\n", v.theoPrice, op.TheoPrice())
		}

	}

}

func TestSetImpliedVolatility(t *testing.T) {
	op := newRawOption(t)

	cases := []struct {
		impliedVol float32
	}{
		{-100.5},
		{-10},
		{0},
		{10},
		{100.5},
	}

	for _, v := range cases {
		err := op.SetImpliedVolatility(v.impliedVol)
		if v.impliedVol < 0 && err == nil {
			t.Errorf("ImpliedVolatility cannot be less than 0: %v", v.impliedVol)
		}

		if v.impliedVol >= 0 && op.ImpliedVolatility() != v.impliedVol {
			t.Errorf("ImpliedVolatility is not set to expected value.\n IP:%v\nOP%v\n", v.impliedVol, op.ImpliedVolatility())
		}

	}

}

func TestSetStrike(t *testing.T) {
	op := newRawOption(t)

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
		err := op.SetStrike(v.strike)
		if v.strike < 0 && err == nil {
			t.Errorf("Strike cannot be less than 0: %v", v.strike)
		}

		if v.strike >= 0 && op.Strike() != v.strike {
			t.Errorf("Strike is not set to expected value.\n IP:%v\nOP%v\n", v.strike, op.Strike())
		}

	}

}

func TestSetDaysToExpiration(t *testing.T) {
	op := newRawOption(t)

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
		err := op.SetDaysToExpiration(v.dte)
		if v.dte < 0 && err == nil {
			t.Errorf("DaysToExpiration cannot be less than 0: %v", v.dte)
		}

		if v.dte >= 0 && op.DaysToExpiration() != v.dte {
			t.Errorf("DaysToExpiration is not set to expected value.\n IP:%v\nOP%v\n", v.dte, op.DaysToExpiration())
		}

	}

}

func TestSetExpirationDate(t *testing.T) {
	o := newRawOption(t)

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
		if err := o.SetExpirationDate(v.d); err == nil && o.ExpirationDate().Before(currDate) {
			t.Errorf("Expiration date is in the past .\nIP Date: %v\nCP Date: %v\n", v.d, currDate)
		}

	}
}
