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
	"strings"
)

// EChartsJSONRenderer 描述了 JSON 渲染器。
type EChartsJSONRenderer struct {
	*BaseRenderer
}

// newEChartsJSONRenderer 创建一个 ECharts JSON 渲染器。
func (lute *Lute) newEChartsJSONRenderer(tree *Tree) Renderer {
	ret := &EChartsJSONRenderer{lute.newBaseRenderer(tree)}
	ret.rendererFuncs[NodeDocument] = ret.renderDocumentEChartsJSON
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraphEChartsJSON
	ret.rendererFuncs[NodeText] = ret.renderTextEChartsJSON
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpanEChartsJSON
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlockEChartsJSON
	ret.rendererFuncs[NodeMathBlock] = ret.renderMathBlockEChartsJSON
	ret.rendererFuncs[NodeInlineMath] = ret.renderInlineMathEChartsJSON
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasisEChartsJSON
	ret.rendererFuncs[NodeStrong] = ret.renderStrongEChartsJSON
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquoteEChartsJSON
	ret.rendererFuncs[NodeHeading] = ret.renderHeadingEChartsJSON
	ret.rendererFuncs[NodeList] = ret.renderListEChartsJSON
	ret.rendererFuncs[NodeListItem] = ret.renderListItemEChartsJSON
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreakEChartsJSON
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreakEChartsJSON
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreakEChartsJSON
	ret.rendererFuncs[NodeHTMLBlock] = ret.renderHTMLEChartsJSON
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTMLEChartsJSON
	ret.rendererFuncs[NodeLink] = ret.renderLinkEChartsJSON
	ret.rendererFuncs[NodeImage] = ret.renderImageEChartsJSON
	ret.rendererFuncs[NodeStrikethrough] = ret.renderStrikethroughEChartsJSON
	ret.rendererFuncs[NodeTaskListItemMarker] = ret.renderTaskListItemMarkerEChartsJSON
	ret.rendererFuncs[NodeTable] = ret.renderTableEChartsJSON
	ret.rendererFuncs[NodeTableHead] = ret.renderTableHeadEChartsJSON
	ret.rendererFuncs[NodeTableRow] = ret.renderTableRowEChartsJSON
	ret.rendererFuncs[NodeTableCell] = ret.renderTableCellEChartsJSON
	ret.rendererFuncs[NodeEmojiUnicode] = ret.renderEmojiUnicodeEChartsJSON
	ret.rendererFuncs[NodeEmojiImg] = ret.renderEmojiImgEChartsJSON

	ret.defaultRendererFunc = ret.renderDefault
	return ret
}

func (r *EChartsJSONRenderer) renderDefault(n *Node, entering bool) (WalkStatus, error) {
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderInlineMathEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Inline Math\nspan", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderMathBlockEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Math Block\ndiv", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderEmojiImgEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Emoji Img\n", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderEmojiUnicodeEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Emoji Unicode\n", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderTableCellEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Table Cell\ntd", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderTableRowEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Table Row\ntr", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderTableHeadEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Table Head\nthead", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderTableEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Table\ntable", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderStrikethroughEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Strikethrough\ndel", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderImageEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.openObj()
		r.val("Image\nimg", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderLinkEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.openObj()
		r.val("Link\na", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderHTMLEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("HTML Block\n", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderInlineHTMLEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Inline HTML\n", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderDocumentEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte(itemOpenBracket)
		r.openObj()
		r.val("Document", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
		r.writeByte(itemCloseBracket)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderParagraphEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.openObj()
		r.val("Paragraph\np", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderTextEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		text := bytesToStr(node.tokens)
		var i int
		summary := ""
		for _, r := range text {
			i++
			summary += string(r)
			if 4 < i {
				summary += "..."
				break
			}
		}
		r.openObj()
		r.val("Text\n"+summary, node)
		r.closeObj(node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderCodeSpanEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Code Span\ncode", node)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderEmphasisEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.openObj()
		r.val("Emphasis\nem", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderStrongEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.openObj()
		r.val("Strong\nstrong", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderBlockquoteEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.openObj()
		r.val("Blockquote\nblockquote", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderHeadingEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.openObj()
		h := "h" + " 123456"[node.headingLevel:node.headingLevel+1]
		r.val("Heading\n"+h, node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderListEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.openObj()
		list := "ul"
		if 1 == node.listData.typ {
			list = "ol"
		}
		r.val("List\n"+list, node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderListItemEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.openObj()
		r.val("List Item\nli "+bytesToStr(node.listData.marker), node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderTaskListItemMarkerEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.openObj()
		check := " "
		if node.taskListItemChecked {
			check = "X"
		}
		r.val("Task List Item Marker\n["+check+"]", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return WalkContinue, nil
}

func (r *EChartsJSONRenderer) renderThematicBreakEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Thematic Break\nhr", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderHardBreakEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Hard Break\nbr", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderSoftBreakEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Soft Break\n", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) renderCodeBlockEChartsJSON(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.leaf("Code Block\npre.code", node)
	}
	return WalkStop, nil
}

func (r *EChartsJSONRenderer) leaf(val string, node *Node) {
	r.openObj()
	r.val(val, node)
	r.closeObj(node)
}

func (r *EChartsJSONRenderer) val(val string, node *Node) {
	val = strings.ReplaceAll(val, "\\", "\\\\")
	val = strings.ReplaceAll(val, "\n", "\\n")
	val = strings.ReplaceAll(val, "\"", "")
	val = strings.ReplaceAll(val, "'", "")
	r.writeString("\"name\":\"" + val + "\"")
}

func (r *EChartsJSONRenderer) openObj() {
	r.writeByte('{')
}

func (r *EChartsJSONRenderer) closeObj(node *Node) {
	r.writeByte('}')
	if !r.ignore(node.next) {
		r.comma()
	}
}

func (r *EChartsJSONRenderer) openChildren(node *Node) {
	if nil != node.firstChild {
		r.writeString(",\"children\":[")
	}
}

func (r *EChartsJSONRenderer) closeChildren(node *Node) {
	if nil != node.firstChild {
		r.writeByte(']')
	}
}

func (r *EChartsJSONRenderer) comma() {
	r.writeString(",")
}

func (r *EChartsJSONRenderer) ignore(node *Node) bool {
	return nil == node ||
		// 以下类型的节点不进行渲染，否则图画出来节点太多
		NodeBlockquoteMarker == node.typ ||
		NodeEmA6kOpenMarker == node.typ || NodeEmA6kCloseMarker == node.typ ||
		NodeEmU8eOpenMarker == node.typ || NodeEmU8eCloseMarker == node.typ ||
		NodeStrongA6kOpenMarker == node.typ || NodeStrongA6kCloseMarker == node.typ ||
		NodeStrongU8eOpenMarker == node.typ || NodeStrongU8eCloseMarker == node.typ ||
		NodeStrikethrough1OpenMarker == node.typ || NodeStrikethrough1CloseMarker == node.typ ||
		NodeStrikethrough2OpenMarker == node.typ || NodeStrikethrough2CloseMarker == node.typ ||
		NodeMathBlockOpenMarker == node.typ || NodeMathBlockContent == node.typ || NodeMathBlockCloseMarker == node.typ ||
		NodeInlineMathOpenMarker == node.typ || NodeInlineMathContent == node.typ || NodeInlineMathCloseMarker == node.typ
}
