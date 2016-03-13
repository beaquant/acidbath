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

	"github.com/marklaczynski/acidbath/lib/types"
)

//StreamerInfo represents the TD StreamerInfo structure
type StreamerInfo struct {
	XMLName      xml.Name `xml:"amtd"`
	Error                 //inline struct
	StreamerInfo streamerInfoXML
}

type streamerInfoXML struct {
	XMLName     xml.Name       `xml:"streamer-info"`
	StreamerURL string         `xml:"streamer-url"`
	Token       string         `xml:"token"`
	Timestamp   types.XMLInt64 `xml:"timestamp"`
	CDDomainID  string         `xml:"cd-domain-id"`
	Usergroup   string         `xml:"usergroup"`
	AccessLevel string         `xml:"access-level"`
	Acl         string         `xml:"acl"`
	Appid       string         `xml:"app-id"`
	Authorized  string         `xml:"authorized"`
	ErrorMsg    string         `xml:"error-msg"`
}
