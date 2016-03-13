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
package financial

import (
	"encoding/xml"
	"fmt"
	"log"
	"math/big"

	"github.com/marklaczynski/acidbath/lib/mjlog"
)

var (
	logInfo  = log.New(mjlog.CreateInfoFile(), "INFO  [financial]: ", log.LstdFlags|log.Lshortfile)
	logDebug = log.New(mjlog.CreateDebugFile(), "DEBUG [financial]: ", log.LstdFlags|log.Lshortfile)
	logError = log.New(mjlog.CreateErrorFile(), "ERROR [financial]: ", log.LstdFlags|log.Lshortfile)
)

//Money represents financial amounts, which is implemented as big.Rat
type Money struct {
	Value *big.Rat
}

//UnmarshalXML unmarshal an Money type from an XML message
func (m *Money) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return fmt.Errorf("XMLFloat64 should come in as a string, got %s", start)
	}

	if m.Value == nil {
		m.Value = big.NewRat(0, 1)
	}

	//logDebug.Printf("s is getting parsed as: :%s:\n", s)
	if s == "" {
		return nil
	}

	var success bool
	m.Value, success = m.Value.SetString(s)
	if !success {
		return fmt.Errorf("invalid rat number %s", s)
	}
	//logDebug.Printf("m.Value: :%s:\n", m.Value)

	return nil

}

func (m Money) String() string {
	return m.Value.FloatString(2)
}

type MoneyArray []Money

func (ma *MoneyArray) Value(i int) float64 {
	f, _ := (*ma)[i].Value.Float64()
	return f
}

func (ma *MoneyArray) SetValue(i int, value float64) {
	(*ma)[i].Value = (&big.Rat{}).SetFloat64(value)
}

func (ma *MoneyArray) Len() int {
	return len(*ma)
}

func (ma *MoneyArray) Less(i, j int) bool {
	return (*ma)[i].Value.Cmp((*ma)[j].Value) < 0
}

func (ma *MoneyArray) Swap(i, j int) {
	(*ma)[i], (*ma)[j] = (*ma)[j], (*ma)[i]
}
