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
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// JSONRenderer 描述了 JSON 渲染器。
type JSONRenderer struct {
	*BaseRenderer
}

// NewJSONRenderer 创建一个 JSON 渲染器。
func NewJSONRenderer(tree *parse.Tree) Renderer {
	ret := &JSONRenderer{NewBaseRenderer(tree)}
	ret.RendererFuncs[ast.NodeDocument] = ret.renderDocument
	ret.RendererFuncs[ast.NodeParagraph] = ret.renderParagraph
	ret.RendererFuncs[ast.NodeText] = ret.renderText
	ret.RendererFuncs[ast.NodeCodeSpan] = ret.renderCodeSpan
	//ret.RendererFuncs[ast.NodeCodeSpanOpenMarker] = ret.renderCodeSpanOpenMarker
	//ret.RendererFuncs[ast.NodeCodeSpanContent] = ret.renderCodeSpanContent
	//ret.RendererFuncs[ast.NodeCodeSpanCloseMarker] = ret.renderCodeSpanCloseMarker
	ret.RendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock
	//ret.RendererFuncs[ast.NodeCodeBlockFenceOpenMarker] = ret.renderCodeBlockOpenMarker
	//ret.RendererFuncs[ast.NodeCodeBlockFenceInfoMarker] = ret.renderCodeBlockInfoMarker
	//ret.RendererFuncs[ast.NodeCodeBlockCode] = ret.renderCodeBlockCode
	//ret.RendererFuncs[ast.NodeCodeBlockFenceCloseMarker] = ret.renderCodeBlockCloseMarker
	ret.RendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock
	//ret.RendererFuncs[ast.NodeMathBlockOpenMarker] = ret.renderMathBlockOpenMarker
	//ret.RendererFuncs[ast.NodeMathBlockContent] = ret.renderMathBlockContent
	//ret.RendererFuncs[ast.NodeMathBlockCloseMarker] = ret.renderMathBlockCloseMarker
	ret.RendererFuncs[ast.NodeInlineMath] = ret.renderInlineMath
	//ret.RendererFuncs[ast.NodeInlineMathOpenMarker] = ret.renderInlineMathOpenMarker
	//ret.RendererFuncs[ast.NodeInlineMathContent] = ret.renderInlineMathContent
	//ret.RendererFuncs[ast.NodeInlineMathCloseMarker] = ret.renderInlineMathCloseMarker
	ret.RendererFuncs[ast.NodeEmphasis] = ret.renderEmphasis
	//ret.RendererFuncs[ast.NodeEmA6kOpenMarker] = ret.renderEmAsteriskOpenMarker
	//ret.RendererFuncs[ast.NodeEmA6kCloseMarker] = ret.renderEmAsteriskCloseMarker
	//ret.RendererFuncs[ast.NodeEmU8eOpenMarker] = ret.renderEmUnderscoreOpenMarker
	//ret.RendererFuncs[ast.NodeEmU8eCloseMarker] = ret.renderEmUnderscoreCloseMarker
	ret.RendererFuncs[ast.NodeStrong] = ret.renderStrong
	//ret.RendererFuncs[ast.NodeStrongA6kOpenMarker] = ret.renderStrongA6kOpenMarker
	//ret.RendererFuncs[ast.NodeStrongA6kCloseMarker] = ret.renderStrongA6kCloseMarker
	//ret.RendererFuncs[ast.NodeStrongU8eOpenMarker] = ret.renderStrongU8eOpenMarker
	//ret.RendererFuncs[ast.NodeStrongU8eCloseMarker] = ret.renderStrongU8eCloseMarker
	ret.RendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote
	//ret.RendererFuncs[ast.NodeBlockquoteMarker] = ret.renderBlockquoteMarker
	ret.RendererFuncs[ast.NodeHeading] = ret.renderHeading
	//ret.RendererFuncs[ast.NodeHeadingC8hMarker] = ret.renderHeadingC8hMarker
	//ret.RendererFuncs[ast.NodeHeadingID] = ret.renderHeadingID
	ret.RendererFuncs[ast.NodeList] = ret.renderList
	ret.RendererFuncs[ast.NodeListItem] = ret.renderListItem
	ret.RendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	ret.RendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak
	ret.RendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak
	ret.RendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	ret.RendererFuncs[ast.NodeInlineHTML] = ret.renderInlineHTML
	ret.RendererFuncs[ast.NodeLink] = ret.renderLink
	ret.RendererFuncs[ast.NodeImage] = ret.renderImage
	//ret.RendererFuncs[ast.NodeBang] = ret.renderBang
	//ret.RendererFuncs[ast.NodeOpenBracket] = ret.renderOpenBracket
	//ret.RendererFuncs[ast.NodeCloseBracket] = ret.renderCloseBracket
	//ret.RendererFuncs[ast.NodeOpenParen] = ret.renderOpenParen
	//ret.RendererFuncs[ast.NodeCloseParen] = ret.renderCloseParen
	//ret.RendererFuncs[ast.NodeLinkText] = ret.renderLinkText
	//ret.RendererFuncs[ast.NodeLinkSpace] = ret.renderLinkSpace
	//ret.RendererFuncs[ast.NodeLinkDest] = ret.renderLinkDest
	//ret.RendererFuncs[ast.NodeLinkTitle] = ret.renderLinkTitle
	ret.RendererFuncs[ast.NodeStrikethrough] = ret.renderStrikethrough
	//ret.RendererFuncs[ast.NodeStrikethrough1OpenMarker] = ret.renderStrikethrough1OpenMarker
	//ret.RendererFuncs[ast.NodeStrikethrough1CloseMarker] = ret.renderStrikethrough1CloseMarker
	//ret.RendererFuncs[ast.NodeStrikethrough2OpenMarker] = ret.renderStrikethrough2OpenMarker
	//ret.RendererFuncs[ast.NodeStrikethrough2CloseMarker] = ret.renderStrikethrough2CloseMarker
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
	ret.RendererFuncs[ast.NodeHTMLEntity] = ret.renderHtmlEntity
	ret.RendererFuncs[ast.NodeYamlFrontMatter] = ret.renderYamlFrontMatter
	//ret.RendererFuncs[ast.NodeYamlFrontMatterOpenMarker] = ret.renderYamlFrontMatterOpenMarker
	//ret.RendererFuncs[ast.NodeYamlFrontMatterContent] = ret.renderYamlFrontMatterContent
	//ret.RendererFuncs[ast.NodeYamlFrontMatterCloseMarker] = ret.renderYamlFrontMatterCloseMarker
	return ret
}

