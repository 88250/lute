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
	"unicode/utf8"
)

// nextLine 返回下一行。
func (l *lexer) nextLine() (ret []byte) {
	if l.offset >= l.length {
		return
	}

	l.ln++
	l.col = 0

	var b, nb byte
	i := l.offset
	for ; i < l.length; i += l.width {
		b = l.input[i]
		l.col++
		if itemNewline == b {
			i++
			break
		} else if itemCarriageReturn == b {
			// 处理 \r
			if i < l.length-1 {
				nb = l.input[i+1]
				if itemNewline == nb {
					l.input = append(l.input[:i], l.input[i+1:]...) // 移除 \r，依靠下一个的 \n 切行
					l.length--                                      // 重新计算总长
				}
			}
			i++
			break
		} else if '\u0000' == b {
			// 将 \u0000 替换为 \uFFFD
			l.input = append(l.input, 0, 0)
			copy(l.input[i+2:], l.input[i:])
			// \uFFFD 的 UTF-8 编码为 \xEF\xBF\xBD 共三个字节
			l.input[i] = '\xEF'
			l.input[i+1] = '\xBF'
			l.input[i+2] = '\xBD'
			l.length += 2 // 重新计算总长
			l.width = 3
			continue
		}

		if utf8.RuneSelf <= b { // 说明占用多个字节
			_, l.width = utf8.DecodeRune(l.input[i:])
		} else {
			l.width = 1
		}
	}
	ret = l.input[l.offset:i]
	l.offset = i
	return
}
