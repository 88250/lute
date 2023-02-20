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
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"
	"github.com/88250/lute/util"
)

// IALStart 判断 kramdown 块级内联属性列表（{: attrs}）是否开始。
func IALStart(t *Tree, container *ast.Node) int {
	if !t.Context.ParseOption.KramdownBlockIAL || t.Context.indented {
		return 0
	}

	if ast.NodeListItem == t.Context.Tip.Type && nil == t.Context.Tip.FirstChild { // 在列表最终化过程中处理
		return 0
	}

	if ial := t.parseKramdownBlockIAL(); nil != ial {
		t.Context.closeUnmatchedBlocks()
		t.Context.offset = t.Context.currentLineLen // 整行过
		if util.IsDocIAL2(ial) {                    // 文档块 IAL
			t.Context.rootIAL = &ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: t.Context.currentLine[t.Context.nextNonspace:]}
			t.Root.KramdownIAL = ial
			t.Root.ID = ial[0][1]
			t.ID = t.Root.ID
			return 2
		}

		lastMatchedContainer := t.Context.lastMatchedContainer
		if t.Context.allClosed {
			if ast.NodeDocument == lastMatchedContainer.Type || ast.NodeListItem == lastMatchedContainer.Type || ast.NodeBlockquote == lastMatchedContainer.Type || ast.NodeSuperBlock == lastMatchedContainer.Type {
				lastMatchedContainer = t.Context.Tip.LastChild // 挂到最后一个子块上
				if nil == lastMatchedContainer {
					lastMatchedContainer = t.Context.lastMatchedContainer
				}
				if (ast.NodeSuperBlockLayoutMarker == lastMatchedContainer.Type || // 三个空块合并的超级块导出模版后使用会变成两个块  https://github.com/siyuan-note/siyuan/issues/4692
					ast.NodeKramdownBlockIAL == lastMatchedContainer.Type) &&
					nil != lastMatchedContainer.Parent { // 两个连续的 IAL
					tokens := IAL2Tokens(ial)
					if !bytes.HasPrefix(lastMatchedContainer.Tokens, tokens) { // 有的块解析已经做过打断处理
						// 在两个连续的 IAL 之间插入空段落，这样能够保持空段落
						p := &ast.Node{Type: ast.NodeParagraph, Tokens: []byte(" ")}
						lastMatchedContainer.InsertAfter(p)
						t.Context.Tip = p
						lastMatchedContainer = p
					}
				} else if ast.NodeBlockquoteMarker == lastMatchedContainer.Type { // 引述块下没有段落子块，需要构建一个空的段落块挂上去
					p := &ast.Node{Type: ast.NodeParagraph, Tokens: []byte(" ")}
					lastMatchedContainer.InsertAfter(p)
					t.Context.Tip = p
					lastMatchedContainer = p
				} else if ast.NodeDocument == lastMatchedContainer.Type {
					// 第一个节点是 IAL 的话需要保留空段落
					p := &ast.Node{Type: ast.NodeParagraph, Tokens: []byte(" ")}
					lastMatchedContainer.AppendChild(p)
					t.Context.Tip = p
					lastMatchedContainer = p
				}
			}
		}
		lastMatchedContainer.KramdownIAL = ial
		ialMap := IAL2MapUnEsc(ial)
		lastMatchedContainer.ID = ialMap["id"]
		node := t.Context.addChild(ast.NodeKramdownBlockIAL)
		node.Tokens = t.Context.currentLine[t.Context.nextNonspace:]
		return 2
	}
	return 0
}

var openCurlyBraceColon = util.StrToBytes("{: ")
var emptyIAL = util.StrToBytes("{:}")

func IAL2Tokens(ial [][]string) []byte {
	buf := bytes.Buffer{}
	buf.WriteString("{: ")
	for i, kv := range ial {
		buf.WriteString(kv[0])
		buf.WriteString("=\"")
		buf.WriteString(kv[1])
		buf.WriteByte('"')
		if i < len(ial)-1 {
			buf.WriteByte(' ')
		}
	}
	buf.WriteByte('}')
	return buf.Bytes()
}

func IALVal(ial *ast.Node, name string) string {
	array := Tokens2IAL(ial.Tokens)
	m := IAL2Map(array)
	return m[name]
}

func IALValMap(ial *ast.Node) (ret map[string]string) {
	ret = map[string]string{}
	array := Tokens2IAL(ial.Tokens)
	ret = IAL2Map(array)
	return
}

func IAL2Map(ial [][]string) (ret map[string]string) {
	ret = map[string]string{}
	for _, kv := range ial {
		ret[kv[0]] = html.UnescapeAttrVal(kv[1])
	}
	return
}

