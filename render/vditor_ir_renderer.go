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

	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// VditorIRRenderer 描述了 Vditor Instant-Rendering DOM 渲染器。
type VditorIRRenderer struct {
	*BaseRenderer
}

// NewVditorIRRenderer 创建一个 Vditor Instant-Rendering DOM 渲染器。
func NewVditorIRRenderer(tree *parse.Tree, options *Options) *VditorIRRenderer {
	ret := &VditorIRRenderer{BaseRenderer: NewBaseRenderer(tree, options)}
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
	ret.RendererFuncs[ast.NodeSup] = ret.renderSup
	ret.RendererFuncs[ast.NodeSupOpenMarker] = ret.renderSupOpenMarker
	ret.RendererFuncs[ast.NodeSupCloseMarker] = ret.renderSupCloseMarker
	ret.RendererFuncs[ast.NodeSub] = ret.renderSub
	ret.RendererFuncs[ast.NodeSubOpenMarker] = ret.renderSubOpenMarker
	ret.RendererFuncs[ast.NodeSubCloseMarker] = ret.renderSubCloseMarker
	ret.RendererFuncs[ast.NodeMark2OpenMarker] = ret.renderMark2OpenMarker
	ret.RendererFuncs[ast.NodeMark2CloseMarker] = ret.renderMark2CloseMarker
	ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL
	ret.RendererFuncs[ast.NodeLinkRefDefBlock] = ret.renderLinkRefDefBlock
	ret.RendererFuncs[ast.NodeLinkRefDef] = ret.renderLinkRefDef
	return ret
}

