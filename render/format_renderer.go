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
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// FormatRenderer 描述了格式化渲染器。
type FormatRenderer struct {
	*BaseRenderer
	NodeWriterStack []*bytes.Buffer // 节点输出缓冲栈
}

// NewFormatRenderer 创建一个格式化渲染器。
func NewFormatRenderer(tree *parse.Tree) *FormatRenderer {
	ret := &FormatRenderer{BaseRenderer: NewBaseRenderer(tree)}
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

func (r *FormatRenderer) renderTag(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderTagOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemCrosshatch)
	return ast.WalkStop
}

func (r *FormatRenderer) renderTagCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemCrosshatch)
	return ast.WalkStop
}

func (r *FormatRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	if !r.Option.KramdownIAL {
		return ast.WalkContinue
	}

	if entering {
		r.Newline()
		r.Write(node.Tokens)
	} else {
		if ast.NodeListItem == node.Parent.Type || ast.NodeList == node.Parent.Type {
			if !node.Parent.Tight {
				r.Newline()
			}
		} else {
			r.Newline()
		}
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) Render() (output []byte) {
	output = r.BaseRenderer.Render()
	if 1 > len(r.Tree.Context.LinkRefDefs) {
		return
	}

	buf := &bytes.Buffer{}
	buf.WriteByte(lex.ItemNewline)
	// 将链接引用定义添加到末尾
	for _, node := range r.Tree.Context.LinkRefDefs {
		label := node.LinkRefLabel
		dest := node.ChildByType(ast.NodeLinkDest).Tokens
		buf.WriteString("[" + util.BytesToStr(label) + "]: " + util.BytesToStr(dest) + "\n")
	}
	output = append(output, buf.Bytes()...)
	return
}

func (r *FormatRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderMark1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("=")
	return ast.WalkStop
}

func (r *FormatRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("=")
	return ast.WalkStop
}

func (r *FormatRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("==")
	return ast.WalkStop
}

func (r *FormatRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("==")
	return ast.WalkStop
}

func (r *FormatRenderer) renderBlockQueryEmbedScript(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
	} else {
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderBlockEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if nil != node.Previous && ast.NodeTaskListItemMarker != node.Previous.Type {
			r.Newline()
		}
	} else {
		r.Newline()
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderBlockEmbedID(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderBlockEmbedSpace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemSpace)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderBlockEmbedText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(node.Tokens)
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkStop
}

func (r *FormatRenderer) renderBlockRef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *FormatRenderer) renderBlockRefID(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Write(node.Tokens)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderBlockRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemSpace)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderBlockRefText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemDoublequote)
		r.Write(node.Tokens)
		r.WriteByte(lex.ItemDoublequote)
	}
	return ast.WalkStop
}

func (r *FormatRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(parse.YamlFrontMatterMarker)
	r.WriteByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *FormatRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	r.WriteByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *FormatRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(parse.YamlFrontMatterMarker)
	r.WriteByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *FormatRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if !entering && !r.isLastNode(r.Tree.Root, node) {
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.HtmlEntityTokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemBackslash)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("[toc]\n\n")
	return ast.WalkStop
}

func (r *FormatRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("[" + util.BytesToStr(node.Tokens) + "]")
	return ast.WalkStop
}

