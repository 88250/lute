// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

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

func (r *FormatRenderer) renderBackslashContent(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.Tokens)
	return WalkStop, nil
}

func (r *FormatRenderer) renderBackslash(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte(itemBackslash)
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderFootnotesRef(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("[" + bytesToStr(node.Tokens) + "]")
	return WalkStop, nil
}

func (r *FormatRenderer) renderFootnotesDef(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("[" + bytesToStr(node.Tokens) + "]: ")
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderEmojiAlias(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.Tokens)
	return WalkStop, nil
}

func (r *FormatRenderer) renderEmojiImg(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *FormatRenderer) renderEmojiUnicode(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *FormatRenderer) renderEmoji(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *FormatRenderer) renderTableCell(node *Node, entering bool) (WalkStatus, error) {
	padding := node.tableCellContentMaxWidth - node.tableCellContentWidth
	if entering {
		r.writeByte(itemPipe)
		r.writeByte(itemSpace)
		switch node.tableCellAlign {
		case 2:
			r.write(bytes.Repeat([]byte{itemSpace}, padding/2))
		case 3:
			r.write(bytes.Repeat([]byte{itemSpace}, padding))
		}
	} else {
		switch node.tableCellAlign {
		case 2:
			r.write(bytes.Repeat([]byte{itemSpace}, padding/2))
		case 3:
		default:
			r.write(bytes.Repeat([]byte{itemSpace}, padding))
		}
		r.writeByte(itemSpace)
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderTableRow(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.writeString("|\n")
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderTableHead(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		headRow := node.FirstChild
		for th := headRow.FirstChild; nil != th; th = th.Next {
			align := th.tableCellAlign
			switch align {
			case 0:
				r.writeString("| -")
				if padding := th.tableCellContentMaxWidth - 1; 0 < padding {
					r.write(bytes.Repeat([]byte{itemHyphen}, padding))
				}
				r.writeByte(itemSpace)
			case 1:
				r.writeString("| :-")
				if padding := th.tableCellContentMaxWidth - 2; 0 < padding {
					r.write(bytes.Repeat([]byte{itemHyphen}, padding))
				}
				r.writeByte(itemSpace)
			case 2:
				r.writeString("| :-")
				if padding := th.tableCellContentMaxWidth - 3; 0 < padding {
					r.write(bytes.Repeat([]byte{itemHyphen}, padding))
				}
				r.writeString(": ")
			case 3:
				r.writeString("| -")
				if padding := th.tableCellContentMaxWidth - 2; 0 < padding {
					r.write(bytes.Repeat([]byte{itemHyphen}, padding))
				}
				r.writeString(": ")
			}
		}
		r.writeString("|\n")
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderTable(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.newline()
		if !r.isLastNode(r.tree.Root, node) {
			r.writeByte(itemNewline)
		}
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderStrikethrough(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *FormatRenderer) renderStrikethrough1OpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemTilde)
	return WalkStop, nil
}

func (r *FormatRenderer) renderStrikethrough1CloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemTilde)
	return WalkStop, nil
}

func (r *FormatRenderer) renderStrikethrough2OpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("~~")
	return WalkStop, nil
}

func (r *FormatRenderer) renderStrikethrough2CloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("~~")
	return WalkStop, nil
}

func (r *FormatRenderer) renderLinkTitle(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("\"")
	r.write(node.Tokens)
	r.writeString("\"")
	return WalkStop, nil
}

func (r *FormatRenderer) renderLinkDest(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.Tokens)
	return WalkStop, nil
}

func (r *FormatRenderer) renderLinkSpace(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemSpace)
	return WalkStop, nil
}

func (r *FormatRenderer) renderLinkText(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.Tokens)
	return WalkStop, nil
}

func (r *FormatRenderer) renderCloseParen(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemCloseParen)
	return WalkStop, nil
}

func (r *FormatRenderer) renderOpenParen(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemOpenParen)
	return WalkStop, nil
}

func (r *FormatRenderer) renderCloseBracket(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemCloseBracket)
	return WalkStop, nil
}

func (r *FormatRenderer) renderOpenBracket(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemOpenBracket)
	return WalkStop, nil
}

