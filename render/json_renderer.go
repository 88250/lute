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

	"github.com/88250/lute/html"
	"github.com/88250/lute/lex"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// JSONRenderer 描述了 JSON 渲染器。
type JSONRenderer struct {
	*BaseRenderer
}

// NewJSONRenderer 创建一个 JSON 渲染器。
func NewJSONRenderer(tree *parse.Tree, options *Options) Renderer {
	ret := &JSONRenderer{NewBaseRenderer(tree, options)}
	ret.RendererFuncs[ast.NodeDocument] = ret.renderDocument
	ret.RendererFuncs[ast.NodeParagraph] = ret.renderParagraph
	ret.RendererFuncs[ast.NodeText] = ret.renderText
	ret.RendererFuncs[ast.NodeCodeSpan] = ret.renderCodeSpan
	ret.RendererFuncs[ast.NodeCodeSpanOpenMarker] = ret.renderCodeSpanOpenMarker
	ret.RendererFuncs[ast.NodeCodeSpanContent] = ret.renderCodeSpanContent
	ret.RendererFuncs[ast.NodeCodeSpanCloseMarker] = ret.renderCodeSpanCloseMarker
	ret.RendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock
	ret.RendererFuncs[ast.NodeCodeBlockFenceOpenMarker] = ret.renderCodeBlockOpenMarker
	ret.RendererFuncs[ast.NodeCodeBlockFenceInfoMarker] = ret.renderCodeBlockInfoMarker
	ret.RendererFuncs[ast.NodeCodeBlockCode] = ret.renderCodeBlockCode
	ret.RendererFuncs[ast.NodeCodeBlockFenceCloseMarker] = ret.renderCodeBlockCloseMarker
	ret.RendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock
	ret.RendererFuncs[ast.NodeMathBlockOpenMarker] = ret.renderMathBlockOpenMarker
	ret.RendererFuncs[ast.NodeMathBlockContent] = ret.renderMathBlockContent
	ret.RendererFuncs[ast.NodeMathBlockCloseMarker] = ret.renderMathBlockCloseMarker
	ret.RendererFuncs[ast.NodeInlineMath] = ret.renderInlineMath
	ret.RendererFuncs[ast.NodeInlineMathOpenMarker] = ret.renderInlineMathOpenMarker
	ret.RendererFuncs[ast.NodeInlineMathContent] = ret.renderInlineMathContent
	ret.RendererFuncs[ast.NodeInlineMathCloseMarker] = ret.renderInlineMathCloseMarker
	ret.RendererFuncs[ast.NodeEmphasis] = ret.renderEmphasis
	ret.RendererFuncs[ast.NodeEmA6kOpenMarker] = ret.renderEmAsteriskOpenMarker
	ret.RendererFuncs[ast.NodeEmA6kCloseMarker] = ret.renderEmAsteriskCloseMarker
	ret.RendererFuncs[ast.NodeEmU8eOpenMarker] = ret.renderEmUnderscoreOpenMarker
	ret.RendererFuncs[ast.NodeEmU8eCloseMarker] = ret.renderEmUnderscoreCloseMarker
	ret.RendererFuncs[ast.NodeStrong] = ret.renderStrong
	ret.RendererFuncs[ast.NodeStrongA6kOpenMarker] = ret.renderStrongA6kOpenMarker
	ret.RendererFuncs[ast.NodeStrongA6kCloseMarker] = ret.renderStrongA6kCloseMarker
	ret.RendererFuncs[ast.NodeStrongU8eOpenMarker] = ret.renderStrongU8eOpenMarker
	ret.RendererFuncs[ast.NodeStrongU8eCloseMarker] = ret.renderStrongU8eCloseMarker
	ret.RendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote
	ret.RendererFuncs[ast.NodeBlockquoteMarker] = ret.renderBlockquoteMarker
	ret.RendererFuncs[ast.NodeHeading] = ret.renderHeading
	ret.RendererFuncs[ast.NodeHeadingC8hMarker] = ret.renderHeadingC8hMarker
	ret.RendererFuncs[ast.NodeHeadingID] = ret.renderHeadingID
	ret.RendererFuncs[ast.NodeList] = ret.renderList
	ret.RendererFuncs[ast.NodeListItem] = ret.renderListItem
	ret.RendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	ret.RendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak
	ret.RendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak
	ret.RendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	ret.RendererFuncs[ast.NodeInlineHTML] = ret.renderInlineHTML
	ret.RendererFuncs[ast.NodeLink] = ret.renderLink
	ret.RendererFuncs[ast.NodeImage] = ret.renderImage
	ret.RendererFuncs[ast.NodeBang] = ret.renderBang
	ret.RendererFuncs[ast.NodeOpenBracket] = ret.renderOpenBracket
	ret.RendererFuncs[ast.NodeCloseBracket] = ret.renderCloseBracket
	ret.RendererFuncs[ast.NodeOpenParen] = ret.renderOpenParen
	ret.RendererFuncs[ast.NodeCloseParen] = ret.renderCloseParen
	ret.RendererFuncs[ast.NodeLess] = ret.renderLess
	ret.RendererFuncs[ast.NodeGreater] = ret.renderGreater
	ret.RendererFuncs[ast.NodeOpenBrace] = ret.renderOpenBrace
	ret.RendererFuncs[ast.NodeCloseBrace] = ret.renderCloseBrace
	ret.RendererFuncs[ast.NodeLinkText] = ret.renderLinkText
	ret.RendererFuncs[ast.NodeLinkSpace] = ret.renderLinkSpace
	ret.RendererFuncs[ast.NodeLinkDest] = ret.renderLinkDest
	ret.RendererFuncs[ast.NodeLinkTitle] = ret.renderLinkTitle
	ret.RendererFuncs[ast.NodeStrikethrough] = ret.renderStrikethrough
	ret.RendererFuncs[ast.NodeStrikethrough1OpenMarker] = ret.renderStrikethrough1OpenMarker
	ret.RendererFuncs[ast.NodeStrikethrough1CloseMarker] = ret.renderStrikethrough1CloseMarker
	ret.RendererFuncs[ast.NodeStrikethrough2OpenMarker] = ret.renderStrikethrough2OpenMarker
	ret.RendererFuncs[ast.NodeStrikethrough2CloseMarker] = ret.renderStrikethrough2CloseMarker
	ret.RendererFuncs[ast.NodeTaskListItemMarker] = ret.renderTaskListItemMarker
	ret.RendererFuncs[ast.NodeTable] = ret.renderTable
	ret.RendererFuncs[ast.NodeTableHead] = ret.renderTableHead
	ret.RendererFuncs[ast.NodeTableRow] = ret.renderTableRow
	ret.RendererFuncs[ast.NodeTableCell] = ret.renderTableCell
	ret.RendererFuncs[ast.NodeEmoji] = ret.renderEmoji
	ret.RendererFuncs[ast.NodeEmojiUnicode] = ret.renderEmojiUnicode
	ret.RendererFuncs[ast.NodeEmojiImg] = ret.renderEmojiImg
	ret.RendererFuncs[ast.NodeEmojiAlias] = ret.renderEmojiAlias
	ret.RendererFuncs[ast.NodeFootnotesDefBlock] = ret.renderFootnotesDefBlock
	ret.RendererFuncs[ast.NodeFootnotesDef] = ret.renderFootnotesDef
	ret.RendererFuncs[ast.NodeFootnotesRef] = ret.renderFootnotesRef
	ret.RendererFuncs[ast.NodeToC] = ret.renderToC
	ret.RendererFuncs[ast.NodeBackslash] = ret.renderBackslash
	ret.RendererFuncs[ast.NodeBackslashContent] = ret.renderBackslashContent
	ret.RendererFuncs[ast.NodeHTMLEntity] = ret.renderHtmlEntity
	ret.RendererFuncs[ast.NodeYamlFrontMatter] = ret.renderYamlFrontMatter
	ret.RendererFuncs[ast.NodeYamlFrontMatterOpenMarker] = ret.renderYamlFrontMatterOpenMarker
	ret.RendererFuncs[ast.NodeYamlFrontMatterContent] = ret.renderYamlFrontMatterContent
	ret.RendererFuncs[ast.NodeYamlFrontMatterCloseMarker] = ret.renderYamlFrontMatterCloseMarker
	ret.RendererFuncs[ast.NodeBlockRef] = ret.renderBlockRef
	ret.RendererFuncs[ast.NodeBlockRefID] = ret.renderBlockRefID
	ret.RendererFuncs[ast.NodeBlockRefSpace] = ret.renderBlockRefSpace
	ret.RendererFuncs[ast.NodeBlockRefText] = ret.renderBlockRefText
	ret.RendererFuncs[ast.NodeBlockRefDynamicText] = ret.renderBlockRefDynamicText
	ret.RendererFuncs[ast.NodeFileAnnotationRef] = ret.renderFileAnnotationRef
	ret.RendererFuncs[ast.NodeFileAnnotationRefID] = ret.renderFileAnnotationRefID
	ret.RendererFuncs[ast.NodeFileAnnotationRefSpace] = ret.renderFileAnnotationRefSpace
	ret.RendererFuncs[ast.NodeFileAnnotationRefText] = ret.renderFileAnnotationRefText
	ret.RendererFuncs[ast.NodeMark] = ret.renderMark
	ret.RendererFuncs[ast.NodeMark1OpenMarker] = ret.renderMark1OpenMarker
	ret.RendererFuncs[ast.NodeMark1CloseMarker] = ret.renderMark1CloseMarker
	ret.RendererFuncs[ast.NodeMark2OpenMarker] = ret.renderMark2OpenMarker
	ret.RendererFuncs[ast.NodeMark2CloseMarker] = ret.renderMark2CloseMarker
	ret.RendererFuncs[ast.NodeSup] = ret.renderSup
	ret.RendererFuncs[ast.NodeSupOpenMarker] = ret.renderSupOpenMarker
	ret.RendererFuncs[ast.NodeSupCloseMarker] = ret.renderSupCloseMarker
	ret.RendererFuncs[ast.NodeSub] = ret.renderSub
	ret.RendererFuncs[ast.NodeSubOpenMarker] = ret.renderSubOpenMarker
	ret.RendererFuncs[ast.NodeSubCloseMarker] = ret.renderSubCloseMarker
	ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL
	ret.RendererFuncs[ast.NodeKramdownSpanIAL] = ret.renderKramdownSpanIAL
	ret.RendererFuncs[ast.NodeBlockQueryEmbed] = ret.renderBlockQueryEmbed
	ret.RendererFuncs[ast.NodeBlockQueryEmbedScript] = ret.renderBlockQueryEmbedScript
	ret.RendererFuncs[ast.NodeTag] = ret.renderTag
	ret.RendererFuncs[ast.NodeTagOpenMarker] = ret.renderTagOpenMarker
	ret.RendererFuncs[ast.NodeTagCloseMarker] = ret.renderTagCloseMarker
	ret.RendererFuncs[ast.NodeLinkRefDefBlock] = ret.renderLinkRefDefBlock
	ret.RendererFuncs[ast.NodeLinkRefDef] = ret.renderLinkRefDef
	ret.RendererFuncs[ast.NodeSuperBlock] = ret.renderSuperBlock
	ret.RendererFuncs[ast.NodeSuperBlockOpenMarker] = ret.renderSuperBlockOpenMarker
	ret.RendererFuncs[ast.NodeSuperBlockLayoutMarker] = ret.renderSuperBlockLayoutMarker
	ret.RendererFuncs[ast.NodeSuperBlockCloseMarker] = ret.renderSuperBlockCloseMarker
	ret.DefaultRendererFunc = ret.renderDefault
	return ret
}

