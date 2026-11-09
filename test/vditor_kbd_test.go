package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestVditorWYSIWYGKbd(t *testing.T) {
	for _, markdown := range []string{
		"<kbd>Ctrl</kbd> + <kbd>C</kbd>",
		"before <kbd>**bold** and `code`</kbd> after",
		"<kbd>outer <kbd>inner</kbd> end</kbd>",
		"<kbd></kbd>",
	} {
		t.Run(markdown, func(t *testing.T) {
			engine := lute.New()
			want := engine.Md2HTML(markdown)
			dom := engine.Md2VditorDOM(markdown)
			for round := 0; round < 3; round++ {
				dom = strings.ReplaceAll(dom, `class="vditor-wysiwyg__kbd-marker"`,
					`class="vditor-wysiwyg__kbd-marker vditor-wysiwyg__kbd-marker--expand"`)
				if !strings.Contains(dom, "<kbd>") || !strings.Contains(dom, "vditor-wysiwyg__kbd-marker") {
					t.Fatalf("missing keyboard preview: %s", dom)
				}
				if got := engine.Md2HTML(engine.VditorDOM2Md(dom)); got != want {
					t.Fatalf("round %d changed keyboard content\nwant: %s\ngot: %s", round, want, got)
				}
				dom = engine.SpinVditorDOM(dom)
			}
		})
	}
}

func TestVditorUnmatchedKbd(t *testing.T) {
	engine := lute.New()
	for _, markdown := range []string{"<kbd>open", "close</kbd>"} {
		dom := engine.Md2VditorDOM(markdown)
		if strings.Contains(dom, "vditor-wysiwyg__kbd-marker") {
			t.Fatalf("unmatched tag became a preview: %s", dom)
		}
		if got := engine.VditorDOM2Md(dom); strings.TrimSpace(got) != markdown {
			t.Fatalf("unmatched tag changed: %s", got)
		}
	}
}

func TestVditorPastedKbd(t *testing.T) {
	engine := lute.New()
	for _, content := range []string{"<kbd>Ctrl</kbd>", "<kbd><strong>bold</strong></kbd>"} {
		dom := engine.SpinVditorDOM(`<p data-block="0">` + content + `</p>`)
		if !strings.Contains(dom, "vditor-wysiwyg__kbd-marker") {
			t.Fatalf("pasted keyboard tag lost its preview: %s", dom)
		}
		if got := engine.Md2HTML(engine.VditorDOM2Md(dom)); got != "<p>"+content+"</p>\n" {
			t.Fatalf("pasted keyboard content changed: %s", got)
		}
	}
}
