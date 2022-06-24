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
	"github.com/88250/lute/html"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// BlockRenderer 描述了 Protyle WYSIWYG Block DOM 渲染器。
type BlockRenderer struct {
	*BaseRenderer
	NodeIndex int
}

// NewBlockRenderer 创建一个 WYSIWYG Block DOM 渲染器。
func NewBlockRenderer(tree *parse.Tree, options *Options) *BlockRenderer {
	ret := &BlockRenderer{BaseRenderer: NewBaseRenderer(tree, options), NodeIndex: options.NodeIndexStart}
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

func (r *BlockRenderer) renderTextMark(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTextMarkOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<span data-type=\"")
		r.Write(node.Tokens)
		r.WriteString("\">")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTextMarkCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</span>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBr(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<br />")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderUnderline(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderUnderlineOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<u>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderUnderlineCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</u>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderKbd(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderKbdOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<kbd>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderKbdCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</kbd>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		script := node.ChildByType(ast.NodeBlockQueryEmbedScript)
		if nil == script {
			return ast.WalkContinue
		}
		var attrs [][]string
		tokens := script.Tokens
		tokens = html.EscapeHTML(bytes.ReplaceAll(tokens, util.CaretTokens, nil))
		attrs = append(attrs, []string{"data-content", util.BytesToStr(tokens)})
		r.blockNodeAttrs(node, &attrs, "render-node")
		r.Tag("div", attrs, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockQueryEmbedScript(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderVideo(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "iframe")
		r.Tag("div", attrs, false)
		r.Tag("div", [][]string{{"class", "iframe-content"}}, false)
		r.WriteString(parse.Zwsp)
		tokens := bytes.ReplaceAll(node.Tokens, util.CaretTokens, nil)
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		dataSrc := r.tagSrc(tokens)
		src := r.LinkPath(dataSrc)
		tokens = r.replaceSrc(tokens, src, dataSrc)
		r.Write(tokens)
	} else {
		r.Tag("span", [][]string{{"class", "protyle-action__drag"}, {"contenteditable", "false"}}, false)
		r.Tag("/span", nil, false)
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderAudio(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "iframe")
		r.Tag("div", attrs, false)
		r.Tag("div", [][]string{{"class", "iframe-content"}}, false)
		tokens := bytes.ReplaceAll(node.Tokens, util.CaretTokens, nil)
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		dataSrc := r.tagSrc(tokens)
		src := r.LinkPath(dataSrc)
		tokens = r.replaceSrc(tokens, src, dataSrc)
		r.Write(tokens)
		r.WriteString(parse.Zwsp)
	} else {
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderWidget(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "iframe")
		attrs = append(attrs, []string{"data-subtype", "widget"})
		r.Tag("div", attrs, false)
		r.Tag("div", [][]string{{"class", "iframe-content"}}, false)
		tokens := bytes.ReplaceAll(node.Tokens, util.CaretTokens, nil)
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		dataSrc := r.tagSrc(tokens)
		src := r.LinkPath(dataSrc)
		tokens = r.replaceSrc(tokens, src, dataSrc)
		r.Write(tokens)
	} else {
		r.Tag("span", [][]string{{"class", "protyle-action__drag"}, {"contenteditable", "false"}}, false)
		r.Tag("/span", nil, false)
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderIFrame(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "iframe")
		r.Tag("div", attrs, false)
		r.Tag("div", [][]string{{"class", "iframe-content"}}, false)
		tokens := bytes.ReplaceAll(node.Tokens, util.CaretTokens, nil)
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		dataSrc := r.tagSrc(tokens)
		src := r.LinkPath(dataSrc)
		tokens = r.replaceSrc(tokens, src, dataSrc)
		r.Write(tokens)
	} else {
		r.Tag("span", [][]string{{"class", "protyle-action__drag"}, {"contenteditable", "false"}}, false)
		r.Tag("/span", nil, false)
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) replaceSrc(tokens, src, dataSrc []byte) []byte {
	src1 := append([]byte(" src=\""), src...)
	src1 = append(src1, []byte("\"")...)
	dataSrc1 := append([]byte(" src=\""), dataSrc...)
	dataSrc1 = append(dataSrc1, []byte("\"")...)
	if bytes.Contains(tokens, []byte("data-src=")) {
		return bytes.ReplaceAll(tokens, dataSrc1, src1)
	}
	src1 = append(src1, []byte(" data-src=\""+util.BytesToStr(dataSrc)+"\"")...)
	return bytes.ReplaceAll(tokens, dataSrc1, src1)
}

func (r *BlockRenderer) renderBlockRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		idNode := node.ChildByType(ast.NodeBlockRefID)
		var refText, subtype string
		refTextNode := node.ChildByType(ast.NodeBlockRefText)
		subtype = "s"
		if nil == refTextNode {
			refTextNode = node.ChildByType(ast.NodeBlockRefDynamicText)
			subtype = "d"
		}
		if nil != refTextNode {
			refText = refTextNode.Text()
		}
		refText = r.escapeRefText(refText)
		attrs := [][]string{{"data-type", "block-ref"}, {"data-subtype", subtype}, {"data-id", idNode.TokensStr()}}
		r.Tag("span", attrs, false)
		r.WriteString(refText)
		r.Tag("/span", nil, false)
		return ast.WalkSkipChildren
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) escapeRefText(refText string) string {
	refText = strings.ReplaceAll(refText, ">", "&gt;")
	refText = strings.ReplaceAll(refText, "<", "&lt;")
	refText = strings.ReplaceAll(refText, "\"", "&quot;")
	refText = strings.ReplaceAll(refText, "'", "&apos;")
	return refText
}

func (r *BlockRenderer) renderBlockRefID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockRefText(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockRefDynamicText(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderFileAnnotationRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		idNode := node.ChildByType(ast.NodeFileAnnotationRefID)
		id := idNode.TokensStr()
		refText := id
		if refTextNode := node.ChildByType(ast.NodeFileAnnotationRefText); nil != refTextNode {
			refText = refTextNode.Text()
		}
		refText = r.escapeRefText(refText)
		attrs := [][]string{{"data-type", "file-annotation-ref"}, {"data-subtype", "s"}, {"data-id", id}}
		r.Tag("span", attrs, false)
		r.WriteString(refText)
		r.Tag("/span", nil, false)
		return ast.WalkSkipChildren
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderFileAnnotationRefID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderFileAnnotationRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderFileAnnotationRefText(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderGitConflictCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderGitConflictContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "git-conflict")
		r.Tag("div", attrs, false)
		attrs = [][]string{{"contenteditable", "false"}, {"spellcheck", "false"}}
		r.Tag("div", attrs, false)

		tokens := bytes.TrimSpace(node.Tokens)
		r.Write(html.EscapeHTML(tokens))
	} else {
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}

	return ast.WalkContinue
}

func (r *BlockRenderer) renderGitConflictOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderGitConflict(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTag(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTagOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		content := node.Parent.Text()
		content = strings.ReplaceAll(content, util.Caret, "")
		r.Tag("span", [][]string{{"data-type", "tag"}, {"data-content", html.EscapeHTMLStr(content)}}, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTagCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSuperBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if nil == node.FirstChild {
			return ast.WalkContinue
		}

		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "sb")
		layout := node.FirstChild.Next.TokensStr()
		if "" == layout {
			layout = "row"
		}
		attrs = append(attrs, []string{"data-sb-layout", layout})
		r.Tag("div", attrs, false)
	} else {
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSuperBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSuperBlockLayoutMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSuperBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLinkRefDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<div data-block=\"0\" data-type=\"link-ref-defs-block\">")
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLinkRefDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if nil == node.FirstChild {
			return ast.WalkContinue
		}

		dest := node.FirstChild.ChildByType(ast.NodeLinkDest).Tokens
		destStr := util.BytesToStr(dest)
		r.WriteString("[" + util.BytesToStr(node.Tokens) + "]:")
		if util.Caret != destStr {
			r.WriteString(" ")
		}
		r.WriteString(destStr + "\n")
	}
	return ast.WalkSkipChildren
}

func (r *BlockRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderKramdownSpanIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMark1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSupOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSupCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSubOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("sub", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderSubCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sub", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		previewTokens := bytes.TrimSpace(node.Tokens)
		codeLen := len(previewTokens)
		codeIsEmpty := 1 > codeLen || (len(util.Caret) == codeLen && util.Caret == string(node.Tokens))
		r.Tag("pre", nil, false)
		r.Tag("code", [][]string{{"data-type", "yaml-front-matter"}}, false)
		if codeIsEmpty {
			r.WriteString(util.FrontEndCaret + "\n")
		} else {
			r.Write(html.EscapeHTML(previewTokens))
		}
		r.WriteString("</code></pre>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString(`<div class="protyle-wysiwyg__block" data-type="yaml-front-matter" data-block="0">`)
	} else {
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	r.WriteString("<span class=\"protyle-wysiwyg__block\" data-type=\"html-entity\">")
	r.Tag("code", [][]string{{"data-type", "html-entity"}, {"style", "display: none"}}, false)
	tokens := node.HtmlEntityTokens
	r.Write(html.EscapeHTML(tokens))
	r.WriteString("</code>")

	r.Tag("span", [][]string{{"class", "protyle-wysiwyg__preview"}, {"data-render", "2"}}, false)
	r.Tag("code", nil, false)
	previewTokens := bytes.ReplaceAll(node.HtmlEntityTokens, util.CaretTokens, nil)
	r.Write(previewTokens)
	r.Tag("/code", nil, false)
	r.Tag("/span", nil, false)
	r.WriteString("</span>")
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *BlockRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	return r.BaseRenderer.renderToC(node, entering)
}

func (r *BlockRenderer) renderFootnotesDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<div class=\"footnotes-defs-div\">")
		r.WriteString("<hr class=\"footnotes-defs-hr\" />\n")
		r.WriteString("<ol class=\"footnotes-defs-ol\">")
	} else {
		r.WriteString("</ol></div>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		// r.WriteString("<li id=\"footnotes-def-" + node.FootnotesRefId + "\">")
		// 在 li 上带 id 后，Pandoc HTML 转换 Docx 会有问题
		r.WriteString("<li>")
		if 0 < len(node.FootnotesRefs) {
			refId := node.FootnotesRefs[0].FootnotesRefId
			node.FirstChild.PrependChild(&ast.Node{Type: ast.NodeInlineHTML, Tokens: []byte("<span id=\"footnotes-def-" + refId + "\"></span>")})
		}
	} else {
		r.WriteString("</li>\n")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *BlockRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	noHighlight := false
	var language string
	if nil != node.FirstChild && nil != node.FirstChild.Next && 0 < len(node.FirstChild.Next.CodeBlockInfo) {
		language = util.BytesToStr(node.FirstChild.Next.CodeBlockInfo)
		language = strings.ReplaceAll(language, util.Caret, "")
		noHighlight = r.NoHighlight(language)
	}

	if entering {
		if noHighlight {
			if nil == node.FirstChild {
				return ast.WalkContinue
			}

			var attrs [][]string
			r.blockNodeAttrs(node, &attrs, "render-node")
			tokens := html.EscapeHTML(node.FirstChild.Next.Next.Tokens)
			tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
			tokens = bytes.TrimSpace(tokens)
			attrs = append(attrs, []string{"data-content", util.BytesToStr(tokens)})
			attrs = append(attrs, []string{"data-subtype", language})
			r.Tag("div", attrs, false)
			r.Tag("div", [][]string{{"spin", "1"}}, false)
			r.Tag("/div", nil, false)
			r.renderIAL(node)
			r.Tag("/div", nil, false)
			return ast.WalkSkipChildren
		}

		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "code-block")
		r.Tag("div", attrs, false)
	} else {
		if noHighlight {
			return ast.WalkSkipChildren
		}
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	r.Tag("div", [][]string{{"class", "protyle-action"}}, false)
	codeLen := len(node.Tokens)
	codeIsEmpty := 1 > codeLen || (len(util.Caret) == codeLen && util.Caret == string(node.Tokens))
	var language string
	caretInInfo := false
	if nil != node.Previous {
		caretInInfo = bytes.Contains(node.Previous.CodeBlockInfo, util.CaretTokens)
		node.Previous.CodeBlockInfo = bytes.ReplaceAll(node.Previous.CodeBlockInfo, util.CaretTokens, nil)
	}

	attrs := [][]string{{"class", "protyle-action--first protyle-action__language"}, {"contenteditable", "false"}}
	if nil != node.Previous && 0 < len(node.Previous.CodeBlockInfo) {
		infoWords := lex.Split(node.Previous.CodeBlockInfo, lex.ItemSpace)
		language = string(infoWords[0])
	}

	r.Tag("span", attrs, false)
	r.WriteString(language)
	r.Tag("/span", nil, false)
	r.WriteString("<span class=\"fn__flex-1\"></span>")
	r.Tag("span", [][]string{{"class", "protyle-icon protyle-icon--first protyle-action__copy"}}, false)
	r.WriteString("<svg><use xlink:href=\"#iconCopy\"></use></svg>")
	r.Tag("/span", nil, false)
	r.WriteString("<span class=\"protyle-icon protyle-icon--last protyle-action__menu\"><svg><use xlink:href=\"#iconMore\"></use></svg></span>")
	r.Tag("/div", nil, false)

	attrs = [][]string{{"class", "hljs"}}
	r.contenteditable(node, &attrs)
	r.spellcheck(&attrs)
	r.Tag("div", attrs, false)
	if codeIsEmpty {
		if caretInInfo {
			r.WriteString(util.FrontEndCaret)
		}
	} else {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	r.Tag("/div", nil, false)
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *BlockRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := html.EscapeHTML(node.Next.Tokens)
		tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
		r.Tag("span", [][]string{{"data-type", "inline-math"}, {"data-subtype", "math"}, {"data-content", util.BytesToStr(tokens)}, {"contenteditable", "false"}, {"class", "render-node"}}, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/span", nil, false)
		if bytes.Contains(node.Previous.Tokens, util.CaretTokens) {
			r.WriteString(util.Caret)
		}
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	if nil == node.FirstChild {
		return ast.WalkContinue
	}

	var attrs [][]string
	r.blockNodeAttrs(node, &attrs, "render-node")
	tokens := html.EscapeHTML(node.FirstChild.Next.Tokens)
	tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
	tokens = bytes.TrimSpace(tokens)
	attrs = append(attrs, []string{"data-content", util.BytesToStr(tokens)})
	attrs = append(attrs, []string{"data-subtype", "math"})
	r.Tag("div", attrs, false)
	r.Tag("div", [][]string{{"spin", "1"}}, false)
	r.Tag("/div", nil, false)
	r.renderIAL(node)
	r.Tag("/div", nil, false)
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
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
		r.spanNodeAttrs(node, &attrs)
		r.Tag(tag, attrs, false)
	} else {
		r.Tag("/"+tag, nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("tr", nil, false)
	} else {
		r.Tag("/tr", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("colgroup", nil, false)
		if colgroup := node.Parent.IALAttr("colgroup"); "" == colgroup {
			if nil != node.FirstChild {
				for th := node.FirstChild.FirstChild; nil != th; th = th.Next {
					if ast.NodeTableCell == th.Type {
						if style := th.IALAttr("style"); "" != style {
							r.Tag("col", [][]string{{"style", style}}, true)
						} else {
							r.Tag("col", nil, true)
						}
					}
				}
			}
		} else {
			cols := strings.Split(colgroup, "|")
			for _, style := range cols {
				if "" != style {
					r.Tag("col", [][]string{{"style", style}}, true)
				} else {
					r.Tag("col", nil, true)
				}
			}
		}
		r.Tag("/colgroup", nil, false)

		r.Tag("thead", nil, false)
	} else {
		r.Tag("/thead", nil, false)
		r.Tag("tbody", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if nil == node.FirstChild {
		return ast.WalkSkipChildren
	}

	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "table")
		r.Tag("div", attrs, false)
		attrs = [][]string{{"contenteditable", "false"}}
		r.Tag("div", attrs, false)
		attrs = [][]string{}
		r.contenteditable(node, &attrs)
		r.spellcheck(&attrs)
		r.Tag("table", attrs, false)
	} else {
		r.Tag("/tbody", nil, false)
		r.Tag("/table", nil, false)
		r.WriteString("<div class=\"protyle-action__table\"><div class=\"table__resize\"></div><div class=\"table__select\"></div></div>")
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("s", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/s", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("s", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/s", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if ast.NodeImage != node.Parent.Type {
			r.Write(html.EscapeHTML(node.Tokens))
		}
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLess(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderGreater(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCloseBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderOpenBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		attrs := [][]string{{"contenteditable", "false"}, {"data-type", "img"}, {"class", "img"}}
		parentStyle := node.IALAttr("parent-style")
		if "" != parentStyle { // 手动设置了位置
			attrs = append(attrs, []string{"style", parentStyle})
		}
		if !strings.Contains(parentStyle, "display") && !strings.Contains(parentStyle, "block") &&
			r.LastOut == '\n' {
			r.WriteString(parse.Zwsp)
		}
		r.Tag("span", attrs, false)
		r.Tag("span", nil, false)
		r.WriteString(" ")
		r.Tag("/span", nil, false)
		r.Tag("span", nil, false)
		r.Tag("span", [][]string{{"class", "protyle-action protyle-icons"}}, false)
		r.WriteString("<span class=\"protyle-icon protyle-icon--only\"><svg class=\"svg\"><use xlink:href=\"#iconMore\"></use></svg></span>")
		r.Tag("/span", nil, false)
	} else {
		destTokens := node.ChildByType(ast.NodeLinkDest).Tokens
		if r.Options.Sanitize {
			destTokens = sanitize(destTokens)
		}
		destTokens = bytes.ReplaceAll(destTokens, util.CaretTokens, nil)
		dataSrcTokens := destTokens
		dataSrc := util.BytesToStr(dataSrcTokens)
		src := util.BytesToStr(r.LinkPath(destTokens))
		attrs := [][]string{{"src", src}, {"data-src", dataSrc}}
		alt := node.ChildByType(ast.NodeLinkText)
		if nil != alt && 0 < len(alt.Tokens) {
			attrs = append(attrs, []string{"alt", util.BytesToStr(alt.Tokens)})
		}

		title := node.ChildByType(ast.NodeLinkTitle)
		var titleTokens []byte
		if nil != title && 0 < len(title.Tokens) {
			titleTokens = title.Tokens
			attrs = append(attrs, []string{"title", r.escapeRefText(string(titleTokens))})
		}

		if style := node.IALAttr("style"); "" != style {
			attrs = append(attrs, []string{"style", style})
		}
		r.Tag("img", attrs, true)

		buf := r.Writer.Bytes()
		idx := bytes.LastIndex(buf, []byte("<img src="))
		imgBuf := buf[idx:]
		if r.Options.Sanitize {
			imgBuf = sanitize(imgBuf)
		}
		imgBuf = r.tagSrcPath(imgBuf)
		r.Writer.Truncate(idx)
		r.Writer.Write(imgBuf)

		r.Tag("span", [][]string{{"class", "protyle-action__drag"}}, false)
		r.Tag("/span", nil, false)

		if r.Options.ProtyleMarkNetImg && !bytes.HasPrefix(dataSrcTokens, []byte("assets/")) {
			r.WriteString("<span class=\"img__net\"><svg><use xlink:href=\"#iconLanguage\"></use></svg></span>")
		}

		attrs = [][]string{{"class", "protyle-action__title"}}
		r.Tag("span", attrs, false)
		r.Writer.Write(html.EscapeHTML(titleTokens))
		r.Tag("/span", nil, false)
		r.Tag("/span", nil, false)
		r.Tag("span", nil, false)
		r.WriteString(" ")
		r.Tag("/span", nil, false)
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		dest := node.ChildByType(ast.NodeLinkDest)
		destTokens := dest.Tokens
		if r.Options.Sanitize {
			destTokens = bytes.TrimSpace(destTokens)
			destTokens = sanitize(destTokens)
			tokens := bytes.ToLower(destTokens)
			if bytes.HasPrefix(tokens, []byte("javascript:")) {
				destTokens = nil
			}
		}
		destTokens = r.LinkPath(destTokens)
		if bytes.HasPrefix(destTokens, []byte("assets/")) {
			if bytes.Contains(destTokens, []byte("?")) {
				idx := bytes.IndexByte(destTokens, '?')
				d := bytes.ReplaceAll(destTokens[:idx], []byte("#"), []byte("%23"))
				destTokens = append(d, destTokens[idx:]...)
			} else {
				destTokens = bytes.ReplaceAll(destTokens, []byte("#"), []byte("%23"))
			}
		}

		caretInDest := bytes.Contains(destTokens, util.CaretTokens)
		if caretInDest {
			text := node.ChildByType(ast.NodeLinkText)
			text.Tokens = append(text.Tokens, util.CaretTokens...)
			destTokens = bytes.ReplaceAll(destTokens, util.CaretTokens, nil)
		}
		attrs := [][]string{{"data-type", "a"}, {"data-href", string(destTokens)}}
		if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
			attrs = append(attrs, []string{"data-title", r.escapeRefText(string(title.Tokens))})
		}
		r.Tag("span", attrs, false)
	} else {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	var attrs [][]string
	r.blockNodeAttrs(node, &attrs, "render-node")
	tokens := node.Tokens
	tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
	attrs = append(attrs, []string{"data-subtype", "block"})
	r.Tag("div", attrs, false)
	r.WriteString("<div class=\"protyle-icons\">")
	r.WriteString("<span class=\"protyle-icon protyle-icon--first protyle-action__edit\"><svg><use xlink:href=\"#iconEdit\"></use></svg></span><span class=\"protyle-icon protyle-action__menu protyle-icon--last\"><svg><use xlink:href=\"#iconMore\"></use></svg></span>")
	r.WriteString("</div><div>")
	attrs = [][]string{{"data-content", util.BytesToStr(html.EscapeHTML(tokens))}}
	r.Tag("protyle-html", attrs, false)
	r.Tag("/protyle-html", nil, false)
	r.WriteString("<span style=\"position: absolute\">" + parse.Zwsp + "</span>")
	r.WriteString("</div>")
	r.renderIAL(node)
	r.Tag("/div", nil, false)
	return ast.WalkContinue
}

func (r *BlockRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	if bytes.Equal(node.Tokens, []byte("<br />")) && node.ParentIs(ast.NodeTableCell) {
		r.Write(node.Tokens)
		return ast.WalkContinue
	}

	if bytes.Equal(node.Tokens, []byte("<u>")) || bytes.Equal(node.Tokens, []byte("</u>")) {
		r.Write(node.Tokens)
		return ast.WalkContinue
	}

	tokens := html.EscapeHTML(node.Tokens)
	r.Write(tokens)
	return ast.WalkContinue
}

func (r *BlockRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "p")
		r.Tag("div", attrs, false)
		attrs = [][]string{}
		r.contenteditable(node, &attrs)
		r.spellcheck(&attrs)
		r.Tag("div", attrs, false)
	} else {
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var tokens []byte
		if r.Options.AutoSpace && ast.NodeKbd != node.Parent.Type {
			tokens = r.Space(node.Tokens)
		} else {
			tokens = node.Tokens
		}
		r.Write(html.EscapeHTML(tokens))
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *BlockRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("code", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := html.EscapeHTML(node.Tokens)
		r.Write(tokens)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</code>")
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.spanNodeAttrs(node.Parent, &attrs)
		r.Tag("strong", attrs, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "bq")
		r.Tag("div", attrs, false)
	} else {
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		if 6 < node.HeadingLevel {
			node.HeadingLevel = 6
		}
		level := headingLevel[node.HeadingLevel : node.HeadingLevel+1]
		attrs = append(attrs, []string{"data-subtype", "h" + level})
		r.blockNodeAttrs(node, &attrs, "h"+level)
		r.Tag("div", attrs, false)
		attrs = [][]string{}
		r.contenteditable(node, &attrs)
		r.spellcheck(&attrs)
		r.Tag("div", attrs, false)
	} else {
		r.Tag("/div", nil, false)
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *BlockRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		switch node.ListData.Typ {
		case 0:
			attrs = append(attrs, []string{"data-subtype", "u"})
		case 1:
			attrs = append(attrs, []string{"data-subtype", "o"})
		case 3:
			attrs = append(attrs, []string{"data-subtype", "t"})
		}
		r.blockNodeAttrs(node, &attrs, "list")
		r.Tag("div", attrs, false)
	} else {
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		class := "li"
		var attrs [][]string
		switch node.ListData.Typ {
		case 0:
			attrs = append(attrs, []string{"data-marker", "*"})
			attrs = append(attrs, []string{"data-subtype", "u"})
		case 1:
			attrs = append(attrs, []string{"data-marker", strconv.Itoa(node.ListData.Num) + "."})
			attrs = append(attrs, []string{"data-subtype", "o"})
		case 3:
			attrs = append(attrs, []string{"data-marker", "*"})
			attrs = append(attrs, []string{"data-subtype", "t"})
			if node.FirstChild != nil && node.FirstChild.TaskListItemChecked {
				class += " protyle-task--done"
			}
		}
		r.blockNodeAttrs(node, &attrs, class)
		r.Tag("div", attrs, false)

		if 0 == node.ListData.Typ {
			attr := [][]string{{"class", "protyle-action"}, {"draggable", "true"}}
			r.Tag("div", attr, false)
			r.WriteString("<svg><use xlink:href=\"#iconDot\"></use></svg>")
			r.Tag("/div", nil, false)
		} else if 1 == node.ListData.Typ {
			attr := [][]string{{"class", "protyle-action protyle-action--order"}, {"contenteditable", "false"}, {"draggable", "true"}}
			r.Tag("div", attr, false)
			r.WriteString(strconv.Itoa(node.ListData.Num) + ".")
			r.Tag("/div", nil, false)
		}
	} else {
		r.renderIAL(node)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if node.TaskListItemChecked {
			r.WriteString("<div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconCheck\"></use></svg></div>")
		} else {
			r.WriteString("<div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div>")
		}
		if nil == node.Next {
			node.InsertAfter(&ast.Node{ID: ast.NewNodeID(), Type: ast.NodeParagraph})
		}
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
		r.blockNodeAttrs(node, &attrs, "hr")
		r.Tag("div", attrs, false)
		r.Tag("div", nil, false)
	} else {
		r.Tag("/div", nil, false)
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	return r.renderBr(node, entering)
}

func (r *BlockRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemNewline)
		if nil != node.Previous && (ast.NodeStrong == node.Previous.Type ||
			ast.NodeEmphasis == node.Previous.Type ||
			ast.NodeTag == node.Previous.Type ||
			ast.NodeStrikethrough == node.Previous.Type ||
			ast.NodeUnderline == node.Previous.Type ||
			ast.NodeKramdownSpanIAL == node.Previous.Type) &&
			nil != node.Next && bytes.Equal(util.CaretTokens, node.Next.Tokens) {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *BlockRenderer) spanNodeAttrs(node *ast.Node, attrs *[][]string) {
	*attrs = append(*attrs, node.KramdownIAL...)
}

func (r *BlockRenderer) blockNodeAttrs(node *ast.Node, attrs *[][]string, class string) {
	r.nodeID(node, attrs)
	r.nodeIndex(node, attrs)
	r.nodeDataType(node, attrs)
	r.nodeClass(node, attrs, class)

	for _, ial := range node.KramdownIAL {
		if "id" == ial[0] {
			continue
		}
		*attrs = append(*attrs, []string{ial[0], strings.ReplaceAll(ial[1], util.IALValEscNewLine, "\n")})
	}
}

func (r *BlockRenderer) nodeClass(node *ast.Node, attrs *[][]string, class string) {
	*attrs = append(*attrs, []string{"class", class})
}

func (r *BlockRenderer) nodeDataType(node *ast.Node, attrs *[][]string) {
	*attrs = append(*attrs, []string{"data-type", node.Type.String()})
}

func (r *BlockRenderer) nodeID(node *ast.Node, attrs *[][]string) {
	*attrs = append(*attrs, []string{"data-node-id", r.NodeID(node)})
}

func (r *BlockRenderer) nodeIndex(node *ast.Node, attrs *[][]string) {
	if nil == node.Parent || ast.NodeDocument != node.Parent.Type {
		return
	}

	*attrs = append(*attrs, []string{"data-node-index", strconv.Itoa(r.NodeIndex)})
	r.NodeIndex++
	return
}

func (r *BlockRenderer) spellcheck(attrs *[][]string) {
	*attrs = append(*attrs, []string{"spellcheck", "false"})
	return
}

func (r *BlockRenderer) contenteditable(node *ast.Node, attrs *[][]string) {
	if contenteditable := node.IALAttr("contenteditable"); "" != contenteditable {
		*attrs = append(*attrs, []string{"contenteditable", contenteditable})
	} else {
		*attrs = append(*attrs, []string{"contenteditable", strconv.FormatBool(r.Options.ProtyleContenteditable)})
	}
	return
}

func (r *BlockRenderer) renderIAL(node *ast.Node) {
	attrs := [][]string{{"class", "protyle-attr"}, {"contenteditable", "false"}}
	r.Tag("div", attrs, false)

	if bookmark := node.IALAttr("bookmark"); "" != bookmark {
		bookmark = strings.ReplaceAll(bookmark, util.IALValEscNewLine, "\n")
		bookmark = html.EscapeHTMLStr(bookmark)
		r.Tag("div", [][]string{{"class", "protyle-attr--bookmark"}}, false)
		r.WriteString(bookmark)
		r.Tag("/div", nil, false)
	}

	if name := node.IALAttr("name"); "" != name {
		name = strings.ReplaceAll(name, util.IALValEscNewLine, "\n")
		name = html.EscapeHTMLStr(name)
		r.Tag("div", [][]string{{"class", "protyle-attr--name"}}, false)
		r.WriteString("<svg><use xlink:href=\"#iconN\"></use></svg>")
		r.WriteString(name)
		r.Tag("/div", nil, false)
	}

	if alias := node.IALAttr("alias"); "" != alias {
		alias = strings.ReplaceAll(alias, util.IALValEscNewLine, "\n")
		alias = html.EscapeHTMLStr(alias)
		r.Tag("div", [][]string{{"class", "protyle-attr--alias"}}, false)
		r.WriteString("<svg><use xlink:href=\"#iconA\"></use></svg>")
		r.WriteString(alias)
		r.Tag("/div", nil, false)
	}

	if memo := node.IALAttr("memo"); "" != memo {
		memo = strings.ReplaceAll(memo, util.IALValEscNewLine, "\n")
		memo = html.EscapeHTMLStr(memo)
		r.Tag("div", [][]string{{"class", "protyle-attr--memo b3-tooltips b3-tooltips__nw"}, {"aria-label", memo}}, false)
		r.WriteString("<svg><use xlink:href=\"#iconM\"></use></svg>")
		r.Tag("/div", nil, false)
	}

	if refCount := node.IALAttr("refcount"); "" != refCount {
		refCount = strings.ReplaceAll(refCount, util.IALValEscNewLine, "\n")
		refCount = html.EscapeHTMLStr(refCount)
		r.Tag("div", [][]string{{"class", "protyle-attr--refcount popover__block"}}, false)
		r.WriteString(refCount)
		r.Tag("/div", nil, false)
	}

	r.WriteString(parse.Zwsp)
	r.Tag("/div", nil, false)
}
