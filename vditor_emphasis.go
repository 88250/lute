package lute

import (
	"bytes"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func (lute *Lute) formatVditorTree(tree *parse.Tree) string {
	options := render.NewOptions()
	options.AutoSpace = false
	options.FixTermTypo = false
	renderer := render.NewFormatRenderer(tree, options, lute.ParseOptions)
	formatted := renderer.Render()
	hasAdjacentEmphasis := false
	ast.Walk(tree.Root, func(node *ast.Node, entering bool) ast.WalkStatus {
		if entering && (node.Type == ast.NodeStrong || node.Type == ast.NodeEmphasis) && nil != node.Next &&
			(node.Next.Type == ast.NodeStrong || node.Next.Type == ast.NodeEmphasis) {
			hasAdjacentEmphasis = true
		}
		return ast.WalkContinue
	})
	if !hasAdjacentEmphasis {
		return string(formatted)
	}
	parsed := parse.Parse("", formatted, lute.ParseOptions)
	expected := render.NewHtmlRenderer(tree, options, lute.ParseOptions).Render()
	actual := render.NewHtmlRenderer(parsed, options, lute.ParseOptions).Render()
	if bytes.Equal(expected, actual) {
		return string(formatted)
	}

	// 连续强调分隔符可能被重新配对，使用 HTML 标签保留原有的格式边界。
	markers := map[ast.NodeType]string{
		ast.NodeEmA6kOpenMarker: "<em>", ast.NodeEmU8eOpenMarker: "<em>",
		ast.NodeEmA6kCloseMarker: "</em>", ast.NodeEmU8eCloseMarker: "</em>",
		ast.NodeStrongA6kOpenMarker: "<strong>", ast.NodeStrongU8eOpenMarker: "<strong>",
		ast.NodeStrongA6kCloseMarker: "</strong>", ast.NodeStrongU8eCloseMarker: "</strong>",
	}
	renderer = render.NewFormatRenderer(tree, options, lute.ParseOptions)
	for nodeType, tag := range markers {
		value := tag
		renderer.ExtRendererFuncs[nodeType] = func(node *ast.Node, entering bool) (string, ast.WalkStatus) {
			if entering {
				return value, ast.WalkContinue
			}
			return "", ast.WalkContinue
		}
	}
	return string(renderer.Render())
}
