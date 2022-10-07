// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"bytes"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

// emoji 将 node 下文本节点和链接文本节点中的 Emoji 别名替换为原生 Unicode 字符。
func (t *Tree) emoji(node *ast.Node) {
	for child := node.FirstChild; nil != child; {
		next := child.Next
		if ast.NodeText == child.Type || ast.NodeLinkText == child.Type {
			t.emoji0(child)
		} else {
			t.emoji(child) // 递归处理子节点
		}
		child = next
	}
}

var EmojiSitePlaceholder = util.StrToBytes("${emojiSite}")
var emojiDot = util.StrToBytes(".")

func (t *Tree) emoji0(node *ast.Node) {
	first := node
	tokens := node.Tokens
	node.Tokens = []byte{} // 先清空，后面逐个添加或者添加 Tokens 或者 Emoji 兄弟节点
	length := len(tokens)
	var token byte
	var maybeEmoji []byte
	var pos int
	for i := 0; i < length; {
		token = tokens[i]
		if i == length-1 {
			node.Tokens = append(node.Tokens, tokens[pos:]...)
			break
		}

		if lex.ItemColon != token {
			i++
			continue
		}

		node.Tokens = append(node.Tokens, tokens[pos:i]...)

		matchCloseColon := false
		for pos = i + 1; pos < length; pos++ {
			token = tokens[pos]
			if lex.IsWhitespace(token) {
				break
			}
			if lex.ItemColon == token {
				matchCloseColon = true
				break
			}
		}
		if !matchCloseColon {
			node.Tokens = append(node.Tokens, tokens[i:pos]...)
			i++
			continue
		}

		maybeEmoji = tokens[i+1 : pos]
		if 1 > len(maybeEmoji) {
			node.Tokens = append(node.Tokens, tokens[pos])
			i++
			continue
		}

		EmojiLock.Lock()
		emoji, ok := t.Context.ParseOption.AliasEmoji[util.BytesToStr(maybeEmoji)]
		EmojiLock.Unlock()
		if ok {
			emojiNode := &ast.Node{Type: ast.NodeEmoji}
			emojiUnicodeOrImg := &ast.Node{Type: ast.NodeEmojiUnicode}
			emojiNode.AppendChild(emojiUnicodeOrImg)
			emojiTokens := util.StrToBytes(emoji)
			if bytes.Contains(emojiTokens, EmojiSitePlaceholder) { // 有的 Emoji 是图片链接，需要单独处理
				alias := util.BytesToStr(maybeEmoji)
				suffix := ".png"
				if "huaji" == alias {
					suffix = ".gif"
				} else if "siyuan" == alias {
					suffix = ".svg"
				}
				src := t.Context.ParseOption.EmojiSite + "/" + alias + suffix
				emojiUnicodeOrImg.Type = ast.NodeEmojiImg
				emojiUnicodeOrImg.Tokens = t.EmojiImgTokens(alias, src)
			} else if bytes.Contains(emojiTokens, emojiDot) { // 自定义 Emoji 路径用 . 判断，包含 . 的认为是图片路径
				alias := util.BytesToStr(maybeEmoji)
				emojiUnicodeOrImg.Type = ast.NodeEmojiImg
				emojiUnicodeOrImg.Tokens = t.EmojiImgTokens(alias, emoji)
			} else {
				emojiUnicodeOrImg.Tokens = emojiTokens
			}

			emojiUnicodeOrImg.AppendChild(&ast.Node{Type: ast.NodeEmojiAlias, Tokens: tokens[i : pos+1]})
			node.InsertAfter(emojiNode)

			if pos+1 < length {
				// 在 Emoji 节点后插入一个内容为空的文本节点，留作下次迭代
				text := &ast.Node{Type: ast.NodeText, Tokens: []byte{}}
				emojiNode.InsertAfter(text)
				node = text
			}
		} else {
			node.Tokens = append(node.Tokens, tokens[i:pos+1]...)
		}

		pos++
		i = pos
	}

	// 丢弃空的文本节点
	if 1 > len(first.Tokens) {
		first.Unlink()
	}
	if 1 > len(node.Tokens) {
		node.Unlink()
	}
}

func (t *Tree) EmojiImgTokens(alias, src string) []byte {
	return util.StrToBytes("<img alt=\"" + alias + "\" class=\"emoji\" src=\"" + src + "\" title=\"" + alias + "\" />")
}
