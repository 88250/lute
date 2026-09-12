package test

import (
	"github.com/88250/lute"
	"strings"
	"testing"
)

func TestCalloutBacktickInput(t *testing.T) {
	l := lute.New()
	l.SetTextMark(true)
	l.SetCallout(true)
	l.SetKramdownIAL(true)
	l.SetSpin(true)
	l.SetProtyleWYSIWYG(true)
	for _, tc := range []struct{ title, want string }{
		{"`<wbr>", "`<wbr>"},
		{"Note `<wbr>", "Note `<wbr>"},
		{"Note ``<wbr>", "Note ``<wbr>"},
		{"Note `code`<wbr>", "Note <span data-type=\"code\">code</span><wbr>"},
		{"Note <span data-type=\"code\">code</span>`<wbr>", "Note <span data-type=\"code\">code</span>`<wbr>"},
	} {
		t.Run(tc.title, func(t *testing.T) {
			dom := "<div data-type=\"NodeCallout\" data-subtype=\"NOTE\" class=\"callout\"><div class=\"callout-info\"><span class=\"callout-icon\">✏️</span><span class=\"callout-title\" contenteditable=\"true\">" + tc.title + "</span></div><div class=\"callout-content\"><div data-type=\"NodeParagraph\"><div contenteditable=\"true\">body</div></div></div></div>"
			got := l.SpinBlockDOM(dom)
			got = strings.NewReplacer("\u200b", "", "\u2060", "").Replace(got)
			if !strings.Contains(got, tc.want) {
				t.Fatalf("%s\nMD: %s", got, l.BlockDOM2Md(dom))
			}
		})
	}
}
