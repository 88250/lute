// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lute

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// VditorRenderer 描述了 Vditor DOM 渲染器。
type VditorRenderer struct {
	*BaseRenderer
}

// newVditorRenderer 创建一个 HTML 渲染器。
func (lute *Lute) newVditorRenderer(tree *Tree) *VditorRenderer {
	ret := &VditorRenderer{BaseRenderer: lute.newBaseRenderer(tree)}
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
	ret.rendererFuncs[NodeFootnotesDef] = ret.renderFootnotesDef
	ret.rendererFuncs[NodeFootnotesRef] = ret.renderFootnotesRef
	ret.rendererFuncs[NodeBackslash] = ret.renderBackslash
	ret.rendererFuncs[NodeBackslashContent] = ret.renderBackslashContent
	return ret
}

func (r *VditorRenderer) renderBackslashContent(node *Node, entering bool) WalkStatus {
	r.write(escapeHTML(node.Tokens))
	return WalkStop
}

func (r *VditorRenderer) renderBackslash(node *Node, entering bool) WalkStatus {
	if entering {
		r.writeString("<span data-type=\"backslash\">")
		r.writeString("<span>")
		r.writeByte(itemBackslash)
		r.writeString("</span>")
	} else {
		r.writeString("</span>")
	}
	return WalkContinue
}

func (r *VditorRenderer) renderFootnotesDef(node *Node, entering bool) WalkStatus {
	if entering {
		r.writeString("[" + bytesToStr(node.Tokens) + "]: ")
	}
	return WalkContinue
}

func (r *VditorRenderer) renderFootnotesRef(node *Node, entering bool) WalkStatus {
	r.writeString("[" + bytesToStr(node.Tokens) + "]")
	return WalkStop
}

