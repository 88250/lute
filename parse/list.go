// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"strconv"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

// ListStart 判断列表、列表项（* - + 1.）或者任务列表项是否开始。
func ListStart(t *Tree, container *ast.Node) int {
	if ast.NodeList != container.Type && t.Context.indented {
		return 0
	}

	data, ial := t.parseListMarker(container)
	if nil == data {
		return 0
	}

	t.Context.closeUnmatchedBlocks()

	listsMatch := container.Type == ast.NodeList && t.Context.listsMatch(container.ListData, data)
	if t.Context.Tip.Type != ast.NodeList || !listsMatch {
		list := t.Context.addChild(ast.NodeList)
		list.ListData = data
	}
	listItem := t.Context.addChild(ast.NodeListItem)
	listItem.ListData = data
	if t.Context.ParseOption.KramdownBlockIAL && nil != ial {
		listItem.KramdownIAL = ial
		listItem.ID = listItem.IALAttr("id")
		t.Context.offset += len(IAL2Tokens(ial))
	}

	listItem.Tokens = data.Marker
	if 1 == listItem.ListData.Typ || (3 == listItem.ListData.Typ && 0 == listItem.ListData.BulletChar) {
		// 修正有序列表项序号
		prev := listItem.Previous
		if nil != prev {
			listItem.ListData.Num = prev.ListData.Num + 1
		} else {
			listItem.ListData.Num = data.Start
		}
	}
	return 1
}

func ListItemContinue(listItem *ast.Node, context *Context) int {
	if context.blank {
		if nil == listItem.FirstChild { // 列表项后面是空的
			return 1
		}

		context.advanceNextNonspace()
	} else if context.indent >= listItem.ListData.MarkerOffset+listItem.ListData.Padding {
		context.advanceOffset(listItem.ListData.MarkerOffset+listItem.ListData.Padding, true)
	} else {
		return 1
	}
	return 0
}

func (context *Context) listFinalize(list *ast.Node) {
	item := list.FirstChild

	// 检查子列表项之间是否包含空行，包含的话说明该列表是非紧凑的，即松散的
	for nil != item {
		if endsWithBlankLine(item) && nil != item.Next {
			list.ListData.Tight = false
			break
		}

		subItem := item.FirstChild
		for nil != subItem {
			if endsWithBlankLine(subItem) && (nil != item.Next || nil != subItem.Next) {
				list.ListData.Tight = false
				break
			}
			subItem = subItem.Next
		}
		item = item.Next
	}

	if context.ParseOption.KramdownBlockIAL {
		for li := list.FirstChild; nil != li; li = li.Next {
			if nil == li.FirstChild {
				if ast.NodeKramdownBlockIAL != li.Type {
					id := ast.NewNodeID()
					ialTokens := []byte("{: id=\"" + id + "\"}")
					li.KramdownIAL = [][]string{{"id", id}}
					li.ID = id
					li.InsertAfter(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: ialTokens})

					id = ast.NewNodeID()
					ialTokens = []byte("{: id=\"" + id + "\"}")
					p := &ast.Node{Type: ast.NodeParagraph, ID: id}
					p.KramdownIAL = [][]string{{"id", id}}
					p.ID = id
					p.InsertAfter(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: ialTokens})
					li.AppendChild(p)
					li = li.Next
				}

				continue
			}

			if 7 < len(li.FirstChild.Tokens) && '{' == li.FirstChild.Tokens[0] {
				if ial := context.parseKramdownIALInListItem(li.FirstChild.Tokens); 0 < len(ial) {
					li.KramdownIAL = ial
					li.ID = ial[0][1]
					ialTokens := IAL2Tokens(ial)
					li.InsertAfter(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: ialTokens})
					tokens := li.FirstChild.Tokens[bytes.Index(li.FirstChild.Tokens, []byte("}"))+1:]
					tokens = lex.TrimWhitespace(tokens)
					li.FirstChild.Tokens = tokens
					li = li.Next
				}
			} else {
				var ialTokens []byte
				if nil == li.KramdownIAL {
					id := ast.NewNodeID()
					ialTokens = []byte("{: id=\"" + id + "\"}")
					li.KramdownIAL = [][]string{{"id", id}}
					li.ID = id
				} else {
					ialTokens = IAL2Tokens(li.KramdownIAL)
				}
				li.InsertAfter(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: ialTokens})
			}
		}
	}
}

var items1 = util.StrToBytes("1")

