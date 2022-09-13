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

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"

	"github.com/88250/lute/lex"
)

// ATXHeadingStart 判断 ATX 标题（#）是否开始。
func ATXHeadingStart(t *Tree, container *ast.Node) int {
	if t.Context.indented {
		return 0
	}

	if ok, markers, content, level := t.parseATXHeading(); ok {
		t.Context.advanceNextNonspace()
		t.Context.advanceOffset(len(content), false)
		t.Context.closeUnmatchedBlocks()
		heading := t.Context.addChild(ast.NodeHeading)
		heading.HeadingLevel = level
		heading.Tokens = content
		crosshatchMarker := &ast.Node{Type: ast.NodeHeadingC8hMarker, Tokens: markers}
		heading.AppendChild(crosshatchMarker)
		t.Context.advanceOffset(t.Context.currentLineLen-t.Context.offset, false)
		return 2
	}
	return 0
}

// SetextHeadingStart 判断 Setext 标题（- =）是否开始。
func SetextHeadingStart(t *Tree, container *ast.Node) int {
	if t.Context.indented || ast.NodeParagraph != container.Type || !t.Context.ParseOption.Setext {
		return 0
	}

	level := t.parseSetextHeading()
	if 0 == level {
		return 0
	}

	if t.Context.ParseOption.GFMTable {
		// 尝试解析表，因为可能出现如下情况：
		//
		//   0
		//   -:
		//   -
		//
		// 前两行可以解析出一个只有一个单元格的表。
		// Empty list following GFM Table makes table broken https://github.com/b3log/lute/issues/9
		table := t.Context.parseTable0(container.Tokens)
		if nil != table {
			// 将该段落节点转成表节点
			container.Type = ast.NodeTable
			container.TableAligns = table.TableAligns
			for tr := table.FirstChild; nil != tr; {
				nextTr := tr.Next
				container.AppendChild(tr)
				tr = nextTr
			}
			container.Tokens = nil
			return 0
		}
	}

	t.Context.closeUnmatchedBlocks()
	// 解析链接引用定义
	for tokens := container.Tokens; 0 < len(tokens) && lex.ItemOpenBracket == tokens[0]; tokens = container.Tokens {
		if remains := t.Context.parseLinkRefDef(tokens); nil != remains {
			container.Tokens = remains
		} else {
			break
		}
	}

	if 0 < len(container.Tokens) {
		child := &ast.Node{Type: ast.NodeHeading, HeadingLevel: level, HeadingSetext: true}
		child.Tokens = lex.TrimWhitespace(container.Tokens)
		container.InsertAfter(child)
		container.Unlink()
		t.Context.Tip = child
		t.Context.advanceOffset(t.Context.currentLineLen-t.Context.offset, false)
		return 2
	}
	return 0
}

func (t *Tree) parseATXHeading() (ok bool, markers, content []byte, level int) {
	tokens := t.Context.currentLine[t.Context.nextNonspace:]
	var startCaret bool
	if (t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.ProtyleWYSIWYG) && bytes.HasPrefix(tokens, editor.CaretTokens) {
		tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
		startCaret = true
	}

	marker := tokens[0]
	if lex.ItemCrosshatch != marker {
		return
	}

	var inCaret bool
	if (t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV) && bytes.Contains(tokens, []byte("#"+editor.Caret+"#")) {
		tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
		inCaret = true
	}

	level = lex.Accept(tokens, lex.ItemCrosshatch)
	if 6 < level {
		return
	}

	var endCaret bool
	if (t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV) && bytes.HasPrefix(tokens[level:], []byte(" "+editor.Caret)) {
		tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
		endCaret = true
	}

	if level < len(tokens) && !lex.IsWhitespace(tokens[level]) {
		return
	}

	markers = t.Context.currentLine[t.Context.nextNonspace : t.Context.nextNonspace+level+1]

	content = make([]byte, 0, 256)
	_, tokens = lex.TrimLeft(tokens)
	_, tokens = lex.TrimLeft(tokens[level:])
	for _, token := range tokens {
		if lex.ItemNewline == token {
			break
		}
		content = append(content, token)
	}
	_, content = lex.TrimRight(content)
	closingCrosshatchIndex := len(content) - 1
	for ; 0 <= closingCrosshatchIndex; closingCrosshatchIndex-- {
		if lex.ItemCrosshatch == content[closingCrosshatchIndex] {
			continue
		}
		if lex.ItemSpace == content[closingCrosshatchIndex] {
			break
		} else {
			closingCrosshatchIndex = len(content)
			break
		}
	}

	if 0 >= closingCrosshatchIndex {
		content = make([]byte, 0, 0)
	} else if 0 < closingCrosshatchIndex {
		content = content[:closingCrosshatchIndex]
		_, content = lex.TrimRight(content)
	}

	if t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.ProtyleWYSIWYG {
		if startCaret || inCaret || endCaret {
			content = append(editor.CaretTokens, content...)
		}
	}
	_, content = lex.TrimRight(content)
	ok = true
	return
}

func (t *Tree) parseSetextHeading() (level int) {
	ln := lex.TrimWhitespace(t.Context.currentLine)
	var caretInLn bool
	if t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.ProtyleWYSIWYG {
		if bytes.Contains(ln, editor.CaretTokens) {
			caretInLn = true
			ln = bytes.ReplaceAll(ln, editor.CaretTokens, nil)
			if 1 > len(ln) {
				return
			}
		}
	}

	start := 0
	marker := ln[start]
	if lex.ItemEqual != marker && lex.ItemHyphen != marker {
		return
	}

	length := len(ln)
	for ; start < length; start++ {
		token := ln[start]
		if lex.ItemEqual != token && lex.ItemHyphen != token {
			return
		}

		if 0 != marker {
			if marker != token {
				return
			}
		} else {
			marker = token
		}
	}

	level = 1
	if lex.ItemHyphen == marker {
		level = 2
	}

	if (t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.ProtyleWYSIWYG) && caretInLn {
		t.Context.oldtip.Tokens = lex.TrimWhitespace(t.Context.oldtip.Tokens)
		t.Context.oldtip.AppendTokens(editor.CaretTokens)
	}
	return
}
