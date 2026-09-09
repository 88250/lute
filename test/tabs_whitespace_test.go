package test

import (
	"html"
	"strings"
	"testing"

	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func TestTabsTitleWhitespace(t *testing.T) {
	l := tabsEngine()
	for _, title := range []string{" leading", "trailing ", "two  spaces", "   ", " <span data-type=\"strong\">bold</span> "} {
		dom := `<div data-node-id="20260909000000-tabs001" data-type="NodeTabs" class="tabs"><div data-node-id="20260909000000-item001" data-type="NodeTabItem" class="tab-item"><div class="tab-item-info"><span class="tab-item-title">` + title + `</span></div><div class="tab-item-content"><div data-node-id="20260909000000-para001" data-type="NodeParagraph" class="p"><div contenteditable="true">Body</div></div></div></div></div>`
		for round := 0; round < 3; round++ {
			tree := l.BlockDOM2Tree(dom)
			_, items := tabNodes(tree.Root)
			if len(items) != 1 {
				t.Fatal(dom)
			}
			inline := parse.Inline("", []byte(items[0].TabItemTitle), l.ParseOptions)
			expected := strings.ReplaceAll(strings.ReplaceAll(title, `<span data-type="strong">`, ""), "</span>", "")
			got := html.UnescapeString(string(render.NewHtmlRenderer(inline, l.RenderOptions, l.ParseOptions).Render()))
			got = strings.TrimSuffix(strings.TrimPrefix(got, "<p>"), "</p>\n")
			got = strings.ReplaceAll(strings.ReplaceAll(got, "<strong>", ""), "</strong>", "")
			got = strings.ReplaceAll(strings.ReplaceAll(got, `<span data-type="strong">`, ""), "</span>", "")
			if got != expected {
				t.Fatalf("round %d: title %q became %q (%q)", round, title, got, items[0].TabItemTitle)
			}
			dom = string(render.NewProtyleRenderer(tree, l.RenderOptions, l.ParseOptions).Render())
			dom = l.SpinBlockDOM(dom)
		}
	}
}
