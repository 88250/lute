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
	"fmt"
)

// newHTMLRenderer 创建一个 HTML 渲染器。
func newHTMLRenderer(option options) (ret *Renderer) {
	ret = &Renderer{rendererFuncs: map[int]RendererFunc{}, option: option}

	// 注册 CommonMark 渲染函数

	ret.rendererFuncs[NodeDocument] = ret.renderDocumentHTML
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraphHTML
	ret.rendererFuncs[NodeText] = ret.renderTextHTML
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpanHTML
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlockHTML
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasisHTML
	ret.rendererFuncs[NodeStrong] = ret.renderStrongHTML
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquoteHTML
	ret.rendererFuncs[NodeHeading] = ret.renderHeadingHTML
	ret.rendererFuncs[NodeList] = ret.renderListHTML
	ret.rendererFuncs[NodeListItem] = ret.renderListItemHTML
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreakHTML
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreakHTML
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreakHTML
	ret.rendererFuncs[NodeHTMLBlock] = ret.renderHTMLHTML
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTMLHTML
	ret.rendererFuncs[NodeLink] = ret.renderLinkHTML
	ret.rendererFuncs[NodeImage] = ret.renderImageHTML

	// 注册 GFM 渲染函数

	ret.rendererFuncs[NodeStrikethrough] = ret.renderStrikethroughHTML
	ret.rendererFuncs[NodeTaskListItemMarker] = ret.renderTaskListItemMarkerHTML
	ret.rendererFuncs[NodeTable] = ret.renderTableHTML
	ret.rendererFuncs[NodeTableHead] = ret.renderTableHeadHTML
	ret.rendererFuncs[NodeTableRow] = ret.renderTableRowHTML
	ret.rendererFuncs[NodeTableCell] = ret.renderTableCellHTML

	return
}

func (r *Renderer) renderTableCellHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	tag := "td"
	if NodeTableHead == node.parent.typ {
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
		r.tag(tag, attrs, false)
	} else {
		r.tag("/"+tag, nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableRowHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("tr", nil, false)
		r.newline()
	} else {
		r.tag("/tr", nil, false)
		r.newline()
		if node == node.parent.lastChild {
			r.tag("/tbody", nil, false)
		}
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableHeadHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("thead", nil, false)
		r.newline()
		r.tag("tr", nil, false)
		r.newline()
	} else {
		r.tag("/tr", nil, false)
		r.newline()
		r.tag("/thead", nil, false)
		r.newline()
		if nil != node.Next() {
			r.tag("tbody", nil, false)
		}
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("table", nil, false)
		r.newline()
	} else {
		r.tag("/table", nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrikethroughHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("del", nil, false)
	} else {
		r.tag("/del", nil, false)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderImageHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		if 0 == r.disableTags {
			r.writeString("<img src=\"")
			r.write(escapeHTML(node.Destination))
			r.writeString("\" alt=\"")
		}
		r.disableTags++
		return WalkContinue, nil
	}

	r.disableTags--
	if 0 == r.disableTags {
		r.writeString("\"")
		if nil != node.Title {
			r.writeString(" title=\"")
			r.write(escapeHTML(node.Title))
			r.writeString("\"")
		}
		r.writeString(" />")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderLinkHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		attrs := [][]string{{"href", fromItems(escapeHTML(node.Destination))}}
		if nil != node.Title {
			attrs = append(attrs, []string{"title", fromItems(escapeHTML(node.Title))})
		}
		r.tag("a", attrs, false)

		return WalkContinue, nil
	}

	r.tag("/a", nil, false)
	return WalkContinue, nil
}

func (r *Renderer) renderHTMLHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.newline()
	r.write(node.Tokens())
	r.newline()
	return WalkContinue, nil
}

func (r *Renderer) renderInlineHTMLHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.write(node.Tokens())
	return WalkContinue, nil
}

func (r *Renderer) renderDocumentHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *Renderer) renderParagraphHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if grandparent := node.parent.parent; nil != grandparent {
		if NodeList == grandparent.typ { // List.ListItem.Paragraph
			if grandparent.tight {
				return WalkContinue, nil
			}
		}
	}

	if entering {
		r.newline()
		r.tag("p", nil, false)
	} else {
		r.tag("/p", nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTextHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.write(escapeHTML(node.tokens))
	return WalkContinue, nil
}

func (r *Renderer) renderCodeSpanHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<code>")
		r.write(escapeHTML(node.Tokens()))
		return WalkSkipChildren, nil
	}
	r.writeString("</code>")
	return WalkContinue, nil
}

func (r *Renderer) renderEmphasisHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("em", nil, false)
	} else {
		r.tag("/em", nil, false)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrongHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<strong>")
		r.write(node.Tokens())
	} else {
		r.writeString("</strong>")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderBlockquoteHTML(n *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.writeString("<blockquote>")
		r.newline()
	} else {
		r.newline()
		r.writeString("</blockquote>")
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderHeadingHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.writeString("<h" + " 123456"[node.HeadingLevel:node.HeadingLevel+1] + ">")
	} else {
		r.writeString("</h" + " 123456"[node.HeadingLevel:node.HeadingLevel+1] + ">")
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	tag := "ul"
	if 1 == node.listData.typ {
		tag = "ol"
	}
	if entering {
		r.newline()
		attrs := [][]string{{"start", fmt.Sprintf("%d", node.start)}}
		if nil == node.bulletChar && 1 != node.start {
			r.tag(tag, attrs, false)
		} else {
			r.tag(tag, nil, false)
		}
		r.newline()
	} else {
		r.newline()
		r.tag("/"+tag, nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListItemHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("li", nil, false)
	} else {
		r.tag("/li", nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTaskListItemMarkerHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		var attrs [][]string
		if node.taskListItemChecked {
			attrs = append(attrs, []string{"checked", ""})
		}
		attrs = append(attrs, []string{"disabled", ""}, []string{"type", "checkbox"})
		r.tag("input", attrs, true)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderThematicBreakHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.tag("hr", nil, true)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderHardBreakHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("br", nil, true)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderSoftBreakHTML(node *BaseNode, entering bool) (WalkStatus, error) {
	if entering {
		if r.option.SoftBreak2HardBreak {
			r.tag("br", nil, true)
			r.newline()
		} else {
			r.newline()
		}
	}
	return WalkContinue, nil
}

func (r *Renderer) tag(name string, attrs [][]string, selfclosing bool) {
	if r.disableTags > 0 {
		return
	}

	r.writeString("<")
	r.write(toItems(name))
	if 0 < len(attrs) {
		for _, attr := range attrs {
			r.writeString(" " + attr[0] + "=\"" + attr[1] + "\"")
		}
	}
	if selfclosing {
		r.writeString(" /")
	}
	r.writeString(">")
}
