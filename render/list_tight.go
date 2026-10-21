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
// 只使用本次已进入的列表或列表项的判定，不扫描渲染范围以外的节点。
func (r *ProtyleExportMdRenderer) inTightList(node *ast.Node) bool {
	list := enclosingList(node)
	if nil == list {
		return false
	}
	if tight, ok := r.tightLists[list]; ok {
		return tight
	}
	if tight, ok := r.tightLists[node]; ok {
		return tight
	}
	// 未进入列表或列表项时，正在单独渲染其段落等片段。片段末尾的换行会被
	// 调用方裁掉，无需扫描未渲染的兄弟块来决定列表整体的松紧。
	return true
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
		if listItemNeedsLoose(item) {
			return true
		}
	}
	return false
}

func listItemNeedsLoose(item *ast.Node) bool {
	layout := listItemLayout{}
	if layout.needsLoose(item) {
		return true
	}
	return 0 == layout.blocks
}

type listItemLayout struct {
	previous   *ast.Node
	blocks     int
	paragraphs int
}

func (layout *listItemLayout) needsLoose(container *ast.Node) bool {
	for child := container.FirstChild; nil != child; child = child.Next {
		if ast.NodeSuperBlock == child.Type {
			// 标准 Markdown 会展开超级块，必须检查展开后的块边界。
			if layout.needsLoose(child) {
				return true
			}
			continue
		}
		if !child.IsBlock() || ast.NodeKramdownBlockIAL == child.Type {
			continue
		}
		layout.blocks++
		if ast.NodeParagraph == child.Type {
			layout.paragraphs++
			if 1 < layout.paragraphs {
				return true
			}
		}
		if nil != layout.previous && blocksNeedBlankLine(layout.previous, child) {
			return true
		}
		layout.previous = child
	}
	return false
}

func blocksNeedBlankLine(previous, next *ast.Node) bool {
	// HTML 块可能延续到下一个空行，表格也不能与相邻的段落或表格拼接。
	// 页签可能展开为标题和正文，保守保留其与相邻块的分隔。
	switch previous.Type {
	case ast.NodeHTMLBlock, ast.NodeTable, ast.NodeTabs:
		return true
	case ast.NodeList, ast.NodeMindmap:
		switch next.Type {
		case ast.NodeParagraph, ast.NodeThematicBreak, ast.NodeList, ast.NodeMindmap:
			return true
		}
	case ast.NodeBlockquote, ast.NodeCallout:
		if ast.NodeBlockquote == next.Type || ast.NodeCallout == next.Type {
			return true
		}
	case ast.NodeParagraph, ast.NodeHeading:
		// 任务列表的首个标题前还会输出任务标记，需保留与段落相同的歧义边界。
		if needsBlankLineAfterParagraph(next) {
			return true
		}
	}
	switch next.Type {
	case ast.NodeHTMLBlock, ast.NodeTable, ast.NodeTabs:
		return true
	}
	return false
}

// needsBlankLineAfterParagraph 判断块是否无法紧跟在段落之后，即无法中断段落。
func needsBlankLineAfterParagraph(block *ast.Node) bool {
	switch block.Type {
	case ast.NodeList, ast.NodeMindmap:
		// 有序列表只有从 1 开始时才能中断段落
		return !listInterruptsParagraph(block)
	case ast.NodeParagraph, ast.NodeThematicBreak:
		// 分隔线导出为 ---，紧跟段落会被解析为 Setext 标题。
		return true
	case ast.NodeTable, ast.NodeHTMLBlock:
		// 表格和 HTML 块无法可靠中断段落，保守留空行
		return true
	}
	return false
}

// listInterruptsParagraph 判断列表能否中断段落：首项必须有内容，
// 有序列表还必须从 1 开始。
func listInterruptsParagraph(list *ast.Node) bool {
	if nil == list.ListData {
		return false
	}
	for item := list.FirstChild; nil != item; item = item.Next {
		if ast.NodeListItem != item.Type && ast.NodeMindmapItem != item.Type {
			continue
		}
		if nil == item.ListData || ((1 == item.ListData.Typ || (3 == item.ListData.Typ && 0 == item.ListData.BulletChar)) && 1 != item.ListData.Num) {
			return false
		}
		for child := item.FirstChild; nil != child; child = child.Next {
			if ast.NodeKramdownBlockIAL == child.Type || ast.NodeTaskListItemMarker == child.Type {
				continue
			}
			if ast.NodeParagraph != child.Type {
				return false
			}
			for inline := child.FirstChild; nil != inline; inline = inline.Next {
				if ast.NodeKramdownBlockIAL == inline.Type || ast.NodeKramdownSpanIAL == inline.Type {
					continue
				}
				if ast.NodeText != inline.Type || "" != strings.TrimSpace(inline.TokensStr()) {
					return true
				}
			}
			return false
		}
		return false
	}
	return false
}
