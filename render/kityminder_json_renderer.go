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
	"bytes"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// KityMinderJSONRenderer 描述了 KityMinder JSON 渲染器。
type KityMinderJSONRenderer struct {
	*BaseRenderer
}

// NewKityMinderJSONRenderer 创建一个 KityMinder JSON 渲染器。
func NewKityMinderJSONRenderer(tree *parse.Tree, options *Options) Renderer {
	ret := &KityMinderJSONRenderer{NewBaseRenderer(tree, options)}
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
	ret.RendererFuncs[ast.NodeMark] = ret.renderMark
	ret.RendererFuncs[ast.NodeSup] = ret.renderSup
	ret.RendererFuncs[ast.NodeSub] = ret.renderSub
	ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL
	ret.RendererFuncs[ast.NodeKramdownSpanIAL] = ret.renderKramdownSpanIAL
	ret.RendererFuncs[ast.NodeBlockEmbed] = ret.renderBlockEmbed
	ret.RendererFuncs[ast.NodeBlockQueryEmbed] = ret.renderBlockQueryEmbed
	ret.DefaultRendererFunc = ret.renderDefault
	return ret
}

func (r *KityMinderJSONRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderKramdownSpanIAL(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if nil == node.Previous {
			return ast.WalkContinue
		}
		id := r.NodeID(node.Previous)
		r.leaf("Span IAL\n{: "+id+"}", node)
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Mark\nmark", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Sup\nsup", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Sub\nsub", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("BlockQueryEmbed\n!{{script}}", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderBlockEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("BlockEmbed\n!((id))", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderBlockRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("BlockRef\n((id))", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderDefault(n *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Front Matter\nYAML", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("HTML Entity\nspan", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Blackslash\ndiv", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("ToC\ndiv", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Footnotes Ref\ndiv", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Footnotes Def\np", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Inline Math\nspan", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Math Block\ndiv", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Emoji Img\n", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Emoji Unicode\n", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Table Cell\ntd", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Table Row\ntr", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Table Head\nthead", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Table\ntable", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Strikethrough\ndel", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Image\nimg", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Link\na", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("HTML Block\n", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Inline HTML\n", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemOpenBrace)
		r.WriteString("\"root\":")
		r.openObj()
		r.dataText("文档名")
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.WriteByte(lex.ItemCloseBrace)
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		md := r.formatNode(node)
		r.dataText(md)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		text := util.BytesToStr(node.Tokens)
		r.openObj()
		r.dataText(text)
		r.closeObj()
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Code Span\ncode", node)
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Emphasis\nem", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Strong\nstrong", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("Blockquote\nblockquote", node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		h := "h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1]
		r.val("Heading\n"+h, node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
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
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val("List Item\nli "+util.BytesToStr(node.ListData.Marker), node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
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
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Thematic Break\nhr", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Hard Break\nbr", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Soft Break\n", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf("Code Block\npre.code", node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) leaf(val string, node *ast.Node) {
	r.openObj()
	r.val(val, node)
	r.closeObj()
}

func (r *KityMinderJSONRenderer) dataText(text string) {
	r.WriteString("\"data\":")
	r.openObj()
	r.WriteString("\"text\":\"" + text + "\"")
	r.closeObj()
}

func (r *KityMinderJSONRenderer) val(val string, node *ast.Node) {
	val = strings.ReplaceAll(val, "\\", "\\\\")
	val = strings.ReplaceAll(val, "\n", "\\n")
	val = strings.ReplaceAll(val, "\"", "")
	val = strings.ReplaceAll(val, "'", "")
	r.WriteString("\"data\":\"{text:\"" + val + "\"}}")
}

func (r *KityMinderJSONRenderer) openObj() {
	r.WriteByte('{')
}

func (r *KityMinderJSONRenderer) closeObj() {
	r.WriteByte('}')
}

func (r *KityMinderJSONRenderer) openChildren(node *ast.Node) {
	if nil != node.FirstChild {
		r.WriteString(",\"children\":[")
	}
}

func (r *KityMinderJSONRenderer) closeChildren(node *ast.Node) {
	if nil != node.FirstChild {
		r.WriteByte(']')
	}
}

func (r *KityMinderJSONRenderer) comma() {
	r.WriteString(",")
}

func (r *KityMinderJSONRenderer) ignore(node *ast.Node) bool {
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

func (r *KityMinderJSONRenderer) formatNode(node *ast.Node) string {
	renderer := NewFormatRenderer(r.Tree, r.Options)
	renderer.Writer = &bytes.Buffer{}
	renderer.NodeWriterStack = append(renderer.NodeWriterStack, renderer.Writer)
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		rendererFunc := renderer.RendererFuncs[n.Type]
		return rendererFunc(n, entering)
	})
	return strings.TrimSpace(renderer.Writer.String())
}
