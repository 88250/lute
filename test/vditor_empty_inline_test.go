package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestVditorEmptyInlineCode(t *testing.T) {
	for _, kind := range []string{"math-inline", "html-inline", "code-inline"} {
		for _, content := range []string{"<code></code>", "\u200b<code></code>", "<code>x</code>"} {
			t.Run(kind+"/"+content, func(t *testing.T) {
				engine := lute.New()
				dom := `<p data-block="0">before <span data-type="` + kind + `">` + content + `</span> after</p>`
				for _, convert := range []func(string) string{engine.VditorDOM2Md, engine.SpinVditorDOM, engine.VditorDOM2HTML} {
					got := convert(dom)
					if !strings.Contains(got, "before") || !strings.Contains(got, "after") {
						t.Fatalf("lost adjacent content: %q", got)
					}
					if strings.Contains(content, ">x<") && !strings.Contains(got, "x") {
						t.Fatalf("lost inline content: %q", got)
					}
				}
			})
		}
	}
}
