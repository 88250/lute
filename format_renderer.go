// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lute

import (
	"bytes"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// FormatRenderer 描述了格式化渲染器。
type FormatRenderer struct {
	*BaseRenderer
	nodeWriterStack []*bytes.Buffer // 节点输出缓冲栈
}

// newFormatRenderer 创建一个格式化渲染器。
func (lute *Lute) newFormatRenderer(tree *Tree) Renderer {
	ret := &FormatRenderer{BaseRenderer: lute.newBaseRenderer(tree)}
	ret.rendererFuncs[NodeDocument] = ret.renderDocument
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraph
	ret.rendererFuncs[NodeText] = ret.renderText
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpan
	ret.rendererFuncs[NodeCodeSpanOpenMarker] = ret.renderCodeSpanOpenMarker
	ret.rendererFuncs[NodeCodeSpanContent] = ret.renderCodeSpanContent
	ret.rendererFuncs[NodeCodeSpanCloseMarker] = ret.renderCodeSpanCloseMarker
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlock
	ret.rendererFuncs[NodeCodeBlockFenceOpenMarker] = ret.renderCodeBlockOpenMarker
	ret.rendererFuncs[NodeCodeBlockFenceInfoMarker] = ret.renderCodeBlockInfoMarker
	ret.rendererFuncs[NodeCodeBlockCode] = ret.renderCodeBlockCode
	ret.rendererFuncs[NodeCodeBlockFenceCloseMarker] = ret.renderCodeBlockCloseMarker
	ret.rendererFuncs[NodeMathBlock] = ret.renderMathBlock
	ret.rendererFuncs[NodeMathBlockOpenMarker] = ret.renderMathBlockOpenMarker
	ret.rendererFuncs[NodeMathBlockContent] = ret.renderMathBlockContent
	ret.rendererFuncs[NodeMathBlockCloseMarker] = ret.renderMathBlockCloseMarker
	ret.rendererFuncs[NodeInlineMath] = ret.renderInlineMath
	ret.rendererFuncs[NodeInlineMathOpenMarker] = ret.renderInlineMathOpenMarker
	ret.rendererFuncs[NodeInlineMathContent] = ret.renderInlineMathContent
	ret.rendererFuncs[NodeInlineMathCloseMarker] = ret.renderInlineMathCloseMarker
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasis
	ret.rendererFuncs[NodeEmA6kOpenMarker] = ret.renderEmAsteriskOpenMarker
	ret.rendererFuncs[NodeEmA6kCloseMarker] = ret.renderEmAsteriskCloseMarker
	ret.rendererFuncs[NodeEmU8eOpenMarker] = ret.renderEmUnderscoreOpenMarker
	ret.rendererFuncs[NodeEmU8eCloseMarker] = ret.renderEmUnderscoreCloseMarker
	ret.rendererFuncs[NodeStrong] = ret.renderStrong
	ret.rendererFuncs[NodeStrongA6kOpenMarker] = ret.renderStrongA6kOpenMarker
	ret.rendererFuncs[NodeStrongA6kCloseMarker] = ret.renderStrongA6kCloseMarker
	ret.rendererFuncs[NodeStrongU8eOpenMarker] = ret.renderStrongU8eOpenMarker
	ret.rendererFuncs[NodeStrongU8eCloseMarker] = ret.renderStrongU8eCloseMarker
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquote
	ret.rendererFuncs[NodeBlockquoteMarker] = ret.renderBlockquoteMarker
	ret.rendererFuncs[NodeHeading] = ret.renderHeading
	ret.rendererFuncs[NodeHeadingC8hMarker] = ret.renderHeadingC8hMarker
	ret.rendererFuncs[NodeList] = ret.renderList
	ret.rendererFuncs[NodeListItem] = ret.renderListItem
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreak
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreak
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreak
	ret.rendererFuncs[NodeHTMLBlock] = ret.renderHTML
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTML
	ret.rendererFuncs[NodeLink] = ret.renderLink
	ret.rendererFuncs[NodeImage] = ret.renderImage
	ret.rendererFuncs[NodeBang] = ret.renderBang
	ret.rendererFuncs[NodeOpenBracket] = ret.renderOpenBracket
	ret.rendererFuncs[NodeCloseBracket] = ret.renderCloseBracket
	ret.rendererFuncs[NodeOpenParen] = ret.renderOpenParen
	ret.rendererFuncs[NodeCloseParen] = ret.renderCloseParen
	ret.rendererFuncs[NodeLinkText] = ret.renderLinkText
	ret.rendererFuncs[NodeLinkSpace] = ret.renderLinkSpace
	ret.rendererFuncs[NodeLinkDest] = ret.renderLinkDest
	ret.rendererFuncs[NodeLinkTitle] = ret.renderLinkTitle
	ret.rendererFuncs[NodeStrikethrough] = ret.renderStrikethrough
	ret.rendererFuncs[NodeStrikethrough1OpenMarker] = ret.renderStrikethrough1OpenMarker
	ret.rendererFuncs[NodeStrikethrough1CloseMarker] = ret.renderStrikethrough1CloseMarker
	ret.rendererFuncs[NodeStrikethrough2OpenMarker] = ret.renderStrikethrough2OpenMarker
	ret.rendererFuncs[NodeStrikethrough2CloseMarker] = ret.renderStrikethrough2CloseMarker
	ret.rendererFuncs[NodeTaskListItemMarker] = ret.renderTaskListItemMarker
	ret.rendererFuncs[NodeTable] = ret.renderTable
	ret.rendererFuncs[NodeTableHead] = ret.renderTableHead
	ret.rendererFuncs[NodeTableRow] = ret.renderTableRow
	ret.rendererFuncs[NodeTableCell] = ret.renderTableCell
	ret.rendererFuncs[NodeEmoji] = ret.renderEmoji
	ret.rendererFuncs[NodeEmojiUnicode] = ret.renderEmojiUnicode
	ret.rendererFuncs[NodeEmojiImg] = ret.renderEmojiImg
	ret.rendererFuncs[NodeEmojiAlias] = ret.renderEmojiAlias
	ret.rendererFuncs[NodeFootnotesDef] = ret.renderFootnotesDef
	ret.rendererFuncs[NodeFootnotesRef] = ret.renderFootnotesRef
	ret.rendererFuncs[NodeBackslash] = ret.renderBackslash
	ret.rendererFuncs[NodeBackslashContent] = ret.renderBackslashContent
	return ret
}

func (r *FormatRenderer) renderBackslashContent(node *Node, entering bool) WalkStatus {
	r.write(node.Tokens)
	return WalkStop
}

func (r *FormatRenderer) renderBackslash(node *Node, entering bool) WalkStatus {
	if entering {
		r.writeByte(itemBackslash)
	}
	return WalkContinue
}

func (r *FormatRenderer) renderFootnotesRef(node *Node, entering bool) WalkStatus {
	r.writeString("[" + bytesToStr(node.Tokens) + "]")
	return WalkStop
}

func (r *FormatRenderer) renderFootnotesDef(node *Node, entering bool) WalkStatus {
	if entering {
		r.writeString("[" + bytesToStr(node.Tokens) + "]: ")
	}
	return WalkContinue
}

func (r *FormatRenderer) renderEmojiAlias(node *Node, entering bool) WalkStatus {
	r.write(node.Tokens)
	return WalkStop
}

func (r *FormatRenderer) renderEmojiImg(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *FormatRenderer) renderEmojiUnicode(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *FormatRenderer) renderEmoji(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *FormatRenderer) renderTableCell(node *Node, entering bool) WalkStatus {
	padding := node.TableCellContentMaxWidth - node.TableCellContentWidth
	if entering {
		r.writeByte(itemPipe)
		r.writeByte(itemSpace)
		switch node.TableCellAlign {
		case 2:
			r.write(bytes.Repeat([]byte{itemSpace}, padding/2))
		case 3:
			r.write(bytes.Repeat([]byte{itemSpace}, padding))
		}
	} else {
		switch node.TableCellAlign {
		case 2:
			r.write(bytes.Repeat([]byte{itemSpace}, padding/2))
		case 3:
		default:
			r.write(bytes.Repeat([]byte{itemSpace}, padding))
		}
		r.writeByte(itemSpace)
	}
	return WalkContinue
}

func (r *FormatRenderer) renderTableRow(node *Node, entering bool) WalkStatus {
	if !entering {
		r.writeString("|\n")
	}
	return WalkContinue
}

func (r *FormatRenderer) renderTableHead(node *Node, entering bool) WalkStatus {
	if !entering {
		headRow := node.FirstChild
		for th := headRow.FirstChild; nil != th; th = th.Next {
			align := th.TableCellAlign
			switch align {
			case 0:
				r.writeString("| -")
				if padding := th.TableCellContentMaxWidth - 1; 0 < padding {
					r.write(bytes.Repeat([]byte{itemHyphen}, padding))
				}
				r.writeByte(itemSpace)
			case 1:
				r.writeString("| :-")
				if padding := th.TableCellContentMaxWidth - 2; 0 < padding {
					r.write(bytes.Repeat([]byte{itemHyphen}, padding))
				}
				r.writeByte(itemSpace)
			case 2:
				r.writeString("| :-")
				if padding := th.TableCellContentMaxWidth - 3; 0 < padding {
					r.write(bytes.Repeat([]byte{itemHyphen}, padding))
				}
				r.writeString(": ")
			case 3:
				r.writeString("| -")
				if padding := th.TableCellContentMaxWidth - 2; 0 < padding {
					r.write(bytes.Repeat([]byte{itemHyphen}, padding))
				}
				r.writeString(": ")
			}
		}
		r.writeString("|\n")
	}
	return WalkContinue
}

func (r *FormatRenderer) renderTable(node *Node, entering bool) WalkStatus {
	if !entering {
		r.newline()
		if !r.isLastNode(r.tree.Root, node) {
			r.writeByte(itemNewline)
		}
	}
	return WalkContinue
}

func (r *FormatRenderer) renderStrikethrough(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *FormatRenderer) renderStrikethrough1OpenMarker(node *Node, entering bool) WalkStatus {
	r.writeByte(itemTilde)
	return WalkStop
}

func (r *FormatRenderer) renderStrikethrough1CloseMarker(node *Node, entering bool) WalkStatus {
	r.writeByte(itemTilde)
	return WalkStop
}

func (r *FormatRenderer) renderStrikethrough2OpenMarker(node *Node, entering bool) WalkStatus {
	r.writeString("~~")
	return WalkStop
}

func (r *FormatRenderer) renderStrikethrough2CloseMarker(node *Node, entering bool) WalkStatus {
	r.writeString("~~")
	return WalkStop
}

func (r *FormatRenderer) renderLinkTitle(node *Node, entering bool) WalkStatus {
	r.writeString("\"")
	r.write(node.Tokens)
	r.writeString("\"")
	return WalkStop
}

func (r *FormatRenderer) renderLinkDest(node *Node, entering bool) WalkStatus {
	r.write(node.Tokens)
	return WalkStop
}

func (r *FormatRenderer) renderLinkSpace(node *Node, entering bool) WalkStatus {
	r.writeByte(itemSpace)
	return WalkStop
}

func (r *FormatRenderer) renderLinkText(node *Node, entering bool) WalkStatus {
	r.write(node.Tokens)
	return WalkStop
}

func (r *FormatRenderer) renderCloseParen(node *Node, entering bool) WalkStatus {
	r.writeByte(itemCloseParen)
	return WalkStop
}

func (r *FormatRenderer) renderOpenParen(node *Node, entering bool) WalkStatus {
	r.writeByte(itemOpenParen)
	return WalkStop
}

func (r *FormatRenderer) renderCloseBracket(node *Node, entering bool) WalkStatus {
	r.writeByte(itemCloseBracket)
	return WalkStop
}

func (r *FormatRenderer) renderOpenBracket(node *Node, entering bool) WalkStatus {
	r.writeByte(itemOpenBracket)
	return WalkStop
}

func (r *FormatRenderer) renderBang(node *Node, entering bool) WalkStatus {
	r.writeByte(itemBang)
	return WalkStop
}

func (r *FormatRenderer) renderImage(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *FormatRenderer) renderLink(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *FormatRenderer) renderHTML(node *Node, entering bool) WalkStatus {
	r.newline()
	r.write(node.Tokens)
	r.newline()
	if !r.isLastNode(r.tree.Root, node) {
		r.writeByte(itemNewline)
	}
	return WalkStop
}

func (r *FormatRenderer) renderInlineHTML(node *Node, entering bool) WalkStatus {
	r.write(node.Tokens)
	return WalkStop
}

func (r *FormatRenderer) renderDocument(node *Node, entering bool) WalkStatus {
	if entering {
		r.writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.writer)
	} else {
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]
		buf := bytes.Trim(r.writer.Bytes(), " \t\n")
		r.writer.Reset()
		r.writeBytes(buf)
		r.writeByte(itemNewline)
	}
	return WalkContinue
}

