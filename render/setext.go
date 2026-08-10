// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"bytes"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/lex"
)

// setextHeadingMarkerOffsetInListParagraph 返回列表段落续行中 Setext 标记符首字符的位置。
func setextHeadingMarkerOffsetInListParagraph(node *ast.Node) int {
	if ast.NodeText != node.Type || nil == node.Parent || ast.NodeParagraph != node.Parent.Type ||
		nil == node.Parent.Parent || ast.NodeListItem != node.Parent.Parent.Type {
		return -1
	}

	lineStart := node
	foundBreak := false
	for previous := node.Previous; nil != previous; previous = previous.Previous {
		if ast.NodeSoftBreak == previous.Type || ast.NodeHardBreak == previous.Type {
			foundBreak = true
			break
		}
		lineStart = previous
	}
	if !foundBreak {
		return -1
	}

	leadingWhitespace := 0
	marker := byte(0)
	markerCount := 0
	trailingWhitespace := false
	markerNode := (*ast.Node)(nil)
	markerOffset := -1
	for current := lineStart; nil != current; current = current.Next {
		if ast.NodeSoftBreak == current.Type || ast.NodeHardBreak == current.Type {
			break
		}
		if ast.NodeVditorCaret == current.Type {
			continue
		}
		if ast.NodeText != current.Type {
			return -1
		}

		for i := 0; i < len(current.Tokens); {
			if bytes.HasPrefix(current.Tokens[i:], editor.CaretTokens) {
				i += len(editor.CaretTokens)
				continue
			}

			token := current.Tokens[i]
			if 0 == marker {
				if lex.ItemSpace == token || lex.ItemTab == token {
					leadingWhitespace++
					if 3 < leadingWhitespace {
						return -1
					}
					i++
					continue
				}
				if lex.ItemHyphen != token && lex.ItemEqual != token {
					return -1
				}
				marker = token
				markerCount = 1
				markerNode = current
				markerOffset = i
				i++
				continue
			}

			if marker == token && !trailingWhitespace {
				markerCount++
				i++
				continue
			}
			if lex.ItemSpace == token || lex.ItemTab == token {
				trailingWhitespace = true
				i++
				continue
			}
			return -1
		}
	}

	if 0 == markerCount || markerNode != node {
		return -1
	}
	return markerOffset
}

func escapeSetextHeadingMarkersInListParagraph(node *ast.Node, tokens []byte) []byte {
	if ast.NodeText != node.Type || nil == node.Parent || ast.NodeParagraph != node.Parent.Type ||
		nil == node.Parent.Parent || ast.NodeListItem != node.Parent.Parent.Type {
		return tokens
	}

	offsets := make([]int, 0, 2)
	if offset := setextHeadingMarkerOffsetInListParagraph(node); 0 <= offset {
		offsets = append(offsets, offset)
	}
	for lineStart := bytes.IndexByte(tokens, lex.ItemNewline) + 1; 0 < lineStart && lineStart <= len(tokens); {
		lineEnd := bytes.IndexByte(tokens[lineStart:], lex.ItemNewline)
		if 0 > lineEnd {
			lineEnd = len(tokens)
		} else {
			lineEnd += lineStart
		}
		if offset := setextHeadingMarkerOffset(tokens[lineStart:lineEnd]); 0 <= offset {
			offsets = append(offsets, lineStart+offset)
		}
		if len(tokens) <= lineEnd {
			break
		}
		lineStart = lineEnd + 1
	}
	if 0 == len(offsets) {
		return tokens
	}

	escaped := append([]byte(nil), tokens...)
	for i := len(offsets) - 1; 0 <= i; i-- {
		offset := offsets[i]
		escaped = append(escaped, 0)
		copy(escaped[offset+1:], escaped[offset:])
		escaped[offset] = lex.ItemBackslash
	}
	return escaped
}

func setextHeadingMarkerOffset(tokens []byte) int {
	leadingWhitespace := 0
	marker := byte(0)
	markerOffset := -1
	trailingWhitespace := false
	markerCount := 0
	for i := 0; i < len(tokens); {
		if bytes.HasPrefix(tokens[i:], editor.CaretTokens) {
			i += len(editor.CaretTokens)
			continue
		}

		token := tokens[i]
		if 0 == marker {
			if lex.ItemSpace == token || lex.ItemTab == token {
				leadingWhitespace++
				if 3 < leadingWhitespace {
					return -1
				}
				i++
				continue
			}
			if lex.ItemHyphen != token && lex.ItemEqual != token {
				return -1
			}
			marker = token
			markerOffset = i
			markerCount = 1
			i++
			continue
		}

		if marker == token && !trailingWhitespace {
			markerCount++
			i++
			continue
		}
		if lex.ItemSpace == token || lex.ItemTab == token {
			trailingWhitespace = true
			i++
			continue
		}
		return -1
	}
	if 0 == markerCount {
		return -1
	}
	return markerOffset
}
