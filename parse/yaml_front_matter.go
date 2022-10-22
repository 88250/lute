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

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

// 判断 YAML Front Matter（---）是否开始。
func YamlFrontMatterStart(t *Tree, container *ast.Node) int {
	if !t.Context.ParseOption.YamlFrontMatter || t.Context.indented || nil != t.Root.FirstChild {
		return 0
	}

	if t.parseYamlFrontMatter() {
		node := &ast.Node{Type: ast.NodeYamlFrontMatter}
		t.Root.AppendChild(node)
		t.Context.Tip = node
		return 2
	}
	return 0
}

func YamlFrontMatterContinue(node *ast.Node, context *Context) int {
	if isYamlFrontMatterClose(context) {
		context.finalize(node)
		return 2
	}
	return 0
}

var YamlFrontMatterMarker = util.StrToBytes("---")
var YamlFrontMatterMarkerNewline = util.StrToBytes("---\n")
var YamlFrontMatterMarkerCaret = util.StrToBytes("---" + editor.Caret)
var YamlFrontMatterMarkerCaretNewline = util.StrToBytes("---" + editor.Caret + "\n")

func (context *Context) yamlFrontMatterFinalize(node *ast.Node) {
	tokens := node.Tokens[3:] // 剔除开头的 ---\n
	tokens = lex.TrimWhitespace(tokens)
	if context.ParseOption.VditorWYSIWYG || context.ParseOption.VditorIR || context.ParseOption.VditorSV {
		if bytes.HasSuffix(tokens, YamlFrontMatterMarkerCaret) {
			// 剔除结尾的 ---‸
			tokens = bytes.TrimSuffix(tokens, YamlFrontMatterMarkerCaret)
			// 把 Vditor 插入符移动到内容末尾
			tokens = append(tokens, editor.CaretTokens...)
		}
	}
	if bytes.HasSuffix(tokens, YamlFrontMatterMarker) {
		tokens = tokens[:len(tokens)-3] // 剔除结尾的 ---
	}
	node.Tokens = tokens
	node.AppendChild(&ast.Node{Type: ast.NodeYamlFrontMatterOpenMarker})
	node.AppendChild(&ast.Node{Type: ast.NodeYamlFrontMatterContent, Tokens: tokens})
	node.AppendChild(&ast.Node{Type: ast.NodeYamlFrontMatterCloseMarker})
}

func (t *Tree) parseYamlFrontMatter() bool {
	if lex.ItemHyphen != t.Context.currentLine[0] {
		return false
	}

	hyphenLength := 0
	for i := 0; i < t.Context.currentLineLen && lex.ItemHyphen == t.Context.currentLine[i]; i++ {
		hyphenLength++
	}
	return 3 == hyphenLength
}

func isYamlFrontMatterClose(context *Context) bool {
	if context.ParseOption.KramdownBlockIAL && simpleCheckIsBlockIAL(context.currentLine) {
		// 判断 IAL 打断
		if ial := context.parseKramdownBlockIAL(context.currentLine); 0 < len(ial) {
			context.Tip.ID = IAL2Map(ial)["id"]
			context.Tip.KramdownIAL = ial
			context.Tip.InsertAfter(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: context.currentLine})
			return true
		}
	}

	if lex.ItemHyphen != context.currentLine[0] {
		return false
	}

	hyphenLength := 0
	for i := 0; i < context.currentLineLen && lex.ItemHyphen == context.currentLine[i]; i++ {
		hyphenLength++
	}
	return 3 == hyphenLength
}
