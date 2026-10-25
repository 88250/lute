package render

import (
	"encoding/json"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
)

// NewDocumentTitle 保存段落或标题块中的全部行级内容，不复用调用方的可变节点。
func NewDocumentTitle(container *ast.Node) (*ast.DocumentTitle, error) {
	if nil == container || (container.Type != ast.NodeParagraph && container.Type != ast.NodeHeading) {
		return nil, ast.ErrInvalidDocumentTitle
	}
	tree := &parse.Tree{Root: &ast.Node{Type: ast.NodeDocument}}
	paragraph := &ast.Node{Type: ast.NodeParagraph}
	tree.Root.AppendChild(paragraph)
	for node := container.FirstChild; nil != node; node = node.Next {
		if node.Type != ast.NodeHeadingC8hMarker {
			paragraph.AppendChild(parse.CloneDocumentTitleInline(node))
		}
	}
	parse.NestedInlines2FlattedSpansHybrid(tree, false)
	var prepare func(*ast.Node)
	prepare = func(node *ast.Node) {
		node.Data, node.TypeStr = string(node.Tokens), node.Type.String()
		node.Properties = ial2Map(node.KramdownIAL)
		delete(node.Properties, "refcount")
		delete(node.Properties, "av-names")
		for child := node.FirstChild; nil != child; child = child.Next {
			prepare(child)
			node.Children = append(node.Children, child)
		}
	}
	nodes := []*ast.Node{}
	for node := paragraph.FirstChild; nil != node; node = node.Next {
		prepare(node)
		nodes = append(nodes, node)
	}
	content, err := json.Marshal(nodes)
	if nil != err {
		return nil, ast.ErrInvalidDocumentTitle
	}
	ret := &ast.DocumentTitle{Spec: ast.DocumentTitleSpec, Format: ast.DocumentTitleFormat, Content: string(content)}
	if _, err = parse.ParseDocumentTitle(ret); nil != err {
		return nil, err
	}
	return ret, nil
}
