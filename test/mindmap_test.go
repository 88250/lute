// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"testing"

	"github.com/88250/lute"
)

var mindmapTests = []parseTest{

	{"1", "```mindmap\n* f\\\\\n```", "<pre><code data-code=\"%7B%22name%22:%20%22f%22%7D\" class=\"language-mindmap\">* f\\\\\n</code></pre>\n"},
	{"0", "```mindmap\n* foo\n  * bar\n  * baz\n```", "<pre><code data-code=\"%7B%22name%22:%20%22foo%22,%20%22children%22:%20%5B%7B%22name%22:%20%22bar%22%7D,%20%7B%22name%22:%20%22baz%22%7D%5D%7D\" class=\"language-mindmap\">* foo\n  * bar\n  * baz\n</code></pre>\n"},
}

func TestMindmap(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range mindmapTests {
		result:= luteEngine.MarkdownStr(test.name, test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}
