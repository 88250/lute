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
	"sort"
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

type ProtyleExportMdRenderer struct {
	*BaseRenderer
	NodeWriterStack []*bytes.Buffer
}

func NewProtyleExportMdRenderer(tree *parse.Tree, options *Options) *ProtyleExportMdRenderer {
	ret := &ProtyleExportMdRenderer{BaseRenderer: NewBaseRenderer(tree, options)}
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
	return ret
}

func (r *ProtyleExportMdRenderer) renderTextMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		marker := r.renderMdMarker(node, entering)
		r.WriteString(marker)
		if !node.IsTextMarkType("a") && !node.IsTextMarkType("inline-memo") && !node.IsTextMarkType("block-ref") && !node.IsTextMarkType("file-annotation-ref") {
			textContent := node.TextMarkTextContent
			if node.IsTextMarkType("code") {
				textContent = html.UnescapeString(textContent)
			}
			if strings.HasPrefix(textContent, " ") {
				r.WriteString(editor.Zwsp) // 填充零宽空格以满足 Markdown 语法 https://github.com/siyuan-note/siyuan/issues/6472
			}
			r.WriteString(textContent)
			if strings.HasSuffix(textContent, " ") {
				r.WriteString(editor.Zwsp) // 填充零宽空格以满足 Markdown 语法 https://github.com/siyuan-note/siyuan/issues/6472
			}
		}
	} else {
		marker := r.renderMdMarker(node, entering)
		r.WriteString(marker)
		if nil != node.Next && ast.NodeTextMark == node.Next.Type {
			r.WriteString(editor.Zwsp) // 通过零宽空格来区隔相邻的 Markdown 标记符
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderMdMarker(node *ast.Node, entering bool) (ret string) {
	types := strings.Split(node.TextMarkType, " ")

	// 将这些元素排到最后，避免输出嵌套标记符影响其语义
	// 这样排序后标记符会在这些元素外部而不是嵌套在内部
	sort.Slice(types, func(i, j int) bool {
		if "code" == types[i] || "inline-math" == types[i] || "kbd" == types[i] || "tag" == types[i] {
			return false
		}
		return true
	})

	if !entering {
		reverse(types)
	}

	for _, typ := range types {
		switch typ {
		case "a":
			if entering {
				ret += "[" + node.TextMarkTextContent + "](" + node.TextMarkAHref
				if "" != node.TextMarkATitle {
					ret += " \"" + node.TextMarkATitle + "\""
				}
				ret += ")"
			}
		case "block-ref":
			if entering {
				node.TextMarkTextContent = strings.ReplaceAll(node.TextMarkTextContent, "'", "&apos;")
				ret += "((" + node.TextMarkBlockRefID
				if "s" == node.TextMarkBlockRefSubtype {
					ret += " \"" + node.TextMarkTextContent + "\""
				} else {
					ret += " '" + node.TextMarkTextContent + "'"
				}
				ret += "))"
			}
		case "file-annotation-ref":
			if entering {
				node.TextMarkTextContent = strings.ReplaceAll(node.TextMarkTextContent, "'", "&apos;")
				ret += "<<" + node.TextMarkFileAnnotationRefID
				ret += " \"" + node.TextMarkTextContent + "\""
				ret += ">>"
			}
		case "inline-math":
			if entering {
				ret += "$" + node.TextMarkInlineMathContent
			} else {
				ret += "$"
			}
		case "inline-memo":
			if entering {
				lastRune, _ := utf8.DecodeLastRuneInString(node.TextMarkTextContent)
				ret += node.TextMarkTextContent
				if isCJK(lastRune) {
					ret += "^（" + node.TextMarkInlineMemoContent + "）^"
				} else {
					ret += "^(" + node.TextMarkInlineMemoContent + ")^"
				}
			}
		case "strong":
			ret += "**"
		case "em":
			ret += "*"
		case "code":
			ret += "`"
		case "tag":
			ret += "#"
		case "s":
			ret += "~~"
		case "mark":
			ret += "=="
		case "u":
			if entering {
				ret += "<u>"
			} else {
				ret += "</u>"
			}
		case "sup":
			if entering {
				ret += "^"
			} else {
				ret += "^"
			}
		case "sub":
			if entering {
				ret += "~"
			} else {
				ret += "~"
			}
		case "kbd":
			if entering {
				ret += "<kbd>"
			} else {
				ret += "</kbd>"
			}
		}
	}
	return
}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

func (r *ProtyleExportMdRenderer) renderBr(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<br />")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderUnderline(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderUnderlineOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<u>")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderUnderlineCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</u>")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderKbd(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderKbdOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("<kbd>")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderKbdCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("</kbd>")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderVideo(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		tokens := node.Tokens
		tokens = r.tagSrcPath(tokens)
		r.Write(tokens)
		r.Newline()
		if !r.isLastNode(r.Tree.Root, node) {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderAudio(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		tokens := node.Tokens
		tokens = r.tagSrcPath(tokens)
		r.Write(tokens)
		r.Newline()
		if !r.isLastNode(r.Tree.Root, node) {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderIFrame(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		tokens := node.Tokens
		tokens = r.tagSrcPath(tokens)
		r.Write(tokens)
		r.Newline()
		if !r.isLastNode(r.Tree.Root, node) {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderWidget(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		tokens := node.Tokens
		tokens = r.tagSrcPath(tokens)
		r.Write(tokens)
		r.Newline()
		if !r.isLastNode(r.Tree.Root, node) {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderGitConflictCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderGitConflictContent(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		r.Write(node.Tokens)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderGitConflictOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderGitConflict(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderSuperBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderSuperBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering && r.Options.SuperBlock {
		r.Write([]byte("{{{"))
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderSuperBlockLayoutMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering && r.Options.SuperBlock {
		r.Write(node.Tokens)
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderSuperBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if r.Options.SuperBlock {
			r.Newline()
			r.Write([]byte("}}}"))
			r.Newline()
		}
		if !r.isLastNode(r.Tree.Root, node) {
			if r.withoutKramdownBlockIAL(node.Parent) {
				r.WriteByte(lex.ItemNewline)
			}
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderLinkRefDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderLinkRefDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemOpenBracket)
		r.Write(node.Tokens)
		r.WriteString("]: ")
	} else {
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderTag(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderTagOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemCrosshatch)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderTagCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemCrosshatch)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	if !r.Options.KramdownBlockIAL {
		return ast.WalkContinue
	}

	if nil != node.Previous && ast.NodeListItem == node.Previous.Type {
		return ast.WalkContinue
	}
	if entering {
		r.Newline()
		if r.Options.KramdownBlockIAL {
			if util.IsDocIAL(node.Tokens) {
				r.WriteByte(lex.ItemNewline)
			}
			r.Write(node.Tokens)
		}
	} else {
		if ast.NodeListItem == node.Parent.Type || ast.NodeList == node.Parent.Type {
			if !node.Parent.ListData.Tight {
				r.Newline()
			}
		} else {
			r.Newline()
		}
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderKramdownSpanIAL(node *ast.Node, entering bool) ast.WalkStatus {
	if !r.Options.KramdownSpanIAL {
		return ast.WalkContinue
	}

	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderMark1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("=")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("=")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("==")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("==")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderSupOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("^")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderSupCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("^")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderSubOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("~")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderSubCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("~")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderBlockQueryEmbedScript(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
	} else {
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderBlockRef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderBlockRefID(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderBlockRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemSpace)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderBlockRefText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		tokens := html.EscapeHTML(node.Tokens)
		tokens = bytes.ReplaceAll(tokens, []byte("'"), []byte("&apos;"))
		r.Write(tokens)
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderBlockRefDynamicText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemSinglequote)
		tokens := html.EscapeHTML(node.Tokens)
		tokens = bytes.ReplaceAll(tokens, []byte("'"), []byte("&apos;"))
		r.Write(tokens)
		r.WriteByte(lex.ItemSinglequote)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderFileAnnotationRef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderFileAnnotationRefID(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderFileAnnotationRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemSpace)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderFileAnnotationRefText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		tokens := html.EscapeHTML(node.Tokens)
		tokens = bytes.ReplaceAll(tokens, []byte("'"), []byte("&apos;"))
		r.Write(tokens)
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(parse.YamlFrontMatterMarker)
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(parse.YamlFrontMatterMarker)
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		if !entering && !r.isLastNode(r.Tree.Root, node) {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.HtmlEntityTokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemBackslash)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("[toc]\n\n")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("[" + util.BytesToStr(node.Tokens) + "]")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderFootnotesDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.NodeWriterStack = append(r.NodeWriterStack, r.Writer)
		r.WriteString("[" + util.BytesToStr(node.Tokens) + "]: ")
	} else {
		writer := r.NodeWriterStack[len(r.NodeWriterStack)-1]
		r.NodeWriterStack = r.NodeWriterStack[:len(r.NodeWriterStack)-1]
		buf := writer.String()
		lines := strings.Split(buf, "\n")
		contentBuf := bytes.Buffer{}
		for i, line := range lines {
			if 0 == i {
				contentBuf.WriteString(line + "\n")
			} else {
				if "" == line {
					contentBuf.WriteString("\n")
				} else {
					contentBuf.WriteString("    " + line + "\n")
				}
			}
		}
		r.NodeWriterStack[len(r.NodeWriterStack)-1].Write(contentBuf.Bytes())
		r.Writer = r.NodeWriterStack[len(r.NodeWriterStack)-1]
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
	padding := node.TableCellContentMaxWidth - node.TableCellContentWidth
	if entering {
		r.WriteByte(lex.ItemPipe)
		if !r.Options.ProtyleWYSIWYG {
			r.WriteByte(lex.ItemSpace)
			switch node.TableCellAlign {
			case 2:
				r.Write(bytes.Repeat([]byte{lex.ItemSpace}, padding/2))
			case 3:
				r.Write(bytes.Repeat([]byte{lex.ItemSpace}, padding))
			}
		}
	} else {
		if !r.Options.ProtyleWYSIWYG {
			switch node.TableCellAlign {
			case 2:
				r.Write(bytes.Repeat([]byte{lex.ItemSpace}, padding/2))
			case 3:
			default:
				r.Write(bytes.Repeat([]byte{lex.ItemSpace}, padding))
			}
			r.WriteByte(lex.ItemSpace)
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		r.WriteString("|\n")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		headRow := node.FirstChild
		for th := headRow.FirstChild; nil != th; th = th.Next {
			if ast.NodeKramdownSpanIAL == th.Type {
				continue
			}

			align := th.TableCellAlign
			switch align {
			case 0:
				r.WriteString("| -")
				if padding := th.TableCellContentMaxWidth - 1; 0 < padding {
					r.Write(bytes.Repeat([]byte{lex.ItemHyphen}, padding))
				}
				if !r.Options.ProtyleWYSIWYG {
					r.WriteByte(lex.ItemSpace)
				}
			case 1:
				r.WriteString("| :-")
				if padding := th.TableCellContentMaxWidth - 2; 0 < padding {
					r.Write(bytes.Repeat([]byte{lex.ItemHyphen}, padding))
				}
				if !r.Options.ProtyleWYSIWYG {
					r.WriteByte(lex.ItemSpace)
				}
			case 2:
				r.WriteString("| :-")
				if padding := th.TableCellContentMaxWidth - 3; 0 < padding {
					r.Write(bytes.Repeat([]byte{lex.ItemHyphen}, padding))
				}
				r.WriteString(": ")
			case 3:
				r.WriteString("| -")
				if padding := th.TableCellContentMaxWidth - 2; 0 < padding {
					r.Write(bytes.Repeat([]byte{lex.ItemHyphen}, padding))
				}
				r.WriteString(": ")
			}
		}
		r.WriteString("|\n")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		// 遍历单元格算出最大宽度

		var cells [][]*ast.Node
		cells = append(cells, []*ast.Node{})

		headRow := node.ChildByType(ast.NodeTableHead)
		if nil == headRow || nil == headRow.FirstChild || nil == node.FirstChild {
			return ast.WalkSkipChildren
		}

		for n := headRow.FirstChild.FirstChild; nil != n; n = n.Next {
			cells[0] = append(cells[0], n)
		}

		i := 1
		for tableRow := node.FirstChild.Next; nil != tableRow; tableRow = tableRow.Next {
			cells = append(cells, []*ast.Node{})
			for n := tableRow.FirstChild; nil != n; n = n.Next {
				cells[i] = append(cells[i], n)
			}
			i++
		}

		var maxWidth int
		for col := 0; col < len(cells[0]); col++ {
			for row := 0; row < len(cells) && col < len(cells[row]); row++ {
				cells[row][col].TableCellContentWidth = cells[row][col].TokenLen()
				// 自动添加空格会导致单元格宽度发生变化
				if r.Options.AutoSpace {
					ret := 0
					// 遍历字节点，将可能会多出来的空格计算出来
					ast.Walk(cells[row][col], func(n *ast.Node, entering bool) ast.WalkStatus {
						if !entering {
							return ast.WalkContinue
						}
						// 空格仅一个字节，可以直接计算长度
						ret += len(r.Space(n.Tokens)) - len(n.Tokens)
						return ast.WalkContinue
					})
					cells[row][col].TableCellContentWidth += ret
				}
				if maxWidth < cells[row][col].TableCellContentWidth {
					maxWidth = cells[row][col].TableCellContentWidth
				}
			}
			for row := 0; row < len(cells) && col < len(cells[row]); row++ {
				cells[row][col].TableCellContentMaxWidth = maxWidth
			}
			maxWidth = 0
		}
	} else {
		r.Newline()
		if !r.isLastNode(r.Tree.Root, node) {
			if r.withoutKramdownBlockIAL(node) {
				r.WriteByte(lex.ItemNewline)
			}
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemTilde)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemTilde)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("~~")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("~~")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(html.EscapeHTML(node.Tokens))
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := node.Tokens
		tokens = r.LinkPath(tokens)
		r.Write(tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemSpace)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var tokens []byte
		if r.Options.AutoSpace {
			tokens = r.Space(node.Tokens)
		} else {
			tokens = node.Tokens
		}
		r.Write(tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemCloseParen)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemOpenParen)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderGreater(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemGreater)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderLess(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemLess)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderCloseBrace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemCloseBrace)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderOpenBrace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemOpenBrace)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemCloseBracket)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemOpenBracket)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemBang)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.LinkTextAutoSpacePrevious(node)
		if 3 == node.LinkType {
			text := node.ChildByType(ast.NodeLinkText).Tokens
			if bytes.Equal(text, node.LinkRefLabel) {
				r.WriteString("[" + util.BytesToStr(text) + "]")
			} else {
				r.WriteString("[" + util.BytesToStr(text) + "][" + util.BytesToStr(node.LinkRefLabel) + "]")
			}
			return ast.WalkSkipChildren
		}
		if 1 == node.LinkType {
			dest := node.ChildByType(ast.NodeLinkDest).Tokens
			r.Write(dest)
			return ast.WalkSkipChildren
		}
	} else {
		r.LinkTextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		tokens := node.Tokens
		tokens = r.tagSrcPath(tokens)
		r.Write(tokens)
		r.Newline()
		if !r.isLastNode(r.Tree.Root, node) {
			if r.withoutKramdownBlockIAL(node) {
				r.WriteByte(lex.ItemNewline)
			}
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.NodeWriterStack = append(r.NodeWriterStack, r.Writer)
	} else {
		r.NodeWriterStack = r.NodeWriterStack[:len(r.NodeWriterStack)-1]
		var buf []byte
		if r.Options.KeepParagraphBeginningSpace {
			buf = bytes.TrimRight(r.Writer.Bytes(), " \t\n")
			buf = bytes.TrimLeft(buf, "\n")
		} else {
			buf = bytes.Trim(r.Writer.Bytes(), " \t\n")
		}
		r.Writer.Reset()
		r.Write(buf)
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if r.Options.ChineseParagraphBeginningSpace && ast.NodeDocument == node.Parent.Type {
			r.WriteString("　　")
		}
	} else {
		if !r.Options.KeepParagraphBeginningSpace && nil != node.FirstChild {
			node.FirstChild.Tokens = bytes.TrimSpace(node.FirstChild.Tokens)
		}

		if node.ParentIs(ast.NodeTableCell) {
			if nil != node.Next && ast.NodeText != node.Next.Type {
				r.WriteString("<br /><br />")
			}
			return ast.WalkContinue
		}

		if r.withoutKramdownBlockIAL(node) {
			r.Newline()
		}

		inTightList := false
		lastListItemLastPara := false
		if parent := node.Parent; nil != parent {
			if ast.NodeListItem == parent.Type { // ListItem.Paragraph
				listItem := parent
				if nil != listItem.Parent && nil != listItem.Parent.ListData {
					// 必须通过列表（而非列表项）上的紧凑标识判断，因为在设置该标识时仅设置了 List.Tight
					// 设置紧凑标识的具体实现可参考函数 List.Finalize()
					inTightList = listItem.Parent.ListData.Tight

					if nextItem := listItem.Next; nil == nextItem {
						nextPara := node.Next
						lastListItemLastPara = nil == nextPara
					}
				} else {
					inTightList = true
				}
			}
		}

		if !inTightList || (lastListItemLastPara) {
			if r.withoutKramdownBlockIAL(node) {
				r.WriteByte(lex.ItemNewline)
			}
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
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
		if (nil == node.Previous || ast.NodeTaskListItemMarker == node.Previous.Type) &&
			nil != node.Parent.Parent && nil != node.Parent.Parent.ListData && 3 == node.Parent.Parent.ListData.Typ {
			if ' ' == r.LastOut {
				tokens = bytes.TrimPrefix(tokens, []byte(" "))
				if bytes.HasPrefix(tokens, []byte(editor.Caret+" ")) {
					tokens = bytes.TrimPrefix(tokens, []byte(editor.Caret+" "))
					tokens = append(editor.CaretTokens, tokens...)
				}
			}
		}
		r.Write(tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *ProtyleExportMdRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemBacktick)
		if 1 < node.Parent.CodeMarkerLen {
			r.WriteByte(lex.ItemBacktick)
			text := util.BytesToStr(node.Next.Tokens)
			firstc, _ := utf8.DecodeRuneInString(text)
			if '`' == firstc {
				r.WriteByte(lex.ItemSpace)
			}
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := node.Tokens
		if node.ParentIs(ast.NodeTableCell) {
			tokens = bytes.ReplaceAll(tokens, []byte("\\|"), []byte("|"))
			tokens = bytes.ReplaceAll(tokens, []byte("|"), []byte("\\|"))
			tokens = bytes.ReplaceAll(tokens, []byte("<br/>"), nil)
		}
		r.Write(tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if 1 < node.Parent.CodeMarkerLen {
			text := util.BytesToStr(node.Previous.Tokens)
			lastc, _ := utf8.DecodeLastRuneInString(text)
			if '`' == lastc {
				r.WriteByte(lex.ItemSpace)
			}
			r.WriteByte(lex.ItemBacktick)
		}
		r.WriteByte(lex.ItemBacktick)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
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
func (r *ProtyleExportMdRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDollar)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := html.UnescapeHTML(node.Tokens)
		r.Write(tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDollar)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(parse.MathBlockMarker)
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := html.UnescapeHTML(node.Tokens)
		r.Write(tokens)
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(parse.MathBlockMarker)
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if !entering && !r.isLastNode(r.Tree.Root, node) {
		if r.withoutKramdownBlockIAL(node) {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		r.Write(node.Tokens)
		r.Newline()
		if !r.isLastNode(r.Tree.Root, node) {
			if r.withoutKramdownBlockIAL(node.Parent) {
				r.WriteByte(lex.ItemNewline)
			}
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		tokens := node.Tokens
		info := node.Parent.ChildByType(ast.NodeCodeBlockFenceInfoMarker)
		if nil != info &&
			bytes.Contains(info.CodeBlockInfo, []byte("flowchart")) ||
			bytes.Contains(info.CodeBlockInfo, []byte("mermaid")) ||
			bytes.Contains(info.CodeBlockInfo, []byte("graphviz")) ||
			bytes.Contains(info.CodeBlockInfo, []byte("plantuml")) {
			tokens = html.UnescapeHTML(tokens)
		}
		r.Write(tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.CodeBlockInfo)
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
		if !node.IsFencedCodeBlock {
			r.Write(bytes.Repeat([]byte{lex.ItemBacktick}, 3))
			r.WriteByte(lex.ItemNewline)
			r.Write(node.FirstChild.Tokens)
			r.Write(bytes.Repeat([]byte{lex.ItemBacktick}, 3))
			r.Newline()
			if !r.isLastNode(r.Tree.Root, node) {
				if r.withoutKramdownBlockIAL(node) {
					r.WriteByte(lex.ItemNewline)
				}
			}
			return ast.WalkSkipChildren
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemAsterisk)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemAsterisk)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemUnderscore)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemUnderscore)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("**")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("**")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("__")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("__")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.NodeWriterStack = append(r.NodeWriterStack, r.Writer)
	} else {
		writer := r.NodeWriterStack[len(r.NodeWriterStack)-1]
		r.NodeWriterStack = r.NodeWriterStack[:len(r.NodeWriterStack)-1]

		blockquoteLines := bytes.Buffer{}
		buf := writer.Bytes()
		lines := bytes.Split(buf, []byte{lex.ItemNewline})
		length := len(lines)
		if 2 < length && lex.IsBlank(lines[length-1]) && lex.IsBlank(lines[length-2]) {
			lines = lines[:length-1]
		}
		if 1 == len(r.NodeWriterStack) { // 已经是根这一层
			length = len(lines)
			if 1 < length && lex.IsBlank(lines[length-1]) {
				lines = lines[:length-1]
			}
		}

		length = len(lines)
		for _, line := range lines {
			if 0 == len(line) {
				blockquoteLines.WriteString(">\n")
				continue
			}

			if lex.ItemGreater == line[0] {
				blockquoteLines.WriteString(">")
			} else {
				blockquoteLines.WriteString("> ")
			}
			blockquoteLines.Write(line)
			blockquoteLines.WriteByte(lex.ItemNewline)
		}
		buf = bytes.TrimSpace(blockquoteLines.Bytes())
		writer.Reset()
		writer.Write(buf)
		r.NodeWriterStack[len(r.NodeWriterStack)-1].Write(writer.Bytes())
		r.Writer = r.NodeWriterStack[len(r.NodeWriterStack)-1]
		buf = bytes.TrimSpace(r.Writer.Bytes())
		r.Writer.Reset()
		r.Write(buf)
		if !node.ParentIs(ast.NodeTableCell) { // 在表格中不能换行，否则会破坏表格的排版 https://github.com/Vanessa219/vditor/issues/368
			if r.withoutKramdownBlockIAL(node) {
				r.WriteString("\n\n")
			}
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if !node.HeadingSetext {
			r.Write(bytes.Repeat([]byte{lex.ItemCrosshatch}, node.HeadingLevel))
			r.WriteByte(lex.ItemSpace)
		}
	} else {
		if node.HeadingSetext {
			r.WriteByte(lex.ItemNewline)
			contentLen := r.setextHeadingLen(node)
			if 1 == node.HeadingLevel {
				r.WriteString(strings.Repeat("=", contentLen))
			} else if 2 == node.HeadingLevel {
				r.WriteString(strings.Repeat("-", contentLen))
			}
		}

		if !node.ParentIs(ast.NodeTableCell) {
			if r.withoutKramdownBlockIAL(node) {
				r.Newline()
				r.WriteByte(lex.ItemNewline)
			}
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString(" {" + util.BytesToStr(node.Tokens) + "}")
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.NodeWriterStack = append(r.NodeWriterStack, r.Writer)
	} else {
		writer := r.NodeWriterStack[len(r.NodeWriterStack)-1]
		r.NodeWriterStack = r.NodeWriterStack[:len(r.NodeWriterStack)-1]
		r.NodeWriterStack[len(r.NodeWriterStack)-1].Write(writer.Bytes())
		r.Writer = r.NodeWriterStack[len(r.NodeWriterStack)-1]
		buf := bytes.TrimSpace(r.Writer.Bytes())
		r.Writer.Reset()
		r.Write(buf)
		if !node.ParentIs(ast.NodeTableCell) {
			if r.withoutKramdownBlockIAL(node) {
				r.WriteString("\n\n")
			}
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.NodeWriterStack = append(r.NodeWriterStack, r.Writer)
		if r.Options.KramdownBlockIAL && nil != node.Next && ast.NodeKramdownBlockIAL == node.Next.Type {
			liIAL := node.Next
			r.Write(liIAL.Tokens)
		}
		if nil != node.FirstChild && ast.NodeList == node.FirstChild.Type {
			r.Newline()
		}
	} else {
		writer := r.NodeWriterStack[len(r.NodeWriterStack)-1]
		r.NodeWriterStack = r.NodeWriterStack[:len(r.NodeWriterStack)-1]
		indent := len(node.ListData.Marker) + 1
		if 1 == node.ListData.Typ || (3 == node.ListData.Typ && 0 == node.ListData.BulletChar) {
			indent++
		}
		indentSpaces := bytes.Repeat([]byte{lex.ItemSpace}, indent)
		indentedLines := bytes.Buffer{}
		buf := writer.Bytes()
		if bytes.HasPrefix(buf, []byte("* ")) {
			// 说明该列表项为空 https://github.com/siyuan-note/siyuan/issues/6206
			buf = append([]byte(" \n\n"), buf...)
		}
		lines := bytes.Split(buf, []byte{lex.ItemNewline})
		for _, line := range lines {
			if 0 == len(line) {
				indentedLines.WriteByte(lex.ItemNewline)
				continue
			}
			indentedLines.Write(indentSpaces)
			indentedLines.Write(line)
			indentedLines.WriteByte(lex.ItemNewline)
		}
		buf = indentedLines.Bytes()
		if indent < len(buf) {
			buf = buf[indent:]
		}

		listItemBuf := bytes.Buffer{}
		if 1 == node.ListData.Typ || (3 == node.ListData.Typ && 0 == node.ListData.BulletChar) {
			listItemBuf.WriteString(strconv.Itoa(node.ListData.Num) + string(node.ListData.Delimiter))
		} else {
			listItemBuf.Write(node.ListData.Marker)
		}
		listItemBuf.WriteByte(lex.ItemSpace)
		buf = append(listItemBuf.Bytes(), buf...)
		if node.ParentIs(ast.NodeTableCell) {
			buf = bytes.ReplaceAll(buf, []byte("\n"), nil)
		}
		writer.Reset()
		writer.Write(buf)
		buf = writer.Bytes()
		if node.ParentIs(ast.NodeTableCell) {
			buf = bytes.ReplaceAll(buf, []byte("\n"), nil)
		}
		r.NodeWriterStack[len(r.NodeWriterStack)-1].Write(buf)
		r.Writer = r.NodeWriterStack[len(r.NodeWriterStack)-1]
		buf = bytes.TrimSpace(r.Writer.Bytes())
		r.Writer.Reset()
		r.Write(buf)
		if !node.ParentIs(ast.NodeTableCell) {
			r.WriteString("\n")
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemOpenBracket)
		if node.TaskListItemChecked {
			r.WriteByte('X')
		} else {
			r.WriteByte(lex.ItemSpace)
		}
		r.WriteByte(lex.ItemCloseBracket)
	} else {
		r.WriteByte(' ')
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if node.ParentIs(ast.NodeTableCell) {
			r.WriteString("<hr/>")
		} else {
			r.WriteString("---")
			if r.withoutKramdownBlockIAL(node) {
				r.WriteByte(lex.ItemNewline)
				r.WriteByte(lex.ItemNewline)
			}
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if !r.Options.SoftBreak2HardBreak {
			r.WriteString("\\\n")
		} else {
			if node.ParentIs(ast.NodeTableCell) {
				r.WriteString("<br/>")
			} else {
				r.WriteByte(lex.ItemNewline)
			}
		}
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *ProtyleExportMdRenderer) withoutKramdownBlockIAL(node *ast.Node) bool {
	return !r.Options.KramdownBlockIAL || 0 == len(node.KramdownIAL) || nil == node.Next || ast.NodeKramdownBlockIAL != node.Next.Type
}
