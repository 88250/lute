package render

import "github.com/88250/lute/ast"

// 普通文本中的连续围栏字符需逐个转义，避免重新解析时开启代码块或行级代码。
func escapeTextCodeFences(node *ast.Node, tokens []byte) []byte {
	var escaped []byte
	copied := 0
	for start := 0; start < len(tokens); {
		marker := tokens[start]
		if marker != '`' && marker != '~' {
			start++
			continue
		}
		end := start + 1
		for end < len(tokens) && tokens[end] == marker {
			end++
		}
		count := end - start
		if start == 0 {
			for previous := node.Previous; count < 3 && previous != nil && previous.Type == ast.NodeText; previous = previous.Previous {
				i := len(previous.Tokens) - 1
				for i >= 0 && previous.Tokens[i] == marker && count < 3 {
					count++
					i--
				}
				if i >= 0 {
					break
				}
			}
		}
		if end == len(tokens) {
			for next := node.Next; count < 3 && next != nil && next.Type == ast.NodeText; next = next.Next {
				i := 0
				for i < len(next.Tokens) && next.Tokens[i] == marker && count < 3 {
					count++
					i++
				}
				if i < len(next.Tokens) {
					break
				}
			}
		}
		if count >= 3 {
			escaped = append(escaped, tokens[copied:start]...)
			for i := start; i < end; i++ {
				escaped = append(escaped, '\\', marker)
			}
			copied = end
		}
		start = end
	}
	if escaped == nil {
		return tokens
	}
	return append(escaped, tokens[copied:]...)
}
