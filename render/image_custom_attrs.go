package render

import (
	"strings"

	"github.com/88250/lute/ast"
)

// imageCustomAttrs 保留图片的自定义行级属性，属性值沿用 AST 中已转义的表示。
func imageCustomAttrs(node *ast.Node) (attrs [][]string) {
	for _, kv := range node.KramdownIAL {
		if len(kv) != 2 || !strings.HasPrefix(kv[0], "custom-") {
			continue
		}
		valid := true
		for _, c := range kv[0] {
			if c <= ' ' || c == 127 || strings.ContainsRune("\"'<>/=", c) {
				valid = false
				break
			}
		}
		if valid {
			attrs = append(attrs, kv)
		}
	}
	return
}
