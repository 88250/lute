package parse

import (
	"html"
	"strconv"
	"strings"

	"github.com/88250/lute/ast"
	lutehtml "github.com/88250/lute/html"
)

// protectTableCellRichStyles 在解析前保护样式实体，解析后区分样式值与代码中的字面内容。
func protectTableCellRichStyles(source string) (string, func(*Tree)) {
	sentinel := tableCellRichStyleSentinel(source)
	var styles, literals, codeLiterals, mathLiterals []string
	var out strings.Builder
	cursor := 0
	for cursor < len(source) {
		start := strings.Index(source[cursor:], `</span>{:`)
		if start < 0 {
			out.WriteString(source[cursor:])
			break
		}
		start += cursor + len(`</span>{:`)
		position := start
		for position < len(source) && strings.ContainsRune(" \t\r\n", rune(source[position])) {
			position++
		}
		if !strings.HasPrefix(source[position:], "style=") || position+7 > len(source) ||
			(source[position+6] != '\'' && source[position+6] != '"') {
			out.WriteString(source[cursor:start])
			cursor = start
			continue
		}
		quote := source[position+6]
		start = position + 7
		end := strings.IndexByte(source[start:], quote)
		if end < 0 {
			out.WriteString(source[cursor:])
			break
		}
		end += start
		out.WriteString(source[cursor:start])
		value := source[start:end]
		for len(value) > 0 {
			index := strings.IndexAny(value, "&`")
			if index < 0 {
				out.WriteString(value)
				break
			}
			out.WriteString(value[:index])
			value = value[index:]
			length := 1
			if value[0] == '&' {
				if semicolon := strings.IndexByte(value[:min(len(value), 32)], ';'); semicolon >= 0 {
					length = semicolon + 1
				}
			}
			literal := value[:length]
			token := sentinel + strconv.Itoa(len(styles)) + "\ue001"
			styles = append(styles, token, html.UnescapeString(literal))
			literals = append(literals, token, literal)
			code := strings.ReplaceAll(literal, "&", "&amp;")
			codeLiterals = append(codeLiterals, token, code)
			mathLiterals = append(mathLiterals, token, strings.ReplaceAll(code, "&", "&amp;"))
			out.WriteString(token)
			value = value[length:]
		}
		cursor = end
	}
	return out.String(), func(tree *Tree) {
		if len(styles) == 0 {
			return
		}
		styleReplacer, literalReplacer := strings.NewReplacer(styles...), strings.NewReplacer(literals...)
		codeReplacer, mathReplacer := strings.NewReplacer(codeLiterals...), strings.NewReplacer(mathLiterals...)
		ast.Walk(tree.Root, func(node *ast.Node, entering bool) ast.WalkStatus {
			if !entering {
				return ast.WalkContinue
			}
			for _, attr := range node.KramdownIAL {
				if attr[0] == "style" && node.IsTextMarkType("text") {
					attr[1] = lutehtml.EscapeAttrVal(styleReplacer.Replace(lutehtml.UnescapeAttrVal(attr[1])))
				} else {
					attr[1] = literalReplacer.Replace(attr[1])
				}
			}
			node.Tokens = []byte(literalReplacer.Replace(string(node.Tokens)))
			node.TextMarkTextContent = codeReplacer.Replace(node.TextMarkTextContent)
			node.TextMarkInlineMathContent = mathReplacer.Replace(node.TextMarkInlineMathContent)
			if ast.NodeKramdownSpanIAL == node.Type && nil != node.Previous && node.Previous.IsTextMarkType("text") {
				node.Tokens = IAL2Tokens(node.Previous.KramdownIAL)
			}
			return ast.WalkContinue
		})
	}
}

func tableCellRichStyleSentinel(source string) string {
	const prefix = "\ue000table-cell-style-"
	used := map[string]bool{}
	for cursor := 0; cursor < len(source); {
		start := strings.Index(source[cursor:], prefix)
		if start < 0 {
			break
		}
		start += cursor + len(prefix)
		end := start
		for end < len(source) && source[end] >= '0' && source[end] <= '9' {
			end++
		}
		if strings.HasPrefix(source[end:], "\ue001") {
			used[source[start:end]] = true
		}
		cursor = end
	}
	for index := 0; ; index++ {
		candidate := strconv.Itoa(index)
		if !used[candidate] {
			return prefix + candidate + "\ue001"
		}
	}
}
