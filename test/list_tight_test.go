// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func newListTightEngine() *lute.Lute {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetDataTask(true)
	engine.SetArbitraryTaskListItemMarker(true)
	engine.SetUnorderedListMarker("-")
	return engine
}

// listLayout 输出块层级与种类组成的结构签名，忽略列表松紧标识和任务标记取值。
func listLayout(engine *lute.Lute, markdown string) string {
	tree := parse.Parse("", []byte(markdown), engine.ParseOptions)
	buf := &strings.Builder{}
	depth := 0
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			if ast.NodeKramdownBlockIAL != n.Type {
				depth--
			}
			return ast.WalkContinue
		}
		if ast.NodeKramdownBlockIAL == n.Type {
			return ast.WalkContinue
		}
		indent := strings.Repeat("-", depth)
		switch n.Type {
		case ast.NodeList:
			typ := 0
			if nil != n.ListData {
				typ = n.ListData.Typ
			}
			fmt.Fprintf(buf, "%sL%d;", indent, typ)
		case ast.NodeListItem, ast.NodeTaskListItemMarker:
			fmt.Fprintf(buf, "%sI;", indent)
		case ast.NodeParagraph:
			fmt.Fprintf(buf, "%sP[%s];", indent, strings.TrimSpace(n.Text()))
		case ast.NodeCodeBlock, ast.NodeTable, ast.NodeBlockquote, ast.NodeCallout, ast.NodeMathBlock, ast.NodeHTMLBlock, ast.NodeThematicBreak:
			fmt.Fprintf(buf, "%s%v;", indent, n.Type)
		case ast.NodeHeading:
			fmt.Fprintf(buf, "%sH%d;", indent, n.HeadingLevel)
		}
		depth++
		return ast.WalkContinue
	})
	return buf.String()
}

// TestListTightExport 校验复制为 Markdown 时的列表松紧：默认紧凑，
// 只有紧凑形式会改变结构时才使用松散形式，松散形式在块之间和列表项之间各留一个空行。
func TestListTightExport(t *testing.T) {
	engine := newListTightEngine()

	for _, test := range []struct {
		name     string
		from, to string
	}{
		{
			"可紧凑的嵌套列表",
			"- 111\n  - 222\n  - 333\n- 444\n  - 555\n  - 666\n",
			"- 111\n  - 222\n  - 333\n- 444\n  - 555\n  - 666\n",
		},
		{
			"有序列表可紧凑",
			"1. 111\n   1. 222\n   2. 333\n2. 444\n   1. 555\n   2. 666\n",
			"1. 111\n   1. 222\n   2. 333\n2. 444\n   1. 555\n   2. 666\n",
		},
		{
			"任务列表可紧凑",
			"- [ ] 111\n  - [x] 222\n  - [ ] 333\n- [ ] 444\n",
			"- [ ] 111\n  - [X] 222\n  - [ ] 333\n- [ ] 444\n",
		},
		{
			"非 1 起始的有序子列表必须松散",
			"- 111\n\n  2. 222\n  3. 333\n- 444\n",
			"- 111\n\n  2. 222\n  3. 333\n\n- 444\n",
		},
		{
			"列表项内多个段落必须松散",
			"- 111\n\n  补充说明\n- 444\n",
			"- 111\n\n  补充说明\n\n- 444\n",
		},
		{
			"混合编号的列表按松散形式统一排版",
			"1. 111\n   1. 222\n   2. 333\n2. 444\n\n   2. 555\n   3. 666\n\n3. 777\n   1. 888\n   2. 999\n",
			"1. 111\n\n   1. 222\n   2. 333\n\n2. 444\n\n   2. 555\n   3. 666\n\n3. 777\n\n   1. 888\n   2. 999\n",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := engine.BlockDOM2StdMd(engine.Md2BlockDOM(test.from, true))
			if test.to != got {
				t.Fatalf("expected\n\t%q\ngot\n\t%q", test.to, got)
			}
			// 导出结果重新解析后必须保持原有层级和内容
			if expected, actual := listLayout(engine, test.from), listLayout(engine, got); expected != actual {
				t.Fatalf("structure changed\nexpected\n\t%q\ngot\n\t%q\nmarkdown\n\t%q", expected, actual, got)
			}
		})
	}
}

