package lute

import (
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

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
				node.TabItemTitle = strings.TrimSpace(lute.HTML2Md(string(util.DomHTML(title))))
			}
		}
	}
	for _, name := range []string{"tabs-active-id", "tabs-position"} {
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
