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

	{"5", "```\nfoo\n```\n", "{\"Type\":\"NodeDocument\",\"Children\":[{\"Type\":\"NodeCodeBlock\",\"IsFencedCodeBlock\":true,\"Children\":[{\"Type\":\"NodeCodeBlockFenceOpenMarker\",\"Val\":\"```\"},{\"Type\":\"NodeCodeBlockFenceInfoMarker\"},{\"Type\":\"NodeCodeBlockCode\",\"Val\":\"foo\\n\"},{\"Type\":\"NodeCodeBlockFenceCloseMarker\",\"Val\":\"```\"}]}]}"},
	{"4", "`foo`\n", "{\"Type\":\"NodeDocument\",\"Children\":[{\"Type\":\"NodeParagraph\",\"Children\":[{\"Type\":\"NodeCodeSpan\",\"Children\":[{\"Type\":\"NodeCodeSpanOpenMarker\",\"Val\":\"`\"},{\"Type\":\"NodeCodeSpanContent\",\"Val\":\"foo\"},{\"Type\":\"NodeCodeSpanCloseMarker\",\"Val\":\"`\"}]}]}]}"},
	{"3", "**foo** __bar__\n", "{\"Type\":\"NodeDocument\",\"Children\":[{\"Type\":\"NodeParagraph\",\"Children\":[{\"Type\":\"NodeStrong\",\"Children\":[{\"Type\":\"NodeStrongA6kOpenMarker\",\"Val\":\"**\"},{\"Type\":\"NodeText\",\"Val\":\"foo\"},{\"Type\":\"NodeStrongA6kCloseMarker\",\"Val\":\"**\"}]},{\"Type\":\"NodeText\",\"Val\":\" \"},{\"Type\":\"NodeStrong\",\"Children\":[{\"Type\":\"NodeStrongU8eOpenMarker\",\"Val\":\"__\"},{\"Type\":\"NodeText\",\"Val\":\"bar\"},{\"Type\":\"NodeStrongU8eCloseMarker\",\"Val\":\"__\"}]}]}]}"},
	{"2", "*foo* _bar_\n", "{\"Type\":\"NodeDocument\",\"Children\":[{\"Type\":\"NodeParagraph\",\"Children\":[{\"Type\":\"NodeEmphasis\",\"Children\":[{\"Type\":\"NodeEmA6kOpenMarker\",\"Val\":\"*\"},{\"Type\":\"NodeText\",\"Val\":\"foo\"},{\"Type\":\"NodeEmA6kCloseMarker\",\"Val\":\"*\"}]},{\"Type\":\"NodeText\",\"Val\":\" \"},{\"Type\":\"NodeEmphasis\",\"Children\":[{\"Type\":\"NodeEmU8eOpenMarker\",\"Val\":\"_\"},{\"Type\":\"NodeText\",\"Val\":\"bar\"},{\"Type\":\"NodeEmU8eCloseMarker\",\"Val\":\"_\"}]}]}]}"},
	{"1", "foo\n\nbar\n", "{\"Type\":\"NodeDocument\",\"Children\":[{\"Type\":\"NodeParagraph\",\"Children\":[{\"Type\":\"NodeText\",\"Val\":\"foo\"}]},{\"Type\":\"NodeParagraph\",\"Children\":[{\"Type\":\"NodeText\",\"Val\":\"bar\"}]}]}"},
	{"0", "\n", "{\"Type\":\"NodeDocument\"}"},
}

func TestJSON(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range jsonTests {
		json := luteEngine.RenderJSON(test.from)
		if test.to != json {
			t.Fatalf("test case [%s] failed\nexpected\n\t%s\ngot\n\t%s\noriginal markdown text\n\t%q", test.name, test.to, json, test.from)
		}

		tree := luteEngine.ParseJSON(json)
		renderer := render.NewFormatRenderer(tree)
		markdown := util.BytesToStr(renderer.Render())
		if test.from != markdown {
			t.Fatalf("test case [%s] failed\nexpected\n\t%s\ngot\n\t%s\noriginal markdown text\n\t%q", test.name, test.from, markdown, test.from)
		}
	}
}