// TestListTightExportKeepsSource 校验导出不会修改源树的列表松紧标识，
// 且对同一棵树重复渲染得到相同结果。
func TestListTightExportKeepsSource(t *testing.T) {
	engine := newListTightEngine()
	markdown := "- 111\n  - 222\n  - 333\n- 444\n\n  补充说明\n"

	tree := engine.BlockDOM2Tree(engine.Md2BlockDOM(markdown, true))
	var lists []*ast.Node
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering && nil != n.ListData && (ast.NodeList == n.Type || ast.NodeListItem == n.Type) {
			lists = append(lists, n)
		}
		return ast.WalkContinue
	})
	if 0 == len(lists) {
		t.Fatal("no list nodes found")
	}

	before := make([]bool, len(lists))
	for i, node := range lists {
		before[i] = node.ListData.Tight
	}

	first := string(render.NewProtyleExportMdRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render())
	for i, node := range lists {
		if before[i] != node.ListData.Tight {
			t.Fatalf("export mutated ListData.Tight: was %v, now %v", before[i], node.ListData.Tight)
		}
	}

	second := string(render.NewProtyleExportMdRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render())
	if first != second {
		t.Fatalf("repeated export changed result\nfirst\n\t%q\nsecond\n\t%q", first, second)
	}
}

// TestListTightExportByNode 校验按节点渲染（内核 ExportNodeStdMd 路径）与整树渲染结果一致。
func TestListTightExportByNode(t *testing.T) {
	engine := newListTightEngine()
	const markdown = "- 111\n  - 222\n  - 333\n- 444\n"

	tree := engine.BlockDOM2Tree(engine.Md2BlockDOM(markdown, true))
	list := tree.Root.FirstChild
	if ast.NodeList != list.Type {
		t.Fatalf("expected a list, got %v", list.Type)
	}

	whole, err := lute.ProtyleExportMdNodeSync(list, engine.ParseOptions, engine.RenderOptions)
	if nil != err {
		t.Fatalf("export list failed: %s", err)
	}
	if expected := "- 111\n  - 222\n  - 333\n- 444"; expected != whole {
		t.Fatalf("expected\n\t%q\ngot\n\t%q", expected, whole)
	}

	// 只渲染首个列表项时也必须按紧凑形式输出，不能退回松散形式
	item, err := lute.ProtyleExportMdNodeSync(list.FirstChild, engine.ParseOptions, engine.RenderOptions)
	if nil != err {
		t.Fatalf("export item failed: %s", err)
	}
	if expected := "- 111\n  - 222\n  - 333"; expected != item {
		t.Fatalf("expected\n\t%q\ngot\n\t%q", expected, item)
	}
}

