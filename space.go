// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lute

import (
	"unicode"
	"unicode/utf8"
)

// space 会把文本节点 textNode 中的中西文之间加上空格。
func (r *BaseRenderer) space(textNode *Node) {
	text := bytesToStr(textNode.Tokens)
	text = space0(text)
	textNode.Tokens = strToBytes(text)
}

func space0(text string) (ret string) {
	runes := []rune(text)
	length := len(runes)
	var r rune
	for i := 0; i < length; {
		r = runes[i]
		if i < length-3 && 'i' == runes[i+1] && 'n' == runes[i+2] && 'g' == runes[i+3] && unicode.Is(unicode.Han, runes[i]) {
			// ing 前不需要空格，如 打码ing https://github.com/88250/lute/issues/9
			ret += string(r) + "ing"
			i += 4
			continue
		}
		ret = addSpaceAtBoundary(ret, r)
		i++
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

	if ('%' == currentChar && nextIsLetter && !nextIsASCII) || (!currentIsASCII && currentIsLetter && '%' == nextChar) {
		return true
	}

	currentIsDigit := '0' <= currentChar && '9' >= currentChar
	nextIsDigit := '0' <= nextChar && '9' >= nextChar

	nextIsSymbol := unicode.IsSymbol(nextChar) && '~' != nextChar
	currentIsPunct := unicode.IsPunct(currentChar) || '~' == currentChar
	nextIsPunct := unicode.IsPunct(nextChar) || '~' == nextChar

	if !currentIsPunct && !currentIsASCII && !unicode.Is(unicode.Han, currentChar) {
		// Emoji 后不应该有空格
		return false
	}

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
		currentIsSymbol := unicode.IsSymbol(currentChar) && '~' != currentChar
		if currentIsSymbol && (nextIsDigit || nextIsPunct || !nextIsASCII) {
			return false
		}
		if currentIsLetter && nextIsPunct {
			return false
		}
		return true
	}
}
