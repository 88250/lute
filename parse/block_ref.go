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

func (t *Tree) parseBlockRef(ctx *InlineContext) *ast.Node {
	if !t.Context.ParseOption.BlockRef {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: []byte("(")}
	}

	tokens := ctx.tokens[ctx.pos:]
	if 5 > len(tokens) || lex.ItemOpenParen != tokens[0] || lex.ItemOpenParen != tokens[1] {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: []byte("(")}
	}

	var id, text []byte
	var subtype string
	savePos := ctx.pos
	ctx.pos += 2
	var ok, matched bool
	var passed, remains []byte
	for { // 这里使用 for 是为了简化逻辑，不是为了循环
		if ok, passed, remains = lex.Spnl(ctx.tokens[ctx.pos:]); !ok {
			break
		}
		ctx.pos += len(passed)
		if passed, remains, id = t.Context.parseBlockRefID(remains); 1 > len(passed) {
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
		// 跟空格的话后续尝试锚文本解析
		if ok, passed, remains = lex.Spnl(remains); !ok {
			break
		}
		ctx.pos += len(passed) + 1
		matched = 2 <= len(remains) && lex.ItemCloseParen == remains[0] && lex.ItemCloseParen == remains[1]
		if matched {
			ctx.pos++
			break
		}
		var validTitle bool
		if validTitle, passed, remains, text, subtype = t.Context.parseBlockRefText(remains); !validTitle {
			break
		}
		ctx.pos += len(passed)
		ok, passed, remains = lex.Spnl(remains)
		ctx.pos += len(passed)
		matched = ok && 1 < len(remains)
		if matched {
			matched = lex.ItemCloseParen == remains[0] && lex.ItemCloseParen == remains[1]
			ctx.pos += 2
		}
		break
	}
	if !matched {
		ctx.pos = savePos + 1
		return &ast.Node{Type: ast.NodeText, Tokens: []byte("(")}
	}

	ret := &ast.Node{Type: ast.NodeBlockRef}
	ret.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
	ret.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
	ret.AppendChild(&ast.Node{Type: ast.NodeBlockRefID, Tokens: id})
	if 0 < len(text) {
		ret.AppendChild(&ast.Node{Type: ast.NodeBlockRefSpace})
		textNode := &ast.Node{Type: ast.NodeBlockRefText, Tokens: text}
		if "d" == subtype {
			textNode.Type = ast.NodeBlockRefDynamicText
		}
		ret.AppendChild(textNode)
	}
	ret.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
	ret.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
	return ret
}

func (context *Context) parseBlockRefID(tokens []byte) (passed, remains, id []byte) {
	remains = tokens
	length := len(tokens)
	if 1 > length {
		return
	}

	var i int
	var token byte
	for ; i < length; i++ {
		token = tokens[i]
		if bytes.Contains(editor.CaretTokens, []byte{token}) {
			continue
		}

		if lex.IsWhitespace(token) || ')' == token || !lex.IsASCIILetterNumHyphen(tokens[i]) {
			break
		}
	}
	remains = tokens[i:]
	id = tokens[:i]
	if 2 > len(remains) || !ast.IsNodeIDPattern(string(id)) {
		return
	}
	passed = make([]byte, 0, 64)
	passed = append(passed, id...)
	if bytes.HasPrefix(remains, editor.CaretTokens) {
		passed = append(passed, editor.CaretTokens...)
		remains = remains[len(editor.CaretTokens):]
	}
	closed := lex.ItemCloseParen == remains[0] && lex.ItemCloseParen == remains[1]
	if closed {
		passed = append(passed, []byte("))")...)
		return
	}

	if !lex.IsWhitespace(remains[0]) {
		passed = nil
		return
	}
	return
}
