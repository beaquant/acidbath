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

//Package optionstrike implements an option strike
package optionstrike

import (
	"errors"

	"github.com/marklaczynski/acidbath/dm/optionchain/option"
)

type index int

const (
	numberOfOptions = 2
)

//OptionStrike represents an option strike
type OptionStrike struct {
	strike float64
	option [numberOfOptions]*option.Option
}

//NewOptionStrike returns a pointer to an OptionStrike
func NewOptionStrike(strike float64) (os *OptionStrike, err error) {
	os = &OptionStrike{}
	if err := os.SetStrike(strike); err != nil {
		return nil, err
	}
	return os, nil
}

//Strike returns the current strike
func (os *OptionStrike) Strike() float64 {
	return os.strike
}

//SetStrike sets the current strike. Must be > 0
func (os *OptionStrike) SetStrike(strike float64) error {
	if strike <= 0 {
		return errors.New("Strike price must be positive")
	}
	os.strike = strike
	return nil
}

//Option returns an option
func (os *OptionStrike) Option(oType option.TypeOfOption) *option.Option {
	return os.option[int(oType)]
}

//AddOption adds an option on the current strike. Use this when adding individual options to a strike
func (os *OptionStrike) AddOption(option *option.Option) error {
	if os.option[option.OptionType()] != nil {
		return errors.New("Option already exists")
	}
	os.option[option.OptionType()] = option
	return nil
}
