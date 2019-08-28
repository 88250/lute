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

import (
	"bytes"
	"fmt"

	"github.com/b3log/gulu"
)

// parse 会将 markdown 原始文本字符数组解析为一颗语法树。
func parse(name string, markdown []byte, option options) (t *Tree, err error) {
	defer gulu.Panic.Recover(&err)

	t = &Tree{Name: name, text: markdown, context: &Context{option: option}}
	t.context.tree = t
	t.lex = lex(t.text)
	t.Root = &Document{&BaseNode{typ: NodeDocument}}
	t.parseBlocks()
	t.parseInlines()
	t.lex = nil

	return
}

// Context 用于维护解析过程中使用到的公共数据。
type Context struct {
	tree   *Tree   // 关联的语法树
	option options // 解析渲染选项

	linkRefDef map[string]*Link // 链接引用定义集

	// 以下变量用于块级解析阶段

	tip                                                      Node  // 末梢节点
	oldtip                                                   Node  // 老的末梢节点
	currentLine                                              items // 当前行
	currentLineLen                                           int   // 当前行长
	offset, column, nextNonspace, nextNonspaceColumn, indent int   // 解析时用到的下标、缩进空格数
	indented, blank, partiallyConsumedTab, allClosed         bool  // 是否是缩进行、空行等标识
	lastMatchedContainer                                     Node  // 最后一个匹配的块节点
}

// InlineContext 描述了行级解析上下文。
type InlineContext struct {
	tokens     items      // 当前解析的 tokens
	tokensLen  int        // 当前解析的 tokens 长度
	pos        int        // 当前解析到的 token 位置
	delimiters *delimiter // 分隔符栈，用于强调解析
	brackets   *delimiter // 括号栈，用于图片和链接解析
}

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
				if charsToTab > count {
					charsToAdvance = count
				} else {
					charsToAdvance = charsToTab
				}
				context.column += charsToAdvance
				if !context.partiallyConsumedTab {
					context.offset += 1
				}
				count -= charsToAdvance
			} else {
				context.partiallyConsumedTab = false
				context.column += charsToTab
				context.offset += 1
				count -= 1
			}
		} else {
			context.partiallyConsumedTab = false
			context.offset += 1
			context.column += 1 // 假定是 ASCII，因为块开始标记都是 ASCII
			count -= 1
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
		if i >= context.currentLineLen {
			fmt.Println(context.currentLine, i)
		}
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

	context.blank = itemNewline == token || itemEnd == token
	context.nextNonspace = i
	context.nextNonspaceColumn = cols
	context.indent = context.nextNonspaceColumn - context.column
	context.indented = context.indent >= 4
}

// closeUnmatchedBlocks 最终化所有未匹配的块节点。
func (context *Context) closeUnmatchedBlocks() {
	if !context.allClosed {
		for context.oldtip != context.lastMatchedContainer {
			parent := context.oldtip.Parent()
			context.finalize(context.oldtip)
			context.oldtip = parent
		}
		context.allClosed = true
	}
}

// finalize 执行 block 的最终化处理。调用该方法会将 context.tip 置为 block 的父节点。
func (context *Context) finalize(block Node) {
	var parent = block.Parent()
	block.Close()
	block.Finalize(context)
	context.tip = parent
}

// addChild 将 child 作为子节点添加到 context.tip 上。如果 tip 节点不能接受子节点（非块级容器不能添加子节点），则最终化该 tip
// 节点并向父节点方向尝试，直到找到一个能接受 child 的节点为止。
func (context *Context) addChild(child Node) {
	for !context.tip.CanContain(child.Type()) {
		context.finalize(context.tip) // 注意调用 finalize 会向父节点方向进行迭代
	}

	context.tip.AppendChild(context.tip, child)
	context.tip = child
}

// listsMatch 用户判断指定的 listData 和 itemData 是否可归属于同一个列表。
func (context *Context) listsMatch(listData, itemData *listData) bool {
	return listData.typ == itemData.typ &&
		listData.delimiter == itemData.delimiter &&
		bytes.Equal(listData.bulletChar, itemData.bulletChar)
}

// Tree 描述了 Markdown 抽象语法树结构。
type Tree struct {
	Name string    // 名称，可以为空
	Root *Document // 根节点

	text          []byte         // 原始的 Markdown 文本
	lex           *lexer         // 词法分析器
	context       *Context       // 语法解析上下文
	inlineContext *InlineContext // 行级解析上下文
}

// render 使用 renderer 对语法树 t 进行渲染，渲染结果以 output 返回。
func (t *Tree) render(renderer *Renderer) (output []byte, err error) {
	err = renderer.render(t.Root)
	if nil != err {
		return
	}
	output = renderer.writer.Bytes()

	return
}
