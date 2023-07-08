// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package html

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"github.com/88250/lute/editor"
	"github.com/88250/lute/lex"
)

var (
	amp  = []byte("&amp;")
	lt   = []byte("&lt;")
	gt   = []byte("&gt;")
	quot = []byte("&quot;")
)

func UnescapeAttrVal(v string) string {
	v = strings.ReplaceAll(v, editor.IALValEscNewLine, "\n")
	v = strings.ReplaceAll(v, "&#123;", "{")
	v = strings.ReplaceAll(v, "&#125;", "}")
	return UnescapeString(v)
}

func EscapeAttrVal(v string) (ret string) {
	ret = string(EscapeHTML([]byte(v)))
	ret = strings.ReplaceAll(ret, "\n", editor.IALValEscNewLine)
	ret = strings.ReplaceAll(ret, "{", "&#123;")
	ret = strings.ReplaceAll(ret, "}", "&#125;")
	return
}

func UnescapeHTMLStr(h string) string {
	return UnescapeString(h)
}

func EscapeHTMLStr(h string) string {
	return string(EscapeHTML([]byte(h)))
}

func UnescapeHTML(h []byte) (ret []byte) {
	return []byte(UnescapeString(string(h)))
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

// DecodeDestination decodes a percent-encoded URL.
// Invalid percent-encoded sequences are left as is.
// Invalid UTF-8 sequences are replaced with U+FFFD.
func DecodeDestination(rawurl []byte) []byte {
	// 鸣谢 https://gitlab.com/golang-commonmark/mdurl

	var buf bytes.Buffer
	i := 0
	const replacement = "\xEF\xBF\xBD"
outer:
	for i < len(rawurl) {
		r, rlen := utf8.DecodeRune(rawurl[i:])
		if r == '%' && i+2 < len(rawurl) && lex.IsHexDigit(rawurl[i+1]) && lex.IsHexDigit(rawurl[i+2]) {
			b := unhex(rawurl[i+1])<<4 | unhex(rawurl[i+2])
			if b < 0x80 {
				buf.WriteByte(b)
				i += 3
				continue
			}
			var n int
			switch {
			case b&0xe0 == 0xc0:
				n = 1
			case b&0xf0 == 0xe0:
				n = 2
			case b&0xf8 == 0xf0:
				n = 3
			}
			if n == 0 {
				buf.WriteString(replacement)
				i += 3
				continue
			}
			rb := make([]byte, n+1)
			rb[0] = b
			j := i + 3
			for k := 0; k < n; k++ {
				b, j = advance(rawurl, j)
				if j > len(rawurl) || b&0xc0 != 0x80 {
					buf.WriteString(replacement)
					i += 3
					continue outer
				}
				rb[k+1] = b
			}
			r, _ := utf8.DecodeRune(rb)
			buf.WriteRune(r)
			i = j
			continue
		}
		buf.WriteRune(r)
		i += rlen
	}
	return buf.Bytes()
}

func unhex(b byte) byte {
	switch {
	case lex.IsDigit(b):
		return b - '0'
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10
	case b >= 'A' && b <= 'F':
		return b - 'A' + 10
	}
	panic("unhex: not a hex digit")
}

func advance(s []byte, pos int) (byte, int) {
	if pos >= len(s) {
		return 0, len(s) + 1
	}
	if s[pos] != '%' {
		return s[pos], pos + 1
	}
	if pos+2 < len(s) &&
		lex.IsHexDigit(s[pos+1]) &&
		lex.IsHexDigit(s[pos+2]) {
		return unhex(s[pos+1])<<4 | unhex(s[pos+2]), pos + 3
	}
	return '%', pos + 1
}
