// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"github.com/88250/lute/lex"
	"strings"
	"unicode/utf8"

	"github.com/88250/lute/html"
)

var (
	amp  = StrToBytes("&amp;")
	lt   = StrToBytes("&lt;")
	gt   = StrToBytes("&gt;")
	quot = StrToBytes("&quot;")
)

func UnescapeHTML(h []byte) (ret []byte) {
	return StrToBytes(html.UnescapeString(BytesToStr(h)))
}

func EscapeHTML(html []byte) (ret []byte) {
	length := len(html)
	var start, i int
	inited := false
	ret = html
	for ; i < length; i++ {
		switch html[i] {
		case lex.ItemAmpersand:
			if !inited { // 通过延迟初始化减少内存分配，下同
				ret = make([]byte, 0, length+128)
				inited = true
			}
			ret = append(ret, html[start:i]...)
			ret = append(ret, amp...)
			start = i + 1
		case lex.ItemLess:
			if !inited {
				ret = make([]byte, 0, length+128)
				inited = true
			}
			ret = append(ret, html[start:i]...)
			ret = append(ret, lt...)
			start = i + 1
		case lex.ItemGreater:
			if !inited {
				ret = make([]byte, 0, length+128)
				inited = true
			}
			ret = append(ret, html[start:i]...)
			ret = append(ret, gt...)
			start = i + 1
		case lex.ItemDoublequote:
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

// EncodeDestination percent-encodes rawurl, avoiding double encoding.
// It doesn't touch:
// - alphanumeric characters ([0-9a-zA-Z]);
// - percent-encoded characters (%[0-9a-fA-F]{2});
// - excluded characters ([;/?:@&=+$,-_.!~*'()#]).
// Invalid UTF-8 sequences are replaced with U+FFFD.
func EncodeDestination(rawurl []byte) (ret []byte) {
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
			if i+2 < len(rawurl) && lex.IsHexDigit(rawurl[i+1]) && lex.IsHexDigit(rawurl[i+2]) {
				ret = append(ret, '%')
				ret = append(ret, lex.TokenToUpper(rawurl[i+1]))
				ret = append(ret, lex.TokenToUpper(rawurl[i+2]))
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
