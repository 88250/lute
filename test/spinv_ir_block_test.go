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

var spinVditorIRBlockDOMTests = []*parseTest{

	{"15", "<p data-block=\"0\" data-node-id=\"1596459249782\">((foo \"text\")<wbr></p>\n", "<p data-block=\"0\" data-node-id=\"\">((foo &quot;text&quot;)<wbr></p>"},
	{"14", "<p data-block=\"0\" data-node-id=\"1596459249782\">((foo \"text\"))<wbr></p>\n", "<p data-block=\"0\" data-node-id=\"\"><a href=\"foo\">text</a><wbr></p>"},
	{"13", "<p data-block=\"0\" data-node-id=\"1596459249782\">((foo))<wbr></p>\n", "<p data-block=\"0\" data-node-id=\"\"><a href=\"foo\">placeholder</a><wbr></p>"},
	{"12", "((foo_bar))\n", "<p data-block=\"0\" data-node-id=\"\">((foo_bar))</p>"},
	{"11", "((000000000000000000000000000000000000000000000000000000000000000001 \"text\"))\n", "<p data-block=\"0\" data-node-id=\"\">((000000000000000000000000000000000000000000000000000000000000000001 &quot;text&quot;))</p>"},
	{"10", "((00000000000000000000000000000000000000000000000000000000000000000 \"text\"))\n", "<p data-block=\"0\" data-node-id=\"\">(<a href=\"0000000000000000000000000000000000000000000000000000000000000000\">text</a></p>"},
	{"9", "((00000000000000000000000000000000000000000000000000000000000000000))\n", "<p data-block=\"0\" data-node-id=\"\">(<a href=\"0000000000000000000000000000000000000000000000000000000000000000\">placeholder</a></p>"},
	{"8", "(( foo  \"text\" ))\n", "<p data-block=\"0\" data-node-id=\"\"><a href=\"foo\">text</a></p>"},
	{"7", "(( foo  ))\n", "<p data-block=\"0\" data-node-id=\"\"><a href=\"foo\">placeholder</a></p>"},
	{"6", "((foo text))\n", "<p data-block=\"0\" data-node-id=\"\">((foo text))</p>"},
	{"5", "((foo\"text\"))\n", "<p data-block=\"0\" data-node-id=\"\">((foo&quot;text&quot;))</p>"},
	{"4", "((foo-bar-123 \"text\"))\n", "<p data-block=\"0\" data-node-id=\"\"><a href=\"foo-bar-123\">text</a></p>"},
	{"3", "((foo-bar))\n", "<p data-block=\"0\" data-node-id=\"\"><a href=\"foo-bar\">placeholder</a></p>"},
	{"2", "((foo))\n", "<p data-block=\"0\" data-node-id=\"\"><a href=\"foo\">placeholder</a></p>"},
	{"1", "<p data-block=\"0\" data-node-id=\"1\">foo</p><p data-block=\"0\"><wbr><br></p>", "<p data-block=\"0\" data-node-id=\"\">foo</p><p data-block=\"0\" data-node-id=\"\"><wbr></p>"},
	{"0", "<p data-block=\"0\" data-node-id=\"1\">foo</p><p data-block=\"0\"><wbr><br></p>", "<p data-block=\"0\" data-node-id=\"\">foo</p><p data-block=\"0\" data-node-id=\"\"><wbr></p>"},
}

func TestSpinVditorIRBlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.BlockRef = true

	for _, test := range spinVditorIRBlockDOMTests {
		html := luteEngine.SpinVditorIRBlockDOM(test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
