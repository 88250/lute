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
	"net/url"
	"strings"
	"unicode"
)

var htmlEscaper = strings.NewReplacer(
	`&`, "&amp;",
	`<`, "&lt;",
	`>`, "&gt;",
	`"`, "&quot;",
)

func escapeHTML(html string) string {
	return htmlEscaper.Replace(html)
}

func unescapeString(str string) string {
	var ret string
	for i := 0; i < len(str); i++ {
		if isBlackslashEscape(str, i) {
			ret = ret[:len(ret)-1]
		}
		ret += string(str[i])
	}

	return ret
}

func isBlackslashEscape(str string, pos int) bool {
	if !unicode.IsPunct(rune(str[pos])) {
		return false
	}

	backslashes := 0
	for i := pos - 1; 0 <= i; i-- {
		if '\\' != str[i] {
			break
		}

		backslashes++
	}

	return 0 != backslashes%2
}

func encodeDestination(destination string) (ret string) {
	destination = unescapeString(destination)
	u, e := url.ParseRequestURI(destination)
	if nil != e {
		return destination
	}

	ret = u.String()
	ret = compatibleJSEncodeURIComponent(ret)

	return ret
}

func compatibleJSEncodeURIComponent(str string) string {
	str = strings.ReplaceAll(str, "+", "%20")
	str = strings.ReplaceAll(str, "%21", "!")
	str = strings.ReplaceAll(str, "%27", "'")
	str = strings.ReplaceAll(str, "%28", "(")
	str = strings.ReplaceAll(str, "%29", ")")
	str = strings.ReplaceAll(str, "%2A", "*")

	return str
}
