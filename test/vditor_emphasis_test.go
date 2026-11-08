package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestVditorMixedEmphasis(t *testing.T) {
	for _, fixture := range []string{
		`1<strong>1<em>2</em></strong><em>1</em>`,
		`1<strong>1<em>2</em></strong><em>1</em>3`,
		`<em>a<strong>b</strong></em><strong>c</strong>`,
		`<strong>a</strong><em>b</em>`,
	} {
		t.Run(fixture, func(t *testing.T) {
			engine := lute.New()
			dom := `<p data-block="0">` + fixture + `</p>`
			want := "<p>" + fixture + "</p>\n"
			for round := 0; round < 3; round++ {
				markdown := engine.VditorDOM2Md(dom)
				if got := engine.Md2HTML(markdown); got != want {
					t.Fatalf("round %d changed emphasis\nmarkdown: %s\nwant: %s\ngot: %s", round, markdown, want, got)
				}
				ir := engine.Md2VditorIRDOM(markdown)
				if got := engine.Md2HTML(engine.VditorIRDOM2Md(ir)); got != want {
					t.Fatalf("IR changed emphasis: %s", got)
				}
				dom = engine.SpinVditorDOM(dom)
				if !strings.Contains(dom, "<strong") || !strings.Contains(dom, "<em") {
					t.Fatalf("formatting is not editable: %s", dom)
				}
			}
		})
	}
}

func TestVditorOrdinaryEmphasisMarkdown(t *testing.T) {
	engine := lute.New()
	for _, markdown := range []string{"**bold**", "*italic*", "**bold *italic* text**"} {
		got := engine.VditorDOM2Md(engine.Md2VditorDOM(markdown))
		if strings.Contains(got, "<strong>") || strings.Contains(got, "<em>") {
			t.Fatalf("ordinary emphasis became HTML: %s", got)
		}
	}
}
