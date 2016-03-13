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

package tdstream

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io"
	"log"
	"math/big"
	"time"

	"github.com/marklaczynski/acidbath/broker/tdapi/tdstream/acctactivityfield"
	"github.com/marklaczynski/acidbath/broker/tdapi/tdstream/optrequestfield"
	"github.com/marklaczynski/acidbath/broker/tdapi/tdstream/quoterequestfield"
	"github.com/marklaczynski/acidbath/dm/optionchain/option"
	"github.com/marklaczynski/acidbath/dm/ordermessage"
	"github.com/marklaczynski/acidbath/lib/financial"
	"github.com/marklaczynski/acidbath/lib/mjlog"
	"github.com/marklaczynski/acidbath/lib/orderconst"
)

var (
	logInfo  = log.New(mjlog.CreateInfoFile(), "INFO  [tdstream]: ", log.LstdFlags|log.Lshortfile)
	logDebug = log.New(mjlog.CreateDebugFile(), "DEBUG [tdstream]: ", log.LstdFlags|log.Lshortfile)
	logError = log.New(mjlog.CreateErrorFile(), "ERROR [tdstream]: ", log.LstdFlags|log.Lshortfile)
)

//Decoder holds stream reader information
type Decoder struct {
	//FUTURE: test this as just an io.Reader... i think it should work
	reader *bufio.Reader
}

//NewDecoder returns a new Decoder
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{bufio.NewReader(r)}

}

//DecodeHeartbeat decodes a TD Heartbeat message
func (d *Decoder) DecodeHeartbeat() {
	logDebug.Printf("Heartbeat\n")
	subType, err := d.reader.ReadByte()
	if err != nil {
		logError.Printf("Error reading subtype: %s\n", err)
	}
	switch subType {
	case 'T':
		//read time
		var t int64
		err := binary.Read(d.reader, binary.BigEndian, &t)
		if err != nil {
			logError.Printf("binary.Read failed: ", err)
		}
		//td sends millisecond, Unix() expects seconds
		logDebug.Printf("Time: %s\n", time.Unix(t/1000, 0))
		return

	case 'H':
		return
	default:
		logError.Printf("Something went wrong\n")
		return
	}
}

/*
	Design decision:
	I'm assuming that I can correctly read from the td stream. If I happen to get
	an error of any type (including EOF) I will panic, because that means I must have
	programmed something incorrectly. I should have an endless stream of bytes, until
	I explicitly close the stream
	In the future, this might become a recoverable error, where I close the stream,
	and re-open it; however, for right now that logic is too much effort.

	Background:
	I ran into a scenario where my parseQuote function was looking for the ending
	delimintor 0xFF; however, the payload I passed to it was just shy of having it.
	Therefore, my code paniced, because there was an EOF error. I think this is
	suitable for now, because these panics mean there is something coded incorrectly.
	I now am making all the parseQuote/parseOption/etc function be able to handle the
	ending delimitor by adding 1 byte to payloadLen, because it should be the very next
	byte, and the payloadLen only returns the size of the payload itself. It's a pain
	to keep track of how many bytes I've read and compare it to payloadLen (although
	maybe if there's ever a redesign, this could be back on the table). For now the
	executive decision is to manually add 1 more byte to payloadLen
*/

// These functions should become public library for myself
func ReadBool(r io.Reader) bool {
	if ReadInt8(r) == 0 {
		return false
	}
	return true
}

func ReadInt8(r io.Reader) int8 {
	var i int8

	err := binary.Read(r, binary.BigEndian, &i)
	if err != nil {
		logError.Printf("binary.Read failed: ", err)
		panic("Unexpected failure")
	}

	return i
}

func ReadInt16(r io.Reader) int16 {
	var i int16
	err := binary.Read(r, binary.BigEndian, &i)
	if err != nil {
		logError.Printf("binary.Read failed: ", err)
		panic("Unexpected failure")
	}

	return i
}

func ReadInt32(r io.Reader) int32 {
	var i int32
	err := binary.Read(r, binary.BigEndian, &i)
	if err != nil {
		logError.Printf("binary.Read failed: ", err)
		panic("Unexpected failure")
	}

	return i
}

