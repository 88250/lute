// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package ast

import (
	"bytes"
	"math/rand"
	"time"
	"unicode/utf8"

	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

// Node 描述了节点结构。
type Node struct {
	// 不用接口实现的原因：
	//   1. 转换节点类型非常方便，只需修改 Type 属性
	//   2. 为了极致的性能而牺牲扩展性

	// 节点基础结构

	ID   string `json:",omitempty"` // 节点的唯一标识
	URL  string `json:"-"`          // 地址部分
	Path string `json:"-"`          // 地址路径部分

	Type       NodeType // 节点类型
	Parent     *Node    `json:"-"`          // 父节点
	Previous   *Node    `json:"-"`          // 前一个兄弟节点
	Next       *Node    `json:"-"`          // 后一个兄弟节点
	FirstChild *Node    `json:"-"`          // 第一个子节点
	LastChild  *Node    `json:"-"`          // 最后一个子节点
	Children   []*Node  `json:",omitempty"` // 所有子节点
	Tokens     []byte   `json:",omitempty"` // 词法分析结果 Tokens，语法分析阶段会继续操作这些 Tokens

	// 解析过程标识

	Close           bool `json:",omitempty"` // 标识是否关闭
	LastLineBlank   bool `json:",omitempty"` // 标识最后一行是否是空行
	LastLineChecked bool `json:",omitempty"` // 标识最后一行是否检查过

	// 代码

	CodeMarkerLen int `json:",omitempty"` // ` 个数，1 或 2

	// 代码块

	IsFencedCodeBlock  bool `json:",omitempty"`
	CodeBlockFenceChar byte `json:",omitempty"`

	CodeBlockFenceLen    int    `json:",omitempty"`
	CodeBlockFenceOffset int    `json:",omitempty"`
	CodeBlockOpenFence   []byte `json:",omitempty"`
	CodeBlockInfo        []byte `json:",omitempty"`
	CodeBlockCloseFence  []byte `json:",omitempty"`

	// HTML 块

	HtmlBlockType int `json:",omitempty"` // 规范中定义的 HTML 块类型（1-7）

	// 列表、列表项

	*ListData `json:",omitempty"`

	// 任务列表项 [ ]、[x] 或者 [X]

	TaskListItemChecked bool `json:",omitempty"` // 是否勾选

	// 表

	TableAligns              []int `json:",omitempty"` // 从左到右每个表格节点的对齐方式，0：默认对齐，1：左对齐，2：居中对齐，3：右对齐
	TableCellAlign           int   `json:",omitempty"` // 表的单元格对齐方式
	TableCellContentWidth    int   `json:",omitempty"` // 表的单元格内容宽度（字节数）
	TableCellContentMaxWidth int   `json:",omitempty"` // 表的单元格内容最大宽度

	// 链接

	LinkType     int    `json:",omitempty"` // 链接类型，0：内联链接 [foo](/bar)，1：链接引用定义 [foo]: /bar，2：自动链接，3：链接引用 [foo]
	LinkRefLabel []byte `json:",omitempty"` // 链接引用 label，[label] 或者 [text][label] 形式，[label] 情况下 text 和 label 相同

	// 标题

	HeadingLevel        int    `json:",omitempty"` // 1~6
	HeadingSetext       bool   `json:",omitempty"` // 是否为 Setext
	HeadingNormalizedID string `json:",omitempty"` // 规范化后的 ID

	// 数学公式块

	MathBlockDollarOffset int `json:",omitempty"`

	// 脚注

	FootnotesRefLabel []byte  `json:",omitempty"` // 脚注引用 label，[^label]
	FootnotesRefId    string  `json:",omitempty"` // 脚注 id
	FootnotesRefs     []*Node `json:",omitempty"` // 脚注引用

	// HTML 实体

	HtmlEntityTokens []byte `json:",omitempty"` // 原始输入的实体 tokens，&amp;

	// Kramdown 内联属性列表
	KramdownIAL [][]string
}

// ListData 用于记录列表或列表项节点的附加信息。
type ListData struct {
	Typ          int    `json:",omitempty"` // 0：无序列表，1：有序列表，3：任务列表
	Tight        bool   `json:",omitempty"` // 是否是紧凑模式
	BulletChar   byte   `json:",omitempty"` // 无序列表标识，* - 或者 +
	Start        int    `json:",omitempty"` // 有序列表起始序号
	Delimiter    byte   `json:",omitempty"` // 有序列表分隔符，. 或者 )
	Padding      int    `json:",omitempty"` // 列表内部缩进空格数（包含标识符长度，即规范中的 W+N）
	MarkerOffset int    `json:",omitempty"` // 标识符（* - + 或者 1 2 3）相对缩进空格数
	Checked      bool   `json:",omitempty"` // 任务列表项是否勾选
	Marker       []byte `json:",omitempty"` // 列表标识符
	Num          int    `json:",omitempty"` // 有序列表项修正过的序号
}

// Testing 标识是否为测试环境。
var Testing bool

func NewNodeID() string {
	if Testing {
		return "20060102150405-1a2b3c4" // 测试环境 ID
	}
	now := time.Now()
	return now.Format("20060102150405") + "-" + randStr(7)
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func randStr(length int) string {
	letter := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func (n *Node) RemoveIALAttr(name string) {
	tmp := n.KramdownIAL[:0]
	for _, kv := range n.KramdownIAL {
		if name != kv[0] {
			tmp = append(tmp, kv)
		}
	}
	n.KramdownIAL = tmp
}

func (n *Node) SetIALAttr(name, value string) {
	for _, kv := range n.KramdownIAL {
		if name == kv[0] {
			kv[1] = value
			return
		}
	}

	if typ := n.IALAttr("type"); "doc" == typ {
		// 让 type="doc" 保持在最后一个
		n.RemoveIALAttr("type")
		n.KramdownIAL = append(n.KramdownIAL, []string{name, value})
		n.KramdownIAL = append(n.KramdownIAL, []string{"type", "doc"})
		return
	}
	n.KramdownIAL = append(n.KramdownIAL, []string{name, value})
}

func (n *Node) IALAttr(name string) string {
	for _, kv := range n.KramdownIAL {
		if name == kv[0] {
			return kv[1]
		}
	}
	return ""
}

// TokensStr 返回 n 的 Tokens 字符串。
func (n *Node) TokensStr() string {
	return util.BytesToStr(n.Tokens)
}

// LastDeepestChild 返回 n 的最后一个最深子节点。
func (n *Node) LastDeepestChild() (ret *Node) {
	if nil == n.LastChild {
		return n
	}
	return n.LastChild.LastDeepestChild()
}

// FirstDeepestChild 返回 n 的第一个最深的子节点。
func (n *Node) FirstDeepestChild() (ret *Node) {
	if nil == n.FirstChild {
		return n
	}
	return n.FirstChild.FirstDeepestChild()
}

// LinkDest 在 n 的子节点中查找 childType 指定类型的第一个子节点。
func (n *Node) ChildByType(childType NodeType) *Node {
	for c := n.FirstChild; nil != c; c = c.Next {
		if c.Type == childType {
			return c
		}
	}
	return nil
}

// ChildrenByType 返回 n 下所有类型为 childType 的子节点。
func (n *Node) ChildrenByType(childType NodeType) (ret []*Node) {
	ret = []*Node{}
	Walk(n, func(n *Node, entering bool) WalkStatus {
		if (childType == n.Type) && entering {
			ret = append(ret, n)
		}
		return WalkContinue
	})
	return
}

// Text 返回 n 及其文本子节点的文本值。
func (n *Node) Text() (ret string) {
	buf := &bytes.Buffer{}
	Walk(n, func(n *Node, entering bool) WalkStatus {
		if !entering {
			return WalkContinue
		}
		switch n.Type {
		case NodeText, NodeLinkText, NodeBlockRefText, NodeBlockEmbedText, NodeFootnotesRef:
			buf.Write(n.Tokens)
		}
		return WalkContinue
	})
	return buf.String()
}

// TextLen 返回 n 及其文本子节点的累计长度。
func (n *Node) TextLen() (ret int) {
	buf := make([]byte, 0, 4096)
	Walk(n, func(n *Node, entering bool) WalkStatus {
		if !entering {
			return WalkContinue
		}
		switch n.Type {
		case NodeText, NodeLinkText, NodeBlockRefText, NodeBlockEmbedText, NodeFootnotesRef:
			buf = append(buf, n.Tokens...)
		}
		return WalkContinue
	})
	return utf8.RuneCount(buf)
}

// Text 返回 n 及其所有内容子节点的文本值。
func (n *Node) Content() (ret string) {
	buf := &bytes.Buffer{}
	Walk(n, func(n *Node, entering bool) WalkStatus {
		if !entering {
			return WalkContinue
		}
		switch n.Type {
		case NodeText, NodeLinkText, NodeBlockRefText, NodeBlockEmbedText, NodeFootnotesRef,
			NodeCodeSpanContent, NodeCodeBlockCode, NodeInlineMathContent, NodeMathBlockContent,
			NodeHTMLEntity, NodeEmojiAlias, NodeEmojiUnicode, NodeBackslashContent, NodeYamlFrontMatterContent,
			NodeGitConflictContent:
			buf.Write(n.Tokens)
		}
		return WalkContinue
	})
	return buf.String()
}

// ContentLen 返回 n 及其所有内容子节点的累计长度。
func (n *Node) ContentLen() (ret int) {
	buf := make([]byte, 0, 8192)
	Walk(n, func(n *Node, entering bool) WalkStatus {
		if !entering {
			return WalkContinue
		}
		switch n.Type {
		case NodeText, NodeLinkText, NodeBlockRefText, NodeBlockEmbedText, NodeFootnotesRef,
			NodeCodeSpanContent, NodeCodeBlockCode, NodeInlineMathContent, NodeMathBlockContent,
			NodeHTMLEntity, NodeEmojiAlias, NodeEmojiUnicode, NodeBackslashContent, NodeYamlFrontMatterContent,
			NodeGitConflictContent:
			buf = append(buf, n.Tokens...)
		}
		return WalkContinue
	})
	return utf8.RuneCount(buf)
}

// TokenLen 返回 n 及其子节点 tokens 累计长度。
func (n *Node) TokenLen() (ret int) {
	Walk(n, func(n *Node, entering bool) WalkStatus {
		if !entering {
			return WalkContinue
		}
		ret += lex.BytesShowLength(n.Tokens)
		return WalkContinue
	})
	return
}

// IsChildBlockOf 用于检查块级节点 n 的父节点是否是 parent 节点，depth 指定层级，0 为任意层级。
// n 如果不是块级节点，则直接返回 false。
func (n *Node) IsChildBlockOf(parent *Node, depth int) bool {
	if "" == n.ID || !n.IsBlock() {
		return false
	}

	if depth == 0 {
		// 任何层级上只要 n 的父节点和 parent 一样就认为是子节点
		for p := n.Parent; nil != p; p = p.Parent {
			if p == parent {
				return true
			}
		}
		return false
	}

	// 只在指定层级上匹配父节点
	nodeParent := n.Parent
	for i := 1; i < depth; i++ {
		if nil == nodeParent {
			break
		}
		nodeParent = nodeParent.Parent
	}
	if parent != nodeParent {
		return false
	}
	return true
}

func (n *Node) NextNodeText() string {
	if nil == n.Next {
		return ""
	}
	return n.Next.Text()
}

func (n *Node) PreviousNodeText() string {
	if nil == n.Previous {
		return ""
	}
	return n.Previous.Text()
}

// Unlink 用于将节点从树上移除，后一个兄弟节点会接替该节点。
func (n *Node) Unlink() {
	if nil != n.Previous {
		n.Previous.Next = n.Next
	} else if nil != n.Parent {
		n.Parent.FirstChild = n.Next
	}
	if nil != n.Next {
		n.Next.Previous = n.Previous
	} else if nil != n.Parent {
		n.Parent.LastChild = n.Previous
	}
	n.Parent = nil
	n.Next = nil
	n.Previous = nil
}

// AppendTokens 添加 Tokens。
func (n *Node) AppendTokens(tokens []byte) {
	n.Tokens = append(n.Tokens, tokens...)
}

// InsertAfter 在当前节点后插入一个兄弟节点。
func (n *Node) InsertAfter(sibling *Node) {
	sibling.Unlink()
	sibling.Next = n.Next
	if nil != sibling.Next {
		sibling.Next.Previous = sibling
	}
	sibling.Previous = n
	n.Next = sibling
	sibling.Parent = n.Parent
	if nil != sibling.Parent && nil == sibling.Next && nil != sibling.Parent.LastChild {
		sibling.Parent.LastChild = sibling
	}
}

// InsertBefore 在当前节点前插入一个兄弟节点。
func (n *Node) InsertBefore(sibling *Node) {
	sibling.Unlink()
	sibling.Previous = n.Previous
	if nil != sibling.Previous {
		sibling.Previous.Next = sibling
	}
	sibling.Next = n
	n.Previous = sibling
	sibling.Parent = n.Parent
	if nil != sibling.Parent && nil == sibling.Previous {
		sibling.Parent.FirstChild = sibling
	}
}

// AppendChild 在 n 的子节点最后再添加一个子节点。
func (n *Node) AppendChild(child *Node) {
	child.Unlink()
	child.Parent = n
	if nil != n.LastChild {
		n.LastChild.Next = child
		child.Previous = n.LastChild
		n.LastChild = child
	} else {
		n.FirstChild = child
		n.LastChild = child
	}
}

// PrependChild 在 n 的子节点最前添加一个子节点。
func (n *Node) PrependChild(child *Node) {
	child.Unlink()
	child.Parent = n
	if nil != n.FirstChild {
		n.FirstChild.Previous = child
		child.Next = n.FirstChild
		n.FirstChild = child
	} else {
		n.FirstChild = child
		n.LastChild = child
	}
}

// List 将 n 及其所有子节点按深度优先遍历添加到结果列表 ret 中。
func (n *Node) List() (ret []*Node) {
	ret = make([]*Node, 0, 512)
	Walk(n, func(n *Node, entering bool) WalkStatus {
		if entering {
			ret = append(ret, n)
		}
		return WalkContinue
	})
	return
}

// ParentIs 判断 n 的类型是否在指定的 nodeTypes 类型列表内。
func (n *Node) ParentIs(nodeType NodeType, nodeTypes ...NodeType) bool {
	types := append(nodeTypes, nodeType)
	for p := n.Parent; nil != p; p = p.Parent {
		for _, pt := range types {
			if pt == p.Type {
				return true
			}
		}
	}
	return false
}

// IsBlock 判断 n 是否为块级节点。
func (n *Node) IsBlock() bool {
	switch n.Type {
	case NodeDocument, NodeParagraph, NodeHeading, NodeThematicBreak, NodeBlockquote, NodeList, NodeListItem, NodeHTMLBlock,
		NodeCodeBlock, NodeTable, NodeMathBlock, NodeFootnotesDefBlock, NodeFootnotesDef, NodeToC, NodeYamlFrontMatter, NodeBlockEmbed, NodeBlockQueryEmbed,
		NodeKramdownBlockIAL, NodeSuperBlock, NodeGitConflict:
		return true
	}
	return false
}

// IsContainerBlock 判断 n 是否为容器块。
func (n *Node) IsContainerBlock() bool {
	switch n.Type {
	case NodeDocument, NodeBlockquote, NodeList, NodeListItem, NodeFootnotesDefBlock, NodeFootnotesDef, NodeSuperBlock:
		return true
	}
	return false
}

// IsMarker 判断 n 是否为标记符节点。
func (n *Node) IsMarker() bool {
	switch n.Type {
	case NodeHeadingC8hMarker, NodeBlockquoteMarker, NodeCodeBlockFenceOpenMarker, NodeCodeBlockFenceCloseMarker, NodeCodeBlockFenceInfoMarker,
		NodeEmA6kOpenMarker, NodeEmA6kCloseMarker, NodeEmU8eOpenMarker, NodeEmU8eCloseMarker, NodeStrongA6kOpenMarker, NodeStrongA6kCloseMarker,
		NodeStrongU8eOpenMarker, NodeStrongU8eCloseMarker, NodeCodeSpanOpenMarker, NodeCodeSpanCloseMarker, NodeTaskListItemMarker,
		NodeStrikethrough1OpenMarker, NodeStrikethrough1CloseMarker, NodeStrikethrough2OpenMarker, NodeStrikethrough2CloseMarker,
		NodeMathBlockOpenMarker, NodeMathBlockCloseMarker, NodeInlineMathOpenMarker, NodeInlineMathCloseMarker, NodeYamlFrontMatterOpenMarker, NodeYamlFrontMatterCloseMarker,
		NodeMark1OpenMarker, NodeMark1CloseMarker, NodeMark2OpenMarker, NodeMark2CloseMarker, NodeTagOpenMarker, NodeTagCloseMarker,
		NodeSuperBlockOpenMarker, NodeSuperBlockLayoutMarker, NodeSuperBlockCloseMarker, NodeSupOpenMarker, NodeSupCloseMarker, NodeSubOpenMarker, NodeSubCloseMarker:
		return true
	}
	return false
}

// AcceptLines 判断是否节点是否可以接受更多的文本行。比如 HTML 块、代码块和段落是可以接受更多的文本行的。
func (n *Node) AcceptLines() bool {
	switch n.Type {
	case NodeParagraph, NodeCodeBlock, NodeHTMLBlock, NodeMathBlock, NodeYamlFrontMatter, NodeBlockEmbed, NodeBlockQueryEmbed, NodeGitConflict, NodeIFrame:
		return true
	}
	return false
}

// CanContain 判断是否能够包含 NodeType 指定类型的节点。 比如列表节点（块级容器）只能包含列表项节点，
// 块引用节点（块级容器）可以包含任意节点；段落节点（叶子块节点）不能包含任何其他块级节点。
func (n *Node) CanContain(nodeType NodeType) bool {
	switch n.Type {
	case NodeCodeBlock, NodeHTMLBlock, NodeParagraph, NodeThematicBreak, NodeTable, NodeMathBlock, NodeYamlFrontMatter, NodeGitConflict, NodeIFrame:
		return false
	case NodeList:
		return NodeListItem == nodeType
	case NodeFootnotesDefBlock:
		return NodeFootnotesDef == nodeType
	case NodeFootnotesDef:
		return NodeFootnotesDef != nodeType
	case NodeSuperBlock:
		if nil != n.LastChild && NodeSuperBlockCloseMarker == n.LastChild.Type {
			// 超级块已经闭合
			return false
		}
		return true
	}
	return NodeListItem != nodeType
}

//go:generate stringer -type=NodeType
type NodeType int

func Str2NodeType(nodeTypeStr string) NodeType {
	for t := NodeDocument; t < NodeTypeMaxVal; t++ {
		if nodeTypeStr == t.String() {
			return t
		}
	}
	return -1
}

const (
	// CommonMark

	NodeDocument                  NodeType = 0  // 根
	NodeParagraph                 NodeType = 1  // 段落
	NodeHeading                   NodeType = 2  // 标题
	NodeHeadingC8hMarker          NodeType = 3  // ATX 标题标记符 #
	NodeThematicBreak             NodeType = 4  // 分隔线
	NodeBlockquote                NodeType = 5  // 块引用
	NodeBlockquoteMarker          NodeType = 6  // 块引用标记符 >
	NodeList                      NodeType = 7  // 列表
	NodeListItem                  NodeType = 8  // 列表项
	NodeHTMLBlock                 NodeType = 9  // HTML 块
	NodeInlineHTML                NodeType = 10 // 内联 HTML
	NodeCodeBlock                 NodeType = 11 // 代码块
	NodeCodeBlockFenceOpenMarker  NodeType = 12 // 开始围栏代码块标记符 ```
	NodeCodeBlockFenceCloseMarker NodeType = 13 // 结束围栏代码块标记符 ```
	NodeCodeBlockFenceInfoMarker  NodeType = 14 // 围栏代码块信息标记符 info string
	NodeCodeBlockCode             NodeType = 15 // 围栏代码块代码
	NodeText                      NodeType = 16 // 文本
	NodeEmphasis                  NodeType = 17 // 强调
	NodeEmA6kOpenMarker           NodeType = 18 // 开始强调标记符 *
	NodeEmA6kCloseMarker          NodeType = 19 // 结束强调标记符 *
	NodeEmU8eOpenMarker           NodeType = 20 // 开始强调标记符 _
	NodeEmU8eCloseMarker          NodeType = 21 // 结束强调标记符 _
	NodeStrong                    NodeType = 22 // 加粗
	NodeStrongA6kOpenMarker       NodeType = 23 // 开始加粗标记符 **
	NodeStrongA6kCloseMarker      NodeType = 24 // 结束加粗标记符 **
	NodeStrongU8eOpenMarker       NodeType = 25 // 开始加粗标记符 __
	NodeStrongU8eCloseMarker      NodeType = 26 // 结束加粗标记符 __
	NodeCodeSpan                  NodeType = 27 // 代码
	NodeCodeSpanOpenMarker        NodeType = 28 // 开始代码标记符 `
	NodeCodeSpanContent           NodeType = 29 // 代码内容
	NodeCodeSpanCloseMarker       NodeType = 30 // 结束代码标记符 `
	NodeHardBreak                 NodeType = 31 // 硬换行
	NodeSoftBreak                 NodeType = 32 // 软换行
	NodeLink                      NodeType = 33 // 链接
	NodeImage                     NodeType = 34 // 图片
	NodeBang                      NodeType = 35 // !
	NodeOpenBracket               NodeType = 36 // [
	NodeCloseBracket              NodeType = 37 // ]
	NodeOpenParen                 NodeType = 38 // (
	NodeCloseParen                NodeType = 39 // )
	NodeLinkText                  NodeType = 40 // 链接文本
	NodeLinkDest                  NodeType = 41 // 链接地址
	NodeLinkTitle                 NodeType = 42 // 链接标题
	NodeLinkSpace                 NodeType = 43 // 链接地址和链接标题之间的空格
	NodeHTMLEntity                NodeType = 44 // HTML 实体
	NodeLinkRefDefBlock           NodeType = 45 // 链接引用定义块
	NodeLinkRefDef                NodeType = 46 // 链接引用定义 [label]:

	// GFM

	NodeTaskListItemMarker        NodeType = 100 // 任务列表项标记符
	NodeStrikethrough             NodeType = 101 // 删除线
	NodeStrikethrough1OpenMarker  NodeType = 102 // 开始删除线标记符 ~
	NodeStrikethrough1CloseMarker NodeType = 103 // 结束删除线标记符 ~
	NodeStrikethrough2OpenMarker  NodeType = 104 // 开始删除线标记符 ~~
	NodeStrikethrough2CloseMarker NodeType = 105 // 结束删除线标记符 ~~
	NodeTable                     NodeType = 106 // 表
	NodeTableHead                 NodeType = 107 // 表头
	NodeTableRow                  NodeType = 108 // 表行
	NodeTableCell                 NodeType = 109 // 表格

	// Emoji

	NodeEmoji        NodeType = 200 // Emoji
	NodeEmojiUnicode NodeType = 201 // Emoji Unicode
	NodeEmojiImg     NodeType = 202 // Emoji 图片
	NodeEmojiAlias   NodeType = 203 // Emoji ASCII

	// 数学公式

	NodeMathBlock             NodeType = 300 // 数学公式块
	NodeMathBlockOpenMarker   NodeType = 301 // 开始数学公式块标记符 $$
	NodeMathBlockContent      NodeType = 302 // 数学公式块内容
	NodeMathBlockCloseMarker  NodeType = 303 // 结束数学公式块标记符 $$
	NodeInlineMath            NodeType = 304 // 内联数学公式
	NodeInlineMathOpenMarker  NodeType = 305 // 开始内联数学公式标记符 $
	NodeInlineMathContent     NodeType = 306 // 内联数学公式内容
	NodeInlineMathCloseMarker NodeType = 307 // 结束内联数学公式标记符 $

	// 转义

	NodeBackslash        NodeType = 400 // 转义反斜杠标记符 \
	NodeBackslashContent NodeType = 401 // 转义反斜杠后的内容

	// Vditor 支持

	NodeVditorCaret NodeType = 405 // 插入符，某些情况下需要使用该节点进行插入符位置调整

	// 脚注

	NodeFootnotesDefBlock NodeType = 410 // 脚注定义块
	NodeFootnotesDef      NodeType = 411 // 脚注定义 [^label]:
	NodeFootnotesRef      NodeType = 412 // 脚注引用 [^label]

	// 目录

	NodeToC NodeType = 415 // 目录 [toc]

	// 标题

	NodeHeadingID NodeType = 420 // 标题 ID # foo {id}

	// YAML Front Matter

	NodeYamlFrontMatter            NodeType = 425 // https://jekyllrb.com/docs/front-matter/
	NodeYamlFrontMatterOpenMarker  NodeType = 426 // 开始 YAML Front Matter 标记符 ---
	NodeYamlFrontMatterContent     NodeType = 427 // YAML Front Matter 内容
	NodeYamlFrontMatterCloseMarker NodeType = 428 // 结束 YAML Front Matter 标记符 ---

	// 内容块引用（Block Reference） https://github.com/88250/lute/issues/82

	NodeBlockRef                    NodeType = 430 // 内容块引用节点
	NodeBlockRefID                  NodeType = 431 // 被引用的内容块（定义块）ID
	NodeBlockRefSpace               NodeType = 432 // 被引用的内容块 ID 和内容块引用锚文本之间的空格
	NodeBlockRefText                NodeType = 433 // 内容块引用锚点本
	NodeBlockRefTextTplRenderResult NodeType = 434 // 内容块引用锚文本模板渲染结果

	// 内容块嵌入（Block Embed） https://github.com/88250/lute/issues/91

	NodeBlockEmbed                    NodeType = 440 // 内容块嵌入节点
	NodeBlockEmbedID                  NodeType = 441 // 被引用的内容块（定义块） ID
	NodeBlockEmbedSpace               NodeType = 442 // 被引用的内容块 ID 和内容块引用锚文本之间的空格
	NodeBlockEmbedText                NodeType = 443 // 内容块嵌入锚点本
	NodeBlockEmbedTextTplRenderResult NodeType = 444 // 内容块嵌入锚点本模板渲染结果

	// ==Mark== 标记语法 https://github.com/88250/lute/issues/84

	NodeMark             NodeType = 450 // 标记
	NodeMark1OpenMarker  NodeType = 451 // 开始标记标记符 =
	NodeMark1CloseMarker NodeType = 452 // 结束标记标记符 =
	NodeMark2OpenMarker  NodeType = 453 // 开始标记标记符 ==
	NodeMark2CloseMarker NodeType = 454 // 结束标记标记符 ==

	// kramdown 内联属性列表 https://github.com/88250/lute/issues/89 and https://github.com/88250/lute/issues/118

	NodeKramdownBlockIAL NodeType = 455 // 块级内联属性列表 {: name="value"}
	NodeKramdownSpanIAL  NodeType = 456 // 行级内联属性列表 *foo*{: name="value"}bar

	// #Tag# 标签语法 https://github.com/88250/lute/issues/92

	NodeTag            NodeType = 460 // 标签
	NodeTagOpenMarker  NodeType = 461 // 开始标签标记符 #
	NodeTagCloseMarker NodeType = 462 // 结束标签标记符 #

	// 内容块查询嵌入（Block Query Embed）语法 https://github.com/88250/lute/issues/96

	NodeBlockQueryEmbed       NodeType = 465 // 内容块查询嵌入
	NodeOpenBrace             NodeType = 466 // {
	NodeCloseBrace            NodeType = 467 // }
	NodeBlockQueryEmbedScript NodeType = 468 // 内容块查询嵌入脚本

	// 超级块语法 https://github.com/88250/lute/issues/111

	NodeSuperBlock             NodeType = 475 // 超级块节点
	NodeSuperBlockOpenMarker   NodeType = 476 // 开始超级块标记符 {{{
	NodeSuperBlockLayoutMarker NodeType = 477 // 超级块布局 row/col
	NodeSuperBlockCloseMarker  NodeType = 478 // 结束超级块标记符 }}}

	// 上标下标语法 https://github.com/88250/lute/issues/113

	NodeSup            NodeType = 485 // 上标
	NodeSupOpenMarker  NodeType = 486 // 开始上标标记符 ^
	NodeSupCloseMarker NodeType = 487 // 结束上标标记符 ^
	NodeSub            NodeType = 490 // 下标
	NodeSubOpenMarker  NodeType = 491 // 开始下标标记符 ~
	NodeSubCloseMarker NodeType = 492 // 结束下标标记符 ~

	// Git 冲突标记 https://github.com/88250/lute/issues/131

	NodeGitConflict            NodeType = 495 // Git 冲突标记
	NodeGitConflictOpenMarker  NodeType = 496 // 开始 Git 冲突标记标记符 <<<<<<<
	NodeGitConflictContent     NodeType = 497 // Git 冲突标记内容
	NodeGitConflictCloseMarker NodeType = 498 // 结束 Git 冲突标记标记符 >>>>>>>

	// <iframe> 标签

	NodeIFrame NodeType = 500 // <iframe> 标签

	// <audio> 标签

	NodeAudio NodeType = 505 // <audio> 标签

	// <video> 标签

	NodeVideo NodeType = 510 // <video> 标签

	NodeTypeMaxVal NodeType = 1024 // 节点类型最大值
)
