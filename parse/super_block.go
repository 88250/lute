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

func SuperBlockContinue(superBlock *ast.Node, context *Context) int {
	ln := context.currentLine
	indent := context.indent
	if superBlock.IsFencedCodeBlock {
		if ok, closeFence := context.isFencedCodeClose(ln[context.nextNonspace:], superBlock.CodeBlockFenceChar, superBlock.CodeBlockFenceLen); indent <= 3 && ok {
			superBlock.CodeBlockCloseFence = closeFence
			context.finalize(superBlock, context.lineNum)
			return 2
		} else {
			// 跳过围栏标记符之前可能存在的空格
			i := superBlock.CodeBlockFenceOffset
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

func (context *Context) superBlockFinalize(superBlock *ast.Node) {
	if superBlock.IsFencedCodeBlock {
		content := superBlock.Tokens
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
		superBlock.Tokens = content[i+1:]
	} else { // 缩进代码块
		superBlock.Tokens = lex.ReplaceNewlineSpace(superBlock.Tokens)
	}
}

func (t *Tree) parseSuperBlock() (ok bool, layout []byte) {
	marker := t.Context.currentLine[t.Context.nextNonspace]
	if lex.ItemOpenBrace != marker {
		return
	}

	fenceChar := marker
	var fenceLen int
	for i := t.Context.nextNonspace; i < t.Context.currentLineLen && fenceChar == t.Context.currentLine[i]; i++ {
		fenceLen++
	}

	if 3 != fenceLen {
		return
	}

	layout = t.Context.currentLine[t.Context.nextNonspace+fenceLen:]
	layout = lex.TrimWhitespace(layout)
	if !bytes.EqualFold(layout, nil) && !bytes.EqualFold(layout, []byte("row")) && !bytes.EqualFold(layout, []byte("col")) {
		return
	}
	return true, layout
}

func (context *Context) isSuperBlockClose(tokens []byte, openMarker byte, num int) (ok bool, closeFence []byte) {
	if context.Option.KramdownIAL && len("{: id=\"") < len(tokens) {
		// 判断 IAL 打断
		if ial := context.parseKramdownIAL(tokens); 0 < len(ial) {
			context.Tip.KramdownIAL = ial
			context.Tip.InsertAfter(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: tokens})
			return true, context.Tip.CodeBlockOpenFence
		}
	}

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
		if endCaret {
			context.Tip.Tokens = bytes.TrimSuffix(context.Tip.Tokens, []byte("\n"))
			context.Tip.Tokens = append(context.Tip.Tokens, util.CaretTokens...)
		}
	}
	for _, token := range tokens {
		if token != openMarker {
			return false, nil
		}
	}
	closeFence = tokens
	return true, closeFence
}
