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

// VditorSVRenderer 描述了 Vditor Split-View DOM 渲染器。
type VditorSVRenderer struct {
	*BaseRenderer
	nodeWriterStack []*bytes.Buffer // 节点输出缓冲栈
}

// NewVditorSVRenderer 创建一个 Vditor Split-View DOM 渲染器
func NewVditorSVRenderer(tree *parse.Tree) *VditorSVRenderer {
	ret := &VditorSVRenderer{BaseRenderer: NewBaseRenderer(tree)}
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

func (r *VditorSVRenderer) Render() (output []byte) {
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

func (r *VditorSVRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemBackslash)
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("[toc]\n\n")
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("[" + util.BytesToStr(node.Tokens) + "]")
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteString("[" + util.BytesToStr(node.Tokens) + "]: ")
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *VditorSVRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		r.WriteString("|\n")
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *VditorSVRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		r.Newline()
		if !r.isLastNode(r.Tree.Root, node) {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemTilde)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemTilde)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("~~")
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("~~")
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("\"")
	r.Write(node.Tokens)
	r.WriteString("\"")
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemSpace)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if r.Option.AutoSpace {
		r.Space(node)
	}
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemCloseParen)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemOpenParen)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemCloseBracket)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemOpenBracket)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemBang)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *VditorSVRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	r.Write(node.Tokens)
	r.Newline()
	if !r.isLastNode(r.Tree.Root, node) {
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]
		buf := bytes.Trim(r.Writer.Bytes(), " \t\n")
		r.Writer.Reset()
		r.Write(buf)
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		r.Newline()

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

		if !inTightList || (lastListItemLastPara) {
			r.WriteByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if r.Option.AutoSpace {
		r.Space(node)
	}
	if r.Option.FixTermTypo {
		r.FixTermTypo(node)
	}
	if r.Option.ChinesePunct {
		r.ChinesePunct(node)
	}
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *VditorSVRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *VditorSVRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *VditorSVRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemDollar)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemDollar)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(parse.MathBlockMarker)
	r.WriteByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	r.WriteByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(parse.MathBlockMarker)
	r.WriteByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	if !entering && !r.isLastNode(r.Tree.Root, node) {
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	r.Write(bytes.Repeat([]byte{lex.ItemBacktick}, node.CodeBlockFenceLen))
	r.Newline()
	if !r.isLastNode(r.Tree.Root, node) {
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.Tokens)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(node.CodeBlockInfo)
	r.WriteByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.Write(bytes.Repeat([]byte{lex.ItemBacktick}, node.CodeBlockFenceLen))
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
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
			r.WriteByte(lex.ItemNewline)
		}
		return ast.WalkStop
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemAsterisk)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemAsterisk)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemUnderscore)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemUnderscore)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.TextAutoSpacePrevious(node)
	} else {
		r.TextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("**")
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("**")
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("__")
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("__")
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		writer := r.nodeWriterStack[len(r.nodeWriterStack)-1]
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]

		blockquoteLines := bytes.Buffer{}
		buf := writer.Bytes()
		lines := bytes.Split(buf, []byte{lex.ItemNewline})
		length := len(lines)
		if 2 < length && lex.IsBlank(lines[length-1]) && lex.IsBlank(lines[length-2]) {
			lines = lines[:length-1]
		}
		if 1 == len(r.nodeWriterStack) { // 已经是根这一层
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
		r.nodeWriterStack[len(r.nodeWriterStack)-1].Write(writer.Bytes())
		r.Writer = r.nodeWriterStack[len(r.nodeWriterStack)-1]
		buf = bytes.TrimSpace(r.Writer.Bytes())
		r.Writer.Reset()
		r.Write(buf)
		r.WriteString("\n\n")
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if !node.HeadingSetext {
			r.Write(bytes.Repeat([]byte{lex.ItemCrosshatch}, node.HeadingLevel))
			r.WriteByte(lex.ItemSpace)
		}
	} else {
		if node.HeadingSetext {
			r.WriteByte(lex.ItemNewline)
			content := node.Text()
			contentLen := 0
			for _, r := range content {
				if utf8.RuneSelf <= r {
					contentLen += 2
				} else {
					contentLen++
				}
			}
			if 1 == node.HeadingLevel {
				r.WriteString(strings.Repeat("=", contentLen))
			} else if 2 == node.HeadingLevel {
				r.WriteString(strings.Repeat("-", contentLen))
			}
		}
		if r.Option.HeadingID && nil != node.HeadingID {
			r.WriteString(" {" + util.BytesToStr(node.HeadingID) + "}")
		}
		r.Newline()
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		writer := r.nodeWriterStack[len(r.nodeWriterStack)-1]
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]
		r.nodeWriterStack[len(r.nodeWriterStack)-1].Write(writer.Bytes())
		r.Writer = r.nodeWriterStack[len(r.nodeWriterStack)-1]
		buf := bytes.TrimSpace(r.Writer.Bytes())
		r.Writer.Reset()
		r.Write(buf)
		r.WriteString("\n\n")
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.Writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.Writer)
	} else {
		writer := r.nodeWriterStack[len(r.nodeWriterStack)-1]
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]
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
		writer.Reset()
		writer.Write(buf)
		r.nodeWriterStack[len(r.nodeWriterStack)-1].Write(writer.Bytes())
		r.Writer = r.nodeWriterStack[len(r.nodeWriterStack)-1]
		buf = bytes.TrimSpace(r.Writer.Bytes())
		r.Writer.Reset()
		r.Write(buf)
		r.WriteString("\n")
	}
	return ast.WalkContinue
}

func (r *VditorSVRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteByte(lex.ItemOpenBracket)
	if node.TaskListItemChecked {
		r.WriteByte('X')
	} else {
		r.WriteByte(lex.ItemSpace)
	}
	r.WriteByte(lex.ItemCloseBracket)
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("---\n\n")
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if !r.Option.SoftBreak2HardBreak {
		r.WriteString("\\\n")
	} else {
		r.WriteByte(lex.ItemNewline)
	}
	return ast.WalkStop
}

func (r *VditorSVRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.Newline()
	return ast.WalkStop
}

func (r *VditorSVRenderer) isLastNode(treeRoot, node *ast.Node) bool {
	if treeRoot == node {
		return true
	}
	if nil != node.Next {
		return false
	}
	if ast.NodeDocument == node.Parent.Type {
		return treeRoot.LastChild == node
	}

	var n *ast.Node
	for n = node.Parent; ; n = n.Parent {
		if ast.NodeDocument == n.Parent.Type {
			break
		}
	}
	return treeRoot.LastChild == n
}
