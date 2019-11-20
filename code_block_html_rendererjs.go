// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

// +build javascript

package lute

// renderCodeBlock 进行代码块 HTML 渲染，不实现语法高亮。
func (r *HTMLRenderer) renderCodeBlock(node *Node, entering bool) (WalkStatus, error) {
	if !node.isFencedCodeBlock {
		// 缩进代码块处理
		r.newline()
		r.writeString("<pre><code>")
		r.write(escapeHTML(node.tokens))
		r.writeString("</code></pre>")
		r.newline()
		return WalkStop, nil
	}
	r.newline()
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderCodeBlockCode(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		tokens := node.tokens
		if 0 < len(node.previous.codeBlockInfo) {
			infoWords := split(node.previous.codeBlockInfo, itemSpace)
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
		return WalkSkipChildren, nil
	}
	r.writeString("</code></pre>")
	r.newline()
	return WalkStop, nil
}
