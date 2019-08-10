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
	"bytes"
	"html"
	"strings"
	"unicode/utf8"
)

func escapeHTML(html items) (ret items) {
	var i int
	var token byte
	tmp := html
	for i < len(tmp) {
		token = tmp[i]
		if itemAmpersand == token {
			if 0 == len(ret) { // 通过延迟初始化减少内存分配，下同
				ret = make(items, len(html), len(html))
				copy(ret, html)
				tmp = ret
			}

			tmp = append(tmp, 0, 0, 0, 0)
			copy(tmp[i+4:], tmp[i:])
			tmp[i] = '&'
			tmp[i+1] = 'a'
			tmp[i+2] = 'm'
			tmp[i+3] = 'p'
			tmp[i+4] = ';'
			i += 5

			continue
		}

		if itemLess == token {
			if 0 == len(ret) {
				ret = make(items, len(html), len(html))
				copy(ret, html)
				tmp = ret
			}

			tmp = append(tmp, 0, 0, 0)
			copy(tmp[i+3:], tmp[i:])
			tmp[i] = '&'
			tmp[i+1] = 'l'
			tmp[i+2] = 't'
			tmp[i+3] = ';'
			i += 4
			continue
		}

		if itemGreater == token {
			if 0 == len(ret) {
				ret = make(items, len(html), len(html))
				copy(ret, html)
				tmp = ret
			}

			tmp = append(tmp, 0, 0, 0)
			copy(tmp[i+3:], tmp[i:])
			tmp[i] = '&'
			tmp[i+1] = 'g'
			tmp[i+2] = 't'
			tmp[i+3] = ';'
			i += 4
			continue
		}

		if itemDoublequote == token {
			if 0 == len(ret) {
				ret = make(items, len(html), len(html))
				copy(ret, html)
				tmp = ret
			}

			tmp = append(tmp, 0, 0, 0, 0, 0)
			copy(tmp[i+5:], tmp[i:])
			tmp[i] = '&'
			tmp[i+1] = 'q'
			tmp[i+2] = 'u'
			tmp[i+3] = 'o'
			tmp[i+4] = 't'
			tmp[i+5] = ';'
			i += 6
			continue
		}

		i++
	}

	if 0 == len(ret) {
		return html
	}

	return tmp
}

func unescapeString(str string) string {
	if "" == str {
		return str
	}

	str = html.UnescapeString(str) // FIXME: 此处应该用内部的实体转义方式
	tokens := toItems(str)
	length := len(tokens)
	retTokens := make(items, 0, length)
	for i := 0; i < length; i++ {
		if tokens.isBackslashEscapePunct(i) {
			retTokens = retTokens[:len(retTokens)-1]
		}
		retTokens = append(retTokens, tokens[i])
	}

	return fromItems(retTokens)
}

// encodeDestination percent-encodes rawurl, avoiding double encoding.
// It doesn't touch:
// - alphanumeric characters ([0-9a-zA-Z]);
// - percent-encoded characters (%[0-9a-fA-F]{2});
// - excluded characters ([;/?:@&=+$,-_.!~*'()#]).
// Invalid UTF-8 sequences are replaced with U+FFFD.
// 鸣谢 https://gitlab.com/golang-commonmark/mdurl
func encodeDestination(rawurl string) string {
	const hexdigit = "0123456789ABCDEF"
	var buf bytes.Buffer
	i := 0
	for i < len(rawurl) {
		r, rlen := utf8.DecodeRuneInString(rawurl[i:])
		if r >= 0x80 {
			for j, n := i, i+rlen; j < n; j++ {
				b := rawurl[j]
				buf.WriteByte('%')
				buf.WriteByte(hexdigit[(b>>4)&0xf])
				buf.WriteByte(hexdigit[b&0xf])
			}
		} else if r == '%' {
			if i+2 < len(rawurl) && isHexDigit(rawurl[i+1]) && isHexDigit(rawurl[i+2]) {
				buf.WriteByte('%')
				buf.WriteByte(tokenToUpper(rawurl[i+1]))
				buf.WriteByte(tokenToUpper(rawurl[i+2]))
				i += 2
			} else {
				buf.WriteString("%25")
			}
		} else if strings.IndexByte("!#$&'()*+,-./0123456789:;=?@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~", byte(r)) == -1 {
			buf.WriteByte('%')
			buf.WriteByte(hexdigit[(r>>4)&0xf])
			buf.WriteByte(hexdigit[r&0xf])
		} else {
			buf.WriteByte(byte(r))
		}
		i += rlen
	}
	return buf.String()
}

//
//func encodeDestination(destination string) (ret string) {
//	if "" == destination {
//		return destination
//	}
//
//	destination = decodeDestination(destination)
//
//	parts := strings.SplitN(destination, ":", 2)
//	var scheme string
//	remains := destination
//	if 1 < len(parts) {
//		scheme = parts[0]
//		remains = parts[1]
//	}
//
//	index := strings.Index(remains, "?")
//	var query string
//	path := remains
//	if 0 <= index {
//		query = remains[index+1:]
//		queries := strings.Split(query, "&")
//		query = ""
//		length := len(queries)
//		for i, q := range queries {
//			kv := strings.Split(q, "=")
//			if 1 < len(kv) {
//				query += url.QueryEscape(kv[0]) + "=" + url.QueryEscape(kv[1])
//			} else {
//				query += url.QueryEscape(kv[0])
//			}
//			if i < length-1 {
//				query += "&"
//			}
//		}
//
//		path = remains[:index]
//	}
//
//	parts = strings.Split(path, "/")
//	path = ""
//	length := len(parts)
//	for i, part := range parts {
//		unescaped := url.PathEscape(part)
//		path += unescaped
//		if i < length-1 {
//			path += "/"
//		}
//	}
//
//	if "" == scheme {
//		ret = path
//	} else {
//		ret = scheme + ":" + path
//	}
//	if "" != query {
//		ret += "?" + query
//	}
//
//	ret = strings.ReplaceAll(ret, "%2A", "*")
//	ret = strings.ReplaceAll(ret, "%29", ")")
//	ret = strings.ReplaceAll(ret, "%28", "(")
//	ret = strings.ReplaceAll(ret, "%23", "#")
//	ret = strings.ReplaceAll(ret, "%2C", ",")
//
//	return
//}

//func decodeDestination(destination string) (ret string) {
//	if "" == destination {
//		return destination
//	}
//
//	parts := strings.SplitN(destination, ":", 2)
//	var scheme string
//	remains := destination
//	if 1 < len(parts) {
//		scheme = parts[0]
//		remains = parts[1]
//	}
//
//	parts = strings.Split(remains, "/")
//	remains = ""
//	length := len(parts)
//	for i, part := range parts {
//		unescaped, err := url.QueryUnescape(part)
//		if nil != err {
//			unescaped = part
//		}
//		remains += unescaped
//		if i < length-1 {
//			remains += "/"
//		}
//	}
//
//	if "" == scheme {
//		ret = remains
//	} else {
//		ret = scheme + ":" + remains
//	}
//
//	return
//}
