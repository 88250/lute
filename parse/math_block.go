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
	"github.com/88250/lute/util"
)

// MathBlockStart 判断数学公式块（$$）是否开始。
func MathBlockStart(t *Tree, container *ast.Node) int {
	if t.Context.indented {
		return 0
	}

	if ok, mathBlockDollarOffset := t.parseMathBlock(); ok {
		t.Context.closeUnmatchedBlocks()
		block := t.Context.addChild(ast.NodeMathBlock)
		block.MathBlockDollarOffset = mathBlockDollarOffset
		t.Context.advanceNextNonspace()
		t.Context.advanceOffset(mathBlockDollarOffset, false)
		return 2
	}
	return 0
}

func MathBlockContinue(mathBlock *ast.Node, context *Context) int {
	ln := context.currentLine
	indent := context.indent
	if 3 >= indent && context.isMathBlockClose(ln[context.nextNonspace:]) {
		context.finalize(mathBlock)
		return 2
	} else {
		// 跳过 $ 之前可能存在的空格
		i := mathBlock.MathBlockDollarOffset
		var token byte
		for i > 0 {
			token = lex.Peek(ln, context.offset)
			if lex.ItemSpace != token && lex.ItemTab != token {
				break
			}
			context.advanceOffset(1, true)
			i--
		}
	}
	return 0
}

var MathBlockMarker = util.StrToBytes("$$")
var MathBlockMarkerNewline = util.StrToBytes("$$\n")
var MathBlockMarkerCaret = util.StrToBytes("$$" + editor.Caret)
var MathBlockMarkerCaretNewline = util.StrToBytes("$$" + editor.Caret + "\n")

func (context *Context) mathBlockFinalize(mathBlock *ast.Node) {
	if 2 > len(mathBlock.Tokens) {
		/*
			- foo

			    $$
			bar
			$$
		*/
		mathBlock.AppendChild(&ast.Node{Type: ast.NodeMathBlockOpenMarker})
		mathBlock.AppendChild(&ast.Node{Type: ast.NodeMathBlockContent})
		mathBlock.AppendChild(&ast.Node{Type: ast.NodeMathBlockCloseMarker})
		return
	}
	tokens := mathBlock.Tokens[2:] // 剔除开头的 $$
	tokens = lex.TrimWhitespace(tokens)
	if context.ParseOption.VditorWYSIWYG || context.ParseOption.VditorIR || context.ParseOption.VditorSV || context.ParseOption.ProtyleWYSIWYG {
		if bytes.HasSuffix(tokens, MathBlockMarkerCaret) {
			// 剔除结尾的 $$‸
			tokens = bytes.TrimSuffix(tokens, MathBlockMarkerCaret)
			// 把 Vditor 插入符移动到内容末尾
			tokens = append(tokens, editor.CaretTokens...)
		}
	}
	if bytes.HasSuffix(tokens, MathBlockMarker) {
		tokens = tokens[:len(tokens)-2] // 剔除结尾的 $$
	}
	if bytes.Contains(tokens, []byte("<span data-type=")) {
		// 行级元素转换为块级元素 https://ld246.com/article/1730804245164
		inlineTree := Inline("", tokens, context.ParseOption)
		if nil != inlineTree {
			tokens = []byte(inlineTree.Root.Content())
		}
	}

	mathBlock.Tokens = nil
	mathBlock.AppendChild(&ast.Node{Type: ast.NodeMathBlockOpenMarker})
	mathBlock.AppendChild(&ast.Node{Type: ast.NodeMathBlockContent, Tokens: tokens})
	mathBlock.AppendChild(&ast.Node{Type: ast.NodeMathBlockCloseMarker})
}

func (t *Tree) parseMathBlock() (ok bool, mathBlockDollarOffset int) {
	marker := t.Context.currentLine[t.Context.nextNonspace]
	if lex.ItemDollar != marker {
		return
	}

	fenceChar := marker
	fenceLength := 0
	for i := t.Context.nextNonspace; i < t.Context.currentLineLen && fenceChar == t.Context.currentLine[i]; i++ {
		fenceLength++
	}

	if 2 > fenceLength {
		return
	}
	return true, t.Context.indent
}

func (context *Context) isMathBlockClose(tokens []byte) bool {
	if context.ParseOption.KramdownBlockIAL && simpleCheckIsBlockIAL(tokens) {
		// 判断 IAL 打断
		if ial := context.parseKramdownBlockIAL(tokens); 0 < len(ial) {
			context.Tip.ID = IAL2Map(ial)["id"]
			context.Tip.KramdownIAL = ial
			context.Tip.InsertAfter(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: tokens})
			return true
		}
	}

	closeMarker := tokens[0]
	if closeMarker != lex.ItemDollar {
		return false
	}
	if 2 > lex.Accept(tokens, closeMarker) {
		return false
	}
	tokens = lex.TrimWhitespace(tokens)
	for _, token := range tokens {
		if token != lex.ItemDollar {
			return false
		}
	}
	return true
}
