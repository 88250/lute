package test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/parse"
)

func TestHTML2TreeEscapedText(t *testing.T) {
	for _, enabled := range []bool{true, false} {
		engine := lute.New()
		engine.SetProtyleWYSIWYG(true)
		engine.SetTextMark(true)
		engine.SetHTMLTag2TextMark(true)
		engine.SetInlineAsterisk(enabled)
		engine.SetInlineUnderscore(enabled)
		engine.SetGFMStrikethrough(enabled)
		engine.SetMark(enabled)
		engine.SetSup(enabled)
		engine.SetSub(enabled)
		for _, tag := range []string{"strong", "em", "s", "mark", "sup", "sub"} {
			t.Run(tag+"/syntax="+strconv.FormatBool(enabled), func(t *testing.T) {
				const source = "a = b \\= \\ *x* _y_ ~z~ $m$ ^n^ <t> `c`"
				tree := engine.HTML2Tree("<" + tag + ">" + string(html.EscapeHTML([]byte(source))) + "</" + tag + ">")
				parse.TextMarks2Inlines(tree)
				parse.NestedInlines2FlattedSpansHybrid(tree, false)
				var actual strings.Builder
				ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
					if entering && n.Type == ast.NodeTextMark {
						if n.TextMarkType != tag {
							t.Errorf("unexpected text mark type %q", n.TextMarkType)
						}
						actual.WriteString(html.UnescapeString(n.TextMarkTextContent))
					}
					return ast.WalkContinue
				})
				if actual.String() != source {
					t.Fatalf("expected %q, got %q", source, actual.String())
				}
			})
		}
	}
}

func TestHTML2MarkdownEscapedText(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	for _, test := range []struct{ name, source, expected string }{
		{"plain", `<p>a = b \=</p>`, "a \\= b \\\\\\=\n"},
		{"rawSpan", `<span>$x=2$</span>`, "$x=2$\n"},
		{"attributedSpan", `<span class="">$x=2$</span>`, "\\$x\\=2\\$\n"},
		{"link", `<a href="https://example.com">a = b</a>`, "[a = b](https://example.com)\n"},
		{"code", `<code>a = b \=</code>`, "`a = b \\=`\n"},
		{"math", `<p>$$x=2$$</p>`, "$$\nx=2\n$$\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			actual, err := engine.HTML2Markdown(test.source)
			if err != nil || actual != test.expected {
				t.Fatalf("expected %q, got %q, err %v", test.expected, actual, err)
			}
		})
	}
}
