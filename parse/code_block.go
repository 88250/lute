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
	"github.com/88250/lute/html"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

func CodeBlockContinue(codeBlock *ast.Node, context *Context) int {
	ln := context.currentLine
	indent := context.indent
	if codeBlock.IsFencedCodeBlock {
		if ok, closeFence := context.isFencedCodeClose(ln[context.nextNonspace:], codeBlock.CodeBlockFenceChar, codeBlock.CodeBlockFenceLen); indent <= 3 && ok {
			codeBlock.CodeBlockCloseFence = closeFence
			context.finalize(codeBlock, context.lineNum)
			return 2
		} else {
			// 跳过围栏标记符之前可能存在的空格
			i := codeBlock.CodeBlockFenceOffset
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
	} else { // 缩进代码块
		if indent >= 4 {
			context.advanceOffset(4, true)
		} else if context.blank {
			context.advanceNextNonspace()
		} else {
			return 1
		}
	}
	return 0
}

func codeBlockFinalize(codeBlock *ast.Node) {
	if codeBlock.IsFencedCodeBlock {
		content := codeBlock.Tokens
		length := len(content)
		if 1 > length {
			return
		}

		var i int
		for ; i < length; i++ {
			if lex.ItemNewline == content[i] {
				break
			}
		}
		codeBlock.Tokens = content[i+1:]
	} else { // 缩进代码块
		codeBlock.Tokens = lex.ReplaceNewlineSpace(codeBlock.Tokens)
	}
}

var codeBlockBacktick = util.StrToBytes("`")

func (t *Tree) parseFencedCode() (ok bool, fenceChar byte, fenceLen int, fenceOffset int, openFence, codeBlockInfo []byte) {
	marker := t.Context.currentLine[t.Context.nextNonspace]
	if lex.ItemBacktick != marker && lex.ItemTilde != marker {
		return
	}

	fenceChar = marker
	for i := t.Context.nextNonspace; i < t.Context.currentLineLen && fenceChar == t.Context.currentLine[i]; i++ {
		fenceLen++
	}

	if 3 > fenceLen {
		return
	}

	openFence = t.Context.currentLine[t.Context.nextNonspace : t.Context.nextNonspace+fenceLen]

	var info []byte
	infoTokens := t.Context.currentLine[t.Context.nextNonspace+fenceLen:]
	if lex.ItemBacktick == marker && bytes.Contains(infoTokens, codeBlockBacktick) {
		// info 部分不能包含 `
		return
	}
	info = lex.TrimWhitespace(infoTokens)
	info = html.UnescapeBytes(info)
	return true, fenceChar, fenceLen, t.Context.indent, openFence, info
}

func (context *Context) isFencedCodeClose(tokens []byte, openMarker byte, num int) (ok bool, closeFence []byte) {
	closeMarker := tokens[0]
	if closeMarker != openMarker {
		return false, nil
	}
	if num > lex.Accept(tokens, closeMarker) {
		return false, nil
	}
	tokens = lex.TrimWhitespace(tokens)
	endCaret := bytes.HasSuffix(tokens, util.CaretTokens)
	if context.Option.VditorWYSIWYG || context.Option.VditorIR || context.Option.VditorSV {
		tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
	}
	if context.Option.VditorSV && endCaret {
		context.Tip.Tokens = bytes.TrimSuffix(context.Tip.Tokens, []byte("\n"))
		context.Tip.Tokens = append(context.Tip.Tokens, util.CaretTokens...)
	}
	for _, token := range tokens {
		if token != openMarker {
			return false, nil
		}
	}
	closeFence = tokens
	return true, closeFence
}
