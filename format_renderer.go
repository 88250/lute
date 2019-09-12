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
func (lute *Lute) newFormatRenderer(treeRoot *Node) Renderer {
	ret := &FormatRenderer{BaseRenderer: &BaseRenderer{rendererFuncs: map[int]RendererFunc{}, option: lute.options, treeRoot: treeRoot}}

	// 注册 CommonMark 渲染函数

	ret.rendererFuncs[NodeDocument] = ret.renderDocumentMarkdown
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraphMarkdown
	ret.rendererFuncs[NodeText] = ret.renderTextMarkdown
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpanMarkdown
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlockMarkdown
	ret.rendererFuncs[NodeMathBlock] = ret.renderMathBlockMarkdown
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasisMarkdown
	ret.rendererFuncs[NodeStrong] = ret.renderStrongMarkdown
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquoteMarkdown
	ret.rendererFuncs[NodeHeading] = ret.renderHeadingMarkdown
	ret.rendererFuncs[NodeList] = ret.renderListMarkdown
	ret.rendererFuncs[NodeListItem] = ret.renderListItemMarkdown
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreakMarkdown
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreakMarkdown
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreakMarkdown
	ret.rendererFuncs[NodeHTMLBlock] = ret.renderHTMLMarkdown
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTMLMarkdown
	ret.rendererFuncs[NodeLink] = ret.renderLinkMarkdown
	ret.rendererFuncs[NodeImage] = ret.renderImageMarkdown

	// 注册 GFM 渲染函数

	ret.rendererFuncs[NodeStrikethrough] = ret.renderStrikethroughMarkdown
	ret.rendererFuncs[NodeTaskListItemMarker] = ret.renderTaskListItemMarkerMarkdown
	ret.rendererFuncs[NodeTable] = ret.renderTableMarkdown
	ret.rendererFuncs[NodeTableHead] = ret.renderTableHeadMarkdown
	ret.rendererFuncs[NodeTableRow] = ret.renderTableRowMarkdown
	ret.rendererFuncs[NodeTableCell] = ret.renderTableCellMarkdown

	// Emoji 渲染函数

	ret.rendererFuncs[NodeEmojiUnicode] = ret.renderEmojiUnicodeMarkdown
	ret.rendererFuncs[NodeEmojiImg] = ret.renderEmojiImgMarkdown

	return ret
}

func (r *FormatRenderer) renderEmojiImgMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.emojiAlias)
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderEmojiUnicodeMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.emojiAlias)
	}
	return WalkStop, nil
}

// TODO: 表的格式化应该按最宽的单元格对齐内容

func (r *FormatRenderer) renderTableCellMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte(itemPipe)
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderTableRowMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.writeString("|\n")
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderTableHeadMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.writeString("|\n")
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

func (r *FormatRenderer) renderTableMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.newline()
		if !r.isLastNode(r.treeRoot, node) {
			r.writeByte(itemNewline)
		}
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderStrikethroughMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(bytes.Repeat(items{node.strongEmDelMarker}, node.strongEmDelMarkenLen))
	} else {
		r.write(bytes.Repeat(items{node.strongEmDelMarker}, node.strongEmDelMarkenLen))
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderImageMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("![")
		r.write(node.firstChild.tokens)
		r.writeString("](")
		r.write(node.destination)
		if nil != node.title {
			r.writeString(" \"")
			r.write(node.title)
			r.writeByte(itemDoublequote)
		}
		r.writeByte(itemCloseParen)
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderLinkMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("[")
		if nil != node.firstChild {
			// FIXME: 未解决链接嵌套，另外还需要考虑链接引用定义
			r.write(node.firstChild.tokens)
		}
		r.writeString("](")
		r.write(node.destination)
		if nil != node.title {
			r.writeString(" \"")
			r.write(node.title)
			r.writeByte(itemDoublequote)
		}
		r.writeByte(itemCloseParen)
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderHTMLMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.write(node.tokens)
		r.newline()
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderInlineHTMLMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderDocumentMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.writer)
	} else {
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]
		buf := bytes.Trim(r.writer.Bytes(), " \t\n")
		r.writer.Reset()
		r.write(buf)
		r.writeByte(itemNewline)
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderParagraphMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.newline()

		inTightList := false
		lastListItemLastPara := false
		if parent := node.parent; nil != parent {
			if NodeListItem == parent.typ { // ListItem.Paragraph
				listItem := parent

				// 必须通过列表（而非列表项）上的紧凑标识判断，因为在设置该标识时仅设置了 List.tight
				// 设置紧凑标识的具体实现可参考函数 List.Finalize()
				inTightList = listItem.parent.tight

				nextItem := listItem.next
				if nil == nextItem {
					nextPara := node.next
					lastListItemLastPara = nil == nextPara
				}
			}
		}

		if !inTightList || (lastListItemLastPara) {
			r.writeByte(itemNewline)
		}
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderTextMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if typ := node.parent.typ; NodeLink != typ && NodeImage != typ {
			r.write(escapeHTML(node.tokens))
		}
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderCodeSpanMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte(itemBacktick)
		if 1 < node.codeMarkerLen {
			r.writeByte(itemBacktick)
			r.writeByte(itemSpace)
		}
		r.write(node.tokens)
		return WalkSkipChildren, nil
	}

	if 1 < node.codeMarkerLen {
		r.writeByte(itemSpace)
		r.writeByte(itemBacktick)
	}
	r.writeByte(itemBacktick)
	return WalkContinue, nil
}

