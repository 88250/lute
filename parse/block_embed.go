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
)

func (t *Tree) parseBlockEmbed() (ret *ast.Node) {
	tokens := t.Context.currentLine[t.Context.nextNonspace:]
	if 6 > t.Context.currentLineLen || !bytes.HasPrefix(tokens, []byte("!((")) {
		return
	}

	tokens = tokens[3:]
	var passed, remains, id, text []byte
	var pos int
	var ok bool
	for { // 这里使用 for 是为了简化逻辑，不是为了循环
		if ok, passed, remains = lex.Spnl(tokens[pos:]); !ok {
			break
		}
		pos += len(passed)
		if passed, remains, id = t.Context.parseBlockRefID(remains); 1 > len(passed) {
			break
		}
		pos += len(passed)
		ok = lex.ItemCloseParen == passed[len(passed)-1] && lex.ItemCloseParen == passed[len(passed)-2]
		if ok {
			break
		}
		if 1 > len(remains) || !lex.IsWhitespace(remains[0]) {
			break
		}
		// 跟空格的话后续尝试 title 解析
		if ok, passed, remains = lex.Spnl(remains); !ok {
			break
		}
		pos += len(passed) + 1
		ok = 2 <= len(remains) && lex.ItemCloseParen == remains[0] && lex.ItemCloseParen == remains[1]
		if ok {
			pos++
			break
		}
		var validTitle bool
		if validTitle, passed, remains, text = t.Context.parseLinkTitle(remains); !validTitle {
			break
		}
		pos += len(passed)
		ok, passed, remains = lex.Spnl(remains)
		pos += len(passed)
		ok = ok && 1 < len(remains)
		if ok {
			ok = lex.ItemCloseParen == remains[0] && lex.ItemCloseParen == remains[1]
			pos += 2
		}
		break
	}
	if !ok {
		return
	}

	ret = &ast.Node{Type: ast.NodeBlockEmbed}
	ret.AppendChild(&ast.Node{Type: ast.NodeBang})
	ret.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
	ret.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
	ret.AppendChild(&ast.Node{Type: ast.NodeBlockEmbedID, Tokens: id})
	if 1 > len(text) {
		text = id
	}
	ret.AppendChild(&ast.Node{Type: ast.NodeBlockEmbedSpace})
	ret.AppendChild(&ast.Node{Type: ast.NodeBlockEmbedText, Tokens: text})
	ret.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
	ret.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
	return
}
