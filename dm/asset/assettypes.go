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

import "fmt"

type AssetType int

/* these exist in balances & position */

const (
	EquityType AssetType = iota
	OptionType
	MutualFundType
	BondType
	MoneyMarketType
	SavingsType
	IndexType
	MaxAssetType
)

func (a AssetType) String() string {
	switch a {
	case EquityType:
		return "Equity or ETF"
	case MutualFundType:
		return "MutualFund"
	case IndexType:
		return "Index"
	case OptionType:
		return "Option"
	case BondType:
		return "Bond"
	case MoneyMarketType:
		return "Money Market"
	case SavingsType:
		return "Savings"
	default:
		panic(fmt.Sprintf("Invalid AssetType %#v\n", a))
	}
}
