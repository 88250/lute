// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lex

import "unicode/utf8"

// Lexer 描述了词法分析器结构。
type Lexer struct {
	input  []byte // 输入的文本字节数组
	length int    // 输入的文本字节数组的长度
	offset int    // 当前读取字节位置
	width  int    // 最新一个字符的长度（字节数）
}

// NewLexer 创建一个词法分析器。
func NewLexer(input []byte) (ret *Lexer) {
	ret = &Lexer{input: input, length: len(input)}
	if 0 < ret.length && ItemNewline != ret.input[ret.length-1] {
		// 以 \n 结尾预处理
		ret.input = append(ret.input, ItemNewline)
		ret.length++
	}
	return
}

// NextLine 返回下一行。
func (l *Lexer) NextLine() (ret []byte) {
	if l.offset >= l.length {
		return
	}

	var b, nb byte
	i := l.offset
	for ; i < l.length; i += l.width {
		b = l.input[i]
		if ItemNewline == b {
			i++
			break
		} else if ItemCarriageReturn == b {
			if i < l.length-1 {
				nb = l.input[i+1]
				if ItemNewline == nb { // \r\n
					l.input = append(l.input[:i], l.input[i+1:]...) // 移除 \r，依靠下一个的 \n 切行
					l.length--                                      // 重新计算总长
				} else { // \rX
					l.input[i] = ItemNewline // 将 \r 替换为 \n
				}
			} else { // \rEOF
				l.input[i] = ItemNewline // 将 \r 替换为 \n
			}
			i++
			break
		} else if '\u0000' == b {
			// 将 \u0000 替换为 \uFFFD
			l.input = append(l.input, 0, 0)
			copy(l.input[i+2:], l.input[i:])
			// \uFFFD 的 UTF-8 编码为 \xEF\xBF\xBD 共三个字节
			l.input[i], l.input[i+1], l.input[i+2] = '\xEF', '\xBF', '\xBD'
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
