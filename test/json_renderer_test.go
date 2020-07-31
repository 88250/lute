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

var jsonRendererTests = []parseTest{

	{"1", "foo\n\nbar\n", "[{\"name\":\"Document\",\"children\":[{\"name\":\"Paragraph\\np\",\"children\":[{\"name\":\"Text\\nfoo\"}]}]}]"},
	{"0", "", "[{\"name\":\"Document\"}]"},
}

func TestJSONRenderer(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range jsonRendererTests {
		html := luteEngine.RenderJSON(test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%s\ngot\n\t%s\noriginal markdown text\n\t%s", test.name, test.to, html, test.from)
		}
	}
}
