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

	{"26", "[^foo] [^404]\n\n[^foo]: bar\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeFootnotesRef\",\"Val\":\"^foo\"},{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\" [^404]\"}]},{\"ID\":\"\",\"Type\":\"NodeFootnotesDef\",\"Val\":\"^foo\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"bar\"}]}]}]}"},
	{"25", "[toc]\n\n# foo\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeToC\"},{\"ID\":\"\",\"Type\":\"NodeHeading\",\"Val\":\"1\",\"HeadingSetext\":false,\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"}]}]}"},
	{"24", "foo\\*\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"},{\"ID\":\"\",\"Type\":\"NodeBackslash\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeBackslashContent\",\"Val\":\"*\"}]}]}]}"},
	{"23", "&hearts;\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeHTMLEntity\",\"Val\":\"&hearts;\"}]}]}"},
	{"22", ":octocat:\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeEmojiImg\",\"Val\":\":octocat:\"}]}]}"},
	{"21", ":heart:\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeEmojiUnicode\",\"Val\":\":heart:\"}]}]}"},
	{"20", "| foo | bar |\n| - | - |\n| baz | baz2 |\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeTable\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeTableHead\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeTableRow\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeTableCell\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"}]},{\"ID\":\"\",\"Type\":\"NodeTableCell\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"bar\"}]}]}]},{\"ID\":\"\",\"Type\":\"NodeTableRow\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeTableCell\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"baz\"}]},{\"ID\":\"\",\"Type\":\"NodeTableCell\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"baz2\"}]}]}]}]}"},
	{"19", "[foo](bar)\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeLink\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeOpenBracket\",\"Val\":\"[\"},{\"ID\":\"\",\"Type\":\"NodeLinkText\",\"Val\":\"foo\"},{\"ID\":\"\",\"Type\":\"NodeCloseBracket\",\"Val\":\"]\"},{\"ID\":\"\",\"Type\":\"NodeOpenParen\",\"Val\":\"(\"},{\"ID\":\"\",\"Type\":\"NodeLinkDest\",\"Val\":\"bar\"},{\"ID\":\"\",\"Type\":\"NodeCloseParen\",\"Val\":\")\"}]}]}]}"},
	{"18", "<span>foo</span>\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeInlineHTML\",\"Val\":\"<span>\"},{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"},{\"ID\":\"\",\"Type\":\"NodeInlineHTML\",\"Val\":\"</span>\"}]}]}"},
	{"17", "<div>foo</div>\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeHTMLBlock\",\"Val\":\"<div>foo</div>\"}]}"},
	{"16", "foo\n\n---\n\nbar\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"}]},{\"ID\":\"\",\"Type\":\"NodeThematicBreak\"},{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"bar\"}]}]}"},
	{"15", "foo\nbar\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"foo\"},{\"ID\":\"\",\"Type\":\"NodeSoftBreak\"},{\"ID\":\"\",\"Type\":\"NodeText\",\"Val\":\"bar\"}]}]}"},
	{"14", "---\nfoo\n---\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeYamlFrontMatter\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeYamlFrontMatterOpenMarker\",\"Val\":\"---\"},{\"ID\":\"\",\"Type\":\"NodeYamlFrontMatterContent\",\"Val\":\"foo\"},{\"ID\":\"\",\"Type\":\"NodeYamlFrontMatterCloseMarker\",\"Val\":\"---\"}]}]}"},
	{"13", "![foo](bar \"baz\")\n", "{\"ID\":\"\",\"Type\":\"NodeDocument\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeParagraph\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeImage\",\"Children\":[{\"ID\":\"\",\"Type\":\"NodeBang\",\"Val\":\"!\"},{\"ID\":\"\",\"Type\":\"NodeOpenBracket\",\"Val\":\"[\"},{\"ID\":\"\",\"Type\":\"NodeLinkText\",\"Val\":\"foo\"},{\"ID\":\"\",\"Type\":\"NodeCloseBracket\",\"Val\":\"]\"},{\"ID\":\"\",\"Type\":\"NodeOpenParen\",\"Val\":\"(\"},{\"ID\":\"\",\"Type\":\"NodeLinkDest\",\"Val\":\"bar\"},{\"ID\":\"\",\"Type\":\"NodeLinkSpace\",\"Val\":\" \"},{\"ID\":\"\",\"Type\":\"NodeLinkTitle\",\"Val\":\"baz\"},{\"ID\":\"\",\"Type\":\"NodeCloseParen\",\"Val\":\")\"}]}]}]}"},
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
	luteEngine.ToC = true

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
