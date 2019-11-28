// Lute - A structured markdown engine.
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
	return ret
}

func (r *FormatRenderer) renderEmojiAlias(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
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

// TODO: 表的格式化应该按最宽的单元格对齐内容

func (r *FormatRenderer) renderTableCell(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte(itemPipe)
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
		table := node.parent
		for i := 0; i < len(table.tableAligns); i++ {
			align := table.tableAligns[i]
			switch align {
			case 0:
				r.writeString("|---")
			case 1:
				r.writeString("|:---")
			case 2:
				r.writeString("|:---:")
			case 3:
				r.writeString("|---:")
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
	r.write(node.tokens)
	r.writeString("\"")
	return WalkStop, nil
}

func (r *FormatRenderer) renderLinkDest(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *FormatRenderer) renderLinkSpace(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemSpace)
	return WalkStop, nil
}

func (r *FormatRenderer) renderLinkText(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
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
	r.write(node.tokens)
	r.newline()
	return WalkStop, nil
}

func (r *FormatRenderer) renderInlineHTML(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
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
		if parent := node.parent; nil != parent {
			if NodeListItem == parent.typ { // ListItem.Paragraph
				listItem := parent
				if nil != listItem.parent && nil != listItem.parent.listData {
					// 必须通过列表（而非列表项）上的紧凑标识判断，因为在设置该标识时仅设置了 List.tight
					// 设置紧凑标识的具体实现可参考函数 List.Finalize()
					inTightList = listItem.parent.tight

					nextItem := listItem.next
					if nil == nextItem {
						nextPara := node.next
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
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *FormatRenderer) renderCodeSpan(node *Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *FormatRenderer) renderCodeSpanOpenMarker(node *Node, entering bool) (WalkStatus, error) {
	r.writeByte(itemBacktick)
	if 1 < node.parent.codeMarkerLen {
		r.writeByte(itemBacktick)
		r.writeByte(itemSpace)
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderCodeSpanContent(node *Node, entering bool) (WalkStatus, error) {
	r.write(node.tokens)
	return WalkStop, nil
}

func (r *FormatRenderer) renderCodeSpanCloseMarker(node *Node, entering bool) (WalkStatus, error) {
	if 1 < node.parent.codeMarkerLen {
		r.writeByte(itemSpace)
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
	r.write(node.tokens)
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
	r.write(node.tokens)
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
	r.write(node.tokens)
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
		r.write(node.tokens)
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
		for _, line := range lines {
			if 0 == len(line) {
				blockquoteLines.WriteString(">\n")
				continue
			}

			blockquoteLines.WriteString("> ")
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
		if 1 == node.listData.typ {
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
		if 1 == node.listData.typ {
			listItemBuf.WriteString(strconv.Itoa(node.num) + ".")
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
	if entering {
		r.writeByte(itemOpenBracket)
		if node.taskListItemChecked {
			r.writeByte('X')
		} else {
			r.writeByte(itemSpace)
		}
		r.writeByte(itemCloseBracket)
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderThematicBreak(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.writeString("---")
		r.newline()
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderHardBreak(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if !r.option.SoftBreak2HardBreak {
			r.writeString("\\\n")
		} else {
			r.writeByte(itemNewline)
		}
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderSoftBreak(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
	}
	return WalkStop, nil
}

func (r *FormatRenderer) isLastNode(treeRoot, node *Node) bool {
	if treeRoot == node {
		return true
	}
	if nil != node.next {
		return false
	}
	if NodeDocument == node.parent.typ {
		return treeRoot.lastChild == node
	}

	var n *Node
	for n = node.parent; ; n = n.parent {
		if NodeDocument == n.parent.typ {
			break
		}
	}
	return treeRoot.lastChild == n
}
