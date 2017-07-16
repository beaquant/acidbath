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

package registry

import (
	"log"
	"sync"

	"github.com/marklaczynski/acidbath/lib/mjlog"
)

var (
	logInfo  = log.New(mjlog.CreateInfoFile(), "INFO  [types]: ", log.LstdFlags|log.Lshortfile)
	logDebug = log.New(mjlog.CreateDebugFile(), "DEBUG [types]: ", log.LstdFlags|log.Lshortfile)
	logError = log.New(mjlog.CreateErrorFile(), "ERROR [types]: ", log.LstdFlags|log.Lshortfile)
	instance *BadSymbols
	once     sync.Once
)

//Dates is a slice of time.Times
type BadSymbols struct {
	isBad map[string]bool
	//optionUpdateChans map[string]chan *option.Option
}

func GetBadSymbols() *BadSymbols {
	once.Do(func() {
		instance = &BadSymbols{isBad: make(map[string]bool)}
	})

	return instance
}

func (bs *BadSymbols) AddBadSymbol(symbol string) {
	bs.isBad[symbol] = true
	//if bs.IsBadSymbol(symbol) == false {
	//	bs.symbols = append(bs.symbols, symbol)
	//}
}

func (bs *BadSymbols) RemoveBadSymbol(symbol string) {
	delete(bs.isBad, symbol)
	/*
		for idx, currSymbol := range bs.symbols {
			if currSymbol == symbol {
				bs.symbols[idx] = bs.symbols[len(bs.symbols)-1]
				bs.symbols = bs.symbols[:len(bs.symbols)-1]
			}
		}
	*/
}

func (bs *BadSymbols) IsBadSymbol(symbol string) bool {
	return bs.isBad[symbol] == true
	/*
		for _, currSymbol := range bs.symbols {
			if currSymbol == symbol {
				return true
			}
		}
		return false
	*/
}

func (bs *BadSymbols) String() string {
	return_string := ""
	for currKey, currValue := range bs.isBad {
		if currValue == true {
			return_string = return_string + "," + currKey
		}
	}
	return return_string
	//return strings.Join(bs.symbols, ",")
}

func (bs *BadSymbols) Count() int {
	return len(bs.isBad)
}