func ReadInt64(r io.Reader) int64 {
	var i int64
	err := binary.Read(r, binary.BigEndian, &i)
	if err != nil {
		logError.Printf("binary.Read failed: ", err)
		panic("Unexpected failure")
	}

	return i
}

func ReadString(r io.Reader, length int) string {
	logDebug.Printf("string len: %d\n", length)

	str := make([]byte, length, length)
	err := binary.Read(r, binary.BigEndian, &str)
	if err != nil {
		logError.Printf("binary.Read failed: ", err)
		panic("Unexpected failure")
	}

	return string(str)

}

func ReadFloat32(r io.Reader) float32 {
	var f float32
	err := binary.Read(r, binary.BigEndian, &f)
	if err != nil {
		logError.Printf("binary.Read failed: ", err)
		panic("Unexpected failure")
	}

	return f
}

func ReadFloat64(r io.Reader) float64 {
	var f float64
	err := binary.Read(r, binary.BigEndian, &f)
	if err != nil {
		logError.Printf("binary.Read failed: ", err)
		panic("Unexpected failure")
	}

	return f
}

//DecodeHeader reads TD header from the stream and returns its value. In case there's an EOF error, it returns 'X'
func (d *Decoder) DecodeHeader() byte {
	logDebug.Printf("DecodeHeader\n")
	header, err := d.reader.ReadByte()

	if err == io.EOF {
		return 'X'
	}

	if err != nil {
		logError.Printf("Error reading header: %s\n", err)
		return 'Y'
	}

	logDebug.Printf("parsed header: %c\n", header)

	return header
}

const delimiterSizeBytes = 1

//SidHandlers is a stuct to hold callback functions
type SidHandlers struct {
	OptionCallback          UpdateOptionAction
	AccountActivityCallback AcctActivityAction
}

//UpdateOptionAction is the function that is called once option data is parsed from the stream
type UpdateOptionAction func(newOptionData *option.Option)

//AcctActivityAction is the function that is called once an account activity (order update) comes in
type AcctActivityAction func(message *ordermessage.Message)

//DecodeCommonStreamingHeader parses the Common Streaming Header from the TD stream
func (d *Decoder) DecodeCommonStreamingHeader(sh *SidHandlers) {
	logDebug.Printf("DecodeCommonStreamingHeader\n")
	payloadLen := ReadInt16(d.reader) + delimiterSizeBytes // this is a difference from snapshot... so common code will need to adjust for this
	logDebug.Printf("Payload Length: %d\n", payloadLen)

	payload := make([]byte, payloadLen, payloadLen)
	err := binary.Read(d.reader, binary.BigEndian, &payload)
	if err != nil {
		logError.Printf("binary.Read failed: ", err)
	}

	// parse the actual payload
	payloadReader := bytes.NewReader(payload)
	parsePayload(payloadReader, sh)

	// parseEnding Delim
	endingDelim := byte(ReadInt8(d.reader))
	if endingDelim != 0x0A {
		logError.Printf("Error reading endingDelim: %s\n", err)
	}
}

//DecodeSnapshotResponse parses the Snapshot Response from the TD stream
func (d *Decoder) DecodeSnapshotResponse(sh *SidHandlers) {
	logDebug.Printf("DecodeSnapshotResponse\n")
	snapshotLen := ReadInt16(d.reader)
	logDebug.Printf("SnapshotID Len: %d\n", snapshotLen)

	//parse SID string (this is just the way Snapshot Response does it)
	//i don't think i need to make any decision
	sidString := ReadString(d.reader, int(snapshotLen))
	logDebug.Printf("SnapshotID: %s\n", sidString)

	// reused code from decode common streaming header
	payloadLen := ReadInt32(d.reader) + delimiterSizeBytes // this is a difference from snapshot... so common code will need to adjust for this
	logDebug.Printf("Payload Length: %d\n", payloadLen)

	payload := make([]byte, payloadLen, payloadLen)
	err := binary.Read(d.reader, binary.BigEndian, &payload)
	if err != nil {
		logError.Printf("binary.Read failed: ", err)
	}

	// parse the actual payload
	payloadReader := bytes.NewReader(payload)
	parsePayload(payloadReader, sh)

	// parseEnding Delim
	endingDelim := byte(ReadInt8(d.reader))
	if endingDelim != 0x0A {
		logError.Printf("Error reading endingDelim: %s\n", err)
	}
	// End resued code

}

