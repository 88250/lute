package parse

import (
	"bytes"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
)

// TabsStart 解析页签围栏，代码等叶子块已经在进入块起始解析前处理。
func TabsStart(t *Tree, container *ast.Node) int {
	context := t.Context
	if !context.ParseOption.Tabs || context.indented {
		return 0
	}
	line := lex.TrimWhitespace(context.currentLine[context.nextNonspace:])
	if bytes.Equal(line, []byte(":::")) {
		var target *ast.Node
		for parent := container; nil != parent; parent = parent.Parent {
			if ast.NodeTabs == parent.Type || ast.NodeTabItem == parent.Type {
				target = parent
				break
			}
		}
		if nil == target {
			return 0
		}
		context.closeUnmatchedBlocks()
		for context.Tip != target {
			context.finalize(context.Tip)
		}
		context.finalize(target)
		return 3
	}
	if bytes.Equal(line, []byte(":::tabs")) {
		context.closeUnmatchedBlocks()
		context.addChild(ast.NodeTabs)
		context.offset = context.currentLineLen - 1
		return 1
	}
	if !bytes.Equal(line, []byte(":::tab")) && !bytes.HasPrefix(line, []byte(":::tab ")) &&
		!bytes.HasPrefix(line, []byte(":::tab\t")) {
		return 0
	}
	var tabs *ast.Node
	for parent := container; nil != parent; parent = parent.Parent {
		if ast.NodeTabs == parent.Type {
			tabs = parent
			break
		}
	}
	if nil == tabs {
		return 0
	}
	context.closeUnmatchedBlocks()
	for context.Tip != tabs {
		context.finalize(context.Tip)
	}
	item := context.addChild(ast.NodeTabItem)
	item.TabItemTitle = strings.TrimSpace(string(line[len(":::tab"):]))
	context.offset = context.currentLineLen - 1
	return 1
}

func (context *Context) tabsFinalize(node *ast.Node) {
	if ast.NodeTabs == node.Type {
		for child := node.FirstChild; nil != child; child = child.Next {
			if ast.NodeTabItem == child.Type {
				return
			}
		}
		item := &ast.Node{Type: ast.NodeTabItem, Close: true}
		item.AppendChild(&ast.Node{Type: ast.NodeParagraph, Close: true})
		node.AppendChild(item)
		return
	}
	for child := node.FirstChild; nil != child; child = child.Next {
		if child.IsBlock() && ast.NodeKramdownBlockIAL != child.Type {
			return
		}
	}
	node.AppendChild(&ast.Node{Type: ast.NodeParagraph, Close: true})
}
