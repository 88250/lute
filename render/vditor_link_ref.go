// Lute - A structured Markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/lex"
)

// renderVditorLinkRefTitle keeps the definition title in the editable Markdown
// text. Escape the Markdown delimiter before escaping the surrounding HTML so
// DOM text extraction preserves quotes, entities and literal HTML in the title.
func (r *BaseRenderer) renderVditorLinkRefTitle(definition *ast.Node) {
	title := definition.FirstChild.ChildByType(ast.NodeLinkTitle)
	if nil == title {
		return
	}
	if "//" == string(definition.Tokens) && "#" == string(definition.FirstChild.ChildByType(ast.NodeLinkDest).Tokens) {
		r.WriteString(" (")
		for i, token := range title.Tokens {
			if ('(' == token || ')' == token) && !lex.IsBackslashEscapePunct(title.Tokens, i) {
				r.WriteByte(lex.ItemBackslash)
			}
			r.Write(html.EscapeHTML([]byte{token}))
		}
		r.WriteByte(')')
		return
	}

	tokens := make([]byte, 0, len(title.Tokens))
	for i, token := range title.Tokens {
		if lex.ItemDoublequote == token && !lex.IsBackslashEscapePunct(title.Tokens, i) {
			tokens = append(tokens, lex.ItemBackslash)
		}
		tokens = append(tokens, token)
	}
	r.WriteString(" \"")
	r.Write(html.EscapeHTML(tokens))
	r.WriteByte(lex.ItemDoublequote)
}
