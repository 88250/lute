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

// VditorRenderer 描述了 Vditor DOM 渲染器。
type VditorRenderer struct {
	*BaseRenderer
}

// newVditorRenderer 创建一个 Vditor DOM 渲染器。
func (lute *Lute) newVditorRenderer(treeRoot *Node) Renderer {
	ret := &VditorRenderer{&BaseRenderer{rendererFuncs: map[int]RendererFunc{}, option: lute.options, treeRoot: treeRoot}}
	ret.rendererFuncs[NodeDocument] = ret.renderDocumentVditor
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraphVditor
	ret.rendererFuncs[NodeText] = ret.renderTextVditor
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpanVditor
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlockVditor
	ret.rendererFuncs[NodeMathBlock] = ret.renderMathBlockVditor
	ret.rendererFuncs[NodeInlineMath] = ret.renderInlineMathHTML
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
	ret.rendererFuncs[NodeStrikethrough] = ret.renderStrikethroughVditor
	ret.rendererFuncs[NodeTaskListItemMarker] = ret.renderTaskListItemMarkerVditor
	ret.rendererFuncs[NodeTable] = ret.renderTableVditor
	ret.rendererFuncs[NodeTableHead] = ret.renderTableHeadVditor
	ret.rendererFuncs[NodeTableRow] = ret.renderTableRowVditor
	ret.rendererFuncs[NodeTableCell] = ret.renderTableCellVditor
	ret.rendererFuncs[NodeEmojiUnicode] = ret.renderEmojiUnicodeVditor
	ret.rendererFuncs[NodeEmojiImg] = ret.renderEmojiImgVditor
	return ret
}

