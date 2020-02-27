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
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// HTMLRenderer 描述了 HTML 渲染器。
type HTMLRenderer struct {
	*BaseRenderer
	needRenderFootnotesDef bool
	headingCnt             int
}

// NewHTMLRenderer 创建一个 HTML 渲染器。
func NewHTMLRenderer(tree *parse.Tree) Renderer {
	ret := &HTMLRenderer{NewBaseRenderer(tree), false, 0}
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
	ret.RendererFuncs[ast.NodeToC] = ret.renderToC
	ret.RendererFuncs[ast.NodeBackslash] = ret.renderBackslash
	ret.RendererFuncs[ast.NodeBackslashContent] = ret.renderBackslashContent
	return ret
}

func (r *HTMLRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(util.EscapeHTML(node.Tokens))
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
	r.WriteString("<div class=\"toc-div\">")
	for i, heading := range headings {
		level := strconv.Itoa(heading.HeadingLevel)
		spaces := (heading.HeadingLevel - 1) * 2
		r.WriteString(strings.Repeat("&emsp;", spaces))
		r.WriteString("<span class=\"toc-h" + level + "\">")
		r.WriteString("<a class=\"toc-a\" href=\"#toc_h" + level + "_" + strconv.Itoa(i) + "\">" + heading.Text() + "</a></span><br>")
	}
	r.WriteString("</div>\n\n")

	return ast.WalkStop
}

