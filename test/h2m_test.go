// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package test

import (
	"testing"

	"github.com/b3log/lute"
)

var h2mTests = []parseTest{

	{"16", "<p><em><strong>foo</strong></em></p>", "***foo***\n"},
	{"15", "<p><strong data-marker=\"__\">foo</strong></p>", "__foo__\n"},
	{"14", "<p><strong data-marker=\"**\">foo</strong></p>", "**foo**\n"},
	{"13", "<h2>foo</h2>\n<p>para<em>em</em></p>", "## foo\n\npara*em*\n"},
	{"12", "<a href=\"/bar\" title=\"baz\">foo</a>", "[foo](/bar \"baz\")\n"},
	{"11", "<img src=\"/bar\" alt=\"foo\" />", "![foo](/bar)\n"},
	{"10", "<img src=\"/bar\" />", "![](/bar)\n"},
	{"9", "<a href=\"/bar\">foo</a>", "[foo](/bar)\n"},
	{"8", "foo<br />bar", "foo\nbar\n"},
	{"7", "<code>foo</code>", "`foo`\n"},
	{"6", "<pre><code>foo</code></pre>", "```\nfoo\n```\n"},
	{"5", "<ul><li>foo</li></ul>", "* foo\n"},
	{"4", "<blockquote>foo</blockquote>", "> foo\n"},
	{"3", "<h2>foo</h2>", "## foo\n"},
	{"2", "<p><strong><em>foo</em></strong></p>", "***foo***\n"},
	{"1", "<p><strong>foo</strong></p>", "**foo**\n"},
	{"0", "<p>foo</p>", "foo\n"},
}

func TestHtml2Md(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range h2mTests {
		md, err := luteEngine.Html2Md(test.from)
		if nil != err {
			t.Fatalf("test case [%s] unexpected: %s", test.name, err)
		}

		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}

func TestHtml2MdSpec(t *testing.T) {
	// TODO: html to markdown spec

	//bytes, err := ioutil.ReadFile("commonmark-spec.json")
	//if nil != err {
	//	t.Fatalf("read spec test cases failed: " + err.Error())
	//}
	//
	//var testcases []testcase
	//if err = json.Unmarshal(bytes, &testcases); nil != err {
	//	t.Fatalf("read spec test caes failed: " + err.Error())
	//}
	//
	//luteEngine := lute.New()
	//luteEngine.GFMTaskListItem = false
	//luteEngine.GFMTable = false
	//luteEngine.GFMAutoLink = false
	//luteEngine.GFMStrikethrough = false
	//luteEngine.SoftBreak2HardBreak = false
	//luteEngine.CodeSyntaxHighlight = false
	//luteEngine.AutoSpace = false
	//luteEngine.FixTermTypo = false
	//luteEngine.Emoji = false
	//
	//for _, test := range testcases {
	//	testName := test.Section + " " + strconv.Itoa(test.Example)
	//	md, err := luteEngine.Html2Md(test.HTML)
	//	if nil != err {
	//		t.Fatalf("test case [%s] unexpected: %s", testName, err)
	//	}
	//
	//	if test.Markdown != md {
	//		t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", testName, test.Markdown, md, test.HTML)
	//	}
	//}
}
