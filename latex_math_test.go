package lute_test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
)

func TestLatexMathMarkdown(t *testing.T) {
	for _, tc := range []struct {
		name, source  string
		inline, block int
	}{
		{"inline", `before \(\hat VY\) after \(z\)`, 2, 0},
		{"table", "| Formula | Value |\n| --- | --- |\n| \\(\\hat VY\\) | 0.60 |", 1, 0},
		{"display same line", `\[x+y\]`, 0, 1},
		{"display multiline", "before\n\\[\n\\begin{aligned}\nx &= 1 \\\\\ny &= 2\n\\end{aligned}\n\\]\nafter", 0, 1},
		{"display blank line", "\\[\nx\n\ny\n\\]", 0, 1},
		{"display row spacing", "\\[\n\\begin{aligned}\nx &= 1 \\\\[1em]\ny &= 2\n\\end{aligned}\n\\]", 0, 1},
		{"display quote", "> \\[\n> x\n> \\]", 0, 1},
		{"display list", "- \\[\n  x\n  \\]", 0, 1},
		{"unclosed", "\\[\nx\n\n# Keep heading\n\\(y", 0, 0},
		{"empty", `\( \) \[\]`, 0, 0},
		{"escaped", `\\(x\\)`, 0, 0},
		{"inline newline", "\\(x\\\ny\\)", 0, 0},
		{"code", "`\\(x\\)`\n\n```tex\n\\[x\\]\n```", 0, 0},
		{"dollar content", `$\text{\(x\)}$`, 1, 0},
		{"escaped close", `\(a\\)b\)`, 1, 0},
		{"braced close", `\(\text{a\)b}+c\)`, 1, 0},
		{"quote boundary", "> \\[\n> x\n\n\\]", 0, 0},
		{"list boundary", "- \\[\n  x\n\n\\]", 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			engine := lute.New()
			engine.SetProtyleWYSIWYG(true)
			engine.SetLatexMath(true)
			dom := engine.Md2BlockDOM(tc.source, false)
			if got := strings.Count(dom, `data-type="inline-math"`); got != tc.inline {
				t.Fatalf("inline count=%d want=%d: %s", got, tc.inline, dom)
			}
			if got := strings.Count(dom, `data-type="NodeMathBlock"`); got != tc.block {
				t.Fatalf("block count=%d want=%d: %s", got, tc.block, dom)
			}
			if tc.name == "unclosed" && !strings.Contains(dom, `data-type="NodeHeading"`) {
				t.Fatalf("unclosed marker swallowed a heading: %s", dom)
			}
			engine.SetLatexMath(false)
			spun := engine.SpinBlockDOM(dom)
			if strings.Count(spun, `data-type="inline-math"`) != tc.inline || strings.Count(spun, `data-type="NodeMathBlock"`) != tc.block {
				t.Fatalf("canonical editing round trip changed formulas: %s", spun)
			}
		})
	}
}

func TestLatexMathHTML(t *testing.T) {
	for _, tc := range []struct {
		name, source  string
		inline, block int
		contains      string
	}{
		{"inline", `<p>before \(x+y\) after \(z\)</p>`, 2, 0, "after"},
		{"table", `<table><tr><th>Formula</th></tr><tr><td>\(\hat VY\)</td></tr></table>`, 1, 0, `data-type="NodeTable"`},
		{"display", `<p>\[x+y\]</p>`, 0, 1, `data-content="x+y"`},
		{"breaks", `<p>\[<br>\begin{aligned}<br>x&amp;=1\\<br>y&amp;=2<br>\end{aligned}<br>\]</p>`, 0, 1, `x&amp;=1`},
		{"split inline", `<p><strong>before \(x</strong><em>+y\) after</em></p>`, 1, 0, `data-content="x+y"`},
		{"split display", `<p><strong>\[x</strong><em>+y\]</em></p>`, 0, 1, `data-content="x+y"`},
		{"styles", `<p><span style="color:#ff0000">before \(x\) after</span></p>`, 1, 0, `color: #ff0000`},
		{"code", `<pre>\[x\]</pre><p><code>\(y\)</code></p>`, 0, 0, "y"},
		{"mixed display", `<p><code>code</code>\[x\]</p>`, 0, 0, "code"},
		{"dollar content", `<p>$\text{\(x\)}$</p>`, 0, 0, "text"},
		{"existing math", `<p><span data-type="inline-math" data-content="\text{\(x\)}"></span></p>`, 1, 0, `data-content="\text{\(x\)}"`},
		{"unclosed", `<p>\(x</p><p>\)</p>`, 0, 0, "x"},
		{"escaped", `<p>\\(x\\)</p>`, 0, 0, "x"},
		{"link", `<p><a href="https://example.com/\(x\)">label</a></p>`, 0, 0, `https://example.com/`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			engine := lute.New()
			engine.SetProtyleWYSIWYG(true)
			engine.SetHTMLTag2TextMark(true)
			engine.SetTextMark(true)
			engine.SetKramdownIAL(true)
			engine.SetLatexMath(true)
			dom := engine.HTML2BlockDOM(tc.source)
			ast.Walk(engine.HTML2Tree(tc.source).Root, func(n *ast.Node, entering bool) ast.WalkStatus {
				if entering && n.Type == ast.NodeMathBlock && (n.Parent == nil || !n.Parent.IsBlock() || n.Parent.Type == ast.NodeParagraph || n.Parent.Type == ast.NodeHeading) {
					t.Fatalf("display formula has an invalid parent: %+v", n.Parent)
				}
				return ast.WalkContinue
			})
			if strings.Count(dom, `data-type="inline-math"`) != tc.inline || strings.Count(dom, `data-type="NodeMathBlock"`) != tc.block || !strings.Contains(dom, tc.contains) {
				t.Fatalf("unexpected conversion: %s", dom)
			}
		})
	}
}

func TestLatexMathOptIn(t *testing.T) {
	engine := lute.New()
	for _, dom := range []string{engine.Md2BlockDOM(`\(x\)`, false), engine.HTML2BlockDOM(`<p>\[x\]</p>`)} {
		if strings.Contains(dom, `data-type="inline-math"`) || strings.Contains(dom, `data-type="NodeMathBlock"`) {
			t.Fatalf("alternate syntax must be opt-in: %s", dom)
		}
	}
	engine.SetLatexMath(true)
	engine.SetInlineMath(false)
	if dom := engine.Md2BlockDOM(`\(x\)`, false); strings.Contains(dom, `data-type="inline-math"`) {
		t.Fatalf("inline math setting ignored: %s", dom)
	}
	if markdown := engine.BlockDOM2StdMd(engine.Md2BlockDOM(`\[x+y\]`, false)); !strings.Contains(markdown, "$$\nx+y\n$$") {
		t.Fatalf("export must use canonical dollar delimiters: %s", markdown)
	}
}