func parsePayload(r io.Reader, sh *SidHandlers) {
	//payload starts with SID
	sid := ReadInt16(r)

	logDebug.Printf("SID: %d: %s\n", sid, StreamingID(sid))

	// once i have sid, i switch on which way to parse rest of data
	switch StreamingID(sid) {
	case Quote:
		parseQuote(r)
	case TimeSale:
	case Response:
		parseResponse(r) // see STREAMER SERVER in documentation
	case Option:
		parseOption(r, sh.OptionCallback)
	case ActivesNYSE:
	case ActivesNASDAQ:
	case ActivesOTCBB:
	case ActivesOptions:
	case News:
	case NewsHistory:
	case AdapNASDAQ:
	case NYSEBook:
	case NYSEChart:
	case NASDAQChart:
	case OpraBook:
	case IndexChart:
	case TotalView:
	case AcctActivity:
		parseAcctActivity(r, sh.AccountActivityCallback)
	case Chart:
	case StreamerServer:
	default:
		logError.Printf("Cannot handle SID: :%d:\n", sid)

	}
}

func parseAcctActivity(r io.Reader, callback AcctActivityAction) {
	logDebug.Printf("parseAcctActivity\n")
	//while column # != 0xFF continue
	buf := ReadInt8(r)
	var key, acctNum, messageType, data string

	for byte(buf) != delimiter {
		columnNum := acctactivityfield.AcctActivityNumber(buf)
		logDebug.Printf("Column num: %d: %s", columnNum, acctactivityfield.AcctActivityNumber(columnNum))

		switch columnNum {
		case acctactivityfield.SubscriptionKey:
			key = ReadString(r, int(ReadInt16(r)))
			logDebug.Printf("key: %s\n", key)
		case acctactivityfield.AccountNumber:
			acctNum = ReadString(r, int(ReadInt16(r)))
			logDebug.Printf("account number: %s\n", acctNum)
		case acctactivityfield.MessageType:
			messageType = ReadString(r, int(ReadInt16(r)))
			logDebug.Printf("message type: %s\n", messageType)
		case acctactivityfield.MessageData:
			data = ReadString(r, int(ReadInt16(r)))
			logDebug.Printf("message data: %s\n", data)

			if len(data) > 0 {
				logDebug.Printf("about to parse data\n")
				switch messageType {
				case string(acctactivityfield.Subscribed):
					//nil

				case string(acctactivityfield.Error):
					//txt
					logDebug.Printf("MESSAGTYPE: %s\n", messageType)

				case string(acctactivityfield.UrOut):
					logDebug.Printf("MESSAGTYPE: %s\n", messageType)
					//xml
					var msg acctactivityfield.UROUTMessage
					err := xml.Unmarshal([]byte(data), &msg)
					if err != nil {
						logError.Printf("Error unmarshaling response: %s\n", err)
					}

					orderMsg := ordermessage.New(msg.Order.OrderKey, orderconst.OrderOut)
					callback(orderMsg)

				case string(acctactivityfield.OrderCancelReplaceRequest):
					logDebug.Printf("MESSAGTYPE: %s\n", messageType)
					//xml
					var msg acctactivityfield.OrderCancelReplaceRequestMessage
					err := xml.Unmarshal([]byte(data), &msg)
					if err != nil {
						logError.Printf("Error unmarshaling response: %s\n", err)
					}

					orderMsg := ordermessage.New(msg.Order.OrderKey, orderconst.OrderCancelReplace)
					callback(orderMsg)

				case string(acctactivityfield.BrokenTrade):
					logDebug.Printf("MESSAGTYPE: %s\n", messageType)
					//xml
					var msg acctactivityfield.BrokenTradeMessage
					err := xml.Unmarshal([]byte(data), &msg)
					if err != nil {
						logError.Printf("Error unmarshaling response: %s\n", err)
					}

					orderMsg := ordermessage.New(msg.Order.OrderKey, orderconst.OrderBroken)
					callback(orderMsg)

				case string(acctactivityfield.ManualExecution):
					logDebug.Printf("MESSAGTYPE: %s\n", messageType)
					//xml
					var msg acctactivityfield.ManualExecutionMessage
					err := xml.Unmarshal([]byte(data), &msg)
					if err != nil {
						logError.Printf("Error unmarshaling response: %s\n", err)
					}

					orderMsg := ordermessage.New(msg.Order.OrderKey, orderconst.OrderManualExecution)
					callback(orderMsg)

				case string(acctactivityfield.OrderActivation):
					logDebug.Printf("MESSAGTYPE: %s\n", messageType)
					//xml
					var msg acctactivityfield.OrderActivationMessage
					err := xml.Unmarshal([]byte(data), &msg)
					if err != nil {
						logError.Printf("Error unmarshaling response: %s\n", err)
					}

					orderMsg := ordermessage.New(msg.Order.OrderKey, orderconst.OrderActivation)
					callback(orderMsg)

				case string(acctactivityfield.OrderCancelRequest):
					logDebug.Printf("MESSAGTYPE: %s\n", messageType)
					//xml
					var msg acctactivityfield.OrderCancelRequestMessage
					err := xml.Unmarshal([]byte(data), &msg)
					if err != nil {
						logError.Printf("Error unmarshaling response: %s\n", err)
					}

					orderMsg := ordermessage.New(msg.Order.OrderKey, orderconst.OrderCancel)
					callback(orderMsg)

				case string(acctactivityfield.OrderEntryRequest):
					logDebug.Printf("MESSAGTYPE: %s\n", messageType)
					//xml

					var msg acctactivityfield.OrderEntryRequestMessage
					err := xml.Unmarshal([]byte(data), &msg)
					if err != nil {
						logError.Printf("Error unmarshaling response: %s\n", err)
					}

					orderMsg := ordermessage.New(msg.Order.OrderKey, orderconst.OrderEntry)
					callback(orderMsg)

				case string(acctactivityfield.OrderFill):
					logDebug.Printf("MESSAGTYPE: %s\n", messageType)
					//xml
					var msg acctactivityfield.OrderFillMessage
					err := xml.Unmarshal([]byte(data), &msg)
					if err != nil {
						logError.Printf("Error unmarshaling response: %s\n", err)
					}

					orderMsg := ordermessage.New(msg.Order.OrderKey, orderconst.OrderFill)
					callback(orderMsg)

				case string(acctactivityfield.OrderPartialFill):
					logDebug.Printf("MESSAGTYPE: %s\n", messageType)
					//xml
					var msg acctactivityfield.OrderPartialFillMessage
					err := xml.Unmarshal([]byte(data), &msg)
					if err != nil {
						logError.Printf("Error unmarshaling response: %s\n", err)
					}

					orderMsg := ordermessage.New(msg.Order.OrderKey, orderconst.OrderPartialFill)
					callback(orderMsg)

				case string(acctactivityfield.OrderRejection):
					logDebug.Printf("MESSAGTYPE: %s\n", messageType)
					//xml
					var msg acctactivityfield.OrderRejectionMessage
					err := xml.Unmarshal([]byte(data), &msg)
					if err != nil {
						logError.Printf("Error unmarshaling response: %s\n", err)
					}

					orderMsg := ordermessage.New(msg.Order.OrderKey, orderconst.OrderRejection)
					callback(orderMsg)

				case string(acctactivityfield.TooLateToCancel):
					logDebug.Printf("MESSAGTYPE: %s\n", messageType)
					//xml
					var msg acctactivityfield.TooLateToCancelMessage
					err := xml.Unmarshal([]byte(data), &msg)
					if err != nil {
						logError.Printf("Error unmarshaling response: %s\n", err)
					}

					orderMsg := ordermessage.New(msg.Order.OrderKey, orderconst.OrderTooLateToCancel)
					callback(orderMsg)
				}
			}
		}

		buf = ReadInt8(r)
	}
}

