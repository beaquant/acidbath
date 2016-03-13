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
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/grd/statistics"

	"github.com/marklaczynski/acidbath/dm/optionchain"
	"github.com/marklaczynski/acidbath/lib/financial"
	"github.com/marklaczynski/acidbath/lib/mjlog"
)

const (
	StockTheta = 0
	StockGamma = 0
	StockVega  = 0
)

var (
	logInfo  = log.New(mjlog.CreateInfoFile(), "INFO  [asset]: ", log.LstdFlags|log.Lshortfile)
	logDebug = log.New(mjlog.CreateDebugFile(), "DEBUG [asset]: ", log.LstdFlags|log.Lshortfile)
	logError = log.New(mjlog.CreateErrorFile(), "ERROR [asset]: ", log.LstdFlags|log.Lshortfile)
)

// Note: if i ever needed sorted by multiple fields: http://stackoverflow.com/questions/19759274/sort-points-structs-by-different-dimensions-in-go-lang
type ImpliedVolatilityTypeSlice []ImpliedVolatilityType

func (ivs ImpliedVolatilityTypeSlice) Len() int {
	return len(ivs)
}

func (ivs ImpliedVolatilityTypeSlice) Less(i, j int) bool {
	return ivs[i].impliedVolatility < ivs[j].impliedVolatility
}

func (ivs ImpliedVolatilityTypeSlice) Swap(i, j int) {
	ivs[i], ivs[j] = ivs[j], ivs[i]
}

type ImpliedVolatilityType struct {
	impliedVolatility float32
	timeStamp         time.Time
}

func NewImpliedVolInstance(iv float32, ts time.Time) ImpliedVolatilityType {
	return ImpliedVolatilityType{
		impliedVolatility: iv,
		timeStamp:         ts,
	}
}

func (iv ImpliedVolatilityType) String() string {
	return fmt.Sprintf("IV: %.2f on %s", iv.impliedVolatility, iv.timeStamp)
}

func (iv ImpliedVolatilityType) ImpliedVolatility() float32 {
	return iv.impliedVolatility
}

func (iv ImpliedVolatilityType) SetImpliedVolatility(newIv float32) {
	iv.impliedVolatility = newIv
}

func (iv ImpliedVolatilityType) TimeStamp() time.Time {
	return iv.timeStamp
}

func (iv ImpliedVolatilityType) SetTimeStamp(newTimeStamp time.Time) {
	iv.timeStamp = newTimeStamp
}

//Stock represent a stock
type Stock struct {
	Quote
	historicalImpliedVol ImpliedVolatilityTypeSlice
	historicalPrice      []PriceHistoryType
	optionChain          *optionchain.OptionChain
	//TODO: make use of the optionChain here... this is a major refactoring, and should be part of the effort where i create a caching system
}

func NewStock(newSymbol string) *Stock {
	s := &Stock{}
	s.NewQuote()
	s.SetSymbol(newSymbol)
	s.optionChain = nil
	return s
}

func (s *Stock) OptionChain() *optionchain.OptionChain {
	return s.optionChain
}

func (s *Stock) SetOptionChain(oc *optionchain.OptionChain) {
	s.optionChain = oc
}

func (s *Stock) DailyCloseChange() statistics.Float64 {
	tmpDailyCloseChange := make(statistics.Float64, 0, 0)
	var previousClose financial.Money
	for i, currPrice := range s.HistoricalPrice() {
		if i != 0 {
			currCloseF, _ := currPrice.Close().Value.Float64()
			prevCloseF, _ := previousClose.Value.Float64()
			tmpDailyCloseChange = append(tmpDailyCloseChange, math.Log(currCloseF/prevCloseF))
		}
		previousClose = currPrice.Close()
	}

	return tmpDailyCloseChange
}

func (s *Stock) Beta(targetDailyCloseChange *statistics.Float64) float64 {
	sourceDailyCloseChange := s.DailyCloseChange()
	//beta = cov( stk, spy ) / var ( spy )
	return statistics.Covariance(&sourceDailyCloseChange, targetDailyCloseChange) / statistics.Variance(targetDailyCloseChange)
}

func (s *Stock) SetHistoricalImpliedVol(newVolArray *ImpliedVolatilityTypeSlice) {
	s.historicalImpliedVol = *newVolArray
}

func (s *Stock) HistoricalImpliedVol() ImpliedVolatilityTypeSlice {
	return s.historicalImpliedVol
}

func (s *Stock) SetHistoricalPrice(newVolArray *[]PriceHistoryType) {
	s.historicalPrice = *newVolArray
}

func (s *Stock) HistoricalPrice() []PriceHistoryType {
	return s.historicalPrice
}

func (s *Stock) CurrentImpliedVolatility() float32 {
	return s.historicalImpliedVol[len(s.historicalImpliedVol)-1].impliedVolatility
}

func (s *Stock) CurrentImpliedVolatilityRank() float32 {
	sortedIv := make(ImpliedVolatilityTypeSlice, len(s.historicalImpliedVol), len(s.historicalImpliedVol))
	copy(sortedIv, s.historicalImpliedVol)
	sort.Sort(sortedIv)
	if len(sortedIv) == 0 {
		panic(fmt.Sprintf("issue with sorted iv %#v vs historicalIv %#v for stock %s", sortedIv, s.historicalImpliedVol, s.Symbol()))
	}
	lowIv := sortedIv[0].impliedVolatility
	highIv := sortedIv[len(sortedIv)-1].impliedVolatility
	rank := (s.CurrentImpliedVolatility() - lowIv) / (highIv - lowIv)

	return rank
}

// TODO: HIGH : FInish
func (s *Stock) Correlation(targetStock *Stock, length int) float32 {
	// corrlate s to targetStock
	return 0.0
}

/* FUTURE
func (s *Stock) CurrentImnpliedVolatailityPercentile() float64 {
	for key, val := range s.historicalImpliedVol {
	}
	return 0.0
}
*/

type PriceHistoryType struct {
	closePrice financial.Money
	timeStamp  time.Time
}

func NewPriceHistoryPoint(c financial.Money, ts time.Time) PriceHistoryType {
	return PriceHistoryType{
		closePrice: c,
		timeStamp:  ts,
	}
}

func (ph PriceHistoryType) Close() financial.Money {
	return ph.closePrice
}

func (ph PriceHistoryType) SetClose(newVal financial.Money) {
	ph.closePrice.Value.Set(newVal.Value)
}

func (ph PriceHistoryType) TimeStamp() time.Time {
	return ph.timeStamp
}

func (ph PriceHistoryType) SetTimeStamp(newVal time.Time) {
	ph.timeStamp = newVal
}

func (ph PriceHistoryType) String() string {
	return fmt.Sprintf("data: Close: %s on %s", ph.Close(), ph.TimeStamp())
}
