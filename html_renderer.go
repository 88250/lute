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
	"strconv"
)

// HTMLRenderer 描述了 HTML 渲染器。
type HTMLRenderer struct {
	*BaseRenderer
}

// newHTMLRenderer 创建一个 HTML 渲染器。
func (lute *Lute) newHTMLRenderer(treeRoot *Node) Renderer {
	ret := &HTMLRenderer{&BaseRenderer{rendererFuncs: map[int]RendererFunc{}, option: lute.options, treeRoot: treeRoot}}
	ret.rendererFuncs[NodeDocument] = ret.renderDocumentHTML
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraphHTML
	ret.rendererFuncs[NodeText] = ret.renderTextHTML
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpanHTML
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlockHTML
	ret.rendererFuncs[NodeMathBlock] = ret.renderMathBlockHTML
	ret.rendererFuncs[NodeInlineMath] = ret.renderInlineMathHTML
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasisHTML
	ret.rendererFuncs[NodeStrong] = ret.renderStrongHTML
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquoteHTML
	ret.rendererFuncs[NodeHeading] = ret.renderHeadingHTML
	ret.rendererFuncs[NodeList] = ret.renderListHTML
	ret.rendererFuncs[NodeListItem] = ret.renderListItemHTML
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreakHTML
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreakHTML
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreakHTML
	ret.rendererFuncs[NodeHTMLBlock] = ret.renderHTMLHTML
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTMLHTML
	ret.rendererFuncs[NodeLink] = ret.renderLinkHTML
	ret.rendererFuncs[NodeImage] = ret.renderImageHTML
	ret.rendererFuncs[NodeStrikethrough] = ret.renderStrikethroughHTML
	ret.rendererFuncs[NodeTaskListItemMarker] = ret.renderTaskListItemMarkerHTML
	ret.rendererFuncs[NodeTable] = ret.renderTableHTML
	ret.rendererFuncs[NodeTableHead] = ret.renderTableHeadHTML
	ret.rendererFuncs[NodeTableRow] = ret.renderTableRowHTML
	ret.rendererFuncs[NodeTableCell] = ret.renderTableCellHTML
	ret.rendererFuncs[NodeEmojiUnicode] = ret.renderEmojiUnicodeHTML
	ret.rendererFuncs[NodeEmojiImg] = ret.renderEmojiImgHTML
	ret.rendererFuncs[NodeVditorCaret] = ret.renderVditorCaretVditor
	return ret
}

func (r *HTMLRenderer) renderVditorCaretVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) renderInlineMathHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		tokens := node.tokens
		tokens = escapeHTML(tokens)
		r.write(tokens)
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) renderMathBlockHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		tokens := node.tokens
		tokens = escapeHTML(tokens)
		r.write(tokens)
		r.newline()
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) renderEmojiImgHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) renderEmojiUnicodeHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) renderTableCellHTML(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderTableRowHTML(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderTableHeadHTML(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderTableHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("table", nil, false)
		r.newline()
	} else {
		r.tag("/table", nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderStrikethroughHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("del", nil, false)
	} else {
		r.tag("/del", nil, false)
	}
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderImageHTML(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderLinkHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		attrs := [][]string{{"href", fromItems(escapeHTML(node.destination))}}
		if nil != node.title {
			attrs = append(attrs, []string{"title", fromItems(escapeHTML(node.title))})
		}
		r.tag("a", attrs, false)

		return WalkContinue, nil
	}

	r.tag("/a", nil, false)
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderHTMLHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.write(node.tokens)
		r.newline()
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) renderInlineHTMLHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) renderDocumentHTML(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderParagraphHTML(node *Node, entering bool) (WalkStatus, error) {
	if grandparent := node.parent.parent; nil != grandparent {
		if NodeList == grandparent.typ { // List.ListItem.Paragraph
			if grandparent.tight {
				return WalkContinue, nil
			}
		}
	}

	if entering {
		r.newline()
		r.tag("p", nil, false)
	} else {
		r.tag("/p", nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderTextHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(escapeHTML(node.tokens))
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) renderCodeSpanHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<code>")
		r.write(escapeHTML(node.tokens))
		return WalkSkipChildren, nil
	}
	r.writeString("</code>")
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderEmphasisHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("em", nil, false)
	} else {
		r.tag("/em", nil, false)
	}
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderStrongHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<strong>")
		r.write(node.tokens)
	} else {
		r.writeString("</strong>")
	}
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderBlockquoteHTML(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderHeadingHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.writeString("<h" + " 123456"[node.headingLevel:node.headingLevel+1] + ">")
	} else {
		r.writeString("</h" + " 123456"[node.headingLevel:node.headingLevel+1] + ">")
		r.newline()
	}
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderListHTML(node *Node, entering bool) (WalkStatus, error) {
	tag := "ul"
	if 1 == node.listData.typ {
		tag = "ol"
	}
	if entering {
		r.newline()
		attrs := [][]string{{"start", strconv.Itoa(node.start)}}
		if nil == node.bulletChar && 1 != node.start {
			r.tag(tag, attrs, false)
		} else {
			r.tag(tag, nil, false)
		}
		r.newline()
	} else {
		r.newline()
		r.tag("/"+tag, nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderListItemHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if 3 == node.listData.typ && "" != r.option.GFMTaskListItemClass {
			r.tag("li", [][]string{{"class", r.option.GFMTaskListItemClass}}, false)
		} else {
			r.tag("li", nil, false)
		}
	} else {
		r.tag("/li", nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderTaskListItemMarkerHTML(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderThematicBreakHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.tag("hr", nil, true)
		r.newline()
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) renderHardBreakHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("br", nil, true)
		r.newline()
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) renderSoftBreakHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if r.option.SoftBreak2HardBreak {
			r.tag("br", nil, true)
			r.newline()
		} else {
			r.newline()
		}
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) tag(name string, attrs [][]string, selfclosing bool) {
	if r.disableTags > 0 {
		return
	}

	r.writeString("<")
	r.writeString(name)
	if 0 < len(attrs) {
		for _, attr := range attrs {
			r.writeString(" " + attr[0] + "=\"" + attr[1] + "\"")
		}
	}
	if selfclosing {
		r.writeString(" /")
	}
	r.writeString(">")
}