func (r *FormatRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("[" + util.BytesToStr(node.Tokens) + "]: ")
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *FormatRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *FormatRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *FormatRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
	padding := node.TableCellContentMaxWidth - node.TableCellContentWidth
	if entering {
		r.WriteByte(lex.ItemPipe)
		r.WriteByte(lex.ItemSpace)
		switch node.TableCellAlign {
		case 2:
			r.Write(bytes.Repeat([]byte{lex.ItemSpace}, padding/2))
		case 3:
			r.Write(bytes.Repeat([]byte{lex.ItemSpace}, padding))
		}
	} else {
		switch node.TableCellAlign {
		case 2:
			r.Write(bytes.Repeat([]byte{lex.ItemSpace}, padding/2))
		case 3:
		default:
			r.Write(bytes.Repeat([]byte{lex.ItemSpace}, padding))
		}
		r.WriteByte(lex.ItemSpace)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		r.WriteString("|\n")
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		headRow := node.FirstChild
		for th := headRow.FirstChild; nil != th; th = th.Next {
			align := th.TableCellAlign
			switch align {
			case 0:
				r.WriteString("| -")
				if padding := th.TableCellContentMaxWidth - 1; 0 < padding {
					r.Write(bytes.Repeat([]byte{lex.ItemHyphen}, padding))
				}
				r.WriteByte(lex.ItemSpace)
			case 1:
				r.WriteString("| :-")
				if padding := th.TableCellContentMaxWidth - 2; 0 < padding {
					r.Write(bytes.Repeat([]byte{lex.ItemHyphen}, padding))
				}
				r.WriteByte(lex.ItemSpace)
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

func (r *FormatRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		r.Newline()
		if !r.isLastNode(r.Tree.Root, node) {
			if r.withoutKramdownIAL(node) {
				r.WriteByte(lex.ItemNewline)
			}
		}
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemTilde)
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemTilde)
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("~~")
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("~~")
	return ast.WalkStop
}

func (r *FormatRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemDoublequote)
	r.Write(node.Tokens)
	r.WriteByte(lex.ItemDoublequote)
	return ast.WalkStop
}

func (r *FormatRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemSpace)
	return ast.WalkStop
}

func (r *FormatRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if r.Option.AutoSpace {
		r.Space(node)
	}
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemCloseParen)
	return ast.WalkStop
}

func (r *FormatRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemOpenParen)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemCloseBracket)
	return ast.WalkStop
}

func (r *FormatRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemOpenBracket)
	return ast.WalkStop
}

func (r *FormatRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemBang)
	return ast.WalkStop
}

