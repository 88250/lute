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

var vditorRendererTests = []parseTest{

	{"10", "Lu\nte\n", "<p data-id=\"0\" data-type=\"1\"><span data-id=\"0\" data-type=\"10\">Lu</span><span data-id=\"0\" data-type=\"15\" /></span>\n<span data-id=\"0\" data-type=\"10\">te</span></p>\n"},
	{"9", "Lu  \nte\n", "<p data-id=\"0\" data-type=\"1\"><span data-id=\"0\" data-type=\"10\">Lu</span><span data-id=\"0\" data-type=\"14\"></span>\n<span data-id=\"0\" data-type=\"10\">te</span></p>\n"},
	{"8", "Lu\\\nte\n", "<p data-id=\"0\" data-type=\"1\"><span data-id=\"0\" data-type=\"10\">Lu</span><span data-id=\"0\" data-type=\"14\"></span>\n<span data-id=\"0\" data-type=\"10\">te</span></p>\n"},
	{"7", "`Lute`\n", "<p data-id=\"0\" data-type=\"1\"><span><span class=\"open\">`</span><code data-id=\"0\" data-type=\"13\">Lute</code><span class=\"close\">`</span></p>\n"},
	{"6", "**Lute**\n", "<p data-id=\"0\" data-type=\"1\"><span><span class=\"open\">**</span><strong data-id=\"0\" data-type=\"12\"><span data-id=\"0\" data-type=\"10\">Lute</span></strong><span class=\"close\">**</span></p>\n"},
	{"5", "*Lute*\n", "<p data-id=\"0\" data-type=\"1\"><span><span class=\"open\">*</span><em data-id=\"0\" data-type=\"11\"><span data-id=\"0\" data-type=\"10\">Lute</span></em><span class=\"close\">*</span></p>\n"},
	{"4", "* Lute\n", "<ul data-id=\"0\" data-type=\"5\">\n<li data-id=\"0\" data-type=\"6\"><span data-id=\"0\" data-type=\"10\">Lute</span></li>\n</ul>\n"},
	{"3", "> Lute\n", "<blockquote data-id=\"0\" data-type=\"4\">\n<p data-id=\"0\" data-type=\"1\"><span data-id=\"0\" data-type=\"10\">Lute</span></p>\n</blockquote>\n"},
	{"2", "---\n", "<hr data-id=\"0\" data-type=\"3\" />\n"},
	{"1", "## Lute\n", "<h2 data-id=\"0\" data-type=\"2\"><span data-id=\"0\" data-type=\"10\">Lute</span></h2>\n"},
	{"0", "Lute\n", "<p data-id=\"0\" data-type=\"1\"><span data-id=\"0\" data-type=\"10\">Lute</span></p>\n"},
}

func TestVditorRenderer(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorRendererTests {
		html, err := luteEngine.RenderVditorDOM(1, test.markdown)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.html != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.html, html, test.markdown)
		}
	}
}