const delimiter = 0xFF

func parseQuote(r io.Reader) {
	logDebug.Printf("parseQuote\n")

	//while column # != 0xFF continue
	buf := ReadInt8(r)
	for byte(buf) != delimiter {
		columnNum := quoterequestfield.QuoteColumnNumber(buf)
		logDebug.Printf("Column num: %d: %s", columnNum, quoterequestfield.QuoteColumnNumber(columnNum))
		switch columnNum {
		case quoterequestfield.Symbol:
			logDebug.Printf("Symbol: %s", ReadString(r, int(ReadInt16(r))))
		case quoterequestfield.Bid:
			bid := ReadFloat32(r)
			logDebug.Printf("Bid: %.2f", bid)
		case quoterequestfield.Ask:
			ask := ReadFloat32(r)
			logDebug.Printf("Ask: %.2f", ask)
		case quoterequestfield.Last:
			ReadFloat32(r)
		case quoterequestfield.BidSize:
			ReadInt32(r)
		case quoterequestfield.AskSize:
			ReadInt32(r)
		case quoterequestfield.BidID:
			// char in td terminology
			ReadInt16(r)
		case quoterequestfield.AskID:
			// char in td terminology
			ReadInt16(r)
		case quoterequestfield.Volume:
			// Long in td terminlogy
			ReadInt64(r)
		case quoterequestfield.LastSize:
			ReadInt32(r)
		case quoterequestfield.TradeTime:
			ReadInt32(r)
		case quoterequestfield.QuoteTime:
			ReadInt32(r)
		case quoterequestfield.High:
			ReadFloat32(r)
		case quoterequestfield.Low:
			ReadFloat32(r)
		case quoterequestfield.Tick:
			// char in td terminology
			ReadInt16(r)
		case quoterequestfield.Close:
			ReadFloat32(r)
		case quoterequestfield.EXChange:
			ReadInt16(r)
		case quoterequestfield.Marginable:
			ReadBool(r)
		case quoterequestfield.Shortable:
			ReadBool(r)
		case quoterequestfield.QuoteDate:
			// # days since 1/1/1970
			ReadInt32(r)
		case quoterequestfield.TradeDate:
			ReadInt32(r)
		case quoterequestfield.Volatility:
			ReadFloat32(r)
		case quoterequestfield.Description:
			ReadString(r, int(ReadInt16(r)))
		case quoterequestfield.TradeID:
			ReadInt16(r)
		case quoterequestfield.Digits:
			ReadInt32(r)
		case quoterequestfield.Open:
			ReadFloat32(r)
		case quoterequestfield.Change:
			ReadFloat32(r)
		case quoterequestfield.WeekHigh52:
			ReadFloat32(r)
		case quoterequestfield.WeekLow52:
			ReadFloat32(r)
		case quoterequestfield.PERatio:
			ReadFloat32(r)
		case quoterequestfield.DividendAmt:
			ReadFloat32(r)
		case quoterequestfield.DividendYield:
			ReadFloat32(r)
		case quoterequestfield.Nav:
			ReadFloat32(r)
		case quoterequestfield.Fund:
			ReadFloat32(r)
		case quoterequestfield.ExchangeName:
			ReadString(r, int(ReadInt16(r)))
		case quoterequestfield.DividendDate:
			ReadString(r, int(ReadInt16(r)))
		case quoterequestfield.LastMarketHours:
			ReadFloat32(r)
		case quoterequestfield.LastSizeMarketHours:
			ReadInt32(r)
		case quoterequestfield.TradeDateMarketHours:
			ReadInt32(r)
		case quoterequestfield.TradeTimeMarketHours:
			ReadInt32(r)
		case quoterequestfield.ChangeMarketHours:
			ReadFloat32(r)
		case quoterequestfield.IsRegularMarketQuote:
			ReadBool(r)
		case quoterequestfield.IsRegularMarketTrade:
			ReadBool(r)

		}

		buf = ReadInt8(r)
		logDebug.Printf("Buf byte %x\n", buf)
	}
	logDebug.Printf("Exit for column loop\n")

}

