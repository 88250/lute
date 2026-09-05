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
	md := ":::: tabs\n@tab **First**\n{: id=\"20260905120000-item001\"}\n\nOne\n::: tabs\n@tab Nested\n```\n@tab literal\n:::\n```\n:::\n@tab:active Second\n{: id=\"20260905120000-item002\"}\n\nTwo\n::::\n{: id=\"20260905120000-tabs001\" tabs-active-id=\"20260905120000-item002\" tabs-position=\"left\"}\n"
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
		"::: tabs\n@tab A\nbody\n:::\n",
		":::\ttabs\n@tab\tA\nbody\n:::\n",
		"> ::: tabs\n> @tab A\n> body\n> :::\n",
		"- ::: tabs\n  @tab A\n  body\n  :::\n",
		"{{{row\n::: tabs\n@tab A\nbody\n:::\n}}}\n",
		"::: tabs\n@tab A\nbody",
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
	for _, md := range []string{
		"@tab stray\nbody\n:::\n",
		"@tab:active stray\nbody\n:::\n",
		":::tabs\n:::tab A\nbody\n:::\n:::\n",
		":: tabs\n@tab A\nbody\n::\n",
		"::: tabs extra\n@tab A\nbody\n:::\n",
		"```\n::: tabs\n@tab A\n:::\n```\n",
	} {
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
	if len(tabs) != 0 || len(items) != 1 {
		t.Fatalf("tabs=%d items=%d\n%s", len(tabs), len(items), spun)
	}
	if items[0].ID != "20260905120000-item001" || !strings.Contains(items[0].TabItemTitle, `data-type="strong">Title`) || items[0].FirstChild.ID != "20260905120000-para001" {
		t.Fatalf("item=%q title=%q body=%q", items[0].ID, items[0].TabItemTitle, items[0].FirstChild.ID)
	}
	md := l.BlockDOM2StdMd(dom)
	if strings.Contains(md, "::: tabs") || strings.Contains(md, "@tab") || !strings.Contains(md, "**Title**") || !strings.Contains(md, "Body") {
		t.Fatal(md)
	}
}

func TestTabsHTMLAndEmptyInput(t *testing.T) {
	l := tabsEngine()
	md := "::: tabs\n@tab **Title**\nBody\n@tab\n:::\n"
	html := l.MarkdownStr("", md)
	back := l.HTML2Md(html)
	tabs, items := tabNodes(parse.Parse("", []byte(back), l.ParseOptions).Root)
	if len(tabs) != 1 || len(items) != 2 {
		t.Fatalf("HTML round trip: tabs=%d items=%d\n%s", len(tabs), len(items), back)
	}
	if !strings.Contains(items[0].TabItemTitle, "Title") || items[1].TabItemTitle != "" || items[1].FirstChild == nil {
		t.Fatalf("HTML round trip: title=%q empty=%q child=%v\n%s", items[0].TabItemTitle, items[1].TabItemTitle, items[1].FirstChild, back)
	}
	if strings.Contains(html, "contenteditable") || strings.Contains(html, "display: none") {
		t.Fatal("HTML must work without enhancement")
	}
	if strings.Contains(l.BlockDOM2StdMd(l.Md2BlockDOM(md, false)), "::: tabs") {
		t.Fatal("standard markdown contains fences")
	}
}