func (r *FormatRenderer) renderBang(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemBang)
	return WalkStop, nil
}

func (r *FormatRenderer) renderImage(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *FormatRenderer) renderLink(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *FormatRenderer) renderHTML(node *Node, entering bool) (WalkStatus, error) {
	r.newline()
	r.write(node.Tokens)
	r.newline()
	if !r.isLastNode(r.tree.Root, node) {
		r.writeByte(itemNewline)
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderInlineHTML(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.Tokens)
	return WalkStop, nil
}

func (r *FormatRenderer) renderDocument(node *Node, entering bool) (WalkStatus, error) {
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
	return WalkContinue, nil
}

func (r *FormatRenderer) renderParagraph(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.newline()

		inTightList := false
		lastListItemLastPara := false
		if parent := node.Parent; nil != parent {
			if NodeListItem == parent.Typ { // ListItem.Paragraph
				listItem := parent
				if nil != listItem.Parent && nil != listItem.Parent.listData {
					// 必须通过列表（而非列表项）上的紧凑标识判断，因为在设置该标识时仅设置了 List.tight
					// 设置紧凑标识的具体实现可参考函数 List.Finalize()
					inTightList = listItem.Parent.tight

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
	return WalkContinue, nil
}

func (r *FormatRenderer) renderText(node *Node, entering bool) (WalkStatus, error) {
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
	return WalkStop, nil
}

func (r *FormatRenderer) renderCodeSpan(node *Node, entering bool) (WalkStatus, error) {
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
	return WalkContinue, nil
}

func (r *FormatRenderer) renderCodeSpanOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemBacktick)
	if 1 < node.Parent.codeMarkerLen {
		r.writeByte(itemBacktick)
		text := bytesToStr(node.Next.Tokens)
		firstc, _ := utf8.DecodeRuneInString(text)
		if '`' == firstc {
			r.writeByte(itemSpace)
		}
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderCodeSpanContent(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.Tokens)
	return WalkStop, nil
}

func (r *FormatRenderer) renderCodeSpanCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	if 1 < node.Parent.codeMarkerLen {
		text := bytesToStr(node.Previous.Tokens)
		lastc, _ := utf8.DecodeLastRuneInString(text)
		if '`' == lastc {
			r.writeByte(itemSpace)
		}
		r.writeByte(itemBacktick)
	}
	r.writeByte(itemBacktick)
	return WalkStop, nil
}

func (r *FormatRenderer) renderInlineMathCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemDollar)
	return WalkStop, nil
}

func (r *FormatRenderer) renderInlineMathContent(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.Tokens)
	return WalkStop, nil
}

func (r *FormatRenderer) renderInlineMathOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemDollar)
	return WalkStop, nil
}

func (r *FormatRenderer) renderInlineMath(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *FormatRenderer) renderMathBlockCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.write(mathBlockMarker)
	r.writeByte(itemNewline)
	return WalkStop, nil
}

func (r *FormatRenderer) renderMathBlockContent(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.Tokens)
	r.writeByte(itemNewline)
	return WalkStop, nil
}

func (r *FormatRenderer) renderMathBlockOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.write(mathBlockMarker)
	r.writeByte(itemNewline)
	return WalkStop, nil
}

func (r *FormatRenderer) renderMathBlock(node *Node, entering bool) (WalkStatus, error) {
	r.newline()
	if !entering && !r.isLastNode(r.tree.Root, node) {
		r.writeByte(itemNewline)
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderCodeBlockCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.newline()
	r.writeBytes(bytes.Repeat([]byte{itemBacktick}, node.codeBlockFenceLen))
	r.newline()
	if !r.isLastNode(r.tree.Root, node) {
		r.writeByte(itemNewline)
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderCodeBlockCode(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.Tokens)
	return WalkStop, nil
}

func (r *FormatRenderer) renderCodeBlockInfoMarker(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.codeBlockInfo)
	r.writeByte(itemNewline)
	return WalkStop, nil
}

func (r *FormatRenderer) renderCodeBlockOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeBytes(bytes.Repeat([]byte{itemBacktick}, node.codeBlockFenceLen))
	return WalkStop, nil
}

