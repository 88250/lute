package test

import (
	"strings"
	"testing"

	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func TestTabsStandardExportLiteralFences(t *testing.T) {
	for name, markdown := range map[string]string{
		"body":  "::: tabs\n@tab First\n\\::: tabs\nLiteral body\n\\:::\n@tab Second\nMore body\n:::\n",
		"title": "::: tabs\n@tab \\::: tabs\nLiteral body\n@tab Second\nMore body\n:::\n",
	} {
		for _, plainDOM := range []bool{false, true} {
			suffix := "/escaped"
			if plainDOM {
				suffix = "/plain"
			}
			t.Run(name+suffix, func(t *testing.T) {
				l := tabsEngine()
				dom := l.Md2BlockDOM(markdown, false)
				if plainDOM {
					// 模拟编辑器中没有转义标记节点的普通文字。
					dom = strings.ReplaceAll(dom, `<span data-type="backslash">:</span>`, ":")
				}
				exported := string(render.NewProtyleExportMdRenderer(l.BlockDOM2Tree(dom), render.NewOptions(), l.ParseOptions).Render())
				if !strings.Contains(exported, "\\::: tabs") || strings.Contains(exported, "\\\\::: tabs") {
					t.Fatalf("literal fence must be escaped exactly once\n%s", exported)
				}
				tabs, _ := tabNodes(parse.Parse("", []byte(exported), l.ParseOptions).Root)
				if len(tabs) != 0 {
					t.Fatalf("flattened export created %d tabs containers\n%s", len(tabs), exported)
				}
				html := l.MarkdownStr("", exported)
				for _, content := range []string{"::: tabs", "Literal body", "Second", "More body"} {
					if !strings.Contains(html, content) {
						t.Fatalf("export lost %q\n%s", content, exported)
					}
				}
			})
		}
	}
}

func TestTabsLiteralMarkerCaretSpin(t *testing.T) {
	for _, marker := range []string{"@tab Literal body", "@tab:active Literal body", "::: tabs", ":::"} {
		t.Run(marker, func(t *testing.T) {
			l := tabsEngine()
			markdown := "::: tabs\n@tab First\n\\" + marker + "\n@tab Second\nMore body\n:::\n"
			dom := l.Md2BlockDOM(markdown, false)
			dom = strings.ReplaceAll(dom, `<span data-type="backslash">`+marker[:1]+`</span>`, "<wbr>"+marker[:1])
			for round := 0; round < 3; round++ {
				dom = l.SpinBlockDOM(dom)
				tabs, items := tabNodes(l.BlockDOM2Tree(dom).Root)
				if len(tabs) != 1 || len(items) != 2 || items[1].TabItemTitle != "Second" {
					t.Fatalf("round %d caret changed structure\n%s", round, dom)
				}
				if strings.Count(dom, "<wbr>") != 1 || strings.Contains(dom, "\\") {
					t.Fatalf("round %d caret lost or escape became visible\n%s", round, dom)
				}
			}
		})
	}
}