func (r *VditorIRRenderer) renderLinkRefDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<div data-block=\"0\" data-type=\"link-ref-defs-block\">")
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderLinkRefDef(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *VditorIRRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderMark1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString("=")
		r.Tag("/span", nil, false)
		r.Tag("mark", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString("=")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString("==")
		r.Tag("/span", nil, false)
		r.Tag("mark", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString("==")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderSupOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString("^")
		r.Tag("/span", nil, false)
		r.Tag("sup", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderSupCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sup", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString("^")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderSubOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString("~")
		r.Tag("/span", nil, false)
		r.Tag("sub", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderSubCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sub", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString("~")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "yaml-front-matter-close-marker"}}, false)
		r.Write(parse.YamlFrontMatterMarker)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		node.Tokens = bytes.TrimSpace(node.Tokens)
		codeLen := len(node.Tokens)
		codeIsEmpty := 1 > codeLen || (len(editor.Caret) == codeLen && editor.Caret == string(node.Tokens))
		r.Tag("pre", [][]string{{"class", "vditor-ir__marker--pre"}}, false)
		r.Tag("code", [][]string{{"data-type", "yaml-front-matter"}, {"class", "language-yaml"}}, false)
		if codeIsEmpty {
			r.WriteString(editor.FrontEndCaret + "\n")
		} else {
			r.Write(html.EscapeHTML(node.Tokens))
		}
		r.WriteString("</code></pre>")
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "yaml-front-matter-open-marker"}}, false)
		r.Write(parse.YamlFrontMatterMarker)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderDivNode(node)
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
		r.Tag("code", [][]string{{"data-newline", "1"}, {"class", "vditor-ir__marker vditor-ir__marker--pre"}, {"data-type", "html-entity"}}, false)
		r.Write(html.EscapeHTML(node.HtmlEntityTokens))
		r.Tag("/code", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__preview"}, {"data-render", "2"}}, false)
		r.Tag("code", nil, false)
		r.Write(node.HtmlEntityTokens)
		r.Tag("/code", nil, false)
		r.Tag("/span", nil, false)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bi"}}, false)
		r.WriteByte(lex.ItemBackslash)
		r.WriteString("</span>")
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	return r.BaseRenderer.renderToC(node, entering)
}

func (r *VditorIRRenderer) renderFootnotesDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<div data-block=\"0\" data-type=\"footnotes-block\">")
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if r.RenderingFootnotes {
			return ast.WalkContinue
		}

		r.WriteString("<div data-type=\"footnotes-def\">")
		r.WriteString("[" + string(node.Tokens) + "]: ")
		for c := node.FirstChild; nil != c; c = c.Next {
			ast.Walk(c, func(n *ast.Node, entering bool) ast.WalkStatus {
				return r.RendererFuncs[n.Type](n, entering)
			})
		}
		r.WriteString("</div>")
		return ast.WalkSkipChildren
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	previousNodeText := node.PreviousNodeText()
	previousNodeText = strings.ReplaceAll(previousNodeText, editor.Caret, "")
	if "" == previousNodeText {
		r.WriteString(editor.Zwsp)
	}
	idx, def := r.Tree.FindFootnotesDef(node.Tokens)
	idxStr := strconv.Itoa(idx)
	label := def.Text()
	attrs := [][]string{{"data-type", "footnotes-ref"}}
	text := node.Text()
	expand := strings.Contains(text, editor.Caret)
	if expand {
		attrs = append(attrs, []string{"class", "vditor-ir__node vditor-ir__node--expand vditor-tooltipped vditor-tooltipped__s"})
	} else {
		attrs = append(attrs, []string{"class", "vditor-ir__node vditor-tooltipped vditor-tooltipped__s"})
	}
	attrs = append(attrs, []string{"aria-label", SubStr(html.EscapeString(label), 24)})
	attrs = append(attrs, []string{"data-footnotes-label", string(node.FootnotesRefLabel)})
	r.Tag("sup", attrs, false)
	r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bracket"}}, false)
	r.WriteByte(lex.ItemOpenBracket)
	r.Tag("/span", nil, false)
	r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--link"}}, false)
	r.Write(node.Tokens)
	r.Tag("/span", nil, false)
	r.Tag("span", [][]string{{"class", "vditor-ir__marker--hide"}, {"data-render", "1"}}, false)
	r.WriteString(idxStr)
	r.Tag("/span", nil, false)
	r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bracket"}}, false)
	r.WriteByte(lex.ItemCloseBracket)
	r.Tag("/span", nil, false)
	r.WriteString("</sup>" + editor.Zwsp)
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "code-block-close-marker"}}, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--info"}, {"data-type", "code-block-info"}}, false)
		r.WriteString(editor.Zwsp)
		r.Write(node.CodeBlockInfo)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "code-block-open-marker"}}, false)
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderDivNode(node)
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	codeLen := len(node.Tokens)
	codeIsEmpty := 1 > codeLen || (len(editor.Caret) == codeLen && editor.Caret == string(node.Tokens))
	isFenced := node.Parent.IsFencedCodeBlock
	caretInInfo := false
	var language string
	if isFenced {
		caretInInfo = bytes.Contains(node.Previous.CodeBlockInfo, editor.CaretTokens)
		node.Previous.CodeBlockInfo = bytes.ReplaceAll(node.Previous.CodeBlockInfo, editor.CaretTokens, nil)
	}
	var attrs [][]string
	if isFenced && 0 < len(node.Previous.CodeBlockInfo) {
		infoWords := lex.Split(node.Previous.CodeBlockInfo, lex.ItemSpace)
		language = string(infoWords[0])
		attrs = append(attrs, []string{"class", "language-" + language})
		if "mindmap" == language {
			dataCode := EChartsMindmap(node.Tokens)
			attrs = append(attrs, []string{"data-code", string(dataCode)})
		}
	}

	class := "vditor-ir__marker--pre"
	if r.Options.VditorCodeBlockPreview {
		class += " vditor-ir__marker"
	}
	r.Tag("pre", [][]string{{"class", class}}, false)
	r.Tag("code", attrs, false)
	if codeIsEmpty {
		if !caretInInfo {
			r.WriteString(editor.FrontEndCaret)
		}
		r.WriteByte(lex.ItemNewline)
	} else {
		r.Write(html.EscapeHTML(node.Tokens))
		r.Newline()
	}
	r.WriteString("</code></pre>")

	if r.Options.VditorCodeBlockPreview {
		r.Tag("pre", [][]string{{"class", "vditor-ir__preview"}, {"data-render", "2"}}, false)
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

func (r *VditorIRRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<span data-render=\"2\">")
		r.Write(node.Tokens)
		r.WriteString("</span>")
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.Write(node.FirstChild.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<span data-render=\"2\">")
		r.Write(node.Tokens)
		r.WriteString("</span>")
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.Write(node.FirstChild.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previousNodeText := node.PreviousNodeText()
		previousNodeText = strings.ReplaceAll(previousNodeText, editor.Caret, "")
		if "" == previousNodeText {
			r.WriteString(editor.Zwsp)
		}
		r.renderSpanNode(node)
	} else {
		r.WriteString("</span>")
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteByte(lex.ItemDollar)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := node.Tokens
		tokens = html.EscapeHTML(tokens)
		r.Write(tokens)
		r.Tag("/code", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__preview"}, {"data-render", "2"}}, false)
		r.Tag("span", [][]string{{"class", "language-math"}}, false)
		tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
		if node.ParentIs(ast.NodeTableCell) {
			// Improve the `|` render in the inline math in the table https://github.com/Vanessa219/vditor/issues/1550
			tokens = bytes.ReplaceAll(tokens, []byte("\\|"), []byte("|"))
		}
		r.Write(tokens)
		r.Tag("/span", nil, false)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteByte(lex.ItemDollar)
		r.Tag("/span", nil, false)
		r.Tag("code", [][]string{{"data-newline", "1"}, {"class", "vditor-ir__marker vditor-ir__marker--pre"}, {"data-type", "math-inline"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "math-block-close-marker"}}, false)
		r.Write(parse.MathBlockMarker)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	node.Tokens = bytes.TrimSpace(node.Tokens)
	codeLen := len(node.Tokens)
	codeIsEmpty := 1 > codeLen || (len(editor.Caret) == codeLen && editor.Caret == string(node.Tokens))
	class := "vditor-ir__marker--pre"
	if r.Options.VditorMathBlockPreview {
		class += " vditor-ir__marker"
	}
	r.Tag("pre", [][]string{{"class", class}}, false)
	r.Tag("code", [][]string{{"data-type", "math-block"}, {"class", "language-math"}}, false)
	if codeIsEmpty {
		r.WriteString(editor.FrontEndCaret + "\n")
	} else {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	r.WriteString("</code></pre>")

	if r.Options.VditorMathBlockPreview {
		r.Tag("pre", [][]string{{"class", "vditor-ir__preview"}, {"data-render", "2"}}, false)
		r.Tag("div", [][]string{{"data-type", "math-block"}, {"class", "language-math"}}, false)
		tokens := node.Tokens
		tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
		r.Write(html.EscapeHTML(tokens))
		r.WriteString("</div></pre>")
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "math-block-open-marker"}}, false)
		r.Write(parse.MathBlockMarker)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderDivNode(node)
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *VditorIRRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("tr", nil, false)
	} else {
		r.Tag("/tr", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *VditorIRRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("table", [][]string{{"data-block", "0"}, {"data-type", "table"}}, false)
	} else {
		if nil != node.FirstChild.Next {
			r.Tag("/tbody", nil, false)
		}
		r.Tag("/table", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString("~")
		r.Tag("/span", nil, false)
		r.Tag("s", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/s", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString("~")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString("~~")
		r.Tag("/span", nil, false)
		r.Tag("s", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/s", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString("~~")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}

		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--title"}}, false)
		r.WriteByte(lex.ItemDoublequote)
		r.Write(node.Tokens)
		r.WriteByte(lex.ItemDoublequote)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}

		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--link"}}, false)
		dest := node.Tokens
		if r.Options.Sanitize {
			tokens := bytes.TrimSpace(dest)
			tokens = bytes.ToLower(tokens)
			if bytes.HasPrefix(tokens, []byte("javascript:")) {
				dest = nil
			}
		}
		dest = html.EscapeHTML(dest)
		r.Write(dest)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}

		r.WriteByte(lex.ItemSpace)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeImage == node.Parent.Type {
			r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bracket"}}, false)
		} else {
			if 3 == node.Parent.LinkType {
				r.Tag("span", nil, false)
			} else {
				r.Tag("span", [][]string{{"class", "vditor-ir__link"}}, false)
			}
		}
		r.Write(node.Tokens)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}

		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--paren"}}, false)
		r.WriteByte(lex.ItemCloseParen)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeLink == node.Parent.Type && 3 == node.Parent.LinkType {
			return ast.WalkContinue
		}

		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--paren"}}, false)
		r.WriteByte(lex.ItemOpenParen)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderCloseBrace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--brace"}}, false)
		r.WriteByte(lex.ItemCloseBrace)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderOpenBrace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--brace"}}, false)
		r.WriteByte(lex.ItemOpenBrace)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bracket"}}, false)
		r.WriteByte(lex.ItemCloseBracket)
		r.Tag("/span", nil, false)

		if 3 == node.Parent.LinkType {
			linkText := node.Parent.ChildByType(ast.NodeLinkText)
			if nil == linkText || !bytes.EqualFold(node.Parent.LinkRefLabel, linkText.Tokens) {
				r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--link"}}, false)
				r.WriteByte(lex.ItemOpenBracket)
				r.Write(node.Parent.LinkRefLabel)
				r.WriteByte(lex.ItemCloseBracket)
				r.Tag("/span", nil, false)
			}
		}
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bracket"}}, false)
		r.WriteByte(lex.ItemOpenBracket)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteByte(lex.ItemBang)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	needResetCaret := nil != node.Next && ast.NodeText == node.Next.Type && bytes.HasPrefix(node.Next.Tokens, editor.CaretTokens)

	if entering {
		if 3 == node.LinkType {
			node.ChildByType(ast.NodeOpenParen).Unlink()
			node.ChildByType(ast.NodeLinkDest).Unlink()
			if linkSpace := node.ChildByType(ast.NodeLinkSpace); nil != linkSpace {
				linkSpace.Unlink()
				node.ChildByType(ast.NodeLinkTitle).Unlink()
			}
			node.ChildByType(ast.NodeCloseParen).Unlink()
		}

		text := r.Text(node)
		class := "vditor-ir__node"
		if strings.Contains(text, editor.Caret) || needResetCaret {
			class += " vditor-ir__node--expand"
		}
		r.Tag("span", [][]string{{"class", class}, {"data-type", "img"}}, false)
	} else {
		if needResetCaret {
			r.WriteString(editor.Caret)
			node.Next.Tokens = bytes.ReplaceAll(node.Next.Tokens, editor.CaretTokens, nil)
		}

		link := node
		if 3 == node.LinkType {
			link = r.Tree.FindLinkRefDefLink(node.LinkRefLabel)
		}
		destTokens := link.ChildByType(ast.NodeLinkDest).Tokens
		destTokens = r.LinkPath(destTokens)
		destTokens = bytes.ReplaceAll(destTokens, editor.CaretTokens, nil)
		attrs := [][]string{{"src", string(destTokens)}}
		alt := node.ChildByType(ast.NodeLinkText)
		if nil != alt && 0 < len(alt.Tokens) {
			altTokens := bytes.ReplaceAll(alt.Tokens, editor.CaretTokens, nil)
			attrs = append(attrs, []string{"alt", string(altTokens)})
		}
		r.Tag("img", attrs, true)

		// XSS 过滤
		buf := r.Writer.Bytes()
		idx := bytes.LastIndex(buf, []byte("<img src="))
		imgBuf := buf[idx:]
		if r.Options.Sanitize {
			imgBuf = sanitize(imgBuf)
		}
		r.Writer.Truncate(idx)
		r.Writer.Write(imgBuf)

		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderDivNode(node)
		tokens := bytes.TrimSpace(node.Tokens)
		r.WriteString("<pre class=\"vditor-ir__marker--pre vditor-ir__marker\">")
		r.Tag("code", [][]string{{"data-type", "html-block"}}, false)
		r.Write(html.EscapeHTML(tokens))
		r.WriteString("</code></pre>")

		r.Tag("pre", [][]string{{"class", "vditor-ir__preview"}, {"data-render", "2"}}, false)
		tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		r.Write(tokens)
		r.WriteString("</pre></div>")
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	openKbd := bytes.Equal(node.Tokens, []byte("<kbd>"))
	closeKbd := bytes.Equal(node.Tokens, []byte("</kbd>"))
	if openKbd || closeKbd {
		if openKbd {
			if r.tagMatchClose("kbd", node) {
				r.renderSpanNode(node)
				r.Tag("code", [][]string{{"class", "vditor-ir__marker"}}, false)
				r.Write(html.EscapeHTML(node.Tokens))
				r.Tag("/code", nil, false)
				r.Tag("/span", nil, false)
				r.Tag("kbd", nil, false)
			} else {
				r.renderSpanNode(node)
				r.Tag("code", [][]string{{"class", "vditor-ir__marker"}}, false)
				r.Write(html.EscapeHTML(node.Tokens))
				r.Tag("/code", nil, false)
				r.Tag("/span", nil, false)
			}
		} else {
			if r.tagMatchOpen("kbd", node) {
				r.Tag("/kbd", nil, false)
				r.renderSpanNode(node)
				r.Tag("code", [][]string{{"class", "vditor-ir__marker"}}, false)
				r.Write(html.EscapeHTML(node.Tokens))
				r.Tag("/code", nil, false)
				r.Tag("/span", nil, false)
			} else {
				r.renderSpanNode(node)
				r.Tag("code", [][]string{{"class", "vditor-ir__marker"}}, false)
				r.Write(html.EscapeHTML(node.Tokens))
				r.Tag("/code", nil, false)
				r.Tag("/span", nil, false)
			}
		}
	} else {
		r.renderSpanNode(node)
		r.Tag("code", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.Write(html.EscapeHTML(node.Tokens))
		r.Tag("/code", nil, false)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) tagMatchClose(tag string, node *ast.Node) bool {
	for n := node.Next; nil != n; n = n.Next {
		if ast.NodeInlineHTML == n.Type && "</"+tag+">" == n.TokensStr() {
			return true
		}
	}
	return false
}

func (r *VditorIRRenderer) tagMatchOpen(tag string, node *ast.Node) bool {
	for n := node.Previous; nil != n; n = n.Previous {
		if ast.NodeInlineHTML == n.Type && "<"+tag+">" == n.TokensStr() {
			return true
		}
	}
	return false
}

func (r *VditorIRRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if grandparent := node.Parent.Parent; nil != grandparent && ast.NodeList == grandparent.Type && grandparent.ListData.Tight { // List.ListItem.Paragraph
		return ast.WalkContinue
	}

	if entering {
		r.Tag("p", [][]string{{"data-block", "0"}}, false)
	} else {
		r.Tag("/p", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := node.Tokens
		if r.Options.FixTermTypo {
			tokens = r.FixTermTypo(tokens)
		}

		// 有的场景需要零宽空格撑起，但如果有其他文本内容的话需要把零宽空格删掉
		if !bytes.EqualFold(tokens, []byte(editor.Caret+editor.Zwsp)) {
			tokens = bytes.ReplaceAll(tokens, []byte(editor.Zwsp), nil)
		}
		r.Write(html.EscapeHTML(tokens))
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		r.WriteString(strings.Repeat("`", node.Parent.CodeMarkerLen))
		if bytes.HasPrefix(node.Next.Tokens, []byte("`")) {
			r.WriteByte(lex.ItemSpace)
		}
		r.Tag("/span", nil, false)
		r.Tag("code", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/code", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker"}}, false)
		if bytes.HasSuffix(node.Previous.Tokens, []byte("`")) {
			r.WriteByte(lex.ItemSpace)
		}
		r.WriteString(strings.Repeat("`", node.Parent.CodeMarkerLen))
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bi"}}, false)
		r.WriteByte(lex.ItemAsterisk)
		r.Tag("/span", nil, false)
		r.Tag("em", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bi"}}, false)
		r.WriteByte(lex.ItemAsterisk)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bi"}}, false)
		r.WriteByte(lex.ItemUnderscore)
		r.Tag("/span", nil, false)
		r.Tag("em", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bi"}}, false)
		r.WriteByte(lex.ItemUnderscore)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.renderSpanNode(node)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bi"}}, false)
		r.WriteString("**")
		r.Tag("/span", nil, false)
		r.Tag("strong", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bi"}}, false)
		r.WriteString("**")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bi"}}, false)
		r.WriteString("__")
		r.Tag("/span", nil, false)
		r.Tag("strong", [][]string{{"data-newline", "1"}}, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
		r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--bi"}}, false)
		r.WriteString("__")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString(`<blockquote data-block="0">`)
	} else {
		r.WriteString("</blockquote>")
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		text := r.Text(node)
		headingID := node.ChildByType(ast.NodeHeadingID)
		if strings.Contains(text, editor.Caret) || (nil != headingID && bytes.Contains(headingID.Tokens, editor.CaretTokens)) {
			r.WriteString("<h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1] + " data-block=\"0\" class=\"vditor-ir__node vditor-ir__node--expand\"")
		} else {
			r.WriteString("<h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1] + " data-block=\"0\" class=\"vditor-ir__node\"")
		}

		var id string
		if nil != headingID {
			id = string(headingID.Tokens)
		}
		if "" == id {
			id = HeadingID(node)
		}
		r.WriteString(" id=\"ir-" + id + "\"")
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

		if !node.HeadingSetext {
			r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--heading"}, {"data-type", "heading-marker"}}, false)
			r.WriteString(strings.Repeat("#", node.HeadingLevel) + " ")
			r.Tag("/span", nil, false)
		}
	} else {
		if node.HeadingSetext {
			r.Tag("span", [][]string{{"class", "vditor-ir__marker vditor-ir__marker--heading"}, {"data-type", "heading-marker"}, {"data-render", "2"}}, false)
			r.Newline()
			contentLen := r.setextHeadingLen(node)
			if 1 == node.HeadingLevel {
				r.WriteString(strings.Repeat("=", contentLen))
			} else {
				r.WriteString(strings.Repeat("-", contentLen))
			}
			r.Tag("/span", nil, false)
		}
		r.WriteString("</h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1] + ">")
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("span", [][]string{{"data-type", "heading-id"}, {"class", "vditor-ir__marker"}}, false)
		r.WriteString(" {" + string(node.Tokens) + "}")
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *VditorIRRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
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
	} else {
		r.Tag("/li", nil, false)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *VditorIRRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("hr", [][]string{{"data-block", "0"}}, true)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("br", nil, true)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *VditorIRRenderer) renderSpanNode(node *ast.Node) {
	text := r.Text(node)
	var attrs [][]string

	switch node.Type {
	case ast.NodeEmphasis:
		attrs = append(attrs, []string{"data-type", "em"})
	case ast.NodeStrong:
		attrs = append(attrs, []string{"data-type", "strong"})
	case ast.NodeStrikethrough:
		attrs = append(attrs, []string{"data-type", "s"})
	case ast.NodeMark:
		attrs = append(attrs, []string{"data-type", "mark"})
	case ast.NodeSup:
		attrs = append(attrs, []string{"data-type", "sup"})
	case ast.NodeSub:
		attrs = append(attrs, []string{"data-type", "sub"})
	case ast.NodeLink:
		if 3 != node.LinkType {
			attrs = append(attrs, []string{"data-type", "a"})
		} else {
			attrs = append(attrs, []string{"data-type", "link-ref"})
		}
	case ast.NodeImage:
		attrs = append(attrs, []string{"data-type", "img"})
	case ast.NodeCodeSpan:
		attrs = append(attrs, []string{"data-type", "code"})
	case ast.NodeEmoji:
		attrs = append(attrs, []string{"data-type", "emoji"})
	case ast.NodeInlineHTML:
		attrs = append(attrs, []string{"data-type", "html-inline"})
	case ast.NodeHTMLEntity:
		attrs = append(attrs, []string{"data-type", "html-entity"})
	case ast.NodeBackslash:
		attrs = append(attrs, []string{"data-type", "backslash"})
	default:
		attrs = append(attrs, []string{"data-type", "inline-node"})
	}

	if strings.Contains(text, editor.Caret) {
		attrs = append(attrs, []string{"class", "vditor-ir__node vditor-ir__node--expand"})
		r.Tag("span", attrs, false)
		return
	}

	preText := node.PreviousNodeText()
	if strings.HasSuffix(preText, editor.Caret) {
		attrs = append(attrs, []string{"class", "vditor-ir__node vditor-ir__node--expand"})
		r.Tag("span", attrs, false)
		return
	}

	nexText := node.NextNodeText()
	if strings.HasPrefix(nexText, editor.Caret) {
		attrs = append(attrs, []string{"class", "vditor-ir__node vditor-ir__node--expand"})
		r.Tag("span", attrs, false)
		return
	}

	attrs = append(attrs, []string{"class", "vditor-ir__node"})
	r.Tag("span", attrs, false)
	return
}

func (r *VditorIRRenderer) renderDivNode(node *ast.Node) {
	text := r.Text(node)
	attrs := [][]string{{"data-block", "0"}}
	switch node.Type {
	case ast.NodeCodeBlock:
		attrs = append(attrs, []string{"data-type", "code-block"})
	case ast.NodeHTMLBlock:
		attrs = append(attrs, []string{"data-type", "html-block"})
	case ast.NodeMathBlock:
		attrs = append(attrs, []string{"data-type", "math-block"})
	case ast.NodeYamlFrontMatter:
		attrs = append(attrs, []string{"data-type", "yaml-front-matter"})
	}

	if strings.Contains(text, editor.Caret) {
		attrs = append(attrs, []string{"class", "vditor-ir__node vditor-ir__node--expand"})
		r.Tag("div", attrs, false)
		return
	}

	attrs = append(attrs, []string{"class", "vditor-ir__node"})
	r.Tag("div", attrs, false)
	return
}

func (r *VditorIRRenderer) Text(node *ast.Node) (ret string) {
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n.Type {
			case ast.NodeText, ast.NodeLinkText, ast.NodeLinkDest, ast.NodeLinkTitle, ast.NodeCodeBlockCode, ast.NodeCodeSpanContent, ast.NodeInlineMathContent, ast.NodeMathBlockContent, ast.NodeYamlFrontMatterContent, ast.NodeHTMLBlock, ast.NodeInlineHTML, ast.NodeEmojiAlias:
				ret += string(n.Tokens)
			case ast.NodeCodeBlockFenceInfoMarker:
				ret += string(n.CodeBlockInfo)
			case ast.NodeLink:
				if 3 == n.LinkType {
					ret += string(n.LinkRefLabel)
				}
			}
		}
		return ast.WalkContinue
	})
	return
}
