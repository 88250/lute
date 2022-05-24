// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
)

const (
	// Zwsp 零宽空格。
	Zwsp = "\u200b"

	// Zwj 零宽连字。
	Zwj = "\u200d"
)

// Parse 会将 markdown 原始文本字节数组解析为一棵语法树。
func Parse(name string, markdown []byte, options *Options) (tree *Tree) {
	tree = &Tree{Name: name, Context: &Context{ParseOption: options}}
	tree.Context.Tree = tree
	tree.lexer = lex.NewLexer(markdown)
	tree.Root = &ast.Node{Type: ast.NodeDocument}
	tree.parseBlocks()
	tree.parseInlines()
	tree.finalParseBlockIAL()
	tree.lexer = nil
	return
}

func (t *Tree) finalParseBlockIAL() {
	if !t.Context.ParseOption.KramdownBlockIAL {
		return
	}

	// 补全空段落
	var appends []*ast.Node

	ast.Walk(t.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering || !n.IsBlock() || ast.NodeKramdownBlockIAL == n.Type {
			return ast.WalkContinue
		}

		if "" == n.ID {
			n.ID = ast.NewNodeID()
		}

		if ast.NodeBlockquote == n.Type && nil != n.FirstChild && nil == n.FirstChild.Next {
			appends = append(appends, n)
		}

		ial := n.Next
		if nil == ial || ast.NodeKramdownBlockIAL != ial.Type {
			if t.Context.ParseOption.ProtyleWYSIWYG {
				n.SetIALAttr("updated", n.ID[:14])
			}
			return ast.WalkContinue
		}

		n.KramdownIAL = Tokens2IAL(ial.Tokens)
		if "" == n.IALAttr("updated") && t.Context.ParseOption.ProtyleWYSIWYG {
			n.SetIALAttr("updated", n.ID[:14])
			ial.Tokens = IAL2Tokens(n.KramdownIAL)
		}
		return ast.WalkContinue
	})

	for _, n := range appends {
		id := ast.NewNodeID()
		ialTokens := []byte("{: id=\"" + id + "\"}")
		p := &ast.Node{Type: ast.NodeParagraph, ID: id}
		p.KramdownIAL = [][]string{{"id", id}, {"updated", id[:14]}}
		p.ID = id
		p.InsertAfter(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: ialTokens})
		n.AppendChild(p)
	}

	var docIAL *ast.Node
	var id string
	if nil != t.Context.rootIAL {
		docIAL = t.Context.rootIAL
	} else {
		id = ast.NewNodeID()
		docIAL = &ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: []byte("{: id=\"" + id + "\" updated=\"" + id[:14] + "\" type=\"doc\"}")}
		t.Root.ID = id
		t.ID = id
	}
	t.Root.AppendChild(docIAL)
}

// Block 会将 markdown 原始文本字节数组解析为一棵语法树，该语法树的第一个块级子节点是段落节点。
func Block(name string, markdown []byte, options *Options) (tree *Tree) {
	tree = &Tree{Name: name, Context: &Context{ParseOption: options}}
	tree.Context.Tree = tree
	tree.lexer = lex.NewLexer(markdown)
	tree.Root = &ast.Node{Type: ast.NodeDocument}
	tree.parseBlocks()
	tree.finalParseBlockIAL()
	tree.lexer = nil
	return
}

// Inline 会将 markdown 原始文本字节数组解析为一棵语法树，该语法树的第一个块级子节点是段落节点。
func Inline(name string, markdown []byte, options *Options) (tree *Tree) {
	tree = &Tree{Name: name, Context: &Context{ParseOption: options}}
	tree.Context.Tree = tree
	tree.Root = &ast.Node{Type: ast.NodeDocument}
	tree.Root.AppendChild(&ast.Node{Type: ast.NodeParagraph, Tokens: markdown})
	tree.parseInlines()
	tree.lexer = nil
	return
}

// Context 用于维护块级元素解析过程中使用到的公共数据。
type Context struct {
	Tree        *Tree    // 关联的语法树
	ParseOption *Options // 解析选项

	Tip                                                      *ast.Node // 末梢节点
	oldtip                                                   *ast.Node // 老的末梢节点
	currentLine                                              []byte    // 当前行
	currentLineLen                                           int       // 当前行长
	offset, column, nextNonspace, nextNonspaceColumn, indent int       // 解析时用到的下标、缩进空格数等
	indented, blank, partiallyConsumedTab, allClosed         bool      // 是否是缩进行、空行等标识
	lastMatchedContainer                                     *ast.Node // 最后一个匹配的块节点

	rootIAL *ast.Node // 根节点 kramdown IAL
}

// InlineContext 描述了行级元素解析上下文。
type InlineContext struct {
	tokens     []byte     // 当前解析的 Tokens
	tokensLen  int        // 当前解析的 Tokens 长度
	pos        int        // 当前解析到的 token 位置
	delimiters *delimiter // 分隔符栈，用于强调解析
	brackets   *delimiter // 括号栈，用于图片和链接解析
}

