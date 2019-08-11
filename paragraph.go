// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

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
	p.tokens = p.tokens.trim()

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

	// 尝试解析任务列表项
	if listItem, ok := p.parent.(*ListItem); ok {
		if 3 == listItem.listData.typ {
			// 如果是任务列表项则添加任务列表标记节点
			taskListItemMarker := &TaskListItemMarker{&BaseNode{typ: NodeTaskListItemMarker}, listItem.listData.checked}
			p.InsertBefore(p, taskListItemMarker)
			p.tokens = p.tokens[3:] // 剔除开头的 [ ]、[x] 或者 [X]
		}
	}

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

func (p *Paragraph) AcceptLines() bool {
	return true
}