func (r *FormatRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *FormatRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.LinkTextAutoSpacePrevious(node)
		if 3 == node.LinkType {
			text := node.ChildByType(ast.NodeLinkText).Tokens
			if bytes.Equal(text, node.LinkRefLabel) {
				r.WriteString("[" + util.BytesToStr(text) + "]")
			} else {
				r.WriteString("[" + util.BytesToStr(text) + "][" + util.BytesToStr(node.LinkRefLabel) + "]")
			}
			return ast.WalkStop
		}
	} else {
		r.LinkTextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	r.Write(node.Tokens)
	r.Newline()
	if !r.isLastNode(r.Tree.Root, node) {
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkStop
}

func (r *FormatRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.NodeWriterStack = append(r.NodeWriterStack, r.Writer)
	} else {
		r.NodeWriterStack = r.NodeWriterStack[:len(r.NodeWriterStack)-1]
		buf := bytes.Trim(r.Writer.Bytes(), " \t\n")
		r.Writer.Reset()
		r.Write(buf)
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if r.Option.KramdownIAL {
			parent := node.Parent
			if ast.NodeListItem == parent.Type && parent.FirstChild == node { // 列表项下第一个段落
				if nil != parent.Next && ast.NodeKramdownBlockIAL == parent.Next.Type {
					liIAL := parent.Next
					r.Write(liIAL.Tokens)
					liIAL.Unlink()
				}
			}
		}
	} else {
		if !node.ParentIs(ast.NodeTableCell) {
			if r.withoutKramdownIAL(node) {
				r.Newline()
			}
		}
		inTightList := false
		lastListItemLastPara := false
		if parent := node.Parent; nil != parent {
			if ast.NodeListItem == parent.Type { // ListItem.Paragraph
				listItem := parent
				if nil != listItem.Parent && nil != listItem.Parent.ListData {
					// 必须通过列表（而非列表项）上的紧凑标识判断，因为在设置该标识时仅设置了 List.Tight
					// 设置紧凑标识的具体实现可参考函数 List.Finalize()
					inTightList = listItem.Parent.Tight

					if nextItem := listItem.Next; nil == nextItem {
						nextPara := node.Next
						lastListItemLastPara = nil == nextPara
					}
				} else {
					inTightList = true
				}
			}
		}

		if (!inTightList || (lastListItemLastPara)) && !node.ParentIs(ast.NodeTableCell) {
			if r.withoutKramdownIAL(node) {
				r.WriteByte(lex.ItemNewline)
			}
		}
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if r.Option.AutoSpace {
		r.Space(node)
	}
	if r.Option.FixTermTypo {
		r.FixTermTypo(node)
	}
	if r.Option.ChinesePunct {
		r.ChinesePunct(node)
	}
	if nil == node.Previous && nil != node.Parent.Parent && nil != node.Parent.Parent.ListData && 3 == node.Parent.Parent.ListData.Typ {
		// 任务列表起始位置使用 `<font>` 标签的预览问题 https://github.com/siyuan-note/siyuan/issues/33
		if !bytes.HasPrefix(node.Tokens, []byte(" ")) && ' ' != r.LastOut {
			node.Tokens = append([]byte(" "), node.Tokens...)
		}
	}
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *FormatRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if node.ParentIs(ast.NodeTableCell) && (bytes.Contains(node.Next.Tokens, []byte("|")) || bytes.Contains(node.Next.Tokens, []byte("`"))) {
		r.WriteString("<code>")
		return ast.WalkStop
	}

	r.WriteByte(lex.ItemBacktick)
	if 1 < node.Parent.CodeMarkerLen {
		r.WriteByte(lex.ItemBacktick)
		text := util.BytesToStr(node.Next.Tokens)
		firstc, _ := utf8.DecodeRuneInString(text)
		if '`' == firstc {
			r.WriteByte(lex.ItemSpace)
		}
	}
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	tokens := node.Tokens
	if node.ParentIs(ast.NodeTableCell) {
		tokens = bytes.ReplaceAll(tokens, []byte("\\|"), []byte("|"))
		tokens = bytes.ReplaceAll(tokens, []byte("|"), []byte("\\|"))
		tokens = bytes.ReplaceAll(tokens, []byte("<br/>"), nil)
	}
	r.Write(tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if node.ParentIs(ast.NodeTableCell) && (bytes.Contains(node.Previous.Tokens, []byte("|")) || bytes.Contains(node.Previous.Tokens, []byte("`"))) {
		r.WriteString("</code>")
		return ast.WalkStop
	}

	if 1 < node.Parent.CodeMarkerLen {
		text := util.BytesToStr(node.Previous.Tokens)
		lastc, _ := utf8.DecodeLastRuneInString(text)
		if '`' == lastc {
			r.WriteByte(lex.ItemSpace)
		}
		r.WriteByte(lex.ItemBacktick)
	}
	r.WriteByte(lex.ItemBacktick)
	return ast.WalkStop
}

func (r *FormatRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}
func (r *FormatRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemDollar)
	return ast.WalkStop
}

func (r *FormatRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemDollar)
	return ast.WalkStop
}

func (r *FormatRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(parse.MathBlockMarker)
	r.WriteByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *FormatRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	r.WriteByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *FormatRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(parse.MathBlockMarker)
	r.WriteByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *FormatRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if !entering && !r.isLastNode(r.Tree.Root, node) {
		if r.withoutKramdownIAL(node) {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	r.Write(node.Tokens)
	r.Newline()
	if !r.isLastNode(r.Tree.Root, node) {
		if r.withoutKramdownIAL(node.Parent) {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.CodeBlockInfo)
	r.WriteByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Newline()
	}
	if !node.IsFencedCodeBlock {
		r.Write(bytes.Repeat([]byte{lex.ItemBacktick}, 3))
		r.WriteByte(lex.ItemNewline)
		r.Write(node.FirstChild.Tokens)
		r.Write(bytes.Repeat([]byte{lex.ItemBacktick}, 3))
		r.Newline()
		if !r.isLastNode(r.Tree.Root, node) {
			if r.withoutKramdownIAL(node) {
				r.WriteByte(lex.ItemNewline)
			}
		}
		return ast.WalkStop
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemAsterisk)
	return ast.WalkStop
}

func (r *FormatRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemAsterisk)
	return ast.WalkStop
}

func (r *FormatRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemUnderscore)
	return ast.WalkStop
}

