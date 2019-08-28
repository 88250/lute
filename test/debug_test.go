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
	"fmt"
	"testing"

	"github.com/b3log/lute"
)

var debugTests = []parseTest{

	{"6", "[https://github.com/b3log/lute](https://github.com/b3log/lute)\n", "<p><a href=\"https://github.com/b3log/lute\">https://github.com/b3log/lute</a></p>\n"},
	{"5", "[1\n--\n", "<h2>[1</h2>\n"},
	{"4", "[1 \n", "<p>[1</p>\n"},
	{"3", "- -\r\n", "<ul>\n<li>\n<ul>\n<li></li>\n</ul>\n</li>\n</ul>\n"},
	{"2", "foo@bar.baz\n", "<p><a href=\"mailto:foo@bar.baz\">foo@bar.baz</a></p>\n"},
	{"1", "B3log https://b3log.org Lute\n", "<p>B3log <a href=\"https://b3log.org\">https://b3log.org</a> Lute</p>\n"},
	{"0", "[https://b3log.org](https://b3log.org)\n", "<p><a href=\"https://b3log.org\">https://b3log.org</a></p>\n"},
}

func TestDebug(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range debugTests {
		fmt.Println("Test [" + test.name + "]")
		html, err := luteEngine.MarkdownStr(test.name, test.markdown)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.html != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.html, html, test.markdown)
		}
	}
}
