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
	ret := &HTMLRenderer{lute.newBaseRenderer(treeRoot)}
	ret.rendererFuncs[NodeDocument] = ret.renderDocument
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraph
	ret.rendererFuncs[NodeText] = ret.renderText
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpan
	ret.rendererFuncs[NodeCodeSpanOpenMarker] = ret.renderCodeSpanOpenMarker
	ret.rendererFuncs[NodeCodeSpanContent] = ret.renderCodeSpanContent
	ret.rendererFuncs[NodeCodeSpanCloseMarker] = ret.renderCodeSpanCloseMarker
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlock
	ret.rendererFuncs[NodeMathBlock] = ret.renderMathBlock
	ret.rendererFuncs[NodeInlineMath] = ret.renderInlineMath
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasis
	ret.rendererFuncs[NodeEmA6kOpenMarker] = ret.renderEmAsteriskOpenMarker
	ret.rendererFuncs[NodeEmA6kCloseMarker] = ret.renderEmAsteriskCloseMarker
	ret.rendererFuncs[NodeEmU8eOpenMarker] = ret.renderEmUnderscoreOpenMarker
	ret.rendererFuncs[NodeEmU8eCloseMarker] = ret.renderEmUnderscoreCloseMarker
	ret.rendererFuncs[NodeStrong] = ret.renderStrong
	ret.rendererFuncs[NodeStrongA6kOpenMarker] = ret.renderStrongA6kOpenMarker
	ret.rendererFuncs[NodeStrongA6kCloseMarker] = ret.renderStrongA6kCloseMarker
	ret.rendererFuncs[NodeStrongU8eOpenMarker] = ret.renderStrongU8eOpenMarker
	ret.rendererFuncs[NodeStrongU8eCloseMarker] = ret.renderStrongU8eCloseMarker
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquote
	ret.rendererFuncs[NodeBlockquoteMarker] = ret.renderBlockquoteMarker
	ret.rendererFuncs[NodeHeading] = ret.renderHeading
	ret.rendererFuncs[NodeHeadingC8hMarker] = ret.renderHeadingC8hMarker
	ret.rendererFuncs[NodeList] = ret.renderList
	ret.rendererFuncs[NodeListItem] = ret.renderListItem
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreak
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreak
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreak
	ret.rendererFuncs[NodeHTMLBlock] = ret.renderHTML
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTML
	ret.rendererFuncs[NodeLink] = ret.renderLink
	ret.rendererFuncs[NodeImage] = ret.renderImage
	ret.rendererFuncs[NodeBang] = ret.renderBang
	ret.rendererFuncs[NodeOpenBracket] = ret.renderOpenBracket
	ret.rendererFuncs[NodeCloseBracket] = ret.renderCloseBracket
	ret.rendererFuncs[NodeOpenParen] = ret.renderOpenParen
	ret.rendererFuncs[NodeCloseParen] = ret.renderCloseParen
	ret.rendererFuncs[NodeLinkText] = ret.renderLinkText
	ret.rendererFuncs[NodeLinkSpace] = ret.renderLinkSpace
	ret.rendererFuncs[NodeLinkDest] = ret.renderLinkDest
	ret.rendererFuncs[NodeLinkTitle] = ret.renderLinkTitle
	ret.rendererFuncs[NodeStrikethrough] = ret.renderStrikethrough
	ret.rendererFuncs[NodeStrikethrough1OpenMarker] = ret.renderStrikethrough1OpenMarker
	ret.rendererFuncs[NodeStrikethrough1CloseMarker] = ret.renderStrikethrough1CloseMarker
	ret.rendererFuncs[NodeStrikethrough2OpenMarker] = ret.renderStrikethrough2OpenMarker
	ret.rendererFuncs[NodeStrikethrough2CloseMarker] = ret.renderStrikethrough2CloseMarker
	ret.rendererFuncs[NodeTaskListItemMarker] = ret.renderTaskListItemMarker
	ret.rendererFuncs[NodeTable] = ret.renderTable
	ret.rendererFuncs[NodeTableHead] = ret.renderTableHead
	ret.rendererFuncs[NodeTableRow] = ret.renderTableRow
	ret.rendererFuncs[NodeTableCell] = ret.renderTableCell
	ret.rendererFuncs[NodeEmoji] = ret.renderEmoji
	ret.rendererFuncs[NodeEmojiUnicode] = ret.renderEmojiUnicode
	ret.rendererFuncs[NodeEmojiImg] = ret.renderEmojiImg
	ret.rendererFuncs[NodeEmojiAlias] = ret.renderEmojiAlias
	ret.rendererFuncs[NodeVditorHidden] = ret.renderVditorHidden
	return ret
}

