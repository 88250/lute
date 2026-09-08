package parse

import (
	"bytes"
	"strings"

	"github.com/88250/lute/ast"
)

// parseTextMarkDelimiters 在相邻文本和样式节点之间配对分隔符，并保留每段文本的样式。
func parseTextMarkDelimiters(tree *Tree) {
	var groups [][]*ast.Node
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering || n.FirstChild == nil {
			return ast.WalkContinue
		}
		var group []*ast.Node
		flush := func() {
			if len(group) > 1 {
				groups = append(groups, group)
			}
			group = nil
		}
		for c := n.FirstChild; c != nil; c = c.Next {
			if c.Type == ast.NodeText || isDelimiterTextMark(c) {
				group = append(group, c)
			} else {
				flush()
			}
		}
		flush()
		return ast.WalkContinue
	})
	for _, group := range groups {
		parseTextMarkDelimiterGroup(tree, group)
	}
}

func isDelimiterTextMark(n *ast.Node) bool {
	if n.Type != ast.NodeTextMark || n.TextMarkFlashcardOcclusionID != "" {
		return false
	}
	for _, typ := range strings.Fields(n.TextMarkType) {
		switch typ {
		case "strong", "em", "s", "mark", "sup", "sub", "u", "text":
		default:
			return false
		}
	}
	return n.TextMarkType != ""
}

func parseTextMarkDelimiterGroup(tree *Tree, group []*ast.Node) {
	var tokens []byte
	var owners []*ast.Node
	hasStyle := false
	for _, n := range group {
		content := n.Tokens
		if n.Type == ast.NodeTextMark {
			content = []byte(n.TextMarkTextContent)
			hasStyle = true
		}
		// 转义、实体和独立语法交由完整行级解析器处理。
		if bytes.ContainsAny(content, "\\`$&<>\n[]()") {
			return
		}
		tokens = append(tokens, content...)
		for range content {
			owners = append(owners, n)
		}
	}
	if !hasStyle || !bytes.ContainsAny(tokens, "*_~=#^") {
		return
	}
	block := &ast.Node{Type: ast.NodeParagraph}
	ctx := &InlineContext{tokens: tokens, tokensLen: len(tokens)}
	markerOwners := map[*ast.Node]*ast.Node{}
	markerLengths := map[*ast.Node]int{}
	appendText := func(start, end int) {
		for start < end {
			next := start + 1
			for next < end && owners[next] == owners[start] {
				next++
			}
			owner := owners[start]
			if owner.Type == ast.NodeTextMark {
				block.AppendChild(&ast.Node{Type: ast.NodeTextMark, TextMarkType: owner.TextMarkType,
					TextMarkTextContent: string(tokens[start:next]), KramdownIAL: owner.KramdownIAL})
			} else {
				block.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: tokens[start:next]})
			}
			start = next
		}
	}
	for ctx.pos < ctx.tokensLen {
		start := ctx.pos
		if strings.ContainsRune("*_~=#^", rune(tokens[start])) {
			tree.handleDelim(block, ctx)
			// 同一分隔符串跨越样式边界时保留原节点，避免误分配剩余标记符的样式。
			for i := start + 1; i < ctx.pos; i++ {
				if owners[i] != owners[start] {
					return
				}
			}
			markerOwners[block.LastChild] = owners[start]
			markerLengths[block.LastChild] = ctx.pos - start
		} else {
			ctx.pos++
			for ctx.pos < ctx.tokensLen && !strings.ContainsRune("*_~=#^", rune(tokens[ctx.pos])) {
				ctx.pos++
			}
			appendText(start, ctx.pos)
		}
	}
	tree.processEmphasis(nil, ctx)
	matched := false
	for n, length := range markerLengths {
		if n.Parent == nil || len(n.Tokens) != length {
			matched = true
		}
		if n.Parent != nil && markerOwners[n].Type == ast.NodeTextMark {
			owner := markerOwners[n]
			n.Type, n.TextMarkType = ast.NodeTextMark, owner.TextMarkType
			n.TextMarkTextContent, n.KramdownIAL = string(n.Tokens), owner.KramdownIAL
			n.Tokens = nil
		}
	}
	if !matched {
		return
	}
	for block.FirstChild != nil {
		group[0].InsertBefore(block.FirstChild)
	}
	for _, n := range group {
		n.Unlink()
	}
}
