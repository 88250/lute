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
	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
	"strings"
)

// EChartsJSONRenderer 描述了 JSON 渲染器。
type EChartsJSONRenderer struct {
	*BaseRenderer
}

// newEChartsJSONRenderer 创建一个 ECharts JSON 渲染器。
func (lute *Lute) newEChartsJSONRenderer(tree *parse.Tree) Renderer {
	ret := &EChartsJSONRenderer{lute.newBaseRenderer(tree)}
	ret.rendererFuncs[ast.NodeDocument] = ret.renderDocument
	ret.rendererFuncs[ast.NodeParagraph] = ret.renderParagraph
	ret.rendererFuncs[ast.NodeText] = ret.renderText
	ret.rendererFuncs[ast.NodeCodeSpan] = ret.renderCodeSpan
	ret.rendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock
	ret.rendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock
	ret.rendererFuncs[ast.NodeInlineMath] = ret.renderInlineMath
	ret.rendererFuncs[ast.NodeEmphasis] = ret.renderEmphasis
	ret.rendererFuncs[ast.NodeStrong] = ret.renderStrong
	ret.rendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote
	ret.rendererFuncs[ast.NodeHeading] = ret.renderHeading
	ret.rendererFuncs[ast.NodeList] = ret.renderList
	ret.rendererFuncs[ast.NodeListItem] = ret.renderListItem
	ret.rendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	ret.rendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak
	ret.rendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak
	ret.rendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	ret.rendererFuncs[ast.NodeInlineHTML] = ret.renderInlineHTML
	ret.rendererFuncs[ast.NodeLink] = ret.renderLink
	ret.rendererFuncs[ast.NodeImage] = ret.renderImage
	ret.rendererFuncs[ast.NodeStrikethrough] = ret.renderStrikethrough
	ret.rendererFuncs[ast.NodeTaskListItemMarker] = ret.renderTaskListItemMarker
	ret.rendererFuncs[ast.NodeTable] = ret.renderTable
	ret.rendererFuncs[ast.NodeTableHead] = ret.renderTableHead
	ret.rendererFuncs[ast.NodeTableRow] = ret.renderTableRow
	ret.rendererFuncs[ast.NodeTableCell] = ret.renderTableCell
	ret.rendererFuncs[ast.NodeEmoji] = ret.renderEmoji
	ret.rendererFuncs[ast.NodeEmojiUnicode] = ret.renderEmojiUnicode
	ret.rendererFuncs[ast.NodeEmojiImg] = ret.renderEmojiImg
	ret.rendererFuncs[ast.NodeEmojiAlias] = ret.renderEmojiAlias

	ret.defaultRendererFunc = ret.renderDefault
	return ret
}

func (r *EChartsJSONRenderer) renderDefault(n *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Inline Math\nspan", node)
	}
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Math Block\ndiv", node)
	}
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Emoji Img\n", node)
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Emoji Unicode\n", node)
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Table Cell\ntd", node)
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Table Row\ntr", node)
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Table Head\nthead", node)
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Table\ntable", node)
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Strikethrough\ndel", node)
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Image\nimg", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Link\na", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("HTML Block\n", node)
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Inline HTML\n", node)
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.writeByte(lex.ItemOpenBracket)
		r.openObj()
		r.val("Document", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
		r.writeByte(lex.ItemCloseBracket)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Paragraph\np", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		text := util.BytesToStr(node.Tokens)
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
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Code Span\ncode", node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Emphasis\nem", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Strong\nstrong", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Blockquote\nblockquote", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		h := "h" + " 123456"[node.HeadingLevel:node.HeadingLevel+1]
		r.val("Heading\n"+h, node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		list := "ul"
		if 1 == node.ListData.Typ || (3 == node.ListData.Typ && 0 == node.ListData.BulletChar) {
			list = "ol"
		}
		r.val("List\n"+list, node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("List Item\nli "+util.BytesToStr(node.ListData.Marker), node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		check := " "
		if node.TaskListItemChecked {
			check = "X"
		}
		r.val("Task List Item Marker\n["+check+"]", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Thematic Break\nhr", node)
	}
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Hard Break\nbr", node)
	}
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Soft Break\n", node)
	}
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Code Block\npre.code", node)
	}
	return ast.WalkStop
}

func (r *EChartsJSONRenderer) leaf(val string, node *ast.Node) {
	r.openObj()
	r.val(val, node)
	r.closeObj(node)
}

func (r *EChartsJSONRenderer) val(val string, node *ast.Node) {
	val = strings.ReplaceAll(val, "\\", "\\\\")
	val = strings.ReplaceAll(val, "\n", "\\n")
	val = strings.ReplaceAll(val, "\"", "")
	val = strings.ReplaceAll(val, "'", "")
	r.writeString("\"name\":\"" + val + "\"")
}

func (r *EChartsJSONRenderer) openObj() {
	r.writeByte('{')
}

func (r *EChartsJSONRenderer) closeObj(node *ast.Node) {
	r.writeByte('}')
	if !r.ignore(node.Next) {
		r.comma()
	}
}

func (r *EChartsJSONRenderer) openChildren(node *ast.Node) {
	if nil != node.FirstChild {
		r.writeString(",\"children\":[")
	}
}

func (r *EChartsJSONRenderer) closeChildren(node *ast.Node) {
	if nil != node.FirstChild {
		r.writeByte(']')
	}
}

func (r *EChartsJSONRenderer) comma() {
	r.writeString(",")
}

func (r *EChartsJSONRenderer) ignore(node *ast.Node) bool {
	return nil == node ||
		// 以下类型的节点不进行渲染，否则图画出来节点太多
		ast.NodeBlockquoteMarker == node.Type ||
		ast.NodeEmA6kOpenMarker == node.Type || ast.NodeEmA6kCloseMarker == node.Type ||
		ast.NodeEmU8eOpenMarker == node.Type || ast.NodeEmU8eCloseMarker == node.Type ||
		ast.NodeStrongA6kOpenMarker == node.Type || ast.NodeStrongA6kCloseMarker == node.Type ||
		ast.NodeStrongU8eOpenMarker == node.Type || ast.NodeStrongU8eCloseMarker == node.Type ||
		ast.NodeStrikethrough1OpenMarker == node.Type || ast.NodeStrikethrough1CloseMarker == node.Type ||
		ast.NodeStrikethrough2OpenMarker == node.Type || ast.NodeStrikethrough2CloseMarker == node.Type ||
		ast.NodeMathBlockOpenMarker == node.Type || ast.NodeMathBlockContent == node.Type || ast.NodeMathBlockCloseMarker == node.Type ||
		ast.NodeInlineMathOpenMarker == node.Type || ast.NodeInlineMathContent == node.Type || ast.NodeInlineMathCloseMarker == node.Type
}
