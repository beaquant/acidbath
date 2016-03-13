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

package portfolio

import (
	"math/big"

	"github.com/marklaczynski/acidbath/dm/asset"
	"github.com/marklaczynski/acidbath/dm/bnp"
	"github.com/marklaczynski/acidbath/dm/optionchain/option"
	"github.com/marklaczynski/acidbath/lib/financial"
)

type Portfolio struct {
	positions        [][]*PositionType
	portfolioBalance *balance.Balance
}

func (p *Portfolio) Balance() *balance.Balance {
	return p.portfolioBalance
}

func (p *Portfolio) SetBalance(b *balance.Balance) {
	p.portfolioBalance = b
}

func NewPortfolio() *Portfolio {
	tmpPositions := make([][]*PositionType, asset.MaxAssetType, asset.MaxAssetType)
	for idx := 0; idx < int(asset.MaxAssetType); idx++ {
		tmpPositions[idx] = make([]*PositionType, 0, 0)
	}

	return &Portfolio{
		positions:        tmpPositions,
		portfolioBalance: balance.New(),
	}
}

func (p *Portfolio) Copy() *Portfolio {
	dst := &Portfolio{}
	dst.positions = make([][]*PositionType, asset.MaxAssetType, asset.MaxAssetType)
	for idx := 0; idx < int(asset.MaxAssetType); idx++ {
		dst.positions[idx] = make([]*PositionType, 0, 0)
		for _, v := range p.positions[idx] {
			dst.positions[idx] = append(dst.positions[idx], v)
		}
	}

	dst.portfolioBalance = p.Balance().Copy()

	return dst
}

func (b *Portfolio) AddPosition(st asset.AssetType, positions *PositionType) {
	b.positions[int(st)] = append(b.positions[int(st)], positions)
}

func (b *Portfolio) Position(st asset.AssetType) []*PositionType {
	return b.positions[st]
}

func (p *Portfolio) BetaWeightedDelta(targetBetaStock *asset.Stock) float64 {
	var tmpPortfolioBetaWeightedDelta float64 = 0.0

	for _, currOptionPosition := range p.Position(asset.OptionType) {
		tmpPortfolioBetaWeightedDelta = tmpPortfolioBetaWeightedDelta + (currOptionPosition.BetaWeightedDelta(targetBetaStock) * currOptionPosition.Quantity())
	}

	for _, currStockPosition := range p.Position(asset.EquityType) {
		tmpPortfolioBetaWeightedDelta = tmpPortfolioBetaWeightedDelta + (currStockPosition.BetaWeightedDelta(targetBetaStock) * currStockPosition.Quantity())
	}

	return tmpPortfolioBetaWeightedDelta
}

func (p *Portfolio) Gamma() float64 {
	var tmpPortfolioGamma float64 = 0.0

	for _, currOptionPosition := range p.Position(asset.OptionType) {
		tmpPortfolioGamma = tmpPortfolioGamma + (currOptionPosition.UnderlyingOption().Gamma() * currOptionPosition.Quantity())
	}

	for _, currStockPosition := range p.Position(asset.EquityType) {
		tmpPortfolioGamma = tmpPortfolioGamma + (asset.StockGamma * currStockPosition.Quantity())
	}

	return tmpPortfolioGamma
}

func (p *Portfolio) Theta() float64 {
	var tmpPortfolioTheta float64 = 0.0

	for _, currOptionPosition := range p.Position(asset.OptionType) {
		tmpPortfolioTheta = tmpPortfolioTheta + (currOptionPosition.UnderlyingOption().Theta() * currOptionPosition.Quantity())
	}

	for _, currStockPosition := range p.Position(asset.EquityType) {
		tmpPortfolioTheta = tmpPortfolioTheta + (asset.StockTheta * currStockPosition.Quantity())
	}

	return tmpPortfolioTheta
}

func (p *Portfolio) Vega() float64 {
	var tmpPortfolioVega float64 = 0.0

	for _, currOptionPosition := range p.Position(asset.OptionType) {
		tmpPortfolioVega = tmpPortfolioVega + (currOptionPosition.UnderlyingOption().Vega() * currOptionPosition.Quantity())
	}

	for _, currStockPosition := range p.Position(asset.EquityType) {
		tmpPortfolioVega = tmpPortfolioVega + (asset.StockVega * currStockPosition.Quantity())
	}

	return tmpPortfolioVega
}

/*
if AssetType == Option
	Stock = retrievesnapshot()
	Option = retrievesnapshot()

if AssetType == Stock
	Stock = retrievesnapshot()
	Option = nil
*/
type PositionType struct {
	quantity         float64
	symbol           string
	assetType        asset.AssetType
	cusip            string
	accountType      int64 // enum
	closePrice       financial.Money
	positionType     string //enum
	averagePrice     financial.Money
	currentValue     financial.Money
	underlyingSymbol string
	putCallIndicator string // enum

	//Supplamental data
	underlyingStock  *asset.Stock
	underlyingOption *option.Option
}

