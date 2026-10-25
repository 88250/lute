package parse

import (
	"encoding/json"
	"strings"

	"github.com/88250/lute/ast"
)

// ParseDocumentTitle 从完整源构造无块身份的临时段落，不修改源或补造正文内容。
func ParseDocumentTitle(title *ast.DocumentTitle) (*ast.Node, error) {
	if err := title.Validate(); nil != err {
		return nil, err
	}
	var nodes []json.RawMessage
	if json.Unmarshal([]byte(title.Content), &nodes) != nil || nil == nodes {
		return nil, ast.ErrInvalidDocumentTitle
	}
	ret := &ast.Node{Type: ast.NodeParagraph}
	count := 0
	var appendNode func(*ast.Node, json.RawMessage, int) error
	appendNode = func(parent *ast.Node, raw json.RawMessage, depth int) error {
		count++
		if count > ast.DocumentTitleMaxNodes || depth > ast.DocumentTitleMaxDepth {
			return ast.ErrInvalidDocumentTitle
		}
		var fields map[string]json.RawMessage
		if json.Unmarshal(raw, &fields) != nil || nil == fields {
			return ast.ErrInvalidDocumentTitle
		}
		for key := range fields {
			switch key {
			case "Type", "Data", "Properties", "Children", "HtmlEntityTokens", "LinkType", "LinkRefLabel", "CodeMarkerLen",
				"TextMarkType", "TextMarkAHref", "TextMarkATitle", "TextMarkInlineMathContent", "TextMarkInlineMemoContent",
				"TextMarkBlockRefID", "TextMarkBlockRefSubtype", "TextMarkFileAnnotationRefID", "TextMarkFlashcardOcclusionID",
				"TextMarkTextContent":
			default:
				return ast.ErrInvalidDocumentTitle
			}
		}
		node := &ast.Node{}
		if json.Unmarshal(raw, node) != nil {
			return ast.ErrInvalidDocumentTitle
		}
		node.Type = ast.Str2NodeType(node.TypeStr)
		_, hasID := node.Properties["id"]
		if !documentTitleChildAllowed(parent.Type, node.Type) || hasID {
			return ast.ErrInvalidDocumentTitle
		}
		if node.Type == ast.NodeTextMark {
			for _, mark := range strings.Fields(node.TextMarkType) {
				switch mark {
				case "text", "strong", "em", "u", "s", "mark", "sup", "sub", "kbd", "code", "tag",
					"a", "block-ref", "inline-math", "inline-memo", "file-annotation-ref":
				default:
					return ast.ErrInvalidDocumentTitle
				}
			}
			if node.IsTextMarkType("block-ref") && (!ast.IsNodeIDPattern(node.TextMarkBlockRefID) ||
				(node.TextMarkBlockRefSubtype != "s" && node.TextMarkBlockRefSubtype != "d")) {
				return ast.ErrInvalidDocumentTitle
			}
		}
		node.Tokens = []byte(node.Data)
		node.Data, node.TypeStr = "", ""
		node.KramdownIAL = Map2IAL(node.Properties)
		node.Properties = nil
		var children []json.RawMessage
		if rawChildren := fields["Children"]; nil != rawChildren && json.Unmarshal(rawChildren, &children) != nil {
			return ast.ErrInvalidDocumentTitle
		}
		node.Children = nil
		parent.AppendChild(node)
		for _, child := range children {
			if err := appendNode(node, child, depth+1); nil != err {
				return err
			}
		}
		if node.Type == ast.NodeImage && !documentTitleImageValid(node) {
			return ast.ErrInvalidDocumentTitle
		}
		if node.Type == ast.NodeBackslash && (node.FirstChild == nil || node.FirstChild != node.LastChild) {
			return ast.ErrInvalidDocumentTitle
		}
		return nil
	}
	for _, node := range nodes {
		if err := appendNode(ret, node, 1); nil != err {
			return nil, err
		}
	}
	return ret, nil
}

func documentTitleChildAllowed(parent, child ast.NodeType) bool {
	switch parent {
	case ast.NodeParagraph:
		switch child {
		case ast.NodeText, ast.NodeTextMark, ast.NodeImage, ast.NodeKramdownSpanIAL,
			ast.NodeSoftBreak, ast.NodeHardBreak, ast.NodeBr, ast.NodeBackslash,
			ast.NodeHTMLEntity, ast.NodeEmoji:
			return true
		}
	case ast.NodeImage:
		switch child {
		case ast.NodeBang, ast.NodeOpenBracket, ast.NodeCloseBracket, ast.NodeOpenParen, ast.NodeCloseParen,
			ast.NodeLinkText, ast.NodeLinkDest, ast.NodeLinkSpace, ast.NodeLinkTitle:
			return true
		}
	case ast.NodeBackslash:
		return child == ast.NodeBackslashContent || child == ast.NodeText
	case ast.NodeEmoji:
		return child == ast.NodeEmojiUnicode || child == ast.NodeEmojiImg
	case ast.NodeEmojiUnicode, ast.NodeEmojiImg:
		return child == ast.NodeEmojiAlias
	}
	return false
}

func documentTitleImageValid(image *ast.Node) bool {
	child := image.FirstChild
	for _, typ := range []ast.NodeType{ast.NodeBang, ast.NodeOpenBracket, ast.NodeLinkText, ast.NodeCloseBracket,
		ast.NodeOpenParen, ast.NodeLinkDest} {
		// Native Markdown omits the alt-text child for ![](dest).
		if typ == ast.NodeLinkText && child != nil && child.Type == ast.NodeCloseBracket {
			continue
		}
		if nil == child || child.Type != typ {
			return false
		}
		child = child.Next
	}
	if nil != child && child.Type == ast.NodeLinkSpace {
		child = child.Next
		if nil != child && child.Type == ast.NodeLinkTitle {
			child = child.Next
		}
	}
	return nil != child && child.Type == ast.NodeCloseParen && nil == child.Next
}

// CloneDocumentTitleInline 复制行级子树和属性，清理只用于块身份的字段，保留原有换行和转义。
func CloneDocumentTitleInline(node *ast.Node) *ast.Node {
	clone := *node
	clone.ID, clone.Box, clone.Path, clone.Spec = "", "", "", ""
	clone.Parent, clone.Previous, clone.Next, clone.FirstChild, clone.LastChild = nil, nil, nil, nil, nil
	clone.Children, clone.Properties = nil, nil
	clone.Tokens = append([]byte(nil), node.Tokens...)
	clone.HtmlEntityTokens = append([]byte(nil), node.HtmlEntityTokens...)
	clone.KramdownIAL = nil
	for _, attribute := range node.KramdownIAL {
		if attribute[0] != "id" {
			clone.KramdownIAL = append(clone.KramdownIAL, append([]string(nil), attribute...))
		}
	}
	for child := node.FirstChild; nil != child; child = child.Next {
		clone.AppendChild(CloneDocumentTitleInline(child))
	}
	return &clone
}
