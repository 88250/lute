package render

import "github.com/88250/lute/ast"

func pairedVditorEmphasisHTML(node *ast.Node) bool {
	tag := string(node.Tokens)
	var counterpart string
	var forward bool
	switch tag {
	case "<strong>":
		counterpart, forward = "</strong>", true
	case "<em>":
		counterpart, forward = "</em>", true
	case "</strong>":
		counterpart = "<strong>"
	case "</em>":
		counterpart = "<em>"
	default:
		return false
	}
	if forward {
		return pairedVditorHTMLTag(node, tag, counterpart)
	}
	return pairedVditorHTMLTag(node, counterpart, tag)
}

func pairedVditorHTMLTag(node *ast.Node, open, close string) bool {
	tag := string(node.Tokens)
	forward := tag == open
	counterpart := open
	if forward {
		counterpart = close
	} else if tag != close {
		return false
	}
	depth := 1
	for sibling := node; nil != sibling; {
		if forward {
			sibling = sibling.Next
		} else {
			sibling = sibling.Previous
		}
		if nil == sibling {
			break
		}
		if sibling.Type != ast.NodeInlineHTML {
			continue
		}
		if string(sibling.Tokens) == tag {
			depth++
		} else if string(sibling.Tokens) == counterpart {
			depth--
			if depth == 0 {
				return true
			}
		}
	}
	return false
}