func (r *HTMLRenderer) headings() (ret []*ast.Node) {
	for n := r.Tree.Root.FirstChild; nil != n; n = n.Next {
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

func (r *HTMLRenderer) RenderFootnotesDefs(context *parse.Context) []byte {
	r.WriteString("<div class=\"footnotes-defs-div\">")
	r.WriteString("<hr class=\"footnotes-defs-hr\" />\n")
	r.WriteString("<ol class=\"footnotes-defs-ol\">")
	for i, def := range context.FootnotesDefs {
		r.WriteString("<li id=\"footnotes-def-" + strconv.Itoa(i+1) + "\">")
		tree := &parse.Tree{Name: "", Context: context}
		tree.Context.Tree = tree
		tree.Root = &ast.Node{Type: ast.NodeDocument}
		tree.Root.AppendChild(def)
		defRenderer := NewHTMLRenderer(tree)
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
		r.Write(defContent)

		r.WriteString("</li>\n")
	}
	r.WriteString("</ol></div>")
	return r.Writer.Bytes()
}

func (r *HTMLRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	idx, _ := r.Tree.Context.FindFootnotesDef(node.Tokens)
	idxStr := strconv.Itoa(idx)
	r.tag("sup", [][]string{{"class", "footnotes-ref"}, {"id", "footnotes-ref-" + node.FootnotesRefId}}, false)
	r.tag("a", [][]string{{"href", "#footnotes-def-" + idxStr}}, false)
	r.WriteString(idxStr)
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
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
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
	r.Write(util.EscapeHTML(node.Tokens))
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
	r.Write(util.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HTMLRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	attrs := [][]string{{"class", "vditor-math"}}
	r.tag("div", attrs, false)
	return ast.WalkStop
}

func (r *HTMLRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
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
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("tr", nil, false)
		r.Newline()
	} else {
		r.tag("/tr", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("thead", nil, false)
		r.Newline()
	} else {
		r.tag("/thead", nil, false)
		r.Newline()
		if nil != node.Next {
			r.tag("tbody", nil, false)
		}
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("table", nil, false)
		r.Newline()
	} else {
		if nil != node.FirstChild.Next {
			r.tag("/tbody", nil, false)
		}
		r.Newline()
		r.tag("/table", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
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
	if r.Option.AutoSpace {
		r.Space(node)
	}
	r.Write(util.EscapeHTML(node.Tokens))
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
		if 0 == r.DisableTags {
			r.WriteString("<img src=\"")
			destTokens := node.ChildByType(ast.NodeLinkDest).Tokens
			destTokens = r.Tree.Context.RelativePath(destTokens)
			r.Write(util.EscapeHTML(destTokens))
			r.WriteString("\" alt=\"")
		}
		r.DisableTags++
		return ast.WalkContinue
	}

	r.DisableTags--
	if 0 == r.DisableTags {
		r.WriteString("\"")
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			r.WriteString(" title=\"")
			r.Write(util.EscapeHTML(title.Tokens))
			r.WriteString("\"")
		}
		r.WriteString(" />")
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.LinkTextAutoSpacePrevious(node)

		dest := node.ChildByType(ast.NodeLinkDest)
		destTokens := dest.Tokens
		destTokens = r.Tree.Context.RelativePath(destTokens)
		attrs := [][]string{{"href", util.BytesToStr(util.EscapeHTML(destTokens))}}
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			attrs = append(attrs, []string{"title", util.BytesToStr(util.EscapeHTML(title.Tokens))})
		}
		r.tag("a", attrs, false)
	} else {
		r.tag("/a", nil, false)

		r.LinkTextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	r.Write(node.Tokens)
	r.Newline()
	return ast.WalkStop
}

func (r *HTMLRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
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
		r.Newline()
		r.tag("p", nil, false)
	} else {
		r.tag("/p", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if r.Option.AutoSpace {
		r.Space(node)
	}
	if r.Option.FixTermTypo {
		r.FixTermTypo(node)
	}
	if r.Option.ChinesePunct {
		r.ChinesePunct(node)
	}
	r.Write(util.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HTMLRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if r.Option.AutoSpace {
			if text := node.PreviousNodeText(); "" != text {
				lastc, _ := utf8.DecodeLastRuneInString(text)
				if unicode.IsLetter(lastc) || unicode.IsDigit(lastc) {
					r.WriteByte(lex.ItemSpace)
				}
			}
		}
	} else {
		if r.Option.AutoSpace {
			if text := node.NextNodeText(); "" != text {
				firstc, _ := utf8.DecodeRuneInString(text)
				if unicode.IsLetter(firstc) || unicode.IsDigit(firstc) {
					r.WriteByte(lex.ItemSpace)
				}
			}
		}
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("<code>")
	return ast.WalkStop
}

func (r *HTMLRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(util.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HTMLRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("</code>")
	return ast.WalkStop
}

func (r *HTMLRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
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
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
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
		r.Newline()
		r.WriteString("<blockquote>")
		r.Newline()
	} else {
		r.Newline()
		r.WriteString("</blockquote>")
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HTMLRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		level := " 123456"[node.HeadingLevel : node.HeadingLevel+1]
		r.WriteString("<h" + level)
		if r.Option.ToC {
			r.WriteString(" id=\"toc_h" + level + "_" + strconv.Itoa(r.headingCnt) + "\"")
			r.headingCnt++
		}
		r.WriteString(">")
		if r.Option.HeadingAnchor {
			anchor := node.Text()
			anchor = strings.ReplaceAll(anchor, " ", "-")
			anchor = strings.ReplaceAll(anchor, ".", "")
			r.tag("a", [][]string{{"id", "vditorAnchor-" + anchor}, {"class", "vditor-anchor"}, {"href", "#" + anchor}}, false)
			r.WriteString(`<svg viewBox="0 0 16 16" version="1.1" width="16" height="16"><path fill-rule="evenodd" d="M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z"></path></svg>`)
			r.tag("/a", nil, false)
		}
	} else {
		r.WriteString("</h" + " 123456"[node.HeadingLevel:node.HeadingLevel+1] + ">")
		r.Newline()
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
		r.Newline()
		attrs := [][]string{{"start", strconv.Itoa(node.Start)}}
		if 0 == node.BulletChar && 1 != node.Start {
			r.tag(tag, attrs, false)
		} else {
			r.tag(tag, nil, false)
		}
		r.Newline()
	} else {
		r.Newline()
		r.tag("/"+tag, nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HTMLRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if 3 == node.ListData.Typ && "" != r.Option.GFMTaskListItemClass &&
			nil != node.FirstChild && nil != node.FirstChild.FirstChild && ast.NodeTaskListItemMarker == node.FirstChild.FirstChild.Type {
			r.tag("li", [][]string{{"class", r.Option.GFMTaskListItemClass}}, false)
		} else {
			r.tag("li", nil, false)
		}
	} else {
		r.tag("/li", nil, false)
		r.Newline()
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
	r.Newline()
	r.tag("hr", nil, true)
	r.Newline()
	return ast.WalkStop
}

func (r *HTMLRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("br", nil, true)
		r.Newline()
	}
	return ast.WalkStop
}

func (r *HTMLRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if r.Option.SoftBreak2HardBreak {
		r.tag("br", nil, true)
		r.Newline()
	} else {
		r.Newline()
	}
	return ast.WalkStop
}

func (r *HTMLRenderer) tag(name string, attrs [][]string, selfclosing bool) {
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
