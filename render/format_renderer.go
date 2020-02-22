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
	nodeWriterStack []*bytes.Buffer // 节点输出缓冲栈
}

// NewFormatRenderer 创建一个格式化渲染器。
func NewFormatRenderer(tree *parse.Tree) Renderer {
	ret := &FormatRenderer{BaseRenderer: newBaseRenderer(tree)}
	ret.rendererFuncs[ast.NodeDocument] = ret.renderDocument
	ret.rendererFuncs[ast.NodeParagraph] = ret.renderParagraph
	ret.rendererFuncs[ast.NodeText] = ret.renderText
	ret.rendererFuncs[ast.NodeCodeSpan] = ret.renderCodeSpan
	ret.rendererFuncs[ast.NodeCodeSpanOpenMarker] = ret.renderCodeSpanOpenMarker
	ret.rendererFuncs[ast.NodeCodeSpanContent] = ret.renderCodeSpanContent
	ret.rendererFuncs[ast.NodeCodeSpanCloseMarker] = ret.renderCodeSpanCloseMarker
	ret.rendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock
	ret.rendererFuncs[ast.NodeCodeBlockFenceOpenMarker] = ret.renderCodeBlockOpenMarker
	ret.rendererFuncs[ast.NodeCodeBlockFenceInfoMarker] = ret.renderCodeBlockInfoMarker
	ret.rendererFuncs[ast.NodeCodeBlockCode] = ret.renderCodeBlockCode
	ret.rendererFuncs[ast.NodeCodeBlockFenceCloseMarker] = ret.renderCodeBlockCloseMarker
	ret.rendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock
	ret.rendererFuncs[ast.NodeMathBlockOpenMarker] = ret.renderMathBlockOpenMarker
	ret.rendererFuncs[ast.NodeMathBlockContent] = ret.renderMathBlockContent
	ret.rendererFuncs[ast.NodeMathBlockCloseMarker] = ret.renderMathBlockCloseMarker
	ret.rendererFuncs[ast.NodeInlineMath] = ret.renderInlineMath
	ret.rendererFuncs[ast.NodeInlineMathOpenMarker] = ret.renderInlineMathOpenMarker
	ret.rendererFuncs[ast.NodeInlineMathContent] = ret.renderInlineMathContent
	ret.rendererFuncs[ast.NodeInlineMathCloseMarker] = ret.renderInlineMathCloseMarker
	ret.rendererFuncs[ast.NodeEmphasis] = ret.renderEmphasis
	ret.rendererFuncs[ast.NodeEmA6kOpenMarker] = ret.renderEmAsteriskOpenMarker
	ret.rendererFuncs[ast.NodeEmA6kCloseMarker] = ret.renderEmAsteriskCloseMarker
	ret.rendererFuncs[ast.NodeEmU8eOpenMarker] = ret.renderEmUnderscoreOpenMarker
	ret.rendererFuncs[ast.NodeEmU8eCloseMarker] = ret.renderEmUnderscoreCloseMarker
	ret.rendererFuncs[ast.NodeStrong] = ret.renderStrong
	ret.rendererFuncs[ast.NodeStrongA6kOpenMarker] = ret.renderStrongA6kOpenMarker
	ret.rendererFuncs[ast.NodeStrongA6kCloseMarker] = ret.renderStrongA6kCloseMarker
	ret.rendererFuncs[ast.NodeStrongU8eOpenMarker] = ret.renderStrongU8eOpenMarker
	ret.rendererFuncs[ast.NodeStrongU8eCloseMarker] = ret.renderStrongU8eCloseMarker
	ret.rendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote
	ret.rendererFuncs[ast.NodeBlockquoteMarker] = ret.renderBlockquoteMarker
	ret.rendererFuncs[ast.NodeHeading] = ret.renderHeading
	ret.rendererFuncs[ast.NodeHeadingC8hMarker] = ret.renderHeadingC8hMarker
	ret.rendererFuncs[ast.NodeList] = ret.renderList
	ret.rendererFuncs[ast.NodeListItem] = ret.renderListItem
	ret.rendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	ret.rendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak
	ret.rendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak
	ret.rendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	ret.rendererFuncs[ast.NodeInlineHTML] = ret.renderInlineHTML
	ret.rendererFuncs[ast.NodeLink] = ret.renderLink
	ret.rendererFuncs[ast.NodeImage] = ret.renderImage
	ret.rendererFuncs[ast.NodeBang] = ret.renderBang
	ret.rendererFuncs[ast.NodeOpenBracket] = ret.renderOpenBracket
	ret.rendererFuncs[ast.NodeCloseBracket] = ret.renderCloseBracket
	ret.rendererFuncs[ast.NodeOpenParen] = ret.renderOpenParen
	ret.rendererFuncs[ast.NodeCloseParen] = ret.renderCloseParen
	ret.rendererFuncs[ast.NodeLinkText] = ret.renderLinkText
	ret.rendererFuncs[ast.NodeLinkSpace] = ret.renderLinkSpace
	ret.rendererFuncs[ast.NodeLinkDest] = ret.renderLinkDest
	ret.rendererFuncs[ast.NodeLinkTitle] = ret.renderLinkTitle
	ret.rendererFuncs[ast.NodeStrikethrough] = ret.renderStrikethrough
	ret.rendererFuncs[ast.NodeStrikethrough1OpenMarker] = ret.renderStrikethrough1OpenMarker
	ret.rendererFuncs[ast.NodeStrikethrough1CloseMarker] = ret.renderStrikethrough1CloseMarker
	ret.rendererFuncs[ast.NodeStrikethrough2OpenMarker] = ret.renderStrikethrough2OpenMarker
	ret.rendererFuncs[ast.NodeStrikethrough2CloseMarker] = ret.renderStrikethrough2CloseMarker
	ret.rendererFuncs[ast.NodeTaskListItemMarker] = ret.renderTaskListItemMarker
	ret.rendererFuncs[ast.NodeTable] = ret.renderTable
	ret.rendererFuncs[ast.NodeTableHead] = ret.renderTableHead
	ret.rendererFuncs[ast.NodeTableRow] = ret.renderTableRow
	ret.rendererFuncs[ast.NodeTableCell] = ret.renderTableCell
	ret.rendererFuncs[ast.NodeEmoji] = ret.renderEmoji
	ret.rendererFuncs[ast.NodeEmojiUnicode] = ret.renderEmojiUnicode
	ret.rendererFuncs[ast.NodeEmojiImg] = ret.renderEmojiImg
	ret.rendererFuncs[ast.NodeEmojiAlias] = ret.renderEmojiAlias
	ret.rendererFuncs[ast.NodeFootnotesDef] = ret.renderFootnotesDef
	ret.rendererFuncs[ast.NodeFootnotesRef] = ret.renderFootnotesRef
	ret.rendererFuncs[ast.NodeBackslash] = ret.renderBackslash
	ret.rendererFuncs[ast.NodeBackslashContent] = ret.renderBackslashContent
	return ret
}

