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

// Vditor DOM Renderer

import (
	"bytes"
	"strconv"
)

// newVditorRenderer 创建一个 Vditor DOM 渲染器。
func newVditorRenderer(option *options) (ret *Renderer) {
	ret = &Renderer{rendererFuncs: map[int]RendererFunc{}, option: option}

	// 注册 CommonMark 渲染函数

	ret.rendererFuncs[NodeDocument] = ret.renderDocumentVditor
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraphVditor
	ret.rendererFuncs[NodeText] = ret.renderTextVditor
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpanVditor
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlockVditor
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasisVditor
	ret.rendererFuncs[NodeStrong] = ret.renderStrongVditor
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquoteVditor
	ret.rendererFuncs[NodeHeading] = ret.renderHeadingVditor
	ret.rendererFuncs[NodeList] = ret.renderListVditor
	ret.rendererFuncs[NodeListItem] = ret.renderListItemVditor
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreakVditor
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreakVditor
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreakVditor
	ret.rendererFuncs[NodeHTMLBlock] = ret.renderHTMLVditor
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTMLVditor
	ret.rendererFuncs[NodeLink] = ret.renderLinkVditor
	ret.rendererFuncs[NodeImage] = ret.renderImageVditor

	// 注册 GFM 渲染函数

	ret.rendererFuncs[NodeStrikethrough] = ret.renderStrikethroughVditor
	ret.rendererFuncs[NodeTaskListItemMarker] = ret.renderTaskListItemMarkerVditor
	ret.rendererFuncs[NodeTable] = ret.renderTableVditor
	ret.rendererFuncs[NodeTableHead] = ret.renderTableHeadVditor
	ret.rendererFuncs[NodeTableRow] = ret.renderTableRowVditor
	ret.rendererFuncs[NodeTableCell] = ret.renderTableCellVditor

	// Emoji 渲染函数

	ret.rendererFuncs[NodeEmojiUnicode] = ret.renderEmojiUnicodeVditor
	ret.rendererFuncs[NodeEmojiImg] = ret.renderEmojiImgVditor

	return
}

