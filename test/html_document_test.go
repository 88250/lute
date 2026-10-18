package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestHTMLDocumentConversion(t *testing.T) {
	for _, prefix := range []string{"", "<!DOCTYPE html>", "<!-- comment -->", "<!-- comment --><!DOCTYPE html>", "<!doctype HTML>\n<!-- comment -->"} {
		t.Run(prefix, func(t *testing.T) {
			engine := lute.New()
			source := prefix + `<html><head><title>Hidden title</title></head><body><h1>Page heading</h1><p>Page content</p></body></html>`
			markdown, err := engine.HTML2Markdown(source)
			if err != nil || !strings.Contains(markdown, "# Page heading") || !strings.Contains(markdown, "Page content") || strings.Contains(markdown, "Hidden title") {
				t.Fatalf("unexpected Markdown: %q, %v", markdown, err)
			}
			text := engine.HTML2Text(source)
			if !strings.Contains(text, "Page heading") || !strings.Contains(text, "Page content") || strings.Contains(text, "Hidden title") {
				t.Fatalf("unexpected text: %q", text)
			}
		})
	}
}

func TestHTMLDocumentEmptyAndFragment(t *testing.T) {
	for _, source := range []string{"", "<!DOCTYPE html>", "<!-- comment -->", "<p>Fragment</p>"} {
		engine := lute.New()
		markdown, err := engine.HTML2Markdown(source)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(source, "Fragment") && !strings.Contains(markdown, "Fragment") {
			t.Fatalf("fragment lost: %q", markdown)
		}
		engine.HTML2Text(source)
	}
}
