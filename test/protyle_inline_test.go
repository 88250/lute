package test

import (
	"strings"
	"testing"

	"github.com/88250/lute/ast"
)

func TestSpinInlineBlockDOM(t *testing.T) {
	engine := newTableRichLute()
	for _, text := range []string{"# heading", "- list", "1. ordered", "&gt; quote", "---", "```code", "| table |", "# first\n\n- second\n", "  spaced  "} {
		input := `<div data-node-id="20261025000000-abcdefg" data-type="NodeParagraph" custom-test="retained"><div contenteditable="true">` + text + `</div><div class="protyle-attr" contenteditable="false"></div></div>`
		output := engine.SpinInlineBlockDOM(input)
		tree := engine.BlockDOM2Tree(output)
		paragraph := tree.Root.FirstChild
		if paragraph.Type != ast.NodeParagraph || paragraph.ID != "20261025000000-abcdefg" || paragraph.IALAttr("custom-test") != "retained" {
			t.Fatalf("paragraph identity lost for %q: %s", text, output)
		}
		for child := paragraph.Next; child != nil; child = child.Next {
			if child.Type != ast.NodeKramdownBlockIAL {
				t.Fatalf("block syntax promoted for %q: %s", text, output)
			}
		}
		if strings.Contains(text, "second") && !strings.Contains(output, "second") {
			t.Fatalf("soft line content lost: %s", output)
		}
		if strings.Contains(text, "second") && !strings.Contains(output, "# first\n\n- second\n") {
			t.Fatalf("soft breaks were trimmed: %s", output)
		}
		if text == "  spaced  " && !strings.Contains(output, text) {
			t.Fatalf("inline spaces were trimmed: %s", output)
		}
		if again := engine.SpinInlineBlockDOM(output); again != output {
			t.Fatalf("unstable spin for %q:\n%s\n%s", text, output, again)
		}
	}
}

func TestSpinInlineBlockDOMRichAndBlockInput(t *testing.T) {
	engine := newTableRichLute()
	input := `<div data-node-id="20261025000000-abcdefg" data-type="NodeParagraph"><div contenteditable="true"># **strong** <span data-type="block-ref" data-id="20261025000001-abcdefg" data-subtype="d">target</span> $x^2$ ![alt](assets/image.png)<br />end</div><div class="protyle-attr" contenteditable="false"></div></div>`
	output := engine.SpinInlineBlockDOM(input)
	for _, expected := range []string{`data-type="NodeParagraph"`, `data-type="strong"`, `data-type="block-ref"`, `data-type="inline-math"`, `data-type="img"`, `assets/image.png`, "end"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("missing %q: %s", expected, output)
		}
	}
	for _, block := range []string{engine.Md2BlockDOM("# heading\n", false), engine.Md2BlockDOM("one\n\ntwo\n", false)} {
		if output := engine.SpinInlineBlockDOM(block); output != block {
			t.Fatalf("non-inline input altered: %s", output)
		}
	}
}

func TestSpinInlineBlockDOMPreservesLiteralEscapes(t *testing.T) {
	engine := newTableRichLute()
	input := `<div data-node-id="20261025000000-abcdefg" data-type="NodeParagraph"><div contenteditable="true"><span data-type="backslash">*</span><span data-type="backslash">*</span>literal<span data-type="backslash">*</span><span data-type="backslash">*</span> **new**</div><div class="protyle-attr" contenteditable="false"></div></div>`
	output := engine.SpinInlineBlockDOM(input)
	if !strings.Contains(output, "literal") || strings.Count(output, `data-type="backslash"`) != 4 || !strings.Contains(output, `data-type="strong">new`) {
		t.Fatalf("existing literal text changed during append: %s", output)
	}
	for _, mark := range []string{"text", "strong", "em", "s", "mark", "sup", "sub", "code", "strong em"} {
		input = `<div data-node-id="20261025000000-abcdefg" data-type="NodeParagraph"><div contenteditable="true"><span data-type="` + mark + `" style="color: red">**literal** [link](url) $x$ #tag#</span> **new**</div><div class="protyle-attr" contenteditable="false"></div></div>`
		output = engine.SpinInlineBlockDOM(input)
		if !strings.Contains(output, "**literal** [link](url) $x$ #tag#") || !strings.Contains(output, `data-type="strong">new`) || !strings.Contains(output, `style="color: red"`) {
			t.Fatalf("existing %s mark content changed during append: %s", mark, output)
		}
	}
}
