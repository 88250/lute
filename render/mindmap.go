// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"bytes"
	"strings"

	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

func EChartsMindmapStr(listContent string) string {
	return util.BytesToStr(echartsMindmap(util.StrToBytes(listContent)))
}

func EChartsMindmap(listContent []byte) []byte {
	return html.EncodeDestination(echartsMindmap(listContent))
}

// echartsMindmap 用于将列表 Markdown 原文转为 ECharts 树图结构，提供给前端渲染脑图。
func echartsMindmap(listContent []byte) []byte {
	listContent = bytes.ReplaceAll(listContent, editor.CaretTokens, nil)
	tree := parse.Parse("", listContent, parse.NewOptions())
	if nil == tree.Root.FirstChild || ast.NodeList != tree.Root.FirstChild.Type {
		// 第一个节点如果不是列表的话直接返回
		return []byte("{}")
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
			if entering {
				if needRoot(n) {
					// 如果根节点下的第一个列表包含多个列表项，则自动生成一个根节点，这些列表项都挂在这个根节点上
					buf.WriteString("{\"name\": \"Root\", \"children\": [")
				}
			} else {
				if needRoot(n) {
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
			return ast.WalkContinue
		}
		return ast.WalkContinue
	})
	return buf.Bytes()
}

// text 返回列表项第一个子节点的文本内容。
func text(listItemFirstChild *ast.Node) (ret string) {
	if nil == listItemFirstChild {
		return ""
	}

	buf := &bytes.Buffer{}
	ast.Walk(listItemFirstChild, func(n *ast.Node, entering bool) ast.WalkStatus {
		if ast.NodeList == n.Type || ast.NodeListItem == n.Type { // 遍历到下一个列表或者列表项时退出
			return ast.WalkContinue
		}

		if (ast.NodeText == n.Type || ast.NodeLinkText == n.Type) && entering {
			buf.Write(n.Tokens)
		}
		return ast.WalkContinue
	})

	ret = buf.String()
	ret = strings.ReplaceAll(ret, "\\", "\\\\")
	ret = strings.ReplaceAll(ret, "\"", "\\\"")
	ret = strings.ReplaceAll(ret, editor.Caret, "")
	return
}

func needRoot(root *ast.Node) bool {
	count := 0

	// 检查根节点下是否包含多个列表
	for c := root.FirstChild; nil != c; c = c.Next {
		if ast.NodeList == c.Type {
			count++
		}
	}
	if 1 < count {
		// 包含多个列表则需要构建一个 Root 节点
		return true
	}
	if 0 == count {
		// 没有列表也需要构建一个 Root 节点
		return true
	}

	count = 0

	// 如果只有一个列表，则检查该列表下是否包含多个列表项
	for c := root.FirstChild.FirstChild; nil != c; c = c.Next {
		if ast.NodeListItem == c.Type {
			count++
		}
	}
	if 1 < count {
		// 包含多个列表项则需要构建一个 Root 节点
		return true
	}
	return false
}