func (r *FormatRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.writeByte(lex.ItemBackslash)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("[" + util.BytesToStr(node.Tokens) + "]")
	return ast.WalkStop
}

func (r *FormatRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.writeString("[" + util.BytesToStr(node.Tokens) + "]: ")
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
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
		r.writeByte(lex.ItemPipe)
		r.writeByte(lex.ItemSpace)
		switch node.TableCellAlign {
		case 2:
			r.write(bytes.Repeat([]byte{lex.ItemSpace}, padding/2))
		case 3:
			r.write(bytes.Repeat([]byte{lex.ItemSpace}, padding))
		}
	} else {
		switch node.TableCellAlign {
		case 2:
			r.write(bytes.Repeat([]byte{lex.ItemSpace}, padding/2))
		case 3:
		default:
			r.write(bytes.Repeat([]byte{lex.ItemSpace}, padding))
		}
		r.writeByte(lex.ItemSpace)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		r.writeString("|\n")
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
				r.writeString("| -")
				if padding := th.TableCellContentMaxWidth - 1; 0 < padding {
					r.write(bytes.Repeat([]byte{lex.ItemHyphen}, padding))
				}
				r.writeByte(lex.ItemSpace)
			case 1:
				r.writeString("| :-")
				if padding := th.TableCellContentMaxWidth - 2; 0 < padding {
					r.write(bytes.Repeat([]byte{lex.ItemHyphen}, padding))
				}
				r.writeByte(lex.ItemSpace)
			case 2:
				r.writeString("| :-")
				if padding := th.TableCellContentMaxWidth - 3; 0 < padding {
					r.write(bytes.Repeat([]byte{lex.ItemHyphen}, padding))
				}
				r.writeString(": ")
			case 3:
				r.writeString("| -")
				if padding := th.TableCellContentMaxWidth - 2; 0 < padding {
					r.write(bytes.Repeat([]byte{lex.ItemHyphen}, padding))
				}
				r.writeString(": ")
			}
		}
		r.writeString("|\n")
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		r.newline()
		if !r.isLastNode(r.tree.Root, node) {
			r.writeByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.textAutoSpacePrevious(node)
	} else {
		r.textAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemTilde)
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemTilde)
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("~~")
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("~~")
	return ast.WalkStop
}

