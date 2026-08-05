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
)

// blockStartsFuncs 是定义好的块起始模式函数切片，每个函数用于判断某种块节点是否可以开始。
// 该切片是只读的，每个行解析都会引用它，所以复用同一个切片以避免每次分配。
// 使用 init() 初始化以避免包级变量初始化顺序导致的循环引用。
var blockStartsFuncs []blockStartFunc

func init() {
	blockStartsFuncs = []blockStartFunc{
		GitConflictStart,
		CalloutStart,
		BlockquoteStart,
		ATXHeadingStart,
		FenceCodeBlockStart,
		// CustomBlockStart, // https://github.com/siyuan-note/siyuan/issues/8418
		SetextHeadingStart,
		HtmlBlockStart,
		YamlFrontMatterStart,
		ThematicBreakStart,
		ListStart,
		MathBlockStart,
		IndentCodeBlockStart,
		FootnotesStart,
		IALStart,
		BlockQueryEmbedStart,
		SuperBlockStart,
	}
}

// blockStarts 返回块起始模式函数切片。
func blockStarts() []blockStartFunc {
	return blockStartsFuncs
}

// blockStartFunc 定义了用于判断块是否开始的函数签名，返回值：
//
//	0：不匹配
//	1：匹配到容器块，需要继续迭代下降
//	2：匹配到叶子块
type blockStartFunc func(t *Tree, container *ast.Node) int
