// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

// +build javascript

package lute

import (
	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
)

// renderCodeBlock 进行代码块 HTML 渲染，不实现语法高亮。
func (r *HTMLRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if !node.IsFencedCodeBlock {
		// 缩进代码块处理
		r.newline()
		r.writeString("<pre><code>")
		r.write(escapeHTML(node.FirstChild.Tokens))
		r.writeString("</code></pre>")
		r.newline()
		return ast.WalkStop
	}
	r.newline()
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.newline()
		tokens := node.Tokens
		if 0 < len(node.Previous.CodeBlockInfo) {
			infoWords := lex.Split(node.Previous.CodeBlockInfo, lex.ItemSpace)
			language := infoWords[0]
			r.writeString("<pre><code class=\"language-")
			r.write(language)
			r.writeString("\">")
			tokens = escapeHTML(tokens)
			r.write(tokens)
		} else {
			r.writeString("<pre><code>")
			tokens = escapeHTML(tokens)
			r.write(tokens)
		}
		return ast.WalkSkipChildren
	}
	r.writeString("</code></pre>")
	r.newline()
	return ast.WalkStop
}