func (r *FormatRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("\"")
	r.write(node.Tokens)
	r.writeString("\"")
	return ast.WalkStop
}

func (r *FormatRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemSpace)
	return ast.WalkStop
}

func (r *FormatRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if r.option.AutoSpace {
		r.space(node)
	}
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemCloseParen)
	return ast.WalkStop
}

func (r *FormatRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemOpenParen)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemCloseBracket)
	return ast.WalkStop
}

func (r *FormatRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemOpenBracket)
	return ast.WalkStop
}

func (r *FormatRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemBang)
	return ast.WalkStop
}

func (r *FormatRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *FormatRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.linkTextAutoSpacePrevious(node)
	} else {
		r.linkTextAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.newline()
	r.write(node.Tokens)
	r.newline()
	if !r.isLastNode(r.tree.Root, node) {
		r.writeByte(lex.ItemNewline)
	}
	return ast.WalkStop
}

func (r *FormatRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.writer)
	} else {
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]
		buf := bytes.Trim(r.writer.Bytes(), " \t\n")
		r.writer.Reset()
		r.writeBytes(buf)
		r.writeByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		r.newline()

		inTightList := false
		lastListItemLastPara := false
		if parent := node.Parent; nil != parent {
			if ast.NodeListItem == parent.Type { // ListItem.Paragraph
				listItem := parent
				if nil != listItem.Parent && nil != listItem.Parent.ListData {
					// 必须通过列表（而非列表项）上的紧凑标识判断，因为在设置该标识时仅设置了 List.Tight
					// 设置紧凑标识的具体实现可参考函数 List.Finalize()
					inTightList = listItem.Parent.Tight

					nextItem := listItem.Next
					if nil == nextItem {
						nextPara := node.Next
						lastListItemLastPara = nil == nextPara
					}
				} else {
					inTightList = true
				}
			}
		}

		if !inTightList || (lastListItemLastPara) {
			r.writeByte(lex.ItemNewline)
		}
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if r.option.AutoSpace {
		r.space(node)
	}
	if r.option.FixTermTypo {
		r.fixTermTypo(node)
	}
	if r.option.ChinesePunct {
		r.chinesePunct(node)
	}
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if r.option.AutoSpace {
			if text := node.PreviousNodeText(); "" != text {
				lastc, _ := utf8.DecodeLastRuneInString(text)
				if unicode.IsLetter(lastc) || unicode.IsDigit(lastc) {
					r.writeByte(lex.ItemSpace)
				}
			}
		}
	} else {
		if r.option.AutoSpace {
			if text := node.NextNodeText(); "" != text {
				firstc, _ := utf8.DecodeRuneInString(text)
				if unicode.IsLetter(firstc) || unicode.IsDigit(firstc) {
					r.writeByte(lex.ItemSpace)
				}
			}
		}
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemBacktick)
	if 1 < node.Parent.CodeMarkerLen {
		r.writeByte(lex.ItemBacktick)
		text := util.BytesToStr(node.Next.Tokens)
		firstc, _ := utf8.DecodeRuneInString(text)
		if '`' == firstc {
			r.writeByte(lex.ItemSpace)
		}
	}
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if 1 < node.Parent.CodeMarkerLen {
		text := util.BytesToStr(node.Previous.Tokens)
		lastc, _ := utf8.DecodeLastRuneInString(text)
		if '`' == lastc {
			r.writeByte(lex.ItemSpace)
		}
		r.writeByte(lex.ItemBacktick)
	}
	r.writeByte(lex.ItemBacktick)
	return ast.WalkStop
}

