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
	"bytes"
	"strconv"
	"strings"
)

// VditorRenderer 描述了 Vditor DOM 渲染器。
type VditorRenderer2 struct {
	*BaseRenderer
}

// newVditorRenderer 创建一个 Vditor JSON 渲染器。
func (lute *Lute) newVditorRenderer2(tree *Tree) Renderer {
	ret := &VditorRenderer2{lute.newBaseRenderer(tree)}
	ret.rendererFuncs[NodeDocument] = ret.renderDocument
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraph
	ret.rendererFuncs[NodeText] = ret.renderText
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpan
	ret.rendererFuncs[NodeCodeSpanOpenMarker] = ret.renderCodeSpanOpenMarker
	ret.rendererFuncs[NodeCodeSpanContent] = ret.renderCodeSpanContent
	ret.rendererFuncs[NodeCodeSpanCloseMarker] = ret.renderCodeSpanCloseMarker
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlock
	ret.rendererFuncs[NodeCodeBlockFenceOpenMarker] = ret.renderCodeBlockOpenMarker
	ret.rendererFuncs[NodeCodeBlockFenceInfoMarker] = ret.renderCodeBlockInfoMarker
	ret.rendererFuncs[NodeCodeBlockCode] = ret.renderCodeBlockCode
	ret.rendererFuncs[NodeCodeBlockFenceCloseMarker] = ret.renderCodeBlockCloseMarker
	ret.rendererFuncs[NodeMathBlock] = ret.renderMathBlock
	ret.rendererFuncs[NodeMathBlockOpenMarker] = ret.renderMathBlockOpenMarker
	ret.rendererFuncs[NodeMathBlockContent] = ret.renderMathBlockContent
	ret.rendererFuncs[NodeMathBlockCloseMarker] = ret.renderMathBlockCloseMarker
	ret.rendererFuncs[NodeInlineMath] = ret.renderInlineMath
	ret.rendererFuncs[NodeInlineMathOpenMarker] = ret.renderInlineMathOpenMarker
	ret.rendererFuncs[NodeInlineMathContent] = ret.renderInlineMathContent
	ret.rendererFuncs[NodeInlineMathCloseMarker] = ret.renderInlineMathCloseMarker
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
	return ret
}

