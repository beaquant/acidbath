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

package asset

import (
	"errors"
	"math/big"

	"github.com/marklaczynski/acidbath/lib/financial"
)

type Quote struct {
	symbol         string
	description    string
	bidPrice       financial.Money
	askPrice       financial.Money
	lastTradePrice financial.Money
	closePrice     financial.Money
}

func (q *Quote) NewQuote() {
	q.bidPrice.Value = big.NewRat(0, 1)
	q.askPrice.Value = big.NewRat(0, 1)
	q.lastTradePrice.Value = big.NewRat(0, 1)
	q.closePrice.Value = big.NewRat(0, 1)
}

func (q *Quote) Symbol() string {
	return q.symbol
}

func (q *Quote) SetSymbol(newValue string) {
	q.symbol = newValue
}

func (q *Quote) Description() string {
	return q.description
}

func (q *Quote) SetDescription(newValue string) {
	q.description = newValue
}

//BidPrice returns the bid price
func (q *Quote) BidPrice() financial.Money {
	return q.bidPrice
}

//SetBidPrice sets the bid price
func (q *Quote) SetBidPrice(newBidPrice financial.Money) error {
	if newBidPrice.Value.Cmp(big.NewRat(0, 1)) < 0 {
		return errors.New("BidPrice cannot be less than 0")
	}

	q.bidPrice.Value.Set(newBidPrice.Value)
	return nil
}

//AskPrice returns the ask price
func (q *Quote) AskPrice() financial.Money {
	return q.askPrice
}

//SetAskPrice sets the ask price
func (q *Quote) SetAskPrice(newAskPrice financial.Money) error {
	if newAskPrice.Value.Cmp(big.NewRat(0, 1)) < 0 {
		return errors.New("AskPrice cannot be less than 0")
	}

	q.askPrice.Value.Set(newAskPrice.Value)
	return nil
}

//Last returns the last trade price
func (q *Quote) LastTradePrice() financial.Money {
	return q.lastTradePrice
}

//SetLast sets the last trade price
func (q *Quote) SetLastTradePrice(newLastTradePrice financial.Money) error {
	if newLastTradePrice.Value.Cmp(big.NewRat(0, 1)) < 0 {
		return errors.New("Last trade price cannot be less than 0")
	}

	q.lastTradePrice.Value.Set(newLastTradePrice.Value)
	return nil
}