func NewPosition() *PositionType {
	tmpPos := &PositionType{}

	tmpPos.closePrice.Value = big.NewRat(0, 1)
	tmpPos.averagePrice.Value = big.NewRat(0, 1)
	tmpPos.currentValue.Value = big.NewRat(0, 1)

	return tmpPos
}

func (p *PositionType) UnderlyingOption() *option.Option {
	return p.underlyingOption
}

func (p *PositionType) SetUnderlyingOption(newOption *option.Option) {
	p.underlyingOption = newOption
}

func (p *PositionType) UnderlyingStock() *asset.Stock {
	return p.underlyingStock
}

func (p *PositionType) SetUnderlyingStock(newStock *asset.Stock) {
	p.underlyingStock = newStock
}

func (p *PositionType) PutCallIndicator() string {
	return p.putCallIndicator
}

func (p *PositionType) SetPutCallIndicator(newPutCallIndicator string) {
	p.putCallIndicator = newPutCallIndicator
}

func (p *PositionType) UnderlyingSymbol() string {
	return p.underlyingSymbol
}

func (p *PositionType) SetUnderlyingSymbol(newUnderlyingSymbol string) {
	p.underlyingSymbol = newUnderlyingSymbol
}

func (p *PositionType) PositionType() string {
	return p.positionType
}

func (p *PositionType) SetPositionType(newPositionType string) {
	p.positionType = newPositionType
}

func (p *PositionType) CurrentValue() financial.Money {
	return p.currentValue
}

func (p *PositionType) SetCurrentValue(newCurrentValue financial.Money) {
	p.currentValue.Value.Set(newCurrentValue.Value)
}

func (p *PositionType) AveragePrice() financial.Money {
	return p.averagePrice
}

func (p *PositionType) SetAveragePrice(newAveragePrice financial.Money) {

	p.averagePrice.Value.Set(newAveragePrice.Value)
}

func (p *PositionType) ClosePrice() financial.Money {
	return p.closePrice
}

func (p *PositionType) SetClosePrice(newClosePrice financial.Money) {
	p.closePrice.Value.Set(newClosePrice.Value)
}

func (p *PositionType) Quantity() float64 {
	return p.quantity
}

func (p *PositionType) SetQuantity(newQuantity float64) {
	p.quantity = newQuantity
}

func (p *PositionType) Symbol() string {
	return p.symbol
}

func (p *PositionType) SetSymbol(newSymbol string) {
	p.symbol = newSymbol
}

func (p *PositionType) AssetType() asset.AssetType {
	return p.assetType
}

func (p *PositionType) SetAssetType(newAssetType asset.AssetType) {
	p.assetType = newAssetType
}

func (p *PositionType) Cusip() string {
	return p.cusip
}

func (p *PositionType) SetCusip(newCusip string) {
	p.cusip = newCusip
}

func (p *PositionType) AccountType() int64 {
	return p.accountType
}

func (p *PositionType) SetAccountType(newAccountType int64) {
	p.accountType = newAccountType
}

/*
	How to calculate beta weighted delta

	ewz: 20.70
	spy (beta target): 202.20
	ewz beta: 1.57
	ewz 19 put delta: 21.90
	beta w spy delta 19 pub: 3.56

	21.9 = 6.15 * 3.56
	(spy) 202.20 / (ewz) 20.70 = 9.8 / (ews beta) 1.57  = 6.22
	spy stock price / pos stock price = notionalAdjFactor (rename)
	notionalAdjFactor / pos beta = delta adjuster

	21.9 / (adj) 6.15 = 3.56
	pos delta / delta adjuster = spy beta weighted delta
*/
// this is for options only
// Assume p.UnderlyingStock != nil && p.UnderlyingStock != nil
func (p *PositionType) BetaWeightedDelta(targetBetaStock *asset.Stock) float64 {
	switch p.assetType {
	case asset.OptionType:
		// NOTE: FUTURE : in the future take the mid between bid ask
		betaTargetLastTradePrice, _ := targetBetaStock.LastTradePrice().Value.Float64()
		currentOptionPositionStockLastTradePrice, _ := p.UnderlyingStock().LastTradePrice().Value.Float64()
		tmpTargetDailyCloseChange := targetBetaStock.DailyCloseChange()

		return p.UnderlyingOption().BetaWeightedDelta(betaTargetLastTradePrice, currentOptionPositionStockLastTradePrice, p.UnderlyingStock().Beta(&tmpTargetDailyCloseChange))
	case asset.EquityType:
		// NOTE: FUTURE : in the future take the mid between bid ask
		const SingleStockDelta = 1

		betaTargetLastTradePrice, _ := targetBetaStock.LastTradePrice().Value.Float64()
		currentStockPositionStockLastTradePrice, _ := p.UnderlyingStock().LastTradePrice().Value.Float64()

		tmpTargetDailyCloseChange := targetBetaStock.DailyCloseChange()
		return (SingleStockDelta / ((betaTargetLastTradePrice / currentStockPositionStockLastTradePrice) / p.UnderlyingStock().Beta(&tmpTargetDailyCloseChange)))
	default:
		panic("Should not reach here")
	}
}
