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

func (t *Tree) parseFileAnnotationRef(ctx *InlineContext) *ast.Node {
	if !t.Context.ParseOption.FileAnnotationRef {
		return nil
	}

	tokens := ctx.tokens[ctx.pos:]
	if 48 > len(tokens) || lex.ItemLess != tokens[0] || lex.ItemLess != tokens[1] {
		return nil
	}

	idPart := tokens[2:48]
	if bytes.ContainsAny(idPart, "<>") {
		return nil
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
		if passed, remains, id = t.Context.parseFileAnnotationRefID(remains); 1 > len(passed) {
			ctx.pos = savePos
			break
		}
		ctx.pos += len(passed)
		matched = lex.ItemGreater == passed[len(passed)-1] && lex.ItemGreater == passed[len(passed)-2]
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
		matched = 2 <= len(remains) && lex.ItemGreater == remains[0] && lex.ItemGreater == remains[1]
		if matched {
			ctx.pos++
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
			matched = lex.ItemGreater == remains[0] && lex.ItemGreater == remains[1]
			ctx.pos += 2
		}
		break
	}
	if !matched {
		return nil
	}

	ret := &ast.Node{Type: ast.NodeFileAnnotationRef}
	ret.AppendChild(&ast.Node{Type: ast.NodeLess})
	ret.AppendChild(&ast.Node{Type: ast.NodeLess})
	ret.AppendChild(&ast.Node{Type: ast.NodeFileAnnotationRefID, Tokens: id})
	if 0 < len(text) {
		ret.AppendChild(&ast.Node{Type: ast.NodeFileAnnotationRefSpace})
		textNode := &ast.Node{Type: ast.NodeFileAnnotationRefText, Tokens: text}
		ret.AppendChild(textNode)
	}
	ret.AppendChild(&ast.Node{Type: ast.NodeGreater})
	ret.AppendChild(&ast.Node{Type: ast.NodeGreater})
	return ret
}

func (context *Context) parseFileAnnotationRefID(tokens []byte) (passed, remains, id []byte) {
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

		if bytes.HasPrefix(tokens[i:], []byte(" \"")) {
			break
		}

		if '>' == token {
			break
		}
	}
	remains = tokens[i:]
	idPart := tokens[:i]
	if !bytes.HasPrefix(idPart, []byte("assets/")) {
		return nil, nil, nil
	}
	idPart = bytes.TrimPrefix(idPart, []byte("assets/"))
	if !bytes.Contains(idPart, []byte("/")) {
		return
	}
	idParts := bytes.Split(idPart, []byte("/"))
	if 2 != len(idParts) {
		return
	}
	filePart := idParts[0]
	if !bytes.Contains(filePart, []byte("-")) || !bytes.HasSuffix(bytes.ToLower(filePart), []byte(".pdf")) {
		return
	}
	fileName := filePart[:len(filePart)-4]
	if 23 > len(fileName) {
		return
	}
	fileID := fileName[len(fileName)-22:]
	if !ast.IsNodeIDPattern(string(fileID)) {
		return
	}
	annotationIDPart := idParts[1]
	if !ast.IsNodeIDPattern(string(annotationIDPart)) {
		return
	}

	id = tokens[:i]
	if 6 > len(remains) {
		return
	}
	passed = make([]byte, 0, 1024)
	passed = append(passed, id...)
	if bytes.HasPrefix(remains, editor.CaretTokens) {
		passed = append(passed, editor.CaretTokens...)
		remains = remains[len(editor.CaretTokens):]
	}
	closed := lex.ItemGreater == remains[0] && lex.ItemGreater == remains[1]
	if closed {
		passed = append(passed, []byte(">>")...)
		return
	}

	if !lex.IsWhitespace(remains[0]) {
		passed = nil
		return
	}
	return
}
