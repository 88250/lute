// Lute - A structural markdown engine.
// Copyright (C) 2019, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package lute

import (
	"fmt"
	"testing"
)

type parseTest struct {
	name   string
	input  string
	result string // what the user would see in an error message.
}

var parseTests = []parseTest{
	{"code", "l`u`te", ``},
	{"str", "lute", ``},
	//{"empty", "", noError, ``},
}

func testParse(t *testing.T) {
	for _, test := range parseTests {
		tree, err := Parse(test.name, test.input)
		if nil != err {
			t.Errorf("%q: unexpected error: %v", test.name, err)
		}

		fmt.Printf("%+v\n", tree)
	}
}

func TestParse(t *testing.T) {
	testParse(t)
}
