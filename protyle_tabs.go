package lute

import (
	"bytes"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// tabTitleMarkdown 用边界文本保护标题空白，并以字符实体保留 Markdown 行首、行尾的空格。
func (lute *Lute) tabTitleMarkdown(title *html.Node, standardHTML bool) string {
	var dom bytes.Buffer
	dom.WriteString("LuteTabTitleStart")
	for child := title.FirstChild; nil != child; child = child.NextSibling {
		dom.Write(util.DomHTML(child))
	}
	dom.WriteString("LuteTabTitleEnd")
	var markdown string
	if standardHTML {
		markdown = lute.HTML2Md(dom.String())
	} else {
		markdown = lute.BlockDOM2Md(dom.String())
	}
	markdown = strings.TrimSpace(markdown)
	markdown = strings.TrimPrefix(markdown, "LuteTabTitleStart")
	markdown = strings.TrimSuffix(markdown, "LuteTabTitleEnd")
	leading := len(markdown) - len(strings.TrimLeft(markdown, " "))
	trailing := len(markdown) - len(strings.TrimRight(markdown, " "))
	if leading == len(markdown) {
		return strings.Repeat("&#32;", leading)
	}
	return strings.Repeat("&#32;", leading) + markdown[leading:len(markdown)-trailing] + strings.Repeat("&#32;", trailing)
}

func hasDOMClass(node *html.Node, name string) bool {
	for _, class := range strings.Fields(util.DomAttrValue(node, "class")) {
		if class == name {
			return true
		}
	}
	return false
}

func directDOMChildByClass(node *html.Node, name string) *html.Node {
	for child := node.FirstChild; nil != child; child = child.NextSibling {
		if hasDOMClass(child, name) {
			return child
		}
	}
	return nil
}

func (lute *Lute) genASTByTabsDOM(dom *html.Node, tree *parse.Tree) bool {
	if !lute.ParseOptions.Tabs || (!hasDOMClass(dom, "tabs") && !hasDOMClass(dom, "tab-item")) {
		return false
	}
	node := &ast.Node{Type: ast.NodeTabs}
	if hasDOMClass(dom, "tab-item") {
		node.Type = ast.NodeTabItem
		if info := directDOMChildByClass(dom, "tab-item-info"); nil != info {
			if title := directDOMChildByClass(info, "tab-item-title"); nil != title {
				node.TabItemTitle = lute.tabTitleMarkdown(title, true)
			}
		}
	}
	for _, name := range []string{"tabs-active-id", "tabs-position", "tabs-task"} {
		if value := util.DomAttrValue(dom, name); "" != value {
			node.SetIALAttr(name, value)
		}
	}
	node.ID = util.DomAttrValue(dom, "data-node-id")
	if "" != node.ID {
		node.SetIALAttr("id", node.ID)
	}
	tree.Context.Tip.AppendChild(node)
	tree.Context.Tip = node
	if ast.NodeTabItem == node.Type {
		if info := directDOMChildByClass(dom, "tab-item-info"); nil != info {
			for child := info.FirstChild; nil != child; child = child.NextSibling {
				if "true" == util.DomAttrValue(child, "tabs-title") {
					lute.genASTByBlockDOM(child, tree)
					break
				}
			}
		}
		if content := directDOMChildByClass(dom, "tab-item-content"); nil != content {
			for child := content.FirstChild; nil != child; child = child.NextSibling {
				lute.genASTByDOM(child, tree)
			}
		}
	} else {
		for child := dom.FirstChild; nil != child; child = child.NextSibling {
			if hasDOMClass(child, "tab-item") {
				lute.genASTByTabsDOM(child, tree)
			}
		}
	}
	tree.Context.ParentTip()
	if 0 < len(node.KramdownIAL) {
		node.InsertAfter(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: parse.IAL2Tokens(node.KramdownIAL)})
	}
	return true
}

// wrapTabItemFragments 为独立编辑的页签项添加临时语法容器，保持内部 Markdown 可以完整解析。
func wrapTabItemFragments(tree *parse.Tree) (wrappers map[string]bool) {
	wrappers = map[string]bool{}
	for node := tree.Root.FirstChild; nil != node; {
		next := node.Next
		if ast.NodeTabItem == node.Type {
			if "" == node.ID {
				node.ID = ast.NewNodeID()
				node.SetIALAttr("id", node.ID)
			}
			wrapper := &ast.Node{Type: ast.NodeTabs, ID: ast.NewNodeID()}
			wrapper.SetIALAttr("id", wrapper.ID)
			node.InsertBefore(wrapper)
			wrapper.AppendChild(node)
			if nil != next && ast.NodeKramdownBlockIAL == next.Type {
				ial := next
				next = ial.Next
				wrapper.AppendChild(ial)
			} else {
				wrapper.AppendChild(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: parse.IAL2Tokens(node.KramdownIAL)})
			}
			wrapper.InsertAfter(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: parse.IAL2Tokens(wrapper.KramdownIAL)})
			wrappers[wrapper.ID] = true
		}
		node = next
	}
	return
}

func unwrapTabItemFragments(tree *parse.Tree, wrappers map[string]bool) {
	for node := tree.Root.FirstChild; nil != node; {
		next := node.Next
		if ast.NodeTabs == node.Type && wrappers[node.ID] {
			for child := node.FirstChild; nil != child; {
				childNext := child.Next
				node.InsertBefore(child)
				child = childNext
			}
			if nil != next && ast.NodeKramdownBlockIAL == next.Type {
				ial := next
				next = next.Next
				ial.Unlink()
			}
			node.Unlink()
		}
		node = next
	}
}
