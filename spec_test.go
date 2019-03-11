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
	"encoding/json"
	"io/ioutil"
	"strconv"
	"testing"
)

type testcase struct {
	EndLine   int    `json:"end_line"`
	Section   string `json:"section"`
	HTML      string `json:"html"`
	Markdown  string `json:"markdown"`
	Example   int    `json:"example"`
	StartLine int    `json:"start_line"`
}

func TestSpec(t *testing.T) {
	bytes, err := ioutil.ReadFile("commonmark-0.28-spec.json")
	if nil != err {
		t.Fatalf("read spec test cases failed: " + err.Error())
	}

	var testcases []testcase
	if err = json.Unmarshal(bytes, &testcases); nil != err {
		t.Fatalf("read spec test caes failed: " + err.Error())
	}

	for _, tc := range testcases {
		tcn := tc.Section+" "+strconv.Itoa(tc.Example)
		tree, err := Parse(tcn, tc.Markdown)
		if nil != err {
			t.Fatalf("parse [%s] failed: %s", tree.name, err.Error())
		}

		html := tree.HTML()
		if tc.HTML != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", tree.name, tc.HTML, html, tc.Markdown)
		}
	}
}
