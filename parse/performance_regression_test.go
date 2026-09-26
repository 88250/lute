package parse

import (
	"bytes"
	"slices"
	"strings"
	"testing"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
)

func TestTableCandidateCompatibility(t *testing.T) {
	context := &Context{ParseOption: NewOptions()}
	if table := context.parseTable0([]byte("foo\n---\nbar")); table == nil {
		t.Fatal("expected a single-column table without pipes")
	}
	for _, tc := range []struct {
		input  string
		setext bool
		want   ast.NodeType
		rows   int
	}{
		{"foo\n---", false, ast.NodeParagraph, 0},
		{"foo\n---", true, ast.NodeHeading, 0},
		{"foo\n-:", false, ast.NodeTable, 1},
		{"foo\n::\nbar", false, ast.NodeTable, 2},
		{"a | b\n--- | ---\nx | y", false, ast.NodeTable, 2},
		{"a | b\n---\nx | y", false, ast.NodeParagraph, 0},
		{"alpha | beta\ngamma | delta", false, ast.NodeParagraph, 0},
		{"前文\na | b\n--- | ---\nx | y", false, ast.NodeParagraph, 2},
	} {
		t.Run(tc.input, func(t *testing.T) {
			o := NewOptions()
			o.Setext = tc.setext
			tree := Parse("", []byte(tc.input), o)
			if tree.Root.FirstChild.Type != tc.want {
				t.Fatalf("unexpected first block: %v", tree.Root.FirstChild.Type)
			}
			rows := 0
			ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
				if entering && n.Type == ast.NodeTableRow {
					rows++
				}
				return ast.WalkContinue
			})
			if rows != tc.rows {
				t.Fatalf("want %d rows, got %d", tc.rows, rows)
			}
		})
	}
}

func TestAutoLinkCandidateCompatibility(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  []string
	}{
		{strings.Repeat("a", 4096), nil},
		{"foo://", nil},
		{"foo://bar", []string{"foo://bar"}},
		{"中文foo://bar", []string{"foo://bar"}},
		{"abc-def://host", []string{"def://host"}},
		{"prefix: foo://bar", []string{"foo://bar"}},
		{"foo://bar baz://qux", []string{"baz://qux"}},
		{"foo://bar https://example.com", []string{"https://example.com"}},
		{"https://example.com foo://bar", []string{"https://example.com", "foo://bar"}},
		{"www.example.com", []string{"http://www.example.com"}},
		{"http://example.com/a).", []string{"http://example.com/a)"}},
		{"https://example.com:8080/a:b ftp://example.com/path", []string{"https://example.com:8080/a:b", "ftp://example.com/path"}},
	} {
		t.Run(tc.input, func(t *testing.T) {
			o := NewOptions()
			o.Emoji = false
			tree := Parse("", []byte(tc.input), o)
			var got []string
			ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
				if entering && n.Type == ast.NodeLinkDest {
					got = append(got, string(n.Tokens))
				}
				return ast.WalkContinue
			})
			if !slices.Equal(got, tc.want) {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}

func TestTaskListProbeCompatibility(t *testing.T) {
	for _, protyle := range []bool{false, true} {
		for _, tc := range []struct {
			input string
			types map[ast.NodeType]int
		}{
			{"- [ ] alpha **beta** and `code`\n", map[ast.NodeType]int{ast.NodeList: 1, ast.NodeParagraph: 1, ast.NodeStrong: 1, ast.NodeCodeSpan: 1}},
			{"- [ ] # heading\n", map[ast.NodeType]int{ast.NodeHeading: 1}},
			{"- [ ] > quote\n", map[ast.NodeType]int{ast.NodeBlockquote: 1, ast.NodeParagraph: 1}},
			{"- [ ] - nested\n", map[ast.NodeType]int{ast.NodeList: 2, ast.NodeParagraph: 1}},
			{"- [ ] ```go\n  code\n  ```\n", map[ast.NodeType]int{ast.NodeCodeBlock: 2}},
		} {
			o := NewOptions()
			o.ProtyleWYSIWYG = protyle
			o.KramdownBlockIAL, o.KramdownSpanIAL = protyle, protyle
			tree := Parse("", []byte(tc.input), o)
			counts := map[ast.NodeType]int{}
			ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
				if entering {
					counts[n.Type]++
				}
				return ast.WalkContinue
			})
			if counts[ast.NodeTaskListItemMarker] != 1 {
				t.Fatalf("task marker lost: %q", tc.input)
			}
			for typ, want := range tc.types {
				if counts[typ] != want {
					t.Fatalf("protyle=%t input=%q type=%v: want %d, got %d", protyle, tc.input, typ, want, counts[typ])
				}
			}
		}
	}
}

func TestReferenceIndexCompatibility(t *testing.T) {
	for _, tc := range []struct{ definitions, label, want string }{
		{"[Foo]: /first\n[foo]: /second\n", "FOO", "/first"},
		{"[K]: /first\n[K]: /second\n", "k", "/first"},
		{"[ſ]: /first\n[S]: /second\n", "s", "/first"},
		{"[ς]: /first\n[Σ]: /second\n", "σ", "/first"},
		{"[ẞ]: /first\n[SS]: /second\n", "ss", "/first"},
		{"[ss]: /first\n[ẞ]: /second\n", "ẞ", "/first"},
		{"[ßs]: /wrong\n[sß]: /right\n", "sß", "/right"},
		{"[ßs]: /wrong\n", "sß", ""},
		{"[foo]: /first\n", "missing", ""},
	} {
		t.Run(tc.label+tc.definitions, func(t *testing.T) {
			tree := Parse("", []byte(tc.definitions), NewOptions())
			want := tc.want
			// JS 版仅执行简单大小写折叠，保留其独立的匹配规则。
			if bytes.Equal(foldBytes([]byte("ẞ")), []byte("ẞ")) && (tc.label == "ss" || tc.label == "ẞ") {
				want = "/second"
			}
			for _, protyle := range []bool{false, true} {
				tree.Context.ParseOption.ProtyleWYSIWYG = protyle
				label := tc.label
				if protyle {
					label += editor.Caret
				}
				link := tree.FindLinkRefDefLink([]byte(label))
				got := ""
				if link != nil {
					got = string(link.ChildByType(ast.NodeLinkDest).Tokens)
				}
				if got != want {
					t.Fatalf("want %q, got %q", want, got)
				}
			}
		})
	}
}

func TestLinkLabelBoundaries(t *testing.T) {
	context := &Context{ParseOption: NewOptions()}
	for _, tc := range []struct {
		input, want, remains string
		consumed             int
	}{
		{"[中文] trailing", "中文", " trailing", 8},
		{"[a\\]b] tail", "a\\]b", " tail", 6},
		{"[a\n  b] tail", "a  b", " tail", 7},
		{"[" + strings.Repeat("a", 999) + "] tail", strings.Repeat("a", 999), " tail", 1001},
	} {
		n, remains, label := context.parseLinkLabel([]byte(tc.input))
		if n != tc.consumed || string(label) != tc.want || string(remains) != tc.remains {
			t.Fatalf("%q: got %d, %q, %q", tc.input, n, label, remains)
		}
	}
	for _, input := range []string{"[a[b]", "[unclosed", "[" + strings.Repeat("a", 1000) + "]"} {
		if n, _, _ := context.parseLinkLabel([]byte(input)); n != 0 {
			t.Fatalf("accepted %q", input)
		}
	}
}
