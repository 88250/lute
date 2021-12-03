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
	"unicode/utf8"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

func (context *Context) parseLinkRefDef(tokens []byte) []byte {
	if !context.ParseOption.LinkRef {
		return nil
	}

	_, tokens = lex.TrimLeft(tokens)
	if 1 > len(tokens) {
		return nil
	}

	n, remains, label := context.parseLinkLabel(tokens)
	if 2 > n || 1 > len(label) {
		return nil
	}

	length := len(remains)
	if 1 > length {
		return nil
	}

	if ':' != remains[0] {
		return nil
	}

	remains = remains[1:]
	whitespaces, remains := lex.TrimLeft(remains)
	newlines, _, _ := lex.StatWhitespace(whitespaces)
	if 1 < newlines {
		return nil
	}

	tokens = remains
	linkDest, remains, destination := context.parseLinkDest(tokens)
	if nil == linkDest {
		return nil
	}

	whitespaces, remains = lex.TrimLeft(remains)
	if nil == whitespaces && 0 < len(remains) {
		return nil
	}
	newlines, spaces1, tabs1 := lex.StatWhitespace(whitespaces)
	if 1 < newlines {
		return nil
	}

	_, tokens = lex.TrimLeft(remains)
	validTitle, _, remains, title := context.parseLinkTitle(tokens)
	if !validTitle && 1 > newlines {
		return nil
	}
	if 0 < spaces1+tabs1 && !lex.IsBlankLine(remains) && lex.ItemNewline != remains[0] {
		return nil
	}

	titleLine := tokens
	whitespaces, tokens = lex.TrimLeft(remains)
	_, spaces2, tabs2 := lex.StatWhitespace(whitespaces)
	if !lex.IsBlankLine(tokens) && 0 < spaces2+tabs2 {
		remains = titleLine
	} else {
		remains = tokens
	}

	link := context.Tree.newLink(ast.NodeLink, label, destination, title, 1)
	def := &ast.Node{Type: ast.NodeLinkRefDef, Tokens: label}
	def.AppendChild(link)
	defBlock := context.Tip
	if ast.NodeLinkRefDefBlock != defBlock.Type {
		defBlock = &ast.Node{Type: ast.NodeLinkRefDefBlock}
	}
	defBlock.AppendChild(def)
	context.Tip.Parent.AppendChild(defBlock)
	return remains
}

func (context *Context) parseLinkTitle(tokens []byte) (validTitle bool, passed, remains, title []byte) {
	if 1 > len(tokens) {
		return true, nil, tokens, nil
	}
	if lex.ItemOpenBracket == tokens[0] {
		return true, nil, tokens, nil
	}

	validTitle, passed, remains, title = context.parseLinkTitleMatch(lex.ItemDoublequote, lex.ItemDoublequote, tokens)
	if !validTitle {
		validTitle, passed, remains, title = context.parseLinkTitleMatch(lex.ItemSinglequote, lex.ItemSinglequote, tokens)
		if !validTitle {
			validTitle, passed, remains, title = context.parseLinkTitleMatch(lex.ItemOpenParen, lex.ItemCloseParen, tokens)
		}
	}
	if nil != title {
		if !context.ParseOption.VditorWYSIWYG && !context.ParseOption.VditorIR && !context.ParseOption.VditorSV && !context.ParseOption.ProtyleWYSIWYG {
			title = html.UnescapeBytes(title)
		}
	}
	return
}

func (context *Context) parseBlockRefText(tokens []byte) (validTitle bool, passed, remains, title []byte, subtype string) {
	if 1 > len(tokens) {
		return true, nil, tokens, nil, ""
	}
	if lex.ItemOpenBracket == tokens[0] {
		return true, nil, tokens, nil, ""
	}

	validTitle, passed, remains, title = context.parseLinkTitleMatch(lex.ItemDoublequote, lex.ItemDoublequote, tokens)
	subtype = "s"
	if !validTitle {
		validTitle, passed, remains, title = context.parseLinkTitleMatch(lex.ItemSinglequote, lex.ItemSinglequote, tokens)
		subtype = "d"
	}
	if nil != title {
		if !context.ParseOption.VditorWYSIWYG && !context.ParseOption.VditorIR && !context.ParseOption.VditorSV && !context.ParseOption.ProtyleWYSIWYG {
			title = html.UnescapeBytes(title)
		}
	}
	return
}

func (context *Context) parseLinkTitleMatch(opener, closer byte, tokens []byte) (validTitle bool, passed, remains, title []byte) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if opener != tokens[0] {
		return
	}

	line := tokens
	length = len(line)
	closed := false
	i := 1
	size := 0
	var r rune
	for ; i < length; i += size {
		token := line[i]
		passed = append(passed, token)
		r, size = utf8.DecodeRune(line[i:])
		for j := 1; j < size; j++ {
			passed = append(passed, tokens[i+j])
		}
		title = append(title, util.StrToBytes(string(r))...)
		if closer == token && !lex.IsBackslashEscapePunct(tokens, i) {
			closed = true
			title = title[:len(title)-1]
			break
		}
	}

	if !closed {
		passed = nil
		return
	}

	validTitle = true
	remains = tokens[i+1:]
	return
}