func TestListTightExportBlockBoundaries(t *testing.T) {
	engine := newListTightEngine()
	engine.SetCallout(true)
	for _, test := range []struct {
		name, markdown string
	}{
		{"thematic break after paragraph", "- text\n\n  ---\n- tail\n"},
		{"empty nested item", "- text\n\n  -\n- tail\n"},
		{"empty first nested item", "- text\n\n  -\n  - child\n- tail\n"},
		{"adjacent blockquotes", "- > first\n\n  > second\n- tail\n"},
		{"blockquote then callout", "- > quote\n\n  > [!NOTE]\n  > body\n- tail\n"},
		{"callout then blockquote", "- > [!NOTE]\n  > body\n\n  > quote\n- tail\n"},
		{"adjacent callouts", "- > [!NOTE]\n  > first\n\n  > [!TIP]\n  > second\n- tail\n"},
		{"HTML then paragraph", "- <div>\n  html\n  </div>\n\n  text\n- tail\n"},
		{"HTML then code", "- <div>\n  html\n  </div>\n\n  ```\n  code\n  ```\n- tail\n"},
		{"HTML then list", "- <div>\n  html\n  </div>\n\n  - child\n- tail\n"},
		{"table then paragraph", "- | a | b |\n  | - | - |\n  | c | d |\n\n  text\n- tail\n"},
		{"adjacent tables", "- | a | b |\n  | - | - |\n  | c | d |\n\n  | e | f |\n  | - | - |\n  | g | h |\n- tail\n"},
		{"nested list then paragraph", "- - child\n\n  text\n- tail\n"},
		{"nested list then thematic break", "- - child\n\n  ---\n- tail\n"},
		{"nested list then table", "- - child\n\n  | a | b |\n  | - | - |\n  | c | d |\n- tail\n"},
		{"task heading then table", "- [ ] # heading\n\n  | a | b |\n  | - | - |\n  | c | d |\n- tail\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			dom := engine.Md2BlockDOM(test.markdown, true)
			expected := listLayout(engine, test.markdown)
			for _, path := range []string{"copy", "tree", "node", "raw"} {
				t.Run(path, func(t *testing.T) {
					var got string
					switch path {
					case "copy":
						got = engine.BlockDOM2StdMd(dom)
					case "tree":
						got = string(render.NewProtyleExportMdRenderer(engine.BlockDOM2Tree(dom), engine.RenderOptions, engine.ParseOptions).Render())
					case "node":
						var err error
						got, err = lute.ProtyleExportMdNodeSync(engine.BlockDOM2Tree(dom).Root, engine.ParseOptions, engine.RenderOptions)
						if nil != err {
							t.Fatal(err)
						}
					case "raw":
						tree := parse.Parse("", []byte(test.markdown), engine.ParseOptions)
						got = string(render.NewProtyleExportMdRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render())
					}
					if actual := listLayout(engine, got); expected != actual {
						t.Fatalf("structure changed\nexpected %q\nactual %q\nmarkdown %q", expected, actual, got)
					}
				})
			}
		})
	}
}

func TestListTightExportFlattenedSuperBlocks(t *testing.T) {
	engine := newListTightEngine()
	engine.SetSuperBlock(true)
	for _, markdown := range []string{
		"- before\n\n  {{{row\n  inside\n  }}}\n- tail\n",
		"- before\n\n  {{{row\n  {{{col\n  inside\n  }}}\n  }}}\n- tail\n",
	} {
		got := engine.BlockDOM2StdMd(engine.Md2BlockDOM(markdown, true))
		const expected = "- before\n\n  inside\n\n- tail\n"
		if got != expected {
			t.Fatalf("expected %q, got %q", expected, got)
		}
		if listLayout(engine, got) != listLayout(engine, expected) {
			t.Fatalf("flattened content changed structure: %q", got)
		}
	}
}

func TestListTightExportKramdownBlockBoundaries(t *testing.T) {
	engine := newListTightEngine()
	engine.SetKramdownIAL(true)
	engine.SetTextMark(true)
	for _, markdown := range []string{
		"- text\n\n  -\n  - child\n- tail\n",
		"- [ ] # heading\n\n  ---\n- tail\n",
		"- [ ] # heading\n\n  -\n- tail\n",
		"- [ ] # heading\n\n  -\n  - child\n- tail\n",
	} {
		expected := listLayout(engine, markdown)
		got := engine.BlockDOM2StdMd(engine.Md2BlockDOM(markdown, true))
		if actual := listLayout(engine, got); expected != actual {
			t.Fatalf("structure changed\nexpected %q\nactual %q\nmarkdown %q", expected, actual, got)
		}
	}
}

func TestListTightExportByNodeUsesRenderedItem(t *testing.T) {
	engine := newListTightEngine()
	tree := engine.BlockDOM2Tree(engine.Md2BlockDOM("- parent\n  - child\n- sibling\n\n  second paragraph\n", true))
	item := tree.Root.FirstChild.FirstChild
	got, err := lute.ProtyleExportMdNodeSync(item, engine.ParseOptions, engine.RenderOptions)
	if nil != err {
		t.Fatal(err)
	}
	if expected := "- parent\n  - child"; got != expected {
		t.Fatalf("unrendered sibling affected item: expected %q, got %q", expected, got)
	}
	paragraph, err := lute.ProtyleExportMdNodeSync(item.FirstChild, engine.ParseOptions, engine.RenderOptions)
	if nil != err || paragraph != "parent" {
		t.Fatalf("unexpected paragraph export %q, error %v", paragraph, err)
	}
}
