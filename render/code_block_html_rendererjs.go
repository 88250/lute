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

package render

import (
	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

// renderCodeBlock 进行代码块 HTML 渲染，不实现语法高亮。
func (r *HtmlRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if !node.IsFencedCodeBlock {
		// 缩进代码块处理
		r.Newline()
		r.WriteString("<pre><code>")
		r.Write(util.EscapeHTML(node.FirstChild.Tokens))
		r.WriteString("</code></pre>")
		r.Newline()
		return ast.WalkStop
	}
	r.Newline()
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		tokens := node.Tokens
		if 0 < len(node.Previous.CodeBlockInfo) {
			infoWords := lex.Split(node.Previous.CodeBlockInfo, lex.ItemSpace)
			language := infoWords[0]
			r.WriteString("<pre><code class=\"language-")
			r.Write(language)
			r.WriteString("\">")
			tokens = util.EscapeHTML(tokens)
			r.Write(tokens)
		} else {
			r.WriteString("<pre><code>")
			tokens = util.EscapeHTML(tokens)
			r.Write(tokens)
		}
		return ast.WalkSkipChildren
	}
	r.WriteString("</code></pre>")
	r.Newline()
	return ast.WalkStop
}
