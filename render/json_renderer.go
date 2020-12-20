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
	"fmt"
	"github.com/88250/lute/lex"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// EChartsJSONRenderer 描述了 JSON 渲染器。
type JSONRenderer struct {
	*BaseRenderer
}

// NewEChartsJSONRenderer 创建一个 JSON 渲染器。
func NewJSONRenderer(tree *parse.Tree) Renderer {
	ret := &JSONRenderer{NewBaseRenderer(tree)}
	ret.RendererFuncs[ast.NodeDocument] = ret.renderDocument
	ret.RendererFuncs[ast.NodeParagraph] = ret.renderParagraph
	//ret.RendererFuncs[ast.NodeText] = ret.renderText
	//ret.RendererFuncs[ast.NodeCodeSpan] = ret.renderCodeSpan
	//ret.RendererFuncs[ast.NodeCodeSpanOpenMarker] = ret.renderCodeSpanOpenMarker
	//ret.RendererFuncs[ast.NodeCodeSpanContent] = ret.renderCodeSpanContent
	//ret.RendererFuncs[ast.NodeCodeSpanCloseMarker] = ret.renderCodeSpanCloseMarker
	//ret.RendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock
	//ret.RendererFuncs[ast.NodeCodeBlockFenceOpenMarker] = ret.renderCodeBlockOpenMarker
	//ret.RendererFuncs[ast.NodeCodeBlockFenceInfoMarker] = ret.renderCodeBlockInfoMarker
	//ret.RendererFuncs[ast.NodeCodeBlockCode] = ret.renderCodeBlockCode
	//ret.RendererFuncs[ast.NodeCodeBlockFenceCloseMarker] = ret.renderCodeBlockCloseMarker
	//ret.RendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock
	//ret.RendererFuncs[ast.NodeMathBlockOpenMarker] = ret.renderMathBlockOpenMarker
	//ret.RendererFuncs[ast.NodeMathBlockContent] = ret.renderMathBlockContent
	//ret.RendererFuncs[ast.NodeMathBlockCloseMarker] = ret.renderMathBlockCloseMarker
	//ret.RendererFuncs[ast.NodeInlineMath] = ret.renderInlineMath
	//ret.RendererFuncs[ast.NodeInlineMathOpenMarker] = ret.renderInlineMathOpenMarker
	//ret.RendererFuncs[ast.NodeInlineMathContent] = ret.renderInlineMathContent
	//ret.RendererFuncs[ast.NodeInlineMathCloseMarker] = ret.renderInlineMathCloseMarker
	//ret.RendererFuncs[ast.NodeEmphasis] = ret.renderEmphasis
	//ret.RendererFuncs[ast.NodeEmA6kOpenMarker] = ret.renderEmAsteriskOpenMarker
	//ret.RendererFuncs[ast.NodeEmA6kCloseMarker] = ret.renderEmAsteriskCloseMarker
	//ret.RendererFuncs[ast.NodeEmU8eOpenMarker] = ret.renderEmUnderscoreOpenMarker
	//ret.RendererFuncs[ast.NodeEmU8eCloseMarker] = ret.renderEmUnderscoreCloseMarker
	//ret.RendererFuncs[ast.NodeStrong] = ret.renderStrong
	//ret.RendererFuncs[ast.NodeStrongA6kOpenMarker] = ret.renderStrongA6kOpenMarker
	//ret.RendererFuncs[ast.NodeStrongA6kCloseMarker] = ret.renderStrongA6kCloseMarker
	//ret.RendererFuncs[ast.NodeStrongU8eOpenMarker] = ret.renderStrongU8eOpenMarker
	//ret.RendererFuncs[ast.NodeStrongU8eCloseMarker] = ret.renderStrongU8eCloseMarker
	//ret.RendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote
	//ret.RendererFuncs[ast.NodeBlockquoteMarker] = ret.renderBlockquoteMarker
	//ret.RendererFuncs[ast.NodeHeading] = ret.renderHeading
	//ret.RendererFuncs[ast.NodeHeadingC8hMarker] = ret.renderHeadingC8hMarker
	//ret.RendererFuncs[ast.NodeHeadingID] = ret.renderHeadingID
	//ret.RendererFuncs[ast.NodeList] = ret.renderList
	//ret.RendererFuncs[ast.NodeListItem] = ret.renderListItem
	//ret.RendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	//ret.RendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak
	//ret.RendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak
	//ret.RendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	//ret.RendererFuncs[ast.NodeInlineHTML] = ret.renderInlineHTML
	//ret.RendererFuncs[ast.NodeLink] = ret.renderLink
	//ret.RendererFuncs[ast.NodeImage] = ret.renderImage
	//ret.RendererFuncs[ast.NodeBang] = ret.renderBang
	//ret.RendererFuncs[ast.NodeOpenBracket] = ret.renderOpenBracket
	//ret.RendererFuncs[ast.NodeCloseBracket] = ret.renderCloseBracket
	//ret.RendererFuncs[ast.NodeOpenParen] = ret.renderOpenParen
	//ret.RendererFuncs[ast.NodeCloseParen] = ret.renderCloseParen
	//ret.RendererFuncs[ast.NodeOpenBrace] = ret.renderOpenBrace
	//ret.RendererFuncs[ast.NodeCloseBrace] = ret.renderCloseBrace
	//ret.RendererFuncs[ast.NodeLinkText] = ret.renderLinkText
	//ret.RendererFuncs[ast.NodeLinkSpace] = ret.renderLinkSpace
	//ret.RendererFuncs[ast.NodeLinkDest] = ret.renderLinkDest
	//ret.RendererFuncs[ast.NodeLinkTitle] = ret.renderLinkTitle
	//ret.RendererFuncs[ast.NodeStrikethrough] = ret.renderStrikethrough
	//ret.RendererFuncs[ast.NodeStrikethrough1OpenMarker] = ret.renderStrikethrough1OpenMarker
	//ret.RendererFuncs[ast.NodeStrikethrough1CloseMarker] = ret.renderStrikethrough1CloseMarker
	//ret.RendererFuncs[ast.NodeStrikethrough2OpenMarker] = ret.renderStrikethrough2OpenMarker
	//ret.RendererFuncs[ast.NodeStrikethrough2CloseMarker] = ret.renderStrikethrough2CloseMarker
	//ret.RendererFuncs[ast.NodeTaskListItemMarker] = ret.renderTaskListItemMarker
	//ret.RendererFuncs[ast.NodeTable] = ret.renderTable
	//ret.RendererFuncs[ast.NodeTableHead] = ret.renderTableHead
	//ret.RendererFuncs[ast.NodeTableRow] = ret.renderTableRow
	//ret.RendererFuncs[ast.NodeTableCell] = ret.renderTableCell
	//ret.RendererFuncs[ast.NodeEmoji] = ret.renderEmoji
	//ret.RendererFuncs[ast.NodeEmojiUnicode] = ret.renderEmojiUnicode
	//ret.RendererFuncs[ast.NodeEmojiImg] = ret.renderEmojiImg
	//ret.RendererFuncs[ast.NodeEmojiAlias] = ret.renderEmojiAlias
	//ret.RendererFuncs[ast.NodeFootnotesDefBlock] = ret.renderFootnotesDefBlock
	//ret.RendererFuncs[ast.NodeFootnotesDef] = ret.renderFootnotesDef
	//ret.RendererFuncs[ast.NodeFootnotesRef] = ret.renderFootnotesRef
	//ret.RendererFuncs[ast.NodeToC] = ret.renderToC
	//ret.RendererFuncs[ast.NodeBackslash] = ret.renderBackslash
	//ret.RendererFuncs[ast.NodeBackslashContent] = ret.renderBackslashContent
	//ret.RendererFuncs[ast.NodeHTMLEntity] = ret.renderHtmlEntity
	//ret.RendererFuncs[ast.NodeYamlFrontMatter] = ret.renderYamlFrontMatter
	//ret.RendererFuncs[ast.NodeYamlFrontMatterOpenMarker] = ret.renderYamlFrontMatterOpenMarker
	//ret.RendererFuncs[ast.NodeYamlFrontMatterContent] = ret.renderYamlFrontMatterContent
	//ret.RendererFuncs[ast.NodeYamlFrontMatterCloseMarker] = ret.renderYamlFrontMatterCloseMarker
	//ret.RendererFuncs[ast.NodeBlockRef] = ret.renderBlockRef
	//ret.RendererFuncs[ast.NodeBlockRefID] = ret.renderBlockRefID
	//ret.RendererFuncs[ast.NodeBlockRefSpace] = ret.renderBlockRefSpace
	//ret.RendererFuncs[ast.NodeBlockRefText] = ret.renderBlockRefText
	//ret.RendererFuncs[ast.NodeMark] = ret.renderMark
	//ret.RendererFuncs[ast.NodeMark1OpenMarker] = ret.renderMark1OpenMarker
	//ret.RendererFuncs[ast.NodeMark1CloseMarker] = ret.renderMark1CloseMarker
	//ret.RendererFuncs[ast.NodeMark2OpenMarker] = ret.renderMark2OpenMarker
	//ret.RendererFuncs[ast.NodeMark2CloseMarker] = ret.renderMark2CloseMarker
	//ret.RendererFuncs[ast.NodeSup] = ret.renderSup
	//ret.RendererFuncs[ast.NodeSupOpenMarker] = ret.renderSupOpenMarker
	//ret.RendererFuncs[ast.NodeSupCloseMarker] = ret.renderSupCloseMarker
	//ret.RendererFuncs[ast.NodeSub] = ret.renderSub
	//ret.RendererFuncs[ast.NodeSubOpenMarker] = ret.renderSubOpenMarker
	//ret.RendererFuncs[ast.NodeSubCloseMarker] = ret.renderSubCloseMarker
	//ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL
	//ret.RendererFuncs[ast.NodeKramdownSpanIAL] = ret.renderKramdownSpanIAL
	//ret.RendererFuncs[ast.NodeBlockQueryEmbed] = ret.renderBlockQueryEmbed
	//ret.RendererFuncs[ast.NodeBlockQueryEmbedScript] = ret.renderBlockQueryEmbedScript
	//ret.RendererFuncs[ast.NodeBlockEmbed] = ret.renderBlockEmbed
	//ret.RendererFuncs[ast.NodeBlockEmbedID] = ret.renderBlockEmbedID
	//ret.RendererFuncs[ast.NodeBlockEmbedSpace] = ret.renderBlockEmbedSpace
	//ret.RendererFuncs[ast.NodeBlockEmbedText] = ret.renderBlockEmbedText
	//ret.RendererFuncs[ast.NodeTag] = ret.renderTag
	//ret.RendererFuncs[ast.NodeTagOpenMarker] = ret.renderTagOpenMarker
	//ret.RendererFuncs[ast.NodeTagCloseMarker] = ret.renderTagCloseMarker
	//ret.RendererFuncs[ast.NodeLinkRefDefBlock] = ret.renderLinkRefDefBlock
	//ret.RendererFuncs[ast.NodeLinkRefDef] = ret.renderLinkRefDef
	//ret.RendererFuncs[ast.NodeSuperBlock] = ret.renderSuperBlock
	//ret.RendererFuncs[ast.NodeSuperBlockOpenMarker] = ret.renderSuperBlockOpenMarker
	//ret.RendererFuncs[ast.NodeSuperBlockLayoutMarker] = ret.renderSuperBlockLayoutMarker
	//ret.RendererFuncs[ast.NodeSuperBlockCloseMarker] = ret.renderSuperBlockCloseMarker

	//ret.RendererFuncs[ast.NodeParagraph] = ret.renderParagraph  // 已测试
	//ret.RendererFuncs[ast.NodeText] = ret.renderText  // TODO 需要酌情清除
	//ret.RendererFuncs[ast.NodeCodeSpanContent] = ret.renderCodeSpanContent  // 解析行内代码块，已测试
	//ret.RendererFuncs[ast.NodeCodeSpan] = ret.renderCodeSpan  // TODO 移除
	////ret.RendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock  // TODO 处理代码块
	////ret.RendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock  // TODO 移除
	//ret.RendererFuncs[ast.NodeMathBlockContent] = ret.renderMathBlockContent  // 解析数学块公式，已测试
	////ret.RendererFuncs[ast.NodeInlineMath] = ret.renderInlineMath  // TODO 移除
	//ret.RendererFuncs[ast.NodeInlineMathContent] = ret.renderInlineMathContent  // 解析行内数学公式，已测试
	//ret.RendererFuncs[ast.NodeEmphasis] = ret.renderEmphasis // 已测试
	//ret.RendererFuncs[ast.NodeStrong] = ret.renderStrong // 已测试
	//ret.RendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote // 已测试
	//ret.RendererFuncs[ast.NodeHeading] = ret.renderHeading  // 已测试 TODO 指定ID
	//ret.RendererFuncs[ast.NodeList] = ret.renderList  // 已测试
	//ret.RendererFuncs[ast.NodeListItem] = ret.renderListItem  // 已测试
	//ret.RendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak  // 已测试
	//ret.RendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak  // TODO 与InlineHTML冲突
	//ret.RendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak  // 已测试
	//ret.RendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML  // TODO HTML块
	//ret.RendererFuncs[ast.NodeInlineHTML] = ret.renderInlineHTML  // TODO 已测试
	//ret.RendererFuncs[ast.NodeLink] = ret.renderLink  // TODO 渲染链接
	//ret.RendererFuncs[ast.NodeImage] = ret.renderImage  // TODO 渲染图片
	//ret.RendererFuncs[ast.NodeStrikethrough] = ret.renderStrikethrough  // 已测试
	//ret.RendererFuncs[ast.NodeTaskListItemMarker] = ret.renderTaskListItemMarker  // 已测试
	//ret.RendererFuncs[ast.NodeTable] = ret.renderTable  // 已测试
	//ret.RendererFuncs[ast.NodeTableHead] = ret.renderTableHead // 已测试
	//ret.RendererFuncs[ast.NodeTableRow] = ret.renderTableRow // 已测试
	//ret.RendererFuncs[ast.NodeTableCell] = ret.renderTableCell  // 已测试
	//ret.RendererFuncs[ast.NodeEmoji] = ret.renderEmoji  // 已测试
	//ret.RendererFuncs[ast.NodeEmojiUnicode] = ret.renderEmojiUnicode  // 已测试
	//ret.RendererFuncs[ast.NodeEmojiImg] = ret.renderEmojiImg  // 已测试
	//ret.RendererFuncs[ast.NodeEmojiAlias] = ret.renderEmojiAlias  // 已测试
	//ret.RendererFuncs[ast.NodeFootnotesDef] = ret.renderFootnotesDef  // TODO
	//ret.RendererFuncs[ast.NodeFootnotesRef] = ret.renderFootnotesRef  // TODO
	//ret.RendererFuncs[ast.NodeToC] = ret.renderToC  // TODO 渲染TOC
	//ret.RendererFuncs[ast.NodeBackslash] = ret.renderBackslash  // TODO 处理中括号，触发条件不明
	//ret.RendererFuncs[ast.NodeBackslashContent] = ret.renderBackslashContent  // TODO 处理中括号，触发条件不明
	//ret.RendererFuncs[ast.NodeHTMLEntity] = ret.renderHtmlEntity  // 已测试
	//ret.RendererFuncs[ast.NodeYamlFrontMatter] = ret.renderYamlFrontMatter  // TODO 处理yaml
	//ret.RendererFuncs[ast.NodeBlockRef] = ret.renderBlockRef
	//ret.RendererFuncs[ast.NodeMark] = ret.renderMark  // TODO 处理高亮
	//ret.RendererFuncs[ast.NodeSup] = ret.renderSup  // TODO 处理上下角标
	//ret.RendererFuncs[ast.NodeSub] = ret.renderSub  // TODO
	//ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL  // TODO 处理Kramdown
	//ret.RendererFuncs[ast.NodeKramdownSpanIAL] = ret.renderKramdownSpanIAL  // TODO 处理Kramdown
	//ret.RendererFuncs[ast.NodeBlockEmbed] = ret.renderBlockEmbed  // TODO 内容块嵌入
	//ret.RendererFuncs[ast.NodeBlockQueryEmbed] = ret.renderBlockQueryEmbed  // TODO 内容块引用
	//ret.DefaultRendererFunc = ret.renderDefault
	return ret
}

