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
	"strings"
	"unicode/utf8"
)

var (
	amp  = strToBytes("&amp;")
	lt   = strToBytes("&lt;")
	gt   = strToBytes("&gt;")
	quot = strToBytes("&quot;")
)

func escapeHTML(html []byte) (ret []byte) {
	length := len(html)
	var start, i int
	inited := false
	ret = html
	for ; i < length; i++ {
		switch html[i] {
		case itemAmpersand:
			if !inited { // 通过延迟初始化减少内存分配，下同
				ret = make([]byte, 0, length+128)
				inited = true
			}
			ret = append(ret, html[start:i]...)
			ret = append(ret, amp...)
			start = i + 1
		case itemLess:
			if !inited {
				ret = make([]byte, 0, length+128)
				inited = true
			}
			ret = append(ret, html[start:i]...)
			ret = append(ret, lt...)
			start = i + 1
		case itemGreater:
			if !inited {
				ret = make([]byte, 0, length+128)
				inited = true
			}
			ret = append(ret, html[start:i]...)
			ret = append(ret, gt...)
			start = i + 1
		case itemDoublequote:
			if !inited {
				ret = make([]byte, 0, length+128)
				inited = true
			}
			ret = append(ret, html[start:i]...)
			ret = append(ret, quot...)
			start = i + 1
		}
	}
	if inited {
		ret = append(ret, html[start:]...)
	}
	return
}

// encodeDestination percent-encodes rawurl, avoiding double encoding.
// It doesn't touch:
// - alphanumeric characters ([0-9a-zA-Z]);
// - percent-encoded characters (%[0-9a-fA-F]{2});
// - excluded characters ([;/?:@&=+$,-_.!~*'()#]).
// Invalid UTF-8 sequences are replaced with U+FFFD.
func encodeDestination(rawurl []byte) (ret []byte) {
	// 鸣谢 https://gitlab.com/golang-commonmark/mdurl

	const hexdigit = "0123456789ABCDEF"
	ret = make([]byte, 0, 256)
	i := 0
	var token byte
	for i < len(rawurl) {
		r, rlen := utf8.DecodeRune(rawurl[i:])
		if utf8.RuneSelf <= r {
			for j, n := i, i+rlen; j < n; j++ {
				b := rawurl[j]
				token = rawurl[j]
				ret = append(ret, '%')
				ret = append(ret, hexdigit[(b>>4)&0xf])
				ret = append(ret, hexdigit[b&0xf])
			}
		} else if r == '%' {
			token = rawurl[i]
			if i+2 < len(rawurl) && isHexDigit(rawurl[i+1]) && isHexDigit(rawurl[i+2]) {
				ret = append(ret, '%')
				ret = append(ret, tokenToUpper(rawurl[i+1]))
				ret = append(ret, tokenToUpper(rawurl[i+2]))
				i += 2
			} else {
				ret = append(ret, '%')
				ret = append(ret, '2')
				ret = append(ret, '5')
			}
		} else if strings.IndexByte("!#$&'()*+,-./0123456789:;=?@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~", byte(r)) == -1 {
			token = rawurl[i]
			ret = append(ret, '%')
			ret = append(ret, hexdigit[(r>>4)&0xf])
			ret = append(ret, hexdigit[r&0xf])
		} else {
			token = rawurl[i]
			ret = append(ret, token)
		}
		i += rlen
	}
	return
}
