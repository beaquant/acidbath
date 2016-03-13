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

package watchlists

import (
	"log"
	"sort"

	"github.com/marklaczynski/acidbath/dm/asset"
	"github.com/marklaczynski/acidbath/lib/mjlog"
)

var (
	logInfo  = log.New(mjlog.CreateInfoFile(), "INFO  [watchlists]: ", log.LstdFlags|log.Lshortfile)
	logDebug = log.New(mjlog.CreateDebugFile(), "DEBUG [watchlists]: ", log.LstdFlags|log.Lshortfile)
	logError = log.New(mjlog.CreateErrorFile(), "ERROR [watchlists]: ", log.LstdFlags|log.Lshortfile)
)

type Watchlists struct {
	watchlist map[int64]*WatchlistType
}

func New() *Watchlists {
	return &Watchlists{
		watchlist: make(map[int64]*WatchlistType),
	}
}

func (wls *Watchlists) Watchlist(id int64) *WatchlistType {
	return wls.watchlist[id]
}

func (wls *Watchlists) AddWatchlist(name string, id int64) {
	wls.watchlist[id] = NewWatchlist(name)
}

func (wls *Watchlists) AddWatchedSymbol(id int64, symbol string, assetType asset.AssetType) {

	var newWatchedSymbol WatchedSymbolType
	switch assetType {
	case asset.EquityType:
		newWatchedSymbol = NewWatchedSymbol(symbol)
	default:
		logError.Printf("unsupported assetType :%s: in watchlist id %d. symbol: %s", assetType, id, symbol)
	}

	wls.watchlist[id].watchedSymbols = append(wls.watchlist[id].watchedSymbols, newWatchedSymbol)
}

func (wl WatchlistType) WatchedSymbolsInIVRankRange(lowestIvRankParam float32, highestIvRankParam float32) WatchedSymbolSlice {
	//logDebug.Printf("watchlist: %#v", wl)
	watchedSymbolsInIvRange := make(WatchedSymbolSlice, 0, 0)
	for _, currWs := range wl.WatchedSymbols() {
		logDebug.Printf("Curr Stock %s IV Rank: %.2f", currWs.Stock().Symbol(), currWs.Stock().CurrentImpliedVolatilityRank())
		if currWs.Stock().CurrentImpliedVolatilityRank() >= lowestIvRankParam && currWs.Stock().CurrentImpliedVolatilityRank() <= highestIvRankParam {
			watchedSymbolsInIvRange = append(watchedSymbolsInIvRange, currWs)
			//logDebug.Printf("added %s to watchedSymbolsInIvRange %#v", currWs.Stock().Symbol(), watchedSymbolsInIvRange)
		}
	}
	return watchedSymbolsInIvRange
}

func (wl WatchlistType) HighestIvRankedSymbol() WatchedSymbolType {
	sortedWatchlistByIv := make(WatchedSymbolSlice, len(wl.WatchedSymbols()), len(wl.WatchedSymbols()))
	copy(sortedWatchlistByIv, wl.WatchedSymbols())
	sort.Sort(sortedWatchlistByIv)

	return sortedWatchlistByIv[len(sortedWatchlistByIv)-1]
}

func (wl WatchlistType) HighestIvRankedSymbolExcluding(symbols []string) WatchedSymbolType {
	sortedWatchlistByIv := make(WatchedSymbolSlice, len(wl.WatchedSymbols()), len(wl.WatchedSymbols()))
	copy(sortedWatchlistByIv, wl.WatchedSymbols())
	sort.Sort(sortedWatchlistByIv)

	//TODO: FINISH loop from the highest to lowest and pick the first occurance which doesn't match any of the symbolsParam
	return sortedWatchlistByIv[len(sortedWatchlistByIv)-1]
}

func (wl WatchlistType) LowestIvRankedSymbol() WatchedSymbolType {
	sortedWatchlistByIv := make(WatchedSymbolSlice, len(wl.WatchedSymbols()), len(wl.WatchedSymbols()))
	copy(sortedWatchlistByIv, wl.WatchedSymbols())
	sort.Sort(sortedWatchlistByIv)

	return sortedWatchlistByIv[0]
}

type WatchlistType struct {
	name           string
	watchedSymbols WatchedSymbolSlice
}

func NewWatchlist(newName string) *WatchlistType {
	return &WatchlistType{
		name:           newName,
		watchedSymbols: make(WatchedSymbolSlice, 0, 0),
	}
}

func (wl WatchlistType) WatchedSymbols() WatchedSymbolSlice {
	return wl.watchedSymbols
}

func (wl WatchlistType) SetWatchedSymbols(newWatchedSymbols *WatchedSymbolSlice) {
	wl.watchedSymbols = *newWatchedSymbols
}

func (wl WatchlistType) Name() string {
	return wl.name
}

func (wl WatchlistType) SetName(newName string) {
	wl.name = newName
}

type WatchedSymbolSlice []WatchedSymbolType

func (wss WatchedSymbolSlice) Len() int {
	return len(wss)
}

func (wss WatchedSymbolSlice) Less(i, j int) bool {
	return wss[i].Stock().CurrentImpliedVolatilityRank() < wss[j].Stock().CurrentImpliedVolatilityRank()
}

func (wss WatchedSymbolSlice) Swap(i, j int) {
	wss[i], wss[j] = wss[j], wss[i]
}

type WatchedSymbolType struct {
	stock *asset.Stock
}

func NewWatchedSymbol(symbol string) WatchedSymbolType {
	return WatchedSymbolType{
		stock: asset.NewStock(symbol),
	}
}

func (ws WatchedSymbolType) Stock() *asset.Stock {
	return ws.stock
}

func (ws WatchedSymbolType) SetStock(newStock *asset.Stock) {
	ws.stock = newStock
}
