package parse

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
)

func TestSetSpanIALStableCustomAttributeOrder(t *testing.T) {
	for _, nodeType := range []ast.NodeType{ast.NodeImage, ast.NodeTableCell} {
		for _, style := range []string{"", "color: red;"} {
			t.Run(fmt.Sprintf("%s/style=%t", nodeType, style != ""), func(t *testing.T) {
				dom := &html.Node{Type: html.ElementNode}
				var want [][]string
				if style != "" {
					dom.Attr = append(dom.Attr, &html.Attribute{Key: "style", Val: style})
					want = append(want, []string{"style", StyleValue(style)})
				}
				// Reverse DOM order also verifies that output is canonical, rather
				// than dependent on map iteration or the source attribute order.
				for i := 15; i >= 0; i-- {
					dom.Attr = append(dom.Attr, &html.Attribute{Key: fmt.Sprintf("custom-%02d", i), Val: "old"})
				}
				for i := 0; i < 16; i++ {
					key := fmt.Sprintf("custom-%02d", i)
					// Preserve last-value-wins and entity unescaping semantics.
					dom.Attr = append(dom.Attr, &html.Attribute{Key: key, Val: "a&amp;b"})
					want = append(want, []string{key, "a&amp;b"})
				}
				for repeat := 0; repeat < 100; repeat++ {
					parent := &ast.Node{Type: ast.NodeParagraph}
					node := &ast.Node{Type: nodeType}
					parent.AppendChild(node)
					SetSpanIAL(node, dom)
					if !reflect.DeepEqual(node.KramdownIAL, want) {
						t.Fatalf("iteration %d: attributes = %v, want %v", repeat, node.KramdownIAL, want)
					}
					ial := node.Next
					if nodeType == ast.NodeTableCell {
						ial = node.FirstChild
					}
					if ial == nil || ial.Type != ast.NodeKramdownSpanIAL || string(ial.Tokens) != string(IAL2Tokens(want)) {
						t.Fatalf("iteration %d: span IAL does not match canonical attributes", repeat)
					}
				}
			})
		}
	}
}
