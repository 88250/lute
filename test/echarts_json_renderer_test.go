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

var echartsJSONRendererTests = []parseTest{

	{"2", "~foo~\n", "[{\"name\":\"Document\",\"children\":[{\"name\":\"Paragraph\\np\",\"children\":[{\"name\":\"Strikethrough\\ndel\"}]}]}]"},
	{"1", "# foo\n*bar*\n", "[{\"name\":\"Document\",\"children\":[{\"name\":\"Heading\\nh1\",\"children\":[{\"name\":\"Text\\nfoo\"}]},{\"name\":\"Paragraph\\np\",\"children\":[{\"name\":\"Emphasis\\nem\",\"children\":[{\"name\":\"Text\\nbar\"}]}]}]}]"},
	{"0", "", "[{\"name\":\"Document\"}]"},
}

func TestEChartsJSONRenderer(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range echartsJSONRendererTests {
		html := luteEngine.RenderEChartsJSON(test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
