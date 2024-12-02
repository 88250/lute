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

type ProtyleExportDocxRenderer struct {
	*BaseRenderer
}

func NewProtyleExportDocxRenderer(tree *parse.Tree, options *Options) *ProtyleExportDocxRenderer {
	ret := &ProtyleExportDocxRenderer{NewBaseRenderer(tree, options)}
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

func (r *ProtyleExportDocxRenderer) Render() (output []byte) {
	output = r.BaseRenderer.Render()
	return
}

func (r *ProtyleExportDocxRenderer) renderCustomBlock(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderAttributeView(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderTextMark(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.WalkContinue
	}

	types := strings.Split(node.TextMarkType, " ")
	for _, typ := range types {
		r.renderHTMLTag0(node, typ, true)
	}

	r.WriteString(r.getTextMarkTextContent(node))

	reverse(types)
	for _, typ := range types {
		r.renderHTMLTag0(node, typ, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderHTMLTag0(node *ast.Node, currentTextMarkType string, entering bool) {
	switch currentTextMarkType {
	case "a":
		if entering {
			attrs := [][]string{{"href", node.TextMarkAHref}}
			if "" != node.TextMarkATitle {
				attrs = append(attrs, []string{"title", node.TextMarkATitle})
			}
			r.Tag("a", attrs, false)

		} else {
			r.WriteString("</a>")
		}
	case "block-ref":
		if entering {
			node.TextMarkTextContent = strings.ReplaceAll(node.TextMarkTextContent, "'", "&apos;")
			r.WriteString("((" + node.TextMarkBlockRefID)
			if "s" == node.TextMarkBlockRefSubtype {
				r.WriteString(" \"" + node.TextMarkTextContent + "\"")
			} else {
				r.WriteString(" '" + node.TextMarkTextContent + "'")
			}
			r.WriteString("))")
		}
	case "file-annotation-ref":
		if entering {
			node.TextMarkTextContent = strings.ReplaceAll(node.TextMarkTextContent, "'", "&apos;")
			r.WriteString("<<" + node.TextMarkFileAnnotationRefID)
			r.WriteString(" \"" + node.TextMarkTextContent + "\"")
			r.WriteString(">>")
		}
	case "inline-memo":
		if entering {
			r.WriteString(node.TextMarkTextContent)
		}

		if node.IsNextSameInlineMemo() {
			return
		}

		lastRune, _ := utf8.DecodeLastRuneInString(node.TextMarkTextContent)
		if isCJK(lastRune) {
			if entering {
				r.WriteString("<sup>（")
			} else {
				r.WriteString("）</sup>")
			}
		} else {
			if entering {
				r.WriteString("<sup>(")
			} else {
				r.WriteString(")</sup>")
			}
		}
	case "inline-math":
		if entering {
			r.WriteString("<span>$")
		} else {
			r.WriteString("$</span>")
		}
	case "strong":
		if entering {
			r.WriteString("<strong>")
		} else {
			r.WriteString("</strong>")
		}
	case "em":
		if entering {
			r.WriteString("<em>")
		} else {
			r.WriteString("</em>")
		}
	case "code":
		if entering {
			r.WriteString("<code>")
		} else {
			r.WriteString("</code>")
		}
	case "tag":
		if entering {
			r.WriteString("<mark>#")
		} else {
			r.WriteString("#</mark>")
		}
	case "s":
		if entering {
			r.WriteString("<s>")
		} else {
			r.WriteString("</s>")
		}
	case "mark":
		if entering {
			r.WriteString("<mark>")
		} else {
			r.WriteString("</mark>")
		}
	case "u":
		if entering {
			r.WriteString("<u>")
		} else {
			r.WriteString("</u>")
		}
	case "sup":
		if entering {
			r.WriteString("<sup>")
		} else {
			r.WriteString("</sup>")
		}
	case "sub":
		if entering {
			r.WriteString("<sub>")
		} else {
			r.WriteString("</sub>")
		}
	case "kbd":
		if entering {
			r.WriteString("<kbd>")
		} else {
			r.WriteString("</kbd>")
		}
	}
	return
}

func (r *ProtyleExportDocxRenderer) getTextMarkTextContent(node *ast.Node) (ret string) {
	ret = node.TextMarkTextContent
	if node.IsTextMarkType("a") || node.IsTextMarkType("block-ref") || node.IsTextMarkType("file-annotation-ref") {
		content := node.TextMarkTextContent
		content = strings.ReplaceAll(content, editor.IALValEscNewLine, " ")
		ret = content
	} else if node.IsTextMarkType("inline-memo") {
		content := node.TextMarkInlineMemoContent
		content = strings.ReplaceAll(content, editor.IALValEscNewLine, " ")
		ret = content
	} else if node.IsTextMarkType("inline-math") {
		ret = node.TextMarkInlineMathContent
	}

	if node.ParentIs(ast.NodeTableCell) {
		// Improve the handling of inline-math containing `|` in the table https://github.com/siyuan-note/siyuan/issues/9227
		ret = strings.ReplaceAll(ret, "|", "&#124;")
		ret = strings.ReplaceAll(ret, "\n", "<br />")
	}
	return ret
}

func (r *ProtyleExportDocxRenderer) renderBr(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<br />")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderUnderline(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderUnderlineOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<u>")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderUnderlineCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</u>")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderKbd(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderKbdOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<kbd>")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderKbdCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</kbd>")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderVideo(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderAudio(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderIFrame(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderWidget(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderGitConflictCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderGitConflictContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderGitConflictOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderGitConflict(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if entering {
		attrs := [][]string{{"class", "language-git-conflict"}}
		attrs = append(attrs, node.KramdownIAL...)
		r.Tag("div", attrs, false)
	} else {
		r.Tag("/div", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderSuperBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderSuperBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *ProtyleExportDocxRenderer) renderSuperBlockLayoutMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *ProtyleExportDocxRenderer) renderSuperBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *ProtyleExportDocxRenderer) renderLinkRefDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *ProtyleExportDocxRenderer) renderLinkRefDef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *ProtyleExportDocxRenderer) renderTag(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderTagOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", node.Parent.KramdownIAL, false)
		r.WriteByte(lex.ItemCrosshatch)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderTagCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemCrosshatch)
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderKramdownSpanIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderMark1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("mark", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("mark", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/mark", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderSupOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderSupCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sup", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderSubOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("sub", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderSubCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/sub", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("div", nil, false)
	} else {
		r.Tag("/div", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderBlockQueryEmbedScript(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(node.Tokens)
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderBlockRef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderBlockRefID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderBlockRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderBlockRefText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(node.Tokens)
	} else {
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderBlockRefDynamicText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemSinglequote)
		r.Write(node.Tokens)
	} else {
		r.WriteByte(lex.ItemSinglequote)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderFileAnnotationRef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderFileAnnotationRefID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderFileAnnotationRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderFileAnnotationRefText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(node.Tokens)
	} else {
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</code></pre>")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		attrs := [][]string{{"class", "vditor-yml-front-matter"}}
		attrs = append(attrs, node.Parent.KramdownIAL...)
		r.Tag("pre", attrs, false)
		r.WriteString("<code class=\"language-yaml\">")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	return r.BaseRenderer.renderToC(node, entering)
}

func (r *ProtyleExportDocxRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderFootnotesDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<div class=\"footnotes-defs-div\">")
		r.WriteString("<hr class=\"footnotes-defs-hr\" />\n")
		r.WriteString("<ol class=\"footnotes-defs-ol\">")
	} else {
		r.WriteString("</ol></div>")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()

	noHighlight := false
	var language string
	if nil != node.FirstChild.Next && 0 < len(node.FirstChild.Next.CodeBlockInfo) {
		language = util.BytesToStr(node.FirstChild.Next.CodeBlockInfo)
		noHighlight = NoHighlight(language)
	}

	if entering {
		if noHighlight {
			var attrs [][]string
			tokens := html.EscapeHTML(node.FirstChild.Next.Next.Tokens)
			tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
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

func (r *ProtyleExportDocxRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/span", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := html.EscapeHTML(node.Next.Tokens)
		content := util.BytesToStr(tokens)
		r.Tag("span", nil, false)
		r.WriteString("$" + content + "$")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if entering {
		tokens := html.EscapeHTML(node.FirstChild.Next.Tokens)
		tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
		tokens = bytes.TrimSpace(tokens)
		content := util.BytesToStr(tokens)
		r.Tag("div", nil, false)
		r.WriteString("$$\n" + content + "\n$$")
		r.Tag("/div", nil, false)
		r.Newline()
	}
	return ast.WalkSkipChildren
}

func (r *ProtyleExportDocxRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("tr", nil, false)
		r.Newline()
	} else {
		r.Tag("/tr", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
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

func (r *ProtyleExportDocxRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("del", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/del", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("del", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/del", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderCloseBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderOpenBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderLess(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderGreater(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := node.Tokens
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		r.Write(tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		var attrs [][]string
		attrs = append(attrs, node.KramdownIAL...)
		r.Tag("p", attrs, false)
		if r.Options.ChineseParagraphBeginningSpace && ast.NodeDocument == node.Parent.Type {
			if !r.ParagraphContainImgOnly(node) {
				r.WriteString("　　")
			}
		}
	} else {
		r.Tag("/p", nil, false)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("code", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(html.EscapeHTML(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/code", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("em", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/em", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("strong", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("strong", node.Parent.KramdownIAL, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("/strong", nil, false)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("blockquote", node.KramdownIAL, false)
		r.Newline()
	} else {
		r.Newline()
		r.WriteString("</blockquote>")
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		level := headingLevel[node.HeadingLevel : node.HeadingLevel+1]
		r.WriteString("<h" + level)
		for _, attr := range node.KramdownIAL {
			r.WriteString(" " + attr[0] + "=\"" + attr[1] + "\"")
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

func (r *ProtyleExportDocxRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var attrs [][]string
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

func (r *ProtyleExportDocxRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Tag("hr", nil, true)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Tag("br", nil, true)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportDocxRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportDocxRenderer) spanNodeAttrs(node *ast.Node, attrs *[][]string) {
	*attrs = append(*attrs, node.KramdownIAL...)
}
