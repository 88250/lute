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
	"slices"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
)

func TestConvertListType(t *testing.T) {
	luteEngine := newRecursiveListLute()
	input := luteEngine.Md2BlockDOM("1. outer\n   1. nested\n      1. deep\n   - [ ] task\n     1. ordered in task", false)
	converted := luteEngine.ConvertListType(input, "u")

	if got, want := listTypes(luteEngine, converted), []int{0, 0, 0, 3, 0}; !slices.Equal(got, want) {
		t.Fatalf("unexpected list types: got %v, want %v", got, want)
	}
	if got := taskMarkerCount(luteEngine, converted); 1 != got {
		t.Fatalf("unexpected task marker count: got %d, want 1", got)
	}
	if got := luteEngine.ConvertListType(input, "invalid"); input != got {
		t.Fatalf("invalid target type should keep the original HTML")
	}
}

func TestConvertListTypePairs(t *testing.T) {
	tests := []struct {
		name       string
		markdown   string
		targetType string
		targetTyp  int
	}{
		{"unordered to ordered", "- outer\n  - nested", "o", 1},
		{"unordered to task", "- outer\n  - nested", "t", 3},
		{"ordered to unordered", "1. outer\n   1. nested", "u", 0},
		{"ordered to task", "1. outer\n   1. nested", "t", 3},
		{"task to unordered", "- [ ] outer\n  - [x] nested", "u", 0},
		{"task to ordered", "- [ ] outer\n  - [x] nested", "o", 1},
	}

	luteEngine := newRecursiveListLute()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := luteEngine.Md2BlockDOM(test.markdown, false)
			converted := luteEngine.ConvertListType(input, test.targetType)
			if got, want := listTypes(luteEngine, converted), []int{test.targetTyp, test.targetTyp}; !slices.Equal(got, want) {
				t.Fatalf("unexpected list types: got %v, want %v", got, want)
			}
			if 3 == test.targetTyp {
				if got := taskMarkerCount(luteEngine, converted); 2 != got {
					t.Fatalf("unexpected task marker count: got %d, want 2", got)
				}
			} else if got := taskMarkerCount(luteEngine, converted); 0 != got {
				t.Fatalf("unexpected task marker count: got %d, want 0", got)
			}
			if 1 == test.targetTyp {
				assertOrderedListNumbers(t, luteEngine, converted)
			}
		})
	}
}

func TestConvertTaskListType(t *testing.T) {
	luteEngine := newRecursiveListLute()
	input := luteEngine.Md2BlockDOM("- [x] outer\n  - [ ] nested\n    1. ordered\n       - [ ] deep task", false)
	converted := luteEngine.ConvertListType(input, "o")

	if got, want := listTypes(luteEngine, converted), []int{1, 1, 1, 1}; !slices.Equal(got, want) {
		t.Fatalf("unexpected list types: got %v, want %v", got, want)
	}
	if got := taskMarkerCount(luteEngine, converted); 0 != got {
		t.Fatalf("unexpected task marker count: got %d, want 0", got)
	}
	assertOrderedListNumbers(t, luteEngine, converted)
}

func TestCancelListRecursively(t *testing.T) {
	luteEngine := newRecursiveListLute()
	input := luteEngine.Md2BlockDOM("- outer\n  - nested\n    1. ordered\n       - deep unordered\n- tail", false)

	shallow := luteEngine.CancelList(input)
	if got, want := listTypes(luteEngine, shallow), []int{0, 1, 0}; !slices.Equal(got, want) {
		t.Fatalf("unexpected shallow list types: got %v, want %v", got, want)
	}

	converted := luteEngine.CancelListRecursively(input)
	if got, want := listTypes(luteEngine, converted), []int{1}; !slices.Equal(got, want) {
		t.Fatalf("unexpected recursive list types: got %v, want %v", got, want)
	}
	if got, want := paragraphTexts(luteEngine, converted), []string{"outer", "nested", "ordered", "deep unordered", "tail"}; !slices.Equal(got, want) {
		t.Fatalf("unexpected paragraph order: got %v, want %v", got, want)
	}
}

func newRecursiveListLute() *lute.Lute {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetDataTask(true)
	luteEngine.SetKramdownIAL(true)
	return luteEngine
}

func listTypes(luteEngine *lute.Lute, html string) (ret []int) {
	tree := luteEngine.BlockDOM2Tree(html)
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering && ast.NodeList == n.Type {
			ret = append(ret, n.ListData.Typ)
		}
		return ast.WalkContinue
	})
	return
}

func taskMarkerCount(luteEngine *lute.Lute, html string) (ret int) {
	tree := luteEngine.BlockDOM2Tree(html)
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering && ast.NodeTaskListItemMarker == n.Type {
			ret++
		}
		return ast.WalkContinue
	})
	return
}

func assertOrderedListNumbers(t *testing.T, luteEngine *lute.Lute, html string) {
	t.Helper()
	tree := luteEngine.BlockDOM2Tree(html)
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering || ast.NodeList != n.Type || 1 != n.ListData.Typ {
			return ast.WalkContinue
		}
		num := 1
		for li := n.FirstChild; nil != li; li = li.Next {
			if ast.NodeListItem != li.Type {
				continue
			}
			if li.ListData.Num != num {
				t.Fatalf("unexpected ordered list number: got %d, want %d", li.ListData.Num, num)
			}
			num++
		}
		return ast.WalkContinue
	})
}

func paragraphTexts(luteEngine *lute.Lute, html string) (ret []string) {
	tree := luteEngine.BlockDOM2Tree(html)
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering && ast.NodeParagraph == n.Type {
			ret = append(ret, n.Text())
		}
		return ast.WalkContinue
	})
	return
}