func parseOption(r io.Reader, callback UpdateOptionAction) {
	logDebug.Printf("parseOption\n")

	var newOptionData *option.Option = option.NewNilOption()

	buf := ReadInt8(r)
	for byte(buf) != delimiter {
		columnNum := optrequestfield.OptionColumnNumber(buf)
		logDebug.Printf("Column num: %d: %s", columnNum, optrequestfield.OptionColumnNumber(columnNum))
		switch columnNum {
		case optrequestfield.Symbol:
			optSymbol := ReadString(r, int(ReadInt16(r)))
			logDebug.Printf("Symbol: %s\n", optSymbol)
			newOptionData.SetOptionTickerSymbol(optSymbol)
		case optrequestfield.Contract:
			ReadString(r, int(ReadInt16(r)))
		case optrequestfield.Bid:
			bid := financial.Money{(&big.Rat{}).SetFloat64(float64(ReadFloat32(r)))}
			logDebug.Printf("Bid: %.2f", bid)
			logDebug.Printf("option: %#v ", newOptionData)
			// this is needed, because it seems like even though we have UNSUBS from all options,
			// there are old lingering options still streaming.
			if newOptionData != nil {
				newOptionData.SetBid(bid)
			}
		case optrequestfield.Ask:
			ask := financial.Money{(&big.Rat{}).SetFloat64(float64(ReadFloat32(r)))}
			logDebug.Printf("Ask: %.2f", ask)
			logDebug.Printf("option: %#v ", newOptionData)
			if newOptionData != nil {
				newOptionData.SetAsk(ask)
			}
		case optrequestfield.Last:
			last := financial.Money{(&big.Rat{}).SetFloat64(float64(ReadFloat32(r)))}
			logDebug.Printf("Last: %.2f", last)
			logDebug.Printf("option: %#v ", newOptionData)
			if newOptionData != nil {
				newOptionData.SetLast(last)
			}

		case optrequestfield.High:
			ReadFloat32(r)
		case optrequestfield.Low:
			ReadFloat32(r)
		case optrequestfield.Close:
			ReadFloat32(r)
		case optrequestfield.Volume:
			ReadInt64(r)
		case optrequestfield.OpenInterest:
			ReadInt32(r)
		case optrequestfield.Volatility:
			ReadFloat32(r)
		case optrequestfield.QuoteTime:
			ReadInt32(r)
		case optrequestfield.TradeTime:
			ReadInt32(r)
		case optrequestfield.InTheMoney:
			ReadFloat32(r)
		case optrequestfield.QuoteDate:
			ReadInt32(r)
		case optrequestfield.TradeDate:
			ReadInt32(r)
		case optrequestfield.Year:
			ReadInt32(r)
		case optrequestfield.Multiplier:
			ReadFloat32(r)
		case optrequestfield.Open:
			ReadFloat32(r)
		case optrequestfield.BidSize:
			ReadInt32(r)
		case optrequestfield.AskSize:
			ReadInt32(r)
		case optrequestfield.LastSize:
			ReadInt32(r)
		case optrequestfield.Change:
			ReadFloat32(r)
		case optrequestfield.Strike:
			ReadFloat32(r)
		case optrequestfield.ContractType:
			ReadInt16(r) //char
		case optrequestfield.Underlying:
			ReadString(r, int(ReadInt16(r)))
		case optrequestfield.Month:
			ReadInt32(r)
		case optrequestfield.Note:
			ReadString(r, int(ReadInt16(r)))
		case optrequestfield.TimeValue:
			ReadFloat32(r)
		case optrequestfield.DaysToExp:
			ReadInt32(r)
		case optrequestfield.DeltaIndex:
			delta := float64(ReadFloat32(r))
			logDebug.Printf("delta: %.2f", delta)
			logDebug.Printf("option: %#v ", newOptionData)
			if newOptionData != nil {
				newOptionData.SetDelta(delta)
			}
		case optrequestfield.GammaIndex:
			gamma := float64(ReadFloat32(r))
			logDebug.Printf("gamma: %.2f", gamma)
			logDebug.Printf("option: %#v ", newOptionData)
			if newOptionData != nil {
				newOptionData.SetGamma(gamma)
			}
		case optrequestfield.ThetaIndex:
			theta := float64(ReadFloat32(r))
			logDebug.Printf("theta: %.2f", theta)
			logDebug.Printf("option: %#v ", newOptionData)
			if newOptionData != nil {
				newOptionData.SetTheta(theta)
			}
		case optrequestfield.VegaIndex:
			vega := float64(ReadFloat32(r))
			logDebug.Printf("vega: %.2f", vega)
			logDebug.Printf("option: %#v ", newOptionData)
			if newOptionData != nil {
				newOptionData.SetVega(vega)
			}
		case optrequestfield.RhoIndex:
			ReadFloat32(r)
		}
		buf = ReadInt8(r)
		logDebug.Printf("Buf byte %x\n", buf)
	}

	callback(newOptionData)
	logDebug.Printf("Exit for column loop\n")
}

