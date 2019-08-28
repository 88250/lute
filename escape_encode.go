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
			tmp[i], tmp[i+1], tmp[i+2], tmp[i+3], tmp[i+4] = '&', 'a', 'm', 'p', ';'
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
			tmp[i], tmp[i+1], tmp[i+2], tmp[i+3] = '&', 'l', 't', ';'
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
			tmp[i], tmp[i+1], tmp[i+2], tmp[i+3] = '&', 'g', 't', ';'
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
			tmp[i], tmp[i+1], tmp[i+2], tmp[i+3], tmp[i+4], tmp[i+5] = '&', 'q', 'u', 'o', 't', ';'
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

func unescapeString(tokens items) (ret items) {
	if nil == tokens {
		return
	}

	tokens = toItems(html.UnescapeString(fromItems(tokens))) // FIXME: 此处应该用内部的实体转义方式
	length := len(tokens)
	ret = make(items, 0, length)
	for i := 0; i < length; i++ {
		if tokens.isBackslashEscapePunct(i) {
			ret = ret[:len(ret)-1]
		}
		ret = append(ret, tokens[i])
	}
	return
}

// encodeDestination percent-encodes rawurl, avoiding double encoding.
// It doesn't touch:
// - alphanumeric characters ([0-9a-zA-Z]);
// - percent-encoded characters (%[0-9a-fA-F]{2});
// - excluded characters ([;/?:@&=+$,-_.!~*'()#]).
// Invalid UTF-8 sequences are replaced with U+FFFD.
func encodeDestination(rawurl items) items {
	// 鸣谢 https://gitlab.com/golang-commonmark/mdurl

	const hexdigit = "0123456789ABCDEF"
	var buf bytes.Buffer
	i := 0
	for i < len(rawurl) {
		r, rlen := utf8.DecodeRune(rawurl[i:])
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
	return buf.Bytes()
}
