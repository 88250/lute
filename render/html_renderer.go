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

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// HtmlRenderer 描述了 HTML 渲染器。
type HtmlRenderer struct {
	*BaseRenderer
}

// NewHtmlRenderer 创建一个 HTML 渲染器。
func NewHtmlRenderer(tree *parse.Tree, options *Options) *HtmlRenderer {
	ret := &HtmlRenderer{NewBaseRenderer(tree, options)}
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
	ret.RendererFuncs[ast.NodeLess] = ret.renderLess
	ret.RendererFuncs[ast.NodeGreater] = ret.renderGreater
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
	ret.RendererFuncs[ast.NodeBlockRef] = ret.renderBlockRef
	ret.RendererFuncs[ast.NodeBlockRefID] = ret.renderBlockRefID
	ret.RendererFuncs[ast.NodeBlockRefSpace] = ret.renderBlockRefSpace
	ret.RendererFuncs[ast.NodeBlockRefText] = ret.renderBlockRefText
	ret.RendererFuncs[ast.NodeBlockRefDynamicText] = ret.renderBlockRefDynamicText
	ret.RendererFuncs[ast.NodeFileAnnotationRef] = ret.renderFileAnnotationRef
	ret.RendererFuncs[ast.NodeFileAnnotationRefID] = ret.renderFileAnnotationRefID
	ret.RendererFuncs[ast.NodeFileAnnotationRefSpace] = ret.renderFileAnnotationRefSpace
	ret.RendererFuncs[ast.NodeFileAnnotationRefText] = ret.renderFileAnnotationRefText
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
	ret.RendererFuncs[ast.NodeKramdownSpanIAL] = ret.renderKramdownSpanIAL
	ret.RendererFuncs[ast.NodeBlockQueryEmbed] = ret.renderBlockQueryEmbed
	ret.RendererFuncs[ast.NodeBlockQueryEmbedScript] = ret.renderBlockQueryEmbedScript
	ret.RendererFuncs[ast.NodeTag] = ret.renderTag
	ret.RendererFuncs[ast.NodeTagOpenMarker] = ret.renderTagOpenMarker
	ret.RendererFuncs[ast.NodeTagCloseMarker] = ret.renderTagCloseMarker
	ret.RendererFuncs[ast.NodeLinkRefDefBlock] = ret.renderLinkRefDefBlock
	ret.RendererFuncs[ast.NodeLinkRefDef] = ret.renderLinkRefDef
	ret.RendererFuncs[ast.NodeSuperBlock] = ret.renderSuperBlock
	ret.RendererFuncs[ast.NodeSuperBlockOpenMarker] = ret.renderSuperBlockOpenMarker
	ret.RendererFuncs[ast.NodeSuperBlockLayoutMarker] = ret.renderSuperBlockLayoutMarker
	ret.RendererFuncs[ast.NodeSuperBlockCloseMarker] = ret.renderSuperBlockCloseMarker
	ret.RendererFuncs[ast.NodeGitConflict] = ret.renderGitConflict
	ret.RendererFuncs[ast.NodeGitConflictOpenMarker] = ret.renderGitConflictOpenMarker
	ret.RendererFuncs[ast.NodeGitConflictContent] = ret.renderGitConflictContent
	ret.RendererFuncs[ast.NodeGitConflictCloseMarker] = ret.renderGitConflictCloseMarker
	ret.RendererFuncs[ast.NodeIFrame] = ret.renderIFrame
	ret.RendererFuncs[ast.NodeWidget] = ret.renderWidget
	ret.RendererFuncs[ast.NodeVideo] = ret.renderVideo
	ret.RendererFuncs[ast.NodeAudio] = ret.renderAudio
	ret.RendererFuncs[ast.NodeKbd] = ret.renderKbd
	ret.RendererFuncs[ast.NodeKbdOpenMarker] = ret.renderKbdOpenMarker
	ret.RendererFuncs[ast.NodeKbdCloseMarker] = ret.renderKbdCloseMarker
	ret.RendererFuncs[ast.NodeUnderline] = ret.renderUnderline
	ret.RendererFuncs[ast.NodeUnderlineOpenMarker] = ret.renderUnderlineOpenMarker
	ret.RendererFuncs[ast.NodeUnderlineCloseMarker] = ret.renderUnderlineCloseMarker
	ret.RendererFuncs[ast.NodeBr] = ret.renderBr
	ret.RendererFuncs[ast.NodeTextMark] = ret.renderTextMark
	ret.RendererFuncs[ast.NodeAttributeView] = ret.renderAttributeView
	ret.RendererFuncs[ast.NodeCustomBlock] = ret.renderCustomBlock
	return ret
}

