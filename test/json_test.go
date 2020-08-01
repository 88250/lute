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
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
)

var jsonTests = []parseTest{

	{"4", "`foo`\n", "[{\"type\":\"NodeParagraph\",\"children\":[{\"type\":\"NodeCodeSpan\",\"children\":[{\"type\":\"NodeCodeSpanOpenMarker\",\"val\":\"`\"},{\"type\":\"NodeCodeSpanContent\",\"val\":\"foo\"},{\"type\":\"NodeCodeSpanCloseMarker\",\"val\":\"`\"}]}]}]"},
	{"3", "**foo** __bar__\n", "[{\"type\":\"NodeParagraph\",\"children\":[{\"type\":\"NodeStrong\",\"children\":[{\"type\":\"NodeStrongA6kOpenMarker\",\"val\":\"**\"},{\"type\":\"NodeText\",\"val\":\"foo\"},{\"type\":\"NodeStrongA6kCloseMarker\",\"val\":\"**\"}]},{\"type\":\"NodeText\",\"val\":\" \"},{\"type\":\"NodeStrong\",\"children\":[{\"type\":\"NodeStrongU8eOpenMarker\",\"val\":\"__\"},{\"type\":\"NodeText\",\"val\":\"bar\"},{\"type\":\"NodeStrongU8eCloseMarker\",\"val\":\"__\"}]}]}]"},
	{"2", "*foo* _bar_\n", "[{\"type\":\"NodeParagraph\",\"children\":[{\"type\":\"NodeEmphasis\",\"children\":[{\"type\":\"NodeEmA6kOpenMarker\",\"val\":\"*\"},{\"type\":\"NodeText\",\"val\":\"foo\"},{\"type\":\"NodeEmA6kCloseMarker\",\"val\":\"*\"}]},{\"type\":\"NodeText\",\"val\":\" \"},{\"type\":\"NodeEmphasis\",\"children\":[{\"type\":\"NodeEmU8eOpenMarker\",\"val\":\"_\"},{\"type\":\"NodeText\",\"val\":\"bar\"},{\"type\":\"NodeEmU8eCloseMarker\",\"val\":\"_\"}]}]}]"},
	{"1", "foo\n\nbar\n", "[{\"type\":\"NodeParagraph\",\"children\":[{\"type\":\"NodeText\",\"val\":\"foo\"}]},{\"type\":\"NodeParagraph\",\"children\":[{\"type\":\"NodeText\",\"val\":\"bar\"}]}]"},
	{"0", "\n", "[]"},
}

func TestJSON(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range jsonTests {
		json := luteEngine.RenderJSON(test.from)
		if test.to != json {
			t.Fatalf("test case [%s] failed\nexpected\n\t%s\ngot\n\t%s\noriginal markdown text\n\t%s", test.name, test.to, json, test.from)
		}

		tree := luteEngine.ParseJSON(json)
		renderer := render.NewFormatRenderer(tree)
		markdown := util.BytesToStr(renderer.Render())
		if test.from != markdown {
			t.Fatalf("test case [%s] failed\nexpected\n\t%s\ngot\n\t%s\noriginal markdown text\n\t%s", test.name, test.from, markdown, test.from)
		}
	}
}
