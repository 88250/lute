package lute_test

import (
	"html"
	"regexp"
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
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
		{"display quote in list", "- > \\[\n  > x\n  > \\]", 0, 1},
		{"display list in quote", "> - \\[\n>   x\n>   \\]", 0, 1},
		{"display nested list", "- Outer\n  - \\[\n    x\n    \\]", 0, 1},
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
	inlineMath := regexp.MustCompile(`data-type="(?:[^" ]+ )*inline-math(?: [^" ]+)*"`)
	for _, tc := range []struct {
		name, source  string
		inline, block int
		contains      string
	}{
		{"inline", `<p>before \(x+y\) after \(z\)</p>`, 2, 0, "after"},
		{"table", `<table><tr><th>Formula</th></tr><tr><td>\(\hat VY\)</td></tr></table>`, 1, 0, `data-type="NodeTable"`},
		{"nested table formula", `<table><tr><th>Formula</th></tr><tr><td><p><span><span data-math-source="x+y" data-math-display="false">\(x+y\)</span></span></p></td></tr></table>`, 1, 0, `data-content="x+y"`},
		{"nested paragraph formula", `<p><span><span>\(x+y\)</span></span></p>`, 1, 0, `data-content="x+y"`},
		{"styled table formulas", `<table><tr><th>Formula</th></tr><tr><td><p><strong><span>\(x+y\)</span></strong></p><p><em><span>\(z\)</span></em></p></td></tr></table>`, 2, 0, `data-content="x+y"`},
		{"bold formula", `<p><strong>\(x+y\)</strong></p>`, 1, 0, `data-type="strong inline-math"`},
		{"nested styled formula", `<p><span data-type="strong">\(x+y\)</span></p>`, 1, 0, `data-content="x+y"`},
		{"nested existing styled formula", `<p><span data-type="strong em"><span data-type="inline-math" data-content="x+y"></span></span></p>`, 1, 0, `data-content="x+y"`},
		{"multiple split formulas", `<p>\(a<strong>+b\) then \(c</strong>+d\)</p>`, 2, 0, "then"},
		{"highlight split formula", `<p>\(x<mark>+y</mark>\)</p>`, 1, 0, `data-content="x+y"`},
		{"highlight display", `<p><mark>\[x+y\]</mark></p>`, 0, 1, `data-content="x+y"`},
		{"heading", `<h2><span>\(x+y\)</span></h2>`, 1, 0, `data-type="NodeHeading"`},
		{"list", `<ul><li><p><span>\(x+y\)</span></p></li></ul>`, 1, 0, `data-type="NodeList"`},
		{"quote", `<blockquote><p><span>\(x+y\)</span></p></blockquote>`, 1, 0, `data-type="NodeBlockquote"`},
		{"empty wrappers before existing math", `<p><span></span><span><span data-type="inline-math" data-content="x+y"></span></span></p>`, 1, 0, `data-content="x+y"`},
		{"link label", `<p><a href="https://example.com">\(x+y\)</a></p>`, 0, 0, "x+y"},
		{"unicode and entities", `<p>前 \(α+β\) 后 <span>\(x &lt; y\)</span></p>`, 2, 0, `data-content="x &lt; y"`},
		{"code boundary", `<p>\(x<code>+y</code>\)</p>`, 0, 0, "x"},
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
			if len(inlineMath.FindAllString(dom, -1)) != tc.inline || strings.Count(dom, `data-type="NodeMathBlock"`) != tc.block || !strings.Contains(dom, tc.contains) {
				t.Fatalf("unexpected conversion: %s", dom)
			}
			for range 3 {
				engine.SetLatexMath(false)
				spun := engine.SpinBlockDOM(dom)
				if len(inlineMath.FindAllString(spun, -1)) != tc.inline || strings.Count(spun, `data-type="NodeMathBlock"`) != tc.block || !strings.Contains(html.UnescapeString(spun), html.UnescapeString(tc.contains)) {
					t.Fatalf("canonical editing round trip changed formulas: %s", spun)
				}
				dom = spun
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

func TestLatexMathRichTableImport(t *testing.T) {
	for _, enabled := range []bool{true, false} {
		engine := lute.New()
		engine.SetProtyleWYSIWYG(true)
		engine.SetTextMark(true)
		engine.SetKramdownIAL(true)
		engine.SetLatexMath(enabled)
		engine.ParseOptions.PreserveTableCellSource = true
		source := "| Formula |\n| --- |\n| <p><span>\\(x+y\\)</span></p><p>\\[z\\]</p><pre><code>\\(code\\)</code></pre> |"
		tree := parse.Parse("", []byte(source), engine.ParseOptions)
		lute.NormalizeMarkdownTableCells(tree)
		var rich *ast.TableCellRich
		ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
			if entering && n.Type == ast.NodeTableCell && n.TableCellRich != nil {
				rich = n.TableCellRich
			}
			return ast.WalkContinue
		})
		if rich == nil {
			t.Fatal("missing rich table cell")
		}
		fragment, err := parse.ParseTableCellRich(rich)
		if err != nil {
			t.Fatal(err)
		}
		inline, block := 0, 0
		ast.Walk(fragment.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
			if entering && n.Type == ast.NodeTextMark && n.ContainTextMarkTypes("inline-math") {
				inline++
			}
			if entering && n.Type == ast.NodeMathBlock {
				block++
			}
			return ast.WalkContinue
		})
		want := 0
		if enabled {
			want = 1
		}
		if inline != want || block != want || !strings.Contains(rich.Content, `\(code\)`) {
			t.Fatalf("enabled=%v inline=%d block=%d content=%s", enabled, inline, block, rich.Content)
		}
		if parse.TableCellRichOptions().LatexMath {
			t.Fatal("import syntax leaked into stored rich table parsing")
		}
	}
}
