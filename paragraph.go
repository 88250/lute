// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lute

import (
	"bytes"
	"github.com/88250/lute/ast"
)

func ParagraphContinue(p *ast.Node, context *Context) int {
	if context.blank {
		return 1
	}
	return 0
}

func paragraphFinalize(p *ast.Node, context *Context) {
	p.Tokens = trimWhitespace(p.Tokens)

	// 尝试解析链接引用定义
	hasReferenceDefs := false
	for tokens := p.Tokens; 0 < len(tokens) && itemOpenBracket == tokens[0]; tokens = p.Tokens {
		if tokens = context.parseLinkRefDef(tokens); nil != tokens {
			p.Tokens = tokens
			hasReferenceDefs = true
			continue
		}
		break
	}
	if hasReferenceDefs && isBlankLine(p.Tokens) {
		p.Unlink()
	}

	if context.option.GFMTaskListItem {
		// 尝试解析任务列表项
		if listItem := p.Parent; nil != listItem && ast.NodeListItem == listItem.Type && listItem.FirstChild == p {
			if 3 == listItem.ListData.Typ {
				isTaskListItem := false
				if !context.option.VditorWYSIWYG {
					isTaskListItem = 3 < len(p.Tokens) && isWhitespace(p.Tokens[3])
				} else {
					isTaskListItem = 3 <= len(p.Tokens)
				}

				if isTaskListItem {
					// 如果是任务列表项则添加任务列表标记符节点
					taskListItemMarker := &ast.Node{Type: ast.NodeTaskListItemMarker, Tokens: p.Tokens[:3], TaskListItemChecked: listItem.ListData.Checked}
					p.PrependChild(taskListItemMarker)
					p.Tokens = p.Tokens[3:] // 剔除开头的 [ ]、[x] 或者 [X]
					if context.option.VditorWYSIWYG {
						if 1 > len(p.Tokens) {
							p.Tokens = []byte(" ")
						} else {
							if !bytes.HasPrefix(p.Tokens, []byte(" ")) {
								p.Tokens = append([]byte(" "), p.Tokens...)
							}
						}
					}
				}
			}
		}
	}

	if context.option.GFMTable {
		if table := context.parseTable(p); nil != table {
			// 将该段落节点转成表节点
			p.Type = ast.NodeTable
			p.TableAligns = table.TableAligns
			for tr := table.FirstChild; nil != tr; {
				nextTr := tr.Next
				p.AppendChild(tr)
				tr = nextTr
			}
			p.Tokens = nil
			return
		}
	}

	if context.option.ToC {
		if toc := context.parseToC(p); nil != toc {
			// 将该段落节点转换成目录节点
			p.Type = ast.NodeToC
			p.Tokens = nil
			return
		}
	}
}
