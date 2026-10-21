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
	"strings"

	"github.com/88250/lute/ast"
)

// listTight 返回列表是否应按紧凑形式渲染，并缓存判定结果。
//
// 判据不依赖 ListData.Tight：由 BlockDOM 转换得到的列表没有紧凑信息，其值恒为 false。
// 这里按实际内容判断，只有紧凑形式会改变文档结构时才使用松散形式；判定结果只保存在渲染器
// 的缓存中，不写回语法树，因此同一棵树可以被多次渲染，且每次结果一致。
func (r *ProtyleExportMdRenderer) listTight(list *ast.Node) bool {
	if nil == list || nil == list.ListData {
		return false
	}
	if tight, ok := r.tightLists[list]; ok {
		return tight
	}

	tight := !listNeedsLoose(list)
	if nil == r.tightLists {
		r.tightLists = map[*ast.Node]bool{}
	}
	r.tightLists[list] = tight
	return tight
}

// inTightList 判断节点所属的列表是否按紧凑形式渲染。
// 按节点渲染时不会经过列表节点，因此需要由列表项反查其所属列表。
func (r *ProtyleExportMdRenderer) inTightList(node *ast.Node) bool {
	return r.listTight(enclosingList(node))
}

// enclosingList 返回节点所属的列表节点，节点本身是列表时返回自身。
func enclosingList(node *ast.Node) *ast.Node {
	if nil == node {
		return nil
	}
	switch node.Type {
	case ast.NodeList, ast.NodeMindmap:
		return node
	case ast.NodeListItem, ast.NodeMindmapItem:
		if nil != node.Parent {
			switch node.Parent.Type {
			case ast.NodeList, ast.NodeMindmap:
				return node.Parent
			}
		}
	}
	return nil
}

// ResetListTight 清除列表紧凑判定缓存，复用渲染器渲染新的节点前调用。
func (r *ProtyleExportMdRenderer) ResetListTight() {
	r.tightLists = nil
}

// listNeedsLoose 判断列表是否必须使用松散形式，即列表项内存在无法紧贴前一块的块，
// 或者列表项本身无法在紧凑形式下保持结构。
func listNeedsLoose(list *ast.Node) bool {
	for item := list.FirstChild; nil != item; item = item.Next {
		if ast.NodeListItem != item.Type && ast.NodeMindmapItem != item.Type {
			continue
		}

		var prev *ast.Node
		blocks, paragraphs := 0, 0
		for child := item.FirstChild; nil != child; child = child.Next {
			if ast.NodeKramdownBlockIAL == child.Type || ast.NodeTaskListItemMarker == child.Type {
				continue
			}

			blocks++
			if ast.NodeParagraph == child.Type {
				// 列表项内第二个及之后的段落必须用空行分隔
				paragraphs++
				if 1 < paragraphs {
					return true
				}
			} else if nil != prev && ast.NodeParagraph == prev.Type && needsBlankLineAfterParagraph(child) {
				return true
			}
			prev = child
		}

		// 空列表项在紧凑形式下会和相邻块合并，丢失原有层级
		if 0 == blocks {
			return true
		}
	}
	return false
}

// needsBlankLineAfterParagraph 判断块是否无法紧跟在段落之后，即无法中断段落。
func needsBlankLineAfterParagraph(block *ast.Node) bool {
	switch block.Type {
	case ast.NodeList, ast.NodeMindmap:
		// 有序列表只有从 1 开始时才能中断段落
		return !listInterruptsParagraph(block)
	case ast.NodeParagraph:
		return true
	case ast.NodeTable, ast.NodeHTMLBlock:
		// 表格和 HTML 块无法可靠中断段落，保守留空行
		return true
	}
	return false
}

// listInterruptsParagraph 判断列表能否中断段落：无序列表和任务列表总是可以，
// 有序列表仅在首项有内容且序号为 1 时可以。
func listInterruptsParagraph(list *ast.Node) bool {
	if nil == list.ListData {
		return false
	}
	if 1 != list.ListData.Typ {
		return true
	}

	for item := list.FirstChild; nil != item; item = item.Next {
		if ast.NodeListItem != item.Type && ast.NodeMindmapItem != item.Type {
			continue
		}
		if nil == item.ListData || 1 != item.ListData.Num {
			return false
		}
		for child := item.FirstChild; nil != child; child = child.Next {
			if ast.NodeKramdownBlockIAL == child.Type || ast.NodeTaskListItemMarker == child.Type {
				continue
			}
			return ast.NodeParagraph == child.Type && "" != strings.TrimSpace(child.Text())
		}
		return false
	}
	return false
}
