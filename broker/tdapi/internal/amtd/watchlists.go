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

//Watchlists represents the TD OptionChain structure
type Watchlists struct {
	XMLName          xml.Name `xml:"amtd"`
	Error                     //inline struct
	WatchlistResults watchlistResultsXML
}

type watchlistResultsXML struct {
	XMLName   xml.Name       `xml:"watchlist-result"`
	Error     string         `xml:"error"`
	AccountID string         `xml:"account-id"`
	Watchlist []watchlistXML `xml:"watchlist"`
}

type watchlistXML struct {
	XMLName    xml.Name       `xml:"watchlist"`
	Name       string         `xml:"name"`
	ID         types.XMLInt64 `xml:"id"`
	SymbolList symbolListXML  `xml:"symbol-list"`
}

type symbolListXML struct {
	XMLName        xml.Name           `xml:"symbol-list"`
	WatchedSymbols []watchedSymbolXML `xml:"watched-symbol"`
}

type watchedSymbolXML struct {
	XMLName      xml.Name             `xml:"watched-symbol"`
	Quantity     types.XMLFloat64     `xml:"quantity"`
	Security     watchlistSecurityXML `xml:"security"`
	PositionType string               `xml:"position-type"` //enum LONG or SHORT
	AveragePrice financial.Money      `xml:"average-price"`
	Commission   financial.Money      `xml:"commission"`
	OpenDate     string               `xml:"open-date"`
}

type watchlistSecurityXML struct {
	XMLName              xml.Name `xml:"security"`
	Symbol               string   `xml:"symbol"`
	SymbolWithTypePrefix string   `xml:"symbol-with-type-prefix"`
	Description          string   `xml:"description"`
	AssetType            string   `xml:"asset-type"` // enum
}