// advanceOffset 用于移动 count 个字符位置，columns 指定了遇到 tab 时是否需要空格进行补偿偏移。
func (context *Context) advanceOffset(count int, columns bool) {
	currentLine := context.currentLine
	var charsToTab, charsToAdvance int
	var c byte
	for 0 < count {
		c = currentLine[context.offset]
		if lex.ItemTab == c {
			charsToTab = 4 - (context.column % 4)
			if columns {
				context.partiallyConsumedTab = charsToTab > count
				if context.partiallyConsumedTab {
					charsToAdvance = count
				} else {
					charsToAdvance = charsToTab
					context.offset++
				}
				context.column += charsToAdvance
				count -= charsToAdvance
			} else {
				context.partiallyConsumedTab = false
				context.column += charsToTab
				context.offset++
				count--
			}
		} else {
			context.partiallyConsumedTab = false
			context.offset++
			context.column++ // 假定是 ASCII，因为块开始标记符都是 ASCII
			count--
		}
	}
}

// advanceNextNonspace 用于预移动到下一个非空字符位置。
func (context *Context) advanceNextNonspace() {
	context.offset = context.nextNonspace
	context.column = context.nextNonspaceColumn
	context.partiallyConsumedTab = false
}

// findNextNonspace 用于查找下一个非空字符。
func (context *Context) findNextNonspace() {
	i := context.offset
	cols := context.column

	var token byte
	for {
		token = context.currentLine[i]
		if lex.ItemSpace == token {
			i++
			cols++
		} else if lex.ItemTab == token {
			i++
			cols += 4 - (cols % 4)
		} else {
			break
		}
	}

	context.blank = lex.ItemNewline == token
	context.nextNonspace = i
	context.nextNonspaceColumn = cols
	context.indent = context.nextNonspaceColumn - context.column
	context.indented = 4 <= context.indent
}

// closeUnmatchedBlocks 最终化所有未匹配的块节点。
func (context *Context) closeUnmatchedBlocks() {
	if !context.allClosed {
		for context.oldtip != context.lastMatchedContainer {
			parent := context.oldtip.Parent
			context.finalize(context.oldtip)
			context.oldtip = parent
		}
		context.allClosed = true
	}
}

// closeSuperBlockChildren 最终化超级块下的子节点。
func (context *Context) closeSuperBlockChildren() {
	for n := context.Tip; nil != n && ast.NodeSuperBlock != n.Type; n = n.Parent {
		context.finalize(n)
	}
}

// finalize 执行 block 的最终化处理。调用该方法会将 context.Tip 置为 block 的父节点。
func (context *Context) finalize(block *ast.Node) {
	parent := block.Parent
	block.Close = true

	// 节点最终化处理。比如围栏代码块提取 info 部分；HTML 代码块剔除结尾空格；段落需要解析链接引用定义等。
	switch block.Type {
	case ast.NodeCodeBlock:
		context.codeBlockFinalize(block)
	case ast.NodeHTMLBlock, ast.NodeIFrame, ast.NodeVideo, ast.NodeAudio, ast.NodeWidget:
		context.htmlBlockFinalize(block)
	case ast.NodeParagraph:
		insertTable := paragraphFinalize(block, context)
		if insertTable {
			return
		}
	case ast.NodeMathBlock:
		context.mathBlockFinalize(block)
	case ast.NodeYamlFrontMatter:
		context.yamlFrontMatterFinalize(block)
	case ast.NodeList:
		context.listFinalize(block)
	case ast.NodeSuperBlock:
		context.superBlockFinalize(block)
	case ast.NodeGitConflict:
		context.gitConflictFinalize(block)
	}

	context.Tip = parent
}

// addChildMarker 将构造一个 NodeType 节点并作为子节点添加到末梢节点 context.Tip 上。
func (context *Context) addChildMarker(nodeType ast.NodeType, tokens []byte) (ret *ast.Node) {
	ret = &ast.Node{Type: nodeType, Tokens: tokens, Close: true}
	context.Tip.AppendChild(ret)
	return
}

// addChild 将构造一个 NodeType 节点并作为子节点添加到末梢节点 context.Tip 上。如果末梢不能接受子节点（非块级容器不能添加子节点），则最终化该末梢
// 节点并向父节点方向尝试，直到找到一个能接受该子节点的节点为止。添加完成后该子节点会被设置为新的末梢节点。
func (context *Context) addChild(nodeType ast.NodeType) (ret *ast.Node) {
	for !context.Tip.CanContain(nodeType) {
		context.finalize(context.Tip) // 注意调用 finalize 会向父节点方向进行迭代
	}

	ret = &ast.Node{Type: nodeType}
	context.Tip.AppendChild(ret)
	context.Tip = ret
	return
}

// listsMatch 用户判断指定的 listData 和 itemData 是否可归属于同一个列表。
func (context *Context) listsMatch(listData, itemData *ast.ListData) bool {
	return listData.Typ == itemData.Typ &&
		((0 == listData.Delimiter && 0 == itemData.Delimiter) || listData.Delimiter == itemData.Delimiter) &&
		listData.BulletChar == itemData.BulletChar
}