func (r *FormatRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemUnderscore)
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("**")
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("**")
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("__")
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("__")
	return ast.WalkStop
}

func (r *FormatRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
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
		buf = blockquoteLines.Bytes()
		writer.Reset()
		writer.Write(buf)
		r.NodeWriterStack[len(r.NodeWriterStack)-1].Write(writer.Bytes())
		r.Writer = r.NodeWriterStack[len(r.NodeWriterStack)-1]
		buf = bytes.TrimSpace(r.Writer.Bytes())
		r.Writer.Reset()
		r.Write(buf)
		if !node.ParentIs(ast.NodeTableCell) { // 在表格中不能换行，否则会破坏表格的排版 https://github.com/Vanessa219/vditor/issues/368
			if r.withoutKramdownIAL(node) {
				r.WriteString("\n\n")
			}
		}
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *FormatRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
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
			if r.withoutKramdownIAL(node) {
				r.Newline()
				r.WriteByte(lex.ItemNewline)
			}
		}
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *FormatRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString(" {" + util.BytesToStr(node.Tokens) + "}")
	return ast.WalkStop
}

func (r *FormatRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
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
			if r.withoutKramdownIAL(node) {
				r.WriteString("\n\n")
			}
		}
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.NodeWriterStack = append(r.NodeWriterStack, r.Writer)
	} else {
		writer := r.NodeWriterStack[len(r.NodeWriterStack)-1]
		r.NodeWriterStack = r.NodeWriterStack[:len(r.NodeWriterStack)-1]
		indent := len(node.Marker) + 1
		if 1 == node.ListData.Typ || (3 == node.ListData.Typ && 0 == node.ListData.BulletChar) {
			indent++
		}
		indentSpaces := bytes.Repeat([]byte{lex.ItemSpace}, indent)
		indentedLines := bytes.Buffer{}
		buf := writer.Bytes()
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
			listItemBuf.WriteString(strconv.Itoa(node.Num) + string(node.ListData.Delimiter))
		} else {
			listItemBuf.Write(node.Marker)
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

func (r *FormatRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if r.Option.KramdownIAL {
		parent := node.Parent
		if nil != parent.Next && ast.NodeKramdownBlockIAL == parent.Next.Type {
			liIAL := parent.Next
			r.Write(liIAL.Tokens)
			liIAL.Unlink()
		}
	}

	r.WriteByte(lex.ItemOpenBracket)
	if node.TaskListItemChecked {
		r.WriteByte('X')
	} else {
		r.WriteByte(lex.ItemSpace)
	}
	r.WriteByte(lex.ItemCloseBracket)
	return ast.WalkStop
}

func (r *FormatRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if node.ParentIs(ast.NodeTableCell) {
		r.WriteString("<hr/>")
	} else {
		r.WriteString("---\n\n")
	}
	return ast.WalkStop
}

func (r *FormatRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if !r.Option.SoftBreak2HardBreak {
		r.WriteString("\\\n")
	} else {
		if node.ParentIs(ast.NodeTableCell) {
			r.WriteString("<br/>")
		} else {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkStop
}

func (r *FormatRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	return ast.WalkStop
}

func (r *FormatRenderer) withoutKramdownIAL(node *ast.Node) bool {
	return !r.Option.KramdownIAL || 0 == len(node.KramdownIAL)
}
