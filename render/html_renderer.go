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
	"github.com/88250/lute/html"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// HtmlRenderer 描述了 HTML 渲染器。
type HtmlRenderer struct {
	*BaseRenderer
	needRenderFootnotesDef bool
}

// NewHtmlRenderer 创建一个 HTML 渲染器。
func NewHtmlRenderer(tree *parse.Tree) *HtmlRenderer {
	ret := &HtmlRenderer{NewBaseRenderer(tree), false}
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
	ret.RendererFuncs[ast.NodeHeadingID] = ret.renderHeadingID
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
	ret.RendererFuncs[ast.NodeHTMLEntity] = ret.renderHtmlEntity
	ret.RendererFuncs[ast.NodeYamlFrontMatter] = ret.renderYamlFrontMatter
	ret.RendererFuncs[ast.NodeYamlFrontMatterOpenMarker] = ret.renderYamlFrontMatterOpenMarker
	ret.RendererFuncs[ast.NodeYamlFrontMatterContent] = ret.renderYamlFrontMatterContent
	ret.RendererFuncs[ast.NodeYamlFrontMatterCloseMarker] = ret.renderYamlFrontMatterCloseMarker
	ret.RendererFuncs[ast.NodeBlockRef] = ret.renderBlockRef
	ret.RendererFuncs[ast.NodeBlockRefID] = ret.renderBlockRefID
	ret.RendererFuncs[ast.NodeBlockRefSpace] = ret.renderBlockRefSpace
	ret.RendererFuncs[ast.NodeBlockRefText] = ret.renderBlockRefText
	ret.RendererFuncs[ast.NodeMark] = ret.renderMark
	ret.RendererFuncs[ast.NodeMark1OpenMarker] = ret.renderMark1OpenMarker
	ret.RendererFuncs[ast.NodeMark1CloseMarker] = ret.renderMark1CloseMarker
	ret.RendererFuncs[ast.NodeMark2OpenMarker] = ret.renderMark2OpenMarker
	ret.RendererFuncs[ast.NodeMark2CloseMarker] = ret.renderMark2CloseMarker
	ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL
	ret.RendererFuncs[ast.NodeBlockQueryEmbed] = ret.renderBlockQueryEmbed
	ret.RendererFuncs[ast.NodeBlockQueryEmbedScript] = ret.renderBlockQueryEmbedScript
	ret.RendererFuncs[ast.NodeBlockEmbed] = ret.renderBlockEmbed
	ret.RendererFuncs[ast.NodeBlockEmbedID] = ret.renderBlockEmbedID
	ret.RendererFuncs[ast.NodeBlockEmbedSpace] = ret.renderBlockEmbedSpace
	ret.RendererFuncs[ast.NodeBlockEmbedText] = ret.renderBlockEmbedText
	ret.RendererFuncs[ast.NodeTag] = ret.renderTag
	ret.RendererFuncs[ast.NodeTagOpenMarker] = ret.renderTagOpenMarker
	ret.RendererFuncs[ast.NodeTagCloseMarker] = ret.renderTagCloseMarker
	return ret
}

