// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"github.com/88250/lute/ast"
	"github.com/88250/lute/util"
)

var openCurlyBraceColon = util.StrToBytes("{: ")
var emptyIAL = util.StrToBytes("{:}")

func (t *Tree) parseKramdownBlockIAL() (ret [][]string) {
	tokens := t.Context.currentLine[t.Context.nextNonspace:]
	return t.Context.parseKramdownBlockIAL(tokens)
}

func (t *Tree) parseKramdownSpanIAL(block *ast.Node, ctx *InlineContext) (ret *ast.Node) {
	if !t.Context.Option.KramdownIAL {
		return nil
	}

	span := block.LastChild
	if nil == span {
		return nil
	}

	switch span.Type {
	case ast.NodeEmphasis, ast.NodeStrong, ast.NodeImage:
		break
	default:
		return nil
	}

	tokens := ctx.tokens[ctx.pos:]
	if pos, ial := t.Context.parseKramdownSpanIAL(tokens); 0 < len(ial) {
		span.KramdownIAL = ial
		tokens = tokens[:pos]
		ret = &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: tokens}
		ctx.pos += pos
		return
	}
	return
}

func (context *Context) parseKramdownBlockIAL(tokens []byte) (ret [][]string) {
	if curlyBracesStart := bytes.Index(tokens, []byte("{:")); 0 == curlyBracesStart {
		tokens = tokens[curlyBracesStart+2:]
		curlyBracesEnd := bytes.Index(tokens, closeCurlyBrace)
		if 3 > curlyBracesEnd {
			return
		}

		if !bytes.Equal(tokens[curlyBracesEnd:], []byte("}\n")) { // IAL 后不能存在其他内容，必须独占一行
			return
		}
		tokens = tokens[:len(tokens)-2]
		for {
			valid, remains, attr, name, val := context.Tree.parseTagAttr(tokens)
			if !valid {
				break
			}

			tokens = remains
			if 1 > len(attr) {
				break
			}

			ret = append(ret, []string{util.BytesToStr(name), util.BytesToStr(val)})
		}
	}
	return
}

func (context *Context) parseKramdownSpanIAL(tokens []byte) (pos int, ret [][]string) {
	if curlyBracesStart := bytes.Index(tokens, []byte("{:")); 0 <= curlyBracesStart {
		tokens = tokens[curlyBracesStart+2:]
		curlyBracesEnd := bytes.Index(tokens, closeCurlyBrace)
		if 3 > curlyBracesEnd {
			return
		}

		pos := bytes.Index(tokens, []byte("}"))
		tokens = tokens[:pos]
		for {
			valid, remains, attr, name, val := context.Tree.parseTagAttr(tokens)
			if !valid {
				break
			}

			tokens = remains
			if 1 > len(attr) {
				break
			}

			ret = append(ret, []string{util.BytesToStr(name), util.BytesToStr(val)})
		}
	}
	return
}

func (context *Context) parseKramdownIALInListItem(tokens []byte) (ret [][]string) {
	if curlyBracesStart := bytes.Index(tokens, []byte("{:")); 0 <= curlyBracesStart {
		tokens = tokens[curlyBracesStart+2:]
		curlyBracesEnd := bytes.Index(tokens, closeCurlyBrace)
		if 3 > curlyBracesEnd {
			return
		}

		tokens = tokens[:bytes.Index(tokens, []byte("}"))]
		for {
			valid, remains, attr, name, val := context.Tree.parseTagAttr(tokens)
			if !valid {
				break
			}

			tokens = remains
			if 1 > len(attr) {
				break
			}

			ret = append(ret, []string{util.BytesToStr(name), util.BytesToStr(val)})
		}
	}
	return
}
