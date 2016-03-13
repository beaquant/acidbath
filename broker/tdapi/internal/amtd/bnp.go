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

package amtd

import (
	"encoding/xml"

	"github.com/marklaczynski/acidbath/lib/financial"
	"github.com/marklaczynski/acidbath/lib/types"
)

//Portfolio represent the TD Balance and Position structure
type Portfolio struct {
	XMLName   xml.Name `xml:"amtd"`
	Error              //inline struct
	Balance   balanceXML
	Positions positionsXML
}

type balanceXML struct {
	XMLName                          xml.Name       `xml:"balance"`
	Error                            string         `xml:"error"`
	AccountID                        string         `xml:"account-id"`
	DayTrader                        bool           `xml:"day-trader"`
	RoundTrips                       types.XMLInt64 `xml:"round-trips"`
	RestictedClosingTransactionsOnly bool           `xml:"resticted-closing-transactions-only"`
	CashBalance                      cashBalanceXML
	MoneyMarketBalance               moneyMarketBalanceXML
	LongStockValue                   longStockValueXML
	LongOptionValue                  longOptionValueXML
	ShortStockValue                  shortStockValueXML
	ShortOptionValue                 shortOptionValueXML
	MutualFundValue                  mutualFundValueXML
	BondValue                        bondValueXML
	AccountValue                     accountValueXML
	PendingDeposits                  pendingDepositsXML
	SavingsBalance                   savingsBalanceXML
	MarginBalance                    marginBalanceXML
	ShortBalance                     shortBalanceXML
	LongMarginableValue              longMarginableValueXML
	ShortMarginableValue             shortMarginableValueXML
	MarginEquity                     marginEquityXML
	EquityPercentage                 equityPercentageXML
	StockBuyingPower                 types.XMLFloat64 `xml:"stock-buying-power"`
	OptionBuyingPower                types.XMLFloat64 `xml:"option-buying-power"`
	DayTradingBuyingPower            types.XMLFloat64 `xml:"day-trading-buying-power"`
	AvailableFundsForTrading         types.XMLFloat64 `xml:"available-funds-for-trading"`
	MaintenanceRequirement           maintenanceRequirementXML
	MaintenanceCallValue             maintenanceCallValueXML
	RegulationTCallValue             regulationTCallValueXML
	DayTradingCallValue              dayTradingCallValueXML
	DayEquiteyCallValue              types.XMLFloat64 `xml:"day-equity-call-value"`
	InCall                           bool             `xml:"in-call"`
	InPotentialCall                  bool             `xml:"in-potential-call"`
}

type cashBalanceXML struct {
	XMLName xml.Name `xml:"cash-balance"`
	icc
}

type moneyMarketBalanceXML struct {
	XMLName xml.Name `xml:"money-market-balance"`
	icc
}

type longStockValueXML struct {
	XMLName xml.Name `xml:"long-stock-value"`
	icc
}

type longOptionValueXML struct {
	XMLName xml.Name `xml:"long-option-value"`
	icc
}

type shortStockValueXML struct {
	XMLName xml.Name `xml:"short-stock-value"`
	icc
}

type shortOptionValueXML struct {
	XMLName xml.Name `xml:"short-option-value"`
	icc
}

type mutualFundValueXML struct {
	XMLName xml.Name `xml:"mutual-fund-value"`
	icc
}

type bondValueXML struct {
	XMLName xml.Name `xml:"bond-value"`
	icc
}

type accountValueXML struct {
	XMLName xml.Name `xml:"account-value"`
	icc
}

type pendingDepositsXML struct {
	XMLName xml.Name `xml:"pending-deposits"`
	icc
}

type savingsBalanceXML struct {
	XMLName xml.Name         `xml:"savings-balance"`
	Current types.XMLFloat64 `xml:"current"`
}

type marginBalanceXML struct {
	XMLName xml.Name `xml:"margin-balance"`
	icc
}

