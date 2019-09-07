package lute

import (
	"bytes"
)

// emoji 将 node 下文本节点中的 Emoji 别名替换为原生 Unicode 字符。
func (t *Tree) emoji(node *Node) {
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

// emojiSite 为图片 Emoji URL 的路径前缀。
var emojiSite = "https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji/"

func emoji0(node *Node) {
	tokens := node.tokens
	node.tokens = items{} // 先清空，后面逐个添加或者添加 tokens 或者 Emoji 兄弟节点
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
			emojiNode := &Node{typ: NodeEmojiUnicode}
			if bytes.Contains(items(emoji), emojiSitePlaceholder) { // 有的 Emoji 是图片链接，需要单独处理
				alias := fromItems(maybeEmoji)
				repl := "<img alt=\"" + alias + "\" class=\"emoji\" src=\"" + emojiSite + alias
				suffix := ".png"
				if "huaji" == alias {
					suffix = ".gif"
				}
				repl += suffix + "\" title=\"" + alias + "\" />"

				emojiNode.typ = NodeEmojiImg
				emojiNode.tokens = toItems(repl)
				emojiNode.emojiImgAlias = tokens[i : pos+1]
			} else {
				emojiNode.tokens = toItems(emoji)
			}

			node.InsertAfter(node, emojiNode)
			// 在 EmojiImg 节点后插入一个内容为空的文本节点，留作下次迭代
			text := &Node{typ: NodeText, tokens: items{}}
			emojiNode.InsertAfter(emojiNode, text)
			node = text
		} else {
			node.tokens = append(node.tokens, tokens[i:pos+1]...)
		}

		pos++
		i = pos
	}
}
