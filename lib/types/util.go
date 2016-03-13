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

package types

import (
	"encoding/xml"
	"fmt"
	"log"
	"strconv"

	"github.com/marklaczynski/acidbath/lib/mjlog"
)

var (
	logInfo  = log.New(mjlog.CreateInfoFile(), "INFO  [types]: ", log.LstdFlags|log.Lshortfile)
	logDebug = log.New(mjlog.CreateDebugFile(), "DEBUG [types]: ", log.LstdFlags|log.Lshortfile)
	logError = log.New(mjlog.CreateErrorFile(), "ERROR [types]: ", log.LstdFlags|log.Lshortfile)
)

//XMLFloat64 is used to parse float64 values from an xml message where it's possible the value may be nil without failing
type XMLFloat64 float64

//XMLInt64 is used to parse float64 values from an xml message where it's possible the value may be nil without failing
type XMLInt64 int64

//XMLUInt64 is used to parse float64 values from an xml message where it's possible the value may be nil without failing
type XMLUInt64 uint64

//XMLByte is used to parse float64 values from an xml message where it's possible the value may be nil without failing
type XMLByte byte

//UnmarshalXML will return nil if there's no value for a float64 element
func (f *XMLFloat64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return fmt.Errorf("XMLFloat64 should come in as a string, got %s", start)
	}

	//logDebug.Printf("s is getting parsed as: :%s:\n", s)
	if s == "" {
		f = nil
		return nil
	}

	tmpFloat, err := strconv.ParseFloat(s, 64)
	//logDebug.Printf("tmpFloat is getting parsed as: :%f:\n", tmpFloat)
	if err != nil {
		return fmt.Errorf("invalid float64 number %s", s)
	}

	*f = XMLFloat64(tmpFloat)
	//logDebug.Printf("f is assigned as: :%f:\n", *f)
	return nil
}

//UnmarshalXML will return nil if there's no value for a int64 element
func (i *XMLInt64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return fmt.Errorf("XMLInt64 should come in as a string, got %s", start)
	}

	//logDebug.Printf("s is getting parsed as: :%s:\n", s)
	if s == "" {
		i = nil
		return nil
	}

	tmpInt, err := strconv.ParseInt(s, 10, 64)
	//logDebug.Printf("tmpInt is getting parsed as: :%f:\n", tmpInt)
	if err != nil {
		return fmt.Errorf("invalid float64 number %s", s)
	}

	*i = XMLInt64(tmpInt)
	//logDebug.Printf("i is assigned as: :%d:\n", *i)
	return nil
}

//UnmarshalXML will return nil if there's no value for a uint64 element
func (i *XMLUInt64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return fmt.Errorf("XMLUInt64 should come in as a string, got %s", start)
	}

	//logDebug.Printf("s is getting parsed as: :%s:\n", s)
	if s == "" {
		i = nil
		return nil
	}

	tmpInt, err := strconv.ParseUint(s, 10, 64)
	//logDebug.Printf("tmpInt is getting parsed as: :%f:\n", tmpInt)
	if err != nil {
		return fmt.Errorf("invalid float64 number %s", s)
	}

	*i = XMLUInt64(tmpInt)
	//logDebug.Printf("i is assigned as: :%d:\n", *i)
	return nil
}

//UnmarshalXML will return nil if there's no value for a byte element
func (b *XMLByte) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return fmt.Errorf("byte should come in as a string, got %s", start)
	}

	//logDebug.Printf("s is getting parsed as: :%s:\n", s)
	if s == "" {
		b = nil
		return nil
	}

	tmpByte, err := strconv.ParseInt(s, 10, 8)
	//logDebug.Printf("tmpByte is getting parsed as: :%f:\n", tmpByte)
	if err != nil {
		return fmt.Errorf("invalid float64 number %s", s)
	}

	*b = XMLByte(tmpByte)
	//logDebug.Printf("b is assigned as: :%d:\n", *b)
	return nil
}