func (r *VditorRenderer) renderCodeBlockCloseMarker(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderCodeBlockInfoMarker(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderCodeBlockOpenMarker(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderEmojiAlias(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderEmojiImg(node *Node, entering bool) WalkStatus {
	r.write(node.Tokens)
	return WalkStop
}

func (r *VditorRenderer) renderEmojiUnicode(node *Node, entering bool) WalkStatus {
	r.write(node.Tokens)
	return WalkStop
}

func (r *VditorRenderer) renderEmoji(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *VditorRenderer) renderInlineMathCloseMarker(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderInlineMathContent(node *Node, entering bool) WalkStatus {
	r.writeString("<span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\">")
	r.tag("code", [][]string{{"data-type", "math-inline"}}, false)
	tokens := bytes.ReplaceAll(node.Tokens, []byte(zwsp), []byte(""))
	tokens = escapeHTML(tokens)
	tokens = append([]byte(zwsp), tokens...)
	r.write(tokens)
	r.writeString("</code></span>" + zwsp)
	return WalkStop
}

func (r *VditorRenderer) renderInlineMathOpenMarker(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderInlineMath(node *Node, entering bool) WalkStatus {
	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, caret, "")
		if "" == previousNodeText {
			r.writeString(zwsp)
		}
	}
	return WalkContinue
}

func (r *VditorRenderer) renderMathBlockCloseMarker(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderMathBlockContent(node *Node, entering bool) WalkStatus {
	node.Tokens = bytes.TrimSpace(node.Tokens)
	codeLen := len(node.Tokens)
	codeIsEmpty := 1 > codeLen || (len(caret) == codeLen && caret == string(node.Tokens))
	r.writeString("<pre>")
	r.tag("code", [][]string{{"data-type", "math-block"}}, false)
	if codeIsEmpty {
		r.writeString("<wbr>\n")
	} else {
		r.write(escapeHTML(node.Tokens))
	}
	r.writeString("</code></pre>")
	return WalkStop
}

func (r *VditorRenderer) renderMathBlockOpenMarker(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderMathBlock(node *Node, entering bool) WalkStatus {
	if entering {
		r.writeString(`<div class="vditor-wysiwyg__block" data-type="math-block" data-block="0">`)
	} else {
		r.writeString("</div>")
	}
	return WalkContinue
}

func (r *VditorRenderer) renderTableCell(node *Node, entering bool) WalkStatus {
	tag := "td"
	if NodeTableHead == node.Parent.Parent.Type {
		tag = "th"
	}
	if entering {
		var attrs [][]string
		switch node.TableCellAlign {
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
	return WalkContinue
}

func (r *VditorRenderer) renderTableRow(node *Node, entering bool) WalkStatus {
	if entering {
		r.tag("tr", nil, false)
	} else {
		r.tag("/tr", nil, false)
	}
	return WalkContinue
}

func (r *VditorRenderer) renderTableHead(node *Node, entering bool) WalkStatus {
	if entering {
		r.tag("thead", nil, false)
	} else {
		r.tag("/thead", nil, false)
		if nil != node.Next {
			r.tag("tbody", nil, false)
		}
	}
	return WalkContinue
}

func (r *VditorRenderer) renderTable(node *Node, entering bool) WalkStatus {
	if entering {
		r.tag("table", [][]string{{"data-block", "0"}}, false)
	} else {
		if nil != node.FirstChild.Next {
			r.tag("/tbody", nil, false)
		}
		r.tag("/table", nil, false)
	}
	return WalkContinue
}

func (r *VditorRenderer) renderStrikethrough(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *VditorRenderer) renderStrikethrough1OpenMarker(node *Node, entering bool) WalkStatus {
	r.tag("s", [][]string{{"data-marker", "~"}}, false)
	return WalkStop
}

func (r *VditorRenderer) renderStrikethrough1CloseMarker(node *Node, entering bool) WalkStatus {
	r.tag("/s", nil, false)
	return WalkStop
}

func (r *VditorRenderer) renderStrikethrough2OpenMarker(node *Node, entering bool) WalkStatus {
	r.tag("s", [][]string{{"data-marker", "~~"}}, false)
	return WalkStop
}

func (r *VditorRenderer) renderStrikethrough2CloseMarker(node *Node, entering bool) WalkStatus {
	r.tag("/s", nil, false)
	return WalkStop
}

func (r *VditorRenderer) renderLinkTitle(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderLinkDest(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderLinkSpace(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderLinkText(node *Node, entering bool) WalkStatus {
	r.write(node.Tokens)
	return WalkStop
}

func (r *VditorRenderer) renderCloseParen(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderOpenParen(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderCloseBracket(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderOpenBracket(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderBang(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderImage(node *Node, entering bool) WalkStatus {
	if entering {
		if 0 == r.disableTags {
			r.writeString("<img src=\"")
			destTokens := node.ChildByType(NodeLinkDest).Tokens
			destTokens = r.tree.context.relativePath(destTokens)
			destTokens = bytes.ReplaceAll(destTokens, []byte(caret), []byte(""))
			r.write(destTokens)
			r.writeString("\" alt=\"")
			if alt := node.ChildByType(NodeLinkText); nil != alt && bytes.Contains(alt.Tokens, []byte(caret)) {
				alt.Tokens = bytes.ReplaceAll(alt.Tokens, []byte(caret), []byte(""))
			}
		}
		r.disableTags++
		return WalkContinue
	}

	r.disableTags--
	if 0 == r.disableTags {
		r.writeString("\"")
		if title := node.ChildByType(NodeLinkTitle); nil != title && nil != title.Tokens {
			r.writeString(" title=\"")
			title.Tokens = bytes.ReplaceAll(title.Tokens, []byte(caret), []byte(""))
			r.write(title.Tokens)
			r.writeString("\"")
		}
		r.writeString(" />")
	}
	return WalkContinue
}

func (r *VditorRenderer) renderLink(node *Node, entering bool) WalkStatus {
	if entering {
		dest := node.ChildByType(NodeLinkDest)
		destTokens := dest.Tokens
		destTokens = r.tree.context.relativePath(destTokens)
		caretInDest := bytes.Contains(destTokens, []byte(caret))
		if caretInDest {
			text := node.ChildByType(NodeLinkText)
			text.Tokens = append(text.Tokens, []byte(caret)...)
			destTokens = bytes.ReplaceAll(destTokens, []byte(caret), []byte(""))
		}
		attrs := [][]string{{"href", string(destTokens)}}
		if title := node.ChildByType(NodeLinkTitle); nil != title && nil != title.Tokens {
			title.Tokens = bytes.ReplaceAll(title.Tokens, []byte(caret), []byte(""))
			attrs = append(attrs, []string{"title", string(title.Tokens)})
		}
		r.tag("a", attrs, false)
	} else {
		r.tag("/a", nil, false)
	}
	return WalkContinue
}

func (r *VditorRenderer) renderHTML(node *Node, entering bool) WalkStatus {
	r.writeString(`<div class="vditor-wysiwyg__block" data-type="html-block" data-block="0">`)
	node.Tokens = bytes.TrimSpace(node.Tokens)
	r.writeString("<pre>")
	r.tag("code", nil, false)
	r.write(escapeHTML(node.Tokens))
	r.writeString("</code></pre></div>")
	return WalkStop
}

func (r *VditorRenderer) renderInlineHTML(node *Node, entering bool) WalkStatus {
	if bytes.Equal(node.Tokens, []byte("<br />")) && node.parentIs(NodeTableCell) {
		r.write(node.Tokens)
		return WalkStop
	}

	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, caret, "")
		if "" == previousNodeText {
			r.writeString(zwsp)
		}
	}

	r.writeString("<span class=\"vditor-wysiwyg__block\" data-type=\"html-inline\">")
	node.Tokens = bytes.TrimSpace(node.Tokens)
	r.tag("code", [][]string{{"data-type", "html-inline"}}, false)
	tokens := bytes.ReplaceAll(node.Tokens, []byte(zwsp), []byte(""))
	tokens = escapeHTML(tokens)
	tokens = append([]byte(zwsp), tokens...)
	r.write(tokens)
	r.writeString("</code></span>" + zwsp)
	return WalkStop
}

func (r *VditorRenderer) renderDocument(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *VditorRenderer) renderParagraph(node *Node, entering bool) WalkStatus {
	if grandparent := node.Parent.Parent; nil != grandparent && NodeList == grandparent.Type && grandparent.Tight { // List.ListItem.Paragraph
		return WalkContinue
	}

	if entering {
		r.tag("p", [][]string{{"data-block", "0"}}, false)
	} else {
		r.writeByte(itemNewline)
		r.tag("/p", nil, false)
	}
	return WalkContinue
}

func (r *VditorRenderer) renderText(node *Node, entering bool) WalkStatus {
	if r.option.AutoSpace {
		r.space(node)
	}
	if r.option.FixTermTypo {
		r.fixTermTypo(node)
	}
	if r.option.ChinesePunct {
		r.chinesePunct(node)
	}

	node.Tokens = bytes.TrimRight(node.Tokens, "\n")
	// 有的场景需要零宽空格撑起，但如果有其他文本内容的话需要把零宽空格删掉
	if !bytes.EqualFold(node.Tokens, []byte(caret+zwsp)) {
		node.Tokens = bytes.ReplaceAll(node.Tokens, []byte(zwsp), []byte(""))
	}
	r.write(escapeHTML(node.Tokens))
	return WalkStop
}

func (r *VditorRenderer) renderCodeSpan(node *Node, entering bool) WalkStatus {
	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, caret, "")
		if "" == previousNodeText {
			r.writeString(zwsp)
		} else {
			lastc, _ := utf8.DecodeLastRuneInString(previousNodeText)
			if unicode.IsLetter(lastc) || unicode.IsDigit(lastc) {
				r.writeByte(itemSpace)
			}
		}
		r.tag("code", [][]string{{"marker", strings.Repeat("`", node.CodeMarkerLen)}}, false)
	}
	return WalkContinue
}

func (r *VditorRenderer) renderCodeSpanOpenMarker(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderCodeSpanContent(node *Node, entering bool) WalkStatus {
	tokens := bytes.ReplaceAll(node.Tokens, []byte(zwsp), []byte(""))
	tokens = escapeHTML(tokens)
	tokens = append([]byte(zwsp), tokens...)
	r.write(tokens)
	return WalkStop
}

func (r *VditorRenderer) renderCodeSpanCloseMarker(node *Node, entering bool) WalkStatus {
	r.writeString("</code>")
	codeSpan := node.Parent
	if codeSpanParent := codeSpan.Parent; nil != codeSpanParent && NodeLink == codeSpanParent.Type {
		return WalkStop
	}
	r.writeString(zwsp)
	return WalkStop
}

func (r *VditorRenderer) renderEmphasis(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *VditorRenderer) renderEmAsteriskOpenMarker(node *Node, entering bool) WalkStatus {
	r.tag("em", [][]string{{"data-marker", "*"}}, false)
	return WalkStop
}

func (r *VditorRenderer) renderEmAsteriskCloseMarker(node *Node, entering bool) WalkStatus {
	r.tag("/em", nil, false)
	return WalkStop
}

func (r *VditorRenderer) renderEmUnderscoreOpenMarker(node *Node, entering bool) WalkStatus {
	r.tag("em", [][]string{{"data-marker", "_"}}, false)
	return WalkStop
}

func (r *VditorRenderer) renderEmUnderscoreCloseMarker(node *Node, entering bool) WalkStatus {
	r.tag("/em", nil, false)
	return WalkStop
}

func (r *VditorRenderer) renderStrong(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *VditorRenderer) renderStrongA6kOpenMarker(node *Node, entering bool) WalkStatus {
	r.tag("strong", [][]string{{"data-marker", "**"}}, false)
	return WalkStop
}

func (r *VditorRenderer) renderStrongA6kCloseMarker(node *Node, entering bool) WalkStatus {
	r.tag("/strong", nil, false)
	return WalkStop
}

func (r *VditorRenderer) renderStrongU8eOpenMarker(node *Node, entering bool) WalkStatus {
	r.tag("strong", [][]string{{"data-marker", "__"}}, false)
	return WalkStop
}

func (r *VditorRenderer) renderStrongU8eCloseMarker(node *Node, entering bool) WalkStatus {
	r.tag("/strong", nil, false)
	return WalkStop
}

func (r *VditorRenderer) renderBlockquote(node *Node, entering bool) WalkStatus {
	if entering {
		r.writeString(`<blockquote data-block="0">`)
	} else {
		r.writeString("</blockquote>")
	}
	return WalkContinue
}

func (r *VditorRenderer) renderBlockquoteMarker(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderHeading(node *Node, entering bool) WalkStatus {
	if entering {
		r.writeString("<h" + " 123456"[node.HeadingLevel:node.HeadingLevel+1] + " data-block=\"0\">")
		if r.option.HeadingAnchor {
			anchor := node.Text()
			anchor = strings.ReplaceAll(anchor, " ", "-")
			anchor = strings.ReplaceAll(anchor, ".", "")
			r.tag("a", [][]string{{"id", "vditorAnchor-" + anchor}, {"class", "vditor-anchor"}, {"href", "#" + anchor}}, false)
			r.writeString(`<svg viewBox="0 0 16 16" version="1.1" width="16" height="16"><path fill-rule="evenodd" d="M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z"></path></svg>`)
			r.tag("/a", nil, false)
		}
	} else {
		r.writeString("</h" + " 123456"[node.HeadingLevel:node.HeadingLevel+1] + ">")
	}
	return WalkContinue
}

func (r *VditorRenderer) renderHeadingC8hMarker(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *VditorRenderer) renderList(node *Node, entering bool) WalkStatus {
	tag := "ul"
	if 1 == node.ListData.Typ || (3 == node.ListData.Typ && 0 == node.ListData.BulletChar) {
		tag = "ol"
	}
	if entering {
		var attrs [][]string
		if node.Tight {
			attrs = append(attrs, []string{"data-tight", "true"})
		}
		if 0 == node.BulletChar {
			if 1 != node.Start {
				attrs = append(attrs, []string{"start", strconv.Itoa(node.Start)})
			}
		} else {
			attrs = append(attrs, []string{"data-marker", string(node.BulletChar)})
		}
		attrs = append(attrs, []string{"data-block", "0"})
		r.tag(tag, attrs, false)
	} else {
		r.tag("/"+tag, nil, false)
	}
	return WalkContinue
}

func (r *VditorRenderer) renderListItem(node *Node, entering bool) WalkStatus {
	if entering {
		var attrs [][]string
		switch node.ListData.Typ {
		case 0:
			attrs = append(attrs, []string{"data-marker", string(node.Marker)})
		case 1:
			attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.Num) + string(node.ListData.Delimiter)})
		case 3:
			if 0 == node.ListData.BulletChar {
				attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.Num) + string(node.ListData.Delimiter)})
			} else {
				attrs = append(attrs, []string{"data-marker", string(node.Marker)})
			}
			if nil != node.FirstChild && nil != node.FirstChild.FirstChild && NodeTaskListItemMarker == node.FirstChild.FirstChild.Type {
				attrs = append(attrs, []string{"class", r.option.GFMTaskListItemClass})
			}
		}
		r.tag("li", attrs, false)
	} else {
		r.tag("/li", nil, false)
	}
	return WalkContinue
}

func (r *VditorRenderer) renderTaskListItemMarker(node *Node, entering bool) WalkStatus {
	var attrs [][]string
	if node.TaskListItemChecked {
		attrs = append(attrs, []string{"checked", ""})
	}
	attrs = append(attrs, []string{"type", "checkbox"})
	r.tag("input", attrs, true)
	return WalkStop
}

func (r *VditorRenderer) renderThematicBreak(node *Node, entering bool) WalkStatus {
	r.tag("hr", [][]string{{"data-block", "0"}}, true)
	if nil != node.Tokens {
		r.tag("p", [][]string{{"data-block", "0"}}, false)
		r.writeBytes(node.Tokens)
		r.writeByte(itemNewline)
		r.tag("/p", nil, false)
	}
	return WalkStop
}

func (r *VditorRenderer) renderHardBreak(node *Node, entering bool) WalkStatus {
	r.tag("br", nil, true)
	return WalkStop
}

func (r *VditorRenderer) renderSoftBreak(node *Node, entering bool) WalkStatus {
	r.writeByte(itemNewline)
	return WalkStop
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

func (r *VditorRenderer) renderCodeBlock(node *Node, entering bool) WalkStatus {
	if entering {
		marker := "```"
		if nil != node.FirstChild {
			marker = string(node.FirstChild.Tokens)
		}
		r.writeString(`<div class="vditor-wysiwyg__block" data-type="code-block" data-block="0" data-marker="` + marker + `">`)
	} else {
		r.writeString("</div>")
	}
	return WalkContinue
}

func (r *VditorRenderer) renderCodeBlockCode(node *Node, entering bool) WalkStatus {
	codeLen := len(node.Tokens)
	codeIsEmpty := 1 > codeLen || (len(caret) == codeLen && caret == string(node.Tokens))
	isFenced := node.Parent.IsFencedCodeBlock
	if isFenced {
		node.Previous.CodeBlockInfo = bytes.ReplaceAll(node.Previous.CodeBlockInfo, []byte(caret), []byte(""))
	}
	var attrs [][]string
	if isFenced && 0 < len(node.Previous.CodeBlockInfo) {
		infoWords := split(node.Previous.CodeBlockInfo, itemSpace)
		language := string(infoWords[0])
		attrs = append(attrs, []string{"class", "language-" + language})
	}
	r.writeString("<pre>")
	r.tag("code", attrs, false)

	if codeIsEmpty {
		r.writeString("<wbr>\n")
	} else {
		r.write(escapeHTML(node.Tokens))
		r.newline()
	}
	r.writeString("</code></pre>")
	return WalkStop
}
