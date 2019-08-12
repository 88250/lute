// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"encoding/json"
	"fmt"
	"github.com/b3log/lute"
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
	bytes, err := ioutil.ReadFile("commonmark-0.29-spec.json")
	if nil != err {
		t.Fatalf("read spec test cases failed: " + err.Error())
	}

	var testcases []testcase
	if err = json.Unmarshal(bytes, &testcases); nil != err {
		t.Fatalf("read spec test caes failed: " + err.Error())
	}

	luteEngine := lute.New(lute.GFM(false))

	for _, test := range testcases {
		testName := test.Section + " " + strconv.Itoa(test.Example)
		fmt.Println("Test [" + testName + "]")

		html, err := luteEngine.MarkdownStr(testName, test.Markdown)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.HTML != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", testName, test.HTML, html, test.Markdown)
		}
	}
}
