// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
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

	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// VditorRenderer 描述了 Vditor WYSIWYG DOM 渲染器。
type VditorRenderer struct {
	*BaseRenderer
	commentStackDepth int
}

// NewVditorRenderer 创建一个 Vditor WYSIWYG DOM 渲染器。
func NewVditorRenderer(tree *parse.Tree, options *Options) *VditorRenderer {
	ret := &VditorRenderer{BaseRenderer: NewBaseRenderer(tree, options)}
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
	ret.RendererFuncs[ast.NodeOpenBrace] = ret.renderOpenBrace
	ret.RendererFuncs[ast.NodeCloseBrace] = ret.renderCloseBrace
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
	ret.RendererFuncs[ast.NodeFootnotesDefBlock] = ret.renderFootnotesDefBlock
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
	ret.RendererFuncs[ast.NodeMark] = ret.renderMark
	ret.RendererFuncs[ast.NodeMark1OpenMarker] = ret.renderMark1OpenMarker
	ret.RendererFuncs[ast.NodeMark1CloseMarker] = ret.renderMark1CloseMarker
	ret.RendererFuncs[ast.NodeMark2OpenMarker] = ret.renderMark2OpenMarker
	ret.RendererFuncs[ast.NodeMark2CloseMarker] = ret.renderMark2CloseMarker
	ret.RendererFuncs[ast.NodeSup] = ret.renderSup
	ret.RendererFuncs[ast.NodeSupOpenMarker] = ret.renderSupOpenMarker
	ret.RendererFuncs[ast.NodeSupCloseMarker] = ret.renderSupCloseMarker
	ret.RendererFuncs[ast.NodeSub] = ret.renderSub
	ret.RendererFuncs[ast.NodeSubOpenMarker] = ret.renderSubOpenMarker
	ret.RendererFuncs[ast.NodeSubCloseMarker] = ret.renderSubCloseMarker
	ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL
	ret.RendererFuncs[ast.NodeLinkRefDefBlock] = ret.renderLinkRefDefBlock
	ret.RendererFuncs[ast.NodeLinkRefDef] = ret.renderLinkRefDef
	return ret
}

func (r *VditorRenderer) renderLinkRefDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<div data-block=\"0\" data-type=\"link-ref-defs-block\">")
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderLinkRefDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		dest := node.FirstChild.ChildByType(ast.NodeLinkDest).Tokens
		destStr := util.BytesToStr(dest)
		r.WriteString("[" + util.BytesToStr(node.Tokens) + "]:")
		if editor.Caret != destStr {
			r.WriteString(" ")
		}
		r.WriteString(destStr + "\n")
	}
	return ast.WalkSkipChildren
}

