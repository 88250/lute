package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
)

func TestMarkdownTableHTMLThemeStyles(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetHTMLTag2TextMark(true)
	engine.SetTextMark(true)
	engine.SetKramdownIAL(true)
	engine.SetAutoSpace(false)
	for _, tc := range []struct {
		name, cell string
		want       []string
	}{
		{"issue", `<span style="color: var(--b3-font-color6);">测试文本 ⭐⭐⭐ / 测试文本</span>⭐⭐⭐<br />测试文本<span style="color: var(--b3-font-color6);">测试文本 </span>⭐`, []string{`style="color: var(--b3-font-color6);"`, `测试文本 ⭐⭐⭐ / 测试文本`, `<br />`, `测试文本 </span>⭐`}},
		{"attributes", `<SPAN title='' class='custom' STYLE='color: var(--b3-font-color10);background:var(--b3-font-background3)'>text</SPAN>`, []string{`color: var(--b3-font-color10);background-color: var(--b3-font-background3);`, `>text</span>`}},
		{"nested", `<span style="color:var(--b3-font-color6)">outer <span style="color:var(--b3-font-color1)">inner</span> tail **bold** [link](https://example.com)</span>`, []string{`style="color: var(--b3-font-color1);">inner</span>`, `style="color: var(--b3-font-color6);"> tail `, `data-type="strong"`, `data-href="https://example.com"`}},
		{"native", `<span data-type="text" style="color: var(--b3-font-color10);">native</span> <span style="color:var(--b3-font-color6)">plain</span>`, []string{`style="color: var(--b3-font-color10);">native</span>`, `style="color: var(--b3-font-color6);">plain</span>`}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			markdown := "| Header |\n| --- |\n| " + tc.cell + " |\n"
			for _, convert := range []func(string, bool) string{engine.Md2BlockDOM, engine.Md2BlockDOMWithAutoLink} {
				dom := convert(markdown, false)
				for round := 0; round < 2; round++ {
					if strings.Contains(dom, "&lt;span") || strings.Count(dom, `data-type="NodeTable"`) != 1 {
						t.Fatalf("unparsed span or broken table: %s", dom)
					}
					for _, want := range tc.want {
						if !strings.Contains(dom, want) {
							t.Errorf("round %d lost %q: %s", round, want, dom)
						}
					}
					dom = engine.SpinBlockDOM(dom)
				}
			}
		})
	}
}

func TestMarkdownTableHTMLThemeStyleBoundaries(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetHTMLTag2TextMark(true)
	engine.SetTextMark(true)
	for _, input := range []string{
		"| Header |\n| --- |\n| `<span style=\"color:var(--b3-font-color6)\">code</span>` |",
		"```html\n<span style=\"color:var(--b3-font-color6)\">code</span>\n```",
		"| Header |\n| --- |\n| <span style=\"color:var(--b3-font-color6)\">unclosed |",
		"| Header |\n| --- |\n| &lt;span style=\"color:var(--b3-font-color6)\"&gt;literal&lt;/span&gt; |",
	} {
		dom, tree := engine.Md2BlockDOMTree(input, false)
		if !strings.Contains(dom, "&lt;span") {
			t.Errorf("literal HTML was lost: %s", dom)
		}
		ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
			if entering && n.IALAttr("style") != "" {
				t.Errorf("unexpected style in %q: %s", input, n.IALAttr("style"))
			}
			return ast.WalkContinue
		})
	}
}
