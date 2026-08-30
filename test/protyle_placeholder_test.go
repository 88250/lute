// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
)

func TestProtyleInlinePlaceholderCodepoints(t *testing.T) {
	assertSingleCodepoint(t, "Zwsp", editor.Zwsp, '\u200b')
	assertSingleCodepoint(t, "WordJoiner", editor.WordJoiner, '\u2060')
	if editor.Zwsp == editor.WordJoiner {
		t.Fatal("the external and internal placeholders must use different codepoints")
	}
}

func TestProtyleInlinePlaceholderRendering(t *testing.T) {
	luteEngine := newProtylePlaceholderLute()
	tests := []struct {
		name     string
		markdown string
		dataType string
		content  string
	}{
		{"code", "a`code`b", "code", "code"},
		{"combined-code", "a**`code`**b", "strong code", "code"},
		{"kbd", "a<kbd>key</kbd>b", "kbd", "key"},
		{"tag", "a#tag#b", "tag", "tag"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dom := luteEngine.Md2BlockDOM(test.markdown, false)
			open := "a" + editor.Zwsp + `<span data-type="` + test.dataType + `">` + editor.WordJoiner + test.content
			close := test.content + "</span>" + editor.Zwsp + "b"
			if !strings.Contains(dom, open) || !strings.Contains(dom, close) {
				t.Fatalf("inline placeholder invariant failed\nexpected open %q and close %q\ngot %q", open, close, dom)
			}
		})
	}
}

func TestBlockDOMInlinePlaceholderCompatibility(t *testing.T) {
	luteEngine := newProtylePlaceholderLute()
	luteWithoutTextMarks := newProtylePlaceholderLute()
	luteWithoutTextMarks.SetTextMark(false)
	markers := []struct {
		name  string
		value string
	}{
		{"zwsp", editor.Zwsp},
		{"zero-width-no-break-space", "\ufeff"},
		{"word-joiner", editor.WordJoiner},
	}
	containers := []struct {
		name  string
		open  string
		close string
	}{
		{"code", `<span data-type="code">`, "</span>"},
		{"combined-code", `<span data-type="strong code">`, "</span>"},
		{"kbd-span", `<span data-type="kbd">`, "</span>"},
		{"kbd", "<kbd>", "</kbd>"},
		{"tag", `<span data-type="tag">`, "</span>"},
		{"code-element", "<code>", "</code>"},
	}

	for _, marker := range markers {
		for _, container := range containers {
			t.Run(marker.name+"/"+container.name, func(t *testing.T) {
				inlines := "outside" + editor.WordJoiner + container.open + marker.value + "value" + editor.WordJoiner + "tail" + container.close
				tree := luteEngine.BlockDOM2Tree(protylePlaceholderBlockDOM(inlines))
				var content string
				ast.Walk(tree.Root, func(node *ast.Node, entering bool) ast.WalkStatus {
					if entering && ast.NodeTextMark == node.Type && "" == content {
						content = node.TextMarkTextContent
					}
					return ast.WalkContinue
				})
				if expected := "value" + editor.WordJoiner + "tail"; expected != content {
					t.Fatalf("expected text mark content %q, got %q", expected, content)
				}
				if expected := "outside" + editor.WordJoiner + content; expected != tree.Root.Text() {
					t.Fatalf("expected AST text %q, got %q", expected, tree.Root.Text())
				}
				withoutTextMarks := luteWithoutTextMarks.BlockDOM2Md(protylePlaceholderBlockDOM(inlines))
				if !strings.Contains(withoutTextMarks, "value"+editor.WordJoiner+"tail") ||
					2 != strings.Count(withoutTextMarks, editor.WordJoiner) || strings.Contains(withoutTextMarks, "\ufeff") ||
					strings.Contains(withoutTextMarks, editor.Zwsp) {
					t.Fatalf("placeholder normalization without TextMark failed: %q", withoutTextMarks)
				}
			})
		}
	}
}