// TODO AST树核心设计概念
// 1. 忽略根节点渲染
// 2. 对于实际只起控制渲染效果的设为flag
// 3. 妥善处理Marker对于树构建过程中的影响
// TODO 4. Flag注释
// Paragraph 新的段

// 处理根节点
func (r *JSONRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemOpenBracket)
	} else {
		r.WriteByte(lex.ItemCloseBracket)
	}
	return ast.WalkContinue
}

// TODO 需要检验KramdownBlockIAL
func (r *JSONRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if nil == node.Previous {
			return ast.WalkContinue
		}
		id := r.NodeID(node.Previous)
		if util.IsDocIAL(node.Tokens) {
			id = r.Tree.ID
		}
		r.leaf(node.Type, "Block IAL\n{: "+id+"}", node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderKramdownSpanIAL(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if nil == node.Previous {
			return ast.WalkContinue
		}
		id := r.NodeID(node.Previous)
		r.leaf(node.Type,"Span IAL\n{: "+id+"}", node)
	}
	return ast.WalkContinue
}

// TODO Mark是指？
func (r *JSONRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		fmt.Println(node.Text())
		r.leaf(node.Type, "Mark\nmark", node)
	}
	return ast.WalkSkipChildren
}

// TODO 处理上下角标
func (r *JSONRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, "Sup\nsup", node)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, "Sub\nsub", node)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type,"BlockQueryEmbed\n!{{script}}", node)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderBlockEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type,"BlockEmbed\n!((id))", node)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderBlockRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type,"BlockRef\n((id))", node)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderDefault(n *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, "Front Matter\nYAML", node)
	}
	return ast.WalkSkipChildren
}

