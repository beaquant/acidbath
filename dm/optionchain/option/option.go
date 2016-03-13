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

//Package option implements a financial option
//reference http://en.wikipedia.org/wiki/Option_%28finance%29
package option

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/marklaczynski/acidbath/lib/date"
	"github.com/marklaczynski/acidbath/lib/financial"
)

//Option represents an option
type Option struct {
	//Quote

	// These are streaming values
	symbol           string
	strike           float64   // Comp Key
	expirationDate   time.Time // Comp Key
	multiplier       float64
	last             financial.Money
	bid              financial.Money
	ask              financial.Money
	delta            float64
	gamma            float64
	theta            float64
	vega             float64
	underlying       string // SPX, SPY, etc. Comp Key
	daysToExpiration int64

	//rho            float64
	//volume         int64
	//impliedVolatility float64 // volatility in TD
	//bidSize          int32
	//askSize          int32
	//openInterest      int32
	//contract          string //unused
	//high              financial.Money //unused
	//low               financial.Money //unused
	//closePrice        financial.Money //unused
	//quotetime         int32   //unused
	//tradetime         int32   //unused
	//inTheMoney        float64 //unused
	//quotedate         int32   //unused
	//tradedate         int32   //unused
	//year              int32   //unused
	//open              float64 //unused
	//lastSize          int32   //unused
	//change            float64 //unused
	//month             int32   //unused
	//note              string  //unused
	//timevalue         float64 //unused

	// this is from OptionChain response not sure how this maps yet
	description string

	// my own utility
	theoPrice          float64
	isOTM              bool
	optionTickerSymbol string       // this probably maps to something else
	optionType         TypeOfOption // Comp Key

	//error stuff
	err error
}

//TypeOfOption is an enum type of option (ie CALL/PUT)
type TypeOfOption int

//defines constants related to TypeOfOption
const (
	CALL TypeOfOption = iota
	PUT
)

var typeOfOptionDescription = [...]string{
	"CALL",
	"PUT",
}

func (t TypeOfOption) String() string {
	return typeOfOptionDescription[t]
}

func NewNilOption() *Option {
	o := &Option{}
	o.bid.Value = big.NewRat(0, 1)
	o.ask.Value = big.NewRat(0, 1)
	o.last.Value = big.NewRat(0, 1)
	//o.high.Value = big.NewRat(0, 1)
	//o.low.Value = big.NewRat(0, 1)
	//o.closePrice.Value = big.NewRat(0, 1)

	return o
}

// NewOption is an Option constructor
func NewOption(ul string, strike float64, expDate time.Time, optType TypeOfOption, multiplier float64) (*Option, error) {
	o := &Option{}
	o.SetUnderlying(ul)
	o.SetOptionType(optType)
	o.SetMultiplier(multiplier)
	if err := o.SetStrike(strike); err != nil {
		o.err = err
		return nil, err
	}

	if err := o.SetExpirationDate(expDate); err != nil {
		o.err = err
		return nil, err
	}

	o.bid.Value = big.NewRat(0, 1)
	o.ask.Value = big.NewRat(0, 1)
	o.last.Value = big.NewRat(0, 1)
	//o.high.Value = big.NewRat(0, 1)
	//o.low.Value = big.NewRat(0, 1)
	//o.closePrice.Value = big.NewRat(0, 1)

	return o, nil
}

func (o *Option) String() string {
	return fmt.Sprintf("Option: %s", o.OptionTickerSymbol())
}

func (o *Option) BetaWeightedDelta(targetBetaStockLastTradePrice float64, underlyingStockLastTradePrice float64, underlyingStockBeta float64) float64 {
	betaTargetLastTradePrice := targetBetaStockLastTradePrice
	currentOptionPositionStockLastTradePrice := underlyingStockLastTradePrice
	betaWeightedDelta := o.Delta() / ((betaTargetLastTradePrice / currentOptionPositionStockLastTradePrice) / underlyingStockBeta)

	return betaWeightedDelta
}

func (o *Option) Error() error {
	return o.err
}

