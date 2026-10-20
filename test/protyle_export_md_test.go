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
	"bytes"
	"fmt"
	"html"
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
)

func TestProtyleExportCodeBlockFence(t *testing.T) {
	for _, tc := range []struct {
		name, content, language, fence string
	}{
		{"ordinary", "hello\n", "text", "```"},
		{"three backticks", "before\n```\nafter\n", "markdown", "````"},
		{"longest run", "```\n`````\n````\n", "text", "``````"},
		{"inline run", "const ticks = \"````\";\n", "js", "`````"},
		{"editor joiner", "\u200d```\n", "text", "````"},
		{"empty", "", "", "```"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			engine := lute.New()
			engine.SetProtyleWYSIWYG(true)
			dom := fmt.Sprintf(`<div data-type="NodeCodeBlock" class="code-block"><div class="protyle-action"><span class="protyle-action__language">%s</span></div><div contenteditable="true">%s</div></div>`,
				tc.language, html.EscapeString(tc.content))
			markdown := engine.BlockDOM2StdMd(dom)
			if !strings.HasPrefix(markdown, tc.fence+tc.language+"\n") || !strings.HasSuffix(strings.TrimSpace(markdown), tc.fence) {
				t.Fatalf("unsafe exported fence: %q", markdown)
			}
			// 前端复制会去掉防输入识别的零宽连接符，清理后仍应保留一个完整代码块。
			cleaned := strings.ReplaceAll(markdown, "\u200d```", "```")
			tree := parse.Parse("", []byte(cleaned), engine.ParseOptions)
			code := tree.Root.FirstChild
			if code == nil || code.Type != ast.NodeCodeBlock || code.Next != nil {
				t.Fatalf("export no longer parses as one code block: %q", cleaned)
			}
			expected := strings.ReplaceAll(tc.content, "\u200d```", "```")
			if got := string(code.ChildByType(ast.NodeCodeBlockCode).Tokens); got != expected {
				t.Fatalf("code content changed: got %q, want %q", got, expected)
			}
		})
	}
}

func TestProtyleExportCodeBlockFencePreservesSource(t *testing.T) {
	for _, source := range []string{
		"~~~~~~text\n````\n~~~~\n~~~~~~\n",
		"``````text\n```\n``````\n",
		"    before\n    ```\n    after\n",
		"> ````text\n> ```\n> ````\n",
		"- ````text\n  ```\n  ````\n",
	} {
		engine := lute.New()
		tree := parse.Parse("", []byte(source), engine.ParseOptions)
		before := render.NewFormatRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()
		exported := render.NewProtyleExportMdRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()
		after := render.NewFormatRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()
		if !bytes.Equal(before, after) {
			t.Fatalf("export mutated the source tree: %q", source)
		}
		originalHTML := engine.MarkdownStr("", source)
		if got := engine.MarkdownStr("", string(exported)); got != originalHTML {
			t.Fatalf("export changed code structure: source %q, export %q, html %q", source, exported, got)
		}
	}
}

var protyleExportMdTests = []parseTest{

	{"19", "[foo & bar](https://baz.com)", "[foo & bar](https://baz.com)\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
	{"18", "![foo](bar \"baz\")", "<img alt=\"foo\" src=\"bar\" title=\"baz\" />\n\n{: id=\"20060102150405-1a2b3c4\" updated=\"20060102150405\" type=\"doc\"}\n"},
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
	luteEngine.SetImgTag(true)

	ast.Testing = true
	for _, test := range protyleExportMdTests {
		tree := parse.Parse("", []byte(test.from), luteEngine.ParseOptions)
		parse.NestedInlines2FlattedSpans(tree, true)
		renderer := render.NewProtyleExportMdRenderer(tree, luteEngine.RenderOptions, luteEngine.ParseOptions)
		output := renderer.Render()
		kmd := util.BytesToStr(output)
		if test.to != kmd {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, kmd, test.from)
		}
	}
}