func (r *HTMLRenderer) renderVditorHidden(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *HTMLRenderer) renderEmojiAlias(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *HTMLRenderer) renderEmojiImg(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderEmojiUnicode(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderEmoji(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderInlineMath(node *Node, entering bool) (WalkStatus, error) {
	attrs := [][]string{{"class", "vditor-math"}}
	r.tag("span", attrs, false)
	r.write(escapeHTML(node.tokens))
	r.tag("/span", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderMathBlock(node *Node, entering bool) (WalkStatus, error) {
	r.newline()
	attrs := [][]string{{"class", "vditor-math"}}
	r.tag("div", attrs, false)
	r.write(escapeHTML(node.tokens))
	r.tag("/div", nil, false)
	r.newline()
	return WalkStop, nil
}

func (r *HTMLRenderer) renderTableCell(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderTableRow(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderTableHead(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderTable(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("table", nil, false)
		r.newline()
	} else {
		r.tag("/table", nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderStrikethrough(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderStrikethrough1OpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("del", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderStrikethrough1CloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/del", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderStrikethrough2OpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("del", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderStrikethrough2CloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/del", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderLinkTitle(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *HTMLRenderer) renderLinkDest(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *HTMLRenderer) renderLinkSpace(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *HTMLRenderer) renderLinkText(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(escapeHTML(node.tokens))
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) renderCloseParen(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *HTMLRenderer) renderOpenParen(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *HTMLRenderer) renderCloseBracket(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *HTMLRenderer) renderOpenBracket(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *HTMLRenderer) renderBang(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *HTMLRenderer) renderImage(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if 0 == r.disableTags {
			r.writeString("<img src=\"")
			r.write(escapeHTML(node.ChildByType(NodeLinkDest).tokens))
			r.writeString("\" alt=\"")
		}
		r.disableTags++
		return WalkContinue, nil
	}

	r.disableTags--
	if 0 == r.disableTags {
		r.writeString("\"")
		if title := node.ChildByType(NodeLinkTitle); nil != title && nil != title.tokens {
			r.writeString(" title=\"")
			r.write(escapeHTML(title.tokens))
			r.writeString("\"")
		}
		r.writeString(" />")
	}
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderLink(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		dest := node.ChildByType(NodeLinkDest)
		attrs := [][]string{{"href", itemsToStr(escapeHTML(dest.tokens))}}
		if title := node.ChildByType(NodeLinkTitle); nil != title && nil != title.tokens {
			attrs = append(attrs, []string{"title", itemsToStr(escapeHTML(title.tokens))})
		}
		r.tag("a", attrs, false)
	} else {
		r.tag("/a", nil, false)
	}
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderHTML(node *Node, entering bool) (WalkStatus, error) {
	r.newline()
	r.write(node.tokens)
	r.newline()
	return WalkStop, nil
}

func (r *HTMLRenderer) renderInlineHTML(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderDocument(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderParagraph(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderText(node *Node, entering bool) (WalkStatus, error) {
	if r.option.AutoSpace {
		r.space(node)
	}
	if r.option.FixTermTypo {
		r.fixTermTypo(node)
	}
	r.write(escapeHTML(node.tokens))
	return WalkStop, nil
}

func (r *HTMLRenderer) renderCodeSpan(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderCodeSpanOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("<code>")
	return WalkStop, nil
}

func (r *HTMLRenderer) renderCodeSpanContent(node *Node, entering bool) (WalkStatus, error) {
	r.write(escapeHTML(node.tokens))
	return WalkStop, nil
}

func (r *HTMLRenderer) renderCodeSpanCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("</code>")
	return WalkStop, nil
}

func (r *HTMLRenderer) renderEmphasis(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderEmAsteriskOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("em", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderEmAsteriskCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/em", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderEmUnderscoreOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("em", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderEmUnderscoreCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/em", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderStrong(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderStrongA6kOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("strong", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderStrongA6kCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/strong", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderStrongU8eOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("strong", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderStrongU8eCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/strong", nil, false)
	return WalkStop, nil
}

func (r *HTMLRenderer) renderBlockquote(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderBlockquoteMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *HTMLRenderer) renderHeading(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.writeString("<h" + " 123456"[node.headingLevel:node.headingLevel+1] + ">")
	} else {
		r.writeString("</h" + " 123456"[node.headingLevel:node.headingLevel+1] + ">")
		r.newline()
	}
	return WalkContinue, nil
}

func (r *HTMLRenderer) renderHeadingC8hMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *HTMLRenderer) renderList(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderListItem(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderTaskListItemMarker(node *Node, entering bool) (WalkStatus, error) {
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

func (r *HTMLRenderer) renderThematicBreak(node *Node, entering bool) (WalkStatus, error) {
	r.newline()
	r.tag("hr", nil, true)
	r.newline()
	return WalkStop, nil
}

func (r *HTMLRenderer) renderHardBreak(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("br", nil, true)
		r.newline()
	}
	return WalkStop, nil
}

func (r *HTMLRenderer) renderSoftBreak(node *Node, entering bool) (WalkStatus, error) {
	if r.option.SoftBreak2HardBreak {
		r.tag("br", nil, true)
		r.newline()
	} else {
		r.newline()
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
