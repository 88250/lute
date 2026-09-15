package test

import (
	"strings"
	"testing"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func TestTabsExportNestedTitleBlocks(t *testing.T) {
	l := tabsEngine()
	tree := parse.Parse("", []byte(":::: tabs\n@tab Outer\n\n::: tabs\n@tab Inner\n\n- Inner\n\n:::\n@tab StringTitle\n\nBody\n\n::::\n"), l.ParseOptions)
	_, items := tabNodes(tree.Root)
	if len(items) != 3 {
		t.Fatalf("expected three tab items, got %d", len(items))
	}
	for _, item := range items[:2] {
		title := parse.Parse("", []byte(item.TabItemTitle), l.ParseOptions).Root.FirstChild
		title.SetIALAttr("tabs-title", "true")
		item.PrependChild(title)
		item.TabItemTitle = ""
	}
	exported := string(render.NewProtyleExportRenderer(tree, render.NewOptions(), l.ParseOptions).Render())
	for content, count := range map[string]int{"Outer": 1, "Inner": 2, "StringTitle": 1, "Body": 1} {
		if strings.Count(exported, content) != count {
			t.Fatalf("expected %q %d times\n%s", content, count, exported)
		}
	}
	if !strings.Contains(exported, `data-type="NodeList"`) || !strings.Contains(exported, `data-type="NodeListItem"`) {
		t.Fatalf("export lost nested list\n%s", exported)
	}
	for _, item := range items[:2] {
		if title := item.TabTitleBlock(); title == nil || title.Type != ast.NodeParagraph {
			t.Fatal("export changed title blocks")
		}
	}
}

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