// AST树核心设计概念
// 1. 忽略根节点渲染
// 2. 对于实际只起控制渲染效果的设为flag
// 3. 妥善处理Marker对于树构建过程中的影响
// 4. 不支持Kramdown

// ==> Flag注释
// Paragraph 新的段
// Emphasis 斜体
// Strong 加粗
// Blockquote 引用块
// ListItem 列表项
// Strikethrough 删除线
// TableHead 表头
// Table 表格
// TableRow 表行
// Mark 高亮
// Sub 下标
// Sup 上标
// Tag 标记
// BlockRef 块引用

// OpenMarker处理
func isOpenMarker(node *ast.Node) bool {
	switch node.Type {
	case ast.NodeCodeBlockFenceOpenMarker,
		ast.NodeEmA6kOpenMarker,
		ast.NodeEmU8eOpenMarker,
		ast.NodeStrongA6kOpenMarker,
		ast.NodeStrongU8eOpenMarker,
		ast.NodeCodeSpanOpenMarker,
		ast.NodeStrikethrough1OpenMarker,
		ast.NodeStrikethrough2OpenMarker,
		ast.NodeMathBlockOpenMarker,
		ast.NodeInlineMathOpenMarker,
		ast.NodeYamlFrontMatterOpenMarker,
		ast.NodeMark1OpenMarker,
		ast.NodeMark2OpenMarker,
		ast.NodeTagOpenMarker,
		ast.NodeSupOpenMarker,
		ast.NodeSubOpenMarker:
		return true
	}
	return false
}