func (r *Renderer) renderEmojiImgVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderEmojiUnicodeVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableCellVditor(node *Node, entering bool) (WalkStatus, error) {
	tag := "td"
	if NodeTableHead == node.parent.typ {
		tag = "th"
	}
	if entering {
		var attrs [][]string
		switch node.tableCellAlign {
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

func (r *Renderer) renderTableRowVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("tr", nil, false)
		r.newline()
	} else {
		r.tag("/tr", nil, false)
		r.newline()
		if node == node.parent.lastChild {
			r.tag("/tbody", nil, false)
		}
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableHeadVditor(node *Node, entering bool) (WalkStatus, error) {
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
		if nil != node.next {
			r.tag("tbody", nil, false)
		}
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("table", nil, false)
		r.newline()
	} else {
		r.tag("/table", nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrikethroughVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("del", nil, false)
	} else {
		r.tag("/del", nil, false)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderImageVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if 0 == r.disableTags {
			r.writeString("<img src=\"")
			r.write(escapeHTML(node.destination))
			r.writeString("\" alt=\"")
		}
		r.disableTags++
		return WalkContinue, nil
	}

	r.disableTags--
	if 0 == r.disableTags {
		r.writeString("\"")
		if nil != node.title {
			r.writeString(" title=\"")
			r.write(escapeHTML(node.title))
			r.writeString("\"")
		}
		r.writeString(" />")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderLinkVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.vditorTag("span", -1, nil, false)
		attrs := [][]string{{"class", "open"}}
		r.vditorTag("span", -1, attrs, false)
		r.writeByte(itemOpenBracket)
		r.vditorTag("/span", -1, nil, false)

		attrs = [][]string{{"href", fromItems(escapeHTML(node.destination))}}
		if nil != node.title {
			attrs = append(attrs, []string{"title", fromItems(escapeHTML(node.title))})
		}
		r.tag("a", attrs, false)
	} else {
		r.tag("/a", nil, false)
		attrs := [][]string{{"class", "close"}}
		r.vditorTag("span", -1, attrs, false)
		r.writeByte(itemCloseBracket)
		r.vditorTag("/span", -1, nil, false)
		attrs = [][]string{{"class", "open"}}
		r.vditorTag("span", -1, attrs, false)
		r.writeByte(itemOpenParen)
		r.vditorTag("/span", -1, nil, false)
		r.vditorTag("span", -1, nil, false)
		r.write(node.destination)
		r.vditorTag("/span", -1, nil, false)
		// TODO: title
		attrs = [][]string{{"class", "close"}}
		r.vditorTag("span", -1, attrs, false)
		r.writeByte(itemCloseParen)
		r.vditorTag("/span", -1, nil, false)
		r.vditorTag("/span", -1, nil, false)
	}

	return WalkContinue, nil
}

func (r *Renderer) renderHTMLVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.write(node.tokens)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderInlineHTMLVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderDocumentVditor(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *Renderer) renderParagraphVditor(node *Node, entering bool) (WalkStatus, error) {
	if grandparent := node.parent.parent; nil != grandparent {
		if NodeList == grandparent.typ { // List.ListItem.Paragraph
			if grandparent.tight {
				return WalkContinue, nil
			}
		}
	}

	if entering {
		r.newline()
		r.vditorTag("p", node.typ, nil, false)
	} else {
		r.vditorTag("/p", node.typ, nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTextVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.vditorTag("span", node.typ, nil, false)
		r.write(escapeHTML(node.tokens))
	} else {
		r.vditorTag("/span", node.typ, nil, false)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderCodeSpanVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.vditorTag("span", -1, nil, false)
		attrs := [][]string{{"class", "open"}}
		r.vditorTag("span", -1, attrs, false)
		r.writeByte(itemBacktick)
		if 1 < node.codeMarkerLen {
			r.writeByte(itemBacktick)
		}
		r.vditorTag("/span", -1, nil, false)
		r.vditorTag("code", node.typ, nil, false)
		r.write(escapeHTML(node.tokens))
	} else {
		r.tag("/code", nil, false)
		attrs := [][]string{{"class", "close"}}
		r.vditorTag("span", -1, attrs, false)
		r.writeByte(itemBacktick)
		if 1 < node.codeMarkerLen {
			r.writeByte(itemBacktick)
		}
		r.vditorTag("/span", -1, nil, false)
	}

	return WalkContinue, nil
}

func (r *Renderer) renderCodeBlockVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		tokens := node.tokens
		if 0 < len(node.codeBlockInfo) {
			infoWords := bytes.Split(node.codeBlockInfo, items(" "))
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
	return WalkContinue, nil
}

func (r *Renderer) renderEmphasisVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.vditorTag("span", -1, nil, false)
		attrs := [][]string{{"class", "open"}}
		r.vditorTag("span", -1, attrs, false)
		r.writeByte(node.strongEmDelMarker)
		r.vditorTag("/span", -1, nil, false)
		r.vditorTag("em", node.typ, nil, false)
	} else {
		r.tag("/em", nil, false)
		attrs := [][]string{{"class", "close"}}
		r.vditorTag("span", -1, attrs, false)
		r.writeByte(node.strongEmDelMarker)
		r.vditorTag("/span", -1, nil, false)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrongVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.vditorTag("span", -1, nil, false)
		attrs := [][]string{{"class", "open"}}
		r.vditorTag("span", -1, attrs, false)
		r.write(items{node.strongEmDelMarker, node.strongEmDelMarker})
		r.vditorTag("/span", -1, nil, false)
		r.vditorTag("strong", node.typ, nil, false)
	} else {
		r.tag("/strong", nil, false)
		attrs := [][]string{{"class", "close"}}
		r.vditorTag("span", -1, attrs, false)
		r.write(items{node.strongEmDelMarker, node.strongEmDelMarker})
		r.vditorTag("/span", -1, nil, false)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderBlockquoteVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.vditorTag("blockquote", node.typ, nil, false)
		r.newline()
	} else {
		r.newline()
		r.writeString("</blockquote>")
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderHeadingVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.vditorTag("h"+" 123456"[node.headingLevel:node.headingLevel+1], node.typ, nil, false)
	} else {
		r.writeString("</h" + " 123456"[node.headingLevel:node.headingLevel+1] + ">")
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListVditor(node *Node, entering bool) (WalkStatus, error) {
	tag := "ul"
	if 1 == node.listData.typ {
		tag = "ol"
	}
	if entering {
		r.newline()
		attrs := [][]string{{"start", strconv.Itoa(node.start)}}
		if nil == node.bulletChar && 1 != node.start {
			r.vditorTag(tag, node.typ, attrs, false)
		} else {
			r.vditorTag(tag, node.typ, nil, false)
		}
		r.newline()
	} else {
		r.newline()
		r.vditorTag("/"+tag, node.typ, nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListItemVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if 3 == node.listData.typ && "" != TaskListItemClass {
			r.vditorTag("li", node.typ, [][]string{{"class", TaskListItemClass}}, false)
		} else {
			r.vditorTag("li", node.typ, nil, false)
		}
	} else {
		r.vditorTag("/li", node.typ, nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTaskListItemMarkerVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		var attrs [][]string
		if node.taskListItemChecked {
			attrs = append(attrs, []string{"checked", ""})
		}
		attrs = append(attrs, []string{"disabled", ""}, []string{"type", "checkbox"})
		r.tag("input", attrs, true)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderThematicBreakVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.vditorTag("hr", node.typ, nil, true)
		r.newline()
	}
	return WalkSkipChildren, nil
}

func (r *Renderer) renderHardBreakVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.vditorTag("span", node.typ, nil, false)
		r.vditorTag("/span", node.typ, nil, false)
		r.newline()
	}
	return WalkSkipChildren, nil
}

func (r *Renderer) renderSoftBreakVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.vditorTag("span", node.typ, nil, true)
		r.vditorTag("/span", node.typ, nil, false)
		r.newline()
	}
	return WalkSkipChildren, nil
}

func (r *Renderer) vditorTag(name string, typ int, attrs [][]string, selfclosing bool) {
	if r.disableTags > 0 {
		return
	}

	r.writeString("<")
	r.writeString(name)

	isClosing := itemSlash == name[0]
	if !isClosing {
		if nil == attrs {
			attrs = [][]string{}
		}
		if -1 != typ {
			attrs = append(attrs, []string{"data-id", "0"}, []string{"data-type", strconv.Itoa(typ)})
		}
		for _, attr := range attrs {
			r.writeString(" " + attr[0] + "=\"" + attr[1] + "\"")
		}
	}
	if selfclosing {
		r.writeString(" /")
	}
	r.writeString(">")
}