func TestBlockDOMInlinePlaceholderNormalizationIsScoped(t *testing.T) {
	luteEngine := newProtylePlaceholderLute()
	inlines := "outside" + editor.WordJoiner + `<span data-type="strong">` + editor.WordJoiner + "strong</span>" +
		`<kbd><span>foo</span>` + editor.WordJoiner + `bar</kbd>` +
		`<kbd>` + editor.WordJoiner + `outer<code>` + editor.WordJoiner + `inner</code></kbd>` +
		`<span data-type="code">` + editor.WordJoiner + editor.WordJoiner + `double</span>` +
		`<span data-type="code-block kbdish hashtag">` + editor.WordJoiner + `similar</span>` +
		`<span data-type="code">&NoBreak;entity</span>`
	tree := luteEngine.BlockDOM2Tree(protylePlaceholderBlockDOM(inlines))
	expected := "outside" + editor.WordJoiner + editor.WordJoiner + "strongfoo" + editor.WordJoiner + "barouterinner" +
		editor.WordJoiner + "double" + editor.WordJoiner + "similarentity"
	if expected != tree.Root.Text() {
		t.Fatalf("expected scoped normalization text %q, got %q", expected, tree.Root.Text())
	}
}

func TestBlockDOMInlinePlaceholderPreservesAttributeWordJoiners(t *testing.T) {
	luteEngine := newProtylePlaceholderLute()
	href := "before" + editor.WordJoiner + "after"
	math := "x" + editor.WordJoiner + "y"
	inlines := `<span data-type="a" data-href="` + href + `">link</span>` +
		`<span data-type="inline-math" data-content="` + math + `"></span>`
	tree := luteEngine.BlockDOM2Tree(protylePlaceholderBlockDOM(inlines))
	var gotHref, gotMath string
	ast.Walk(tree.Root, func(node *ast.Node, entering bool) ast.WalkStatus {
		if entering && ast.NodeTextMark == node.Type {
			if node.IsTextMarkType("a") {
				gotHref = node.TextMarkAHref
			}
			if node.IsTextMarkType("inline-math") {
				gotMath = node.TextMarkInlineMathContent
			}
		}
		return ast.WalkContinue
	})
	if href != gotHref || math != gotMath {
		t.Fatalf("expected word joiners in attributes to be preserved, got href %q and math %q", gotHref, gotMath)
	}
}

func TestSpinBlockDOMInlinePlaceholderCanonicalAndIdempotent(t *testing.T) {
	luteEngine := newProtylePlaceholderLute()
	ast.Testing = true
	defer func() {
		ast.Testing = false
	}()

	inlines := `a<span data-type="code">` + "\ufeff" + `code</span>b` +
		`c<span data-type="strong code">` + editor.Zwsp + `both</span>d` +
		`e<span data-type="kbd">` + editor.WordJoiner + `key</span>f` +
		`g<span data-type="tag">` + editor.Zwsp + `tag</span>h`
	spun := luteEngine.SpinBlockDOM(protylePlaceholderBlockDOM(inlines))
	checks := []string{
		"a" + editor.Zwsp + `<span data-type="code">` + editor.WordJoiner + "code</span>" + editor.Zwsp + "b",
		"c" + editor.Zwsp + `<span data-type="strong code">` + editor.WordJoiner + "both</span>" + editor.Zwsp + "d",
		"e" + editor.Zwsp + `<span data-type="kbd">` + editor.WordJoiner + "key</span>" + editor.Zwsp + "f",
		"g" + editor.Zwsp + `<span data-type="tag">` + editor.WordJoiner + "tag</span>" + editor.Zwsp + "h",
	}
	for _, check := range checks {
		if !strings.Contains(spun, check) {
			t.Fatalf("canonical placeholder sequence %q not found in %q", check, spun)
		}
	}
	if strings.Contains(spun, "\ufeff") {
		t.Fatalf("legacy placeholder leaked into spun DOM: %q", spun)
	}
	if spunAgain := luteEngine.SpinBlockDOM(spun); spun != spunAgain {
		t.Fatalf("SpinBlockDOM is not idempotent\nfirst  %q\nsecond %q", spun, spunAgain)
	}
}