func parseResponse(r io.Reader) {
	logDebug.Printf("parseResponse\n")

	columnNum := ReadInt8(r)
	logDebug.Printf("Column num: %d", columnNum)

	sid := ReadInt16(r)
	logDebug.Printf("Service ID: %d\n", sid)

	columnNum = ReadInt8(r)
	logDebug.Printf("Column num: %d", columnNum)

	returnCode := ReadInt16(r)
	logDebug.Printf("Return Code: %d", returnCode)

	columnNum = ReadInt8(r)
	logDebug.Printf("Column num: %d", columnNum)

	descriptionLen := ReadInt16(r)
	logDebug.Printf("Description Len: %d", descriptionLen)

	description := ReadString(r, int(descriptionLen))
	logDebug.Printf("Description: %s", description)

	// parseLastField
	lastField := byte(ReadInt8(r))
	if lastField != delimiter {
		logError.Printf("Error reading lastfield\n")
	}

}

//StreamingID (aka SID) is an enum type that represents the what is being streamed (ie QUOTE, OPTION, etc)
type StreamingID int16

//The various constants that represent the Streaming ID codes specified by TD API
const (
	Quote          StreamingID = 1
	TimeSale                   = 5
	Response                   = 10
	Option                     = 18
	ActivesNYSE                = 23
	ActivesNASDAQ              = 25
	ActivesOTCBB               = 26
	ActivesOptions             = 35
	News                       = 27
	NewsHistory                = 28
	AdapNASDAQ                 = 62
	NYSEBook                   = 81
	NYSEChart                  = 82
	NASDAQChart                = 83
	OpraBook                   = 84
	IndexChart                 = 85
	TotalView                  = 87
	AcctActivity               = 90
	Chart                      = 91
	StreamerServer             = 100
)

func (id StreamingID) String() string {
	switch id {
	case Quote:
		return "QUOTE"
	case TimeSale:
		return "TIMESALE"
	case Response:
		return "RESPONSE"
	case Option:
		return "OPTION"
	case ActivesNYSE:
		return "ACTIVES_NYSE"
	case ActivesNASDAQ:
		return "ACTIVES_NASDAQ"
	case ActivesOTCBB:
		return "ACTIVES_OTCBB"
	case ActivesOptions:
		return "ACTIVES_OPTIONS"
	case News:
		return "NEWS"
	case NewsHistory:
		return "NEWS_HISTORY"
	case AdapNASDAQ:
		return "ADAP_NASDAQ"
	case NYSEBook:
		return "NYSE_BOOK"
	case NYSEChart:
		return "NYSE_CHART"
	case NASDAQChart:
		return "NASDAQ_CHART"
	case OpraBook:
		return "OPRA_BOOK"
	case IndexChart:
		return "INDEX_CHART"
	case TotalView:
		return "TOTAL_VIEW"
	case AcctActivity:
		return "ACCT_ACTIVITY"
	case Chart:
		return "CHART"
	case StreamerServer:
		return "STREAMER_SERVER"

	}
	return ""
}