func TestTabsEmptyParagraphIDs(t *testing.T) {
	l := tabsEngine()
	md := "::: tabs\n@tab\n{: id=\"20260905120000-item001\"}\n\n{: id=\"20260905120000-empty01\"}\n\n{: id=\"20260905120000-empty02\"}\n:::\n{: id=\"20260905120000-tabs001\"}\n"
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

func TestTabsSlashPreset(t *testing.T) {
	l := tabsEngine()
	markdown := "::: tabs\n@tab\n‸\n@tab\n\n:::\n"
	dom := l.Md2BlockDOM(markdown, false)
	if strings.Contains(dom, ":::") {
		t.Fatal(dom)
	}
	tabs, items := tabNodes(l.BlockDOM2Tree(dom).Root)
	if len(tabs) != 1 || len(items) != 2 || nil == items[0].FirstChild || nil == items[1].FirstChild {
		t.Fatalf("tabs=%d items=%d\n%s", len(tabs), len(items), dom)
	}
	paragraph := `<div data-node-id="20260905134000-paragraph" data-type="NodeParagraph" class="p"><div contenteditable="true">` + markdown + `</div><div class="protyle-attr" contenteditable="false"></div></div>`
	spun := l.SpinBlockDOM(paragraph)
	if strings.Contains(spun, ":::") || !strings.Contains(spun, "<wbr>") {
		t.Fatal(spun)
	}
}

func TestTabsNestedFenceLengthsAndIndent(t *testing.T) {
	for _, indent := range []string{"", "  ", "    "} {
		t.Run("indent"+strings.ReplaceAll(indent, " ", "_"), func(t *testing.T) {
			l := tabsEngine()
			md := "::::: tabs\n@tab Outer\n" + indent + ":::: tabs\n" + indent + "@tab Middle\n" + indent + indent + "::: tabs\n" + indent + indent + "@tab Inner\n" + indent + indent + "body\n" + indent + indent + ":::\n" + indent + "@tab Middle sibling\n" + indent + "body\n" + indent + "::::\n@tab Outer sibling\nbody\n:::::\n"
			for round := 0; round < 3; round++ {
				tabs, items := tabNodes(parse.Parse("", []byte(md), l.ParseOptions).Root)
				if len(tabs) != 3 || len(items) != 5 {
					t.Fatalf("round %d tabs=%d items=%d\n%s", round, len(tabs), len(items), md)
				}
				if tabs[1].Parent != items[0] || tabs[2].Parent != items[1] || items[3].Parent != tabs[1] || items[4].Parent != tabs[0] {
					t.Fatalf("round %d invalid nesting\n%s", round, md)
				}
				md = l.FormatStr("", md)
			}
		})
	}
	for _, md := range []string{
		"::: tabs\n@tab A\n::: tabs\nbody\n@tab B\n:::\n",
		"::: tabs\n@tab A\n:::: tabs\nbody\n@tab B\n:::\n",
		":::: tabs\n@tab A\n:::\nbody\n@tab B\n::::\n",
		"::: tabs\n@tab A\n::::\nbody\n@tab B\n:::\n",
	} {
		l := tabsEngine()
		tabs, items := tabNodes(parse.Parse("", []byte(md), l.ParseOptions).Root)
		if len(tabs) != 1 || len(items) != 2 || items[1].TabItemTitle != "B" {
			t.Fatalf("mismatched fence changed nesting: tabs=%d items=%d\n%s", len(tabs), len(items), md)
		}
	}
}

func TestTabsCanonicalFenceLengths(t *testing.T) {
	l := tabsEngine()
	for _, test := range []struct {
		markdown string
		fences   []string
	}{
		{":::::::: tabs\n@tab A\nbody\n::::::::\n", []string{"::: tabs", ":::"}},
		{":::::::: tabs\n@tab A\n:::::: tabs\n@tab B\nbody\n::::::\n::::::::\n", []string{":::: tabs", "::: tabs", ":::", "::::"}},
		{":::::::: tabs\n@tab A\n:::::: tabs\n@tab B\n:::: tabs\n@tab C\nbody\n::::\n::::::\n::::::::\n", []string{"::::: tabs", ":::: tabs", "::: tabs", ":::", "::::", ":::::"}},
	} {
		formatted := l.FormatStr("", test.markdown)
		var fences []string
		for _, line := range strings.Split(formatted, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, ":::") {
				fences = append(fences, line)
			}
		}
		if strings.Join(fences, "|") != strings.Join(test.fences, "|") {
			t.Fatalf("fences=%v want=%v\n%s", fences, test.fences, formatted)
		}
	}
}

