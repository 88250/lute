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
	"fmt"
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestFootnoteLabelInlineTokens(t *testing.T) {
	const ref = `<sup class="footnotes-ref" id="footnotes-ref-1"><a href="#footnotes-def-1">1</a></sup>`
	for _, protyle := range []bool{false, true} {
		for _, sup := range []bool{false, true} {
			for _, label := range []string{"a_b", "a*b", "a&b", "a<b", "a=b", "a~b", "a^b", "a$b", "a#b", "a!b", "a(b)", "a`b", "a&amp;b"} {
				t.Run(fmt.Sprintf("protyle=%v/sup=%v/label=%s", protyle, sup, label), func(t *testing.T) {
					engine := lute.New()
					engine.SetProtyleWYSIWYG(protyle)
					engine.SetSup(sup)
					for _, scenario := range []struct{ source, paragraph string }{
						{"Text[^%s] and *after*.", "<p>Text" + ref + " and <em>after</em>.</p>\n"},
						{"*before[^%s] after*.", "<p><em>before" + ref + " after</em>.</p>\n"},
						{"Text[caption_*][^%s].", "<p>Text" + ref + ".</p>\n"},
					} {
						markdown := fmt.Sprintf(scenario.source, label) + "\n\n[^" + label + "]: body\n"
						output := engine.MarkdownStr("", markdown)
						if !strings.HasPrefix(output, scenario.paragraph) {
							t.Fatalf("label altered surrounding content:\nsource: %q\nwant: %s\ngot: %s", markdown, scenario.paragraph, output)
						}
					}
				})
			}
		}
	}
}

var fnTests = []parseTest{

	{"2", "[foo][^label]\n[^label]: bar\n", "<p><sup class=\"footnotes-ref\" id=\"footnotes-ref-1\"><a href=\"#footnotes-def-1\">1</a></sup></p>\n<div class=\"footnotes-defs-div\"><hr class=\"footnotes-defs-hr\" />\n<ol class=\"footnotes-defs-ol\"><li id=\"footnotes-def-1\"><p>bar <a href=\"#footnotes-ref-1\" class=\"vditor-footnotes__goto-ref\">↩</a></p>\n</li>\n</ol></div>"},
	{"1", "foo[^label]\n[^label]:bar\n    * baz", "<p>foo<sup class=\"footnotes-ref\" id=\"footnotes-ref-1\"><a href=\"#footnotes-def-1\">1</a></sup></p>\n<div class=\"footnotes-defs-div\"><hr class=\"footnotes-defs-hr\" />\n<ol class=\"footnotes-defs-ol\"><li id=\"footnotes-def-1\"><p>bar</p>\n<ul>\n<li>baz <a href=\"#footnotes-ref-1\" class=\"vditor-footnotes__goto-ref\">↩</a></li>\n</ul>\n</li>\n</ol></div>"},
	{"0", "foo[^1]\n[^1]:bar\n    * baz", "<p>foo<sup class=\"footnotes-ref\" id=\"footnotes-ref-1\"><a href=\"#footnotes-def-1\">1</a></sup></p>\n<div class=\"footnotes-defs-div\"><hr class=\"footnotes-defs-hr\" />\n<ol class=\"footnotes-defs-ol\"><li id=\"footnotes-def-1\"><p>bar</p>\n<ul>\n<li>baz <a href=\"#footnotes-ref-1\" class=\"vditor-footnotes__goto-ref\">↩</a></li>\n</ul>\n</li>\n</ol></div>"},
}

func TestFootnotes(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetSup(true)

	for _, test := range fnTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
