package parse

import (
	"bytes"

	"github.com/88250/lute/ast"
)

// MatchLatexMath matches one complete, nonempty alternate-delimited formula at the start of tokens.
// It returns a view of the original TeX, without decoding escapes or parsing Markdown inside it.
func MatchLatexMath(tokens []byte, block bool) (content []byte, consumed int) {
	content, consumed, _ = matchLatexMath(tokens, block)
	return
}

func matchLatexMath(tokens []byte, block bool) (content []byte, consumed int, invalid bool) {
	open, close := byte('('), byte(')')
	if block {
		open, close = '[', ']'
	}
	if len(tokens) < 5 || tokens[0] != '\\' || tokens[1] != open {
		return
	}
	depth := 0
	for i := 2; i < len(tokens); i++ {
		switch tokens[i] {
		case '\n', '\r':
			if !block {
				return nil, 0, true
			}
		case '{':
			depth++
		case '}':
			if depth > 0 {
				depth--
			}
		case '\\':
			if i+1 == len(tokens) {
				return nil, 0, false
			}
			next := tokens[i+1]
			if !block && (next == '\n' || next == '\r') {
				return nil, 0, true
			}
			if depth == 0 && next == close {
				if len(bytes.TrimSpace(tokens[2:i])) == 0 {
					return nil, 0, true
				}
				return tokens[2:i], i + 2, false
			}
			if depth == 0 && next == open {
				return nil, 0, true
			}
			i++ // TeX escapes, including \\, \{ and \}, are indivisible.
		}
	}
	return nil, 0, false
}

// Look ahead without advancing the real lexer. Unclosed delimiters remain ordinary Markdown.
// Container prefixes are checked so a delimiter outside a quote or list cannot close its formula.
func (t *Tree) latexMathBlockClosed(container *ast.Node) bool {
	first := t.Context.currentLine[t.Context.nextNonspace:]
	if !bytes.HasPrefix(first, []byte(`\[`)) {
		return false
	}
	var containers []*ast.Node
	for n := container; n != nil; n = n.Parent {
		if n.Type == ast.NodeBlockquote || n.Type == ast.NodeListItem {
			containers = append(containers, n)
		}
	}
	lookahead := *t.lexer
	tokens := append([]byte(nil), first...)
	checked := 2
	for {
		if bytes.Contains(tokens[checked:], []byte(`\]`)) || bytes.Contains(tokens[checked:], []byte(`\[`)) {
			_, consumed, invalid := matchLatexMath(tokens, true)
			if consumed > 0 {
				return len(bytes.TrimSpace(tokens[consumed:])) == 0
			}
			// A new opener prevents an unmatched opener from consuming later formulas.
			if invalid {
				return false
			}
		}
		checked = len(tokens)
		line := lookahead.NextLine()
		if line == nil {
			return false
		}
		for i := len(containers) - 1; i >= 0; i-- {
			c := containers[i]
			if c.Type == ast.NodeBlockquote {
				line = bytes.TrimLeft(line, " \t")
				if len(line) == 0 || line[0] != '>' {
					return false
				}
				line = line[1:]
				if len(line) > 0 && line[0] == ' ' {
					line = line[1:]
				}
			} else if len(bytes.TrimSpace(line)) > 0 {
				columns, offset := 0, 0
				for offset < len(line) && (line[offset] == ' ' || line[offset] == '\t') {
					if line[offset] == '\t' {
						columns += 4 - columns%4
					} else {
						columns++
					}
					offset++
					if columns >= c.ListData.MarkerOffset+c.ListData.Padding {
						break
					}
				}
				if columns < c.ListData.MarkerOffset+c.ListData.Padding {
					return false
				}
				line = line[offset:]
			}
		}
		tokens = append(tokens, line...)
	}
}
