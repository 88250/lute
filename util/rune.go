// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package util

import (
	"strings"
	"unicode/utf8"

	"github.com/88250/lute/editor"
)

func ContainsStr(str string, strs []string) bool {
	if 0 == len(strs) {
		return false
	}
	for _, s := range strs {
		if strings.Contains(str, s) {
			return true
		}
	}
	return false
}

func IsEmptyStr(str string) bool {
	str = strings.ReplaceAll(str, editor.Zwsp, "")
	str = strings.ReplaceAll(str, editor.Zwj, "")
	return 0 == len(strings.TrimSpace(str))
}

func WordCount(str string) (runeCount, wordCount int) {
	words := strings.Fields(str)
	for _, word := range words {
		r, w := wordCount0(word)
		runeCount += r
		wordCount += w
	}
	return
}

func wordCount0(str string) (runeCount, wordCount int) {
	runes := []rune(str)
	length := len(runes)
	if 1 > length {
		return
	}

	runeCount, wordCount = 1, 1
	isAscii := runes[0] < utf8.RuneSelf
	for i := 1; i < length; i++ {
		r := runes[i]
		runeCount++
		if r >= utf8.RuneSelf {
			wordCount++
			isAscii = false
			continue
		}

		if r < utf8.RuneSelf == isAscii {
			continue
		}
		wordCount++
		isAscii = !isAscii
	}
	return
}

func IsDigit(str string) bool {
	if 0 == len(str) {
		return false
	}
	for _, r := range str {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
