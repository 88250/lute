package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func TestInlineMemoEscapedContentRoundTrip(t *testing.T) {
	engine := lute.New()
	engine.SetTextMark(true)
	engine.SetHTMLTag2TextMark(true)
	engine.SetKramdownIAL(true)
	engine.SetSpin(true)
	engine.SetProtyleWYSIWYG(true)
	for _, memo := range []string{
		"&lt;a", "before&lt;a after", "vector&lt;T&gt;", "&amp;lt;a", "&amp;#60;a",
		"&lt;img src=x onerror=alert(1)&gt;", "&quot;quote&quot; &amp; text",
		`<b>legacy</b><br><a href="https://example.com/?a=1&amp;b=2">link</a>`,
	} {
		t.Run(memo, func(t *testing.T) {
			dom := `<div data-node-id="20261018120000-abcdefg" data-type="NodeParagraph"><div contenteditable="true"><span data-type="inline-memo" data-inline-memo-content="` + html.EscapeHTMLStr(memo) + `">note</span></div></div>`
			for round := 0; round < 3; round++ {
				tree := engine.BlockDOM2Tree(dom)
				assertInlineMemoContent(t, tree, memo)
				for name, output := range map[string]string{
					"HTML":     engine.BlockDOM2HTML(dom),
					"Markdown": engine.BlockDOM2StdMd(dom),
					"preview":  engine.ProtylePreview(tree, engine.RenderOptions, engine.ParseOptions),
					"DOCX":     string(render.NewProtyleExportDocxRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()),
				} {
					if !strings.Contains(output, memo) {
						t.Fatalf("%s export lost memo on round %d: %q", name, round, output)
					}
				}
				dom = engine.SpinBlockDOM(dom)
				assertInlineMemoContent(t, engine.BlockDOM2Tree(dom), memo)
				markdown := engine.BlockDOM2Md(dom)
				dom = engine.Md2BlockDOM(markdown, false)
				assertInlineMemoContent(t, engine.BlockDOM2Tree(dom), memo)
			}
		})
	}
}

func TestInlineMemoLegacyTreeRoundTrip(t *testing.T) {
	engine := lute.New()
	engine.SetTextMark(true)
	engine.SetHTMLTag2TextMark(true)
	engine.SetProtyleWYSIWYG(true)
	for _, memo := range []string{`<b>legacy</b>`, `foo&quot;bar&quot;baz`, `plain & safe`} {
		root := &ast.Node{Type: ast.NodeDocument, ID: "20261018120000-abcdefg"}
		paragraph := &ast.Node{Type: ast.NodeParagraph, ID: "20261018120001-abcdefg"}
		paragraph.AppendChild(&ast.Node{Type: ast.NodeTextMark, TextMarkType: "inline-memo", TextMarkTextContent: "note", TextMarkInlineMemoContent: memo})
		root.AppendChild(paragraph)
		tree := &parse.Tree{Root: root}
		dom := engine.Tree2BlockDOM(tree, engine.RenderOptions, engine.ParseOptions)
		assertInlineMemoContent(t, engine.BlockDOM2Tree(dom), memo)
	}
}

func TestInlineMemoExportWhitespace(t *testing.T) {
	engine := lute.New()
	engine.SetTextMark(true)
	engine.SetProtyleWYSIWYG(true)
	dom := `<div data-node-id="20261018120000-abcdefg" data-type="NodeParagraph"><div contenteditable="true"><span data-type="inline-memo" data-inline-memo-content="first` + "\n" + `second">note</span></div></div>`
	if output := engine.BlockDOM2HTML(dom); !strings.Contains(output, "<sup>(first second)</sup>") {
		t.Fatalf("memo newline was not normalized: %q", output)
	}
}

func assertInlineMemoContent(t *testing.T, tree *parse.Tree, want string) {
	t.Helper()
	found := false
	ast.Walk(tree.Root, func(node *ast.Node, entering bool) ast.WalkStatus {
		if entering && node.IsTextMarkType("inline-memo") {
			found = true
			if node.TextMarkInlineMemoContent != want {
				t.Fatalf("memo = %q, want %q", node.TextMarkInlineMemoContent, want)
			}
		}
		return ast.WalkContinue
	})
	if !found {
		t.Fatal("inline memo was lost")
	}
}
