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
	"github.com/b3log/lute"
	"testing"
)

var code = `go
package main
import "fmt"
func main() {
fmt.Println("Hello, World!")
}
`

var codeSyntaxHighlightGoTests = []parseTest{

	{"0", "```" + code + "```\n", "<pre><code class=\"language-go highlight-chroma\"><span class=\"highlight-kn\">package</span> <span class=\"highlight-nx\">main</span>\n\n<span class=\"highlight-kn\">import</span> <span class=\"highlight-s\">&#34;fmt&#34;</span>\n\n<span class=\"highlight-kd\">func</span> <span class=\"highlight-nf\">main</span><span class=\"highlight-p\">()</span> <span class=\"highlight-p\">{</span>\n\t<span class=\"highlight-nx\">fmt</span><span class=\"highlight-p\">.</span><span class=\"highlight-nf\">Println</span><span class=\"highlight-p\">(</span><span class=\"highlight-s\">&#34;Hello, World!&#34;</span><span class=\"highlight-p\">)</span>\n<span class=\"highlight-p\">}</span>\n</code></pre>\n"},
}

func TestCodeSyntaxHighlightGo(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range codeSyntaxHighlightGoTests {
		html, err := luteEngine.MarkdownStr(test.name, test.from)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

