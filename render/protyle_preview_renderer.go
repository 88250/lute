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
	"unicode"
	"unicode/utf8"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

type ProtylePreviewRenderer struct {
	*BaseRenderer
}

func NewProtylePreviewRenderer(tree *parse.Tree, options *Options) *ProtylePreviewRenderer {
	ret := &ProtylePreviewRenderer{NewBaseRenderer(tree, options)}
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
	ret.RendererFuncs[ast.NodeTextMarkOpenMarker] = ret.renderTextMarkOpenMarker
	ret.RendererFuncs[ast.NodeTextMarkCloseMarker] = ret.renderTextMarkCloseMarker
	return ret
}

func (r *ProtylePreviewRenderer) renderTextMark(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderTextMarkOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<span data-type=\"")
		r.Write(node.Tokens)
		r.WriteString("\">")
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderTextMarkCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</span>")
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderBr(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<br />")
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderUnderline(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderUnderlineOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<u>")
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderUnderlineCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</u>")
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderKbd(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderKbdOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<kbd>")
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderKbdCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</kbd>")
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderVideo(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderAudio(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderIFrame(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderWidget(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) Render() (output []byte) {
	output = r.BaseRenderer.Render()
	return
}

func (r *ProtylePreviewRenderer) renderGitConflictCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderGitConflictContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderGitConflictOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderGitConflict(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderSuperBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderSuperBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *ProtylePreviewRenderer) renderSuperBlockLayoutMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *ProtylePreviewRenderer) renderSuperBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *ProtylePreviewRenderer) renderLinkRefDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *ProtylePreviewRenderer) renderLinkRefDef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *ProtylePreviewRenderer) renderTag(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderTagOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", node.Parent.KramdownIAL, false)
		r.WriteByte(lex.ItemCrosshatch)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderTagCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemCrosshatch)
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderKramdownSpanIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderMark1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("mark", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("mark", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderSupOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderSupCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderSubOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("sub", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderSubCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sub", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("div", nil, false)
	} else {
		r.Tag("/div", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderBlockQueryEmbedScript(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(node.Tokens)
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderBlockRef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderBlockRefID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderBlockRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderBlockRefText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(node.Tokens)
	} else {
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderBlockRefDynamicText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemSinglequote)
		r.Write(node.Tokens)
	} else {
		r.WriteByte(lex.ItemSinglequote)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderFileAnnotationRef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderFileAnnotationRefID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderFileAnnotationRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderFileAnnotationRefText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(node.Tokens)
	} else {
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</code></pre>")
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		attrs := [][]string{{"class", "vditor-yml-front-matter"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("pre", attrs, false)
		r.WriteString("<code class=\"language-yaml\">")
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	return r.BaseRenderer.renderToC(node, entering)
}

func (r *ProtylePreviewRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		idx, _ := r.Tree.FindFootnotesDef(node.Tokens)
		idxStr := strconv.Itoa(idx)
		r.Tag("sup", [][]string{{"class", "footnotes-ref"}, {"id", "footnotes-ref-" + node.FootnotesRefId}}, false)
		r.Tag("a", [][]string{{"href", "#footnotes-def-" + idxStr}}, false)
		r.WriteString(idxStr)
		r.Tag("/a", nil, false)
		r.Tag("/sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderFootnotesDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<div class=\"footnotes-defs-div\">")
		r.WriteString("<hr class=\"footnotes-defs-hr\" />\n")
		r.WriteString("<ol class=\"footnotes-defs-ol\">")
	} else {
		r.WriteString("</ol></div>")
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		// r.WriteString("<li id=\"footnotes-def-" + node.FootnotesRefId + "\">")
		// 在 li 上带 id 后，Pandoc HTML 转换 Docx 会有问题
		r.WriteString("<li>")
		if 0 < len(node.FootnotesRefs) && nil != node.FirstChild {
			refId := node.FootnotesRefs[0].FootnotesRefId
			node.FirstChild.PrependChild(&ast.Node{Type: ast.NodeInlineHTML, Tokens: []byte("<span id=\"footnotes-def-" + refId + "\"></span>")})
		}
	} else {
		r.WriteString("</li>\n")
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()

	noHighlight := false
	var language string
	if nil != node.FirstChild.Next && 0 < len(node.FirstChild.Next.CodeBlockInfo) {
		language = util.BytesToStr(node.FirstChild.Next.CodeBlockInfo)
		noHighlight = r.NoHighlight(language)
	}

	if entering {
		if noHighlight {
			var attrs [][]string
			tokens := html.EscapeHTML(node.FirstChild.Next.Next.Tokens)
			tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
			tokens = bytes.TrimSpace(tokens)
			attrs = append(attrs, []string{"data-content", util.BytesToStr(tokens)})
			attrs = append(attrs, []string{"data-subtype", language})
			r.Tag("div", attrs, false)
			r.Tag("div", [][]string{{"spin", "1"}}, false)
			r.Tag("/div", nil, false)
			r.Tag("/div", nil, false)
			return ast.WalkSkipChildren
		}

		attrs := [][]string{{"class", "code-block"}, {"data-language", language}}
		r.Tag("pre", attrs, false)
		r.WriteString("<code class=\"hljs\">")
	} else {
		if noHighlight {
			return ast.WalkSkipChildren
		}

		r.Tag("/code", nil, false)
		r.Tag("/pre", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := html.EscapeHTML(node.Next.Tokens)
		r.Tag("span", [][]string{{"data-subtype", "math"}, {"data-content", util.BytesToStr(tokens)}}, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if entering {
		var attrs [][]string
		tokens := html.EscapeHTML(node.FirstChild.Next.Tokens)
		tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
		tokens = bytes.TrimSpace(tokens)
		attrs = append(attrs, []string{"data-content", util.BytesToStr(tokens)})
		attrs = append(attrs, []string{"data-subtype", "math"})
		r.Tag("div", attrs, false)
		r.Tag("div", [][]string{{"spin", "1"}}, false)
		r.Tag("/div", nil, false)
		r.Tag("/div", nil, false)
		r.Newline()
	}
	return ast.WalkSkipChildren
}

func (r *ProtylePreviewRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("tr", nil, false)
		r.Newline()
	} else {
		r.Tag("/tr", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("del", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/del", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("del", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/del", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderCloseBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderOpenBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderLess(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderGreater(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if 0 == r.DisableTags {
			attrs := [][]string{{"class", "img"}}
			if style := node.IALAttr("parent-style"); "" != style {
				attrs = append(attrs, []string{"style", style})
			}
			r.Tag("span", attrs, false)
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
		title := node.ChildByType(ast.NodeLinkTitle)
		var titleTokens []byte
		if nil != title && nil != title.Tokens {
			titleTokens = html.EscapeHTML(title.Tokens)
			r.WriteString(" title=\"")
			r.Write(titleTokens)
			r.WriteByte(lex.ItemDoublequote)
		}
		ial := r.NodeAttrsStr(node)
		if "" != ial {
			r.WriteString(" " + ial)
		}
		r.WriteString(" />")
		if 0 < len(titleTokens) {
			r.Tag("span", [][]string{{"class", "protyle-action__title"}}, false)
			r.Write(titleTokens)
			r.Tag("/span", nil, false)
		}
		r.Tag("/span", nil, false)

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

func (r *ProtylePreviewRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := node.Tokens
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		r.Write(tokens)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.handleKramdownBlockIAL(node)
		var attrs [][]string
		attrs = append(attrs, node.KramdownIAL...)
		r.Tag("p", attrs, false)
		if r.Options.ChineseParagraphBeginningSpace && ast.NodeDocument == node.Parent.Type {
			r.WriteString("　　")
		}
	} else {
		r.Tag("/p", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("code", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/code", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("strong", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("strong", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.handleKramdownBlockIAL(node)
		attrs = append(attrs, node.KramdownIAL...)
		if 3 == node.ListData.Typ && nil != node.FirstChild && ((ast.NodeTaskListItemMarker == node.FirstChild.Type) ||
			(nil != node.FirstChild.FirstChild && ast.NodeTaskListItemMarker == node.FirstChild.FirstChild.Type)) {
			taskListItemMarker := node.FirstChild.FirstChild
			if nil == taskListItemMarker {
				taskListItemMarker = node.FirstChild
			}
			taskClass := "protyle-task"
			if taskListItemMarker.TaskListItemChecked {
				taskClass += " protyle-task--done"
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

func (r *ProtylePreviewRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("hr", nil, true)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("br", nil, true)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtylePreviewRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtylePreviewRenderer) handleKramdownBlockIAL(node *ast.Node) {
	if r.Options.KramdownBlockIAL && "id" != r.Options.KramdownIALIDRenderName && 0 < len(node.KramdownIAL) {
		// 第一项必须是 ID
		node.KramdownIAL[0][0] = r.Options.KramdownIALIDRenderName
	}
}
