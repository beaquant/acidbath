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

//Package tdapi supports an interface to the TDAmeritrade API.
package tdapi

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/marklaczynski/acidbath/dm/order"
	"github.com/marklaczynski/acidbath/lib/orderconst"
)

type tdOrder struct {
	accountID string
	order     *order.Order
}

func (o *tdOrder) orderString() string {
	return fmt.Sprintf("action=%s~accountid=%s~action=%s~actprice=%s~expire=%s~exday=%s~exmonth=%s~exyear=%s~ordtype=%s~price=%s~quantity=%d~routing=%s~spinstructions=%s~symbol=%s",
		o.order.Action(),
		o.accountID,
		o.order.Action(),
		o.order.ActivatePrice().Value.FloatString(2),
		o.order.Expire(),
		o.order.ExDay(),
		o.order.ExMonth(),
		o.order.ExYear(),
		o.order.OrderType(),
		o.order.Price().Value.FloatString(2),
		o.order.Quantity(),
		o.order.Routing(),
		o.order.SpecialInstructions(),
		o.order.Symbol())
}

//validate is the internal API call, all other validate* are helper functions
func (o *tdOrder) validate() error {
	return o.validateNewOptionTrade()
}

func (o *tdOrder) validateNewOptionTrade() error {

	// action, symbol, ordtype, quantity, accountid, and expire are required parameters
	// symbol, quantity, account id are checked through logic
	// action, ordtype are defaulted via enum
	if o.order.Symbol() == "" {
		return errors.New("Symbol is required")
	}

	if o.order.Quantity() < 1 {
		return errors.New("Quantity is required")
	}

	if o.accountID == "" {
		return errors.New("AccountID is required")
	}

	if o.order.OrderType() == orderconst.Market {
		if o.order.Expire() != orderconst.Day {
			return errors.New("Market ordertype must use Day expiry")
		}
		if o.order.SpecialInstructions() != orderconst.None {
			return errors.New("Market ordertype must use no speical instructions")
		}
		if o.order.Price().Value.Cmp(big.NewRat(0, 1)) != 0 {
			return errors.New("Market ordertype must have null or 0 price")
		}
		if o.order.ActivatePrice().Value.Cmp(big.NewRat(0, 1)) != 0 {
			return errors.New("Market ordertype must have null or 0 activatePrice")
		}
	}

	if o.order.OrderType() == orderconst.Limit {
		if o.order.Price().Value.Cmp(big.NewRat(0, 1)) < 0 {
			return errors.New("Price must be greater than 0")
		}

		if o.order.ActivatePrice().Value.Cmp(big.NewRat(0, 1)) != 0 {
			return errors.New("Market ordertype must have null or 0 activatePrice")
		}

		if o.order.Expire() != orderconst.GTC && o.order.Expire() != orderconst.Day {
			return errors.New("Limit must use either GTC or Day expiration")
		}

		if o.order.Expire() == orderconst.GTC {
			if o.order.SpecialInstructions() != orderconst.None && o.order.SpecialInstructions() != orderconst.Aon {
				return errors.New("GTC Limit ordertype must None or Aon as specialInstructions")
			}
		}

		if o.order.Expire() == orderconst.Day {
			if o.order.SpecialInstructions() != orderconst.None && o.order.SpecialInstructions() != orderconst.Aon && o.order.SpecialInstructions() != orderconst.Fok {
				return errors.New("GTC Limit ordertype must None, Aon, or Fok as specialInstructions")
			}
		}

	}

	if o.order.OrderType() == orderconst.StopMarket {
		if o.order.ActivatePrice().Value.Cmp(big.NewRat(0, 1)) < 0 {
			return errors.New("Active price must be greater than 0")
		}

		if o.order.Price().Value.Cmp(big.NewRat(0, 1)) != 0 {
			return errors.New("Stop Market ordertype must have null or 0 price")
		}

		if o.order.Expire() != orderconst.GTC && o.order.Expire() != orderconst.Day {
			return errors.New("Stop Market must use either GTC or Day expiration")
		}

		if o.order.SpecialInstructions() != orderconst.None && o.order.SpecialInstructions() != orderconst.Aon {
			return errors.New("Stop Market ordertype must None or Aon as specialInstructions")
		}
	}

	if o.order.OrderType() == orderconst.StopLimit {
		if o.order.Price().Value.Cmp(big.NewRat(0, 1)) < 0 {
			return errors.New("Price must be greater than 0")
		}

		if o.order.ActivatePrice().Value.Cmp(big.NewRat(0, 1)) < 0 {
			return errors.New("Active price must be greater than 0")
		}

		if o.order.Expire() != orderconst.GTC && o.order.Expire() != orderconst.Day {
			return errors.New("Stop Limit must use either GTC or Day expiration")
		}

		if o.order.SpecialInstructions() != orderconst.None && o.order.SpecialInstructions() != orderconst.Aon {
			return errors.New("Stop Market ordertype must None or Aon as specialInstructions")
		}
	}

	logDebug.Printf("order: %#v\n", o)
	logDebug.Printf("order: %#v\n", o.order)
	if o.order.ExDay() != 0 || o.order.ExMonth() != 0 || o.order.ExYear() != 0 {
		if o.order.ExDay() < 1 || o.order.ExDay() > 31 {
			return errors.New("ExDay is out of range [1,31]")
		}

		if o.order.ExMonth() < 1 || o.order.ExMonth() > 12 {
			return errors.New("ExMonth is out of range [1,12]")
		}

		if int(o.order.ExYear()) < time.Now().Year() {
			return fmt.Errorf("Selected year %d is before current year %d", o.order.ExYear(), time.Now().Year())
		}

		//validation documentation lists gtc_ext, but i didn't see it defined in list of domain values
		if o.order.Expire() != orderconst.GTC {
			return errors.New("Expiry must be GTC if using ex-Date")
		}

		if time.Now().Sub(time.Date(int(o.order.ExYear()), time.Month(int(o.order.ExMonth())), int(o.order.ExDay()), 0, 0, 0, 0, time.UTC)).Hours() < 24 {
			return fmt.Errorf("Selected date %s needs to be in future. Current date %s", time.Date(int(o.order.ExYear()), time.Month(int(o.order.ExMonth())), int(o.order.ExDay()), 0, 0, 0, 0, time.UTC), time.Now())
		}

		//the date cannot be after the last day of the following month
		if time.Now().AddDate(0, 1, 1).After(time.Date(int(o.order.ExYear()), time.Month(int(o.order.ExMonth())), int(o.order.ExDay()), 0, 0, 0, 0, time.UTC)) {
			return fmt.Errorf("Selected date %s cannot be after last day of next month. Future Date %s", time.Date(int(o.order.ExYear()), time.Month(int(o.order.ExMonth())), int(o.order.ExDay()), 0, 0, 0, 0, time.UTC), time.Now().AddDate(0, 1, 1))
		}
	}

	if o.order.Expire() == orderconst.GTC {
		if o.order.ExDay() < 1 || o.order.ExDay() > 31 {
			return errors.New("ExDay is out of range [1,31]")
		}

		if o.order.ExMonth() < 1 || o.order.ExMonth() > 12 {
			return errors.New("ExMonth is out of range [1,12]")
		}

		if int(o.order.ExYear()) < time.Now().Year() {
			return fmt.Errorf("Selected year %d is before current year %d", o.order.ExYear(), time.Now().Year())
		}
	}

	return nil
}
