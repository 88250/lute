// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/88250/lute"
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
	bytes, err := os.ReadFile("commonmark-spec.json")
	if nil != err {
		t.Fatalf("read spec test cases failed: " + err.Error())
	}

	var testcases []testcase
	if err = json.Unmarshal(bytes, &testcases); nil != err {
		t.Fatalf("read spec test caes failed: " + err.Error())
	}

	luteEngine := lute.New()
	luteEngine.ParseOptions.GFMTaskListItem = false
	luteEngine.ParseOptions.GFMTable = false
	luteEngine.ParseOptions.GFMAutoLink = false
	luteEngine.ParseOptions.GFMStrikethrough = false
	luteEngine.RenderOptions.SoftBreak2HardBreak = false
	luteEngine.RenderOptions.CodeSyntaxHighlight = false
	luteEngine.ParseOptions.HeadingID = false
	luteEngine.RenderOptions.HeadingID = false
	luteEngine.RenderOptions.AutoSpace = false
	luteEngine.RenderOptions.FixTermTypo = false
	luteEngine.ParseOptions.Emoji = false
	luteEngine.ParseOptions.YamlFrontMatter = false

	for _, test := range testcases {
		testName := test.Section + " " + strconv.Itoa(test.Example)
		html := luteEngine.MarkdownStr(testName, test.Markdown)
		if test.HTML != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", testName, test.HTML, html, test.Markdown)
		}
	}
}
