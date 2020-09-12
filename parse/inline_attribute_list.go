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

var openCurlyBraceColon = util.StrToBytes("{:")
var emptyIAL = util.StrToBytes("{:}")

func (t *Tree) parseKramdownIALs() {
	if !t.Context.Option.KramdownIAL {
		return
	}

	t.parseKramdownIAL0(t.Root)
}

func (t *Tree) parseKramdownIAL0(node *ast.Node) {
	if ast.NodeText == node.Type {
		if curlyBracesStart := bytes.Index(node.Tokens, []byte("{:")); 0 <= curlyBracesStart {
			content := node.Tokens[curlyBracesStart+2:]
			curlyBracesEnd := bytes.Index(content, closeCurlyBrace)
			if 3 > curlyBracesEnd {
				goto Continue
			}

			content = content[:len(content)-1]
			for {
				valid, remains, attr, name, val := t.parseTagAttr(content)
				if !valid {
					break
				}

				content = remains
				if 1 > len(attr) {
					break
				}

				node.Parent.KramdownIAL = append(node.Parent.KramdownIAL, []string{util.BytesToStr(name), util.BytesToStr(val)})
				node.Tokens = bytes.Replace(node.Tokens, attr, nil, 1)
			}

			if bytes.Equal(emptyIAL, node.Tokens) {
				if nil != node.Previous && ast.NodeSoftBreak == node.Previous.Type {
					node.Previous.Unlink()
				}
				parent := node.Parent
				node.Unlink()
				if nil != parent && nil == parent.FirstChild { // 如果父节点已经没有子节点，说明这个父节点应该指向它的前一个兄弟节点
					parent.Previous.KramdownIAL = parent.KramdownIAL
					parent.Unlink()
				}
			}
		}
	}

Continue: // 遍历处理子节点
	for child := node.FirstChild; nil != child; child = child.Next {
		t.parseKramdownIAL0(child)
	}
}
