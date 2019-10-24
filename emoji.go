// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

// emoji 将 node 下文本节点中的 Emoji 别名替换为原生 Unicode 字符。
func (t *Tree) emoji(node *Node) {
	for child := node.firstChild; nil != child; {
		next := child.next
		if NodeText == child.typ {
			t.emoji0(child)
		} else {
			t.emoji(child) // 递归处理子节点
		}
		child = next
	}
}

var emojiSitePlaceholder = strToItems("${emojiSite}")
var emojiDot = strToItems(".")

func (t *Tree) emoji0(node *Node) {
	tokens := node.tokens
	node.tokens = items{} // 先清空，后面逐个添加或者添加 tokens 或者 Emoji 兄弟节点
	length := len(tokens)
	var token byte
	var maybeEmoji []byte
	var pos int
	for i := 0; i < length; {
		token = tokens[i].term()
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
			token = tokens[pos].term()
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

		maybeEmoji = itemsToBytes(tokens[i+1 : pos])
		if 1 > len(maybeEmoji) {
			node.tokens = append(node.tokens, tokens[pos])
			i++
			continue
		}

		if emoji, ok := t.context.option.Emojis[bytesToStr(maybeEmoji)]; ok {
			emojiNode := &Node{typ: NodeEmojiUnicode}
			emojiTokens := strToItems(emoji)
			if contains(emojiTokens, emojiSitePlaceholder) { // 有的 Emoji 是图片链接，需要单独处理
				alias := bytesToStr(maybeEmoji)
				repl := "<img alt=\"" + alias + "\" class=\"emoji\" src=\"" + t.context.option.EmojiSite + "/" + alias
				suffix := ".png"
				if "huaji" == alias {
					suffix = ".gif"
				}
				repl += suffix + "\" title=\"" + alias + "\" />"

				emojiNode.typ = NodeEmojiImg
				emojiNode.tokens = strToItems(repl)
			} else if contains(emojiTokens, emojiDot) { // 自定义 Emoji 路径用 . 判断，包含 . 的认为是图片路径
				alias := bytesToStr(maybeEmoji)
				repl := "<img alt=\"" + alias + "\" class=\"emoji\" src=\"" + emoji + "\" title=\"" + alias + "\" />"
				emojiNode.typ = NodeEmojiImg
				emojiNode.tokens = strToItems(repl)
			} else {
				emojiNode.tokens = emojiTokens
			}

			emojiNode.AppendChild(&Node{typ: NodeEmojiAlias, tokens: tokens[i : pos+1]})
			node.InsertAfter(emojiNode)

			if pos+1 < length {
				// 在 Emoji 节点后插入一个内容为空的文本节点，留作下次迭代
				text := &Node{typ: NodeText, tokens: items{}}
				emojiNode.InsertAfter(text)
				node = text
			}
		} else {
			node.tokens = append(node.tokens, tokens[i:pos+1]...)
		}

		pos++
		i = pos
	}

	if 1 > len(node.tokens) {
		node.Unlink() // 丢弃空的文本节点
	}
}
