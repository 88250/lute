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
	"github.com/88250/lute/ast"
	"github.com/88250/lute/util"
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
	ret.rendererFuncs[ast.NodeDocument] = ret.renderDocument
	ret.rendererFuncs[ast.NodeParagraph] = ret.renderParagraph
	ret.rendererFuncs[ast.NodeText] = ret.renderText
	ret.rendererFuncs[ast.NodeCodeSpan] = ret.renderCodeSpan
	ret.rendererFuncs[ast.NodeCodeSpanOpenMarker] = ret.renderCodeSpanOpenMarker
	ret.rendererFuncs[ast.NodeCodeSpanContent] = ret.renderCodeSpanContent
	ret.rendererFuncs[ast.NodeCodeSpanCloseMarker] = ret.renderCodeSpanCloseMarker
	ret.rendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock
	ret.rendererFuncs[ast.NodeCodeBlockFenceOpenMarker] = ret.renderCodeBlockOpenMarker
	ret.rendererFuncs[ast.NodeCodeBlockFenceInfoMarker] = ret.renderCodeBlockInfoMarker
	ret.rendererFuncs[ast.NodeCodeBlockCode] = ret.renderCodeBlockCode
	ret.rendererFuncs[ast.NodeCodeBlockFenceCloseMarker] = ret.renderCodeBlockCloseMarker
	ret.rendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock
	ret.rendererFuncs[ast.NodeMathBlockOpenMarker] = ret.renderMathBlockOpenMarker
	ret.rendererFuncs[ast.NodeMathBlockContent] = ret.renderMathBlockContent
	ret.rendererFuncs[ast.NodeMathBlockCloseMarker] = ret.renderMathBlockCloseMarker
	ret.rendererFuncs[ast.NodeInlineMath] = ret.renderInlineMath
	ret.rendererFuncs[ast.NodeInlineMathOpenMarker] = ret.renderInlineMathOpenMarker
	ret.rendererFuncs[ast.NodeInlineMathContent] = ret.renderInlineMathContent
	ret.rendererFuncs[ast.NodeInlineMathCloseMarker] = ret.renderInlineMathCloseMarker
	ret.rendererFuncs[ast.NodeEmphasis] = ret.renderEmphasis
	ret.rendererFuncs[ast.NodeEmA6kOpenMarker] = ret.renderEmAsteriskOpenMarker
	ret.rendererFuncs[ast.NodeEmA6kCloseMarker] = ret.renderEmAsteriskCloseMarker
	ret.rendererFuncs[ast.NodeEmU8eOpenMarker] = ret.renderEmUnderscoreOpenMarker
	ret.rendererFuncs[ast.NodeEmU8eCloseMarker] = ret.renderEmUnderscoreCloseMarker
	ret.rendererFuncs[ast.NodeStrong] = ret.renderStrong
	ret.rendererFuncs[ast.NodeStrongA6kOpenMarker] = ret.renderStrongA6kOpenMarker
	ret.rendererFuncs[ast.NodeStrongA6kCloseMarker] = ret.renderStrongA6kCloseMarker
	ret.rendererFuncs[ast.NodeStrongU8eOpenMarker] = ret.renderStrongU8eOpenMarker
	ret.rendererFuncs[ast.NodeStrongU8eCloseMarker] = ret.renderStrongU8eCloseMarker
	ret.rendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote
	ret.rendererFuncs[ast.NodeBlockquoteMarker] = ret.renderBlockquoteMarker
	ret.rendererFuncs[ast.NodeHeading] = ret.renderHeading
	ret.rendererFuncs[ast.NodeHeadingC8hMarker] = ret.renderHeadingC8hMarker
	ret.rendererFuncs[ast.NodeList] = ret.renderList
	ret.rendererFuncs[ast.NodeListItem] = ret.renderListItem
	ret.rendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	ret.rendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak
	ret.rendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak
	ret.rendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	ret.rendererFuncs[ast.NodeInlineHTML] = ret.renderInlineHTML
	ret.rendererFuncs[ast.NodeLink] = ret.renderLink
	ret.rendererFuncs[ast.NodeImage] = ret.renderImage
	ret.rendererFuncs[ast.NodeBang] = ret.renderBang
	ret.rendererFuncs[ast.NodeOpenBracket] = ret.renderOpenBracket
	ret.rendererFuncs[ast.NodeCloseBracket] = ret.renderCloseBracket
	ret.rendererFuncs[ast.NodeOpenParen] = ret.renderOpenParen
	ret.rendererFuncs[ast.NodeCloseParen] = ret.renderCloseParen
	ret.rendererFuncs[ast.NodeLinkText] = ret.renderLinkText
	ret.rendererFuncs[ast.NodeLinkSpace] = ret.renderLinkSpace
	ret.rendererFuncs[ast.NodeLinkDest] = ret.renderLinkDest
	ret.rendererFuncs[ast.NodeLinkTitle] = ret.renderLinkTitle
	ret.rendererFuncs[ast.NodeStrikethrough] = ret.renderStrikethrough
	ret.rendererFuncs[ast.NodeStrikethrough1OpenMarker] = ret.renderStrikethrough1OpenMarker
	ret.rendererFuncs[ast.NodeStrikethrough1CloseMarker] = ret.renderStrikethrough1CloseMarker
	ret.rendererFuncs[ast.NodeStrikethrough2OpenMarker] = ret.renderStrikethrough2OpenMarker
	ret.rendererFuncs[ast.NodeStrikethrough2CloseMarker] = ret.renderStrikethrough2CloseMarker
	ret.rendererFuncs[ast.NodeTaskListItemMarker] = ret.renderTaskListItemMarker
	ret.rendererFuncs[ast.NodeTable] = ret.renderTable
	ret.rendererFuncs[ast.NodeTableHead] = ret.renderTableHead
	ret.rendererFuncs[ast.NodeTableRow] = ret.renderTableRow
	ret.rendererFuncs[ast.NodeTableCell] = ret.renderTableCell
	ret.rendererFuncs[ast.NodeEmoji] = ret.renderEmoji
	ret.rendererFuncs[ast.NodeEmojiUnicode] = ret.renderEmojiUnicode
	ret.rendererFuncs[ast.NodeEmojiImg] = ret.renderEmojiImg
	ret.rendererFuncs[ast.NodeEmojiAlias] = ret.renderEmojiAlias
	ret.rendererFuncs[ast.NodeFootnotesDef] = ret.renderFootnotesDef
	ret.rendererFuncs[ast.NodeFootnotesRef] = ret.renderFootnotesRef
	ret.rendererFuncs[ast.NodeBackslash] = ret.renderBackslash
	ret.rendererFuncs[ast.NodeBackslashContent] = ret.renderBackslashContent
	return ret
}