func (r *HtmlRenderer) Render() (output []byte) {
	output = r.BaseRenderer.Render()
	output = append(output, r.RenderFootnotes()...)
	return
}

func (r *HtmlRenderer) renderCustomBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("div", [][]string{
			{"data-type", "NodeCustomBlock"},
			{"data-info", node.CustomBlockInfo},
			{"data-content", string(html.EscapeHTML(node.Tokens))},
		}, false)
		r.WriteString("</div>")
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderAttributeView(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("div", [][]string{
			{"data-type", "NodeAttributeView"},
			{"data-av-id", node.AttributeViewID},
			{"data-av-type", node.AttributeViewType},
		}, false)
		r.WriteString("</div>")
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderTextMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		textContent := node.TextMarkTextContent
		if node.ParentIs(ast.NodeTableCell) {
			if node.IsTextMarkType("code") {
				textContent = strings.ReplaceAll(textContent, "|", "&#124;")
			} else {
				textContent = strings.ReplaceAll(textContent, "\\|", "|")
			}
			textContent = strings.ReplaceAll(textContent, "\n", "<br />")
		}

		if node.IsTextMarkType("a") {
			attrs := [][]string{{"href", node.TextMarkAHref}}
			if "" != node.TextMarkATitle {
				attrs = append(attrs, []string{"title", node.TextMarkATitle})
			}
			r.spanNodeAttrs(node, &attrs)
			r.Tag("a", attrs, false)
			r.WriteString(textContent)
			r.WriteString("</a>")
		} else if node.IsTextMarkType("inline-memo") {
			r.WriteString(textContent)
			lastRune, _ := utf8.DecodeLastRuneInString(node.TextMarkTextContent)
			if isCJK(lastRune) {
				r.WriteString("<sup>（")
				memo := node.TextMarkInlineMemoContent
				memo = strings.ReplaceAll(memo, editor.IALValEscNewLine, " ")
				r.WriteString(memo)
				r.WriteString("）</sup>")
			} else {
				r.WriteString("<sup>(")
				r.WriteString(node.TextMarkInlineMemoContent)
				r.WriteString(")</sup>")
			}
		} else {
			attrs := r.renderTextMarkAttrs(node)
			r.spanNodeAttrs(node, &attrs)
			r.Tag("span", attrs, false)
			r.WriteString(textContent)
			r.WriteString("</span>")
		}
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBr(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<br />")
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderUnderline(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderUnderlineOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<u>")
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderUnderlineCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</u>")
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderKbd(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderKbdOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<kbd>")
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderKbdCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</kbd>")
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderVideo(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("div", [][]string{{"class", "iframe"}}, false)
		tokens := node.Tokens
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		tokens = r.tagSrcPath(tokens)
		r.Write(tokens)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderAudio(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("div", [][]string{{"class", "iframe"}}, false)
		tokens := node.Tokens
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		tokens = r.tagSrcPath(tokens)
		r.Write(tokens)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderIFrame(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("div", [][]string{{"class", "iframe"}}, false)
		tokens := node.Tokens
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		tokens = r.tagSrcPath(tokens)
		r.Write(tokens)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderWidget(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("div", [][]string{{"class", "iframe"}}, false)
		tokens := node.Tokens
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		tokens = r.tagSrcPath(tokens)
		r.Write(tokens)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderGitConflictCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderGitConflictContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderGitConflictOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderGitConflict(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if entering {
		attrs := [][]string{{"class", "language-git-conflict"}}
		r.handleKramdownBlockIAL(node)
		attrs = append(attrs, node.KramdownIAL...)
		r.Tag("div", attrs, false)
	} else {
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderSuperBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderSuperBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *HtmlRenderer) renderSuperBlockLayoutMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *HtmlRenderer) renderSuperBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *HtmlRenderer) renderLinkRefDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *HtmlRenderer) renderLinkRefDef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
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
	if entering {
		r.Tag("em", node.Parent.KramdownIAL, false)
		r.WriteByte(lex.ItemCrosshatch)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderTagCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemCrosshatch)
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderKramdownSpanIAL(node *ast.Node, entering bool) ast.WalkStatus {
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
	if entering {
		r.Tag("mark", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("mark", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderSupOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderSupCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderSubOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("sub", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderSubCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sub", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("div", nil, false)
	} else {
		r.Tag("/div", nil, false)
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
	} else {
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockRefDynamicText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemSinglequote)
		r.Write(html.EscapeHTML(node.Tokens))
	} else {
		r.WriteByte(lex.ItemSinglequote)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderFileAnnotationRef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderFileAnnotationRefID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderFileAnnotationRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderFileAnnotationRefText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(html.EscapeHTML(node.Tokens))
	} else {
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</code></pre>")
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		attrs := [][]string{{"class", "vditor-yml-front-matter"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("pre", attrs, false)
		r.WriteString("<code class=\"language-yaml\">")
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	return r.BaseRenderer.renderToC(node, entering)
}

func (r *HtmlRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		idx, _ := r.Tree.FindFootnotesDef(node.Tokens)
		idxStr := strconv.Itoa(idx)
		r.Tag("sup", [][]string{{"class", "footnotes-ref"}, {"id", "footnotes-ref-" + node.FootnotesRefId}}, false)
		r.Tag("a", [][]string{{"href", r.Options.LinkBase + "#footnotes-def-" + idxStr}}, false)
		r.WriteString(idxStr)
		r.Tag("/a", nil, false)
		r.Tag("/sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderFootnotesDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) RenderFootnotes() []byte {
	if 1 > len(r.FootnotesDefs) {
		return nil
	}

	buf := bytes.Buffer{}
	buf.WriteString("<div class=\"footnotes-defs-div\">")
	buf.WriteString("<hr class=\"footnotes-defs-hr\" />\n")
	buf.WriteString("<ol class=\"footnotes-defs-ol\">")
	for i, def := range r.FootnotesDefs {
		buf.WriteString("<li id=\"footnotes-def-" + strconv.Itoa(i+1) + "\">")
		footnotesTree := &parse.Tree{Name: "", Context: r.Tree.Context}
		footnotesTree.Context.Tree = footnotesTree
		footnotesTree.Root = &ast.Node{Type: ast.NodeDocument}
		footnotesTree.Root.AppendChild(def)
		defRenderer := NewHtmlRenderer(footnotesTree, r.Options)
		lc := footnotesTree.Root.LastDeepestChild()
		for i = len(def.FootnotesRefs) - 1; 0 <= i; i-- {
			ref := def.FootnotesRefs[i]
			gotoRef := " <a href=\"#footnotes-ref-" + ref.FootnotesRefId + "\" class=\"vditor-footnotes__goto-ref\">↩</a>"
			link := &ast.Node{Type: ast.NodeInlineHTML, Tokens: util.StrToBytes(gotoRef)}
			lc.InsertAfter(link)
		}
		defRenderer.RenderingFootnotes = true
		defContent := defRenderer.Render()
		buf.Write(defContent)
		buf.WriteString("</li>\n")
	}
	buf.WriteString("</ol></div>")
	return buf.Bytes()
}

func (r *HtmlRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if !r.RenderingFootnotes {
			var found bool
			for _, n := range r.FootnotesDefs {
				if bytes.EqualFold(node.Tokens, n.Tokens) {
					found = true
					break
				}
			}
			if !found {
				r.FootnotesDefs = append(r.FootnotesDefs, node)
			}
			return ast.WalkSkipChildren
		}
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := node.Tokens
		if node.ParentIs(ast.NodeTableCell) {
			// Improve the `|` render in the inline math in the table https://github.com/Vanessa219/vditor/issues/1550
			tokens = bytes.ReplaceAll(tokens, []byte("\\|"), []byte("|"))
		}
		r.Write(html.EscapeHTML(tokens))
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		attrs := [][]string{{"class", "language-math"}}
		r.Tag("span", attrs, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if entering {
		attrs := [][]string{{"class", "language-math"}}
		r.handleKramdownBlockIAL(node)
		attrs = append(attrs, node.KramdownIAL...)
		r.Tag("div", attrs, false)
	}
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
		r.Tag(tag, attrs, false)
	} else {
		r.Tag("/"+tag, nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("tr", nil, false)
		r.Newline()
	} else {
		r.Tag("/tr", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("thead", nil, false)
		r.Newline()
	} else {
		r.Tag("/thead", nil, false)
		r.Newline()
		if nil != node.Next {
			r.Tag("tbody", nil, false)
		}
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.handleKramdownBlockIAL(node)
		r.Tag("table", node.KramdownIAL, false)
		r.Newline()
	} else {
		if nil != node.FirstChild.Next {
			r.Tag("/tbody", nil, false)
		}
		r.Newline()
		r.Tag("/table", nil, false)
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
	if entering {
		r.Tag("del", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/del", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("del", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/del", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var tokens []byte
		if r.Options.AutoSpace {
			tokens = r.Space(node.Tokens)
		} else {
			tokens = node.Tokens
		}
		r.Write(html.EscapeHTML(tokens))
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderCloseBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderOpenBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderLess(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderGreater(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if 0 == r.DisableTags {
			if style := node.IALAttr("parent-style"); "" != style {
				r.Tag("span", [][]string{{"style", style}}, false)
			}

			r.WriteString("<img src=\"")
			destTokens := node.ChildByType(ast.NodeLinkDest).Tokens
			destTokens = r.LinkPath(destTokens)
			if "" != r.Options.ImageLazyLoading {
				r.Write(html.EscapeHTML(util.StrToBytes(r.Options.ImageLazyLoading)))
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
		ial := r.NodeAttrsStr(node)
		if "" != ial {
			r.WriteString(" " + ial)
		}
		r.WriteString(" />")
		if style := node.IALAttr("style"); "" != style {
			r.Tag("/span", nil, false)
		}

		if r.Options.Sanitize {
			buf := r.Writer.Bytes()
			idx := bytes.LastIndex(buf, []byte("<img src="))
			imgBuf := buf[idx:]
			if r.Options.Sanitize {
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
		if r.Options.Sanitize {
			tokens := bytes.TrimSpace(destTokens)
			tokens = bytes.ToLower(tokens)
			if bytes.HasPrefix(tokens, []byte("javascript:")) {
				destTokens = nil
			}
		}
		destTokens = r.LinkPath(destTokens)
		attrs := [][]string{{"href", util.BytesToStr(html.EscapeHTML(destTokens))}}
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			attrs = append(attrs, []string{"title", util.BytesToStr(html.EscapeHTML(title.Tokens))})
		}
		r.Tag("a", attrs, false)
	} else {
		r.Tag("/a", nil, false)

		r.LinkTextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		tokens := node.Tokens
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		tokens = r.tagSrcPath(tokens)
		r.Write(tokens)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := node.Tokens
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		r.Write(tokens)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if grandparent := node.Parent.Parent; nil != grandparent && ast.NodeList == grandparent.Type && grandparent.ListData.Tight { // List.ListItem.Paragraph
		return ast.WalkContinue
	}

	if entering {
		r.Newline()
		r.handleKramdownBlockIAL(node)
		var attrs [][]string
		attrs = append(attrs, node.KramdownIAL...)
		if r.Options.ChineseParagraphBeginningSpace && ast.NodeDocument == node.Parent.Type {
			if !r.ParagraphContainImgOnly(node) {
				attrs = append(attrs, []string{"class", "indent--2"})
			}
		}
		r.Tag("p", attrs, false)
	} else {
		r.Tag("/p", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var tokens []byte
		if r.Options.AutoSpace {
			tokens = r.Space(node.Tokens)
		} else {
			tokens = node.Tokens
		}

		if r.Options.FixTermTypo {
			tokens = r.FixTermTypo(tokens)
		}
		r.Write(html.EscapeHTML(tokens))
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if r.Options.AutoSpace {
			if text := node.PreviousNodeText(); "" != text {
				lastc, _ := utf8.DecodeLastRuneInString(text)
				if unicode.IsLetter(lastc) || unicode.IsDigit(lastc) {
					r.WriteByte(lex.ItemSpace)
				}
			}
		}
	} else {
		if r.Options.AutoSpace {
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
	if entering {
		r.Tag("code", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/code", nil, false)
	}
	return ast.WalkContinue
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
	if entering {
		r.Tag("em", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
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
	if entering {
		r.Tag("strong", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("strong", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.handleKramdownBlockIAL(node)
		r.Tag("blockquote", node.KramdownIAL, false)
		r.Newline()
	} else {
		r.Newline()
		r.WriteString("</blockquote>")
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

const headingLevel = " 123456"

func (r *HtmlRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		level := headingLevel[node.HeadingLevel : node.HeadingLevel+1]
		r.WriteString("<h" + level)
		id := HeadingID(node)
		if r.Options.ToC || r.Options.HeadingID || r.Options.KramdownBlockIAL {
			r.WriteString(" id=\"" + id + "\"")
			if r.Options.KramdownBlockIAL {
				if "id" != r.Options.KramdownIALIDRenderName && 0 < len(node.KramdownIAL) {
					r.WriteString(" " + r.Options.KramdownIALIDRenderName + "=\"" + node.KramdownIAL[0][1] + "\"")
				}
				if 1 < len(node.KramdownIAL) {
					exceptID := node.KramdownIAL[1:]
					for _, attr := range exceptID {
						r.WriteString(" " + attr[0] + "=\"" + attr[1] + "\"")
					}
				}
			}
		}
		r.WriteString(">")
	} else {
		if r.Options.HeadingAnchor {
			id := HeadingID(node)
			r.Tag("a", [][]string{{"id", "vditorAnchor-" + id}, {"class", "vditor-anchor"}, {"href", "#" + id}}, false)
			r.WriteString(`<svg viewBox="0 0 16 16" version="1.1" width="16" height="16"><path fill-rule="evenodd" d="M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z"></path></svg>`)
			r.Tag("/a", nil, false)
		}
		r.WriteString("</h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1] + ">")
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
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
		if 0 == node.ListData.BulletChar && 1 != node.ListData.Start {
			attrs = append(attrs, []string{"start", strconv.Itoa(node.ListData.Start)})
		}
		r.handleKramdownBlockIAL(node)
		attrs = append(attrs, node.KramdownIAL...)
		r.Tag(tag, attrs, false)
		r.Newline()
	} else {
		r.Newline()
		r.Tag("/"+tag, nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.handleKramdownBlockIAL(node)
		attrs = append(attrs, node.KramdownIAL...)
		if 3 == node.ListData.Typ && "" != r.Options.GFMTaskListItemClass && nil != node.FirstChild &&
			((ast.NodeTaskListItemMarker == node.FirstChild.Type) ||
				(nil != node.FirstChild.FirstChild && ast.NodeTaskListItemMarker == node.FirstChild.FirstChild.Type)) {
			taskListItemMarker := node.FirstChild.FirstChild
			if nil == taskListItemMarker {
				taskListItemMarker = node.FirstChild
			}
			taskClass := r.Options.GFMTaskListItemClass
			if taskListItemMarker.TaskListItemChecked {
				taskClass += " vditor-task--done"
			}
			attrs = append(attrs, []string{"class", taskClass})
		}
		r.Tag("li", attrs, false)
	} else {
		r.Tag("/li", nil, false)
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
		r.Tag("input", attrs, true)
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("hr", nil, true)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("br", nil, true)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if r.Options.SoftBreak2HardBreak {
			r.Tag("br", nil, true)
			r.Newline()
		} else {
			r.Newline()
		}
	}
	return ast.WalkContinue
}

func (r *HtmlRenderer) handleKramdownBlockIAL(node *ast.Node) {
	if r.Options.KramdownBlockIAL && "id" != r.Options.KramdownIALIDRenderName && 0 < len(node.KramdownIAL) {
		// 第一项必须是 ID
		node.KramdownIAL[0][0] = r.Options.KramdownIALIDRenderName
	}
}

func (r *HtmlRenderer) renderTextMarkAttrs(node *ast.Node) (attrs [][]string) {
	attrs = [][]string{{"data-type", node.TextMarkType}}

	types := strings.Split(node.TextMarkType, " ")
	for _, typ := range types {
		if "block-ref" == typ {
			attrs = append(attrs, []string{"data-subtype", node.TextMarkBlockRefSubtype})
			attrs = append(attrs, []string{"data-id", node.TextMarkBlockRefID})
		} else if "a" == typ {
			href := node.TextMarkAHref
			href = string(r.LinkPath([]byte(href)))

			attrs = append(attrs, []string{"data-href", href})
			if "" != node.TextMarkATitle {
				attrs = append(attrs, []string{"data-title", node.TextMarkATitle})
			}
		} else if "inline-math" == typ {
			attrs = append(attrs, []string{"data-subtype", "math"})
			content := node.TextMarkInlineMathContent
			if node.ParentIs(ast.NodeTableCell) {
				// Improve the handling of inline-math containing `|` in the table https://github.com/siyuan-note/siyuan/issues/9227
				content = strings.ReplaceAll(content, "|", "&#124;")
				content = strings.ReplaceAll(content, "\n", "<br/>")
			}
			content = strings.ReplaceAll(content, editor.IALValEscNewLine, "\n")
			attrs = append(attrs, []string{"data-content", content})
			attrs = append(attrs, []string{"contenteditable", "false"})
			attrs = append(attrs, []string{"class", "render-node"})
		} else if "file-annotation-ref" == typ {
			attrs = append(attrs, []string{"data-id", node.TextMarkFileAnnotationRefID})
		} else if "inline-memo" == typ {
			content := node.TextMarkInlineMemoContent
			content = strings.ReplaceAll(content, editor.IALValEscNewLine, "\n")
			attrs = append(attrs, []string{"data-inline-memo-content", content})
		}
	}
	return
}

func (r *HtmlRenderer) spanNodeAttrs(node *ast.Node, attrs *[][]string) {
	*attrs = append(*attrs, node.KramdownIAL...)
}
