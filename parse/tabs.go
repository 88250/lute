package parse

import (
	"bytes"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
)

type tabsParseState struct {
	fenceLen int
	indent   int
	active   *ast.Node
}

// TabsStart 解析页签围栏和分隔符，代码等叶子块由块解析流程保护。
func TabsStart(t *Tree, container *ast.Node) int {
	context := t.Context
	if !context.ParseOption.Tabs {
		return 0
	}
	line := lex.TrimWhitespace(context.currentLine[context.nextNonspace:])
	fenceLen := 0
	for fenceLen < len(line) && ':' == line[fenceLen] {
		fenceLen++
	}
	if 3 <= fenceLen {
		if fenceLen == len(line) {
			if context.indented {
				return 0
			}
			for target := tabsContainer(container); nil != target; target = tabsContainer(target.Parent) {
				if state := context.tabs[target]; nil != state && state.fenceLen == fenceLen {
					context.closeUnmatchedBlocks()
					for context.Tip != target {
						context.finalize(context.Tip)
					}
					context.finalize(target)
					return 3
				}
			}
			return 0
		}
		if (' ' == line[fenceLen] || '\t' == line[fenceLen]) &&
			bytes.Equal(bytes.TrimSpace(line[fenceLen:]), []byte("tabs")) {
			nested := false
			// 包括引述和列表中的子组，所有后代围栏都必须短于外层。
			for parent := container; nil != parent; parent = parent.Parent {
				if state := context.tabs[parent]; nil != state {
					nested = true
					if fenceLen >= state.fenceLen {
						return 0
					}
				}
			}
			if context.indented && !nested {
				return 0
			}
			context.closeUnmatchedBlocks()
			node := context.addChild(ast.NodeTabs)
			if nil == context.tabs {
				context.tabs = map[*ast.Node]*tabsParseState{}
			}
			context.tabs[node] = &tabsParseState{fenceLen: fenceLen, indent: context.indent}
			return 3
		}
		return 0
	}
	if context.indented {
		return 0
	}
	marker := "@tab"
	active := bytes.HasPrefix(line, []byte("@tab:active"))
	if active {
		marker = "@tab:active"
	}
	if !bytes.Equal(line, []byte(marker)) && !bytes.HasPrefix(line, []byte(marker+" ")) &&
		!bytes.HasPrefix(line, []byte(marker+"\t")) {
		return 0
	}
	tabs := tabsContainer(container)
	if nil == tabs {
		return 0
	}
	context.closeUnmatchedBlocks()
	for context.Tip != tabs {
		context.finalize(context.Tip)
	}
	item := context.addChild(ast.NodeTabItem)
	item.TabItemTitle = strings.TrimSpace(string(line[len(marker):]))
	context.tabsItemIAL = item
	if state := context.tabs[tabs]; active && nil != state && nil == state.active {
		state.active = item
	}
	return 3
}

// tabsContainer 不跨越正文中的显式容器前缀识别页签分隔符。
func tabsContainer(container *ast.Node) *ast.Node {
	for parent := container; nil != parent; parent = parent.Parent {
		switch parent.Type {
		case ast.NodeTabs:
			return parent
		case ast.NodeParagraph, ast.NodeTabItem, ast.NodeList, ast.NodeKramdownBlockIAL:
			continue
		default:
			return nil
		}
	}
	return nil
}

func (context *Context) tabsContinue(node *ast.Node) int {
	if state := context.tabs[node]; nil != state {
		// 每层仅消费本层开围栏的相对缩进，正文额外的四格缩进仍表示代码。
		indent := state.indent
		if context.indent < indent {
			indent = context.indent
		}
		context.advanceOffset(indent, true)
	}
	return 0
}

func (t *Tree) finalParseTabs() {
	for tabs, state := range t.Context.tabs {
		active := state.active
		if nil == active {
			for item := tabs.FirstChild; nil != item; item = item.Next {
				if ast.NodeTabItem == item.Type && "" != item.ID && item.ID == tabs.IALAttr("tabs-active-id") {
					active = item
					break
				}
			}
		}
		if nil == active {
			for item := tabs.FirstChild; nil != item; item = item.Next {
				if ast.NodeTabItem == item.Type {
					active = item
					break
				}
			}
		}
		if nil == active {
			continue
		}
		if "" == active.ID {
			active.ID = ast.NewNodeID()
			active.SetIALAttr("id", active.ID)
		}
		tabs.SetIALAttr("tabs-active-id", active.ID)
		if nil != tabs.Next && ast.NodeKramdownBlockIAL == tabs.Next.Type {
			tabs.Next.Tokens = IAL2Tokens(tabs.KramdownIAL)
		}
	}
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
	if 0 < len(node.KramdownIAL) && (nil == node.Next || ast.NodeKramdownBlockIAL != node.Next.Type) {
		// 标题后的属性直到页签最终化时才挂为后继，避免阻断打开的正文容器。
		node.InsertAfter(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: IAL2Tokens(node.KramdownIAL), Close: true})
	}
	for child := node.FirstChild; nil != child; child = child.Next {
		if child.IsBlock() && ast.NodeKramdownBlockIAL != child.Type {
			return
		}
	}
	node.AppendChild(&ast.Node{Type: ast.NodeParagraph, Close: true})
}
