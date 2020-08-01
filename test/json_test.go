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

	{"12", "$\nfoo\n$\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeInlineMath\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeInlineMathOpenMarker\",\"Val\":\"$\"},{\"ID\":\"\",\"Type\":\"NodeInlineMathContent\",\"Val\":\"\\nfoo\\n\"},{\"ID\":\"\",\"Type\":\"NodeInlineMathCloseMarker\",\"Val\":\"$\"}]}]}]}"},
	{"11", "$$\nfoo\n$$\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeMathBlock\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeMathBlockOpenMarker\",\"Val\":\"$$\"},{\"ID\":\"\",\"Type\":\"NodeMathBlockContent\",\"Val\":\"foo\"},{\"ID\":\"\",\"Type\":\"NodeMathBlockCloseMarker\",\"Val\":\"$$\"}]}]}"},
	{"10", "~foo~\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeStrikethrough\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeStrikethrough1OpenMarker\",\"Val\":\"~\"},{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"},{\"ID\":\"\",\"Type\":\"NodeStrikethrough1CloseMarker\",\"Val\":\"~\"}]}]}]}"},
	{"9", "> foo\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeBlockquote\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"}]}]}]}"},
	{"8", "* foo\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeList\",\"Val\":\"0\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeListItem\",\"Val\":\"*\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"}]}]}]}]}"},
	{"7", "foo\n---\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeHeading\",\"Val\":\"2\",\"HeadingSetext\":true,\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"}]}]}"},
	{"6", "# foo\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeHeading\",\"Val\":\"1\",\"HeadingSetext\":false,\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"}]}]}"},
	{"5", "```\nfoo\n```\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeCodeBlock\",\"IsFencedCodeBlock\":true,\"Children\":[{\"ID\":\"\",\"Type\":\"NodeCodeBlockFenceOpenMarker\",\"Val\":\"```\"},{\"ID\":\"\",\"Type\":\"NodeCodeBlockFenceInfoMarker\"},{\"ID\":\"\",\"Type\":\"NodeCodeBlockCode\",\"Val\":\"foo\\n\"},{\"ID\":\"\",\"Type\":\"NodeCodeBlockFenceCloseMarker\",\"Val\":\"```\"}]}]}"},
	{"4", "`foo`\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeCodeSpan\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeCodeSpanOpenMarker\",\"Val\":\"`\"},{\"ID\":\"\",\"Type\":\"NodeCodeSpanContent\",\"Val\":\"foo\"},{\"ID\":\"\",\"Type\":\"NodeCodeSpanCloseMarker\",\"Val\":\"`\"}]}]}]}"},
	{"3", "**foo** __bar__\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeStrong\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeStrongA6kOpenMarker\",\"Val\":\"**\"},{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"},{\"ID\":\"\",\"Type\":\"NodeStrongA6kCloseMarker\",\"Val\":\"**\"}]},{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\" \"},{\"ID\":\"\",\"Type\":\"NodeStrong\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeStrongU8eOpenMarker\",\"Val\":\"__\"},{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"bar\"},{\"ID\":\"\",\"Type\":\"NodeStrongU8eCloseMarker\",\"Val\":\"__\"}]}]}]}"},
	{"2", "*foo* _bar_\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeEmphasis\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeEmA6kOpenMarker\",\"Val\":\"*\"},{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"},{\"ID\":\"\",\"Type\":\"NodeEmA6kCloseMarker\",\"Val\":\"*\"}]},{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\" \"},{\"ID\":\"\",\"Type\":\"NodeEmphasis\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeEmU8eOpenMarker\",\"Val\":\"_\"},{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"bar\"},{\"ID\":\"\",\"Type\":\"NodeEmU8eCloseMarker\",\"Val\":\"_\"}]}]}]}"},
	{"1", "foo\n\nbar\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"}]},{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"bar\"}]}]}"},
	{"0", "\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\"}"},
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