func (r *FormatRenderer) renderMathBlockMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.write(node.tokens)
		return WalkSkipChildren, nil
	}
	if !r.isLastNode(r.treeRoot, node) {
		r.writeString("\n\n")
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderCodeBlockMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if !node.isFencedCodeBlock {
		node.codeBlockFenceLen = 3
	}
	if entering {
		r.write(bytes.Repeat(items{itemBacktick}, node.codeBlockFenceLen))
		r.write(node.codeBlockInfo)
		r.writeByte(itemNewline)
		r.write(node.tokens)
		return WalkSkipChildren, nil
	}
	r.write(bytes.Repeat(items{itemBacktick}, node.codeBlockFenceLen))
	r.newline()
	if !r.isLastNode(r.treeRoot, node) {
		r.writeByte(itemNewline)
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderEmphasisMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(bytes.Repeat(items{node.strongEmDelMarker}, node.strongEmDelMarkenLen))
	} else {
		r.write(bytes.Repeat(items{node.strongEmDelMarker}, node.strongEmDelMarkenLen))
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderStrongMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(bytes.Repeat(items{node.strongEmDelMarker}, node.strongEmDelMarkenLen))
	} else {
		r.write(bytes.Repeat(items{node.strongEmDelMarker}, node.strongEmDelMarkenLen))
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderBlockquoteMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writer = &bytes.Buffer{}
		r.nodeWriterStack = append(r.nodeWriterStack, r.writer)
	} else {
		writer := r.nodeWriterStack[len(r.nodeWriterStack)-1]
		r.nodeWriterStack = r.nodeWriterStack[:len(r.nodeWriterStack)-1]

		blockquoteLines := bytes.Buffer{}
		buf := writer.Bytes()
		lines := bytes.Split(buf, items{itemNewline})
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
		r.write(buf)
		r.writeString("\n\n")
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderHeadingMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(bytes.Repeat(items{itemCrosshatch}, node.headingLevel)) // 统一使用 ATX 标题，不使用 Setext 标题
		r.writeByte(itemSpace)
	} else {
		r.newline()
		r.writeByte(itemNewline)
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderListMarkdown(node *Node, entering bool) (WalkStatus, error) {
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
		r.write(buf)
		r.writeString("\n\n")
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderListItemMarkdown(node *Node, entering bool) (WalkStatus, error) {
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
		indentSpaces := bytes.Repeat(items{itemSpace}, indent)
		indentedLines := bytes.Buffer{}
		buf := writer.Bytes()
		lines := bytes.Split(buf, items{itemNewline})
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
		buf = buf[indent:]

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
		r.writer = r.nodeWriterStack[len(r.nodeWriterStack)-1]
		buf = bytes.TrimSpace(r.writer.Bytes())
		r.writer.Reset()
		r.write(buf)
		r.writeString("\n")
	}
	return WalkContinue, nil
}

func (r *FormatRenderer) renderTaskListItemMarkerMarkdown(node *Node, entering bool) (WalkStatus, error) {
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

func (r *FormatRenderer) renderThematicBreakMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.writeString("---")
		r.newline()
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderHardBreakMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if !r.option.SoftBreak2HardBreak {
			r.writeString("\\\n")
		} else {
			r.writeByte(itemNewline)
		}
	}
	return WalkStop, nil
}

func (r *FormatRenderer) renderSoftBreakMarkdown(node *Node, entering bool) (WalkStatus, error) {
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