func (r *VditorRenderer) renderEmojiImgVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmojiUnicodeVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderTableCellVditor(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderTableRowVditor(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderTableHeadVditor(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderTableVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("table", node, nil, false)
	} else {
		r.tag("/table", node, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderStrikethroughVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("del", node, nil, false)
	} else {
		r.tag("/del", node, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderImageVditor(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderLinkVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("span", nil, nil, false)
		attrs := [][]string{{"class", "marker"}}
		r.tag("span", nil, attrs, false)
		r.writeByte(itemOpenBracket)
		r.tag("/span", nil, nil, false)

		attrs = [][]string{{"href", fromItems(escapeHTML(node.destination))}}
		if nil != node.title {
			attrs = append(attrs, []string{"title", fromItems(escapeHTML(node.title))})
		}
		r.tag("a", nil, attrs, false)
	} else {
		r.tag("/a", node, nil, false)
		attrs := [][]string{{"class", "marker"}}
		r.tag("span", nil, attrs, false)
		r.writeByte(itemCloseBracket)
		r.tag("/span", nil, nil, false)
		attrs = [][]string{{"class", "marker"}}
		r.tag("span", nil, attrs, false)
		r.writeByte(itemOpenParen)
		r.tag("/span", nil, nil, false)
		r.tag("span", nil, nil, false)
		r.write(node.destination)
		r.tag("/span", nil, nil, false)
		// TODO: title
		attrs = [][]string{{"class", "marker"}}
		r.tag("span", nil, attrs, false)
		r.writeByte(itemCloseParen)
		r.tag("/span", nil, nil, false)
		r.tag("/span", nil, nil, false)
	}

	return WalkContinue, nil
}

func (r *VditorRenderer) renderHTMLVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderInlineHTMLVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderDocumentVditor(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer) renderParagraphVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("p", node, nil, false)
	} else {
		r.writeString("<span><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></span>")
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderTextVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("span", node, nil, false)
		r.write(escapeHTML(node.tokens))
	} else {
		r.tag("/span", node, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderCodeSpanVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("span", nil, nil, false)
		attrs := [][]string{{"class", "marker"}}
		r.tag("span", nil, attrs, false)
		r.writeByte(itemBacktick)
		if 1 < node.codeMarkerLen {
			r.writeByte(itemBacktick)
		}
		r.tag("/span", nil, nil, false)
		r.tag("code", node, nil, false)
		r.write(escapeHTML(node.tokens))
	} else {
		r.tag("/code", node, nil, false)
		attrs := [][]string{{"class", "marker"}}
		r.tag("span", nil, attrs, false)
		r.writeByte(itemBacktick)
		if 1 < node.codeMarkerLen {
			r.writeByte(itemBacktick)
		}
		r.tag("/span", nil, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderInlineMathHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		tokens := node.tokens
		tokens = escapeHTML(tokens)
		r.write(tokens)
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderMathBlockVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		tokens := node.tokens
		tokens = escapeHTML(tokens)
		r.write(tokens)
		r.newline()
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderCodeBlockVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
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
	return WalkContinue, nil
}

func (r *VditorRenderer) renderEmphasisVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		attrs := [][]string{{"class", "node"}}
		r.tag("span", node, attrs, false)
		attrs = [][]string{{"class", "marker"}}
		r.tag("span", nil, attrs, false)
		r.writeByte(node.strongEmDelMarker)
		r.tag("/span", nil, nil, false)
		r.tag("em", node, nil, false)
	} else {
		r.tag("/em", node, nil, false)
		attrs := [][]string{{"class", "marker"}}
		r.tag("span", nil, attrs, false)
		r.writeByte(node.strongEmDelMarker)
		r.tag("/span", nil, nil, false)
		r.tag("/span", nil, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderStrongVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		attrs := [][]string{{"class", "node"}}
		r.tag("span", nil, attrs, false)
		attrs = [][]string{{"class", "marker"}}
		r.tag("span", nil, attrs, false)
		r.write(items{node.strongEmDelMarker, node.strongEmDelMarker})
		r.tag("/span", nil, nil, false)
		r.tag("strong", node, nil, false)
	} else {
		r.tag("/strong", node, nil, false)
		attrs := [][]string{{"class", "marker"}}
		r.tag("span", nil, attrs, false)
		r.write(items{node.strongEmDelMarker, node.strongEmDelMarker})
		r.tag("/span", nil, nil, false)
		r.tag("/span", nil, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderBlockquoteVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		attrs := [][]string{{"class", "node"}}
		r.tag("span", nil, attrs, false)
		attrs = [][]string{{"class", "marker"}}
		r.tag("span", nil, attrs, false)
		r.writeString("&gt;")
		r.tag("/span", nil, nil, false)
		r.tag("blockquote", node, nil, false)
	} else {
		r.tag("/blockquote", node, nil, false)
		r.tag("/span", nil, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderHeadingVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("h"+" 123456"[node.headingLevel:node.headingLevel+1], node, nil, false)
	} else {
		r.writeString("</h" + " 123456"[node.headingLevel:node.headingLevel+1] + ">")
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderListVditor(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderListItemVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if 3 == node.listData.typ && "" != r.option.GFMTaskListItemClass {
			r.tag("li", node, [][]string{{"class", r.option.GFMTaskListItemClass}}, false)
		} else {
			r.tag("li", node, nil, false)
		}
		attrs := [][]string{{"class", "node"}}
		r.tag("span", nil, attrs, false)
		attrs = [][]string{{"class", "marker"}}
		r.tag("span", nil, attrs, false)

		marker := node.listData.marker
		if 0 != node.listData.delimiter {
			marker = append(marker, node.listData.delimiter)
		}
		r.writeString(fromItems(marker) + " ")
		r.tag("/span", nil, nil, false)
		r.tag("/span", nil, nil, false)
		r.tag("p", nil, nil, false)
	} else {
		r.tag("/p", nil, nil, false)
		r.tag("/li", node, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderTaskListItemMarkerVditor(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderThematicBreakVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("hr", node, nil, true)
	}
	return WalkSkipChildren, nil
}

func (r *VditorRenderer) renderHardBreakVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("span", node, nil, false)
		r.tag("br", nil, nil, false)
		attrs := [][]string{{"class", "newline"}}
		r.tag("span", node, attrs, true)
		r.writeByte(itemNewline)
		r.tag("/span", node, nil, false)
		r.tag("/span", node, nil, false)
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderSoftBreakVditor(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("span", node, nil, false)
		r.tag("br", nil, nil, false)
		attrs := [][]string{{"class", "newline"}}
		r.tag("span", node, attrs, true)
		r.writeByte(itemNewline)
		r.tag("/span", node, nil, false)
		r.tag("/span", node, nil, false)
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
			attrs = append(attrs, []string{"data-ntype", strconv.Itoa(node.typ)})
			attrs = append(attrs, []string{"data-mtype", r.mtype(node.typ)})
			attrs = append(attrs, []string{"data-pos-start", strconv.Itoa(node.srcPosStartLine) + ":" + strconv.Itoa(node.srcPosStartCol)})
			attrs = append(attrs, []string{"data-pos-end", strconv.Itoa(node.srcPosEndLine) + ":" + strconv.Itoa(node.srcPosEndCol)})
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

// mtype 返回节点类型 nodeType 对应的 Markdown 元素类型。
//   0：叶子块元素
//   1：容器块元素
//   2：行级元素
func (r *VditorRenderer) mtype(nodeType int) string {
	switch nodeType {
	case NodeThematicBreak, NodeHeading, NodeCodeBlock, NodeMathBlock, NodeHTMLBlock, NodeParagraph:
		return "0"
	case NodeBlockquote, NodeList, NodeListItem:
		return "1"
	default:
		return "2"
	}
}