func TestTabsDeepNesting(t *testing.T) {
	l := tabsEngine()
	const depth = 12
	var markdown strings.Builder
	for level := 0; level < depth; level++ {
		indent := strings.Repeat(" ", 2*level)
		markdown.WriteString(indent + strings.Repeat(":", depth-level+2) + " tabs\n")
		markdown.WriteString(indent + "@tab Level " + string(rune('A'+level)) + "\n")
	}
	markdown.WriteString(strings.Repeat(" ", 2*(depth-1)) + "deep body\n")
	for level := depth - 1; level >= 0; level-- {
		markdown.WriteString(strings.Repeat(" ", 2*level) + strings.Repeat(":", depth-level+2) + "\n")
	}
	md := markdown.String()
	for round := 0; round < 3; round++ {
		tabs, items := tabNodes(parse.Parse("", []byte(md), l.ParseOptions).Root)
		if len(tabs) != depth || len(items) != depth {
			t.Fatalf("round %d tabs=%d items=%d\n%s", round, len(tabs), len(items), md)
		}
		for level := 1; level < depth; level++ {
			if tabs[level].Parent != items[level-1] {
				t.Fatalf("round %d wrong parent at level %d\n%s", round, level, md)
			}
		}
		if !strings.Contains(items[depth-1].Content(), "deep body") {
			t.Fatalf("round %d deepest body lost\n%s", round, md)
		}
		md = l.FormatStr("", md)
	}
}

func TestTabsLiteralMarkers(t *testing.T) {
	l := tabsEngine()
	for name, body := range map[string]string{
		"backtick":      "```\n@tab Literal\n:::\n```\n",
		"tilde":         "~~~\n@tab Literal\n:::\n~~~\n",
		"indented code": "    @tab Literal\n    :::\n",
		"blockquote":    "> @tab Literal\n> :::\n",
		"list":          "- @tab Literal\n  :::\n",
		"escaped":       "\\@tab Literal\n\\@tab:active Literal\n",
		"escaped fence": "\\::: tabs\n\\:::\nLiteral\n",
		"soft break":    "Literal first line\n\\@tab Literal\n\\@tab:active Literal\n",
		"similar token": "@table Literal\n@tab:activeX Literal\n@tab:other Literal\n",
	} {
		t.Run(name, func(t *testing.T) {
			md := "::: tabs\n@tab A\n\n" + body + "\n@tab B\nbody\n:::\n"
			for round := 0; round < 3; round++ {
				tabs, items := tabNodes(parse.Parse("", []byte(md), l.ParseOptions).Root)
				if len(tabs) != 1 || len(items) != 2 || items[1].TabItemTitle != "B" {
					t.Fatalf("round %d literal changed structure: tabs=%d items=%d\n%s", round, len(tabs), len(items), md)
				}
				if !strings.Contains(items[0].Content(), "Literal") {
					t.Fatalf("round %d literal lost\n%s", round, md)
				}
				md = l.FormatStr("", md)
			}
		})
	}
}

func TestTabsLiteralMarkersSpin(t *testing.T) {
	l := tabsEngine()
	for _, body := range []string{
		"\\@tab Literal\n\\@tab:active Literal\n",
		"Literal first line\n\\@tab Literal\n\\@tab:active Literal\n",
		"\\::: tabs\n\\:::\nLiteral\n",
		"Literal first line\n\\::: tabs\n\\:::\n",
	} {
		md := "::: tabs\n@tab A\n\n" + body + "\n@tab B\nbody\n:::\n"
		dom := l.Md2BlockDOM(md, false)
		var expected string
		for round := 0; round < 4; round++ {
			tabs, items := tabNodes(l.BlockDOM2Tree(dom).Root)
			if len(tabs) != 1 || len(items) != 2 || items[1].TabItemTitle != "B" {
				t.Fatalf("round %d tabs=%d items=%d\n%s", round, len(tabs), len(items), dom)
			}
			if round == 0 {
				expected = items[0].Content()
			} else if items[0].Content() != expected {
				t.Fatalf("round %d literal changed: %q want %q\n%s", round, items[0].Content(), expected, dom)
			}
			dom = l.SpinBlockDOM(dom)
		}
	}
}

