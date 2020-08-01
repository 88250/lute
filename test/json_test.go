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

	{"7", "foo\n---\n", "{\"Id\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeHeading\",\"Val\":\"2\",\"HeadingSetext\":true,\"Children\":[{\"Id\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"}]}]}"},
	{"6", "# foo\n", "{\"Id\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeHeading\",\"Val\":\"1\",\"HeadingSetext\":false,\"Children\":[{\"Id\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"}]}]}"},
	{"5", "```\nfoo\n```\n", "{\"Id\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeCodeBlock\",\"IsFencedCodeBlock\":true,\"Children\":[{\"Id\":\"\",\"Type\":\"NodeCodeBlockFenceOpenMarker\",\"Val\":\"```\"},{\"Id\":\"\",\"Type\":\"NodeCodeBlockFenceInfoMarker\"},{\"Id\":\"\",\"Type\":\"NodeCodeBlockCode\",\"Val\":\"foo\\n\"},{\"Id\":\"\",\"Type\":\"NodeCodeBlockFenceCloseMarker\",\"Val\":\"```\"}]}]}"},
	{"4", "`foo`\n", "{\"Id\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeCodeSpan\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeCodeSpanOpenMarker\",\"Val\":\"`\"},{\"Id\":\"\",\"Type\":\"NodeCodeSpanContent\",\"Val\":\"foo\"},{\"Id\":\"\",\"Type\":\"NodeCodeSpanCloseMarker\",\"Val\":\"`\"}]}]}]}"},
	{"3", "**foo** __bar__\n", "{\"Id\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeStrong\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeStrongA6kOpenMarker\",\"Val\":\"**\"},{\"Id\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"},{\"Id\":\"\",\"Type\":\"NodeStrongA6kCloseMarker\",\"Val\":\"**\"}]},{\"Id\":\"\",\"Type\":\"NodeText\",\"Val\":\" \"},{\"Id\":\"\",\"Type\":\"NodeStrong\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeStrongU8eOpenMarker\",\"Val\":\"__\"},{\"Id\":\"\",\"Type\":\"NodeText\",\"Val\":\"bar\"},{\"Id\":\"\",\"Type\":\"NodeStrongU8eCloseMarker\",\"Val\":\"__\"}]}]}]}"},
	{"2", "*foo* _bar_\n", "{\"Id\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeEmphasis\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeEmA6kOpenMarker\",\"Val\":\"*\"},{\"Id\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"},{\"Id\":\"\",\"Type\":\"NodeEmA6kCloseMarker\",\"Val\":\"*\"}]},{\"Id\":\"\",\"Type\":\"NodeText\",\"Val\":\" \"},{\"Id\":\"\",\"Type\":\"NodeEmphasis\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeEmU8eOpenMarker\",\"Val\":\"_\"},{\"Id\":\"\",\"Type\":\"NodeText\",\"Val\":\"bar\"},{\"Id\":\"\",\"Type\":\"NodeEmU8eCloseMarker\",\"Val\":\"_\"}]}]}]}"},
	{"1", "foo\n\nbar\n", "{\"Id\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"}]},{\"Id\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"Id\":\"\",\"Type\":\"NodeText\",\"Val\":\"bar\"}]}]}"},
	{"0", "\n", "{\"Id\":\"\",\"Type\":\"NodeDocument\"}"},
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
