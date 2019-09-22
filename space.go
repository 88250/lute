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

import (
	"unicode"
	"unicode/utf8"
)

// space 会把 node 下文本节点中的中西文之间加上空格。
func (t *Tree) space(node *Node) {
	for child := node.firstChild; nil != child; {
		next := child.next
		if NodeText == child.typ && nil != child.parent &&
			NodeLink != child.parent.typ /* 不处理链接 label */ {
			text := fromItems(child.tokens)
			text = space0(text)
			child.tokens = toItems(text)
		} else {
			t.space(child) // 递归处理子节点
		}
		child = next
	}
}

func space0(text string) (ret string) {
	for _, r := range text {
		ret = addSpaceAtBoundary(ret, r)
	}
	return
}

func addSpaceAtBoundary(prefix string, nextChar rune) string {
	if 0 == len(prefix) {
		return string(nextChar)
	}
	if unicode.IsSpace(nextChar) {
		return prefix + string(nextChar)
	}

	currentChar, _ := utf8.DecodeLastRuneInString(prefix)
	if isAllowSpace(currentChar, nextChar) {
		return prefix + " " + string(nextChar)
	}
	return prefix + string(nextChar)
}

func isAllowSpace(currentChar, nextChar rune) bool {
	if unicode.IsSpace(currentChar) {
		return false
	}

	if utf8.RuneSelf <= currentChar { // 当前字符不是 ASCII 字符
		if '℃' == currentChar {
			// 摄氏度符号后必须跟空格
			return true
		}
		if '%' == nextChar {
			return true
		}
		if unicode.IsPunct(nextChar) {
			// 后面的字符如果是标点符号的话不需要后跟空格
			return false
		}
		return utf8.RuneSelf > nextChar && !unicode.IsPunct(currentChar)
	} else { // 当前字符是 ASCII 字符
		if unicode.IsDigit(currentChar) && '℃' == nextChar {
			// 数字后更摄氏度符号不需要空格
			return false
		}
		if '%' == currentChar {
			return true
		}
		if unicode.IsPunct(currentChar) {
			// 当前字符如果是 ASCII 标点符号的话不需要后跟空格
			return false
		}
		return utf8.RuneSelf <= nextChar && !unicode.IsPunct(nextChar)
	}
}
