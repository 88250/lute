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
	"strconv"
	"strings"
)

// newFormatRenderer 创建一个格式化渲染器。
func newFormatRenderer(option options) (ret *Renderer) {
	ret = &Renderer{rendererFuncs: map[int]RendererFunc{}, option: option}

	// 注册 CommonMark 渲染函数

	ret.rendererFuncs[NodeDocument] = ret.renderDocumentMarkdown
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraphMarkdown
	ret.rendererFuncs[NodeText] = ret.renderTextMarkdown
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpanMarkdown
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlockMarkdown
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

	return
}

// TODO: 表的格式化应该按最宽的单元格对齐内容

func (r *Renderer) renderTableCellMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte('|')
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableRowMarkdown(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.writeString("|\n")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableHeadMarkdown(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.writeString("|\n")
		n := node.(*TableHead)
		table := n.Parent().(*Table)
		for i := 0; i < len(table.Aligns); i++ {
			align := table.Aligns[i]
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

func (r *Renderer) renderTableMarkdown(node Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *Renderer) renderStrikethroughMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("~~")
	} else {
		r.writeString("~~")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderImageMarkdown(node Node, entering bool) (WalkStatus, error) {
	n := node.(*Image)
	if entering {
		r.writeString("![")
		r.write(n.firstChild.Tokens())
		r.writeString("](" + n.Destination + "")
		if "" != n.Title {
			r.writeString(" \"" + n.Title + "\"")
		}
		r.writeString(")")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderLinkMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		n := node.(*Link)
		r.writeString("[")
		r.write(n.firstChild.Tokens()) // FIXME: 未解决链接嵌套，另外还需要考虑链接引用定义
		r.writeString("](" + n.Destination + "")
		if "" != n.Title {
			r.writeString(" \"" + n.Title + "\"")
		}
		r.writeString(")")
	}

	return WalkContinue, nil
}

func (r *Renderer) renderHTMLMarkdown(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.newline()
	r.write(node.Tokens())
	r.newline()
	return WalkContinue, nil
}

func (r *Renderer) renderInlineHTMLMarkdown(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.write(node.Tokens())
	return WalkContinue, nil
}

func (r *Renderer) renderDocumentMarkdown(node Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *Renderer) renderParagraphMarkdown(node Node, entering bool) (WalkStatus, error) {
	listPadding := 0
	inTightList := false
	if grandparent := node.Parent().Parent(); nil != grandparent {
		if list, ok := grandparent.(*List); ok { // List.ListItem.Paragraph
			inTightList = list.tight
			if node.Parent().FirstChild() != node {
				listPadding = list.padding
			}
		}
	}

	if entering {
		r.writeString(strings.Repeat(" ", listPadding))
	} else {
		r.newline()
		if !inTightList {
			r.writeString("\n")
		}
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTextMarkdown(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	if typ := node.Parent().Type(); NodeLink != typ && NodeImage != typ {
		r.write(node.Tokens())
	}
	return WalkContinue, nil
}

func (r *Renderer) renderCodeSpanMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte('`')
		r.write(node.Tokens())
		return WalkSkipChildren, nil
	}

	r.writeByte('`')
	return WalkContinue, nil
}

func (r *Renderer) renderCodeBlockMarkdown(node Node, entering bool) (WalkStatus, error) {
	listPadding := 0
	if grandparent := node.Parent().Parent(); nil != grandparent {
		if list, ok := grandparent.(*List); ok { // List.ListItem.CodeBlock
			if node.Parent().FirstChild() != node {
				listPadding = list.padding
			}
		}
	}

	n := node.(*CodeBlock)
	if entering {
		r.newline()
		if 0 < listPadding {
			r.writeString(strings.Repeat(" ", listPadding))
		}
		r.writeString(strings.Repeat("`", n.fenceLength))
		r.writeString(n.info + "\n")
		if 0 < listPadding {
			lines := n.tokens.split(itemNewline)
			length := len(lines)
			for i, line := range lines {
				r.writeString(strings.Repeat(" ", listPadding))
				r.write(line)
				if i < length-1 {
					r.writeByte('\n')
				}
			}
		} else {
			r.write(n.tokens)
		}

		r.newline()
		return WalkSkipChildren, nil
	}

	if 0 < listPadding {
		r.writeString(strings.Repeat(" ", listPadding))
	}
	strings.Repeat("`", n.fenceLength)
	r.writeString("\n\n")
	return WalkContinue, nil
}

func (r *Renderer) renderEmphasisMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte('*')
	} else {
		r.writeByte('*')
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrongMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("**")
	} else {
		r.writeString("**")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderBlockquoteMarkdown(n Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.writeString("> ") // 带个空格更好一些
	} else {
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderHeadingMarkdown(node Node, entering bool) (WalkStatus, error) {
	n := node.(*Heading)
	if entering {
		r.writeString(strings.Repeat("#", n.Level) + " ") // 统一使用 ATX 标题，不使用 Setext 标题
	} else {
		r.newline()
		r.writeString("\n")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListMarkdown(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.newline()
		r.writeString("\n")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListItemMarkdown(node Node, entering bool) (WalkStatus, error) {
	n := node.(*ListItem)
	if entering {
		r.writeString(strings.Repeat(" ", n.margin))
		if 1 == n.listData.typ {
			r.writeString(strconv.Itoa(n.num) + ".")
		} else {
			r.write(n.marker)
		}
		r.writeByte(' ')
	} else {
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTaskListItemMarkerMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		n := node.(*TaskListItemMarker)
		r.writeString("[")
		if n.checked {
			r.writeByte('X')
		} else {
			r.writeByte(' ')
		}
		r.writeString("]")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderThematicBreakMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("\n---\n")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderHardBreakMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("\n\n")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderSoftBreakMarkdown(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
	}
	return WalkContinue, nil
}