func (r *HtmlRenderer) renderTag(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderTagOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("em", nil, false)
	r.WriteByte(lex.ItemCrosshatch)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderTagCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemCrosshatch)
	r.tag("/em", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderMark1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("mark", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/mark", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("mark", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/mark", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.tag("div", nil, false)
	} else {
		r.tag("/div", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockQueryEmbedScript(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(node.Tokens)
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.tag("div", nil, false)
	} else {
		r.tag("/div", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockEmbedID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockEmbedSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockEmbedText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(html.EscapeHTML(node.Tokens))
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkStop
}

func (r *HtmlRenderer) renderBlockRef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockRefID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockRefText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(html.EscapeHTML(node.Tokens))
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkStop
}

func (r *HtmlRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/div", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(html.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HtmlRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	attrs := [][]string{{"class", "vditor-yml-front-matter"}}
	r.tag("div", attrs, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(html.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HtmlRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(html.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HtmlRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	headings := r.headings()
	length := len(headings)
	if 1 > length {
		return ast.WalkStop
	}
	r.WriteString("<div class=\"vditor-toc\">")
	for _, heading := range headings {
		level := strconv.Itoa(heading.HeadingLevel)
		spaces := (heading.HeadingLevel - 1) * 2
		r.WriteString(strings.Repeat("&emsp;", spaces))
		r.WriteString("<span class=\"toc-h" + level + "\">")
		r.WriteString("<a class=\"toc-a\" href=\"#" + HeadingID(heading) + "\">" + heading.Text() + "</a></span><br>")
	}
	r.WriteString("</div>")
	return ast.WalkStop
}

func (r *HtmlRenderer) RenderFootnotesDefs(context *parse.Context) []byte {
	r.WriteString("<div class=\"footnotes-defs-div\">")
	r.WriteString("<hr class=\"footnotes-defs-hr\" />\n")
	r.WriteString("<ol class=\"footnotes-defs-ol\">")
	for i, def := range context.FootnotesDefs {
		r.WriteString("<li id=\"footnotes-def-" + strconv.Itoa(i+1) + "\">")
		tree := &parse.Tree{Name: "", Context: context}
		tree.Context.Tree = tree
		tree.Root = &ast.Node{Type: ast.NodeDocument}
		tree.Root.AppendChild(def)
		defRenderer := NewHtmlRenderer(tree)
		lc := tree.Root.LastDeepestChild()
		for i = len(def.FootnotesRefs) - 1; 0 <= i; i-- {
			ref := def.FootnotesRefs[i]
			gotoRef := " <a href=\"#footnotes-ref-" + ref.FootnotesRefId + "\" class=\"vditor-footnotes__goto-ref\">↩</a>"
			link := &ast.Node{Type: ast.NodeInlineHTML, Tokens: util.StrToBytes(gotoRef)}
			lc.InsertAfter(link)
		}
		defRenderer.needRenderFootnotesDef = true
		defContent := defRenderer.Render()
		r.Write(defContent)

		r.WriteString("</li>\n")
	}
	r.WriteString("</ol></div>")
	return r.Writer.Bytes()
}

func (r *HtmlRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	idx, _ := r.Tree.Context.FindFootnotesDef(node.Tokens)
	idxStr := strconv.Itoa(idx)
	r.tag("sup", [][]string{{"class", "footnotes-ref"}, {"id", "footnotes-ref-" + node.FootnotesRefId}}, false)
	r.tag("a", [][]string{{"href", "#footnotes-def-" + idxStr}}, false)
	r.WriteString(idxStr)
	r.tag("/a", nil, false)
	r.tag("/sup", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if !r.needRenderFootnotesDef {
		return ast.WalkStop
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/span", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(html.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HtmlRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	attrs := [][]string{{"class", "vditor-math"}}
	r.tag("span", attrs, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/div", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(html.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HtmlRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	attrs := [][]string{{"class", "vditor-math"}}
	r.tag("div", attrs, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *HtmlRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("tr", nil, false)
		r.Newline()
	} else {
		r.tag("/tr", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *HtmlRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *HtmlRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("del", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/del", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("del", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/del", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if r.Option.AutoSpace {
		r.Space(node)
	}
	r.Write(html.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HtmlRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if 0 == r.DisableTags {
			r.WriteString("<img src=\"")
			destTokens := node.ChildByType(ast.NodeLinkDest).Tokens
			destTokens = r.Tree.Context.LinkPath(destTokens)
			if "" != r.Option.ImageLazyLoading {
				r.Write(html.EscapeHTML(util.StrToBytes(r.Option.ImageLazyLoading)))
				r.WriteString("\" data-src=\"")
			}
			r.Write(html.EscapeHTML(destTokens))
			r.WriteString("\" alt=\"")
		}
		r.DisableTags++
		return ast.WalkContinue
	}

	r.DisableTags--
	if 0 == r.DisableTags {
		r.WriteByte(lex.ItemDoublequote)
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			r.WriteString(" title=\"")
			r.Write(html.EscapeHTML(title.Tokens))
			r.WriteByte(lex.ItemDoublequote)
		}
		r.WriteString(" />")

		if r.Option.Sanitize {
			buf := r.Writer.Bytes()
			idx := bytes.LastIndex(buf, []byte("<img src="))
			imgBuf := buf[idx:]
			if r.Option.Sanitize {
				imgBuf = sanitize(imgBuf)
			}
			r.Writer.Truncate(idx)
			r.Writer.Write(imgBuf)
		}
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.LinkTextAutoSpacePrevious(node)

		dest := node.ChildByType(ast.NodeLinkDest)
		destTokens := dest.Tokens
		destTokens = r.Tree.Context.LinkPath(destTokens)
		attrs := [][]string{{"href", util.BytesToStr(html.EscapeHTML(destTokens))}}
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			attrs = append(attrs, []string{"title", util.BytesToStr(html.EscapeHTML(title.Tokens))})
		}
		r.tag("a", attrs, false)
	} else {
		r.tag("/a", nil, false)

		r.LinkTextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	tokens := node.Tokens
	if r.Option.Sanitize {
		tokens = sanitize(tokens)
	}
	r.Write(tokens)
	r.Newline()
	return ast.WalkStop
}

func (r *HtmlRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	tokens := node.Tokens
	if r.Option.Sanitize {
		tokens = sanitize(tokens)
	}
	r.Write(tokens)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if grandparent := node.Parent.Parent; nil != grandparent && ast.NodeList == grandparent.Type && grandparent.Tight { // List.ListItem.Paragraph
		return ast.WalkContinue
	}

	if entering {
		r.Newline()
		r.tag("p", node.KramdownIAL, false)
		if r.Option.ChineseParagraphBeginningSpace && ast.NodeDocument == node.Parent.Type {
			r.WriteString("&emsp;&emsp;")
		}
	} else {
		r.tag("/p", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if r.Option.AutoSpace {
		r.Space(node)
	}
	if r.Option.FixTermTypo {
		r.FixTermTypo(node)
	}
	if r.Option.ChinesePunct {
		r.ChinesePunct(node)
	}
	r.Write(html.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HtmlRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *HtmlRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("<code>")
	return ast.WalkStop
}

func (r *HtmlRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(html.EscapeHTML(node.Tokens))
	return ast.WalkStop
}

func (r *HtmlRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("</code>")
	return ast.WalkStop
}

func (r *HtmlRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("em", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/em", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("em", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/em", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("strong", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/strong", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("strong", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.tag("/strong", nil, false)
	return ast.WalkStop
}

func (r *HtmlRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.tag("blockquote", node.KramdownIAL, false)
		r.Newline()
	} else {
		r.Newline()
		r.WriteString("</blockquote>")
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

var headingLevel = " 123456"

func (r *HtmlRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		level := headingLevel[node.HeadingLevel : node.HeadingLevel+1]
		r.WriteString("<h" + level)
		id := HeadingID(node)
		if r.Option.ToC || r.Option.HeadingID {
			r.WriteString(" id=\"" + id + "\"")
		}
		r.WriteString(">")
	} else {
		if r.Option.HeadingAnchor {
			id := HeadingID(node)
			r.tag("a", [][]string{{"id", "vditorAnchor-" + id}, {"class", "vditor-anchor"}, {"href", "#" + id}}, false)
			r.WriteString(`<svg viewBox="0 0 16 16" version="1.1" width="16" height="16"><path fill-rule="evenodd" d="M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z"></path></svg>`)
			r.tag("/a", nil, false)
		}
		r.WriteString("</h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1] + ">")
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *HtmlRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	tag := "ul"
	if 1 == node.ListData.Typ || (3 == node.ListData.Typ && 0 == node.ListData.BulletChar) {
		tag = "ol"
	}
	if entering {
		r.Newline()
		var attrs [][]string
		r.renderListStyle(node, &attrs)
		if 0 == node.BulletChar && 1 != node.Start {
			attrs = append(attrs, []string{"start", strconv.Itoa(node.Start)})
		}
		attrs = append(attrs, node.KramdownIAL...)
		r.tag(tag, attrs, false)
		r.Newline()
	} else {
		r.Newline()
		r.tag("/"+tag, nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		attrs = append(attrs, node.KramdownIAL...)
		if 3 == node.ListData.Typ && "" != r.Option.GFMTaskListItemClass &&
			nil != node.FirstChild && nil != node.FirstChild.FirstChild && ast.NodeTaskListItemMarker == node.FirstChild.FirstChild.Type {
			attrs = append(attrs, []string{"class", r.Option.GFMTaskListItemClass})
		}
		r.tag("li", attrs, false)
	} else {
		r.tag("/li", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *HtmlRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	r.tag("hr", nil, true)
	r.Newline()
	return ast.WalkStop
}

func (r *HtmlRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.tag("br", nil, true)
		r.Newline()
	}
	return ast.WalkStop
}

func (r *HtmlRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if r.Option.SoftBreak2HardBreak {
		r.tag("br", nil, true)
		r.Newline()
	} else {
		r.Newline()
	}
	return ast.WalkStop
}

func (r *HtmlRenderer) tag(name string, attrs [][]string, selfclosing bool) {
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
