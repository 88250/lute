// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// VditorRenderer 描述了 Vditor DOM 渲染器。
type VditorRenderer struct {
	*BaseRenderer
}

// NewVditorRenderer 创建一个 HTML 渲染器。
func NewVditorRenderer(tree *parse.Tree) *VditorRenderer {
	ret := &VditorRenderer{BaseRenderer: NewBaseRenderer(tree)}
	ret.RendererFuncs[ast.NodeDocument] = ret.renderDocument
	ret.RendererFuncs[ast.NodeParagraph] = ret.renderParagraph
	ret.RendererFuncs[ast.NodeText] = ret.renderText
	ret.RendererFuncs[ast.NodeCodeSpan] = ret.renderCodeSpan
	ret.RendererFuncs[ast.NodeCodeSpanOpenMarker] = ret.renderCodeSpanOpenMarker
	ret.RendererFuncs[ast.NodeCodeSpanContent] = ret.renderCodeSpanContent
	ret.RendererFuncs[ast.NodeCodeSpanCloseMarker] = ret.renderCodeSpanCloseMarker
	ret.RendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock
	ret.RendererFuncs[ast.NodeCodeBlockFenceOpenMarker] = ret.renderCodeBlockOpenMarker
	ret.RendererFuncs[ast.NodeCodeBlockFenceInfoMarker] = ret.renderCodeBlockInfoMarker
	ret.RendererFuncs[ast.NodeCodeBlockCode] = ret.renderCodeBlockCode
	ret.RendererFuncs[ast.NodeCodeBlockFenceCloseMarker] = ret.renderCodeBlockCloseMarker
	ret.RendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock
	ret.RendererFuncs[ast.NodeMathBlockOpenMarker] = ret.renderMathBlockOpenMarker
	ret.RendererFuncs[ast.NodeMathBlockContent] = ret.renderMathBlockContent
	ret.RendererFuncs[ast.NodeMathBlockCloseMarker] = ret.renderMathBlockCloseMarker
	ret.RendererFuncs[ast.NodeInlineMath] = ret.renderInlineMath
	ret.RendererFuncs[ast.NodeInlineMathOpenMarker] = ret.renderInlineMathOpenMarker
	ret.RendererFuncs[ast.NodeInlineMathContent] = ret.renderInlineMathContent
	ret.RendererFuncs[ast.NodeInlineMathCloseMarker] = ret.renderInlineMathCloseMarker
	ret.RendererFuncs[ast.NodeEmphasis] = ret.renderEmphasis
	ret.RendererFuncs[ast.NodeEmA6kOpenMarker] = ret.renderEmAsteriskOpenMarker
	ret.RendererFuncs[ast.NodeEmA6kCloseMarker] = ret.renderEmAsteriskCloseMarker
	ret.RendererFuncs[ast.NodeEmU8eOpenMarker] = ret.renderEmUnderscoreOpenMarker
	ret.RendererFuncs[ast.NodeEmU8eCloseMarker] = ret.renderEmUnderscoreCloseMarker
	ret.RendererFuncs[ast.NodeStrong] = ret.renderStrong
	ret.RendererFuncs[ast.NodeStrongA6kOpenMarker] = ret.renderStrongA6kOpenMarker
	ret.RendererFuncs[ast.NodeStrongA6kCloseMarker] = ret.renderStrongA6kCloseMarker
	ret.RendererFuncs[ast.NodeStrongU8eOpenMarker] = ret.renderStrongU8eOpenMarker
	ret.RendererFuncs[ast.NodeStrongU8eCloseMarker] = ret.renderStrongU8eCloseMarker
	ret.RendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote
	ret.RendererFuncs[ast.NodeBlockquoteMarker] = ret.renderBlockquoteMarker
	ret.RendererFuncs[ast.NodeHeading] = ret.renderHeading
	ret.RendererFuncs[ast.NodeHeadingC8hMarker] = ret.renderHeadingC8hMarker
	ret.RendererFuncs[ast.NodeList] = ret.renderList
	ret.RendererFuncs[ast.NodeListItem] = ret.renderListItem
	ret.RendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	ret.RendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak
	ret.RendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak
	ret.RendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	ret.RendererFuncs[ast.NodeInlineHTML] = ret.renderInlineHTML
	ret.RendererFuncs[ast.NodeLink] = ret.renderLink
	ret.RendererFuncs[ast.NodeImage] = ret.renderImage
	ret.RendererFuncs[ast.NodeBang] = ret.renderBang
	ret.RendererFuncs[ast.NodeOpenBracket] = ret.renderOpenBracket
	ret.RendererFuncs[ast.NodeCloseBracket] = ret.renderCloseBracket
	ret.RendererFuncs[ast.NodeOpenParen] = ret.renderOpenParen
	ret.RendererFuncs[ast.NodeCloseParen] = ret.renderCloseParen
	ret.RendererFuncs[ast.NodeLinkText] = ret.renderLinkText
	ret.RendererFuncs[ast.NodeLinkSpace] = ret.renderLinkSpace
	ret.RendererFuncs[ast.NodeLinkDest] = ret.renderLinkDest
	ret.RendererFuncs[ast.NodeLinkTitle] = ret.renderLinkTitle
	ret.RendererFuncs[ast.NodeStrikethrough] = ret.renderStrikethrough
	ret.RendererFuncs[ast.NodeStrikethrough1OpenMarker] = ret.renderStrikethrough1OpenMarker
	ret.RendererFuncs[ast.NodeStrikethrough1CloseMarker] = ret.renderStrikethrough1CloseMarker
	ret.RendererFuncs[ast.NodeStrikethrough2OpenMarker] = ret.renderStrikethrough2OpenMarker
	ret.RendererFuncs[ast.NodeStrikethrough2CloseMarker] = ret.renderStrikethrough2CloseMarker
	ret.RendererFuncs[ast.NodeTaskListItemMarker] = ret.renderTaskListItemMarker
	ret.RendererFuncs[ast.NodeTable] = ret.renderTable
	ret.RendererFuncs[ast.NodeTableHead] = ret.renderTableHead
	ret.RendererFuncs[ast.NodeTableRow] = ret.renderTableRow
	ret.RendererFuncs[ast.NodeTableCell] = ret.renderTableCell
	ret.RendererFuncs[ast.NodeEmoji] = ret.renderEmoji
	ret.RendererFuncs[ast.NodeEmojiUnicode] = ret.renderEmojiUnicode
	ret.RendererFuncs[ast.NodeEmojiImg] = ret.renderEmojiImg
	ret.RendererFuncs[ast.NodeEmojiAlias] = ret.renderEmojiAlias
	ret.RendererFuncs[ast.NodeFootnotesDef] = ret.renderFootnotesDef
	ret.RendererFuncs[ast.NodeFootnotesRef] = ret.renderFootnotesRef
	ret.RendererFuncs[ast.NodeBackslash] = ret.renderBackslash
	ret.RendererFuncs[ast.NodeBackslashContent] = ret.renderBackslashContent
	return ret
}