type shortBalanceXML struct {
	XMLName xml.Name `xml:"short-balance"`
	icc
}

type longMarginableValueXML struct {
	XMLName xml.Name `xml:"long-marginable-value"`
	icc
}

type shortMarginableValueXML struct {
	XMLName xml.Name `xml:"short-marginable-value"`
	icc
}

type marginEquityXML struct {
	XMLName xml.Name `xml:"margin-equity"`
	icc
}

type equityPercentageXML struct {
	XMLName xml.Name `xml:"equity-percentage"`
	icc
}

type maintenanceRequirementXML struct {
	XMLName xml.Name `xml:"maintenance-requirement"`
	icc
}

type maintenanceCallValueXML struct {
	XMLName xml.Name `xml:"maintenance-call-value"`
	icp
}

type regulationTCallValueXML struct {
	XMLName xml.Name `xml:"regulation-t-call-value"`
	icp
}

type dayTradingCallValueXML struct {
	XMLName xml.Name `xml:"day-trading-call-value"`
	ip
}

type icc struct {
	Initial types.XMLFloat64 `xml:"initial"`
	Current types.XMLFloat64 `xml:"current"`
	Change  types.XMLFloat64 `xml:"change"`
}

type icp struct {
	Initial   types.XMLFloat64 `xml:"initial"`
	Current   types.XMLFloat64 `xml:"current"`
	Potential types.XMLFloat64 `xml:"potential"`
}

type ip struct {
	Potential types.XMLFloat64 `xml:"potential"`
	Initial   types.XMLFloat64 `xml:"initial"`
}

type positionsXML struct {
	XMLName     xml.Name `xml:"positions"`
	Error       string   `xml:"error"`
	AccountID   string   `xml:"account-id"`
	Stocks      stocksXML
	Options     optionsXML
	Funds       fundsXML
	Bonds       bondsXML
	MoneyMarket moneyMarketXML
	Savings     savingsXML
}

type stocksXML struct {
	XMLName  xml.Name      `xml:"stocks"`
	Position []positionXML `xml:"position"`
}

type optionsXML struct {
	XMLName  xml.Name      `xml:"options"`
	Position []positionXML `xml:"position"`
}

type fundsXML struct {
	XMLName  xml.Name      `xml:"funds"`
	Position []positionXML `xml:"position"`
}

type bondsXML struct {
	XMLName  xml.Name      `xml:"bonds"`
	Position []positionXML `xml:"position"`
}

type moneyMarketXML struct {
	XMLName  xml.Name      `xml:"money-market"`
	Position []positionXML `xml:"position"`
}

type savingsXML struct {
	XMLName  xml.Name      `xml:"savings"`
	Position []positionXML `xml:"position"`
}

type positionXML struct {
	XMLName                xml.Name         `xml:"position"`
	Error                  string           `xml:"error"`
	Quantity               types.XMLFloat64 `xml:"quantity"`
	Security               securityXML
	AccountType            types.XMLInt64   `xml:"account-type"`
	ClosePrice             financial.Money  `xml:"close-price"`
	PositionType           string           `xml:"position-type"`
	AveragePrice           financial.Money  `xml:"average-price"`
	CurrentValue           financial.Money  `xml:"current-value"`
	UnderlyingSymbol       string           `xml:"underlying-symbol"`
	PutCall                string           `xml:"put-call"`
	MaintenanceRequirement financial.Money  `xml:"maintenance-requirement"`
	BondFactor             types.XMLFloat64 `xml:"bond-factor"`
	Quote                  QuoteXML         `xml:"quote"`
}

type securityXML struct {
	XMLName              xml.Name `xml:"security"`
	Symbol               string   `xml:"symbol"`
	SymbolWithTypePrefix string   `xml:"symbol-with-type-prefix"`
	Description          string   `xml:"description"`
	AssetType            string   `xml:"asset-type"`
	Cusip                string   `xml:"cusip"`
}
