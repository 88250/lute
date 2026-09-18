package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"
	"github.com/88250/lute/util"
)

func TestTextMarkEscapedKramdownRoundTrip(t *testing.T) {
	for _, mark := range []string{"em", "strong", "s", "mark", "sup", "sub", "u", "kbd", "code", "text", "a em"} {
		t.Run(mark, func(t *testing.T) {
			engine := lute.New()
			engine.SetTextMark(true)
			engine.SetProtyleWYSIWYG(true)
			engine.SetKramdownIAL(true)
			engine.SetAutoSpace(false)
			engine.SetSpin(true)
			const content = `<vitae> & &lt; "quote"`
			span := `<span data-type="` + mark + `" data-href="https://example.com">` + html.EscapeHTMLStr(content) + `</span>`
			dom := `<div data-node-id="20260918120000-abcdefg" data-type="NodeParagraph"><div contenteditable="true">` + span + `</div></div>`
			check := func(got string) {
				t.Helper()
				text := util.DomText(util.ParseHTML(got))
				text = strings.NewReplacer(editor.Zwsp, "", editor.WordJoiner, "").Replace(text)
				if text != content {
					t.Fatalf("expected %q, got %q in %s", content, text, got)
				}
			}
			for i := 0; i < 3; i++ {
				check(engine.SpinBlockDOM(dom))
				if text := strings.ReplaceAll(engine.BlockDOM2Content(dom), editor.Zwsp, ""); text != content {
					t.Fatalf("plain content changed: %q", text)
				}
				markdown := engine.BlockDOM2Md(dom)
				dom = engine.Md2BlockDOM(markdown, false)
				check(dom)
				check(engine.InlineMd2BlockDOM(span))
			}
		})
	}
}