func TestBlockDOMInlinePlaceholderDoesNotLeakToExports(t *testing.T) {
	luteEngine := newProtylePlaceholderLute()
	markers := []struct {
		name  string
		value string
	}{
		{"zwsp", editor.Zwsp},
		{"zero-width-no-break-space", "\ufeff"},
		{"word-joiner", editor.WordJoiner},
	}

	for _, marker := range markers {
		t.Run(marker.name, func(t *testing.T) {
			dom := protylePlaceholderBlockDOM(`before <span data-type="code">` + marker.value + `code</span> after`)
			tree := luteEngine.BlockDOM2Tree(dom)
			exports := map[string]string{
				"ast":       tree.Root.Text(),
				"markdown":  luteEngine.BlockDOM2Md(dom),
				"standard":  luteEngine.BlockDOM2StdMd(dom),
				"html":      luteEngine.BlockDOM2HTML(dom),
				"rich-html": luteEngine.BlockDOM2RichHTML(dom),
			}
			for name, output := range exports {
				if strings.Contains(output, marker.value) {
					t.Fatalf("%s export leaked %s placeholder: %q", name, marker.name, output)
				}
			}
		})
	}
}

func TestBlockDOMInlinePlaceholderPreservesRealWordJoiners(t *testing.T) {
	luteEngine := newProtylePlaceholderLute()
	markers := []struct {
		name  string
		value string
	}{
		{"zwsp", editor.Zwsp},
		{"zero-width-no-break-space", "\ufeff"},
		{"word-joiner", editor.WordJoiner},
	}

	for _, marker := range markers {
		t.Run(marker.name, func(t *testing.T) {
			inlines := "before" + editor.WordJoiner + `<span data-type="code">` + marker.value + editor.WordJoiner + "co" + editor.WordJoiner + "de</span>after"
			dom := protylePlaceholderBlockDOM(inlines)
			exports := map[string]string{
				"ast":       luteEngine.BlockDOM2Tree(dom).Root.Text(),
				"markdown":  luteEngine.BlockDOM2Md(dom),
				"standard":  luteEngine.BlockDOM2StdMd(dom),
				"html":      luteEngine.BlockDOM2HTML(dom),
				"rich-html": luteEngine.BlockDOM2RichHTML(dom),
			}
			for name, output := range exports {
				if count := strings.Count(output, editor.WordJoiner); 3 != count {
					t.Fatalf("expected %s to preserve 3 real word joiners, got %d in %q", name, count, output)
				}
			}
		})
	}
}

func newProtylePlaceholderLute() *lute.Lute {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetKramdownIAL(true)
	luteEngine.SetTag(true)
	luteEngine.SetTextMark(true)
	luteEngine.SetHTMLTag2TextMark(true)
	luteEngine.SetAutoSpace(false)
	luteEngine.SetSpin(true)
	return luteEngine
}

func protylePlaceholderBlockDOM(inlines string) string {
	return `<div data-node-id="20260830120000-abcdefg" data-node-index="1" data-type="NodeParagraph" class="p" updated="20260830120000"><div contenteditable="true" spellcheck="false">` +
		inlines + `</div><div class="protyle-attr" contenteditable="false">` + editor.Zwsp + `</div></div>`
}

func assertSingleCodepoint(t *testing.T, name, value string, expected rune) {
	t.Helper()
	codepoints := []rune(value)
	if 1 != len(codepoints) || expected != codepoints[0] {
		t.Fatalf("expected %s to be U+%04X, got %U", name, expected, codepoints)
	}
}