func (r *FormatRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemDollar)
	return ast.WalkStop
}

func (r *FormatRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemDollar)
	return ast.WalkStop
}

func (r *FormatRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *FormatRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(parse.MathBlockMarker)
	r.writeByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *FormatRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	r.writeByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *FormatRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(parse.MathBlockMarker)
	r.writeByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *FormatRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.newline()
	if !entering && !r.isLastNode(r.tree.Root, node) {
		r.writeByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.newline()
	r.writeBytes(bytes.Repeat([]byte{lex.ItemBacktick}, node.CodeBlockFenceLen))
	r.newline()
	if !r.isLastNode(r.tree.Root, node) {
		r.writeByte(lex.ItemNewline)
	}
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.Tokens)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.write(node.CodeBlockInfo)
	r.writeByte(lex.ItemNewline)
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeBytes(bytes.Repeat([]byte{lex.ItemBacktick}, node.CodeBlockFenceLen))
	return ast.WalkStop
}

func (r *FormatRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.newline()
	}
	if !node.IsFencedCodeBlock {
		r.writeBytes(bytes.Repeat([]byte{lex.ItemBacktick}, 3))
		r.writeByte(lex.ItemNewline)
		r.write(node.FirstChild.Tokens)
		r.writeBytes(bytes.Repeat([]byte{lex.ItemBacktick}, 3))
		r.newline()
		if !r.isLastNode(r.tree.Root, node) {
			r.writeByte(lex.ItemNewline)
		}
		return ast.WalkStop
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.textAutoSpacePrevious(node)
	} else {
		r.textAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemAsterisk)
	return ast.WalkStop
}

func (r *FormatRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemAsterisk)
	return ast.WalkStop
}

func (r *FormatRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemUnderscore)
	return ast.WalkStop
}

func (r *FormatRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemUnderscore)
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.textAutoSpacePrevious(node)
	} else {
		r.textAutoSpaceNext(node)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("**")
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("**")
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("__")
	return ast.WalkStop
}

func (r *FormatRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("__")
	return ast.WalkStop
}

func (r *FormatRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.writer)
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
		r.writer = r.nodeWriterStack[len(r.nodeWriterStack)-1]
		buf = bytes.TrimSpace(r.writer.Bytes())
		r.writer.Reset()
		r.writeBytes(buf)
		r.writeString("\n\n")
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *FormatRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.writeBytes(bytes.Repeat([]byte{lex.ItemCrosshatch}, node.HeadingLevel)) // 统一使用 ATX 标题，不使用 Setext 标题
		r.writeByte(lex.ItemSpace)
	} else {
		r.newline()
		r.writeByte(lex.ItemNewline)
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *FormatRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.writer)
	} else {
		writer := r.nodeWriterStack[len(r.nodeWriterStack)-1]
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]
		r.nodeWriterStack[len(r.nodeWriterStack)-1].Write(writer.Bytes())
		r.writer = r.nodeWriterStack[len(r.nodeWriterStack)-1]
		buf := bytes.TrimSpace(r.writer.Bytes())
		r.writer.Reset()
		r.writeBytes(buf)
		r.writeString("\n\n")
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.writer)
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
		r.writer = r.nodeWriterStack[len(r.nodeWriterStack)-1]
		buf = bytes.TrimSpace(r.writer.Bytes())
		r.writer.Reset()
		r.writeBytes(buf)
		r.writeString("\n")
	}
	return ast.WalkContinue
}

func (r *FormatRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeByte(lex.ItemOpenBracket)
	if node.TaskListItemChecked {
		r.writeByte('X')
	} else {
		r.writeByte(lex.ItemSpace)
	}
	r.writeByte(lex.ItemCloseBracket)
	return ast.WalkStop
}

func (r *FormatRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("---\n\n")
	return ast.WalkStop
}

func (r *FormatRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if !r.option.SoftBreak2HardBreak {
		r.writeString("\\\n")
	} else {
		r.writeByte(lex.ItemNewline)
	}
	return ast.WalkStop
}

func (r *FormatRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.newline()
	return ast.WalkStop
}

func (r *FormatRenderer) isLastNode(treeRoot, node *ast.Node) bool {
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
