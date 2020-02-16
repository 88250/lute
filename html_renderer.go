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
	"github.com/88250/lute/ast"
	"github.com/88250/lute/util"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// HTMLRenderer 描述了 HTML 渲染器。
type HTMLRenderer struct {
	*BaseRenderer
	needRenderFootnotesDef bool
	headingCnt             int
}

// newHTMLRenderer 创建一个 HTML 渲染器。
func (lute *Lute) newHTMLRenderer(tree *Tree) Renderer {
	ret := &HTMLRenderer{lute.newBaseRenderer(tree), false, 0}
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
	ret.rendererFuncs[ast.NodeToC] = ret.renderToC
	ret.rendererFuncs[ast.NodeBackslash] = ret.renderBackslash
	ret.rendererFuncs[ast.NodeBackslashContent] = ret.renderBackslashContent
	return ret
}

func (r *HTMLRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(escapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HTMLRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	headings := r.headings()
	length := len(headings)
	if 1 > length {
		return ast.WalkStop
	}
	r.writeString("<div class=\"toc-div\">")
	for i, heading := range headings {
		level := strconv.Itoa(heading.HeadingLevel)
		spaces := (heading.HeadingLevel - 1) * 2
		r.writeString(strings.Repeat("&emsp;", spaces))
		r.writeString("<span class=\"toc-h" + level + "\">")
		r.writeString("<a class=\"toc-a\" href=\"#toc_h" + level + "_" + strconv.Itoa(i) + "\">" + heading.Text() + "</a></span><br>")
	}
	r.writeString("</div>\n\n")

	return ast.WalkStop
}

func (r *HTMLRenderer) headings() (ret []*ast.Node) {
	for n := r.tree.Root.FirstChild; nil != n; n = n.Next {
		r.headings0(n, &ret)
	}
	return
}

func (r *HTMLRenderer) headings0(n *ast.Node, headings *[]*ast.Node) {
	if ast.NodeHeading == n.Type {
		*headings = append(*headings, n)
		return
	}
	if ast.NodeList == n.Type || ast.NodeListItem == n.Type || ast.NodeBlockquote == n.Type {
		for c := n.FirstChild; nil != c; c = c.Next {
			r.headings0(c, headings)
		}
	}
}

func (r *HTMLRenderer) renderFootnotesDefs(lute *Lute, context *Context) []byte {
	r.writeString("<div class=\"footnotes-defs-div\">")
	r.writeString("<hr class=\"footnotes-defs-hr\" />\n")
	r.writeString("<ol class=\"footnotes-defs-ol\">")
	for i, def := range context.footnotesDefs {
		r.writeString("<li id=\"footnotes-def-" + strconv.Itoa(i+1) + "\">")
		tree := &Tree{Name: "", context: &Context{option: lute.options}}
		tree.context.tree = tree
		tree.Root = &ast.Node{Type: ast.NodeDocument}
		tree.Root.AppendChild(def)
		defRenderer := lute.newHTMLRenderer(tree)
		lc := tree.Root.LastDeepestChild()
		for i = len(def.FootnotesRefs) - 1; 0 <= i; i-- {
			ref := def.FootnotesRefs[i]
			gotoRef := " <a href=\"#footnotes-ref-" + ref.FootnotesRefId + "\" class=\"footnotes-goto-ref\">↩</a>"
			link := &ast.Node{Type: ast.NodeInlineHTML, Tokens: util.StrToBytes(gotoRef)}
			lc.InsertAfter(link)
		}
		defRenderer.(*HTMLRenderer).needRenderFootnotesDef = true
		defContent, err := defRenderer.Render()
		if nil != err {
			break
		}
		r.write(defContent)

		r.writeString("</li>\n")
	}
	r.writeString("</ol></div>")
	return r.writer.Bytes()
}

func (r *HTMLRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	idx, _ := r.tree.context.findFootnotesDef(node.Tokens)
	idxStr := strconv.Itoa(idx)
	r.tag("sup", [][]string{{"class", "footnotes-ref"}, {"id", "footnotes-ref-" + node.FootnotesRefId}}, false)
	r.tag("a", [][]string{{"href", "#footnotes-def-" + idxStr}}, false)
	r.writeString(idxStr)
	r.tag("/a", nil, false)
	r.tag("/sup", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if !r.needRenderFootnotesDef {
		return ast.WalkStop
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/span", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(escapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HTMLRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	attrs := [][]string{{"class", "vditor-math"}}
	r.tag("span", attrs, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/div", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(escapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HTMLRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	attrs := [][]string{{"class", "vditor-math"}}
	r.tag("div", attrs, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.newline()
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
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
		r.newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("tr", nil, false)
		r.newline()
	} else {
		r.tag("/tr", nil, false)
		r.newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("thead", nil, false)
		r.newline()
	} else {
		r.tag("/thead", nil, false)
		r.newline()
		if nil != node.Next {
			r.tag("tbody", nil, false)
		}
		r.newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("table", nil, false)
		r.newline()
	} else {
		if nil != node.FirstChild.Next {
			r.tag("/tbody", nil, false)
		}
		r.newline()
		r.tag("/table", nil, false)
		r.newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("del", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/del", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("del", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/del", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if r.option.AutoSpace {
		r.space(node)
	}
	r.write(escapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HTMLRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if 0 == r.disableTags {
			r.writeString("<img src=\"")
			destTokens := node.ChildByType(ast.NodeLinkDest).Tokens
			destTokens = r.tree.context.relativePath(destTokens)
			r.write(escapeHTML(destTokens))
			r.writeString("\" alt=\"")
		}
		r.disableTags++
		return ast.WalkContinue
	}

	r.disableTags--
	if 0 == r.disableTags {
		r.writeString("\"")
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			r.writeString(" title=\"")
			r.write(escapeHTML(title.Tokens))
			r.writeString("\"")
		}
		r.writeString(" />")
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		dest := node.ChildByType(ast.NodeLinkDest)
		destTokens := dest.Tokens
		destTokens = r.tree.context.relativePath(destTokens)
		attrs := [][]string{{"href", util.BytesToStr(escapeHTML(destTokens))}}
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			attrs = append(attrs, []string{"title", util.BytesToStr(escapeHTML(title.Tokens))})
		}
		r.tag("a", attrs, false)
	} else {
		r.tag("/a", nil, false)
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.newline()
	r.write(node.Tokens)
	r.newline()
	return ast.WalkStop
}

func (r *HTMLRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if grandparent := node.Parent.Parent; nil != grandparent && ast.NodeList == grandparent.Type && grandparent.Tight { // List.ListItem.Paragraph
		return ast.WalkContinue
	}

	if entering {
		r.newline()
		r.tag("p", nil, false)
	} else {
		r.tag("/p", nil, false)
		r.newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if r.option.AutoSpace {
		r.space(node)
	}
	if r.option.FixTermTypo {
		r.fixTermTypo(node)
	}
	if r.option.ChinesePunct {
		r.chinesePunct(node)
	}
	r.write(escapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HTMLRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if r.option.AutoSpace {
			if text := node.PreviousNodeText(); "" != text {
				lastc, _ := utf8.DecodeLastRuneInString(text)
				if unicode.IsLetter(lastc) || unicode.IsDigit(lastc) {
					r.writeByte(itemSpace)
				}
			}
		}
	} else {
		if r.option.AutoSpace {
			if text := node.NextNodeText(); "" != text {
				firstc, _ := utf8.DecodeRuneInString(text)
				if unicode.IsLetter(firstc) || unicode.IsDigit(firstc) {
					r.writeByte(itemSpace)
				}
			}
		}
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("<code>")
	return ast.WalkStop
}

func (r *HTMLRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(escapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HTMLRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("</code>")
	return ast.WalkStop
}

func (r *HTMLRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("em", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/em", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("em", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/em", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("strong", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/strong", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("strong", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/strong", nil, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.newline()
		r.writeString("<blockquote>")
		r.newline()
	} else {
		r.newline()
		r.writeString("</blockquote>")
		r.newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.newline()
		level := " 123456"[node.HeadingLevel : node.HeadingLevel+1]
		r.writeString("<h" + level)
		if r.option.ToC {
			r.writeString(" id=\"toc_h" + level + "_" + strconv.Itoa(r.headingCnt) + "\"")
			r.headingCnt++
		}
		r.writeString(">")
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
		r.newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	tag := "ul"
	if 1 == node.ListData.Typ || (3 == node.ListData.Typ && 0 == node.ListData.BulletChar) {
		tag = "ol"
	}
	if entering {
		r.newline()
		attrs := [][]string{{"start", strconv.Itoa(node.Start)}}
		if 0 == node.BulletChar && 1 != node.Start {
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
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if 3 == node.ListData.Typ && "" != r.option.GFMTaskListItemClass &&
			nil != node.FirstChild && nil != node.FirstChild.FirstChild && ast.NodeTaskListItemMarker == node.FirstChild.FirstChild.Type {
			r.tag("li", [][]string{{"class", r.option.GFMTaskListItemClass}}, false)
		} else {
			r.tag("li", nil, false)
		}
	} else {
		r.tag("/li", nil, false)
		r.newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		if node.TaskListItemChecked {
			attrs = append(attrs, []string{"checked", ""})
		}
		attrs = append(attrs, []string{"disabled", ""}, []string{"type", "checkbox"})
		r.tag("input", attrs, true)
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.newline()
	r.tag("hr", nil, true)
	r.newline()
	return ast.WalkStop
}

func (r *HTMLRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("br", nil, true)
		r.newline()
	}
	return ast.WalkStop
}

func (r *HTMLRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if r.option.SoftBreak2HardBreak {
		r.tag("br", nil, true)
		r.newline()
	} else {
		r.newline()
	}
	return ast.WalkStop
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
