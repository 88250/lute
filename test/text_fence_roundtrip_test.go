package test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func TestLiteralTextFencesPreserveFollowingBlocks(t *testing.T) {
	for _, parts := range [][]string{
		{"<pre>\n  ```Java\nexample\n</pre>"},
		{"```"},
		{"~~~~text"},
		{"before\n   ````text\nafter"},
		{"literal ``` and ~~~~ markers"},
		{"`", "`", "`Java"},
		{"~", "~~text"},
		{"before\n ", "``", "`Java"},
	} {
		for _, export := range []bool{false, true} {
			t.Run(strings.Join(parts, "")+map[bool]string{false: "/format", true: "/export"}[export], func(t *testing.T) {
				engine := lute.New()
				engine.SetProtyleWYSIWYG(true)
				engine.SetAutoSpace(false)
				engine.SetParagraphBeginningSpace(true)
				tree := parse.Parse("", []byte("placeholder\n\n```java\noriginal();\n```\n\n# Following\n\n![gif](assets/animation.gif)\n\n![png](assets/picture.png)\n"), engine.ParseOptions)
				paragraph := tree.Root.FirstChild
				paragraph.FirstChild.Unlink()
				for _, part := range parts {
					paragraph.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(part)})
				}
				for pass := 0; pass < 3; pass++ {
					var markdown []byte
					if export {
						markdown = render.NewProtyleExportMdRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()
					} else {
						renderer := render.NewFormatRenderer(tree, engine.RenderOptions, engine.ParseOptions)
						renderer.PreserveTextCodeFences = true
						markdown = renderer.Render()
					}
					if pass == 0 {
						var source []byte
						for n := paragraph.FirstChild; n != nil; n = n.Next {
							source = append(source, n.Tokens...)
						}
						if string(source) != strings.Join(parts, "") {
							t.Fatalf("rendering mutated the source text: %q", source)
						}
					}
					tree = parse.Parse("", markdown, engine.ParseOptions)
					counts := map[ast.NodeType]int{}
					var code []byte
					assets := map[string]bool{}
					ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
						if entering {
							counts[n.Type]++
							if n.Type == ast.NodeCodeBlockCode {
								code = n.Tokens
							}
							if n.Type == ast.NodeImage {
								assets[string(n.ChildByType(ast.NodeLinkDest).Tokens)] = true
							}
						}
						return ast.WalkContinue
					})
					if counts[ast.NodeCodeBlock] != 1 || !bytes.Equal(code, []byte("original();\n")) || counts[ast.NodeHeading] != 1 ||
						!assets["assets/animation.gif"] || !assets["assets/picture.png"] {
						t.Fatalf("pass %d: following blocks changed: %v\n%s", pass, counts, markdown)
					}
					var text strings.Builder
					ast.Walk(tree.Root.FirstChild, func(n *ast.Node, entering bool) ast.WalkStatus {
						if entering && (n.Type == ast.NodeText || n.Type == ast.NodeBackslashContent) {
							text.Write(n.Tokens)
						}
						if entering && n.Type == ast.NodeSoftBreak {
							text.WriteByte('\n')
						}
						return ast.WalkContinue
					})
					if text.String() != strings.Join(parts, "") {
						t.Fatalf("pass %d: literal text changed: %q\n%s", pass, text.String(), markdown)
					}
				}
			})
		}
	}
}
