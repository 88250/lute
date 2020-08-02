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
	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
)

func (t *Tree) parseBlockRef(ctx *InlineContext) *ast.Node {
	tokens := ctx.tokens
	if 12 > len(tokens) || lex.ItemOpenParen != tokens[0] || lex.ItemOpenParen != ctx.tokens[1] {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: []byte("(")}
	}
	var id, text []byte
	savePos := ctx.pos
	ctx.pos += 2
	var ok, matched bool
	var passed, remains []byte
	for { // 这里使用 for 是为了简化逻辑，不是为了循环
		if ok, passed, remains = lex.Spnl(ctx.tokens[ctx.pos:]); !ok {
			break
		}
		ctx.pos += len(passed)
		if passed, remains, id = t.Context.parseBlockRefID(remains); 8 > len(passed) {
			break
		}
		ctx.pos += len(passed)
		matched = lex.ItemCloseParen == passed[len(passed)-1] && lex.ItemCloseParen == passed[len(passed)-2]
		if matched {
			break
		}
		if 1 > len(remains) || !lex.IsWhitespace(remains[0]) {
			break
		}
		// 跟空格的话后续尝试 title 解析
		if ok, passed, remains = lex.Spnl(remains); !ok {
			break
		}
		ctx.pos += len(passed) + 1
		matched = 2 <= len(remains) && lex.ItemCloseParen == remains[0] && lex.ItemCloseParen == remains[1]
		if matched {
			break
		}
		var validTitle bool
		if validTitle, passed, remains, text = t.Context.parseLinkTitle(remains); !validTitle {
			break
		}
		ctx.pos += len(passed)
		ok, passed, remains = lex.Spnl(remains)
		ctx.pos += len(passed)
		matched = ok && 1 < len(remains)
		if matched {
			// TODO: Vditor 内容块引用输入优化
			//if t.Context.Option.VditorWYSIWYG || t.Context.Option.VditorIR || t.Context.Option.VditorSV {
			//	if bytes.HasPrefix(remains, []byte(util.Caret+")")) {
			//		if 0 < len(title) {
			//			// 将 ‸) 换位为 )‸
			//			remains = remains[len([]byte(util.Caret+")")):]
			//			remains = append([]byte(")"+util.Caret), remains...)
			//			copy(ctx.tokens[ctx.pos-1:], remains) // 同时也将 tokens 换位，后续解析从插入符位置开始
			//		} else {
			//			// 将 ""‸ 换位为 "‸"
			//			title = util.CaretTokens
			//			remains = remains[len(util.CaretTokens):]
			//			ctx.pos += 3
			//		}
			//	} else if bytes.HasPrefix(remains, []byte(")"+util.Caret)) {
			//		if 0 == len(title) {
			//			// 将 "")‸ 换位为 "‸")
			//			title = util.CaretTokens
			//			remains = bytes.ReplaceAll(remains, util.CaretTokens, nil)
			//			ctx.pos += 3
			//		}
			//	}
			//}
			matched = lex.ItemCloseParen == remains[0] && lex.ItemCloseParen == remains[1]
		}
		break
	}
	if !matched {
		ctx.pos = savePos + 1
		return &ast.Node{Type: ast.NodeText, Tokens: []byte("(")}
	}

	ctx.pos+=2
	ret := &ast.Node{Type: ast.NodeBlockRef}
	ret.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
	ret.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
	ret.AppendChild(&ast.Node{Type: ast.NodeBlockRefID, Tokens: id})
	if 0 < len(text) {
		ret.AppendChild(&ast.Node{Type: ast.NodeBlockRefText, Tokens: text})
	}
	ret.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
	ret.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
	return ret
}

func (context *Context) parseBlockRefID(tokens []byte) (passed, remains, id []byte) {
	remains = tokens
	length := len(tokens)
	if 10 > length {
		return
	}

	passed = make([]byte, 0, 64)

	id = tokens[:8]
	for i := 0; i < 8; i++ {
		if !lex.IsASCIILetterNum(id[i]) {
			return
		}
	}

	passed = append(passed, id...)
	closed := lex.ItemCloseParen == tokens[8] && lex.ItemCloseParen == tokens[9]
	if closed {
		passed = append(passed, []byte("))")...)
		remains = tokens[10:]
		return
	}

	if lex.ItemSpace != tokens[8] && lex.ItemNewline != tokens[8] {
		passed = nil
		return
	}
	remains = tokens[8:]
	return
}
