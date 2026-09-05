package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
)

func tabsEngine() *lute.Lute {
	l := lute.New()
	l.SetTabs(true)
	l.SetKramdownBlockIAL(true)
	l.SetBlockRef(true)
	l.SetSuperBlock(true)
	l.SetTextMark(true)
	l.SetProtyleWYSIWYG(true)
	return l
}

func tabNodes(root *ast.Node) (tabs, items []*ast.Node) {
	ast.Walk(root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			if n.Type == ast.NodeTabs {
				tabs = append(tabs, n)
			}
			if n.Type == ast.NodeTabItem {
				items = append(items, n)
			}
		}
		return ast.WalkContinue
	})
	return
}

func TestTabsRoundTrip(t *testing.T) {
	l := tabsEngine()
	md := ":::tabs\n:::tab **First**\nOne\n:::tabs\n:::tab Nested\n```\n:::tab literal\n:::\n```\n:::\n:::\n:::\n{: id=\"20260905120000-item001\"}\n:::tab Second\nTwo\n:::\n{: id=\"20260905120000-item002\"}\n:::\n{: id=\"20260905120000-tabs001\" tabs-active-id=\"20260905120000-item002\" tabs-position=\"left\"}\n"
	dom := l.Md2BlockDOM(md, false)
	for round := 0; round < 3; round++ {
		tree := l.BlockDOM2Tree(dom)
		tabs, items := tabNodes(tree.Root)
		if len(tabs) != 2 || len(items) != 3 {
			t.Fatalf("round %d structure: %s", round, dom)
		}
		if !strings.Contains(items[0].TabItemTitle, `data-type="strong">First`) || items[2].ID != "20260905120000-item002" || tabs[0].IALAttr("tabs-active-id") != items[2].ID || tabs[0].IALAttr("tabs-position") != "left" {
			t.Fatalf("round %d metadata: title=%q item=%q active=%q position=%q", round, items[0].TabItemTitle, items[2].ID, tabs[0].IALAttr("tabs-active-id"), tabs[0].IALAttr("tabs-position"))
		}
		if !strings.Contains(dom, "literal") || !strings.Contains(dom, "Two") {
			t.Fatal(dom)
		}
		dom = l.SpinBlockDOM(dom)
	}
}

func TestTabsContainersAndFences(t *testing.T) {
	l := tabsEngine()
	for _, md := range []string{
		":::tabs\n:::tab A\nbody\n:::\n:::\n",
		"> :::tabs\n> :::tab A\n> body\n> :::\n> :::\n",
		"- :::tabs\n  :::tab A\n  body\n  :::\n  :::\n",
		"{{{row\n:::tabs\n:::tab A\nbody\n:::\n:::\n}}}\n",
		":::tabs\n:::tab A\nbody",
	} {
		tree := parse.Parse("", []byte(md), l.ParseOptions)
		tabs, items := tabNodes(tree.Root)
		if len(tabs) != 1 || len(items) != 1 || items[0].Parent != tabs[0] {
			t.Fatalf("invalid structure: %q", md)
		}
		formatted := l.FormatStr("", md)
		tree = parse.Parse("", []byte(formatted), l.ParseOptions)
		tabs, items = tabNodes(tree.Root)
		if len(tabs) != 1 || len(items) != 1 {
			t.Fatalf("format lost structure: %q", formatted)
		}
	}
	for _, md := range []string{":::tab stray\nbody\n:::\n", "```\n:::tabs\n:::tab A\n:::\n```\n"} {
		tabs, _ := tabNodes(parse.Parse("", []byte(md), l.ParseOptions).Root)
		if len(tabs) != 0 {
			t.Fatalf("unexpected tabs: %q", md)
		}
	}
}

func TestTabItemFragmentAndExport(t *testing.T) {
	l := tabsEngine()
	dom := `<div data-node-id="20260905120000-item001" data-type="NodeTabItem" class="tab-item"><div class="tab-item-info callout-info"><span class="tab-item-title callout-title"><span data-type="strong">Title</span></span></div><div class="tab-item-content"><div data-node-id="20260905120000-para001" data-type="NodeParagraph" class="p"><div contenteditable="true">Body</div></div></div></div>`
	spun := l.SpinBlockDOM(dom)
	tabs, items := tabNodes(l.BlockDOM2Tree(spun).Root)
	if len(tabs) != 0 || len(items) != 1 || !strings.Contains(items[0].TabItemTitle, `data-type="strong">Title`) || items[0].FirstChild.ID != "20260905120000-para001" {
		t.Fatalf("tabs=%d items=%d title=%q body=%q", len(tabs), len(items), items[0].TabItemTitle, items[0].FirstChild.ID)
	}
	md := l.BlockDOM2StdMd(dom)
	if strings.Contains(md, ":::tab") || !strings.Contains(md, "**Title**") || !strings.Contains(md, "Body") {
		t.Fatal(md)
	}
}

func TestTabsHTMLAndEmptyInput(t *testing.T) {
	l := tabsEngine()
	md := ":::tabs\n:::tab **Title**\nBody\n:::\n:::tab\n:::\n:::\n"
	html := l.MarkdownStr("", md)
	back := l.HTML2Md(html)
	tabs, items := tabNodes(parse.Parse("", []byte(back), l.ParseOptions).Root)
	if len(tabs) != 1 || len(items) != 2 || !strings.Contains(items[0].TabItemTitle, "Title") || items[1].TabItemTitle != "" || items[1].FirstChild == nil {
		t.Fatalf("HTML round trip: tabs=%d items=%d title=%q empty=%q child=%v\n%s", len(tabs), len(items), items[0].TabItemTitle, items[1].TabItemTitle, items[1].FirstChild, back)
	}
	if strings.Contains(html, "contenteditable") || strings.Contains(html, "display: none") {
		t.Fatal("HTML must work without enhancement")
	}
	if strings.Contains(l.BlockDOM2StdMd(l.Md2BlockDOM(md, false)), ":::tab") {
		t.Fatal("standard markdown contains fences")
	}
}

func TestTabsEmptyParagraphIDs(t *testing.T) {
	l := tabsEngine()
	md := ":::tabs\n:::tab\n\n{: id=\"20260905120000-empty01\"}\n\n{: id=\"20260905120000-empty02\"}\n:::\n{: id=\"20260905120000-item001\"}\n:::\n{: id=\"20260905120000-tabs001\"}\n"
	dom := l.Md2BlockDOM(md, false)
	for round := 0; round < 3; round++ {
		tree := l.BlockDOM2Tree(dom)
		tabs, items := tabNodes(tree.Root)
		if len(tabs) != 1 || len(items) != 1 || tabs[0].ID != "20260905120000-tabs001" || items[0].ID != "20260905120000-item001" {
			t.Fatalf("round %d containers: %s", round, dom)
		}
		var paragraphs []string
		for child := items[0].FirstChild; child != nil; child = child.Next {
			if child.Type == ast.NodeParagraph {
				paragraphs = append(paragraphs, child.ID)
			}
		}
		if strings.Join(paragraphs, ",") != "20260905120000-empty01,20260905120000-empty02" {
			t.Fatalf("round %d empty paragraph IDs: %v\n%s", round, paragraphs, dom)
		}
		dom = l.SpinBlockDOM(dom)
	}
}