// CloseMarker处理
func isCloseMarker(node *ast.Node) bool {
	switch node.Type {
	case ast.NodeCodeBlockFenceCloseMarker,
		ast.NodeEmA6kCloseMarker,
		ast.NodeEmU8eCloseMarker,
		ast.NodeStrongA6kCloseMarker,
		ast.NodeStrongU8eCloseMarker,
		ast.NodeCodeSpanCloseMarker,
		ast.NodeStrikethrough1CloseMarker,
		ast.NodeStrikethrough2CloseMarker,
		ast.NodeMathBlockCloseMarker,
		ast.NodeInlineMathCloseMarker,
		ast.NodeYamlFrontMatterCloseMarker,
		ast.NodeMark1CloseMarker,
		ast.NodeMark2CloseMarker,
		ast.NodeTagCloseMarker,
		ast.NodeSupCloseMarker,
		ast.NodeSubCloseMarker:
		return true
	}
	return false
}

// 是否是标识
func isFlag(node *ast.Node) bool {
	switch node.Type {
	case ast.NodeParagraph,
		ast.NodeEmphasis,
		ast.NodeStrong,
		ast.NodeBlockquote,
		ast.NodeListItem,
		ast.NodeStrikethrough,
		ast.NodeTableHead,
		ast.NodeTable,
		ast.NodeTableRow,
		ast.NodeMark,
		ast.NodeSub,
		ast.NodeSup,
		ast.NodeTag,
		ast.NodeBlockRef,
		ast.NodeFileAnnotationRef:
		return true
	}
	return false
}

