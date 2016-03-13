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

//Package generic provides an apis for various structures. This one is related to the Broker
package generic

import (
	"github.com/marklaczynski/acidbath/dm/asset"
	"github.com/marklaczynski/acidbath/dm/optionchain/option"
	"github.com/marklaczynski/acidbath/dm/order"
	"github.com/marklaczynski/acidbath/dm/orderbook"
	"github.com/marklaczynski/acidbath/dm/ordermessage"
	"github.com/marklaczynski/acidbath/dm/portfolio"
	"github.com/marklaczynski/acidbath/dm/watchlists"
	"github.com/marklaczynski/acidbath/eventproc/factory"
)

/*
Notes: I'de like to make this RESTful one day
*/

//Broker is a stock/option electronic broker you can interface with
//As a convention, use a string tickerSymbol to pass information about which instrument to key
type Broker interface {
	Login(loginid string, pass string) error
	Logout() error
	RetrieveSnapshot(symbol string, assetType asset.AssetType, security interface{}) error
	RetrieveImpliedVolatilityHistory(stockSymbol string, stock *asset.Stock) error
	RetrievePriceHistory(stockSymbol string, stock *asset.Stock) error
	RetrievePortfolio(newPortfolio *portfolio.Portfolio) error
	AddStockOptionsToStream(stock *asset.Stock) error
	RemoveStockOptionsFromStream(stock *asset.Stock) error
	AddOptionToStrategy(opt *option.Option, strategy factory.StrategyType) ([]string, error)
	RemoveOptionFromStrategy(opt *option.Option, strategy factory.StrategyType) ([]string, error)
	SendSingleLegOptionTrade(order *order.Order) error
	CancelOrder(orderids []string) error
	RetrieveOrderBook(accountid string, ob *orderbook.OrderBook) error
	RegisterOptionUpdateChan(id string) chan *option.Option
	DeregisterOptionUpdateChan(id string)
	RegisterPortfolioUpdateChan(id string) chan *portfolio.Portfolio
	DeregisterPortfolioUpdateChan(id string)
	RegisterOrderUpdateChan(id string) chan *ordermessage.Message
	DeregisterOrderUpdateChan(id string)

	RetrieveWatchlists(wls *watchlists.Watchlists) error
}
