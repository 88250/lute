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
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
)

var protyleExportMdTests = []parseTest{

	{"17", "` ``foo`` `", "`` `foo` ``\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"16", "`` `foo` ``", "`` `foo` ``\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"15", "`foo`", "`foo`\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"14", "[foo](bar baz.txt)", "[foo](bar%20baz.txt)\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"13", "| `foo\\\\|bar` |\n| -- |", "|`foo\\\\|bar`|\n| -|\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"12", "| `foo\\&#124;bar` |\n| -- |", "|`foo\\|bar`|\n| -|\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"11", "| `foo\\|bar` |\n| -- |", "|`foo\\|bar`|\n| -|\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"10", "| $foo\\\\\\|bar$ |\n| -- |", "|$foo\\\\\\\\|bar$|\n| -|\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"9", "| $foo\\\\|bar$ |\n| -- |", "|$foo\\\\\\|bar$|\n| -|\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"8", "| $foo\\|bar$ |\n| -- |", "|$foo\\\\|bar$|\n| -|\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"7", "[~\\~foo\\~foo\\~~](bar)", "<sub>[\\~foo\\~foo\\~](bar)</sub>\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"6", "[^\\^foo\\^foo\\^^](bar)", "<sup>[\\^foo\\^foo\\^](bar)</sup>\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"5", "[==\\=foo\\=foo\\===](bar)", "==[\\=foo\\=foo\\=](bar)==\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"4", "[~~\\~foo\\~foo\\~~~](bar)", "~~[\\~foo\\~foo\\~](bar)~~\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"3", "[*\\*foo\\*foo\\**](bar)", "*[\\*foo\\*foo\\*](bar)*\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"2", "[**\\*foo\\*foo\\***](bar)", "**[\\*foo\\*foo\\*](bar)**\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"1", "[`foo`](bar \"baz\")", "[`foo`](bar \"baz\")\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"0", "foo**bar**{: style=\"color: red;\"}baz", "foo**bar**{: style=\"color: red;\"}baz\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
}

func TestProtyleExportMd(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.ParseOptions.Mark = true
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.SetKramdownIAL(true)
	luteEngine.SetSuperBlock(true)
	luteEngine.SetTag(true)
	luteEngine.SetSub(true)
	luteEngine.SetSup(true)
	luteEngine.SetGitConflict(true)
	luteEngine.SetIndentCodeBlock(false)
	luteEngine.SetEmojiSite("http://127.0.0.1:6806/stage/protyle/images/emoji")
	luteEngine.SetAutoSpace(true)
	luteEngine.SetParagraphBeginningSpace(true)
	luteEngine.SetFileAnnotationRef(true)
	luteEngine.SetInlineMathAllowDigitAfterOpenMarker(true)
	luteEngine.SetTextMark(true)
	luteEngine.SetImgPathAllowSpace(true)

	ast.Testing = true
	for _, test := range protyleExportMdTests {
		tree := parse.Parse("", []byte(test.from), luteEngine.ParseOptions)
		parse.NestedInlines2FlattedSpans(tree, true)
		renderer := render.NewProtyleExportMdRenderer(tree, luteEngine.RenderOptions)
		output := renderer.Render()
		kmd := util.BytesToStr(output)
		if test.to != kmd {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, kmd, test.from)
		}
	}
}