func (r *FormatRenderer) renderCodeBlock(node *Node, entering bool) (WalkStatus, error) {
	if !node.isFencedCodeBlock {
		r.writeBytes(bytes.Repeat([]byte{itemBacktick}, 3))
		r.writeByte(itemNewline)
		r.write(node.FirstChild.Tokens)
		r.writeBytes(bytes.Repeat([]byte{itemBacktick}, 3))
		r.newline()
		if !r.isLastNode(r.tree.Root, node) {
			r.writeByte(itemNewline)
		}
		return WalkStop, nil
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderEmphasis(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *FormatRenderer) renderEmAsteriskOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemAsterisk)
	return WalkStop, nil
}

func (r *FormatRenderer) renderEmAsteriskCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemAsterisk)
	return WalkStop, nil
}

func (r *FormatRenderer) renderEmUnderscoreOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemUnderscore)
	return WalkStop, nil
}

func (r *FormatRenderer) renderEmUnderscoreCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemUnderscore)
	return WalkStop, nil
}

func (r *FormatRenderer) renderStrong(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *FormatRenderer) renderStrongA6kOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("**")
	return WalkStop, nil
}

func (r *FormatRenderer) renderStrongA6kCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("**")
	return WalkStop, nil
}

func (r *FormatRenderer) renderStrongU8eOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("__")
	return WalkStop, nil
}

func (r *FormatRenderer) renderStrongU8eCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("__")
	return WalkStop, nil
}

func (r *FormatRenderer) renderBlockquote(node *Node, entering bool) (WalkStatus, error) {
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
	return WalkContinue, nil
}

func (r *FormatRenderer) renderBlockquoteMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *FormatRenderer) renderHeading(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeBytes(bytes.Repeat([]byte{itemCrosshatch}, node.headingLevel)) // 统一使用 ATX 标题，不使用 Setext 标题
		r.writeByte(itemSpace)
	} else {
		r.newline()
		r.writeByte(itemNewline)
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderHeadingC8hMarker(node *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *FormatRenderer) renderList(node *Node, entering bool) (WalkStatus, error) {
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
	return WalkContinue, nil
}

func (r *FormatRenderer) renderListItem(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.writer)
	} else {
		writer := r.nodeWriterStack[len(r.nodeWriterStack)-1]
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]
		indent := len(node.marker) + 1
		if 1 == node.listData.typ || (3 == node.listData.typ && 0 == node.listData.bulletChar) {
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
		if 1 == node.listData.typ || (3 == node.listData.typ && 0 == node.listData.bulletChar) {
			listItemBuf.WriteString(strconv.Itoa(node.num) + string(node.listData.delimiter))
		} else {
			listItemBuf.Write(node.marker)
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
	return WalkContinue, nil
}

func (r *FormatRenderer) renderTaskListItemMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemOpenBracket)
	if node.taskListItemChecked {
		r.writeByte('X')
	} else {
		r.writeByte(itemSpace)
	}
	r.writeByte(itemCloseBracket)
	return WalkStop, nil
}

func (r *FormatRenderer) renderThematicBreak(node *Node, entering bool) (WalkStatus, error) {
	r.writeString("---\n\n")
	return WalkStop, nil
}

func (r *FormatRenderer) renderHardBreak(node *Node, entering bool) (WalkStatus, error) {
	if !r.option.SoftBreak2HardBreak {
		r.writeString("\\\n")
	} else {
		r.writeByte(itemNewline)
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderSoftBreak(node *Node, entering bool) (WalkStatus, error) {
	r.newline()
	return WalkStop, nil
}

func (r *FormatRenderer) isLastNode(treeRoot, node *Node) bool {
	if treeRoot == node {
		return true
	}
	if nil != node.Next {
		return false
	}
	if NodeDocument == node.Parent.Typ {
		return treeRoot.LastChild == node
	}

	var n *Node
	for n = node.Parent; ; n = n.Parent {
		if NodeDocument == n.Parent.Typ {
			break
		}
	}
	return treeRoot.LastChild == n
}
