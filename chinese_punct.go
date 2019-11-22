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

// chinesePunct 会把文本节点 textNode 中的中文间的英文逗号句号换成中文的逗号句号。
func (r *BaseRenderer) chinesePunct(textNode *Node) {
	text := bytesToStr(textNode.tokens)
	text = chinesePunct0(text)
	textNode.tokens = strToBytes(text)
}

func chinesePunct0(text string) (ret string) {
	for _, r := range text {
		ret = commaPeriod(ret, r)
	}
	return
}

func commaPeriod(prefix string, nextChar rune) string {
	if 0 == len(prefix) {
		return string(nextChar)
	}
	nextCharIsEnglishComma := ',' == nextChar
	nextCharIsEnglishPeriod := '.' == nextChar
	nextCharIsEnglishColon := ':' == nextChar

	if !nextCharIsEnglishComma && !nextCharIsEnglishPeriod && !nextCharIsEnglishColon {
		return prefix + string(nextChar)
	}

	currentChar, _ := utf8.DecodeLastRuneInString(prefix)
	if !unicode.Is(unicode.Han, currentChar) {
		return prefix + string(nextChar)
	}

	if nextCharIsEnglishComma {
		return prefix + "，"
	} else if nextCharIsEnglishPeriod {
		return prefix + "。"
	} else {
		return prefix + "："
	}
}
