// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

import "bytes"

// Paragraph 描述了段落节点结构。
type Paragraph struct {
	*BaseNode
}

func (p *Paragraph) CanContain(nodeType int) bool {
	return false
}

func (p *Paragraph) Continue(context *Context) int {
	if context.blank {
		return 1
	}
	return 0
}

func (p *Paragraph) Finalize(context *Context) {
	p.tokens = bytes.TrimSpace(p.tokens)

	// 尝试解析链接引用定义
	hasReferenceDefs := false
	for tokens := p.tokens; 0 < len(tokens) && itemOpenBracket == tokens[0]; tokens = p.tokens {
		if tokens = context.parseLinkRefDef(tokens); nil != tokens {
			p.tokens = tokens
			hasReferenceDefs = true
			continue
		}
		break
	}
	if hasReferenceDefs && p.tokens.isBlankLine() {
		p.Unlink()
	}

	if context.option.GFMTaskListItem {
		// 尝试解析任务列表项
		if listItem, ok := p.parent.(*ListItem); ok {
			if 3 == listItem.listData.typ {
				// 如果是任务列表项则添加任务列表标记节点
				taskListItemMarker := &TaskListItemMarker{&BaseNode{typ: NodeTaskListItemMarker}, listItem.listData.checked}
				p.InsertBefore(p, taskListItemMarker)
				p.tokens = p.tokens[3:] // 剔除开头的 [ ]、[x] 或者 [X]
			}
		}
	}

	if context.option.GFMTable {
		// 尝试解析表
		lines := p.tokens.lines()
		table := context.parseTable(lines)
		if nil != table {
			p.InsertBefore(p, table)
			// 移除该段落所有内容 tokens，但段落节点本身暂时保留在语法树上
			// 在行级解析中，如果段落内容为空则从语法树上移除该段落节点
			// 这样处理的目的是让块级解析保持简单，在关闭未匹配的节点时只用判断段落类型
			p.tokens = nil
		}
	}
}

func (p *Paragraph) AcceptLines() bool {
	return true
}
