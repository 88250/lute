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
	"strings"
)

// VditorRenderer 描述了 Vditor DOM 渲染器。
type VditorRenderer struct {
	*BaseRenderer
}

// newVditorRenderer 创建一个 HTML 渲染器。
func (lute *Lute) newVditorRenderer(tree *Tree) Renderer {
	ret := &VditorRenderer{lute.newBaseRenderer(tree)}
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

func (r *VditorRenderer) renderCodeBlockCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderCodeBlockInfoMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderCodeBlockOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmojiAlias(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmojiImg(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmojiUnicode(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmoji(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer) renderInlineMathCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderInlineMathContent(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderInlineMathOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderInlineMath(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer) renderMathBlockCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderMathBlockContent(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<div class=\"vditor-wysiwyg__block\" data-type=\"math-block\">")
		r.writeString("<pre><code data-type=\"math-block\">")
		r.write(node.tokens)
	} else {
		r.writeString("</code></pre></div>")
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderMathBlockOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderMathBlock(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer) renderTableCell(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderTableRow(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("tr", nil, false)
	} else {
		r.tag("/tr", nil, false)
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderTableHead(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderTable(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderStrikethrough(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer) renderStrikethrough1OpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("s", [][]string{{"data-marker", "~"}}, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrikethrough1CloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/s", nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrikethrough2OpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("s", [][]string{{"data-marker", "~~"}}, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrikethrough2CloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/s", nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderLinkTitle(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderLinkDest(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderLinkSpace(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderLinkText(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderCloseParen(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderOpenParen(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderCloseBracket(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderOpenBracket(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderBang(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderImage(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderLink(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderHTML(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderInlineHTML(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderDocument(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer) renderParagraph(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderText(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderCodeSpan(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer) renderCodeSpanOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	if !strings.HasSuffix(node.parent.PreviousNodeText(), " ") {
		r.writeByte(itemSpace)
	}
	r.writeString("<code>")
	return WalkStop, nil
}

func (r *VditorRenderer) renderCodeSpanContent(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *VditorRenderer) renderCodeSpanCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("</code>")
	if !strings.HasPrefix(node.parent.NextNodeText(), " ") {
		r.writeByte(itemSpace)
	}
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmphasis(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer) renderEmAsteriskOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("em", [][]string{{"data-marker", "*"}}, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmAsteriskCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/em", nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmUnderscoreOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("em", [][]string{{"data-marker", "_"}}, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderEmUnderscoreCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/em", nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrong(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *VditorRenderer) renderStrongA6kOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("strong", [][]string{{"data-marker", "**"}}, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrongA6kCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/strong", nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrongU8eOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("strong", [][]string{{"data-marker", "__"}}, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderStrongU8eCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.tag("/strong", nil, false)
	return WalkStop, nil
}

func (r *VditorRenderer) renderBlockquote(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<blockquote>")
	} else {
		r.writeString("</blockquote>")
	}
	return WalkContinue, nil
}

func (r *VditorRenderer) renderBlockquoteMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderHeading(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderHeadingC8hMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *VditorRenderer) renderList(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderListItem(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderTaskListItemMarker(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderThematicBreak(node *Node, entering bool) (WalkStatus, error) {
	r.tag("hr", nil, true)
	return WalkStop, nil
}

func (r *VditorRenderer) renderHardBreak(node *Node, entering bool) (WalkStatus, error) {
	r.tag("br", nil, true)
	return WalkStop, nil
}

func (r *VditorRenderer) renderSoftBreak(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemNewline)
	return WalkStop, nil
}

func (r *VditorRenderer) tag(name string, attrs [][]string, selfclosing bool) {
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

func (r *VditorRenderer) renderCodeBlock(node *Node, entering bool) (WalkStatus, error) {
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

func (r *VditorRenderer) renderCodeBlockCode(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if 0 < len(node.previous.codeBlockInfo) {
			infoWords := split(node.previous.codeBlockInfo, itemSpace)
			language := bytesToStr(infoWords[0])
			r.writeString("<pre><code class=\"language-" + language + "\">")
		} else {
			r.writeString("<pre><code>")
		}
		tokens := node.tokens
		if 1 > len(tokens) {
			tokens = append(tokens, itemNewline)
		}
		r.write(tokens)
		return WalkSkipChildren, nil
	}
	r.writeString("</code></pre>")
	return WalkStop, nil
}
