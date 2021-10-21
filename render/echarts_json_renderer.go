// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
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

// EChartsJSONRenderer 描述了 ECharts JSON 渲染器。
type EChartsJSONRenderer struct {
	*BaseRenderer
}

// NewEChartsJSONRenderer 创建一个 ECharts JSON 渲染器。
func NewEChartsJSONRenderer(tree *parse.Tree, options *Options) Renderer {
	ret := &EChartsJSONRenderer{NewBaseRenderer(tree, options)}
	ret.RendererFuncs[ast.NodeDocument] = ret.renderDocument
	ret.RendererFuncs[ast.NodeParagraph] = ret.renderParagraph
	ret.RendererFuncs[ast.NodeText] = ret.renderText
	ret.RendererFuncs[ast.NodeCodeSpan] = ret.renderCodeSpan
	ret.RendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock
	ret.RendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock
	ret.RendererFuncs[ast.NodeInlineMath] = ret.renderInlineMath
	ret.RendererFuncs[ast.NodeEmphasis] = ret.renderEmphasis
	ret.RendererFuncs[ast.NodeStrong] = ret.renderStrong
	ret.RendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote
	ret.RendererFuncs[ast.NodeHeading] = ret.renderHeading
	ret.RendererFuncs[ast.NodeList] = ret.renderList
	ret.RendererFuncs[ast.NodeListItem] = ret.renderListItem
	ret.RendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	ret.RendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak
	ret.RendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak
	ret.RendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	ret.RendererFuncs[ast.NodeInlineHTML] = ret.renderInlineHTML
	ret.RendererFuncs[ast.NodeLink] = ret.renderLink
	ret.RendererFuncs[ast.NodeImage] = ret.renderImage
	ret.RendererFuncs[ast.NodeStrikethrough] = ret.renderStrikethrough
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
	ret.RendererFuncs[ast.NodeBlockRef] = ret.renderBlockRef
	ret.RendererFuncs[ast.NodeFileAnnotationRef] = ret.renderFileAnnotationRef
	ret.RendererFuncs[ast.NodeMark] = ret.renderMark
	ret.RendererFuncs[ast.NodeSup] = ret.renderSup
	ret.RendererFuncs[ast.NodeSub] = ret.renderSub
	ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL
	ret.RendererFuncs[ast.NodeKramdownSpanIAL] = ret.renderKramdownSpanIAL
	ret.RendererFuncs[ast.NodeBlockQueryEmbed] = ret.renderBlockQueryEmbed
	ret.DefaultRendererFunc = ret.renderDefault
	return ret
}

func (r *EChartsJSONRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if nil == node.Previous {
			return ast.WalkContinue
		}
		id := r.NodeID(node.Previous)
		if util.IsDocIAL(node.Tokens) {
			id = r.Tree.ID
		}
		r.leaf("Block IAL\n{: "+id+"}", node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderKramdownSpanIAL(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if nil == node.Previous {
			return ast.WalkContinue
		}
		id := r.NodeID(node.Previous)
		r.leaf("Span IAL\n{: "+id+"}", node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Mark\nmark", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Sup\nsup", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Sub\nsub", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("BlockQueryEmbed\n{{script}}", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderBlockRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("BlockRef\n((id))", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderFileAnnotationRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("FileAnnotationRef\n<<id>>", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderDefault(n *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Front Matter\nYAML", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("HTML Entity\nspan", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Blackslash\ndiv", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("ToC\ndiv", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Footnotes Ref\ndiv", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *EChartsJSONRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Inline Math\nspan", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Math Block\ndiv", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Emoji Img\n", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Emoji Unicode\n", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Table Cell\ntd", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Table Row\ntr", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Table Head\nthead", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Table\ntable", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *EChartsJSONRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Strikethrough\ndel", node)
	}
	return ast.WalkSkipChildren
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
	if entering {
		r.leaf("HTML Block\n", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Inline HTML\n", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemOpenBracket)
		r.openObj()
		r.val("Document", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
		r.WriteByte(lex.ItemCloseBracket)
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
	return ast.WalkSkipChildren
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
		h := "h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1]
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
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Hard Break\nbr", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Soft Break\n", node)
	}
	return ast.WalkSkipChildren
}

func (r *EChartsJSONRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Code Block\npre.code", node)
	}
	return ast.WalkSkipChildren
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
	r.WriteString("\"name\":\"" + val + "\"")
}

func (r *EChartsJSONRenderer) openObj() {
	r.WriteByte('{')
}

func (r *EChartsJSONRenderer) closeObj(node *ast.Node) {
	r.WriteByte('}')
	if !r.ignore(node.Next) {
		r.comma()
	}
}

func (r *EChartsJSONRenderer) openChildren(node *ast.Node) {
	if nil != node.FirstChild {
		r.WriteString(",\"children\":[")
	}
}

func (r *EChartsJSONRenderer) closeChildren(node *ast.Node) {
	if nil != node.FirstChild {
		r.WriteByte(']')
	}
}

func (r *EChartsJSONRenderer) comma() {
	r.WriteString(",")
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
		ast.NodeInlineMathOpenMarker == node.Type || ast.NodeInlineMathContent == node.Type || ast.NodeInlineMathCloseMarker == node.Type ||
		ast.NodeYamlFrontMatterOpenMarker == node.Type || ast.NodeYamlFrontMatterCloseMarker == node.Type || ast.NodeYamlFrontMatterContent == node.Type
}
