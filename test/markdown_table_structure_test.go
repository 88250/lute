package test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
)

func tableStructure(tree *parse.Tree) []string {
	var ret []string
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if n.Type == ast.NodeKramdownBlockIAL || n.Type == ast.NodeKramdownSpanIAL {
			return ast.WalkSkipChildren
		}
		if n.IsBlock() {
			if entering {
				ret = append(ret, n.Type.String())
			} else {
				ret = append(ret, "/"+n.Type.String())
			}
		}
		if entering {
			switch n.Type {
			case ast.NodeTaskListItemMarker:
				ret = append(ret, "task:"+n.EffectiveTaskListItemMarker())
			case ast.NodeCodeBlockCode, ast.NodeMathBlockContent:
				ret = append(ret, n.Type.String()+":"+n.TokensStr())
			}
		}
		return ast.WalkContinue
	})
	return ret
}

func TestMarkdownTableStructureRoundTrip(t *testing.T) {
	engine := newTableRichLute()
	for _, source := range []string{
		"- [ ] first\n- [x] checked\n- [/] progress\n- [X] upper",
		"- outer\n  - nested\n    - [ ] task\n- tail",
		"first paragraph\n\nsecond **bold** paragraph\n\n> quote\n\n# heading",
		"```go\na | b\n<span>literal</span> & &#10;\n```",
		"$$\na | b\n$$",
		"7. seven\n8. eight",
		"- [ ] <span style=\"color:var(--b3-font-color6)\">colored</span>\n\n![image](assets/a.png)",
	} {
		t.Run(source, func(t *testing.T) {
			tree := parse.Parse("", []byte("| head |\n| --- |\n| old |"), engine.ParseOptions)
			cell := tree.Root.FirstChild.LastChild.FirstChild
			cell.TableCellRich = &ast.TableCellRich{Spec: 1, Format: "kramdown", Content: source}
			if err := parse.ApplyTableCellRichProjection(cell); err != nil {
				t.Fatal(err)
			}
			original, _ := parse.ParseTableCellRich(cell.TableCellRich)
			dom := engine.Tree2BlockDOM(tree, engine.RenderOptions, engine.ParseOptions)
			for round := 0; round < 3; round++ {
				markdown := engine.BlockDOM2StdMd(dom)
				if strings.Contains(markdown, "table-cell-rich") || len(strings.Split(strings.TrimSpace(markdown), "\n")) != 3 {
					t.Fatalf("nonstandard metadata or broken row: %s", markdown)
				}
				dom = engine.Md2BlockDOM(markdown, false)
				restored := firstRichCell(engine.BlockDOM2Tree(dom))
				if restored == nil {
					t.Fatalf("rich cell missing after %s: %s", markdown, dom)
				}
				fragment, err := parse.ParseTableCellRich(restored.TableCellRich)
				if err != nil {
					t.Fatal(err)
				}
				if got, want := tableStructure(fragment), tableStructure(original); !reflect.DeepEqual(got, want) {
					t.Fatalf("round %d\ngot %q\nwant %q\nsource: %s", round, got, want, restored.TableCellRich.Content)
				}
			}
		})
	}
}

func TestMarkdownTableLegacyTasks(t *testing.T) {
	engine := newTableRichLute()
	for _, source := range []string{
		"- [ ] task",
		"- [x] done<br />- [ ] todo",
		"- [/] custom<br/>  - [ ] nested",
		"- [ ] `<br />` literal",
	} {
		dom := engine.Md2BlockDOM("| head |\n| --- |\n| "+source+" |", false)
		cell := firstRichCell(engine.BlockDOM2Tree(dom))
		if cell == nil || !strings.Contains(dom, `type="checkbox"`) {
			t.Fatalf("task not restored: %s", dom)
		}
		if strings.Contains(source, "`<br />`") && !strings.Contains(cell.TableCellRich.Content, "`<br />`") {
			t.Fatalf("code break was expanded: %s", cell.TableCellRich.Content)
		}
	}
}

func TestHTMLTableStructureImport(t *testing.T) {
	engine := newTableRichLute()
	engine.SetProtyleWYSIWYG(true)
	dom := engine.HTML2BlockDOM(`<table><tr><th>head</th></tr><tr><td colspan="2"><ul><li>outer<ul><li>nested</li></ul></li></ul><pre><code class="language-go">a | b
c
</code></pre></td></tr></table>`)
	cell := firstRichCell(engine.BlockDOM2Tree(dom))
	if cell == nil || cell.IALAttr("colspan") != "2" {
		t.Fatalf("merged rich cell missing: %s", dom)
	}
	if !strings.Contains(cell.TableCellRich.Content, "a | b\nc\n") || !strings.Contains(dom, "<ul>") {
		t.Fatalf("block structure lost: %s", dom)
	}
}

func TestMarkdownTableStructureBoundaries(t *testing.T) {
	engine := newTableRichLute()
	for _, source := range []string{
		"`- [ ] literal`", `\- [ ] literal`, "- \\[ ] literal", "ordinary", "- literal<br />1. literal",
		"`<ul><li>code</li></ul>`", "&lt;ul&gt;&lt;li&gt;literal&lt;/li&gt;&lt;/ul&gt;",
		"<ul><li>unclosed", "<ul><li>mismatch</ul></li>", "<table><tr><td>nested</td></tr></table>",
		"<p>keep</p><script>alert(1)</script>", "<p onclick='alert(1)'>keep</p>",
	} {
		dom := engine.Md2BlockDOM("| head |\n| --- |\n| "+source+" |", false)
		if firstRichCell(engine.BlockDOM2Tree(dom)) != nil {
			t.Errorf("unexpected rich conversion of %s: %s", source, dom)
		}
	}
	dom := `<div data-type="NodeTable"><div contenteditable="true"><table><thead><tr><th>head</th></tr></thead><tbody><tr><td>- [ ] existing text</td></tr></tbody></table></div></div>`
	if firstRichCell(engine.BlockDOM2Tree(engine.SpinBlockDOM(dom))) != nil {
		t.Fatal("existing ordinary cell was reinterpreted")
	}
}