func (r *FormatRenderer) renderParagraph(node *Node, entering bool) WalkStatus {
	if !entering {
		r.newline()

		inTightList := false
		lastListItemLastPara := false
		if parent := node.Parent; nil != parent {
			if NodeListItem == parent.Type { // ListItem.Paragraph
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
			r.writeByte(itemNewline)
		}
	}
	return WalkContinue
}

func (r *FormatRenderer) renderText(node *Node, entering bool) WalkStatus {
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
	return WalkStop
}

func (r *FormatRenderer) renderCodeSpan(node *Node, entering bool) WalkStatus {
	if entering {
		if r.option.AutoSpace {
			if text := node.PreviousNodeText(); "" != text {
				lastc, _ := utf8.DecodeLastRuneInString(text)
				if unicode.IsLetter(lastc) || unicode.IsDigit(lastc) {
					r.writeByte(itemSpace)
				}
			}
		}
	} else {
		if r.option.AutoSpace {
			if text := node.NextNodeText(); "" != text {
				firstc, _ := utf8.DecodeRuneInString(text)
				if unicode.IsLetter(firstc) || unicode.IsDigit(firstc) {
					r.writeByte(itemSpace)
				}
			}
		}
	}
	return WalkContinue
}

func (r *FormatRenderer) renderCodeSpanOpenMarker(node *Node, entering bool) WalkStatus {
	r.writeByte(itemBacktick)
	if 1 < node.Parent.CodeMarkerLen {
		r.writeByte(itemBacktick)
		text := bytesToStr(node.Next.Tokens)
		firstc, _ := utf8.DecodeRuneInString(text)
		if '`' == firstc {
			r.writeByte(itemSpace)
		}
	}
	return WalkStop
}

func (r *FormatRenderer) renderCodeSpanContent(node *Node, entering bool) WalkStatus {
	r.write(node.Tokens)
	return WalkStop
}

func (r *FormatRenderer) renderCodeSpanCloseMarker(node *Node, entering bool) WalkStatus {
	if 1 < node.Parent.CodeMarkerLen {
		text := bytesToStr(node.Previous.Tokens)
		lastc, _ := utf8.DecodeLastRuneInString(text)
		if '`' == lastc {
			r.writeByte(itemSpace)
		}
		r.writeByte(itemBacktick)
	}
	r.writeByte(itemBacktick)
	return WalkStop
}

func (r *FormatRenderer) renderInlineMathCloseMarker(node *Node, entering bool) WalkStatus {
	r.writeByte(itemDollar)
	return WalkStop
}

func (r *FormatRenderer) renderInlineMathContent(node *Node, entering bool) WalkStatus {
	r.write(node.Tokens)
	return WalkStop
}

func (r *FormatRenderer) renderInlineMathOpenMarker(node *Node, entering bool) WalkStatus {
	r.writeByte(itemDollar)
	return WalkStop
}

func (r *FormatRenderer) renderInlineMath(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *FormatRenderer) renderMathBlockCloseMarker(node *Node, entering bool) WalkStatus {
	r.write(mathBlockMarker)
	r.writeByte(itemNewline)
	return WalkStop
}

func (r *FormatRenderer) renderMathBlockContent(node *Node, entering bool) WalkStatus {
	r.write(node.Tokens)
	r.writeByte(itemNewline)
	return WalkStop
}

func (r *FormatRenderer) renderMathBlockOpenMarker(node *Node, entering bool) WalkStatus {
	r.write(mathBlockMarker)
	r.writeByte(itemNewline)
	return WalkStop
}

func (r *FormatRenderer) renderMathBlock(node *Node, entering bool) WalkStatus {
	r.newline()
	if !entering && !r.isLastNode(r.tree.Root, node) {
		r.writeByte(itemNewline)
	}
	return WalkContinue
}

func (r *FormatRenderer) renderCodeBlockCloseMarker(node *Node, entering bool) WalkStatus {
	r.newline()
	r.writeBytes(bytes.Repeat([]byte{itemBacktick}, node.CodeBlockFenceLen))
	r.newline()
	if !r.isLastNode(r.tree.Root, node) {
		r.writeByte(itemNewline)
	}
	return WalkStop
}

func (r *FormatRenderer) renderCodeBlockCode(node *Node, entering bool) WalkStatus {
	r.write(node.Tokens)
	return WalkStop
}

func (r *FormatRenderer) renderCodeBlockInfoMarker(node *Node, entering bool) WalkStatus {
	r.write(node.CodeBlockInfo)
	r.writeByte(itemNewline)
	return WalkStop
}

func (r *FormatRenderer) renderCodeBlockOpenMarker(node *Node, entering bool) WalkStatus {
	r.writeBytes(bytes.Repeat([]byte{itemBacktick}, node.CodeBlockFenceLen))
	return WalkStop
}

func (r *FormatRenderer) renderCodeBlock(node *Node, entering bool) WalkStatus {
	if !node.IsFencedCodeBlock {
		r.writeBytes(bytes.Repeat([]byte{itemBacktick}, 3))
		r.writeByte(itemNewline)
		r.write(node.FirstChild.Tokens)
		r.writeBytes(bytes.Repeat([]byte{itemBacktick}, 3))
		r.newline()
		if !r.isLastNode(r.tree.Root, node) {
			r.writeByte(itemNewline)
		}
		return WalkStop
	}
	return WalkContinue
}

func (r *FormatRenderer) renderEmphasis(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *FormatRenderer) renderEmAsteriskOpenMarker(node *Node, entering bool) WalkStatus {
	r.writeByte(itemAsterisk)
	return WalkStop
}

func (r *FormatRenderer) renderEmAsteriskCloseMarker(node *Node, entering bool) WalkStatus {
	r.writeByte(itemAsterisk)
	return WalkStop
}

func (r *FormatRenderer) renderEmUnderscoreOpenMarker(node *Node, entering bool) WalkStatus {
	r.writeByte(itemUnderscore)
	return WalkStop
}

func (r *FormatRenderer) renderEmUnderscoreCloseMarker(node *Node, entering bool) WalkStatus {
	r.writeByte(itemUnderscore)
	return WalkStop
}

func (r *FormatRenderer) renderStrong(node *Node, entering bool) WalkStatus {
	return WalkContinue
}

func (r *FormatRenderer) renderStrongA6kOpenMarker(node *Node, entering bool) WalkStatus {
	r.writeString("**")
	return WalkStop
}

func (r *FormatRenderer) renderStrongA6kCloseMarker(node *Node, entering bool) WalkStatus {
	r.writeString("**")
	return WalkStop
}

func (r *FormatRenderer) renderStrongU8eOpenMarker(node *Node, entering bool) WalkStatus {
	r.writeString("__")
	return WalkStop
}

func (r *FormatRenderer) renderStrongU8eCloseMarker(node *Node, entering bool) WalkStatus {
	r.writeString("__")
	return WalkStop
}

func (r *FormatRenderer) renderBlockquote(node *Node, entering bool) WalkStatus {
	if entering {
		r.writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.writer)
	} else {
		writer := r.nodeWriterStack[len(r.nodeWriterStack)-1]
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]

		blockquoteLines := bytes.Buffer{}
		buf := writer.Bytes()
		lines := bytes.Split(buf, []byte{itemNewline})
		length := len(lines)
		if 2 < length && isBlank(lines[length-1]) && isBlank(lines[length-2]) {
			lines = lines[:length-1]
		}
		if 1 == len(r.nodeWriterStack) { // 已经是根这一层
			length = len(lines)
			if 1 < length && isBlank(lines[length-1]) {
				lines = lines[:length-1]
			}
		}

		length = len(lines)
		for _, line := range lines {
			if 0 == len(line) {
				blockquoteLines.WriteString(">\n")
				continue
			}

			if itemGreater == line[0] {
				blockquoteLines.WriteString(">")
			} else {
				blockquoteLines.WriteString("> ")
			}
			blockquoteLines.Write(line)
			blockquoteLines.WriteByte(itemNewline)
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
	return WalkContinue
}

func (r *FormatRenderer) renderBlockquoteMarker(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *FormatRenderer) renderHeading(node *Node, entering bool) WalkStatus {
	if entering {
		r.writeBytes(bytes.Repeat([]byte{itemCrosshatch}, node.HeadingLevel)) // 统一使用 ATX 标题，不使用 Setext 标题
		r.writeByte(itemSpace)
	} else {
		r.newline()
		r.writeByte(itemNewline)
	}
	return WalkContinue
}

func (r *FormatRenderer) renderHeadingC8hMarker(node *Node, entering bool) WalkStatus {
	return WalkStop
}

func (r *FormatRenderer) renderList(node *Node, entering bool) WalkStatus {
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
	return WalkContinue
}

func (r *FormatRenderer) renderListItem(node *Node, entering bool) WalkStatus {
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
		indentSpaces := bytes.Repeat([]byte{itemSpace}, indent)
		indentedLines := bytes.Buffer{}
		buf := writer.Bytes()
		lines := bytes.Split(buf, []byte{itemNewline})
		for _, line := range lines {
			if 0 == len(line) {
				indentedLines.WriteByte(itemNewline)
				continue
			}
			indentedLines.Write(indentSpaces)
			indentedLines.Write(line)
			indentedLines.WriteByte(itemNewline)
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
		listItemBuf.WriteByte(itemSpace)
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
	return WalkContinue
}

func (r *FormatRenderer) renderTaskListItemMarker(node *Node, entering bool) WalkStatus {
	r.writeByte(itemOpenBracket)
	if node.TaskListItemChecked {
		r.writeByte('X')
	} else {
		r.writeByte(itemSpace)
	}
	r.writeByte(itemCloseBracket)
	return WalkStop
}

func (r *FormatRenderer) renderThematicBreak(node *Node, entering bool) WalkStatus {
	r.writeString("---\n\n")
	return WalkStop
}

func (r *FormatRenderer) renderHardBreak(node *Node, entering bool) WalkStatus {
	if !r.option.SoftBreak2HardBreak {
		r.writeString("\\\n")
	} else {
		r.writeByte(itemNewline)
	}
	return WalkStop
}

func (r *FormatRenderer) renderSoftBreak(node *Node, entering bool) WalkStatus {
	r.newline()
	return WalkStop
}

func (r *FormatRenderer) isLastNode(treeRoot, node *Node) bool {
	if treeRoot == node {
		return true
	}
	if nil != node.Next {
		return false
	}
	if NodeDocument == node.Parent.Type {
		return treeRoot.LastChild == node
	}

	var n *Node
	for n = node.Parent; ; n = n.Parent {
		if NodeDocument == n.Parent.Type {
			break
		}
	}
	return treeRoot.LastChild == n
}
