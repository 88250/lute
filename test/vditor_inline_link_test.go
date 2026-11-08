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
	"github.com/88250/lute/editor"
	"github.com/88250/lute/parse"
)

func TestVditorInlineLinkRoundTrip(t *testing.T) {
	for _, mode := range []string{"wysiwyg", "ir"} {
		for _, fixture := range []struct {
			name, markdown string
		}{
			{"empty-angle", `[](<url> "title")`},
			{"empty-space", `[](<url 1> "title")`},
			{"empty-bare", `[](url "title")`},
			{"empty-untitled", `[](<url 1>)`},
			{"text-space", `[text](<url 1> "title")`},
			{"encoded-space", `[text](url%201 "title")`},
			{"mixed-spaces", `[text](<url%201 2> "title")`},
			{"surrounding-text", `before [](<url 1> "title") after`},
			{"image-link", `[![alt](image.png)](<url 1> "title")`},
		} {
			t.Run(mode+"/"+fixture.name, func(t *testing.T) {
				standard := lute.New()
				want := standard.Md2HTML(fixture.markdown)
				engine := lute.New()
				toDOM, toMarkdown, toHTML, spin := engine.Md2VditorDOM, engine.VditorDOM2Md, engine.VditorDOM2HTML, engine.SpinVditorDOM
				if mode == "wysiwyg" {
					engine.SetVditorWYSIWYG(true)
				} else {
					engine.SetVditorIR(true)
					toDOM, toMarkdown, toHTML, spin = engine.Md2VditorIRDOM, engine.VditorIRDOM2Md, engine.VditorIRDOM2HTML, engine.SpinVditorIRDOM
				}
				dom := toDOM(fixture.markdown)
				for round := 0; round < 3; round++ {
					markdown := toMarkdown(dom)
					if got := standard.Md2HTML(markdown); got != want {
						t.Fatalf("round %d changed link semantics\nwant: %q\ngot: %q\nMarkdown: %q", round, want, got, markdown)
					}
					if got := toHTML(dom); got != want {
						t.Fatalf("round %d changed exported HTML\nwant: %q\ngot: %q", round, want, got)
					}
					dom = spin(dom)
				}
			})
		}
	}
}

func TestVditorEmptyLinkCaret(t *testing.T) {
	engine := lute.New()
	engine.SetVditorWYSIWYG(true)
	dom := engine.Md2VditorDOM(`[](<url ` + editor.Caret + `1> "title")`)
	if !strings.Contains(dom, editor.Caret) || !strings.Contains(dom, `href="url%201"`) {
		t.Fatalf("empty link lost its caret or destination: %s", dom)
	}
	dom = strings.ReplaceAll(dom, editor.Caret, editor.FrontEndCaret)
	if got := engine.VditorDOM2Md(engine.SpinVditorDOM(dom)); got != "[](url%201 \"title\")\n" {
		t.Fatalf("empty link changed after caret removal: %q", got)
	}
}

func TestVditorEmptyAnchorAndEditedDestination(t *testing.T) {
	engine := lute.New()
	engine.SetVditorWYSIWYG(true)
	for _, label := range []string{"", "text"} {
		dom := `<p><a href="url 1" title="title">` + label + `</a></p>`
		want := "[" + label + "](url%201 \"title\")\n"
		if got := engine.VditorDOM2Md(dom); got != want {
			t.Fatalf("edited anchor changed\nwant: %q\ngot: %q", want, got)
		}
	}
	if got := engine.VditorDOM2Md(`<p><a></a><a href=""></a></p>`); got != "\n" {
		t.Fatalf("destination-less empty anchors produced Markdown: %q", got)
	}
}

func TestVditorIncompleteInlineLink(t *testing.T) {
	for _, mode := range []string{"wysiwyg", "ir"} {
		engine := lute.New()
		engine.SetVditorWYSIWYG(mode == "wysiwyg")
		engine.SetVditorIR(mode == "ir")
		for _, markdown := range []string{"[]()", "[](<>)", "[](", `[](<url> "title"`} {
			tree := parse.Parse("", []byte(markdown), engine.ParseOptions)
			ast.Walk(tree.Root, func(node *ast.Node, entering bool) ast.WalkStatus {
				if entering && node.Type == ast.NodeLink {
					t.Fatalf("%s parsed an incomplete link: %q", mode, markdown)
				}
				return ast.WalkContinue
			})
		}
	}
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	tree := parse.Parse("", []byte(`[](<url> "title")`), engine.ParseOptions)
	ast.Walk(tree.Root, func(node *ast.Node, entering bool) ast.WalkStatus {
		if entering && node.Type == ast.NodeLink {
			t.Fatal("Protyle parsed an empty link label")
		}
		return ast.WalkContinue
	})
}
