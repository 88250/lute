package render

import (
	"bytes"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
)

// tableCellRichInlineHTML 在 Markdown 单元格中保留行内格式，用硬换行承接片段的段落和代码行。
func tableCellRichInlineHTML(cell *ast.Node, options *Options, parseOptions *parse.Options) []byte {
	inline := parse.CloneTableCellInline(cell)
	inline.TableCellRich = nil
	opt := *options
	opt.ProtyleWYSIWYG, opt.KramdownBlockIAL, opt.KramdownSpanIAL = false, false, false
	opt.Sanitize, opt.SoftBreak2HardBreak = true, true
	renderer := NewHtmlRenderer(&parse.Tree{Root: inline}, &opt, parseOptions)
	renderer.SetTextMarkStandardTag()
	renderer.RendererFuncs[ast.NodeTableCell] = func(*ast.Node, bool) ast.WalkStatus { return ast.WalkContinue }
	return bytes.ReplaceAll(renderer.Render(), []byte("|"), []byte("&#124;"))
}

// TableCellRichHTML 输出片段的静态 HTML，不把内部节点作为可编辑文档块暴露。
func TableCellRichHTML(rich *ast.TableCellRich, options *Options) ([]byte, error) {
	tree, err := parse.ParseTableCellRich(rich)
	if nil != err {
		return nil, err
	}
	var remove []*ast.Node
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering && n.IsBlock() {
			n.ID = ""
			n.RemoveIALAttr("id")
			n.RemoveIALAttr("updated")
			if ast.NodeKramdownBlockIAL == n.Type {
				remove = append(remove, n)
			}
		}
		return ast.WalkContinue
	})
	for _, n := range remove {
		n.Unlink()
	}
	opt := *options
	opt.ProtyleWYSIWYG = false
	opt.KramdownBlockIAL, opt.KramdownSpanIAL = false, false
	opt.CodeSyntaxHighlight = false
	opt.Sanitize = true
	opt.HeadingID = false
	return NewHtmlRenderer(tree, &opt, tree.Context.ParseOption).Render(), nil
}
