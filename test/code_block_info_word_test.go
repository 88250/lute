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
	"testing"

	"github.com/88250/lute"
)

// The language of a fenced code block is the first word of the info string
// (CommonMark 4.5). Words are separated by any whitespace, but only a space
// ended the word, so a tab in the info string ended up in the class attribute.
var codeBlockInfoWordTests = []parseTest{
	{"tab", "```js\tfoo\ncode\n```\n", "<pre><code class=\"language-js\">code\n</code></pre>\n"},
	{"space", "```js foo\ncode\n```\n", "<pre><code class=\"language-js\">code\n</code></pre>\n"},
	{"leading-tab", "```\tjs\ncode\n```\n", "<pre><code class=\"language-js\">code\n</code></pre>\n"},
	{"tilde-tabs", "~~~\truby\tstartline=3\ncode\n~~~\n", "<pre><code class=\"language-ruby\">code\n</code></pre>\n"},
}

func TestCodeBlockInfoWord(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.RenderOptions.CodeSyntaxHighlight = false

	for _, test := range codeBlockInfoWordTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