//Bid returns the bid price
func (o *Option) Bid() financial.Money {
	return o.bid
}

//SetBid sets the bid price
func (o *Option) SetBid(bid financial.Money) error {
	if bid.Value.Cmp(big.NewRat(0, 1)) < 0 {
		o.err = errors.New("Bid cannot be less than 0")
		return o.err
	}
	o.bid.Value.Set(bid.Value)
	return nil
}

//Ask returns the ask price
func (o *Option) Ask() financial.Money {
	return o.ask
}

//SetAsk sets the ask price
func (o *Option) SetAsk(ask financial.Money) error {
	if ask.Value.Cmp(big.NewRat(0, 1)) < 0 {
		o.err = errors.New("Ask cannot be less than 0")
		return o.err
	}

	o.ask.Value.Set(ask.Value)
	return nil
}

/*
//BidSize retuns the bid size
func (o *Option) BidSize() int32 {
	return o.bidSize
}

//SetBidSize sets the bid size
func (o *Option) SetBidSize(bidSize int32) error {
	if bidSize <= 0 {
		o.err = errors.New("BidSize cannot be <= 0")
		return o.err
	}

	o.bidSize = bidSize
	return nil
}
*/

/*
//AskSize retuns the ask size
func (o *Option) AskSize() int32 {
	return o.askSize
}

//SetAskSize sets the ask size
func (o *Option) SetAskSize(askSize int32) error {
	if askSize <= 0 {
		o.err = errors.New("AskSize cannot be <= 0")
		return o.err
	}

	o.askSize = askSize
	return nil
}
*/

//Last returns the last trade price
func (o *Option) Last() financial.Money {
	return o.last
}

//SetLast sets the last trade price
func (o *Option) SetLast(last financial.Money) error {
	if last.Value.Cmp(big.NewRat(0, 1)) < 0 {
		o.err = errors.New("Last price cannot be less than 0")
		return o.err
	}

	o.last.Value.Set(last.Value)
	return nil
}

/*
//Volume returns the volume
func (o *Option) Volume() int64 {
	return o.volume
}

//SetVolume sets the volume
func (o *Option) SetVolume(volume int64) error {
	if volume < 0 {
		o.err = errors.New("Volume cannot be less than 0")
		return o.err
	}

	o.volume = volume
	return nil
}
*/

/*
//OpenInterest returns the open interest
func (o *Option) OpenInterest() int32 {
	return o.openInterest
}

//SetOpenInterest sets the open interest
func (o *Option) SetOpenInterest(openInterest int32) error {
	if openInterest < 0 {
		o.err = errors.New("Open Intereste cannot be less than 0")
		return o.err
	}

	o.openInterest = openInterest
	return nil
}
*/

//Delta returns the delta
func (o *Option) Delta() float64 {
	return o.delta
}

//SetDelta sets the delta
func (o *Option) SetDelta(delta float64) {
	o.delta = delta
}

//Vega returns the vega
func (o *Option) Vega() float64 {
	return o.vega
}

//SetVega sets the vega
func (o *Option) SetVega(vega float64) {
	o.vega = vega
}

//Gamma returns the gamma
func (o *Option) Gamma() float64 {
	return o.gamma
}

//SetGamma sets the gamma
func (o *Option) SetGamma(gamma float64) {
	o.gamma = gamma
}

//Theta returns the theta
func (o *Option) Theta() float64 {
	return o.theta
}

//SetTheta sets the theta
func (o *Option) SetTheta(theta float64) {
	o.theta = theta
}

/*
//Rho returns the rho
func (o *Option) Rho() float64 {
	return o.rho
}

//SetRho sets the rho
func (o *Option) SetRho(rho float64) {
	o.rho = rho
}
*/

//Multiplier retuns the multiplier
func (o *Option) Multiplier() float64 {
	return o.multiplier
}

//SetMultiplier sets the multiplier
func (o *Option) SetMultiplier(multiplier float64) {
	o.multiplier = multiplier
}

//TheoPrice retuns the theoretical price
func (o *Option) TheoPrice() float64 {
	return o.theoPrice
}