// 处理根节点
func (r *JSONRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemOpenBracket)
	} else {
		r.WriteByte(lex.ItemCloseBracket)
	}
	return ast.WalkContinue
}

// 处理段落 - 标识
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

// 处理文本
func (r *JSONRenderer) renderText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(node.Type, util.BytesToStr(node.Tokens))
	} else {
		r.closeObj(node)
	}
	return ast.WalkSkipChildren
}

// 处理行内代码根节点
func (r *JSONRenderer) renderCodeSpan(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

// TODO 处理renderCodeSpanOpenMarker
func (r *JSONRenderer) renderCodeSpanOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

// 渲染行内代码
func (r *JSONRenderer) renderCodeSpanContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(ast.NodeCodeSpan, util.BytesToStr(node.Tokens))
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// TODO 处理renderCodeSpanCloseMarker
func (r *JSONRenderer) renderCodeSpanCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

// 渲染代码块
func (r *JSONRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}
func (r *JSONRenderer) renderCodeBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}
func (r *JSONRenderer) renderCodeBlockInfoMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderCodeBlockCode(node *ast.Node, entering bool) ast.WalkStatus {
	var language string
	if 0 < len(node.Previous.CodeBlockInfo) {
		infoWords := lex.Split(node.Previous.CodeBlockInfo, lex.ItemSpace)
		language = util.BytesToStr(infoWords[0])
	}
	if entering {
		r.openObj()
		tokens := node.Tokens
		if 0 < len(node.Previous.CodeBlockInfo) {
			r.language(ast.NodeCodeBlock, util.BytesToStr(tokens), language)

			if "mindmap" == language {
				// 给echarts调用
				json := EChartsMindmap(tokens)
				r.mindMap(util.BytesToStr(json))
			}
		} else {
			tokens = html.EscapeHTML(tokens)
			r.language(ast.NodeCodeBlock, util.BytesToStr(tokens), language)
		}
	} else {
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderCodeBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderMathBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

// 解析数学块
func (r *JSONRenderer) renderMathBlockContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(ast.NodeMathBlock, util.BytesToStr(node.Tokens))
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderMathBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderInlineMath(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderInlineMathOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderInlineMathCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderEmAsteriskOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderEmUnderscoreOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

// 斜体标志存在
func (r *JSONRenderer) renderEmphasis(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderEmAsteriskCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderEmUnderscoreCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderStrongA6kOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderStrongU8eOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderStrongA6kCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderStrongU8eCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

// 引用标识存在 - flag
func (r *JSONRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderBlockquoteMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

// 渲染标题
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

func (r *JSONRenderer) renderHeadingC8hMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderHeadingID(node *ast.Node, entering bool) ast.WalkStatus {
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

// flag ListItem标识
func (r *JSONRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
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

// 分割线
func (r *JSONRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.leaf(node.Type, "hr", node)
	}
	return ast.WalkSkipChildren
}

// 硬换行
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

// 行内HTML
func (r *JSONRenderer) renderInlineHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		tokens := node.Tokens
		if r.Options.Sanitize {
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

// 处理链接
func (r *JSONRenderer) renderLink(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		dest := node.ChildByType(ast.NodeLinkDest)
		destTokens := dest.Tokens
		destTokens = r.LinkPath(destTokens)
		r.val(node.Type, util.BytesToStr(destTokens))
	} else {
		r.closeObj(node)
	}
	return ast.WalkContinue
}

// 处理图片
func (r *JSONRenderer) renderImage(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		destTokens := node.ChildByType(ast.NodeLinkDest).Tokens
		destTokens = r.LinkPath(destTokens)
		r.val(node.Type, util.BytesToStr(destTokens))
	} else {
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderBang(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderOpenBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderCloseBracket(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderOpenParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderCloseParen(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderLess(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderGreater(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderOpenBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderCloseBrace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderLinkTitle(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderLinkSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderLinkText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		var tokens []byte
		if r.Options.AutoSpace {
			tokens = r.Space(node.Tokens)
		} else {
			tokens = node.Tokens
		}
		r.WriteString(",\"title\":\"" + util.BytesToStr(tokens) + "\"")
	}
	return ast.WalkContinue
}

// 删除线标识
func (r *JSONRenderer) renderStrikethrough(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderStrikethrough1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderStrikethrough1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderStrikethrough2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderStrikethrough2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderTaskListItemMarker(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		check := "false"
		if node.TaskListItemChecked {
			check = "true"
		}
		r.val(node.Type, check)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj(node)
	}
	return ast.WalkContinue
}

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
		r.flag(node)
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
		r.flag(node)
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
		r.flag(node)
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

func (r *JSONRenderer) renderFootnotesDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

// 不支持脚注
func (r *JSONRenderer) renderFootnotesRef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderFootnotesDef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
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

// Yaml处理
func (r *JSONRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
	} else {
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r JSONRenderer) renderYamlFrontMatterOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}
func (r JSONRenderer) renderYamlFrontMatterContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.val(ast.NodeYamlFrontMatter, util.BytesToStr(node.Tokens))
	}
	return ast.WalkSkipChildren
}
func (r JSONRenderer) renderYamlFrontMatterCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderBackslashContent(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(ast.NodeBackslash, util.BytesToStr(node.Tokens))
	} else {
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderBackslash(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	// 忽略大纲渲染
	return ast.WalkContinue
}

func (r *JSONRenderer) renderMark(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderMark1OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderMark1CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderMark2OpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderMark2CloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderSup(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderSupOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}
func (r *JSONRenderer) renderSupCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderSub(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderSubOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderSubCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderKramdownSpanIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderBlockQueryEmbedScript(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.val(ast.NodeBlockQueryEmbed, util.BytesToStr(node.Tokens))
	} else {
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderTag(node *ast.Node, entering bool) ast.WalkStatus {
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

func (r *JSONRenderer) renderTagOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderTagCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderSuperBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderSuperBlockOpenMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderSuperBlockLayoutMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderSuperBlockCloseMarker(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

// 不支持链接引用
func (r *JSONRenderer) renderLinkRefDefBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderLinkRefDef(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) renderBlockRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.flag(node)
	} else {
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderBlockRefID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderBlockRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderBlockRefText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.val(ast.NodeBlockRef, util.BytesToStr(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderBlockRefDynamicText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.val(ast.NodeBlockRef, util.BytesToStr(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderFileAnnotationRef(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.flag(node)
	} else {
		r.closeObj(node)
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderFileAnnotationRefID(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderFileAnnotationRefSpace(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderFileAnnotationRefText(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.val(ast.NodeFileAnnotationRef, util.BytesToStr(node.Tokens))
	}
	return ast.WalkContinue
}

func (r *JSONRenderer) renderDefault(n *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *JSONRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		tokens := node.Tokens
		if r.Options.Sanitize {
			tokens = sanitize(tokens)
		}
		r.val(node.Type, util.BytesToStr(tokens))
	} else {
		r.closeObj(node)
	}
	return ast.WalkSkipChildren
}

func (r *JSONRenderer) leaf(nodeType ast.NodeType, val string, node *ast.Node) {
	r.openObj()
	r.val(nodeType, val)
	r.closeObj(node)
}

func (r *JSONRenderer) val(nodeType ast.NodeType, val string) {
	val = strings.ReplaceAll(val, "\\", "\\\\")
	val = strings.ReplaceAll(val, "\n", "\\n")
	val = strings.ReplaceAll(val, "\"", "\\\"")
	val = strings.ReplaceAll(val, "'", "\\'")
	// 要写入Type和Value
	r.WriteString("\"type\":\"" + nodeType.String()[4:] + "\"" + ",")
	r.WriteString("\"value\":\"" + val + "\"")
}

func (r *JSONRenderer) language(nodeType ast.NodeType, val string, language string) {
	val = strings.ReplaceAll(val, "\\", "\\\\")
	val = strings.ReplaceAll(val, "\n", "\\n")
	val = strings.ReplaceAll(val, "\"", "\\\"")
	val = strings.ReplaceAll(val, "'", "\\'")
	r.WriteString("\"type\":\"" + nodeType.String()[4:] + "\"" + ",")
	r.WriteString("\"value\":\"" + val + "\"")
	r.WriteString(",\"language\":\"" + language + "\"")
}

// 写入 mindmap
func (r *JSONRenderer) mindMap(val string) {
	r.WriteString(",\"mindmap\":\"" + val + "\"")
}

func (r *JSONRenderer) flag(node *ast.Node) {
	r.WriteString("\"flag\":\"" + node.Type.String()[4:] + "\"")
}

// 打开对象
func (r *JSONRenderer) openObj() {
	r.WriteByte('{')
}

func (r *JSONRenderer) closeObj(node *ast.Node) {
	r.WriteByte('}')
	if !r.ignore(node.Next) {
		r.comma()
	}
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

// 忽略函数
func (r *JSONRenderer) ignore(node *ast.Node) bool {
	return nil == node || isOpenMarker(node) || isCloseMarker(node)
}
