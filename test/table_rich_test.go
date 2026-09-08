package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func newTableRichLute() *lute.Lute {
	ret := lute.New()
	ret.ParseOptions = parse.TableCellRichOptions()
	ret.ParseOptions.DisableTableCellRich = false
	ret.RenderOptions.ProtyleWYSIWYG = true
	ret.SetKramdownIAL(true)
	ret.SetTextMark(true)
	return ret
}

func firstRichCell(tree *parse.Tree) (ret *ast.Node) {
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering && n.TableCellRich != nil {
			ret = n
			return ast.WalkStop
		}
		return ast.WalkContinue
	})
	return
}

func TestTableRichRoundTrip(t *testing.T) {
	engine := newTableRichLute()
	sources := []string{
		"- first\n- **second**",
		"# title\n\n> quote\n\n```go\nfmt.Println(\"a|b\")\n```\n\n$$\nx^2\n$$",
		"- [ ] task\n  - nested\n\n![image](assets/picture.png)\n\n((20260101000000-abcdefg 'reference'))",
		"literal \\*marker\\* and `a|b`\n\nnext paragraph",
		`<span data-type="text">font</span>{: style="font-family: var(--b3-font-family-emoji-reset), 'A&#92;&#92;B &gt; C', var(--b3-font-family-editor), var(--b3-font-family);"}`,
		`<span data-type="code">&lt;span data-type="text"&gt;font&lt;/span&gt;{: style="color: var(--b3-font-color1);"}</span>`,
		"",
	}
	for _, source := range sources {
		t.Run(source, func(t *testing.T) {
			tree := parse.Parse("", []byte("| title |\n| --- |\n| old |"), engine.ParseOptions)
			cell := tree.Root.FirstChild.LastChild.FirstChild
			cell.TableCellRich = &ast.TableCellRich{Spec: 1, Format: "kramdown", Content: source}
			cell.SetIALAttr("style", "color: red;")
			if err := parse.ApplyTableCellRichProjection(cell); err != nil {
				t.Fatal(err)
			}
			dom := engine.Tree2BlockDOM(tree, engine.RenderOptions, engine.ParseOptions)
			for iteration := 0; iteration < 3; iteration++ {
				dom = engine.SpinBlockDOM(dom)
				read := firstRichCell(engine.BlockDOM2Tree(dom))
				if read == nil || read.TableCellRich.Content != source || read.IALAttr("style") != "color: red;" {
					t.Fatalf("round trip lost source or cell style: %s", dom)
				}
				ast.Walk(read, func(n *ast.Node, entering bool) ast.WalkStatus {
					if entering && n.IsBlock() {
						t.Errorf("cell projection contains a block: %s", n.Type)
					}
					return ast.WalkContinue
				})
			}
			md := engine.BlockDOM2StdMd(dom)
			if strings.Contains(md, ast.TableCellRichIAL) || strings.Contains(md, ast.TableCellRichAttribute) {
				t.Fatalf("standard Markdown leaked internal metadata: %s", md)
			}
			jsonData := render.NewJSONRenderer(engine.BlockDOM2Tree(dom), engine.RenderOptions, engine.ParseOptions).Render()
			if !strings.Contains(string(jsonData), "\"TableCellRich\"") {
				t.Fatalf("JSON dropped rich source: %s", jsonData)
			}
		})
	}
}

func TestTableRichClipboardAndMergedExport(t *testing.T) {
	engine := newTableRichLute()
	tree := parse.Parse("", []byte("| first | second |\n| --- | --- |\n| old | hidden |"), engine.ParseOptions)
	cell := tree.Root.FirstChild.LastChild.FirstChild
	cell.TableCellRich = &ast.TableCellRich{Spec: 1, Format: "kramdown", Content: "- first\n- second\n\n```text\na | b\nc\n```"}
	cell.SetIALAttr("colspan", "2")
	cell.Next.SetIALAttr("class", "fn__none")
	if err := parse.ApplyTableCellRichProjection(cell); nil != err {
		t.Fatal(err)
	}
	dom := engine.Tree2BlockDOM(tree, engine.RenderOptions, engine.ParseOptions)
	clipboard := engine.HTML2BlockDOM(dom)
	restored := firstRichCell(engine.BlockDOM2Tree(clipboard))
	if restored == nil || restored.TableCellRich.Content != cell.TableCellRich.Content || restored.IALAttr("colspan") != "2" {
		t.Fatalf("HTML clipboard lost merged rich cell: %s", clipboard)
	}
	markdown := engine.BlockDOM2StdMd(dom)
	if !strings.Contains(markdown, "first") || !strings.Contains(markdown, "<code>a &#124; b<br />c</code>") || strings.Contains(markdown, "table-cell-rich") || strings.Contains(markdown, "<ul>") {
		t.Fatalf("merged Markdown export did not use its readable inline projection: %s", markdown)
	}
	html := engine.BlockDOM2RichHTML(dom)
	if !strings.Contains(html, "<ul>") || !strings.Contains(html, "<code") || strings.Contains(html, "data-sy-table-cell-rich") {
		t.Fatalf("HTML export lost the fragment structure: %s", html)
	}
}

