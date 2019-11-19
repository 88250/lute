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

func (p *Node) paragraphContinue(context *Context) int {
	if context.blank {
		return 1
	}
	return 0
}

func (p *Node) paragraphFinalize(context *Context) {
	p.tokens = trimWhitespace(p.tokens)

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
	if hasReferenceDefs && isBlankLine(p.tokens) {
		p.Unlink()
	}

	if context.option.GFMTaskListItem {
		// 尝试解析任务列表项
		listItem := p.parent
		if nil != listItem && NodeListItem == listItem.typ {
			if 3 == listItem.listData.typ && 3 < len(p.tokens) && isWhitespace(p.tokens[3]) {
				// 如果是任务列表项则添加任务列表标记符节点
				taskListItemMarker := &Node{typ: NodeTaskListItemMarker, tokens: p.tokens[:3], taskListItemChecked: listItem.listData.checked}
				p.InsertBefore(taskListItemMarker)
				p.tokens = p.tokens[3:] // 剔除开头的 [ ]、[x] 或者 [X]
			}
		}
	}

	if context.option.GFMTable {
		table := context.parseTable(p)
		if nil != table {
			// 将该段落节点转成表节点
			p.typ = NodeTable
			p.tableAligns = table.tableAligns
			for tr := table.firstChild; nil != tr; {
				nextTr := tr.next
				p.AppendChild(tr)
				tr = nextTr
			}
			p.tokens = nil
		}
	}
}
