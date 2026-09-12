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
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/parse"
)

var fileAnnotationRefTests = []parseTest{

	{"8", "<<assets/foo bar-20211115212742-80mbhnk.pdf/20211115213157-zifdhvu \"The phrase \\\"regular expression\\\" is often \">>", "<p>\"The phrase &quot;regular expression&quot; is often \"</p>\n"},
	{"7", "<<assets/foo bar.pdf/20210911230820-lhiaysx \"注解<锚文本>\">>", "<p>\"注解&lt;锚文本&gt;\"</p>\n"},
	{"6", "<<assets/foo bar.pdf/20210911230820-lhiaysx \"注解锚文本\">>", "<p>\"注解锚文本\"</p>\n"},
	{"5", "<<foo bar-20210911230735-pzlpdt.txt/20210911230820-lhiaysx \"注解锚文本\">>", "<p>&lt;&lt;foo bar-20210911230735-pzlpdt.txt/20210911230820-lhiaysx &quot;注解锚文本&quot;&gt;&gt;</p>\n"},
	{"4", "<<assets/foo bar-20210911230735-pzlpdt.txt/20210911230820-lhiaysx \"注解锚文本\">>", "<p>&lt;&lt;assets/foo bar-20210911230735-pzlpdt.txt/20210911230820-lhiaysx &quot;注解锚文本&quot;&gt;&gt;</p>\n"},
	{"3", "<<assets/foo bar-20210911230735-pzlpdt.pdf/20210911230820-lhiaysx \"注解锚文本\">>", "<p>\"注解锚文本\"</p>\n"},
	{"2", "foo<<<bar>>>bazbazbazbazbazbazbazbazbazbazbazbazbazbazbazbaz", "<p>foo&lt;&lt;<bar>&gt;&gt;bazbazbazbazbazbazbazbazbazbazbazbazbazbazbazbaz</p>\n"},
	{"1", "<<assets/foo bar-20210911230735-pzlpdtf.pdf/20210911230820-lhiaysx \"注解锚文本\">>", "<p>\"注解锚文本\"</p>\n"},
	{"0", "<<assets/文件名-20210911230735-pzlpdtf.pdf/20210911230820-lhiaysx \"注解锚文本\">>", "<p>\"注解锚文本\"</p>\n"},
}

func TestFileAnnotationRefPaths(t *testing.T) {
	l := lute.New()
	l.SetFileAnnotationRef(true)
	l.SetTextMark(true)
	l.SetProtyleWYSIWYG(true)
	for _, reference := range []string{
		"assets/a.pdf/20210911230820-lhiaysx",
		"assets/document.pdf/20210911230820-lhiaysx",
		"assets/文档 空格.PDF/20210911230820-lhiaysx",
		"assets/folder/a%20b.pdf/20210911230820-lhiaysx",
		"assets/document-20210911230735-pzlpdtf.pdf/20210911230820-lhiaysx",
		"assets/document.pdf/20210911230820-lhiaysx?box=20210911230735-pzlpdtf&dataPath=/docs/a.sy",
		"assets/document.pdf/20210911230820-lhiaysx?dataPath=%2Fdocs%2Fa%20b.sy#view",
		"assets/a&copy;.pdf/20210911230820-lhiaysx?dataPath=/a.sy&copy;=yes",
	} {
		t.Run(reference, func(t *testing.T) {
			for _, anchor := range []string{"", ` ""`, ` "a"`, ` "<anchor>"`, ` "a \"quote\""`} {
				markdown := "<<" + reference + anchor + ">>"
				tree := parse.Parse("", []byte(markdown), l.ParseOptions)
				node := tree.Root.FirstChild.FirstChild
				if node.Type != ast.NodeFileAnnotationRef || node.ChildByType(ast.NodeFileAnnotationRefID).TokensStr() != reference {
					t.Fatalf("reference was not parsed intact: %q", markdown)
				}
				if anchor == ` "a"` {
					dom := l.Md2BlockDOM(markdown, false)
					if !strings.Contains(dom, `data-id="`+html.EscapeHTMLStr(reference)+`"`) {
						t.Fatalf("reference missing from DOM: %s", dom)
					}
					if got := strings.TrimSpace(l.BlockDOM2StdMd(dom)); got != markdown {
						t.Fatalf("round trip: got %q, want %q", got, markdown)
					}
				}
			}
		})
	}
}

func TestFileAnnotationRefInvalid(t *testing.T) {
	l := lute.New()
	l.SetFileAnnotationRef(true)
	for _, reference := range []string{
		"a.pdf/20210911230820-lhiaysx",
		"https://example.com/a.pdf/20210911230820-lhiaysx",
		"assets/a.txt/20210911230820-lhiaysx",
		"assets/a.pdf/20210911230820-lhiays",
		"assets/a.pdf/20210911230820-lhiays_",
		"assets/a.pdf/20210911230820-LHIAYSX",
		"assets/a.pdf/20210911230820-lhiaysx/extra",
		"assets/../a.pdf/20210911230820-lhiaysx",
		"assets//a.pdf/20210911230820-lhiaysx",
		"assets/<a>.pdf/20210911230820-lhiaysx",
	} {
		markdown := "<<" + reference + ` "anchor">>`
		tree := parse.Parse("", []byte(markdown), l.ParseOptions)
		ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
			if entering && n.Type == ast.NodeFileAnnotationRef {
				t.Errorf("invalid reference parsed: %q", markdown)
			}
			return ast.WalkContinue
		})
	}
	for _, markdown := range []string{"<<", `<<assets/a.pdf/20210911230820-lhiaysx "unterminated`, `<<assets/a.pdf/20210911230820-lhiaysx "a">`} {
		if got := l.MarkdownStr("", markdown); !strings.Contains(got, "&lt;&lt;") {
			t.Errorf("incomplete reference lost its prefix: %q => %q", markdown, got)
		}
	}
}

func TestFileAnnotationRef(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetFileAnnotationRef(true)
	for _, test := range fileAnnotationRefTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
