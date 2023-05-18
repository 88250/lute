// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"bytes"
	"unicode/utf8"

	"github.com/88250/lute/html"

	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

func (context *Context) parseInlineLinkDest(tokens []byte) (passed, remains, destination []byte) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	passed = make([]byte, 0, 256)
	destination = make([]byte, 0, 256)

	isPointyBrackets := lex.ItemLess == tokens[1]
	if isPointyBrackets {
		matchEnd := false
		passed = append(passed, tokens[0], tokens[1])
		i := 2
		size := 1
		var r rune
		var dest, runes []byte
		for ; i < length; i += size {
			size = 1
			token := tokens[i]
			if lex.ItemNewline == token {
				passed = nil
				return
			}

			if token < utf8.RuneSelf {
				passed = append(passed, token)
				dest = []byte{token}
			} else {
				dest = []byte{}
				r, size = utf8.DecodeRune(tokens[i:])
				runes = util.StrToBytes(string(r))
				passed = append(passed, runes...)
				dest = append(dest, runes...)
			}
			destination = append(destination, dest...)
			if lex.ItemGreater == token && !lex.IsBackslashEscapePunct(tokens, i) {
				destination = destination[:len(destination)-1]
				matchEnd = true
				break
			}
		}

		if !matchEnd || length <= i+1 {
			passed = nil
			return
		}

		if lex.ItemGreater == tokens[i+1] || lex.ItemCloseParen == tokens[i+1] {
			passed = append(passed, tokens[i+1])
			remains = tokens[i+2:]
		} else { // 后跟空格的情况
			remains = tokens[i+1:]
		}
	} else {
		var openParens int
		i := 0
		size := 1
		var r rune
		var dest, runes []byte
		destStarted := false
		for ; i < length; i += size {
			size = 1
			token := tokens[i]
			if token < utf8.RuneSelf {
				passed = append(passed, token)
				dest = []byte{token}
			} else {
				dest = []byte{}
				r, size = utf8.DecodeRune(tokens[i:])
				runes = util.StrToBytes(string(r))
				passed = append(passed, runes...)
				dest = append(dest, runes...)
			}
			destination = append(destination, dest...)
			if !destStarted && !lex.IsWhitespace(token) && 0 < i {
				destStarted = true
				destination = destination[1:]
				destination = lex.TrimWhitespace(destination)
			}
			if !context.ParseOption.ImgPathAllowSpace {
				if destStarted && (lex.IsWhitespace(token) || lex.IsControl(token)) {
					destination = destination[:len(destination)-size]
					passed = passed[:len(passed)-1]
					openParens--
					break
				}
			} else {
				if destStarted && lex.IsWhitespace(token) && i+1 < length {
					nextToken := tokens[i+1]
					if '"' == nextToken || '\'' == nextToken {
						destination = destination[:len(destination)-size]
						passed = passed[:len(passed)-1]
						openParens--
						break
					}
				}
			}
			if lex.ItemOpenParen == token && !lex.IsBackslashEscapePunct(tokens, i) {
				openParens++
			}
			if lex.ItemCloseParen == token && !lex.IsBackslashEscapePunct(tokens, i) {
				openParens--
				if 1 > openParens {
					if lex.ItemOpenParen == destination[0] {
						destination = destination[1:]
					}
					destination = destination[:len(destination)-1]
					break
				}
			}
		}

		remains = tokens[i:]
		if length > i && (lex.ItemCloseParen != tokens[i] && lex.ItemSpace != tokens[i] && lex.ItemNewline != tokens[i]) {
			passed = nil
			return
		}

		if 0 != openParens {
			passed = nil
			return
		}
	}

	if (context.ParseOption.ProtyleWYSIWYG || !context.ParseOption.DataImage) && bytes.HasPrefix(bytes.ToLower(destination), []byte("data:image")) {
		return nil, nil, nil
	}

	if nil != passed {
		if (!context.ParseOption.VditorWYSIWYG && !context.ParseOption.VditorIR && !context.ParseOption.VditorSV && !context.ParseOption.ProtyleWYSIWYG) &&
			!context.ParseOption.ImgPathAllowSpace {
			destination = html.EncodeDestination(html.UnescapeBytes(destination))
		}
	}
	return
}
