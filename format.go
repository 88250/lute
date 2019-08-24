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

package lute

import (
	"strings"
)

// newFormatRenderer 创建一个格式化渲染器。
func newFormatRenderer(option options) (ret *Renderer) {
	ret = &Renderer{rendererFuncs: map[int]RendererFunc{}, option: option}

	// 注册 CommonMark 渲染函数

	ret.rendererFuncs[NodeDocument] = ret.renderDocumentMarkdown
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraphMarkdown
	ret.rendererFuncs[NodeText] = ret.renderTextMarkdown
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpanMarkdown
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlockMarkdown
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasisMarkdown
	ret.rendererFuncs[NodeStrong] = ret.renderStrongMarkdown
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquoteMarkdown
	ret.rendererFuncs[NodeHeading] = ret.renderHeadingMarkdown
	ret.rendererFuncs[NodeList] = ret.renderListMarkdown
	ret.rendererFuncs[NodeListItem] = ret.renderListItemMarkdown
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreakMarkdown
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreakMarkdown
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreakMarkdown
	ret.rendererFuncs[NodeHTMLBlock] = ret.renderHTMLMarkdown
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTMLMarkdown
	ret.rendererFuncs[NodeLink] = ret.renderLinkMarkdown
	ret.rendererFuncs[NodeImage] = ret.renderImageMarkdown

	// 注册 GFM 渲染函数

	ret.rendererFuncs[NodeStrikethrough] = ret.renderStrikethroughMarkdown
	ret.rendererFuncs[NodeTaskListItemMarker] = ret.renderTaskListItemMarkerMarkdown
	ret.rendererFuncs[NodeTable] = ret.renderTableMarkdown
	ret.rendererFuncs[NodeTableHead] = ret.renderTableHeadMarkdown
	ret.rendererFuncs[NodeTableRow] = ret.renderTableRowMarkdown
	ret.rendererFuncs[NodeTableCell] = ret.renderTableCellMarkdown

	return
}

func (r *Renderer) renderTableCellMarkdown(node Node, entering bool) (WalkStatus, error) {
	tag := "td"
	if NodeTableHead == node.Parent().Type() {
		tag = "th"
	}
	if entering {
		cell := node.(*TableCell)
		var attrs [][]string
		switch cell.Aligns {
		case 1:
			attrs = append(attrs, []string{"align", "left"})
		case 2:
			attrs = append(attrs, []string{"align", "center"})
		case 3:
			attrs = append(attrs, []string{"align", "right"})
		}
		r.tag(tag, attrs, false)
	} else {
		r.tag("/"+tag, nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableRowMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("tr", nil, false)
		r.newline()
	} else {
		r.tag("/tr", nil, false)
		r.newline()
		if node == node.Parent().LastChild() {
			r.tag("/tbody", nil, false)
		}
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableHeadMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("thead", nil, false)
		r.newline()
		r.tag("tr", nil, false)
		r.newline()
	} else {
		r.tag("/tr", nil, false)
		r.newline()
		r.tag("/thead", nil, false)
		r.newline()
		if nil != node.Next() {
			r.tag("tbody", nil, false)
		}
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("table", nil, false)
		r.newline()
	} else {
		r.tag("/table", nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrikethroughMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("del", nil, false)
	} else {
		r.tag("/del", nil, false)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderImageMarkdown(node Node, entering bool) (WalkStatus, error) {
	n := node.(*Image)
	if entering {
		if 0 == r.disableTags {
			r.writeString("<img src=\"")
			r.write(escapeHTML(toItems(n.Destination)))
			r.writeString("\" alt=\"")
		}
		r.disableTags++
		return WalkContinue, nil
	}

	r.disableTags--
	if 0 == r.disableTags {
		r.writeString("\"")
		if "" != n.Title {
			r.writeString(" title=\"")
			r.write(escapeHTML(toItems(n.Title)))
			r.writeString("\"")
		}
		r.writeString(" />")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderLinkMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		n := node.(*Link)
		attrs := [][]string{{"href", fromItems(escapeHTML(toItems(n.Destination)))}}
		if "" != n.Title {
			attrs = append(attrs, []string{"title", fromItems(escapeHTML(toItems(n.Title)))})
		}
		r.tag("a", attrs, false)

		return WalkContinue, nil
	}

	r.tag("/a", nil, false)

	return WalkContinue, nil
}

func (r *Renderer) renderHTMLMarkdown(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.newline()
	r.write(node.Tokens())
	r.newline()

	return WalkContinue, nil
}

func (r *Renderer) renderInlineHTMLMarkdown(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.write(node.Tokens())

	return WalkContinue, nil
}

func (r *Renderer) renderDocumentMarkdown(node Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *Renderer) renderParagraphMarkdown(node Node, entering bool) (WalkStatus, error) {
	if grandparent := node.Parent().Parent(); nil != grandparent {
		if list, ok := grandparent.(*List); ok { // List.ListItem.Paragraph
			if list.tight {
				return WalkContinue, nil
			}
		}
	}

	if entering {
		r.newline()
	} else {
		r.newline()
	}

	return WalkContinue, nil
}

func (r *Renderer) renderTextMarkdown(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.write(escapeHTML(node.Tokens()))

	return WalkContinue, nil
}

func (r *Renderer) renderCodeSpanMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<code>")
		r.write(escapeHTML(node.Tokens()))
		return WalkSkipChildren, nil
	}
	r.writeString("</code>")
	return WalkContinue, nil
}

func (r *Renderer) renderCodeBlockMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		n := node.(*CodeBlock)
		tokens := n.tokens
		if "" != n.info {
			infoWords := strings.Fields(n.info)
			language := infoWords[0]
			r.writeString("<pre><code class=\"language-" + language + "\">")
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
	return WalkContinue, nil
}

func (r *Renderer) renderEmphasisMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("em", nil, false)
	} else {
		r.tag("/em", nil, false)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrongMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<strong>")
		r.write(node.Tokens())
	} else {
		r.writeString("</strong>")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderBlockquoteMarkdown(n Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.writeString("<blockquote>")
		r.newline()
	} else {
		r.newline()
		r.writeString("</blockquote>")
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderHeadingMarkdown(node Node, entering bool) (WalkStatus, error) {
	n := node.(*Heading)
	if entering {
		r.newline()
		r.writeString(strings.Repeat("#", n.Level) + " ") // 统一使用 ATX 标题，不使用 Setext 标题
	} else {
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
	} else {
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListItemMarkdown(node Node, entering bool) (WalkStatus, error) {
	n := node.(*ListItem)
	if entering {
		r.writeString(strings.Repeat(" ", n.margin))
		r.write(n.marker)
		if 1 == n.listData.typ {
			r.write(n.bulletChar)
			r.writeByte(n.delimiter)
		}
		r.writeByte(' ')
	} else {
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTaskListItemMarkerMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		n := node.(*TaskListItemMarker)
		var attrs [][]string
		if n.checked {
			attrs = append(attrs, []string{"checked", ""})
		}
		attrs = append(attrs, []string{"disabled", ""}, []string{"type", "checkbox"})
		r.tag("input", attrs, true)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderThematicBreakMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.tag("hr", nil, true)
		r.newline()
	}

	return WalkContinue, nil
}

func (r *Renderer) renderHardBreakMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("br", nil, true)
		r.newline()
	}

	return WalkContinue, nil
}

func (r *Renderer) renderSoftBreakMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		if r.option.SoftBreak2HardBreak {
			r.tag("br", nil, true)
			r.newline()
		} else {
			r.newline()
		}
	}

	return WalkContinue, nil
}