func (r *VditorRenderer2) renderCodeBlockCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderCodeBlockInfoMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderCodeBlockOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderEmojiAlias(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderEmojiImg(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *VditorRenderer2) renderEmojiUnicode(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *VditorRenderer2) renderEmoji(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderInlineMathCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderInlineMathContent(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if !strings.HasSuffix(node.parent.PreviousNodeText(), " ") {
			r.writeByte(itemSpace)
		}
		r.writeString("<span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\">")
		r.writeString("<code data-type=\"math-inline\">")
		r.write(node.tokens)
	} else {
		r.writeString("</code></span>")
		if !strings.HasPrefix(node.parent.NextNodeText(), " ") {
			r.writeByte(itemSpace)
		}
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderInlineMathOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderInlineMath(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderMathBlockCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderMathBlockContent(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<div class=\"vditor-wysiwyg__block\" data-type=\"math-block\">")
		r.writeString("<pre><code data-type=\"math-block\">")
		r.write(node.tokens)
	} else {
		r.writeString("</code></pre></div>")
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderMathBlockOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderMathBlock(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderTableCell(node *Node, entering bool) (WalkStatus, error) {
	tag := "td"
	if NodeTableHead == node.parent.parent.typ {
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
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderTableRow(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("tr", nil, false)
	} else {
		r.tag("/tr", nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderTableHead(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("thead", nil, false)
	} else {
		r.tag("/thead", nil, false)
		if nil != node.next {
			r.tag("tbody", nil, false)
		}
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderTable(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("table", nil, false)
	} else {
		if nil != node.firstChild.next {
			r.tag("/tbody", nil, false)
		}
		r.tag("/table", nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderStrikethrough(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderStrikethrough1OpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("s", [][]string{{"data-marker", "~"}}, false)
	if caret == string(node.next.tokens) {
		node.next.tokens = append([]byte("\n"), node.next.tokens...)
	}
	return WalkStop, nil
}

func (r *VditorRenderer2) renderStrikethrough1CloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/s", nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer2) renderStrikethrough2OpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("s", [][]string{{"data-marker", "~~"}}, false)
	if caret == string(node.next.tokens) {
		node.next.tokens = append([]byte("\n"), node.next.tokens...)
	}
	return WalkStop, nil
}

func (r *VditorRenderer2) renderStrikethrough2CloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/s", nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer2) renderLinkTitle(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderLinkDest(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderLinkSpace(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderLinkText(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkStop, nil
}

func (r *VditorRenderer2) renderCloseParen(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderOpenParen(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderCloseBracket(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderOpenBracket(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderBang(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderImage(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if 0 == r.disableTags {
			r.writeString("<img src=\"")
			r.write(node.ChildByType(NodeLinkDest).tokens)
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
			r.write(title.tokens)
			r.writeString("\"")
		}
		r.writeString(" />")
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderLink(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		dest := node.ChildByType(NodeLinkDest)
		attrs := [][]string{{"href", bytesToStr(dest.tokens)}}
		if title := node.ChildByType(NodeLinkTitle); nil != title && nil != title.tokens {
			attrs = append(attrs, []string{"title", bytesToStr(title.tokens)})
		}
		r.tag("a", attrs, false)
	} else {
		r.tag("/a", nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\">")
		r.writeString("<pre><code data-type=\"html-block\">")
		r.write(node.tokens)
	} else {
		r.writeString("</code></pre>")
		r.writeString("</div>")
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderInlineHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if !strings.HasSuffix(node.PreviousNodeText(), " ") {
			r.writeByte(itemSpace)
		}
		r.writeString("<span class=\"vditor-wysiwyg__block\" data-type=\"html-inline\">")
		r.writeString("<code data-type=\"html-inline\">")
		r.write(node.tokens)
	} else {
		r.writeString("</code></span>")
		if !strings.HasPrefix(node.NextNodeText(), " ") {
			r.writeByte(itemSpace)
		}
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderDocument(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte(itemOpenBracket)
	} else {
		r.writeByte(itemCloseBracket)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderParagraph(node *Node, entering bool) (WalkStatus, error) {
	if grandparent := node.parent.parent; nil != grandparent {
		if NodeList == grandparent.typ { // List.ListItem.Paragraph
			if grandparent.tight {
				return WalkContinue, nil
			}
		}
	}

	if entering {
		r.tag("p", nil, false)
	} else {
		if nil != node.firstChild && NodeText == node.firstChild.typ && caret == bytesToStr(node.firstChild.tokens) {
			r.writeByte(itemNewline)
		}
		r.tag("/p", nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderText(node *Node, entering bool) (WalkStatus, error) {
	if r.option.AutoSpace {
		r.space(node)
	}
	if r.option.FixTermTypo {
		r.fixTermTypo(node)
	}
	if r.option.ChinesePunct {
		r.chinesePunct(node)
	}
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *VditorRenderer2) renderCodeSpan(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderCodeSpanOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	if !strings.HasSuffix(node.parent.PreviousNodeText(), " ") {
		r.writeByte(itemSpace)
	}
	r.writeString("<code>")
	return WalkStop, nil
}

func (r *VditorRenderer2) renderCodeSpanContent(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *VditorRenderer2) renderCodeSpanCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("</code>")
	if !strings.HasPrefix(node.parent.NextNodeText(), " ") {
		r.writeByte(itemSpace)
	}
	return WalkStop, nil
}

func (r *VditorRenderer2) renderEmphasis(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderEmAsteriskOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("em", [][]string{{"data-marker", "*"}}, false)
	if caret == string(node.next.tokens) {
		node.next.tokens = append([]byte("\n"), node.next.tokens...)
	}
	return WalkStop, nil
}

func (r *VditorRenderer2) renderEmAsteriskCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/em", nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer2) renderEmUnderscoreOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("em", [][]string{{"data-marker", "_"}}, false)
	if caret == string(node.next.tokens) {
		node.next.tokens = append([]byte("\n"), node.next.tokens...)
	}
	return WalkStop, nil
}

func (r *VditorRenderer2) renderEmUnderscoreCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/em", nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer2) renderStrong(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderStrongA6kOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("strong", [][]string{{"data-marker", "**"}}, false)
	if caret == string(node.next.tokens) {
		node.next.tokens = append([]byte("\n"), node.next.tokens...)
	}
	return WalkStop, nil
}

func (r *VditorRenderer2) renderStrongA6kCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/strong", nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer2) renderStrongU8eOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("strong", [][]string{{"data-marker", "__"}}, false)
	if caret == string(node.next.tokens) {
		node.next.tokens = append([]byte("\n"), node.next.tokens...)
	}
	return WalkStop, nil
}

func (r *VditorRenderer2) renderStrongU8eCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/strong", nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer2) renderBlockquote(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<blockquote>")
	} else {
		r.writeString("</blockquote>")
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderBlockquoteMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderHeading(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<h" + " 123456"[node.headingLevel:node.headingLevel+1] + ">")
		if r.option.HeadingAnchor {
			anchor := node.Text()
			anchor = strings.ReplaceAll(anchor, " ", "-")
			r.tag("a", [][]string{{"id", "vditorAnchor-" + anchor}, {"class", "vditor-anchor"}, {"href", "#" + anchor}}, false)
			r.writeString(`<svg viewBox="0 0 16 16" version="1.1" width="16" height="16"><path fill-rule="evenodd" d="M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z"></path></svg>`)
			r.tag("/a", nil, false)
		}
	} else {
		r.writeString("</h" + " 123456"[node.headingLevel:node.headingLevel+1] + ">")
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderHeadingC8hMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer2) renderList(node *Node, entering bool) (WalkStatus, error) {
	tag := "ul"
	if 1 == node.listData.typ {
		tag = "ol"
	}
	if entering {
		var attrs [][]string
		if node.tight {
			attrs = append(attrs, []string{"data-tight", "true"})
		}
		if nil == node.bulletChar && 1 != node.start {
			attrs = append(attrs, []string{"start", strconv.Itoa(node.start)})
		}
		r.tag(tag, attrs, false)
	} else {
		r.tag("/"+tag, nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderListItem(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		var attrs [][]string
		switch node.listData.typ {
		case 0:
			attrs = append(attrs, []string{"data-marker", bytesToStr(node.marker)})
		case 1:
			attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.num) + "."})
		case 3:
			attrs = append(attrs, []string{"data-marker", bytesToStr(node.marker)})
			attrs = append(attrs, []string{"class", r.option.GFMTaskListItemClass})
		}
		r.tag("li", attrs, false)
	} else {
		r.tag("/li", nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderTaskListItemMarker(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		var attrs [][]string
		if node.taskListItemChecked {
			attrs = append(attrs, []string{"checked", ""})
		}
		attrs = append(attrs, []string{"type", "checkbox"})
		r.tag("input", attrs, true)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderThematicBreak(node *Node, entering bool) (WalkStatus, error) {
	r.tag("hr", nil, true)
	return WalkStop, nil
}

func (r *VditorRenderer2) renderHardBreak(node *Node, entering bool) (WalkStatus, error) {
	r.tag("br", nil, true)
	return WalkStop, nil
}

func (r *VditorRenderer2) renderSoftBreak(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemNewline)
	return WalkStop, nil
}

func (r *VditorRenderer2) tag(name string, attrs [][]string, selfclosing bool) {
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

func (r *VditorRenderer2) renderCodeBlock(node *Node, entering bool) (WalkStatus, error) {
	if !node.isFencedCodeBlock {
		// TODO: 移除缩进代码块处理
		r.writeString("<pre><code>")
		r.write(node.tokens)
		r.writeString("</code></pre>")
		return WalkStop, nil
	}
	if entering {
		r.writeString("<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\">")
	} else {
		r.writeString("</div>")
	}
	return WalkContinue, nil
}

func (r *VditorRenderer2) renderCodeBlockCode(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if bytes.Contains(node.previous.codeBlockInfo, []byte(caret)) {
			node.previous.codeBlockInfo = bytes.ReplaceAll(node.previous.codeBlockInfo, []byte(caret), []byte(""))
			node.tokens = append(node.tokens, []byte("\n"+caret)...)
		}

		if 0 < len(node.previous.codeBlockInfo) {
			infoWords := split(node.previous.codeBlockInfo, itemSpace)
			language := bytesToStr(infoWords[0])
			r.writeString("<pre><code class=\"language-" + language + "\">")
		} else {
			r.writeString("<pre><code>")
		}
		r.write(node.tokens)
		return WalkSkipChildren, nil
	}
	r.writeString("</code></pre>")
	return WalkStop, nil
}

func (r *VditorRenderer2) html() (ret string) {
	ret = r.writer.String()
	r.writer.Reset()
	return
}

func (r *VditorRenderer2) val(val string, node *Node) {
	val = strings.ReplaceAll(val, "\"", "\\\"")
	r.writeString("\"name\":\"" + val + "\"")
}

func (r *VditorRenderer2) openObj() {
	r.writeByte('{')
}

func (r *VditorRenderer2) closeObj(node *Node) {
	r.writeByte('}')
	if !r.ignore(node.next) {
		r.comma()
	}
}

func (r *VditorRenderer2) comma() {
	r.writeString(",")
}

func (r *VditorRenderer2) ignore(node *Node) bool {
	return nil == node ||
		// 以下类型的节点不进行渲染，否则图画出来节点太多
		NodeBlockquoteMarker == node.typ ||
		NodeEmA6kOpenMarker == node.typ || NodeEmA6kCloseMarker == node.typ ||
		NodeEmU8eOpenMarker == node.typ || NodeEmU8eCloseMarker == node.typ ||
		NodeStrongA6kOpenMarker == node.typ || NodeStrongA6kCloseMarker == node.typ ||
		NodeStrongU8eOpenMarker == node.typ || NodeStrongU8eCloseMarker == node.typ ||
		NodeStrikethrough1OpenMarker == node.typ || NodeStrikethrough1CloseMarker == node.typ ||
		NodeStrikethrough2OpenMarker == node.typ || NodeStrikethrough2CloseMarker == node.typ ||
		NodeMathBlockOpenMarker == node.typ || NodeMathBlockContent == node.typ || NodeMathBlockCloseMarker == node.typ ||
		NodeInlineMathOpenMarker == node.typ || NodeInlineMathContent == node.typ || NodeInlineMathCloseMarker == node.typ
}
