package lute

import (
	"bytes"
)

// emoji 将 node 下文本节点中的 Emoji 别名替换为原生 Unicode 字符。
func (t *Tree) emoji(node *Node) {
	if nil == node {
		return
	}

	for child := node.firstChild; nil != child; {
		next := child.next
		if NodeText == child.typ && nil != child.parent &&
			NodeLink != child.parent.typ /* 不处理链接 label */ {
			emoji0(child)
		} else {
			t.emoji(child) // 递归处理子节点
		}
		child = next
	}
}

var emojiSitePlaceholder = toItems("${imgStaticPath}")

func emoji0(node *Node) {
	tokens := node.tokens
	node.tokens = items{}
	length := len(tokens)
	var token byte
	var maybeEmoji items
	var pos int
	for i := 0; i < length; {
		token = tokens[i]
		if i == length-1 {
			node.tokens = append(node.tokens, tokens[pos:]...)
			break
		}

		if itemColon != token {
			i++
			continue
		}

		node.tokens = append(node.tokens, tokens[pos:i]...)

		matchCloseColon := false
		for pos = i + 1; pos < length; pos++ {
			token = tokens[pos]
			if isWhitespace(token) {
				break
			}
			if itemColon == token {
				matchCloseColon = true
				break
			}
		}
		if !matchCloseColon {
			node.tokens = append(node.tokens, tokens[i:pos]...)
			i++
			continue
		}

		maybeEmoji = tokens[i+1 : pos]
		if 1 > len(maybeEmoji) {
			node.tokens = append(node.tokens, tokens[pos])
			i++
			continue
		}

		if emoji, ok := emojis[fromItems(maybeEmoji)]; ok {
			if bytes.Contains(items(emoji), emojiSitePlaceholder) { // 有的 Emoji 是图片链接，需要单独处理
				alias := fromItems(maybeEmoji)
				repl := "<img alt=\"" + alias + "\" class=\"emoji\" src=\"" + EmojiSite + alias
				suffix := ".png"
				if "huaji" == alias {
					suffix = ".gif"
				}
				repl += suffix + "\" title=\"" + alias + "\" />"
				img := &Node{typ: NodeInlineHTML, tokens: toItems(repl)}
				node.InsertAfter(node, img)
				text := &Node{typ: NodeText, tokens: items{}} // 生成一个新文本节点
				img.InsertAfter(node, text)
				node = text
			} else {
				node.tokens = append(node.tokens, emoji...)
			}
		} else {
			node.tokens = append(node.tokens, tokens[i:pos+1]...)
		}

		pos++
		i = pos
	}
}
