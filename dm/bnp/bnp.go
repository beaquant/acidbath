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

package balance

//Balance contains the balance information for the account connected to the broker session
type Balance struct {
	netLiquidity      float64
	optionBuyingPower float64
}

func New() *Balance {
	return &Balance{}
}

//SetNetLiquidity sets the current liquidity of the account connected to the broker session
func (b *Balance) SetNetLiquidity(liq float64) {
	b.netLiquidity = liq
}

//SetOptionBuyingPower sets the current option buying power of the account connected to the broker session
func (b *Balance) SetOptionBuyingPower(obp float64) {
	b.optionBuyingPower = obp
}

//NetLiquidity pulls the current liquidity of the account connected to the broker session
func (b *Balance) NetLiquidity() float64 {
	return b.netLiquidity
}

//OptionBuyingPower pulls the current option buying power of the account connected to the broker session
func (b *Balance) OptionBuyingPower() float64 {
	return b.optionBuyingPower
}

//Copy returns a copy of the balance structure
//TODO: MEDIUM : update to do deep copy of dereferenced pointer values
func (b *Balance) Copy() *Balance {
	dst := &Balance{}

	dst.netLiquidity = b.netLiquidity
	dst.optionBuyingPower = b.optionBuyingPower

	return dst
}
