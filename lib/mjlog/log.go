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

package mjlog

import "os"

//CreateDebugFile create a debug file named "debug.log" or it appends to it if it exists
func CreateDebugFile() *os.File {
	debugFile := "./debuginfo.log"

	f, err := os.OpenFile(debugFile, os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		f, err = os.Create(debugFile)
		if err != nil {
			return nil
		}

		return f
	}

	return f
}

//CreateErrorFile create a error file named "error.log" or it appends to it if it exists
func CreateErrorFile() *os.File {
	debugFile := "./debuginfo.log"

	f, err := os.OpenFile(debugFile, os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		f, err = os.Create(debugFile)
		if err != nil {
			return nil
		}

		return f
	}

	return f
}

//CreateInfoFile create an info file named "info.log" or it appends to it if it exists
func CreateInfoFile() *os.File {
	debugFile := "./debuginfo.log"

	f, err := os.OpenFile(debugFile, os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		f, err = os.Create(debugFile)
		if err != nil {
			return nil
		}

		return f
	}

	return f
}
