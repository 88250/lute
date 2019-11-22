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

// space 会把文本节点 textNode 中的中西文之间加上空格。
func (r *BaseRenderer) space(textNode *Node) {
	text := bytesToStr(textNode.tokens)
	text = space0(text)
	textNode.tokens = strToBytes(text)
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
	if unicode.IsSpace(nextChar) || !unicode.IsPrint(nextChar) {
		return prefix + string(nextChar)
	}

	currentChar, _ := utf8.DecodeLastRuneInString(prefix)
	if allowSpace(currentChar, nextChar) {
		return prefix + " " + string(nextChar)
	}
	return prefix + string(nextChar)
}

func allowSpace(currentChar, nextChar rune) bool {
	if unicode.IsSpace(currentChar) || !unicode.IsPrint(currentChar) {
		return false
	}

	currentIsASCII := utf8.RuneSelf > currentChar
	nextIsASCII := utf8.RuneSelf > nextChar
	currentIsLetter := unicode.IsLetter(currentChar)
	nextIsLetter := unicode.IsLetter(nextChar)
	if currentIsASCII == nextIsASCII && currentIsLetter && nextIsLetter {
		return false
	}

	if (currentIsLetter && ('￥' == nextChar || '℃' == nextChar)) || (('￥' == currentChar || '℃' == currentChar) && nextIsLetter) {
		return true
	}

	if ('%' == currentChar && nextIsLetter) || (currentIsLetter && '%' == nextChar) {
		return true
	}

	currentIsDigit := '0' <= currentChar && '9' >= currentChar
	nextIsDigit := '0' <= nextChar && '9' >= nextChar

	nextIsSymbol := unicode.IsSymbol(nextChar)
	currentIsPunct := unicode.IsPunct(currentChar)
	nextIsPunct := unicode.IsPunct(nextChar)
	if currentIsASCII {
		if currentIsDigit && nextIsSymbol {
			return false
		}
		if currentIsPunct && nextIsLetter {
			return false
		}
		if nextIsPunct || nextIsSymbol {
			return false
		}
		return !nextIsASCII
	} else {
		if currentIsPunct {
			return false
		}
		if nextIsSymbol {
			return true
		}
		currentIsSymbol := unicode.IsSymbol(currentChar)
		if currentIsSymbol && (nextIsDigit || nextIsPunct || !nextIsASCII) {
			return false
		}
		if currentIsLetter && nextIsPunct {
			return false
		}
		return true
	}
}