func (r *VditorRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(escapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *VditorRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.writeString("<span data-type=\"backslash\">")
		r.writeString("<span>")
		r.writeByte(itemBackslash)
		r.writeString("</span>")
	} else {
		r.writeString("</span>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.writeString("[" + util.BytesToStr(node.Tokens) + "]: ")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("[" + util.BytesToStr(node.Tokens) + "]")
	return ast.WalkStop
}

func (r *VditorRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("<span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\">")
	r.tag("code", [][]string{{"data-type", "math-inline"}}, false)
	tokens := bytes.ReplaceAll(node.Tokens, []byte(zwsp), []byte(""))
	tokens = escapeHTML(tokens)
	tokens = append([]byte(zwsp), tokens...)
	r.write(tokens)
	r.writeString("</code></span>" + zwsp)
	return ast.WalkStop
}

func (r *VditorRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, caret, "")
		if "" == previousNodeText {
			r.writeString(zwsp)
		}
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
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
	return ast.WalkStop
}

func (r *VditorRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.writeString(`<div class="vditor-wysiwyg__block" data-type="math-block" data-block="0">`)
	} else {
		r.writeString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
	tag := "td"
	if ast.NodeTableHead == node.Parent.Parent.Type {
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
	return ast.WalkContinue
}

func (r *VditorRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("tr", nil, false)
	} else {
		r.tag("/tr", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("thead", nil, false)
	} else {
		r.tag("/thead", nil, false)
		if nil != node.Next {
			r.tag("tbody", nil, false)
		}
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("table", [][]string{{"data-block", "0"}}, false)
	} else {
		if nil != node.FirstChild.Next {
			r.tag("/tbody", nil, false)
		}
		r.tag("/table", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("s", [][]string{{"data-marker", "~"}}, false)
	return ast.WalkStop
}

func (r *VditorRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/s", nil, false)
	return ast.WalkStop
}

func (r *VditorRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("s", [][]string{{"data-marker", "~~"}}, false)
	return ast.WalkStop
}

func (r *VditorRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/s", nil, false)
	return ast.WalkStop
}

func (r *VditorRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if 0 == r.disableTags {
			r.writeString("<img src=\"")
			destTokens := node.ChildByType(ast.NodeLinkDest).Tokens
			destTokens = r.tree.context.relativePath(destTokens)
			destTokens = bytes.ReplaceAll(destTokens, []byte(caret), []byte(""))
			r.write(destTokens)
			r.writeString("\" alt=\"")
			if alt := node.ChildByType(ast.NodeLinkText); nil != alt && bytes.Contains(alt.Tokens, []byte(caret)) {
				alt.Tokens = bytes.ReplaceAll(alt.Tokens, []byte(caret), []byte(""))
			}
		}
		r.disableTags++
		return ast.WalkContinue
	}

	r.disableTags--
	if 0 == r.disableTags {
		r.writeString("\"")
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			r.writeString(" title=\"")
			title.Tokens = bytes.ReplaceAll(title.Tokens, []byte(caret), []byte(""))
			r.write(title.Tokens)
			r.writeString("\"")
		}
		r.writeString(" />")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		dest := node.ChildByType(ast.NodeLinkDest)
		destTokens := dest.Tokens
		destTokens = r.tree.context.relativePath(destTokens)
		caretInDest := bytes.Contains(destTokens, []byte(caret))
		if caretInDest {
			text := node.ChildByType(ast.NodeLinkText)
			text.Tokens = append(text.Tokens, []byte(caret)...)
			destTokens = bytes.ReplaceAll(destTokens, []byte(caret), []byte(""))
		}
		attrs := [][]string{{"href", string(destTokens)}}
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			title.Tokens = bytes.ReplaceAll(title.Tokens, []byte(caret), []byte(""))
			attrs = append(attrs, []string{"title", string(title.Tokens)})
		}
		r.tag("a", attrs, false)
	} else {
		r.tag("/a", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString(`<div class="vditor-wysiwyg__block" data-type="html-block" data-block="0">`)
	node.Tokens = bytes.TrimSpace(node.Tokens)
	r.writeString("<pre>")
	r.tag("code", nil, false)
	r.write(escapeHTML(node.Tokens))
	r.writeString("</code></pre></div>")
	return ast.WalkStop
}

func (r *VditorRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if bytes.Equal(node.Tokens, []byte("<br />")) && node.ParentIs(ast.NodeTableCell) {
		r.write(node.Tokens)
		return ast.WalkStop
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
	return ast.WalkStop
}

func (r *VditorRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if grandparent := node.Parent.Parent; nil != grandparent && ast.NodeList == grandparent.Type && grandparent.Tight { // List.ListItem.Paragraph
		return ast.WalkContinue
	}

	if entering {
		r.tag("p", [][]string{{"data-block", "0"}}, false)
	} else {
		r.writeByte(itemNewline)
		r.tag("/p", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
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
	return ast.WalkStop
}

func (r *VditorRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
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
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	tokens := bytes.ReplaceAll(node.Tokens, []byte(zwsp), []byte(""))
	tokens = escapeHTML(tokens)
	tokens = append([]byte(zwsp), tokens...)
	r.write(tokens)
	return ast.WalkStop
}

func (r *VditorRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("</code>")
	codeSpan := node.Parent
	if codeSpanParent := codeSpan.Parent; nil != codeSpanParent && ast.NodeLink == codeSpanParent.Type {
		return ast.WalkStop
	}
	r.writeString(zwsp)
	return ast.WalkStop
}

func (r *VditorRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("em", [][]string{{"data-marker", "*"}}, false)
	return ast.WalkStop
}

func (r *VditorRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/em", nil, false)
	return ast.WalkStop
}

func (r *VditorRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("em", [][]string{{"data-marker", "_"}}, false)
	return ast.WalkStop
}

func (r *VditorRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/em", nil, false)
	return ast.WalkStop
}

func (r *VditorRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("strong", [][]string{{"data-marker", "**"}}, false)
	return ast.WalkStop
}

func (r *VditorRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/strong", nil, false)
	return ast.WalkStop
}

func (r *VditorRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("strong", [][]string{{"data-marker", "__"}}, false)
	return ast.WalkStop
}

func (r *VditorRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/strong", nil, false)
	return ast.WalkStop
}

func (r *VditorRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.writeString(`<blockquote data-block="0">`)
	} else {
		r.writeString("</blockquote>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
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
	return ast.WalkContinue
}

func (r *VditorRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
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
	return ast.WalkContinue
}

func (r *VditorRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
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
			if nil != node.FirstChild && nil != node.FirstChild.FirstChild && ast.NodeTaskListItemMarker == node.FirstChild.FirstChild.Type {
				attrs = append(attrs, []string{"class", r.option.GFMTaskListItemClass})
			}
		}
		r.tag("li", attrs, false)
	} else {
		r.tag("/li", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	var attrs [][]string
	if node.TaskListItemChecked {
		attrs = append(attrs, []string{"checked", ""})
	}
	attrs = append(attrs, []string{"type", "checkbox"})
	r.tag("input", attrs, true)
	return ast.WalkStop
}

func (r *VditorRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("hr", [][]string{{"data-block", "0"}}, true)
	if nil != node.Tokens {
		r.tag("p", [][]string{{"data-block", "0"}}, false)
		r.writeBytes(node.Tokens)
		r.writeByte(itemNewline)
		r.tag("/p", nil, false)
	}
	return ast.WalkStop
}

func (r *VditorRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("br", nil, true)
	return ast.WalkStop
}

func (r *VditorRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(itemNewline)
	return ast.WalkStop
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

func (r *VditorRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		marker := "```"
		if nil != node.FirstChild {
			marker = string(node.FirstChild.Tokens)
		}
		r.writeString(`<div class="vditor-wysiwyg__block" data-type="code-block" data-block="0" data-marker="` + marker + `">`)
	} else {
		r.writeString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
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
	return ast.WalkStop
}
