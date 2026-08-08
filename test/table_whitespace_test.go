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
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestVditorTableInlineWhitespaceRoundTrip(t *testing.T) {
	luteEngine := lute.New()
	markdown := "| a *em* b **strong** c `code` d [link](https://example.com) e ~~del~~ f |\n| --- |\n| x |\n"
	expected := "| a *em* b **strong** c `code` d [link](https://example.com) e ~~del~~ f |"
	tests := []struct {
		name      string
		roundTrip func(string) string
	}{
		{"wysiwyg", func(markdown string) string {
			return luteEngine.VditorDOM2Md(luteEngine.Md2VditorDOM(markdown))
		}},
		{"ir", func(markdown string) string {
			return luteEngine.VditorIRDOM2Md(luteEngine.Md2VditorIRDOM(markdown))
		}},
	}
	for _, test := range tests {
		actual := strings.SplitN(test.roundTrip(markdown), "\n", 2)[0]
		if expected != actual {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q", test.name, expected, actual)
		}
	}
}

func TestHTML2MdTableInlineWhitespace(t *testing.T) {
	luteEngine := lute.New()
	html := "<table>\n<thead><tr><th>\n a <em>em</em> b <strong>strong</strong> c <code>code</code> d " +
		"<a href=\"https://example.com\">link</a> e <del>del</del> f \n</th></tr></thead>\n" +
		"<tbody><tr><td>x</td></tr></tbody>\n</table>"
	expected := "| a *em* b **strong** c `code` d [link](https://example.com) e ~~del~~ f |"
	actual := strings.SplitN(luteEngine.HTML2Md(html), "\n", 2)[0]
	if expected != actual {
		t.Fatalf("expected\n\t%q\ngot\n\t%q", expected, actual)
	}
}