func (r *VditorRenderer) Render() (output []byte, err error) {
	output, err = r.BaseRenderer.Render()
	if nil != err || 1 > len(r.Tree.Context.LinkRefDefs) {
		return
	}

	buf := &bytes.Buffer{}
	// 将链接引用定义添加到末尾
	buf.WriteString("<p data-block=\"0\" data-type=\"link-ref-defs\">")
	for _, node := range r.Tree.Context.LinkRefDefs {
		label := node.LinkRefLabel
		dest := node.ChildByType(ast.NodeLinkDest).Tokens
		buf.WriteString("[" + util.BytesToStr(label) + "]: " + util.BytesToStr(dest) + "\n")
	}
	buf.WriteString("</p>")
	output = append(output, buf.Bytes()...)
	return
}

func (r *VditorRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(util.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *VditorRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<span data-type=\"backslash\">")
		r.WriteString("<span>")
		r.WriteByte(lex.ItemBackslash)
		r.WriteString("</span>")
	} else {
		r.WriteString("</span>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("[" + util.BytesToStr(node.Tokens) + "]: ")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("[" + util.BytesToStr(node.Tokens) + "]")
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
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("<span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\">")
	r.tag("code", [][]string{{"data-type", "math-inline"}}, false)
	tokens := bytes.ReplaceAll(node.Tokens, []byte(parse.Zwsp), []byte(""))
	tokens = util.EscapeHTML(tokens)
	tokens = append([]byte(parse.Zwsp), tokens...)
	r.Write(tokens)
	r.WriteString("</code></span>" + parse.Zwsp)
	return ast.WalkStop
}

func (r *VditorRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, parse.Caret, "")
		if "" == previousNodeText {
			r.WriteString(parse.Zwsp)
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
	codeIsEmpty := 1 > codeLen || (len(parse.Caret) == codeLen && parse.Caret == string(node.Tokens))
	r.WriteString("<pre>")
	r.tag("code", [][]string{{"data-type", "math-block"}}, false)
	if codeIsEmpty {
		r.WriteString("<wbr>\n")
	} else {
		r.Write(util.EscapeHTML(node.Tokens))
	}
	r.WriteString("</code></pre>")
	return ast.WalkStop
}

func (r *VditorRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString(`<div class="vditor-wysiwyg__block" data-type="math-block" data-block="0">`)
	} else {
		r.WriteString("</div>")
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
		if nil == node.FirstChild {
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
		} else if bytes.Equal(node.FirstChild.Tokens, []byte(parse.Caret)) {
			node.FirstChild.Tokens = []byte(parse.Caret + " ")
		} else {
			node.FirstChild.Tokens = bytes.TrimSpace(node.FirstChild.Tokens)
		}
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
	r.Write(node.Tokens)
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
		if 0 == r.DisableTags {
			r.WriteString("<img src=\"")
			destTokens := node.ChildByType(ast.NodeLinkDest).Tokens
			destTokens = r.Tree.Context.RelativePath(destTokens)
			destTokens = bytes.ReplaceAll(destTokens, []byte(parse.Caret), []byte(""))
			r.Write(destTokens)
			r.WriteString("\" alt=\"")
			if alt := node.ChildByType(ast.NodeLinkText); nil != alt && bytes.Contains(alt.Tokens, []byte(parse.Caret)) {
				alt.Tokens = bytes.ReplaceAll(alt.Tokens, []byte(parse.Caret), []byte(""))
			}
		}
		r.DisableTags++
		return ast.WalkContinue
	}

	r.DisableTags--
	if 0 == r.DisableTags {
		r.WriteString("\"")
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			r.WriteString(" title=\"")
			title.Tokens = bytes.ReplaceAll(title.Tokens, []byte(parse.Caret), []byte(""))
			r.Write(title.Tokens)
			r.WriteString("\"")
		}
		r.WriteString(" />")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if 3 == node.LinkType {
			text := string(node.ChildByType(ast.NodeLinkText).Tokens)
			label := string(node.LinkRefLabel)
			attrs := [][]string{{"data-type", "link-ref"}, {"data-link-text", text}, {"data-link-label", label}}
			r.tag("span", attrs, false)
			r.WriteString(text)
			r.tag("/span", nil, false)
			return ast.WalkStop
		}

		dest := node.ChildByType(ast.NodeLinkDest)
		destTokens := dest.Tokens
		destTokens = r.Tree.Context.RelativePath(destTokens)
		caretInDest := bytes.Contains(destTokens, []byte(parse.Caret))
		if caretInDest {
			text := node.ChildByType(ast.NodeLinkText)
			text.Tokens = append(text.Tokens, []byte(parse.Caret)...)
			destTokens = bytes.ReplaceAll(destTokens, []byte(parse.Caret), []byte(""))
		}
		attrs := [][]string{{"href", string(destTokens)}}
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			title.Tokens = bytes.ReplaceAll(title.Tokens, []byte(parse.Caret), []byte(""))
			attrs = append(attrs, []string{"title", string(title.Tokens)})
		}
		r.tag("a", attrs, false)
	} else {
		r.tag("/a", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString(`<div class="vditor-wysiwyg__block" data-type="html-block" data-block="0">`)
	node.Tokens = bytes.TrimSpace(node.Tokens)
	r.WriteString("<pre>")
	r.tag("code", nil, false)
	r.Write(util.EscapeHTML(node.Tokens))
	r.WriteString("</code></pre></div>")
	return ast.WalkStop
}

func (r *VditorRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if bytes.Equal(node.Tokens, []byte("<br />")) && node.ParentIs(ast.NodeTableCell) {
		r.Write(node.Tokens)
		return ast.WalkStop
	}

	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, parse.Caret, "")
		if "" == previousNodeText {
			r.WriteString(parse.Zwsp)
		}
	}

	r.WriteString("<span class=\"vditor-wysiwyg__block\" data-type=\"html-inline\">")
	node.Tokens = bytes.TrimSpace(node.Tokens)
	r.tag("code", [][]string{{"data-type", "html-inline"}}, false)
	tokens := bytes.ReplaceAll(node.Tokens, []byte(parse.Zwsp), []byte(""))
	tokens = util.EscapeHTML(tokens)
	tokens = append([]byte(parse.Zwsp), tokens...)
	r.Write(tokens)
	r.WriteString("</code></span>" + parse.Zwsp)
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
		r.WriteByte(lex.ItemNewline)
		r.tag("/p", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if r.Option.AutoSpace {
		r.Space(node)
	}
	if r.Option.FixTermTypo {
		r.FixTermTypo(node)
	}
	if r.Option.ChinesePunct {
		r.ChinesePunct(node)
	}

	node.Tokens = bytes.TrimRight(node.Tokens, "\n")
	// 有的场景需要零宽空格撑起，但如果有其他文本内容的话需要把零宽空格删掉
	if !bytes.EqualFold(node.Tokens, []byte(parse.Caret+parse.Zwsp)) {
		node.Tokens = bytes.ReplaceAll(node.Tokens, []byte(parse.Zwsp), []byte(""))
	}
	r.Write(util.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *VditorRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, parse.Caret, "")
		if "" == previousNodeText {
			r.WriteString(parse.Zwsp)
		} else {
			lastc, _ := utf8.DecodeLastRuneInString(previousNodeText)
			if unicode.IsLetter(lastc) || unicode.IsDigit(lastc) {
				r.WriteByte(lex.ItemSpace)
			}
		}
		r.tag("code", [][]string{{"data-marker", strings.Repeat("`", node.CodeMarkerLen)}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	tokens := bytes.ReplaceAll(node.Tokens, []byte(parse.Zwsp), []byte(""))
	tokens = util.EscapeHTML(tokens)
	tokens = append([]byte(parse.Zwsp), tokens...)
	r.Write(tokens)
	return ast.WalkStop
}

func (r *VditorRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("</code>")
	codeSpan := node.Parent
	if codeSpanParent := codeSpan.Parent; nil != codeSpanParent && ast.NodeLink == codeSpanParent.Type {
		return ast.WalkStop
	}
	r.WriteString(parse.Zwsp)
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
		r.WriteString(`<blockquote data-block="0">`)
	} else {
		r.WriteString("</blockquote>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1] + " data-block=\"0\"")
		if !node.HeadingSetext {
			r.WriteString(" data-marker=\"#\">")
		} else {
			if 1 == node.HeadingLevel {
				r.WriteString(" data-marker=\"=\">")
			} else {
				r.WriteString(" data-marker=\"-\">")
			}
		}
		if r.Option.HeadingAnchor {
			id := r.headingID(node)
			r.tag("a", [][]string{{"id", "vditorAnchor-" + id}, {"class", "vditor-anchor"}, {"href", "#" + id}}, false)
			r.WriteString(`<svg viewBox="0 0 16 16" version="1.1" width="16" height="16"><path fill-rule="evenodd" d="M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z"></path></svg>`)
			r.tag("/a", nil, false)
		}
	} else {
		if nil != node.HeadingID {
			r.WriteString(" {" + string(node.HeadingID) + "}")
		}
		r.WriteString("</h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1] + ">")
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
				attrs = append(attrs, []string{"class", r.Option.GFMTaskListItemClass})
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
		r.Write(node.Tokens)
		r.WriteByte(lex.ItemNewline)
		r.tag("/p", nil, false)
	}
	return ast.WalkStop
}

func (r *VditorRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("br", nil, true)
	return ast.WalkStop
}

func (r *VditorRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *VditorRenderer) tag(name string, attrs [][]string, selfclosing bool) {
	if r.DisableTags > 0 {
		return
	}

	r.WriteString("<")
	r.WriteString(name)
	if 0 < len(attrs) {
		for _, attr := range attrs {
			r.WriteString(" " + attr[0] + "=\"" + attr[1] + "\"")
		}
	}
	if selfclosing {
		r.WriteString(" /")
	}
	r.WriteString(">")
}

func (r *VditorRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		marker := "```"
		if nil != node.FirstChild {
			marker = string(node.FirstChild.Tokens)
		}
		r.WriteString(`<div class="vditor-wysiwyg__block" data-type="code-block" data-block="0" data-marker="` + marker + `">`)
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	codeLen := len(node.Tokens)
	codeIsEmpty := 1 > codeLen || (len(parse.Caret) == codeLen && parse.Caret == string(node.Tokens))
	isFenced := node.Parent.IsFencedCodeBlock
	if isFenced {
		node.Previous.CodeBlockInfo = bytes.ReplaceAll(node.Previous.CodeBlockInfo, []byte(parse.Caret), []byte(""))
	}
	var attrs [][]string
	if isFenced && 0 < len(node.Previous.CodeBlockInfo) {
		infoWords := lex.Split(node.Previous.CodeBlockInfo, lex.ItemSpace)
		language := string(infoWords[0])
		attrs = append(attrs, []string{"class", "language-" + language})
	}
	r.WriteString("<pre>")
	r.tag("code", attrs, false)

	if codeIsEmpty {
		r.WriteString("<wbr>\n")
	} else {
		r.Write(util.EscapeHTML(node.Tokens))
		r.Newline()
	}
	r.WriteString("</code></pre>")
	return ast.WalkStop
}