func (r *JSONRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Front Matter\nYAML", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("HTML Entity\nspan", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *JSONRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Blackslash\ndiv", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("ToC\ndiv", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Footnotes Ref\ndiv", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Footnotes Def\np", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Inline Math\nspan", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Math Block\ndiv", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Emoji Img\n", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Emoji Unicode\n", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkStop
}

func (r *JSONRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Table Cell\ntd", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Table Row\ntr", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Table Head\nthead", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Table\ntable", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Strikethrough\ndel", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("HTML Block\n", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Inline HTML\n", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemOpenBracket)
		r.openObj()
		r.val("", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
		r.WriteByte(lex.ItemCloseBracket)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		text := util.BytesToStr(node.Tokens)
		r.openObj()
		r.val(text, node)
		r.closeObj(node)
	}
	return ast.WalkStop
}

func (r *JSONRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Code Span\ncode", node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		h := "h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1]
		r.val("Heading\n"+h, node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Thematic Break\nhr", node)
	}
	return ast.WalkStop
}

func (r *JSONRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Hard Break\nbr", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Soft Break\n", node)
	return ast.WalkStop
}

func (r *JSONRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	r.leaf("Code Block\npre.code", node)
	return ast.WalkStop
}

func (r *JSONRenderer) leaf(val string, node *ast.Node) {
	r.openObj()
	r.val(val, node)
	r.closeObj(node)
}

func (r *JSONRenderer) val(val string, node *ast.Node) {
	typ := node.Type.String()
	typ = typ[len("Node"):]
	r.WriteString("\"type\":\"" + typ + "\"")
	if "" != val {
		r.WriteString(",")
	} else {
		return
	}
	val = strings.ReplaceAll(val, "\\", "\\\\")
	val = strings.ReplaceAll(val, "\n", "\\n")
	val = strings.ReplaceAll(val, "\"", "\\\"")
	val = strings.ReplaceAll(val, "'", "\\'")
	r.WriteString("\"val\":\"" + val + "\"")
}

func (r *JSONRenderer) openObj() {
	r.WriteByte('{')
}

func (r *JSONRenderer) closeObj(node *ast.Node) {
	r.WriteByte('}')
	if nil != node.Next {
		r.WriteString(",")
	}
}

func (r *JSONRenderer) openChildren(node *ast.Node) {
	if nil != node.FirstChild {
		r.WriteString(",\"children\":[")
	}
}

func (r *JSONRenderer) closeChildren(node *ast.Node) {
	if nil != node.FirstChild {
		r.WriteByte(']')
	}
}
