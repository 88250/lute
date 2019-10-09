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

// Vditor DOM Renderer

import (
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// VditorRenderer 描述了 Vditor DOM 渲染器。
type VditorRenderer struct {
	*BaseRenderer
	lastOut string
}

// newVditorRenderer 创建一个 Vditor DOM 渲染器。
func (lute *Lute) newVditorRenderer(treeRoot *Node) *VditorRenderer {
	ret := &VditorRenderer{BaseRenderer: lute.newBaseRenderer(treeRoot)}
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
	ret.rendererFuncs[NodeEmojiUnicode] = ret.renderEmojiUnicode
	ret.rendererFuncs[NodeEmojiImg] = ret.renderEmojiImg
	return ret
}

func (r *VditorRenderer) renderEmojiImg(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmojiUnicode(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderTableCell(node *Node, entering bool) (WalkStatus, error) {
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
		r.tag(tag, node, attrs, false)
	} else {
		r.tag("/"+tag, node, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderTableRow(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("tr", node, nil, false)
	} else {
		r.tag("/tr", node, nil, false)
		if node == node.parent.lastChild {
			r.tag("/tbody", node, nil, false)
		}
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderTableHead(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("thead", node, nil, false)
		r.tag("tr", node, nil, false)
	} else {
		r.tag("/tr", node, nil, false)
		r.tag("/thead", node, nil, false)
		if nil != node.next {
			r.tag("tbody", node, nil, false)
		}
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderTable(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("table", node, nil, false)
	} else {
		r.tag("/table", node, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderStrikethrough(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("del", node, [][]string{{"class", "node"}}, false)
	} else {
		r.tag("/del", nil, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderStrikethrough1OpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeString("~")
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrikethrough1CloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeString("~")
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrikethrough2OpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeString("~~")
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrikethrough2CloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeString("~~")
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderLinkTitle(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeString(" \"")
	r.write(node.tokens)
	r.writeString("\"")
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderLinkDest(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.write(node.tokens)
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderLinkText(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, nil, false)
	r.write(escapeHTML(node.tokens))
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderCloseParen(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeByte(itemCloseParen)
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderOpenParen(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeByte(itemOpenParen)
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderCloseBracket(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeByte(itemCloseBracket)
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderOpenBracket(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeByte(itemOpenBracket)
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderBang(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderImage(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderLink(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		dest := node.ChildByType(NodeLinkDest)
		attrs := [][]string{{"class", "node"}, {"href", itemsToStr(escapeHTML(dest.tokens))}}
		if title := node.ChildByType(NodeLinkTitle); nil != title && nil != title.tokens {
			attrs = append(attrs, []string{"title", itemsToStr(escapeHTML(title.tokens))})
		}
		r.tag("a", node, attrs, false)
	} else {
		r.tag("/a", nil, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderInlineHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderDocument(node *Node, entering bool) (WalkStatus, error) {
	if nil == node.firstChild {
		r.writeString("<p data-ntype=\"" + NodeParagraph.String() + "\" data-mtype=\"" + r.mtype(NodeParagraph) + "\">" +
			"<span data-ntype=\"" + NodeParagraph.String() + " data-mtype=\"2\" data-cso=\"0\" data-ceo=\"0\"></span>" +
			"<span class=\"newline\">\n</span></p>")
		return WalkStop, nil
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderParagraph(node *Node, entering bool) (WalkStatus, error) {
	if parent := node.parent; nil != parent {
		if NodeListItem == parent.typ { // ListItem.Paragraph
			if parent.tight || parent.parent.tight {
				return WalkContinue, nil
			}
		}
	}

	if entering {
		r.tag("p", node, nil, false)
	} else {
		r.writeString("<span class=\"newline\">\n\n</span>")
		r.tag("/p", nil, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderText(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("span", node, nil, false)
		r.write(escapeHTML(node.tokens))
		r.tag("/span", nil, nil, false)
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderCodeSpan(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("code", node, [][]string{{"class", "node"}}, false)
	} else {
		r.tag("/code", nil, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderCodeSpanOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.write(node.tokens)
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderCodeSpanContent(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(escapeHTML(node.tokens))
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderCodeSpanCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.write(node.tokens)
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderInlineMath(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		tokens := node.tokens
		tokens = escapeHTML(tokens)
		r.write(tokens)
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderMathBlock(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		tokens := node.tokens
		tokens = escapeHTML(tokens)
		r.write(tokens)
		r.newline()
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderCodeBlock(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		tokens := node.tokens
		if 0 < len(node.codeBlockInfo) {
			infoWords := split(node.codeBlockInfo, itemSpace)
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
	return WalkContinue, nil
}

func (r *VditorRenderer) renderEmphasis(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("em", node, [][]string{{"class", "node"}}, false)
	} else {
		r.tag("/em", nil, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderEmAsteriskOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeByte(itemAsterisk)
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmAsteriskCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeByte(itemAsterisk)
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmUnderscoreOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeByte(itemUnderscore)
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmUnderscoreCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeByte(itemUnderscore)
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrong(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("strong", node, [][]string{{"class", "node"}}, false)
	} else {
		r.tag("/strong", nil, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderStrongA6kOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeString("**")
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrongA6kCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeString("**")
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrongU8eOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeString("__")
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrongU8eCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.writeString("__")
	r.tag("/span", nil, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderBlockquote(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.tag("/blockquote", node, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderBlockquoteMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("blockquote", node.parent, [][]string{{"class", "node node--block"}}, false)
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.write(escapeHTML(node.tokens))
	r.tag("/span", node, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderHeading(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("h"+" 123456"[node.headingLevel:node.headingLevel+1], node, [][]string{{"class", "node"}}, false)
	} else {
		r.writeString("</h" + " 123456"[node.headingLevel:node.headingLevel+1] + ">")
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderHeadingC8hMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("span", node, [][]string{{"class", "marker"}}, false)
	r.write(node.tokens)
	r.tag("/span", node, nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderList(node *Node, entering bool) (WalkStatus, error) {
	tag := "ul"
	if 1 == node.listData.typ {
		tag = "ol"
	}
	if entering {
		attrs := [][]string{{"start", strconv.Itoa(node.start)}}
		if nil == node.bulletChar && 1 != node.start {
			r.tag(tag, node, attrs, false)
		} else {
			r.tag(tag, node, nil, false)
		}
	} else {
		r.tag("/"+tag, node, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderListItem(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if 3 == node.listData.typ && "" != r.option.GFMTaskListItemClass {
			r.tag("li", node, [][]string{{"class", "node node--block " + r.option.GFMTaskListItemClass}}, false)
		} else {
			r.tag("li", node, [][]string{{"class", "node node--block"}}, false)
		}
		r.tag("span", nil, [][]string{{"class", "marker"}}, false)

		marker := node.listData.marker
		if !isNilItem(node.listData.delimiter) {
			marker = append(marker, node.listData.delimiter)
		}
		r.writeString(itemsToStr(marker) + " ")
		r.tag("/span", nil, nil, false)
	} else {
		if node.tight || node.parent.tight {
			r.writeString("<span class=\"newline\">\n</span>")
		}
		r.tag("/li", node, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderTaskListItemMarker(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		var attrs [][]string
		if node.taskListItemChecked {
			attrs = append(attrs, []string{"checked", ""})
		}
		attrs = append(attrs, []string{"disabled", ""}, []string{"type", "checkbox"})
		r.tag("input", node, attrs, true)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderThematicBreak(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("hr", node, nil, true)
	}
	return WalkSkipChildren, nil
}

func (r *VditorRenderer) renderHardBreak(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte(itemNewline)
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderSoftBreak(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte(itemNewline)
	}
	return WalkStop, nil
}

func (r *VditorRenderer) tag(name string, node *Node, attrs [][]string, selfclosing bool) {
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
		if nil != node {
			if !r.containClass(&attrs, "marker") && "span" != name {
				attrs = append(attrs, []string{"data-ntype", node.typ.String()})
				attrs = append(attrs, []string{"data-mtype", r.mtype(node.typ)})
			}
			if "" != node.caretStartOffset {
				attrs = append(attrs, []string{"data-cso", node.caretStartOffset})
			}
			if "" != node.caretEndOffset {
				attrs = append(attrs, []string{"data-ceo", node.caretEndOffset})
			}
			if node.expand && r.containClass(&attrs, "node") {
				r.appendClass(&attrs, "node--expand")
			}
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

func (r *VditorRenderer) appendClass(attrs *[][]string, class string) {
	for _, attr := range *attrs {
		if "class" == attr[0] && !strings.Contains(attr[1], class) {
			attr[1] += " " + class
			return
		}
	}
	*attrs = append(*attrs, []string{"class", class})
}

func (r *VditorRenderer) containClass(attrs *[][]string, class string) bool {
	for _, attr := range *attrs {
		if "class" == attr[0] && strings.Contains(attr[1], class) {
			return true
		}
	}
	return false
}

// mtype 返回节点类型 nodeType 对应的 Markdown 元素类型。
//   0：叶子块元素
//   1：容器块元素
//   2：行级元素
func (r *VditorRenderer) mtype(nodeType nodeType) string {
	switch nodeType {
	case NodeThematicBreak, NodeHeading, NodeCodeBlock, NodeMathBlock, NodeHTMLBlock, NodeParagraph:
		return "0"
	case NodeDocument, NodeBlockquote, NodeList, NodeListItem:
		return "1"
	default:
		return "2"
	}
}

// mapSelection 用于映射文本选段。
// 根据 Markdown 原文中的 startOffset 和 endOffset 位置从根节点 root 上遍历查找该范围内的节点，找到对应节点后进行标记以支持后续渲染。
func (r *VditorRenderer) mapSelection(root *Node, startOffset, endOffset int) {
	var nodes []*Node
	for c := root.firstChild; nil != c; c = c.next {
		r.findSelection(c, startOffset, endOffset, &nodes)
	}

	if 1 > len(nodes) {
		// 当且仅当渲染空 Markdown 时
		nodes = append(nodes, root)
	}

	// TODO: 仅实现了光标插入位置节点获取，选段映射待实现
	en := r.nearest(nodes, endOffset)
	sn := en
	base := 0
	if 0 < len(sn.tokens) {
		base = sn.tokens[0].Offset()
	}

	startOffset = startOffset - base
	endOffset = endOffset - base
	startOffset, endOffset = r.runeOffset(itemsToBytes(sn.tokens), startOffset, endOffset)
	sn.caretStartOffset = strconv.Itoa(startOffset)
	en.caretEndOffset = strconv.Itoa(endOffset)
	r.expand(sn)
}

// expand 用于在 node 上或者 node 的祖先节点上标记展开。
func (r *VditorRenderer) expand(node *Node) {
	for p := node; nil != p; p = p.parent {
		switch p.typ {
		case NodeEmphasis, NodeStrong, NodeBlockquote, NodeListItem, NodeCodeSpan, NodeHeading, NodeLink:
			p.expand = true
			return
		}
	}
}

func (tokens items) Len() int           { return len(tokens) }
func (tokens items) Swap(i, j int)      { tokens[i], tokens[j] = tokens[j], tokens[i] }
func (tokens items) Less(i, j int) bool { return tokens[i].Offset() < tokens[j].Offset() }

// findSelection 在 node 上递归查找 startOffset 和 endOffset 选段，选中节点累计到 selected 中。
func (r *VditorRenderer) findSelection(node *Node, startOffset, endOffset int, selected *[]*Node) {
	nodes := node.List()
	length := len(nodes)
	var n *Node
	tokens := make(items, 0, len(nodes)*4)
	for i := 0; i < length; i++ {
		n = nodes[i]
		for i, _ := range n.tokens {
			n.tokens[i].node = n
		}
		tokens = append(tokens, n.tokens...)
	}

	sort.Sort(tokens)

	var token item
	var startToken, endToken *item
	length = len(tokens)
	for i := 0; i < length; i++ {
		token = tokens[i]
		if nil == startToken && startOffset <= token.Offset()+1 {
			startToken = &token
		}
		if endOffset <= token.Offset()+1 {
			endToken = &token
			break
		}
	}

	if nil != startToken {
		*selected = append(*selected, startToken.node)
		if startToken.node != endToken.node {
			*selected = append(*selected, endToken.node)
		}
	}
}

// nearest 在 selected 节点列表中查找离 offset 最近的节点。
func (r *VditorRenderer) nearest(selected []*Node, offset int) (ret *Node) {
	dist := 16
	for i := 0; i < len(selected); i++ {
		n := selected[i]
		s, e := n.Range()
		tmpSDis := s - offset
		if 0 > tmpSDis {
			tmpSDis = -tmpSDis
		}
		tmpEDis := e - offset
		if 0 > tmpEDis {
			tmpEDis = -tmpEDis
		}
		minDist := tmpSDis
		if minDist > tmpEDis {
			minDist = tmpEDis
		}
		if minDist < dist {
			dist = minDist
			ret = n
		}
	}
	return
}

// byteOffset 返回字符偏移位置在 str 中考虑字符编码情况下的字节偏移位置。
func (r *VditorRenderer) byteOffset(str string, runeStartOffset, runeEndOffset int) (startOffset, endOffset int) {
	runes := 0
	for i, _ := range str {
		runes++
		if runes > runeStartOffset {
			startOffset = i
		}
		if runes > runeEndOffset {
			endOffset = i
			return
		}
	}
	return
}

// runeOffset 返回字节偏移位置在 bytes 中考虑字符编码情况下的字符偏移位置。
func (r *VditorRenderer) runeOffset(bytes []byte, byteStartOffset, byteEndOffset int) (runeStartOffset, runeEndOffset int) {
	length := len(bytes)
	var i, size int
	for ; i < length; i += size {
		_, size = utf8.DecodeRune(bytes[i:])
		if i < byteStartOffset {
			runeStartOffset++
		}
		if i < byteEndOffset {
			runeEndOffset++
		} else {
			return
		}
	}
	return
}
