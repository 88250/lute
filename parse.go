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

// parse 会将 markdown 原始文本字符数组解析为一颗语法树。
func (lute *Lute) parse(name string, markdown []byte) (tree *Tree, err error) {
	defer recoverPanic(&err)

	tree = &Tree{Name: name, context: &Context{option: lute.options}}
	tree.context.tree = tree
	tree.lexer = newLexer(markdown)
	tree.Root = &Node{typ: NodeDocument}
	tree.parseBlocks()
	tree.parseInlines()
	tree.lexer = nil

	return
}

// Context 用于维护块级元素解析过程中使用到的公共数据。
type Context struct {
	tree   *Tree    // 关联的语法树
	option *options // 解析渲染选项

	linkRefDef map[string]*Node // 链接引用定义集

	tip                                                               *Node  // 末梢节点
	oldtip                                                            *Node  // 老的末梢节点
	currentLine                                                       []byte // 当前行
	currentLineLen                                                    int    // 当前行长
	lineNum, offset, column, nextNonspace, nextNonspaceColumn, indent int    // 解析时用到的行号、下标、缩进空格数等
	indented, blank, partiallyConsumedTab, allClosed                  bool   // 是否是缩进行、空行等标识
	lastMatchedContainer                                              *Node  // 最后一个匹配的块节点
}

// InlineContext 描述了行级元素解析上下文。
type InlineContext struct {
	tokens     []byte     // 当前解析的 tokens
	tokensLen  int        // 当前解析的 tokens 长度
	pos        int        // 当前解析到的 token 位置
	lineNum    int        // 当前解析的起始行号
	columnNum  int        // 当前解析的起始列号
	delimiters *delimiter // 分隔符栈，用于强调解析
	brackets   *delimiter // 括号栈，用于图片和链接解析
}

// advanceOffset 用于移动 count 个字符位置，columns 指定了遇到 tab 时是否需要空格进行补偿偏移。
func (context *Context) advanceOffset(count int, columns bool) {
	var currentLine = context.currentLine
	var charsToTab, charsToAdvance int
	var c byte
	for 0 < count {
		c = currentLine[context.offset]
		if itemTab == c {
			charsToTab = 4 - (context.column % 4)
			if columns {
				context.partiallyConsumedTab = charsToTab > count
				if context.partiallyConsumedTab {
					charsToAdvance = count
				} else {
					charsToAdvance = charsToTab
					context.offset++
				}
				context.column += charsToAdvance
				count -= charsToAdvance
			} else {
				context.partiallyConsumedTab = false
				context.column += charsToTab
				context.offset++
				count--
			}
		} else {
			context.partiallyConsumedTab = false
			context.offset++
			context.column++ // 假定是 ASCII，因为块开始标记符都是 ASCII
			count--
		}
	}
}

// advanceNextNonspace 用于预移动到下一个非空字符位置。
func (context *Context) advanceNextNonspace() {
	context.offset = context.nextNonspace
	context.column = context.nextNonspaceColumn
	context.partiallyConsumedTab = false
}

// findNextNonspace 用于查找下一个非空字符。
func (context *Context) findNextNonspace() {
	i := context.offset
	cols := context.column

	var token byte
	for {
		token = context.currentLine[i]
		if itemSpace == token {
			i++
			cols++
		} else if itemTab == token {
			i++
			cols += 4 - (cols % 4)
		} else {
			break
		}
	}

	context.blank = itemNewline == token
	context.nextNonspace = i
	context.nextNonspaceColumn = cols
	context.indent = context.nextNonspaceColumn - context.column
	context.indented = context.indent >= 4
}

// closeUnmatchedBlocks 最终化所有未匹配的块节点。
func (context *Context) closeUnmatchedBlocks() {
	if !context.allClosed {
		for context.oldtip != context.lastMatchedContainer {
			parent := context.oldtip.parent
			context.finalize(context.oldtip, context.lineNum-1)
			context.oldtip = parent
		}
		context.allClosed = true
	}
}

// finalize 执行 block 的最终化处理。调用该方法会将 context.tip 置为 block 的父节点。
func (context *Context) finalize(block *Node, lineNum int) {
	var parent = block.parent
	block.close = true
	block.Finalize(context)
	context.tip = parent
}

// addChildMarker 将构造一个 nodeType 节点并作为子节点添加到末梢节点 context.tip 上。
func (context *Context) addChildMarker(nodeType nodeType, tokens []byte) (ret *Node) {
	ret = &Node{typ: nodeType, tokens: tokens, close: true}
	context.tip.AppendChild(ret)
	return ret
}

// addChild 将构造一个 nodeType 节点并作为子节点添加到末梢节点 context.tip 上。如果末梢不能接受子节点（非块级容器不能添加子节点），则最终化该末梢
// 节点并向父节点方向尝试，直到找到一个能接受该子节点的节点为止。添加完成后该子节点会被设置为新的末梢节点。
func (context *Context) addChild(nodeType nodeType, offset int) (ret *Node) {
	for !context.tip.CanContain(nodeType) {
		context.finalize(context.tip, context.lineNum-1) // 注意调用 finalize 会向父节点方向进行迭代
	}

	ret = &Node{typ: nodeType}
	context.tip.AppendChild(ret)
	context.tip = ret
	return ret
}

// listsMatch 用户判断指定的 listData 和 itemData 是否可归属于同一个列表。
func (context *Context) listsMatch(listData, itemData *listData) bool {
	return listData.typ == itemData.typ &&
		((0 == listData.delimiter && 0 == itemData.delimiter) || listData.delimiter == itemData.delimiter) &&
		bytes.Equal(listData.bulletChar, itemData.bulletChar)
}

// Tree 描述了 Markdown 抽象语法树结构。
type Tree struct {
	Name          string         // 名称，可以为空
	Root          *Node          // 根节点
	lexer         *lexer         // 词法分析器
	context       *Context       // 块级解析上下文
	inlineContext *InlineContext // 行级解析上下文
}
