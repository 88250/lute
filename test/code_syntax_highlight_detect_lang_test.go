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

func TestCodeSyntaxHighlightDetectLang(t *testing.T) {
	// 围栏代码块自动探测语言 https://github.com/88250/lute/issues/22

	luteEngine := lute.New()
	luteEngine.SetCodeSyntaxHighlightDetectLang(true)
	for _, test := range issue22Tests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var issue22Tests = []parseTest{
	{"0", issue22Case0[0], issue22Case0[1]},
}

var issue22Case0 = []string{
	"```" + `
package lute

func test() {
	return
}
` + "```",
	"<pre><code class=\"language-go highlight-chroma\"><span class=\"highlight-line\"><span class=\"highlight-cl\"><span class=\"highlight-kn\">package</span> <span class=\"highlight-nx\">lute</span>\n</span></span><span class=\"highlight-line\"><span class=\"highlight-cl\">\n</span></span><span class=\"highlight-line\"><span class=\"highlight-cl\"><span class=\"highlight-kd\">func</span> <span class=\"highlight-nf\">test</span><span class=\"highlight-p\">()</span> <span class=\"highlight-p\">{</span>\n</span></span><span class=\"highlight-line\"><span class=\"highlight-cl\">\t<span class=\"highlight-k\">return</span>\n</span></span><span class=\"highlight-line\"><span class=\"highlight-cl\"><span class=\"highlight-p\">}</span>\n</span></span></code></pre>\n",
}
