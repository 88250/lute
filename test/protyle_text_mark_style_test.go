package test

import (
	"reflect"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
)

func TestProtyleTextMarkStyleBoundaries(t *testing.T) {
	engine := lute.New()
	engine.SetTextMark(true)
	engine.SetHTMLTag2TextMark(true)
	engine.SetKramdownIAL(true)
	engine.SetProtyleWYSIWYG(true)
	for _, typ := range []string{"text", "text strong", "text a"} {
		for _, style := range []string{
			"background-color: var(--b3-font-background5);", "color: red;", "font-size: 2em;",
		} {
			for _, separator := range []string{"", "\u200b", "\u200b‸"} {
				t.Run(typ+"/"+style+"/"+separator, func(t *testing.T) {
					span := func(text, style string) string {
						return `<span data-type="` + typ + `" data-href="https://example.com" style="` + style + `">` + text + `</span>`
					}
					dom := `<div data-node-id="20261020000000-abcdefg" data-type="NodeParagraph"><div contenteditable="true">` +
						span("before", "") + separator + span("marked", style) + span("after", "") + `</div></div>`
					for round := 0; round < 3; round++ {
						for name, converted := range map[string]string{
							"input": dom, "spin": engine.SpinBlockDOM(dom),
							"markdown": engine.Md2BlockDOM(engine.BlockDOM2Md(dom), false),
						} {
							var styled []string
							ast.Walk(engine.BlockDOM2Tree(converted).Root, func(node *ast.Node, entering bool) ast.WalkStatus {
								if entering && node.Type == ast.NodeTextMark && node.IALAttr("style") != "" {
									styled = append(styled, node.TextMarkTextContent)
								}
								return ast.WalkContinue
							})
							if !reflect.DeepEqual(styled, []string{"marked"}) {
								t.Fatalf("%s round %d styled content = %q", name, round, styled)
							}
						}
						dom = engine.Md2BlockDOM(engine.BlockDOM2Md(engine.SpinBlockDOM(dom)), false)
					}
				})
			}
		}
	}
}

func TestProtyleTextMarkMergesMatchingStyles(t *testing.T) {
	engine := lute.New()
	engine.SetTextMark(true)
	engine.SetKramdownIAL(true)
	engine.SetProtyleWYSIWYG(true)
	for _, style := range []string{"", "color: red;"} {
		dom := `<div data-node-id="20261020000000-abcdefg" data-type="NodeParagraph"><div contenteditable="true">` +
			`<span data-type="text" style="` + style + `">first</span>` +
			`<span data-type="text" style="` + style + `">second</span></div></div>`
		var content []string
		ast.Walk(engine.BlockDOM2Tree(dom).Root, func(node *ast.Node, entering bool) ast.WalkStatus {
			if entering && node.Type == ast.NodeTextMark {
				content = append(content, node.TextMarkTextContent)
			}
			return ast.WalkContinue
		})
		if !reflect.DeepEqual(content, []string{"firstsecond"}) {
			t.Fatalf("matching style %q did not merge: %q", style, content)
		}
	}
}
