// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"bytes"
	"github.com/88250/lute/util"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
)

func ParagraphContinue(p *ast.Node, context *Context) int {
	if context.blank {
		return 1
	}
	return 0
}

func paragraphFinalize(p *ast.Node, context *Context) (insertTable bool) {
	p.Tokens = lex.TrimWhitespace(p.Tokens)

	// 尝试解析链接引用定义
	hasReferenceDefs := false
	for tokens := p.Tokens; 0 < len(tokens) && lex.ItemOpenBracket == tokens[0]; tokens = p.Tokens {
		if tokens = context.parseLinkRefDef(tokens); nil != tokens {
			p.Tokens = tokens
			hasReferenceDefs = true
			continue
		}
		break
	}
	if hasReferenceDefs && lex.IsBlankLine(p.Tokens) {
		p.Unlink()
	}

	if context.Option.GFMTaskListItem {
		// 尝试解析任务列表项
		if listItem := p.Parent; nil != listItem && ast.NodeListItem == listItem.Type && listItem.FirstChild == p {
			if 3 == listItem.ListData.Typ {
				isTaskListItem := false
				if !(context.Option.VditorWYSIWYG || context.Option.VditorIR || context.Option.VditorSV) {
					isTaskListItem = 3 < len(p.Tokens) && lex.IsWhitespace(p.Tokens[3])
				} else {
					isTaskListItem = 3 <= len(p.Tokens)
				}

				if isTaskListItem {
					// 如果是任务列表项则添加任务列表标记符节点
					tokens := p.Tokens
					var caretStartText, caretAfterCloseBracket, caretInBracket bool
					if context.Option.VditorWYSIWYG || context.Option.VditorIR || context.Option.VditorSV {
						closeBracket := bytes.IndexByte(tokens, lex.ItemCloseBracket)
						if bytes.HasPrefix(tokens, util.CaretTokens) {
							tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
							caretStartText = true
						} else if bytes.HasPrefix(tokens[closeBracket+1:], util.CaretTokens) {
							tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
							caretAfterCloseBracket = true
						} else if bytes.Contains(tokens[1:closeBracket], util.CaretTokens) {
							tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
							caretInBracket = true
						}
					}
					taskListItemMarker := &ast.Node{Type: ast.NodeTaskListItemMarker, Tokens: tokens[:3], TaskListItemChecked: listItem.ListData.Checked}
					p.PrependChild(taskListItemMarker)
					p.Tokens = tokens[3:] // 剔除开头的 [ ]、[x] 或者 [X]
					if context.Option.VditorWYSIWYG || context.Option.VditorIR || context.Option.VditorSV {
						p.Tokens = bytes.TrimSpace(p.Tokens)
						if caretStartText || caretAfterCloseBracket || caretInBracket {
							p.Tokens = append([]byte(" "+util.Caret), p.Tokens...)
						} else {
							p.Tokens = append([]byte(" "), p.Tokens...)
						}
					}
				}
			}
		}
	}

	if context.Option.GFMTable {
		if paragraph, table := context.parseTable(p); nil != table {
			if nil != paragraph {
				p.Tokens = paragraph.Tokens
				p.InsertAfter(table)
				// 设置末梢及其状态
				table.Close = true
				context.Tip = table
				return true
			} else {
				// 将该段落节点转成表节点
				p.Type = ast.NodeTable
				p.TableAligns = table.TableAligns
				for tr := table.FirstChild; nil != tr; {
					nextTr := tr.Next
					p.AppendChild(tr)
					tr = nextTr
				}
			}
			return
		}
	}

	if context.Option.ToC {
		if toc := context.parseToC(p); nil != toc {
			// 将该段落节点转换成目录节点
			p.Type = ast.NodeToC
			p.Tokens = toc.Tokens
			return
		}
	}
	return
}