func TestTableRichPreservesInvalidEmptyMetadata(t *testing.T) {
	engine := newTableRichLute()
	dom := `<div data-type="NodeTable"><div contenteditable="true"><table><thead><tr><th data-sy-table-cell-rich="">preserve</th></tr></thead></table></div></div>`
	for iteration := 0; iteration < 3; iteration++ {
		dom = engine.SpinBlockDOM(dom)
		cell := firstRichCell(engine.BlockDOM2Tree(dom))
		if cell == nil || cell.TableCellRich.Validate() == nil {
			t.Fatalf("corrupt metadata was silently discarded: %s", dom)
		}
	}
}

func TestTableRichCodeProjectionEscapesMarkup(t *testing.T) {
	engine := newTableRichLute()
	tree := parse.Parse("", []byte("| title |\n| --- |\n| value |"), engine.ParseOptions)
	cell := tree.Root.FirstChild.LastChild.FirstChild
	cell.TableCellRich = &ast.TableCellRich{Spec: 1, Format: "kramdown", Content: "```html\n<script>alert(1)</script>\na | b & c\n```"}
	if err := parse.ApplyTableCellRichProjection(cell); err != nil {
		t.Fatal(err)
	}
	dom := engine.Tree2BlockDOM(tree, engine.RenderOptions, engine.ParseOptions)
	inline := engine.BlockDOM2InlineBlockDOM(dom)
	markdown := engine.BlockDOM2StdMd(dom)
	for name, output := range map[string]string{"preview": dom, "inline conversion": inline, "Markdown": markdown} {
		if strings.Contains(output, "<script>") || !strings.Contains(output, "&lt;script&gt;") {
			t.Fatalf("%s interpreted literal code as markup: %s", name, output)
		}
	}
	if !strings.Contains(markdown, "a &#124; b &amp; c</code>") || strings.Contains(markdown, "\na |") {
		t.Fatalf("code content broke its Markdown table cell: %s", markdown)
	}
}

func TestTableRichDoesNotReinterpretLegacyCells(t *testing.T) {
	engine := newTableRichLute()
	dom := engine.Md2BlockDOM("| title |\n| --- |\n| - literal<br />1. literal |", false)
	if firstRichCell(engine.BlockDOM2Tree(engine.SpinBlockDOM(dom))) != nil {
		t.Fatal("legacy cell was implicitly converted")
	}
	if strings.Contains(dom, "<ul>") || strings.Contains(dom, "<ol>") {
		t.Fatal("legacy text became a list")
	}
}

func TestTableRichRejectsUnsupportedPayloads(t *testing.T) {
	for _, rich := range []*ast.TableCellRich{
		{Spec: 2, Format: "kramdown", Content: "preserve"},
		{Spec: 1, Format: "future", Content: "preserve"},
		{Spec: 1, Format: "kramdown", Content: "| nested |\n| --- |\n| table |"},
		{Spec: 1, Format: "kramdown", Content: "```mermaid\ngraph LR\n```"},
		ast.DecodeTableCellRich("invalid%payload"),
		{Spec: 1, Format: "kramdown", Content: "invalid\xfftext"},
	} {
		before := rich.Encode()
		if _, err := parse.ParseTableCellRich(rich); err == nil {
			t.Fatalf("accepted unsupported source: %#v", rich)
		}
		if rich.Encode() != before {
			t.Fatal("validation changed the original payload")
		}
	}
}
