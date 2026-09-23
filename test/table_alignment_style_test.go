package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func TestTableCellAlignmentUsesInlineStyle(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetKramdownIAL(true)
	legacy := `| Header | Other |` + "\n" + `| :---: | ---: |` + "\n" + `| value | next |`
	dom := engine.Md2BlockDOM(legacy, false)
	if strings.Contains(dom, ` align="`) || !strings.Contains(dom, `text-align: center`) || !strings.Contains(dom, `text-align: right`) {
		t.Fatalf("legacy Markdown alignment was not rendered as inline style: %s", dom)
	}
	for name, output := range map[string]string{
		"HTML": engine.BlockDOM2RichHTML(dom),
		"Word": string(render.NewProtyleExportDocxRenderer(engine.BlockDOM2Tree(dom), engine.RenderOptions, engine.ParseOptions).Render()),
	} {
		if strings.Contains(output, ` align="`) || !strings.Contains(output, `text-align: center`) || !strings.Contains(output, `text-align: right`) {
			t.Fatalf("%s lost inline alignment: %s", name, output)
		}
	}
}

func TestTableCellInlineStyleOverridesLegacyAlign(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetKramdownIAL(true)
	dom := `<div data-type="NodeTable"><div contenteditable="true"><table><thead><tr>` +
		`<th align="left" style="background-color: red; text-align: right;">Header</th>` +
		`</tr></thead><tbody><tr><td align="center" style="vertical-align: middle;">value</td></tr>` +
		`</tbody></table></div></div>`
	tree := engine.BlockDOM2Tree(dom)
	head := tree.Root.FirstChild.FirstChild.FirstChild.FirstChild
	if head.TableCellAlign != 3 {
		t.Fatalf("inline style did not override legacy align: %d", head.TableCellAlign)
	}
	output := string(render.NewProtyleExportDocxRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render())
	body := head.Parent.Parent.Next.FirstChild
	if got := body.IALAttr("style"); got != "vertical-align: middle;" {
		t.Fatalf("rendering changed the source cell style: %q", got)
	}
	for _, expected := range []string{`background-color: red`, `text-align: right`, `vertical-align: middle`, `text-align: center`} {
		if !strings.Contains(output, expected) {
			t.Fatalf("Word output omitted %q: %s", expected, output)
		}
	}
	if strings.Contains(output, ` align="`) {
		t.Fatalf("Word output retained a legacy align attribute: %s", output)
	}
	markdown := engine.BlockDOM2Md(dom)
	roundTrip := engine.Md2BlockDOM(markdown, false)
	if !strings.Contains(roundTrip, `text-align: right`) || !strings.Contains(roundTrip, `vertical-align: middle`) {
		t.Fatalf("cell styles were lost during Markdown round trip: %s", markdown)
	}
}

func TestHTMLTableCellAlignmentDifferencesSurviveMarkdown(t *testing.T) {
	engine := lute.New()
	input := `<table><thead><tr><th style="text-align: center">Header</th><th>Other</th></tr></thead>` +
		`<tbody><tr><td style="text-align: right">different</td><td>same</td></tr></tbody></table>`
	markdown, err := engine.HTML2Markdown(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(markdown, `text-align: right;`) || !strings.Contains(markdown, `| :-`) {
		t.Fatalf("per-cell alignment or column delimiter was lost: %s", markdown)
	}
	dom := engine.Md2BlockDOM(markdown, false)
	if !strings.Contains(dom, `text-align: right;`) || !strings.Contains(dom, `text-align: center;`) {
		t.Fatalf("per-cell alignment was lost after Markdown parsing: %s", dom)
	}
}

func TestLegacyTableCellAlignmentSurvivesMarkdownExport(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetKramdownIAL(true)
	tree := parse.Parse("", []byte("| Header |\n| :---: |\n| value |\n"), engine.ParseOptions)
	table := tree.Root.FirstChild
	cell := table.FirstChild.Next.FirstChild
	cell.TableCellAlign = 3
	markdown := string(render.NewProtyleExportMdRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render())
	if !strings.Contains(markdown, "text-align: right;") {
		t.Fatalf("legacy per-cell alignment was lost: %s", markdown)
	}
	roundTrip := parse.Parse("", []byte(markdown), engine.ParseOptions)
	if got := roundTrip.Root.FirstChild.FirstChild.Next.FirstChild.TableCellAlign; got != 3 {
		t.Fatalf("legacy per-cell alignment changed after export: %d", got)
	}
}
