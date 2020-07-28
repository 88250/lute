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
	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

func MathBlockContinue(mathBlock *ast.Node, context *Context) int {
	ln := context.currentLine
	indent := context.indent
	if 3 >= indent && isMathBlockClose(ln[context.nextNonspace:]) {
		context.finalize(mathBlock, context.lineNum)
		return 2
	} else {
		// 跳过 $ 之前可能存在的空格
		var i = mathBlock.MathBlockDollarOffset
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
var MathBlockMarkerCaret = util.StrToBytes("$$" + util.Caret)
var MathBlockMarkerCaretNewline = util.StrToBytes("$$" + util.Caret + "\n")

func (context *Context) mathBlockFinalize(mathBlock *ast.Node) {
	tokens := mathBlock.Tokens[2:] // 剔除开头的 $$
	tokens = lex.TrimWhitespace(tokens)
	if context.Option.VditorWYSIWYG || context.Option.VditorIR || context.Option.VditorSV {
		if bytes.HasSuffix(tokens, MathBlockMarkerCaret) {
			// 剔除结尾的 $$‸
			tokens = bytes.TrimSuffix(tokens, MathBlockMarkerCaret)
			// 把 Vditor 光标移动到内容末尾
			tokens = append(tokens, util.CaretTokens...)
		}
	}
	if bytes.HasSuffix(tokens, MathBlockMarker) {
		tokens = tokens[:len(tokens)-2] // 剔除结尾的 $$
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

func isMathBlockClose(tokens []byte) bool {
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
