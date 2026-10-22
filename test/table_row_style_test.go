package test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/render"
)

func TestTableRowStyleRoundTrip(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetKramdownIAL(true)
	for name, cells := range map[string]string{
		"ordinary": `<td>A</td><td>B</td>`,
		"merged":   `<td rowspan="2">A</td><td>B</td>`,
	} {
		t.Run(name, func(t *testing.T) {
			dom := `<div data-type="NodeTable"><div><table><thead><tr style="height: 60px;"><th>H1</th><th>H2</th></tr></thead>` +
				`<tbody><tr style="height: 80px; background-color: red;">` + cells + `</tr><tr><td>C</td><td>D</td></tr></tbody></table></div></div>`
			for round := 0; round < 3; round++ {
				dom = engine.SpinBlockDOM(dom)
				tree := engine.BlockDOM2Tree(dom)
				rows := 0
				ast.Walk(tree.Root, func(node *ast.Node, entering bool) ast.WalkStatus {
					if entering && ast.NodeTableRow == node.Type {
						expected := []string{"height: 60px;", "height: 80px; background-color: red;", ""}[rows]
						if got := node.IALAttr("style"); got != expected {
							t.Fatalf("row %d style = %q, want %q", rows, got, expected)
						}
						rows++
					}
					return ast.WalkContinue
				})
				if rows != 3 {
					t.Fatalf("row count = %d", rows)
				}
				data := render.NewJSONRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()
				var saved struct {
					Children []json.RawMessage
				}
				if err := json.Unmarshal(data, &saved); err != nil || !strings.Contains(string(data), `height: 80px; background-color: red;`) {
					t.Fatalf("JSON did not preserve row properties: %s, %v", data, err)
				}
				for format, output := range map[string]string{
					"HTML":         engine.BlockDOM2RichHTML(dom),
					"preview":      string(render.NewProtylePreviewRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()),
					"Word":         string(render.NewProtyleExportDocxRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()),
					"generic HTML": string(render.NewHtmlRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()),
				} {
					if !strings.Contains(output, `<tr style="height: 60px;">`) || !strings.Contains(output, `height: 80px; background-color: red;`) {
						t.Fatalf("%s lost row styles: %s", format, output)
					}
				}
				markdown := engine.BlockDOM2Md(dom)
				if !strings.Contains(markdown, `{: row-style="height: 60px;"}`) {
					t.Fatalf("internal Markdown lost row style: %s", markdown)
				}
				dom = engine.Md2BlockDOM(markdown, false)
			}
			standard := engine.BlockDOM2StdMd(dom)
			if strings.Contains(standard, "row-style=") {
				t.Fatalf("standard Markdown leaked internal row metadata: %s", standard)
			}
			if name == "ordinary" && strings.Contains(standard, "height:") {
				t.Fatalf("ordinary standard Markdown retained layout: %s", standard)
			}
		})
	}
}

func TestTableRowStyleLiteralAndLegacy(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetKramdownIAL(true)
	for _, body := range []string{"| a | b |", `| {: row-style="height: 80px;"} | b |`, `| a | b |{: row-style="80px"} literal`} {
		dom := engine.Md2BlockDOM("| H1 | H2 |\n| --- | --- |\n"+body, false)
		if strings.Contains(dom, `<tr style=`) {
			t.Fatalf("literal or legacy table gained a row style: %s", dom)
		}
	}
}

func TestTableRowHeightHTMLImport(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetKramdownIAL(true)
	dom := engine.HTML2BlockDOM(`<table><thead><tr style="height: 60px; color: red;"><th>Header</th></tr></thead><tbody><tr style="height: 80px;"><td>Body</td></tr></tbody></table>`)
	for _, expected := range []string{`<tr style="height: 60px;">`, `<tr style="height: 80px;">`} {
		if !strings.Contains(dom, expected) {
			t.Fatalf("HTML import lost row height %s: %s", expected, dom)
		}
	}
}