func IAL2MapUnEsc(ial [][]string) (ret map[string]string) {
	ret = map[string]string{}
	for _, kv := range ial {
		ret[kv[0]] = kv[1]
	}
	return
}

func Map2IAL(properties map[string]string) (ret [][]string) {
	ret = [][]string{}
	for k, v := range properties {
		ret = append(ret, []string{k, v})
	}
	return
}

func simpleCheckIsBlockIAL(tokens []byte) bool {
	if len("{: id=\"") >= len(tokens) {
		return false
	}
	return bytes.Contains(tokens, []byte("id=\""))
}

func Tokens2IAL(tokens []byte) (ret [][]string) {
	// tokens 开头必须是空格
	tokens = bytes.TrimRight(tokens, " \n")
	tokens = bytes.TrimPrefix(tokens, []byte("{:"))
	tokens = bytes.TrimSuffix(tokens, []byte("}"))
	tokens = bytes.ReplaceAll(tokens, []byte("\n"), []byte(editor.IALValEscNewLine))
	for {
		valid, remains, attr, name, val := TagAttr(tokens)
		if !valid {
			break
		}

		tokens = remains
		if 1 > len(attr) {
			break
		}

		val = bytes.ReplaceAll(val, []byte(editor.IALValEscNewLine), []byte("\n"))
		ret = append(ret, []string{util.BytesToStr(name), util.BytesToStr(val)})
	}
	return
}

func (t *Tree) parseKramdownBlockIAL() (ret [][]string) {
	tokens := t.Context.currentLine[t.Context.nextNonspace:]
	return t.Context.parseKramdownBlockIAL(tokens)
}

func (t *Tree) parseKramdownSpanIAL() {
	ast.Walk(t.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		switch n.Type {
		case ast.NodeEmphasis, ast.NodeStrong, ast.NodeCodeSpan, ast.NodeStrikethrough, ast.NodeTag, ast.NodeMark, ast.NodeImage, ast.NodeTextMark:
			break
		default:
			return ast.WalkContinue
		}

		if nil == n.Next || ast.NodeText != n.Next.Type {
			return ast.WalkContinue
		}

		tokens := n.Next.Tokens
		if pos, ial := t.Context.parseKramdownSpanIAL(tokens); 0 < len(ial) {
			n.KramdownIAL = ial
			n.Next.Tokens = tokens[pos+1:]
			if 1 > len(n.Next.Tokens) {
				n.Next.Unlink() // 移掉空的文本节点 {: ial}
			}
			spanIAL := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: tokens[:pos+1]}
			n.InsertAfter(spanIAL)
		}
		return ast.WalkContinue
	})
	return
}

func (context *Context) parseKramdownBlockIAL(tokens []byte) (ret [][]string) {
	if curlyBracesStart := bytes.Index(tokens, []byte("{:")); 0 == curlyBracesStart {
		tokens = tokens[curlyBracesStart+2:]
		curlyBracesEnd := bytes.LastIndex(tokens, closeCurlyBrace)
		if 3 > curlyBracesEnd {
			return
		}

		if !bytes.Equal(tokens[curlyBracesEnd:], []byte("}\n")) { // IAL 后不能存在其他内容，必须独占一行
			return
		}
		ret = Tokens2IAL(tokens)
	}
	return
}

func (context *Context) parseKramdownSpanIAL(tokens []byte) (pos int, ret [][]string) {
	pos = bytes.Index(tokens, closeCurlyBrace)
	if curlyBracesStart := bytes.Index(tokens, []byte("{:")); 0 == curlyBracesStart && curlyBracesStart+2 < pos {
		tokens = tokens[curlyBracesStart+2:]
		curlyBracesEnd := bytes.Index(tokens, closeCurlyBrace)
		if 3 > curlyBracesEnd {
			return
		}

		tokens = tokens[:curlyBracesEnd]
		for {
			valid, remains, attr, name, val := TagAttr(tokens)
			if !valid {
				break
			}

			tokens = remains
			if 1 > len(attr) {
				break
			}

			nameStr := strings.ReplaceAll(util.BytesToStr(name), editor.Caret, "")
			valStr := strings.ReplaceAll(util.BytesToStr(val), editor.Caret, "")
			ret = append(ret, []string{nameStr, valStr})
		}
	}
	return
}

func (context *Context) parseKramdownIALInListItem(tokens []byte) (ret [][]string) {
	if curlyBracesStart := bytes.Index(tokens, []byte("{:")); 0 == curlyBracesStart {
		tokens = tokens[curlyBracesStart+2:]
		curlyBracesEnd := bytes.Index(tokens, closeCurlyBrace)
		if 3 > curlyBracesEnd {
			return
		}

		tokens = tokens[:bytes.Index(tokens, []byte("}"))]
		for {
			valid, remains, attr, name, val := TagAttr(tokens)
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
