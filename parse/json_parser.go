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
	"encoding/json"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/util"
)

func ParseJSONWithoutFix(jsonData []byte, options *Options) (ret *Tree, err error) {
	root := &ast.Node{}
	err = json.Unmarshal(jsonData, root)
	if nil != err {
		return
	}

	ret = &Tree{Name: "", ID: root.ID, Root: &ast.Node{Type: ast.NodeDocument, ID: root.ID}, Context: &Context{ParseOption: options}}
	ret.Root.KramdownIAL = Map2IAL(root.Properties)
	ret.Context.Tip = ret.Root
	if nil == root.Children {
		return
	}

	idMap := map[string]bool{}
	for _, child := range root.Children {
		genTreeByJSON(child, ret, &idMap, nil, true)
	}
	return
}

func ParseJSON(jsonData []byte, options *Options) (ret *Tree, needFix bool, err error) {
	root := &ast.Node{}
	err = json.Unmarshal(jsonData, root)
	if nil != err {
		return
	}

	ret = &Tree{Name: "", ID: root.ID, Root: &ast.Node{Type: ast.NodeDocument, ID: root.ID, Spec: root.Spec}, Context: &Context{ParseOption: options}}
	ret.Root.KramdownIAL = Map2IAL(root.Properties)
	for _, kv := range ret.Root.KramdownIAL {
		if strings.Contains(kv[1], "\n") {
			val := kv[1]
			val = strings.ReplaceAll(val, "\n", editor.IALValEscNewLine)
			ret.Root.SetIALAttr(kv[0], val)
			needFix = true
		}
	}

	ret.Context.Tip = ret.Root
	if nil == root.Children {
		newPara := &ast.Node{Type: ast.NodeParagraph, ID: ast.NewNodeID()}
		newPara.SetIALAttr("id", newPara.ID)
		ret.Root.AppendChild(newPara)
		needFix = true
		return
	}

	idMap := map[string]bool{}
	for _, child := range root.Children {
		genTreeByJSON(child, ret, &idMap, &needFix, false)
	}

	if nil == ret.Root.FirstChild {
		// 如果是空文档的话挂一个空段落上去
		newP := NewParagraph()
		ret.Root.AppendChild(newP)
		ret.Root.SetIALAttr("updated", newP.ID[:14])
	}
	return
}

func NewParagraph() (ret *ast.Node) {
	newID := ast.NewNodeID()
	ret = &ast.Node{ID: newID, Type: ast.NodeParagraph}
	ret.SetIALAttr("id", newID)
	ret.SetIALAttr("updated", newID[:14])
	return
}

func genTreeByJSON(node *ast.Node, tree *Tree, idMap *map[string]bool, needFix *bool, ignoreFix bool) {
	node.Tokens, node.Type = util.StrToBytes(node.Data), ast.Str2NodeType(node.TypeStr)
	node.Data, node.TypeStr = "", ""
	node.KramdownIAL = Map2IAL(node.Properties)
	node.Properties = nil

	if !ignoreFix {
		// 历史数据订正

		if -1 == node.Type {
			*needFix = true
			node.Type = ast.NodeParagraph
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: node.Tokens})
			node.Children = nil
		}

		switch node.Type {
		case ast.NodeList:
			if 1 > len(node.Children) {
				*needFix = true
				return // 忽略空列表
			}
		case ast.NodeListItem:
			if 1 > len(node.Children) {
				*needFix = true
				return // 忽略空列表项
			}
		case ast.NodeBlockquote:
			if 2 > len(node.Children) {
				*needFix = true
				return // 忽略空引述
			}
		case ast.NodeSuperBlock:
			if 4 > len(node.Children) {
				*needFix = true
				return // 忽略空超级块
			}
		case ast.NodeMathBlock:
			if 1 > len(node.Children) {
				*needFix = true
				return // 忽略空公式
			}
		case ast.NodeBlockQueryEmbed:
			if 1 > len(node.Children) {
				*needFix = true
				return // 忽略空查询嵌入块
			}
		}

		fixLegacyData(tree.Context.Tip, node, idMap, needFix)
	}

	tree.Context.Tip.AppendChild(node)
	tree.Context.Tip = node
	defer tree.Context.ParentTip()
	if nil == node.Children {
		return
	}
	for _, child := range node.Children {
		genTreeByJSON(child, tree, idMap, needFix, ignoreFix)
	}
	node.Children = nil
}

func fixLegacyData(tip, node *ast.Node, idMap *map[string]bool, needFix *bool) {
	if node.IsBlock() && "" == node.ID {
		node.ID = ast.NewNodeID()
		node.SetIALAttr("id", node.ID)
		*needFix = true
	}
	if "" != node.ID {
		if _, ok := (*idMap)[node.ID]; ok {
			node.ID = ast.NewNodeID()
			node.SetIALAttr("id", node.ID)
			*needFix = true
		}
		(*idMap)[node.ID] = true
	}

	switch node.Type {
	case ast.NodeIFrame:
		if bytes.Contains(node.Tokens, util.StrToBytes("iframe-content")) {
			start := bytes.Index(node.Tokens, util.StrToBytes("<iframe"))
			end := bytes.Index(node.Tokens, util.StrToBytes("</iframe>"))
			node.Tokens = node.Tokens[start : end+9]
			*needFix = true
		}
	case ast.NodeWidget:
		if bytes.Contains(node.Tokens, util.StrToBytes("http://127.0.0.1:6806")) {
			node.Tokens = bytes.ReplaceAll(node.Tokens, []byte("http://127.0.0.1:6806"), nil)
			*needFix = true
		}
	case ast.NodeList:
		if nil != node.ListData && 3 != node.ListData.Typ && 0 < len(node.Children) &&
			nil != node.Children[0].ListData && 3 == node.Children[0].ListData.Typ {
			node.ListData.Typ = 3
			*needFix = true
		}
	case ast.NodeMark:
		if 3 == len(node.Children) && "NodeText" == node.Children[1].TypeStr {
			if strings.HasPrefix(node.Children[1].Data, " ") || strings.HasSuffix(node.Children[1].Data, " ") {
				node.Children[1].Data = strings.TrimSpace(node.Children[1].Data)
				*needFix = true
			}
		}
	case ast.NodeHeading:
		if 6 < node.HeadingLevel {
			node.HeadingLevel = 6
			*needFix = true
		}
	case ast.NodeLinkDest:
		if bytes.HasPrefix(node.Tokens, []byte("assets/")) && bytes.HasSuffix(node.Tokens, []byte(" ")) {
			node.Tokens = bytes.TrimSpace(node.Tokens)
			*needFix = true
		}
	case ast.NodeText:
		if nil != tip.LastChild && ast.NodeTagOpenMarker == tip.LastChild.Type && 1 > len(node.Tokens) {
			node.Tokens = []byte("Untitled")
			*needFix = true
		}
	case ast.NodeTagCloseMarker:
		if nil != tip.LastChild {
			if ast.NodeTagOpenMarker == tip.LastChild.Type {
				tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte("Untitled")})
				*needFix = true
			} else if "" == tip.LastChild.Text() {
				tip.LastChild.Type = ast.NodeText
				tip.LastChild.Tokens = []byte("Untitled")
				*needFix = true
			}
		}
	}

	for _, kv := range node.KramdownIAL {
		if strings.Contains(kv[1], "\n") {
			val := kv[1]
			val = strings.ReplaceAll(val, "\n", editor.IALValEscNewLine)
			node.SetIALAttr(kv[0], val)
			*needFix = true
		}
	}
}
