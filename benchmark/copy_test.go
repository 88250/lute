// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package main

import (
	"strconv"
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
)

var (
	copyEngine    *lute.Lute
	copySmallDOM  string
	copyNestedDOM string
	copyLooseDOM  string
	copyDeepDOM   string
	copyMixedDOM  string
)

func copyNestedMarkdown(top, sub int) string {
	buf := &strings.Builder{}
	for i := 0; i < top; i++ {
		buf.WriteString("- 条目 ")
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString("\n")
		for j := 0; j < sub; j++ {
			buf.WriteString("  - 子条目 ")
			buf.WriteString(strconv.Itoa(j))
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func copyLooseMarkdown(top int) string {
	buf := &strings.Builder{}
	for i := 0; i < top; i++ {
		buf.WriteString("- 条目 ")
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString("\n\n  补充说明 ")
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString("\n\n")
	}
	return buf.String()
}

func copyDeepMarkdown(top int) string {
	buf := &strings.Builder{}
	for i := 0; i < top; i++ {
		buf.WriteString("1. 一级 ")
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString("\n   1. 二级\n      - 三级\n        - 四级\n")
	}
	return buf.String()
}

func copyMixedMarkdown(top int) string {
	buf := &strings.Builder{}
	for i := 0; i < top; i++ {
		buf.WriteString("- 条目 ")
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString("\n  - 子条目 A\n  - 子条目 B\n")
		if 0 == i%5 {
			buf.WriteString("  ```go\n  x := 1\n  ```\n")
		}
	}
	buf.WriteString("\n正文段落。\n")
	return buf.String()
}

func init() {
	copyEngine = lute.New()
	copyEngine.SetProtyleWYSIWYG(true)
	copyEngine.SetDataTask(true)
	copyEngine.SetArbitraryTaskListItemMarker(true)
	copyEngine.SetUnorderedListMarker("-")
	copySmallDOM = copyEngine.Md2BlockDOM("1. 111\n   1. 222\n   2. 333\n2. 444\n   1. 555\n   2. 666\n3. 777\n", true)
	copyNestedDOM = copyEngine.Md2BlockDOM(copyNestedMarkdown(200, 3), true)
	copyLooseDOM = copyEngine.Md2BlockDOM(copyLooseMarkdown(100), true)
	copyDeepDOM = copyEngine.Md2BlockDOM(copyDeepMarkdown(50), true)
	copyMixedDOM = copyEngine.Md2BlockDOM(copyMixedMarkdown(200), true)
}

// BenchmarkCopySmall 测量小块复制的固定开销。
func BenchmarkCopySmall(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = copyEngine.BlockDOM2StdMd(copySmallDOM)
	}
}

// BenchmarkCopyNested 测量大量可紧凑嵌套列表的复制开销。
func BenchmarkCopyNested(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = copyEngine.BlockDOM2StdMd(copyNestedDOM)
	}
}

// BenchmarkCopyLoose 测量必须使用松散形式的列表复制开销。
func BenchmarkCopyLoose(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = copyEngine.BlockDOM2StdMd(copyLooseDOM)
	}
}

// BenchmarkCopyDeep 测量多层级嵌套列表的复制开销。
func BenchmarkCopyDeep(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = copyEngine.BlockDOM2StdMd(copyDeepDOM)
	}
}

// BenchmarkCopyMixed 测量列表与代码块混合内容的复制开销。
func BenchmarkCopyMixed(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = copyEngine.BlockDOM2StdMd(copyMixedDOM)
	}
}

// BenchmarkExportListNodes 覆盖内核逐块索引时逐个导出列表项和段落的路径。
// 每次迭代导出整份列表中的全部节点，以观察总耗时随列表规模的增长。
func BenchmarkExportListNodes(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		for _, typ := range []ast.NodeType{ast.NodeListItem, ast.NodeParagraph} {
			b.Run(strconv.Itoa(size)+"/"+typ.String(), func(b *testing.B) {
				engine := lute.New()
				tree := parse.Parse("", []byte(strings.Repeat("- item\n", size)), engine.ParseOptions)
				var nodes []*ast.Node
				for item := tree.Root.FirstChild.FirstChild; nil != item; item = item.Next {
					if ast.NodeListItem == typ {
						nodes = append(nodes, item)
					} else {
						nodes = append(nodes, item.FirstChild)
					}
				}
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					for _, node := range nodes {
						if _, err := lute.ProtyleExportMdNodeSync(node, engine.ParseOptions, engine.RenderOptions); nil != err {
							b.Fatal(err)
						}
					}
				}
			})
		}
	}
}
