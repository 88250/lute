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

// +build js

package lute

import (
	"unicode/utf8"
)

// nextLine 返回下一行。
func (l *lexer) nextLine() (ret items) {
	if l.offset >= l.length {
		return
	}
	ret = make(items, 0, 256)

	l.ln++
	l.col = 0

	var b, nb byte
	i := l.offset
	for ; i < l.length; i += l.width {
		b = l.input[i]
		l.col++
		if itemNewline == b {
			i++
			ret = append(ret, newItem(b, l.ln, l.col))
			break
		} else if itemCarriageReturn == b {
			// 按照规范定义的 line ending (https://spec.commonmark.org/0.29/#line-ending) 处理 \r
			ret = append(ret, newItem(b, l.ln, l.col))
			if i < l.length-1 {
				nb = l.input[i+1]
				if itemNewline == nb {
					l.input = append(l.input[:i], l.input[i+1:]...) // 移除 \r，依靠下一个的 \n 切行
					l.length--                                      // 重新计算总长
					ret = ret[:len(ret)-1]
					ret = append(ret, newItem(nb, l.ln, l.col))
				}
			}
			i++
			break
		} else if '\u0000' == b {
			// 将 \u0000 替换为 \uFFFD https://spec.commonmark.org/0.29/#insecure-characters

			l.input = append(l.input, 0, 0)
			copy(l.input[i+2:], l.input[i:])
			// \uFFFD 的 UTF-8 编码为 \xEF\xBF\xBD 共三个字节
			l.input[i] = '\xEF'
			l.input[i+1] = '\xBF'
			l.input[i+2] = '\xBD'
			l.length += 2 // 重新计算总长
			l.width = 3
			ret = append(ret, newItem(l.input[i], l.ln, l.col))
			ret = append(ret, newItem(l.input[i+1], l.ln, l.col))
			ret = append(ret, newItem(l.input[i+2], l.ln, l.col))
			continue
		}

		if utf8.RuneSelf <= b { // 说明占用多个字节
			_, l.width = utf8.DecodeRune(l.input[i:])
			for j := 0; j < l.width; j++ {
				ret = append(ret, newItem(l.input[i+j], l.ln, l.col))
			}
		} else {
			l.width = 1
			ret = append(ret, newItem(b, l.ln, l.col))
		}
	}
	l.offset = i
	return
}