func TestTabsNestedSlashPreset(t *testing.T) {
	l := tabsEngine()
	markdown := "::: tabs\n@tab\n\n‸\n\n@tab\n\n:::\n"
	paragraph := `<div data-node-id="20260905134000-para001" data-type="NodeParagraph" class="p"><div contenteditable="true">` + markdown + `</div><div class="protyle-attr" contenteditable="false"></div></div>`
	innerDOM := l.SpinBlockDOM(paragraph)
	innerTabs, innerItems := tabNodes(l.BlockDOM2Tree(innerDOM).Root)
	if len(innerTabs) != 1 || len(innerItems) != 2 || !strings.Contains(innerDOM, "<wbr>") {
		t.Fatalf("slash fragment tabs=%d items=%d\n%s", len(innerTabs), len(innerItems), innerDOM)
	}
	dom := `<div data-node-id="20260905120000-tabs001" data-type="NodeTabs" class="tabs"><div data-node-id="20260905120000-item001" data-type="NodeTabItem" class="tab-item"><div class="tab-item-info callout-info"><span class="tab-item-title callout-title">Outer</span></div><div class="tab-item-content">` + innerDOM + `</div></div></div>`
	for round := 0; round < 3; round++ {
		dom = l.SpinBlockDOM(dom)
		tabs, items := tabNodes(l.BlockDOM2Tree(dom).Root)
		if len(tabs) != 2 || len(items) != 3 || tabs[1].Parent != items[0] || !strings.Contains(dom, "<wbr>") {
			t.Fatalf("round %d tabs=%d items=%d\n%s", round, len(tabs), len(items), dom)
		}
	}
}

func TestTabsEmptyAndLeadingBody(t *testing.T) {
	l := tabsEngine()
	for _, test := range []struct {
		markdown string
		items    int
		leading  string
	}{
		{"::: tabs\n:::\n", 1, ""},
		{"::: tabs", 1, ""},
		{"::: tabs\n@tab", 1, ""},
		{"::: tabs\n@tab\n@tab\n:::\n", 2, ""},
		{"::: tabs\nLeading body\n@tab Titled\nbody\n:::\n", 2, "Leading body"},
		{"::: tabs\nLeading body", 1, "Leading body"},
	} {
		md := test.markdown
		for round := 0; round < 3; round++ {
			tabs, items := tabNodes(parse.Parse("", []byte(md), l.ParseOptions).Root)
			if len(tabs) != 1 || len(items) != test.items {
				t.Fatalf("round %d tabs=%d items=%d want=%d\n%s", round, len(tabs), len(items), test.items, md)
			}
			if test.leading != "" && (items[0].TabItemTitle != "" || !strings.Contains(items[0].Text(), test.leading)) {
				t.Fatalf("round %d leading content lost\n%s", round, md)
			}
			for _, item := range items {
				if item.FirstChild == nil {
					t.Fatalf("round %d item has no content block\n%s", round, md)
				}
			}
			md = l.FormatStr("", md)
		}
	}
}