func (r *VditorRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, editor.Caret, "")
		if "" == previousNodeText {
			r.WriteString(editor.Zwsp)
		}
	} else {
		r.WriteString(editor.Zwsp)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderMark1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("mark", [][]string{{"data-marker", "="}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("mark", [][]string{{"data-marker", "=="}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderSupOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("sup", [][]string{{"data-marker", "^"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderSupCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderSubOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("sub", [][]string{{"data-marker", "~"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderSubCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sub", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previewTokens := bytes.TrimSpace(node.Tokens)
		codeLen := len(previewTokens)
		codeIsEmpty := 1 > codeLen || (len(editor.Caret) == codeLen && editor.Caret == string(node.Tokens))
		r.Tag("pre", nil, false)
		r.Tag("code", [][]string{{"data-type", "yaml-front-matter"}}, false)
		if codeIsEmpty {
			r.WriteString(editor.FrontEndCaret + "\n")
		} else {
			r.Write(html.EscapeHTML(previewTokens))
		}
		r.WriteString("</code></pre>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString(`<div class="vditor-wysiwyg__block" data-type="yaml-front-matter" data-block="0">`)
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	previousNodeText := node.PreviousNodeText()
	previousNodeText = strings.ReplaceAll(previousNodeText, editor.Caret, "")
	if "" == previousNodeText {
		r.WriteString(editor.Zwsp)
	}

	r.WriteString("<span class=\"vditor-wysiwyg__block\" data-type=\"html-entity\">")
	r.Tag("code", [][]string{{"data-type", "html-entity"}, {"style", "display: none"}}, false)
	tokens := append([]byte(editor.Zwsp), node.HtmlEntityTokens...)
	r.Write(html.EscapeHTML(tokens))
	r.WriteString("</code>")

	r.Tag("span", [][]string{{"class", "vditor-wysiwyg__preview"}, {"data-render", "2"}}, false)
	r.Tag("code", nil, false)
	previewTokens := bytes.ReplaceAll(node.HtmlEntityTokens, editor.CaretTokens, nil)
	r.Write(previewTokens)
	r.Tag("/code", nil, false)
	r.Tag("/span", nil, false)
	r.WriteString("</span>" + editor.Zwsp)
	return ast.WalkContinue
}

func (r *VditorRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
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

func (r *VditorRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	return r.BaseRenderer.renderToC(node, entering)
}

func (r *VditorRenderer) renderFootnotesDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<div data-block=\"0\" data-type=\"footnotes-block\">")
		r.WriteString("<ol data-type=\"footnotes-defs-ol\">")
	} else {
		r.WriteString("</ol></div>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if r.RenderingFootnotes {
		return ast.WalkContinue
	}

	if entering {
		if nil != node.Previous && bytes.EqualFold(node.Previous.Tokens, node.Tokens) {
			return ast.WalkContinue
		}

		r.WriteString("<li data-type=\"footnotes-li\" data-marker=\"" + string(node.Tokens) + "\">")
		for c := node.FirstChild; nil != c; c = c.Next {
			ast.Walk(c, func(n *ast.Node, entering bool) ast.WalkStatus {
				return r.RendererFuncs[n.Type](n, entering)
			})
		}
		r.WriteString("</li>")
		return ast.WalkSkipChildren
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, editor.Caret, "")
		if "" == previousNodeText {
			r.WriteString(editor.Zwsp)
		}
		idx, def := r.Tree.FindFootnotesDef(node.Tokens)
		idxStr := strconv.Itoa(idx)
		label := def.Text()
		r.Tag("sup", [][]string{{"data-type", "footnotes-ref"}, {"data-footnotes-label", string(node.FootnotesRefLabel)},
			{"class", "vditor-tooltipped vditor-tooltipped__s"}, {"aria-label", SubStr(html.EscapeString(label), 24)}}, false)
		r.WriteString(idxStr)
		r.WriteString("</sup>" + editor.Zwsp)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	tokens := bytes.ReplaceAll(node.Tokens, []byte(editor.Zwsp), nil)
	previewTokens := tokens
	codeAttrs := [][]string{{"data-type", "math-inline"}}
	if !bytes.Contains(previewTokens, editor.CaretTokens) {
		codeAttrs = append(codeAttrs, []string{"style", "display: none"})
	}
	r.WriteString("<span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\">")
	r.Tag("code", codeAttrs, false)
	tokens = html.EscapeHTML(tokens)
	tokens = append([]byte(editor.Zwsp), tokens...)
	r.Write(tokens)
	r.WriteString("</code>")

	r.Tag("span", [][]string{{"class", "vditor-wysiwyg__preview"}, {"data-render", "2"}}, false)
	r.Tag("span", [][]string{{"class", "language-math"}}, false)
	previewTokens = bytes.ReplaceAll(previewTokens, editor.CaretTokens, nil)
	if node.ParentIs(ast.NodeTableCell) {
		// Improve the `|` render in the inline math in the table https://github.com/Vanessa219/vditor/issues/1550
		previewTokens = bytes.ReplaceAll(previewTokens, []byte("\\|"), []byte("|"))
	}
	r.Write(html.EscapeHTML(previewTokens))
	r.Tag("/span", nil, false)
	r.Tag("/span", nil, false)
	r.WriteString("</span>" + editor.Zwsp)
	return ast.WalkContinue
}

func (r *VditorRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, editor.Caret, "")
		if "" == previousNodeText {
			r.WriteString(editor.Zwsp)
		}
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	previewTokens := bytes.TrimSpace(node.Tokens)
	var preAttrs [][]string
	if !bytes.Contains(previewTokens, editor.CaretTokens) && r.Options.VditorMathBlockPreview {
		preAttrs = append(preAttrs, []string{"style", "display: none"})
	}
	codeLen := len(previewTokens)
	codeIsEmpty := 1 > codeLen || (len(editor.Caret) == codeLen && editor.Caret == string(node.Tokens))
	r.Tag("pre", preAttrs, false)
	r.Tag("code", [][]string{{"data-type", "math-block"}}, false)
	if codeIsEmpty {
		r.WriteString(editor.FrontEndCaret + "\n")
	} else {
		r.Write(html.EscapeHTML(previewTokens))
	}
	r.WriteString("</code></pre>")

	if r.Options.VditorMathBlockPreview {
		r.Tag("pre", [][]string{{"class", "vditor-wysiwyg__preview"}, {"data-render", "2"}}, false)
		r.Tag("div", [][]string{{"data-type", "math-block"}, {"class", "language-math"}}, false)
		tokens := node.Tokens
		tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
		r.Write(html.EscapeHTML(tokens))
		r.WriteString("</div></pre>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
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
		r.Tag(tag, attrs, false)
		if nil == node.FirstChild {
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
		} else if bytes.Equal(node.FirstChild.Tokens, editor.CaretTokens) {
			node.FirstChild.Tokens = []byte(editor.Caret + " ")
		} else {
			node.FirstChild.Tokens = bytes.TrimSpace(node.FirstChild.Tokens)
		}
	} else {
		r.Tag("/"+tag, nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("tr", nil, false)
	} else {
		r.Tag("/tr", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("thead", nil, false)
	} else {
		r.Tag("/thead", nil, false)
		if nil != node.Next {
			r.Tag("tbody", nil, false)
		}
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("table", [][]string{{"data-block", "0"}}, false)
	} else {
		if nil != node.FirstChild.Next {
			r.Tag("/tbody", nil, false)
		}
		r.Tag("/table", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("s", [][]string{{"data-marker", "~"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/s", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("s", [][]string{{"data-marker", "~~"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/s", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCloseBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderOpenBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if 3 == node.LinkType {
			previousNodeText := node.PreviousNodeText()
			previousNodeText = strings.ReplaceAll(previousNodeText, editor.Caret, "")
			if "" == previousNodeText {
				r.WriteString(editor.Zwsp)
			}
			r.WriteString("<img src=\"")
			link := r.Tree.FindLinkRefDefLink(node.LinkRefLabel)
			destTokens := link.ChildByType(ast.NodeLinkDest).Tokens
			destTokens = r.LinkPath(destTokens)
			destTokens = bytes.ReplaceAll(destTokens, editor.CaretTokens, nil)
			r.Write(destTokens)
			r.WriteString("\" alt=\"")
			if alt := node.ChildByType(ast.NodeLinkText); nil != alt {
				alt.Tokens = bytes.ReplaceAll(alt.Tokens, editor.CaretTokens, nil)
				r.Write(alt.Tokens)
			}
			r.WriteByte(lex.ItemDoublequote)
			if title := link.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
				r.WriteString(" title=\"")
				title.Tokens = bytes.ReplaceAll(title.Tokens, editor.CaretTokens, nil)
				r.Write(title.Tokens)
				r.WriteByte(lex.ItemDoublequote)
			}
			r.WriteString(" data-type=\"link-ref\" data-link-label=\"" + string(node.LinkRefLabel) + "\"")
			r.WriteString(" />")

			// XSS 过滤
			buf := r.Writer.Bytes()
			idx := bytes.LastIndex(buf, []byte("<img src="))
			imgBuf := buf[idx:]
			if r.Options.Sanitize {
				imgBuf = sanitize(imgBuf)
			}
			r.Writer.Truncate(idx)
			r.Writer.Write(imgBuf)
			return ast.WalkSkipChildren
		}

		if 0 == r.DisableTags {
			r.WriteString("<img src=\"")
			destTokens := node.ChildByType(ast.NodeLinkDest).Tokens
			destTokens = r.LinkPath(destTokens)
			destTokens = bytes.ReplaceAll(destTokens, editor.CaretTokens, nil)
			r.Write(destTokens)
			r.WriteString("\" alt=\"")
			if alt := node.ChildByType(ast.NodeLinkText); nil != alt && bytes.Contains(alt.Tokens, editor.CaretTokens) {
				alt.Tokens = bytes.ReplaceAll(alt.Tokens, editor.CaretTokens, nil)
			}
		}
		r.DisableTags++
		return ast.WalkContinue
	}

	r.DisableTags--
	if 0 == r.DisableTags {
		r.WriteByte(lex.ItemDoublequote)
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			r.WriteString(" title=\"")
			title.Tokens = bytes.ReplaceAll(title.Tokens, editor.CaretTokens, nil)
			r.Write(title.Tokens)
			r.WriteByte(lex.ItemDoublequote)
		}
		r.WriteString(" />")

		// XSS 过滤
		buf := r.Writer.Bytes()
		idx := bytes.LastIndex(buf, []byte("<img src="))
		imgBuf := buf[idx:]
		if r.Options.Sanitize {
			imgBuf = sanitize(imgBuf)
		}
		r.Writer.Truncate(idx)
		r.Writer.Write(imgBuf)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if 3 == node.LinkType {
		if entering {
			previousNodeText := node.PreviousNodeText()
			previousNodeText = strings.ReplaceAll(previousNodeText, editor.Caret, "")
			if "" == previousNodeText {
				r.WriteString(editor.Zwsp)
			}

			linkText := node.ChildrenByType(ast.NodeLinkText)
			var text string
			if 0 < len(linkText) {
				text = string(linkText[0].Tokens)
			}
			label := string(node.LinkRefLabel)
			attrs := [][]string{{"data-type", "link-ref"}, {"data-link-label", label}}
			r.Tag("span", attrs, false)
			r.WriteString(text)
			r.Tag("/span", nil, false)
			r.WriteString(editor.Zwsp)
			return ast.WalkSkipChildren
		} else {
			return ast.WalkContinue
		}
	}

	if entering {
		dest := node.ChildByType(ast.NodeLinkDest)
		destTokens := dest.Tokens
		if r.Options.Sanitize {
			tokens := bytes.TrimSpace(destTokens)
			tokens = bytes.ToLower(tokens)
			if bytes.HasPrefix(tokens, []byte("javascript:")) {
				destTokens = nil
			}
		}
		destTokens = r.LinkPath(destTokens)
		caretInDest := bytes.Contains(destTokens, editor.CaretTokens)
		if caretInDest {
			text := node.ChildByType(ast.NodeLinkText)
			text.Tokens = append(text.Tokens, editor.CaretTokens...)
			destTokens = bytes.ReplaceAll(destTokens, editor.CaretTokens, nil)
		}
		attrs := [][]string{{"href", string(destTokens)}}
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			title.Tokens = bytes.ReplaceAll(title.Tokens, editor.CaretTokens, nil)
			attrs = append(attrs, []string{"title", string(title.Tokens)})
		}
		r.Tag("a", attrs, false)
	} else {
		r.Tag("/a", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	r.WriteString(`<div class="vditor-wysiwyg__block" data-type="html-block" data-block="0">`)
	tokens := bytes.TrimSpace(node.Tokens)
	r.WriteString("<pre>")
	r.Tag("code", nil, false)
	r.Write(html.EscapeHTML(tokens))
	r.WriteString("</code></pre>")

	r.Tag("pre", [][]string{{"class", "vditor-wysiwyg__preview"}, {"data-render", "2"}}, false)
	tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
	if r.Options.Sanitize {
		tokens = sanitize(tokens)
	}
	r.Write(tokens)
	r.WriteString("</pre></div>")
	return ast.WalkContinue
}

func (r *VditorRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	if bytes.Equal(node.Tokens, []byte("<br />")) && node.ParentIs(ast.NodeTableCell) {
		r.Write(node.Tokens)
		return ast.WalkContinue
	}

	if bytes.Contains(node.Tokens, []byte("<span class=\"vditor-comment")) {
		r.commentStackDepth++
		r.Write(node.Tokens)
		return ast.WalkContinue
	}

	if bytes.Equal(node.Tokens, []byte("</span>")) {
		if 0 < r.commentStackDepth {
			r.commentStackDepth--
			r.Write(node.Tokens)
			return ast.WalkContinue
		}
	}

	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, editor.Caret, "")
		if editor.Zwsp == previousNodeText || "" == previousNodeText {
			r.WriteString(editor.Zwsp)
		}
	}

	tokens := bytes.ReplaceAll(node.Tokens, []byte(editor.Zwsp), nil)
	tokens = append([]byte(editor.Zwsp), tokens...)

	node.Tokens = bytes.TrimSpace(node.Tokens)
	r.Tag("code", [][]string{{"data-type", "html-inline"}}, false)
	tokens = html.EscapeHTML(tokens)
	r.Write(tokens)
	r.WriteString("</code>")
	return ast.WalkContinue
}

func (r *VditorRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if grandparent := node.Parent.Parent; nil != grandparent && ast.NodeList == grandparent.Type && grandparent.ListData.Tight { // List.ListItem.Paragraph
		return ast.WalkContinue
	}

	if entering {
		attr := [][]string{{"data-block", "0"}}
		attr = append(attr, node.KramdownIAL...)
		r.Tag("p", attr, false)
	} else {
		r.Tag("/p", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := node.Tokens
		if r.Options.FixTermTypo {
			tokens = r.FixTermTypo(tokens)
		}

		tokens = bytes.TrimRight(tokens, "\n")
		// 有的场景需要零宽空格撑起，但如果有其他文本内容的话需要把零宽空格删掉
		if !bytes.EqualFold(tokens, []byte(editor.Caret+editor.Zwsp)) {
			tokens = bytes.ReplaceAll(tokens, []byte(editor.Zwsp), nil)
		}
		r.Write(html.EscapeHTML(tokens))
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, editor.Caret, "")
		if "" == previousNodeText {
			r.WriteString(editor.Zwsp)
		} else {
			lastc, _ := utf8.DecodeLastRuneInString(previousNodeText)
			if unicode.IsLetter(lastc) || unicode.IsDigit(lastc) {
				r.WriteByte(lex.ItemSpace)
			}
		}
		r.Tag("code", [][]string{{"data-marker", strings.Repeat("`", node.CodeMarkerLen)}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := bytes.ReplaceAll(node.Tokens, []byte(editor.Zwsp), nil)
		tokens = html.EscapeHTML(tokens)
		tokens = append([]byte(editor.Zwsp), tokens...)
		r.Write(tokens)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</code>")
		codeSpan := node.Parent
		if codeSpanParent := codeSpan.Parent; nil != codeSpanParent && ast.NodeLink == codeSpanParent.Type {
			return ast.WalkContinue
		}
		r.WriteString(editor.Zwsp)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", [][]string{{"data-marker", "*"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", [][]string{{"data-marker", "_"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("strong", [][]string{{"data-marker", "**"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("strong", [][]string{{"data-marker", "__"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
	}
	return ast.WalkContinue
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
	return ast.WalkContinue
}

func (r *VditorRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1] + " data-block=\"0\"")
		var id string
		headingID := node.ChildByType(ast.NodeHeadingID)
		if nil != headingID {
			id = string(headingID.Tokens)
		}
		if r.Options.HeadingID && "" != id {
			r.WriteString(" data-id=\"" + id + "\"")
		}
		if "" == id {
			id = HeadingID(node)
		}
		r.WriteString(" id=\"wysiwyg-" + id + "\"")
		if !node.HeadingSetext {
			r.WriteString(" data-marker=\"#\">")
		} else {
			if 1 == node.HeadingLevel {
				r.WriteString(" data-marker=\"=\">")
			} else {
				r.WriteString(" data-marker=\"-\">")
			}
		}
		if r.Options.HeadingAnchor {
			id := HeadingID(node)
			r.Tag("a", [][]string{{"id", "vditorAnchor-" + id}, {"class", "vditor-anchor"}, {"href", "#" + id}}, false)
			r.WriteString(`<svg viewBox="0 0 16 16" version="1.1" width="16" height="16"><path fill-rule="evenodd" d="M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z"></path></svg>`)
			r.Tag("/a", nil, false)
		}
	} else {
		r.WriteString("</h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1] + ">")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	tag := "ul"
	if 1 == node.ListData.Typ || (3 == node.ListData.Typ && 0 == node.ListData.BulletChar) {
		tag = "ol"
	}
	if entering {
		var attrs [][]string
		if node.ListData.Tight {
			attrs = append(attrs, []string{"data-tight", "true"})
		}
		if 0 == node.ListData.BulletChar {
			if 1 != node.ListData.Start {
				attrs = append(attrs, []string{"start", strconv.Itoa(node.ListData.Start)})
			}
		}
		switch node.ListData.Typ {
		case 0:
			attrs = append(attrs, []string{"data-marker", string(node.ListData.Marker)})
		case 1:
			attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.ListData.Num) + string(node.ListData.Delimiter)})
		case 3:
			if 0 == node.ListData.BulletChar {
				attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.ListData.Num) + string(node.ListData.Delimiter)})
			} else {
				attrs = append(attrs, []string{"data-marker", string(node.ListData.Marker)})
			}
		}
		attrs = append(attrs, []string{"data-block", "0"})
		r.renderListStyle(node, &attrs)
		r.Tag(tag, attrs, false)
	} else {
		r.Tag("/"+tag, nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		switch node.ListData.Typ {
		case 0:
			attrs = append(attrs, []string{"data-marker", string(node.ListData.Marker)})
		case 1:
			attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.ListData.Num) + string(node.ListData.Delimiter)})
		case 3:
			if 0 == node.ListData.BulletChar {
				attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.ListData.Num) + string(node.ListData.Delimiter)})
			} else {
				attrs = append(attrs, []string{"data-marker", string(node.ListData.Marker)})
			}
			if nil != node.FirstChild && nil != node.FirstChild.FirstChild && ast.NodeTaskListItemMarker == node.FirstChild.FirstChild.Type { // li.p.task
				attrs = append(attrs, []string{"class", r.Options.GFMTaskListItemClass})
			}
		}
		r.Tag("li", attrs, false)
		if nil == node.FirstChild {
			r.WriteString(editor.Zwsp)
		}
	} else {
		r.Tag("/li", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		if node.TaskListItemChecked {
			attrs = append(attrs, []string{"checked", ""})
		}
		attrs = append(attrs, []string{"type", "checkbox"})
		r.Tag("input", attrs, true)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("hr", [][]string{{"data-block", "0"}}, true)
		if nil != node.Tokens {
			r.Tag("p", [][]string{{"data-block", "0"}}, false)
			r.Write(node.Tokens)
			r.WriteByte(lex.ItemNewline)
			r.Tag("/p", nil, false)
		}
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("br", nil, true)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		marker := "```"
		if nil != node.FirstChild && bytes.HasPrefix(node.FirstChild.Tokens, []byte(marker)) {
			marker = string(node.FirstChild.Tokens)
		}
		r.WriteString(`<div class="vditor-wysiwyg__block" data-type="code-block" data-block="0" data-marker="` + marker + `">`)
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	codeLen := len(node.Tokens)
	codeIsEmpty := 1 > codeLen || (len(editor.Caret) == codeLen && editor.Caret == string(node.Tokens))
	isFenced := node.Parent.IsFencedCodeBlock
	var language string
	var caretInInfo bool
	var attrs [][]string
	if isFenced && 0 < len(node.Previous.CodeBlockInfo) {
		if bytes.Contains(node.Previous.CodeBlockInfo, editor.CaretTokens) {
			caretInInfo = true
			node.Previous.CodeBlockInfo = bytes.ReplaceAll(node.Previous.CodeBlockInfo, editor.CaretTokens, nil)
		}
		if 0 < len(node.Previous.CodeBlockInfo) {
			infoWords := lex.Split(node.Previous.CodeBlockInfo, lex.ItemSpace)
			language = string(infoWords[0])
			attrs = append(attrs, []string{"class", "language-" + language})
			if "mindmap" == language {
				dataCode := EChartsMindmap(node.Tokens)
				attrs = append(attrs, []string{"data-code", string(dataCode)})
			}
		}
	}
	preAttrs := [][]string{{"class", "vditor-wysiwyg__pre"}}
	if !bytes.Contains(node.Tokens, editor.CaretTokens) && !caretInInfo && r.Options.VditorCodeBlockPreview {
		preAttrs = append(preAttrs, []string{"style", "display: none"})
	}
	r.Tag("pre", preAttrs, false)
	r.Tag("code", attrs, false)

	if codeIsEmpty {
		r.WriteString(editor.FrontEndCaret + "\n")
	} else {
		if caretInInfo {
			r.WriteString(editor.FrontEndCaret)
		}
		r.Write(html.EscapeHTML(node.Tokens))
		r.Newline()
	}
	r.WriteString("</code></pre>")

	if r.Options.VditorCodeBlockPreview {
		r.Tag("pre", [][]string{{"class", "vditor-wysiwyg__preview"}, {"data-render", "2"}}, false)
		preDiv := NoHighlight(language)
		if preDiv {
			r.Tag("div", attrs, false)
		} else {
			r.Tag("code", attrs, false)
		}
		tokens := node.Tokens
		tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
		r.Write(html.EscapeHTML(tokens))
		if preDiv {
			r.WriteString("</div></pre>")
		} else {
			r.WriteString("</code></pre>")
		}
	}
	return ast.WalkContinue
}
