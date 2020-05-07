// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package main

import (
	"bytes"
	"github.com/88250/lute/util"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/gopherjs/gopherjs/js"
)

func New(options map[string]map[string]*js.Object) *js.Object {
	engine := lute.New()
	engine.SetJSRenderers(options)
	return js.MakeWrapper(engine)
}

func main() {
	js.Global.Set("Lute", map[string]interface{}{
		"Version":          lute.Version,
		"New":              New,
		"WalkStop":         ast.WalkStop,
		"WalkSkipChildren": ast.WalkSkipChildren,
		"WalkContinue":     ast.WalkContinue,
		"GetHeadingID":     render.HeadingID,
		"RenderMindmap":    renderMindmap,
	})
}

// renderMindmap 用于将列表 Markdown 原文转为 ECharts 树图结构，提供给前端渲染脑图。
func renderMindmap(listContent string) string {
	tree := parse.Parse("", []byte(listContent), lute.NewOptions())
	if nil == tree.Root.FirstChild || ast.NodeList != tree.Root.FirstChild.Type {
		// 第一个节点如果不是列表的话直接返回
		return "{}"
	}

	// 移除非列表节点
	var toRemoved []*ast.Node
	for c := tree.Root.FirstChild; nil != c; c = c.Next {
		if ast.NodeList != c.Type {
			toRemoved = append(toRemoved, c)
		}
	}
	for _, c := range toRemoved {
		c.Unlink()
	}

	buf := &bytes.Buffer{}
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		switch n.Type {
		case ast.NodeDocument:
			listItems := countListItem(n)
			if entering {
				if 0 < listItems {
					// 如果根节点下的第一个列表包含多个列表项，则自动生成一个根节点，这些列表项都挂在这个根节点上
					buf.WriteString("{\"name\": \"Root\", \"children\": [")
				}
			} else {
				if 0 < listItems {
					buf.WriteString("]}")
				}
			}
			return ast.WalkContinue
		case ast.NodeList:
			return ast.WalkContinue
		case ast.NodeListItem:
			children := nil != n.ChildByType(ast.NodeList)
			if entering {
				buf.WriteString("{\"name\": \"" + text(n.FirstChild) + "\"")
				if children {
					buf.WriteString(", \"children\": [")
				}
			} else {
				if children {
					buf.WriteString("]")
				}
				buf.WriteString("}")
				if nil != n.Next || nil != n.Parent.Next {
					buf.WriteString(", ")
				}
			}
		default:
			return ast.WalkStop
		}
		return ast.WalkContinue
	})
	return buf.String()
}

// text 返回列表项第一个子节点的文本内容。
func text(listItemFirstChild *ast.Node) (ret string) {
	if nil == listItemFirstChild {
		return ""
	}

	ast.Walk(listItemFirstChild, func(n *ast.Node, entering bool) ast.WalkStatus {
		if ast.NodeList == n.Type || ast.NodeListItem == n.Type { // 遍历到下一个列表或者列表项时退出
			return ast.WalkStop
		}

		if (ast.NodeText == n.Type || ast.NodeLinkText == n.Type) && entering {
			ret += util.BytesToStr(n.Tokens)
		}
		return ast.WalkContinue
	})
	return
}

func countListItem(n *ast.Node) (ret int) {
	if nil == n {
		return 0
	}

	for c := n.FirstChild; nil != c; c = c.Next {
		if ast.NodeList == c.Type || ast.NodeListItem == c.Type {
			ret++
		}
	}

	return
}