// parseListMarker 用于解析泛列表（列表、列表项或者任务列表）标记符。
func (t *Tree) parseListMarker(container *ast.Node) (data *ast.ListData, ial [][]string) {
	if 4 <= t.Context.indent {
		return nil, nil
	}

	ln := t.Context.currentLine
	tokens := ln[t.Context.nextNonspace:]
	data = &ast.ListData{
		Typ:          0,                // 默认无序列表
		Tight:        true,             // 默认紧凑模式
		MarkerOffset: t.Context.indent, // 设置前置相对缩进
		Num:          -1,               // 假设有序列表起始为 -1，后面会进行计算赋值
	}

	markerLength := 1
	marker := []byte{tokens[0]}
	var delim byte
	if lex.ItemPlus == marker[0] || lex.ItemHyphen == marker[0] || lex.ItemAsterisk == marker[0] {
		data.BulletChar = marker[0]
	} else if marker, delim = t.parseOrderedListMarker(tokens); nil != marker {
		if container.Type != ast.NodeParagraph || bytes.Equal(items1, marker) {
			data.Typ = 1 // 有序列表
			data.Start, _ = strconv.Atoi(util.BytesToStr(marker))
			markerLength = len(marker) + 1
			data.Delimiter = delim
		} else {
			return nil, nil
		}
	} else {
		return nil, nil
	}

	data.Marker = marker

	token := ln[t.Context.nextNonspace+markerLength]

	// 列表项标记符后必须是空白字符
	if !lex.IsWhitespace(token) {
		return nil, nil
	}

	// 如果要打断段落，则列表项内容部分不能为空
	if container.Type == ast.NodeParagraph && lex.ItemNewline == token {
		return nil, nil
	}

	// 到这里说明满足列表规则，开始解析并计算内部缩进空格数
	t.Context.advanceNextNonspace()             // 把起始下标移动到标记符起始位置
	t.Context.advanceOffset(markerLength, true) // 把结束下标移动到标记符结束位置
	spacesStartCol := t.Context.column
	spacesStartOffset := t.Context.offset
	for {
		t.Context.advanceOffset(1, true)
		token = lex.Peek(ln, t.Context.offset)
		if t.Context.column-spacesStartCol >= 5 || 0 == (token) || (lex.ItemSpace != token && lex.ItemTab != token) {
			break
		}
	}

	token = lex.Peek(ln, t.Context.offset)
	isBlankItem := 0 == token || lex.ItemNewline == token
	spacesAfterMarker := t.Context.column - spacesStartCol
	if spacesAfterMarker >= 5 || spacesAfterMarker < 1 || isBlankItem {
		data.Padding = markerLength + 1
		t.Context.column = spacesStartCol
		t.Context.offset = spacesStartOffset
		if token = lex.Peek(ln, t.Context.offset); lex.ItemSpace == token || lex.ItemTab == token {
			t.Context.advanceOffset(1, true)
		}
	} else {
		data.Padding = markerLength + spacesAfterMarker
	}

	if !isBlankItem {
		// 判断是否是任务列表项

		tokens := ln[t.Context.offset:]
		if t.Context.ParseOption.KramdownBlockIAL {
			if ial = t.Context.parseKramdownIALInListItem(tokens); 0 < len(ial) {
				tokens = tokens[bytes.Index(tokens, []byte("}"))+1:]
			}
		}

		if t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV {
			tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
		}

		if 3 <= len(tokens) { // 至少需要 [ ] 或者 [x] 3 个字符
			if lex.ItemOpenBracket == tokens[0] && ('x' == tokens[1] || 'X' == tokens[1] || lex.ItemSpace == tokens[1]) && lex.ItemCloseBracket == tokens[2] {
				data.Typ = 3
				data.Checked = 'x' == tokens[1] || 'X' == tokens[1]
			}
		}
	}
	return
}

func (t *Tree) parseOrderedListMarker(tokens []byte) (marker []byte, delimiter byte) {
	length := len(tokens)
	var i int
	var token byte
	for ; i < length; i++ {
		token = tokens[i]
		if !lex.IsDigit(token) || 8 < i {
			delimiter = token
			break
		}
		marker = append(marker, token)
	}

	if 1 > len(marker) || (lex.ItemDot != delimiter && lex.ItemCloseParen != delimiter) {
		return nil, 0
	}

	return
}

// endsWithBlankLine 判断块节点 block 是否是空行结束。如果 block 是列表或者列表项则迭代下降进入判断。
func endsWithBlankLine(block *ast.Node) bool {
	for nil != block {
		if block.LastLineBlank {
			return true
		}
		t := block.Type
		if !block.LastLineChecked && (t == ast.NodeList || t == ast.NodeListItem) {
			block.LastLineChecked = true
			block = block.LastChild
		} else {
			block.LastLineChecked = true
			break
		}
	}

	return false
}
