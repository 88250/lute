package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestVditorTODOComment(t *testing.T) {
	for _, mode := range []string{"wysiwyg", "ir"} {
		for _, comment := range []string{"TODO: remember to report this issue", "TODO: check \\(nested\\) & text", "TODO: <img src=x onerror=alert(1)>"} {
			t.Run(mode+"/"+comment, func(t *testing.T) {
				engine := lute.New()
				toDOM, toMarkdown, spin := engine.Md2VditorDOM, engine.VditorDOM2Md, engine.SpinVditorDOM
				if mode == "ir" {
					toDOM, toMarkdown, spin = engine.Md2VditorIRDOM, engine.VditorIRDOM2Md, engine.SpinVditorIRDOM
				}
				definition := "[//]: # (" + comment + ")"
				dom := toDOM("before\n\n" + definition + "\n\nafter")
				for round := 0; round < 3; round++ {
					got := toMarkdown(dom)
					if !strings.Contains(got, definition) || !strings.Contains(got, "before") || !strings.Contains(got, "after") {
						t.Fatalf("lost TODO syntax or adjacent content: %q", got)
					}
					if strings.Contains(dom, "<img") {
						t.Fatalf("comment became HTML: %s", dom)
					}
					dom = spin(dom)
				}
			})
		}
	}
}
