// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// space 会把 text 中的中西文之间加上空格。
func space(text string) (ret string) {
	// 鸣谢 https://github.com/studygolang/autocorrect

	for _, r := range text {
		ret = addSpaceAtBoundary(ret, r)
	}
	return
}

func addSpaceAtBoundary(prefix string, nextChar rune) string {
	if len(prefix) == 0 {
		return string(nextChar)
	}

	r, size := utf8.DecodeLastRuneInString(prefix)
	if isLatin(size) != isLatin(utf8.RuneLen(nextChar)) &&
		isAllowSpace(nextChar) && isAllowSpace(r) {
		return prefix + " " + string(nextChar)
	}

	return prefix + string(nextChar)
}

func isLatin(size int) bool {
	return size == 1
}

func isAllowSpace(r rune) bool {
	return !unicode.IsSpace(r) && !strings.ContainsRune("，。；「」：《》『』、[]（）*_", r)
}
