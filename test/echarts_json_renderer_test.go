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

var echartsJSONRendererTests = []parseTest{

	{"3", ":smile:", "[{\"name\":\"Document\",\"children\":[{\"name\":\"Paragraph\\np\",\"children\":[{\"name\":\"Emoji Unicode\\n\"}]}]}]"},
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
