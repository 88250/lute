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
	"fmt"
	"github.com/b3log/lute"
	"testing"
)

var spaceTests = []parseTest{

	{"space0", "Lute是一款结构化的Markdown引擎，完整实现了最新的GFM / CommonMark规范，对中文语境支持更好。\n", "<p>Lute 是一款结构化的 Markdown 引擎，完整实现了最新的 GFM / CommonMark 规范，对中文语境支持更好。</p>\n"},
}

func TestAutoSpace(t *testing.T) {
	luteEngine := lute.New() // 默认已经开启自动空格优化

	for _, test := range spaceTests {
		fmt.Println("Test [" + test.name + "]")
		html, err := luteEngine.MarkdownStr(test.name, test.markdown)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.html != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.html, html, test.markdown)
		}
		t.Log(html)
	}
}