func (context *Context) parseLinkDest(tokens []byte) (ret, remains, destination []byte) {
	ret, remains, destination = context.parseLinkDest1(tokens) // <autolink>
	if nil == ret {
		ret, remains, destination = context.parseLinkDest2(tokens) // [label](/url)
	}
	if nil != ret {
		if !context.ParseOption.VditorWYSIWYG && !context.ParseOption.VditorIR && !context.ParseOption.VditorSV && !context.ParseOption.ProtyleWYSIWYG {
			destination = html.EncodeDestination(html.UnescapeBytes(destination))
		}
	}
	return
}

func (context *Context) parseLinkDest2(tokens []byte) (ret, remains, destination []byte) {
	remains = tokens
	length := len(tokens)
	if 1 > length {
		return
	}

	ret = make([]byte, 0, 256)
	destination = make([]byte, 0, 256)

	var openParens int
	i := 0
	size := 0
	var r rune
	for i < length {
		token := tokens[i]
		ret = append(ret, token)
		r, size = utf8.DecodeRune(tokens[i:])
		for j := 1; j < size; j++ {
			ret = append(ret, tokens[i+j])
		}
		destination = append(destination, util.StrToBytes(string(r))...)
		if lex.IsWhitespace(token) || lex.IsControl(token) {
			destination = destination[:len(destination)-1]
			ret = ret[:len(ret)-1]
			break
		}

		if lex.ItemOpenParen == token && !lex.IsBackslashEscapePunct(tokens, i) {
			openParens++
		}
		if lex.ItemCloseParen == token && !lex.IsBackslashEscapePunct(tokens, i) {
			openParens--
			if 1 > openParens {
				i++
				break
			}
		}

		i += size
	}

	remains = tokens[i:]
	if length > i && !lex.IsWhitespace(tokens[i]) {
		ret = nil
		return
	}

	return
}

func (context *Context) parseLinkDest1(tokens []byte) (ret, remains, destination []byte) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if lex.ItemLess != tokens[0] {
		return
	}

	ret = make([]byte, 0, 256)
	destination = make([]byte, 0, 256)

	closed := false
	i := 0
	size := 0
	var r rune
	for ; i < length; i += size {
		token := tokens[i]
		ret = append(ret, token)
		size = 1
		if 0 < i {
			r, size = utf8.DecodeRune(tokens[i:])
			for j := 1; j < size; j++ {
				ret = append(ret, tokens[i+j])
			}
			destination = append(destination, util.StrToBytes(string(r))...)
			if lex.ItemLess == token && !lex.IsBackslashEscapePunct(tokens, i) {
				ret = nil
				return
			}
		}

		if lex.ItemGreater == token && !lex.IsBackslashEscapePunct(tokens, i) {
			closed = true
			destination = destination[0 : len(destination)-1]
			break
		}
	}

	if !closed {
		ret = nil
		return
	}

	remains = tokens[i+1:]

	return
}

func (context *Context) parseLinkLabel(tokens []byte) (n int, remains, label []byte) {
	length := len(tokens)
	if 2 > length {
		return
	}

	if lex.ItemOpenBracket != tokens[0] {
		return
	}

	passed := make([]byte, 0, len(tokens))
	passed = append(passed, tokens[0])

	closed := false
	i := 1
	for i < length {
		token := tokens[i]
		passed = append(passed, token)
		r, size := utf8.DecodeRune(tokens[i:])
		for j := 1; j < size; j++ {
			passed = append(passed, tokens[i+j])
		}
		label = append(label, util.StrToBytes(string(r))...)
		if lex.ItemCloseBracket == token && !lex.IsBackslashEscapePunct(tokens, i) {
			closed = true
			label = label[0 : len(label)-1]
			remains = tokens[i+1:]
			break
		}
		if lex.ItemOpenBracket == token && !lex.IsBackslashEscapePunct(tokens, i) {
			passed = nil
			return
		}
		i += size
	}

	if !closed || nil == lex.TrimWhitespace(label) || 999 < len(label) {
		passed = nil
		return
	}

	label = lex.TrimWhitespace(label)
	if !context.ParseOption.VditorWYSIWYG && !context.ParseOption.VditorIR && !context.ParseOption.VditorSV && !context.ParseOption.ProtyleWYSIWYG {
		label = lex.ReplaceAll(label, lex.ItemNewline, lex.ItemSpace)
		length := len(label)
		var token byte
		for i := 0; i < length; i++ {
			token = label[i]
			if token == lex.ItemSpace && i < length-1 && label[i+1] == lex.ItemSpace {
				label = append(label[:i], label[i+1:]...)
				length--
			}
		}
	}
	n = len(passed)
	return
}
