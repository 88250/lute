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

var codeSyntaxHighlightTests = []parseTest{

	{"0", "```java\nint i;\n```\n", "<pre><code class=\"language-java\"><span class=\"highlight-kt\">int</span> <span class=\"highlight-nf\">i</span><span class=\"highlight-p\">;</span>\n</code></pre>\n"},
}

func TestCodeSyntaxHighlight(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range codeSyntaxHighlightTests {
		html, err := luteEngine.MarkdownStr(test.name, test.markdown)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.html != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.html, html, test.markdown)
		}
	}
}

var codeSyntaxHighlightInlineTests = []parseTest{

	{"0", "```java\nint i;\n```\n", "<pre><code class=\"language-java\"><span style=\"color:#8be9fd\">int</span> <span style=\"color:#50fa7b\">i</span>;\n</code></pre>\n"},
}

func TestCodeSyntaxHighlightInline(t *testing.T) {
	luteEngine := lute.New(lute.CodeSyntaxHighlight(true, true, "dracula"))

	for _, test := range codeSyntaxHighlightInlineTests {
		html, err := luteEngine.MarkdownStr(test.name, test.markdown)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.html != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.html, html, test.markdown)
		}
	}
}