// HTML实体符号处理
func (r *JSONRenderer) renderHtmlEntity(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(node.Type, util.BytesToStr(node.Tokens))
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, "Blackslash\ndiv", node)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, "ToC\ndiv", node)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, "Footnotes Ref\ndiv", node)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(node.Type, node.Text())
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// 解析行内数学公式
func (r *JSONRenderer) renderInlineMathContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(ast.NodeInlineMath, util.BytesToStr(node.Tokens))
	} else {
		r.closeObj(node)
	}
	return ast.WalkSkipChildren
}

// 解析数学块
func (r *JSONRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(ast.NodeMathBlock,util.BytesToStr(node.Tokens))
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderEmojiImg(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, util.BytesToStr(node.Tokens), node)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderEmojiUnicode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, util.BytesToStr(node.Tokens), node)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderEmojiAlias(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderEmoji(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

// TODO 处理表格渲染
// TODO Cell渲染函数才能拿到对齐
func (r *JSONRenderer) renderTableCell(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var align string
		switch node.TableCellAlign {
		case 1:
			align = "left"
		case 2:
			align = "center"
		case 3:
			align = "right"
		default:
			align = "left"
		}
		r.openObj()
		r.val(node.Type, align)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// 表格新行标识
func (r *JSONRenderer) renderTableRow(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(node.Type, "")
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// 表头标识
func (r *JSONRenderer) renderTableHead(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(node.Type, "")
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// 表格标识
func (r *JSONRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(node.Type, "")
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// 删除线标识
func (r *JSONRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(node.Type, "")
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	// TODO 处理图片渲染
	//if entering {
	//	if 0 == r.DisableTags {
	//		r.WriteString("<img src=\"")
	//		destTokens := node.ChildByType(ast.NodeLinkDest).Tokens
	//		destTokens = r.Tree.Context.LinkPath(destTokens)
	//		if "" != r.Option.ImageLazyLoading {
	//			r.Write(html.EscapeHTML(util.StrToBytes(r.Option.ImageLazyLoading)))
	//			r.WriteString("\" data-src=\"")
	//		}
	//		r.Write(html.EscapeHTML(destTokens))
	//		r.WriteString("\" alt=\"")
	//	}
	//	r.DisableTags++
	//	return ast.WalkContinue
	//}
	//
	//r.DisableTags--
	//if 0 == r.DisableTags {
	//	r.WriteByte(lex.ItemDoublequote)
	//	if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
	//		r.WriteString(" title=\"")
	//		r.Write(html.EscapeHTML(title.Tokens))
	//		r.WriteByte(lex.ItemDoublequote)
	//	}
	//	ial := r.NodeAttrsStr(node)
	//	if "" != ial {
	//		r.WriteString(" " + ial)
	//	}
	//	r.WriteString(" />")
	//
	//	if r.Option.Sanitize {
	//		buf := r.Writer.Bytes()
	//		idx := bytes.LastIndex(buf, []byte("<img src="))
	//		imgBuf := buf[idx:]
	//		if r.Option.Sanitize {
	//			imgBuf = sanitize(imgBuf)
	//		}
	//		r.Writer.Truncate(idx)
	//		r.Writer.Write(imgBuf)
	//	}
	//}
	//return ast.WalkContinue
	if entering {
		r.openObj()
		r.val(node.Type, node.Text())
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
    // TODO 处理链接渲染
	//if entering {
	//	r.LinkTextAutoSpacePrevious(node)
	//
	//	dest := node.ChildByType(ast.NodeLinkDest)
	//	destTokens := dest.Tokens
	//	destTokens = r.Tree.Context.LinkPath(destTokens)
	//	attrs := [][]string{{"href", util.BytesToStr(html.EscapeHTML(destTokens))}}
	//	if title := node.ChildByType(ast.NodeLinkTitle); nil != title && nil != title.Tokens {
	//		attrs = append(attrs, []string{"title", util.BytesToStr(html.EscapeHTML(title.Tokens))})
	//	}
	//	r.Tag("a", attrs, false)
	//} else {
	//	r.Tag("/a", nil, false)
	//
	//	r.LinkTextAutoSpaceNext(node)
	//}
	//return ast.WalkContinue
	if entering {
		r.openObj()
		r.val(node.Type, node.Text())
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, "HTML Block\n", node)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		tokens := node.Tokens
		if r.Option.Sanitize {
			tokens = sanitize(tokens)
		}
		r.val(node.Type, util.BytesToStr(tokens))
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.flag(node)
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
		r.text(text)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.flag(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// 渲染行内代码
func (r *JSONRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(node.Type, util.BytesToStr(node.Tokens))
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// 斜体标志存在
func (r *JSONRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(node.Type, "")
		//r.flag(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// 加粗标志存在
func (r *JSONRenderer) renderStrong(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.flag(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// 引用标识存在
func (r *JSONRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(node.Type, "")
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// 标题存在，返回标题等级
func (r *JSONRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		h := "h" + headingLevel[node.HeadingLevel:node.HeadingLevel+1]
		// 标题值为标题等级
		r.val(node.Type, h)
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
		// 值为列表的类型
		r.val(node.Type, list)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// TODO 修复可能的BUG
func (r *JSONRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(node.Type, node.Text())
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
		check := "false"
		if node.TaskListItemChecked {
			check = "true"
		}
		// 返回勾选情况
		r.val(node.Type, check)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// 分割线
func (r *JSONRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, "hr", node)
	}
	return ast.WalkSkipChildren
}

// TODO 硬换行，与InlineHTML冲突
func (r *JSONRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, "br", node)
	}
	return ast.WalkSkipChildren
}

// \n换行
func (r *JSONRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, "\n", node)
	}
	return ast.WalkSkipChildren
}

// TODO 重新考虑处理代码内容和语言类型
//func (r *JSONRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
//	if entering {
//		//r.leaf(node.Type, "Code Block\npre.code", node)
//		r.openObj()
//		fmt.Println(node.Children)
//		r.val(node.Type, util.BytesToStr(node.FirstChild.Tokens), node)
//	} else {
//		r.closeObj(node)
//	}
//	return ast.WalkContinue
//}

func (r *JSONRenderer) leaf(nodeType ast.NodeType, val string, node *ast.Node) {
	r.openObj()
	r.val(nodeType, val)
	r.closeObj(node)
}

func (r *JSONRenderer) val(nodeType ast.NodeType , val string) {
	val = strings.ReplaceAll(val, "\\", "\\\\")
	val = strings.ReplaceAll(val, "\n", "\\n")
	val = strings.ReplaceAll(val, "\"", "")
	val = strings.ReplaceAll(val, "'", "")
	// 要写入Type和Value
	r.WriteString("\"type\":\"" + nodeType.String()[4:] + "\"" + ",")
	r.WriteString("\"value\":\"" + val + "\"")
}

// 写入如果节点是Text
func (r *JSONRenderer) text(val string) {
	r.WriteString(",\"text\":\"" + val + "\"")
}

func (r *JSONRenderer) flag(node *ast.Node) {
	//if !checkIfFlag(node.Parent.Type) {
	//	// 父节点不是flag
	//	r.openObj()
	//	r.WriteString("\"flags\":\"" + node.Type.String()[4:])
	//} else {
	//	//	父节点是flag
	//	r.WriteString(node.Type.String()[4:])
	//}
	//// 判断是不是Marker
	////	不需要闭合的条件 --> 有一个子节点，并且子节点还是flag类型
	////fmt.Println(node.FirstChild)
	//if len(node.Children) == 1 && checkIfFlag(node.FirstChild.Type) {
	//	r.WriteString("|")
	//} else {
	//	r.WriteString("\"")
	//}
	r.WriteString("\"flag\":\"" + node.Type.String()[4:]+"\"")
}

// 打开对象
func (r *JSONRenderer) openObj() {
	r.WriteByte('{')
}

// 关闭对象
func (r *JSONRenderer) closeObj(node *ast.Node) {
	r.WriteByte('}')
	r.comma()
}

// 打开子节点
func (r *JSONRenderer) openChildren(node *ast.Node) {
	if nil != node.FirstChild {
		r.WriteString(",\"children\":[")
	}
}

// 关闭子节点
func (r *JSONRenderer) closeChildren(node *ast.Node) {
	if nil != node.FirstChild {
		r.WriteByte(']')
	}
}

// 逗号
func (r *JSONRenderer) comma() {
	r.WriteString(",")
}