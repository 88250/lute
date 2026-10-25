package lute

import (
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

// SpinInlineBlockDOM parses editable Markdown inline syntax in a single paragraph.
// Block markers stay literal and paragraph identity, attributes and soft breaks
// are retained. Other block structures are returned unchanged.
func (lute *Lute) SpinInlineBlockDOM(ivHTML string) string {
	keepEscaped := lute.ParseOptions.KeepEscaped
	lute.ParseOptions.KeepEscaped = true
	defer func() { lute.ParseOptions.KeepEscaped = keepEscaped }()

	fragment := lute.BlockDOM2Tree(ivHTML)
	var paragraph *ast.Node
	for child := fragment.Root.FirstChild; child != nil; child = child.Next {
		if child.Type == ast.NodeKramdownBlockIAL {
			continue
		}
		if child.Type != ast.NodeParagraph || paragraph != nil {
			return ivHTML
		}
		paragraph = child
	}
	if paragraph == nil {
		return ivHTML
	}
	normalizeSpinCaretNewline(fragment)
	options := render.NewOptions()
	options.AutoSpace, options.FixTermTypo = false, false
	options.KramdownSpanIAL, options.KeepParagraphBeginningSpace, options.ProtyleWYSIWYG = true, true, true
	renderer := render.NewFormatRenderer(fragment, options, lute.ParseOptions)
	// The normal document/paragraph renderers trim text and add block boundaries.
	// Inline parsing needs the original inline stream, including its final break.
	renderer.RendererFuncs[ast.NodeDocument] = func(_ *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			renderer.NodeWriterStack = append(renderer.NodeWriterStack, renderer.Writer)
		} else {
			renderer.NodeWriterStack = renderer.NodeWriterStack[:len(renderer.NodeWriterStack)-1]
		}
		return ast.WalkContinue
	}
	renderer.RendererFuncs[ast.NodeParagraph] = func(_ *ast.Node, _ bool) ast.WalkStatus { return ast.WalkContinue }
	renderer.RendererFuncs[ast.NodeKramdownBlockIAL] = func(_ *ast.Node, _ bool) ast.WalkStatus { return ast.WalkSkipChildren }
	markdown := strings.ReplaceAll(string(renderer.Render()), editor.Zwsp, "")
	tree := parse.Inline("", []byte(markdown), lute.ParseOptions)
	tree.Root.FirstChild.ID = paragraph.ID
	tree.Root.FirstChild.KramdownIAL = paragraph.KramdownIAL
	// Existing semantic text marks already carry their source content. Parsing
	// their text again would turn a saved literal "**example**" into new styles
	// when another part of the title is edited. Raw input was parsed by Inline.
	parse.NestedInlines2FlattedSpansHybrid(tree, false)
	return lute.Tree2BlockDOM(tree, lute.RenderOptions, lute.ParseOptions)
}