//SetTheoPrice sets the theoretical price
func (o *Option) SetTheoPrice(theoPrice float64) error {
	if theoPrice < 0 {
		o.err = errors.New("Theoretical Price cannot be less than 0")
		return o.err
	}

	o.theoPrice = theoPrice
	return nil
}

//IsOTM returns if the option is Out-Of-The-Money
func (o *Option) IsOTM() bool {
	return o.isOTM
}

//SetIsOTM sets if the option is Out-Of-The-Money
func (o *Option) SetIsOTM(isOTM bool) {
	o.isOTM = isOTM
}

//Description returns the description
func (o *Option) Description() string {
	return o.description
}

//SetDescription sets the description
func (o *Option) SetDescription(description string) {
	o.description = description
}

/*
//ImpliedVolatility returns the implied volatility
func (o *Option) ImpliedVolatility() float64 {
	return o.impliedVolatility
}

//SetImpliedVolatility sets the implied volatility
func (o *Option) SetImpliedVolatility(impliedVolatility float64) error {
	if impliedVolatility < 0 {
		o.err = errors.New("Implied Volatility cannot be less than 0")
		return o.err
	}

	o.impliedVolatility = impliedVolatility
	return nil
}
*/

//Symbol retuns the symbol
func (o *Option) Symbol() string {
	return o.symbol
}

//SetSymbol sets the symbol
func (o *Option) SetSymbol(symbol string) {
	o.symbol = symbol
}

//Strike returns the current strike price
func (o *Option) Strike() float64 {
	return o.strike
}

//SetStrike sets the current strike price
func (o *Option) SetStrike(strike float64) error {
	if strike < 0 {
		o.err = errors.New("Strike cannot be less than 0")
		return o.err
	}

	o.strike = strike
	return nil
}

//Underlying returns the underlying associated with the option
func (o *Option) Underlying() string {
	return o.underlying
}

//SetUnderlying sets the underlying associated with the option
func (o *Option) SetUnderlying(underlying string) {
	o.underlying = underlying
}

//DaysToExpiration returns the number of days to expiration
func (o *Option) DaysToExpiration() int64 {
	return o.daysToExpiration
}

//SetDaysToExpiration sets the number of days to expiration
func (o *Option) SetDaysToExpiration(daysToExpiration int64) error {
	if daysToExpiration < 0 {
		o.err = errors.New("Days to expiration cannot be less than 0")
		return o.err
	}

	o.daysToExpiration = daysToExpiration
	return nil
}

//ExpirationDate returns the expiration date as time.Time, but with no timestamp value
func (o *Option) ExpirationDate() time.Time {
	return o.expirationDate
}

//SetExpirationDate returns the expiration date as time.Time, but with no timestamp value
//expiration date cannot be in the past
func (o *Option) SetExpirationDate(expirationDate time.Time) error {
	expDate, err := date.New(expirationDate)
	if err != nil {
		o.err = err
		return o.err
	}

	/* for some reason this causes issues late in the even
	currDate, err := date.New(time.Now())
	if err != nil {
		return err
	}

		if expDate.Before(currDate) {
			return errors.New("Time is in past")
		}
	*/

	o.expirationDate = expDate
	return nil
}

//OptionTickerSymbol returns the option ticker symbol
func (o *Option) OptionTickerSymbol() string {
	return o.optionTickerSymbol
}

//SetOptionTickerSymbol sets the option ticker symbol
func (o *Option) SetOptionTickerSymbol(optionTickerSymbol string) {
	o.optionTickerSymbol = optionTickerSymbol
}

//OptionType returns the type of option (Call/Put)
func (o *Option) OptionType() TypeOfOption {
	return o.optionType
}

//SetOptionType sets the type of option (Call/Put)
func (o *Option) SetOptionType(optionType TypeOfOption) {
	o.optionType = optionType
}

//Copy returns a new copy/clone of the option
func (o *Option) Copy() *Option {
	dst := &Option{}
	*dst = *o
	return dst
}