// Tree 描述了 Markdown 抽象语法树结构。
type Tree struct {
	Root          *ast.Node      // 根节点
	Context       *Context       // 块级解析上下文
	lexer         *lex.Lexer     // 词法分析器
	inlineContext *InlineContext // 行级解析上下文

	Name    string   // 名称
	ID      string   // ID
	Box     string   // 容器
	Path    string   // 路径
	HPath   string   // 人类可读的路径
	Marks   []string // 文本标记
	Created int64    // 创建时间
	Updated int64    // 更新时间
	Hash    string   // 内容哈希
}

// Options 描述了解析选项。
type Options struct {
	// GFMTable 设置是否打开“GFM 表”支持。
	GFMTable bool
	// GFMTaskListItem 设置是否打开“GFM 任务列表项”支持。
	GFMTaskListItem bool
	// GFMStrikethrough 设置是否打开“GFM 删除线”支持。
	GFMStrikethrough bool
	// GFMAutoLink 设置是否打开“GFM 自动链接”支持。
	GFMAutoLink bool
	// Footnotes 设置是否打开“脚注”支持。
	Footnotes bool
	// HeadingID 设置是否打开“自定义标题 ID”支持。
	HeadingID bool
	// ToC 设置是否打开“目录”支持。
	ToC bool
	// Emoji 设置是否对 Emoji 别名替换为原生 Unicode 字符。
	Emoji bool
	// AliasEmoji 存储 ASCII 别名到表情 Unicode 映射。
	AliasEmoji map[string]string
	// EmojiAlias 存储表情 Unicode 到 ASCII 别名映射。
	EmojiAlias map[string]string
	// EmojiSite 设置图片 Emoji URL 的路径前缀。
	EmojiSite string
	// Vditor 所见即所得支持。
	VditorWYSIWYG bool
	// Vditor 即时渲染支持。
	VditorIR bool
	// Vditor 分屏预览支持。
	VditorSV bool
	// Protyle 所见即所得支持。
	ProtyleWYSIWYG bool
	// InlineMathAllowDigitAfterOpenMarker 设置内联数学公式是否允许起始 $ 后紧跟数字 https://github.com/b3log/lute/issues/38
	InlineMathAllowDigitAfterOpenMarker bool
	// Setext 设置是否解析 Setext 标题 https://github.com/88250/lute/issues/50
	Setext bool
	// YamlFrontMatter 设置是否开启 YAML Front Matter 支持。
	YamlFrontMatter bool
	// BlockRef 设置是否开启内容块引用支持。
	BlockRef bool
	// FileAnnotationRef 设置是否开启文件注解引用支持。
	FileAnnotationRef bool
	// Mark 设置是否打开 ==标记== 支持。
	Mark bool
	// KramdownBlockIAL 设置是否打开 kramdown 块级内联属性列表支持。 https://kramdown.gettalong.org/syntax.html#inline-attribute-lists
	KramdownBlockIAL bool
	// KramdownSpanIAL 设置是否打开 kramdown 行级内联属性列表支持。
	KramdownSpanIAL bool
	// Tag 设置是否开启 #标签# 支持。
	Tag bool
	// ImgPathAllowSpace 设置是否支持图片路径带空格。
	ImgPathAllowSpace bool
	// SuperBlock 设置是否支持超级块。 https://github.com/88250/lute/issues/111
	SuperBlock bool
	// Sup 设置是否打开 ^上标^ 支持。
	Sup bool
	// Sub 设置是否打开 ~下标~ 支持。
	Sub bool
	// GitConflict 设置是否打开 Git 冲突标记支持。
	GitConflict bool
	// LinkRef 设置是否打开“链接引用”支持。
	LinkRef bool
	// IndentCodeBlock 设置是否打开“缩进代码块”支持。
	IndentCodeBlock bool
	// ParagraphBeginningSpace 设置是否打开“段首空格”支持。
	ParagraphBeginningSpace bool
	// DataImage 设置是否打开 ![foo](data:image...) 形式的图片支持。
	DataImage bool
}

func NewOptions() *Options {
	return &Options{
		GFMTable:         true,
		GFMTaskListItem:  true,
		GFMStrikethrough: true,
		GFMAutoLink:      true,
		Footnotes:        true,
		Emoji:            true,
		AliasEmoji:       EmojiAliasUnicode,
		EmojiAlias:       EmojiUnicodeAlias,
		//EmojiSite:         "https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji",
		EmojiSite:         "https://unpkg.com/vditor/dist/images/emoji", // 表情图片从 `cdn.jsdelivr.net` 切换到 `unpkg.com` https://github.com/88250/lute/issues/171
		Setext:            true,
		YamlFrontMatter:   true,
		BlockRef:          false,
		FileAnnotationRef: false,
		Mark:              false,
		KramdownBlockIAL:  false,
		HeadingID:         true,
		LinkRef:           true,
		IndentCodeBlock:   true,
		DataImage:         true,
	}
}

func (context *Context) ParentTip() {
	if tip := context.Tip.Parent; nil != tip {
		context.Tip = context.Tip.Parent
	}
}

func (context *Context) TipAppendChild(child *ast.Node) {
	context.Tip.AppendChild(child)
}
