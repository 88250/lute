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

import "bytes"

// Lexer 描述了词法分析器结构。
type Lexer struct {
	input  []byte // 输入的文本字节数组
	length int    // 输入的文本字节数组的长度
	offset int    // 当前读取字节位置
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

	start := l.offset
	hasNUL := false
	i := l.offset
	for ; i < l.length; i++ {
		b := l.input[i]
		if ItemNewline == b {
			i++
			break
		} else if ItemCarriageReturn == b {
			// 返回当前行时规范化行尾，并跳过 CRLF 中的 LF，不移动后续行。
			l.input[i] = ItemNewline
			i++
			ret = l.input[start:i]
			if i < l.length && ItemNewline == l.input[i] {
				i++
			}
			break
		} else if 0 == b {
			hasNUL = true
		}
	}
	if nil == ret {
		ret = l.input[start:i]
	}
	if hasNUL {
		// 只为包含 NUL 的当前行分配替换缓冲，保留其他字节（包括无效 UTF-8）。
		ret = bytes.ReplaceAll(ret, []byte{0}, []byte("\uFFFD"))
	}
	l.offset = i
	return
}
