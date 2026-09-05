package test

import (
	"strings"
	"testing"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func TestTabsOriginalTitleBlock(t *testing.T) {
	l := tabsEngine()
	dom := `<div data-node-id="20260905120000-tabs001" data-type="NodeTabs" class="tabs" tabs-active-id="20260905120000-item001"><div data-node-id="20260905120000-item001" data-type="NodeTabItem" class="tab-item"><div class="tab-item-info callout-info"><div data-node-id="20260905120000-title01" data-type="NodeParagraph" class="p" tabs-title="true" custom-test="kept"><div class="tab-item-title callout-title" contenteditable="true"><span data-type="strong">abcdefghijklmnop</span> <span data-type="block-ref" data-id="20260905120000-ref0001" data-subtype="d">Reference</span></div></div></div><div class="tab-item-content"><div data-node-id="20260905120000-para001" data-type="NodeParagraph" class="p"><div contenteditable="true">Body</div></div></div></div></div>`
	for round := 0; round < 4; round++ {
		tree := l.BlockDOM2Tree(dom)
		_, items := tabNodes(tree.Root)
		if len(items) != 1 {
			t.Fatal(dom)
		}
		title := items[0].TabTitleBlock()
		if title == nil || title.ID != "20260905120000-title01" || title.IALAttr("custom-test") != "kept" {
			t.Fatalf("round %d title identity lost\n%s", round, dom)
		}
		count := 0
		ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
			if entering && n.ID == title.ID {
				count++
			}
			return ast.WalkContinue
		})
		if count != 1 {
			t.Fatalf("duplicate title block: %d", count)
		}
		md := l.BlockDOM2StdMd(dom)
		if strings.Count(md, "abcdefghijklmnop") != 1 || !strings.Contains(md, "@tab:active **abcdefghijklmnop**") || strings.Contains(md, "custom-test") || strings.Contains(md, "tabs-title") || strings.Contains(md, "{: ") {
			t.Fatalf("clipboard lost structure or duplicated title\n%s", md)
		}
		copiedTabs, copiedItems := tabNodes(parse.Parse("", []byte(md), l.ParseOptions).Root)
		if len(copiedTabs) != 1 || len(copiedItems) != 1 || !strings.Contains(copiedItems[0].TabItemTitle, "abcdefghijklmnop") {
			t.Fatal(md)
		}
		options := render.NewOptions()
		flat := string(render.NewProtyleExportMdRenderer(tree, options, l.ParseOptions).Render())
		if strings.Count(flat, "abcdefghijklmnop") != 1 || strings.Contains(flat, "@tab") {
			t.Fatalf("flat export\n%s", flat)
		}
		html := l.BlockDOM2HTML(dom)
		if strings.Count(html, "abcdefghijklmnop") != 1 || !strings.Contains(html, "Body") {
			t.Fatalf("HTML duplicated or lost title\n%s", html)
		}
		back := l.HTML2Md(html)
		_, htmlItems := tabNodes(parse.Parse("", []byte(back), l.ParseOptions).Root)
		if len(htmlItems) != 1 || (!strings.Contains(htmlItems[0].TabItemTitle, "abcdefghijklmnop") && htmlItems[0].TabTitleBlock() == nil) {
			t.Fatalf("HTML round trip lost title\n%s", back)
		}
		if strings.Contains(htmlItems[0].TabItemTitle, "{: ") || strings.Count(back, "abcdefghijklmnop") != 1 {
			t.Fatalf("HTML title leaked block attributes or duplicated content\n%s", back)
		}
		dom = l.SpinBlockDOM(dom)
		if strings.Count(dom, `data-node-id="20260905120000-title01"`) != 1 || !strings.Contains(dom, `class="tab-item-title callout-title"`) {
			t.Fatal(dom)
		}
	}
}

func TestTabsClipboardNestingAndLiteralIAL(t *testing.T) {
	l := tabsEngine()
	md := ":::: tabs\n@tab First\n\n::: tabs\n@tab Inner\nA\n@tab:active Selected\nB\n:::\n@tab:active Second\n\n```\n{: id=\"literal\"}\n@tab literal\n::: tabs\n```\n::::\n"
	copied := l.BlockDOM2StdMd(l.Md2BlockDOM(md, false))
	tabs, items := tabNodes(parse.Parse("", []byte(copied), l.ParseOptions).Root)
	if len(tabs) != 2 || len(items) != 4 || !strings.Contains(copied, ":::: tabs") || strings.Count(copied, "@tab:active") != 2 {
		t.Fatal(copied)
	}
	for _, group := range tabs {
		var last *ast.Node
		for n := group.FirstChild; nil != n; n = n.Next {
			if n.Type == ast.NodeTabItem {
				last = n
			}
		}
		if last == nil || group.IALAttr("tabs-active-id") != last.ID {
			t.Fatal("clipboard lost nested selection")
		}
	}
	if !strings.Contains(copied, "{: id=\"literal\"}\n@tab literal\n::: tabs") || strings.Count(copied, "{: ") != 1 {
		t.Fatal("clipboard must preserve literal IAL in code without leaking block attributes", copied)
	}
}

func TestTabsEmptyAndMultilineTitleIdentity(t *testing.T) {
	l := tabsEngine()
	for _, content := range []string{"", "First<br />Second", "😀😀😀😀😀😀😀😀😀😀😀😀😀", "<span data-type=\"code\">@tab literal</span>"} {
		dom := `<div data-node-id="20260905120000-item001" data-type="NodeTabItem" class="tab-item"><div class="tab-item-info"><div data-node-id="20260905120000-title01" data-type="NodeParagraph" class="p" tabs-title="true"><div class="tab-item-title callout-title" contenteditable="true">` + content + `</div></div></div><div class="tab-item-content"><div data-node-id="20260905120000-empty01" data-type="NodeParagraph" class="p" tabs-placeholder="true"><div contenteditable="true"></div></div></div></div>`
		for round := 0; round < 3; round++ {
			dom = l.SpinBlockDOM(dom)
			_, items := tabNodes(l.BlockDOM2Tree(dom).Root)
			if len(items) != 1 || items[0].TabTitleBlock() == nil || items[0].TabTitleBlock().ID != "20260905120000-title01" || !strings.Contains(dom, `tabs-placeholder="true"`) {
				t.Fatalf("title or empty body lost in round %d\n%s", round, dom)
			}
		}
	}
}