func TestTabsActiveSelection(t *testing.T) {
	for _, ial := range []bool{false, true} {
		for _, blockOnly := range []bool{false, true} {
			for _, test := range []struct {
				name     string
				markers  [3]string
				groupIAL string
				selected int
			}{
				{"default", [3]string{"@tab", "@tab", "@tab"}, "", 0},
				{"second", [3]string{"@tab", "@tab:active", "@tab"}, "", 1},
				{"first marker wins", [3]string{"@tab", "@tab:active", "@tab:active"}, "", 1},
				{"marker overrides IAL", [3]string{"@tab", "@tab:active", "@tab"}, ` tabs-active-id="20260905120000-item003"`, 1},
				{"valid IAL", [3]string{"@tab", "@tab", "@tab"}, ` tabs-active-id="20260905120000-item003"`, 2},
				{"invalid IAL", [3]string{"@tab", "@tab", "@tab"}, ` tabs-active-id="20260905120000-missing"`, 0},
			} {
				name := test.name
				if ial {
					name += "/IAL"
				}
				if blockOnly {
					name += "/Block"
				}
				t.Run(name, func(t *testing.T) {
					l := tabsEngine()
					l.SetKramdownBlockIAL(ial)
					md := "::: tabs\n"
					for index, marker := range test.markers {
						md += marker + " Title\n"
						if ial {
							md += "{: id=\"20260905120000-item00" + string(rune('1'+index)) + "\"}\n"
						}
						md += "\nbody\n"
					}
					md += ":::\n"
					if ial {
						md += "{: id=\"20260905120000-tabs001\"" + test.groupIAL + "}\n"
					}
					parser := parse.Parse
					if blockOnly {
						parser = parse.Block
					}
					tabs, items := tabNodes(parser("", []byte(md), l.ParseOptions).Root)
					if len(tabs) != 1 || len(items) != 3 {
						t.Fatalf("tabs=%d items=%d\n%s", len(tabs), len(items), md)
					}
					selected := test.selected
					if !ial && test.name == "valid IAL" {
						selected = 0
					}
					if items[selected].ID == "" || tabs[0].IALAttr("tabs-active-id") != items[selected].ID {
						t.Fatalf("active=%q expected item %d ID=%q", tabs[0].IALAttr("tabs-active-id"), selected, items[selected].ID)
					}
					formatted := l.FormatStr("", md)
					if strings.Count(formatted, "@tab:active") != 1 {
						t.Fatalf("expected exactly one active marker\n%s", formatted)
					}
					for round := 0; round < 3; round++ {
						tabs, items = tabNodes(parser("", []byte(formatted), l.ParseOptions).Root)
						if len(tabs) != 1 || len(items) != 3 || items[selected].ID == "" || tabs[0].IALAttr("tabs-active-id") != items[selected].ID {
							t.Fatalf("round %d selection changed\n%s", round, formatted)
						}
						formatted = l.FormatStr("", formatted)
					}
				})
			}
		}
	}
}

func TestTabsNestedActiveSelection(t *testing.T) {
	l := tabsEngine()
	md := ":::: tabs\n@tab Outer first\n::: tabs\n@tab Inner first\nbody\n@tab:active Inner second\nbody\n:::\n@tab:active Outer second\nbody\n::::\n"
	for round := 0; round < 3; round++ {
		tabs, items := tabNodes(parse.Parse("", []byte(md), l.ParseOptions).Root)
		if len(tabs) != 2 || len(items) != 4 {
			t.Fatalf("round %d tabs=%d items=%d\n%s", round, len(tabs), len(items), md)
		}
		if tabs[0].IALAttr("tabs-active-id") != items[3].ID || tabs[1].IALAttr("tabs-active-id") != items[2].ID {
			t.Fatalf("round %d outer active=%q inner active=%q\n%s", round, tabs[0].IALAttr("tabs-active-id"), tabs[1].IALAttr("tabs-active-id"), md)
		}
		md = l.FormatStr("", md)
	}
}

func TestTabsItemAttributePlacement(t *testing.T) {
	l := tabsEngine()
	md := "::: tabs\n@tab A\n{: id=\"20260905120000-item001\" custom-item=\"kept\"}\n\nbody\n{: id=\"20260905120000-para001\" custom-body=\"kept\"}\n@tab B\n{: id=\"20260905120000-item002\"}\n\n{: id=\"20260905120000-empty01\"}\n:::\n{: id=\"20260905120000-tabs001\" tabs-position=\"left\"}\n"
	for round := 0; round < 3; round++ {
		tabs, items := tabNodes(parse.Parse("", []byte(md), l.ParseOptions).Root)
		if len(tabs) != 1 || len(items) != 2 {
			t.Fatalf("round %d tabs=%d items=%d\n%s", round, len(tabs), len(items), md)
		}
		if tabs[0].ID != "20260905120000-tabs001" || tabs[0].IALAttr("tabs-position") != "left" || items[0].ID != "20260905120000-item001" || items[0].IALAttr("custom-item") != "kept" || items[1].ID != "20260905120000-item002" {
			t.Fatalf("round %d container attributes changed\n%s", round, md)
		}
		if items[0].FirstChild.ID != "20260905120000-para001" || items[0].FirstChild.IALAttr("custom-body") != "kept" || items[1].FirstChild.ID != "20260905120000-empty01" {
			t.Fatalf("round %d body attributes changed\n%s", round, md)
		}
		md = l.FormatStr("", md)
	}
}
