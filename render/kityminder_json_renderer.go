// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"bytes"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
)

// KityMinderJSONRenderer 描述了 KityMinder JSON 渲染器。
type KityMinderJSONRenderer struct {
	*BaseRenderer
}

// NewKityMinderJSONRenderer 创建一个 KityMinder JSON 渲染器。
func NewKityMinderJSONRenderer(tree *parse.Tree, options *Options) Renderer {
	ret := &KityMinderJSONRenderer{NewBaseRenderer(tree, options)}
	ret.RendererFuncs[ast.NodeDocument] = ret.renderDocument
	ret.RendererFuncs[ast.NodeParagraph] = ret.renderParagraph
	ret.RendererFuncs[ast.NodeCodeBlock] = ret.renderCodeBlock
	ret.RendererFuncs[ast.NodeMathBlock] = ret.renderMathBlock
	ret.RendererFuncs[ast.NodeBlockquote] = ret.renderBlockquote
	ret.RendererFuncs[ast.NodeSuperBlock] = ret.renderSuperBlock
	ret.RendererFuncs[ast.NodeHeading] = ret.renderHeading
	ret.RendererFuncs[ast.NodeList] = ret.renderList
	ret.RendererFuncs[ast.NodeListItem] = ret.renderListItem
	ret.RendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	ret.RendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	ret.RendererFuncs[ast.NodeTable] = ret.renderTable
	ret.RendererFuncs[ast.NodeToC] = ret.renderToC
	ret.RendererFuncs[ast.NodeYamlFrontMatter] = ret.renderYamlFrontMatter
	ret.RendererFuncs[ast.NodeBlockQueryEmbed] = ret.renderBlockQueryEmbed
	ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL
	ret.DefaultRendererFunc = ret.renderDefault
	return ret
}

func (r *KityMinderJSONRenderer) renderDefault(n *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if grandparent := node.Parent.Parent; nil != grandparent && ast.NodeList == grandparent.Type && grandparent.ListData.Tight { // List.ListItem.Paragraph
		if node.Parent.FirstChild == node && node.Parent.LastChild == node { // ListItem 下面只有一个段落时不渲染该段落
			return ast.WalkContinue
		}
	}

	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderSuperBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)

		for c := node.FirstChild; nil != c; c = c.Next {
			c.Unlink()
		}

		children := headingChildren(node)
		for _, c := range children {
			node.AppendChild(c)
		}
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.comma(node)
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderKramdownBlockIAL(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderDocument(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.WriteByte(lex.ItemOpenBrace)
		r.WriteString("\"root\":")
		r.openObj()
		r.data(node)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.WriteByte(lex.ItemCloseBrace)
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) data(node *ast.Node) {
	r.WriteString("\"data\":")
	r.openObj()

	var text string
	switch node.Type {
	case ast.NodeDocument:
		text = r.Tree.Name
	case ast.NodeList:
		if 0 == node.ListData.Typ {
			r.WriteString("\"priority\": \"iconList\",")
		} else if 1 == node.ListData.Typ {
			r.WriteString("\"priority\": \"iconOrderedList\",")
		} else {
			r.WriteString("\"priority\": \"iconCheck\",")
		}
	case ast.NodeBlockquote:
		r.WriteString("\"priority\": \"iconQuote\",")
	case ast.NodeSuperBlock:
		r.WriteString("\"priority\": \"iconSuper\",")
	default:
		buf := &bytes.Buffer{}
		ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
			if !entering {
				return ast.WalkContinue
			}

			if ast.NodeTag == n.Type {
				buf.WriteString("#" + n.Text() + "#")
				return ast.WalkSkipChildren
			}

			if ast.NodeText == n.Type || ast.NodeLinkText == n.Type || ast.NodeBlockRefText == n.Type || ast.NodeBlockRefDynamicText == n.Type ||
				ast.NodeCodeSpanContent == n.Type || ast.NodeCodeBlockCode == n.Type || ast.NodeLinkTitle == n.Type || ast.NodeMathBlockContent == n.Type ||
				ast.NodeInlineMathContent == n.Type || ast.NodeYamlFrontMatterContent == n.Type {
				buf.Write(n.Tokens)
			}
			return ast.WalkContinue
		})
		text = buf.String()
	}

	replacer := strings.NewReplacer("\\", "",
		"\n", "",
		"\"", "",
		"\t", "",
		"'", "")
	text = replacer.Replace(text)

	text = strings.ReplaceAll(text, "'", "")
	if 16 < utf8.RuneCountInString(text) {
		text = SubStr(text, 16) + "..."
	}
	if ast.NodeDocument == node.Type {
		r.WriteString("\"layout\":\"right\",")
	}
	r.WriteString("\"text\":\"" + text + "\",")
	r.WriteString("\"id\":\"" + node.IALAttr("id") + "\",")
	r.WriteString("\"type\":\"" + node.Type.String() + "\",")
	r.WriteString("\"isContainer\":" + strconv.FormatBool(node.IsContainerBlock()))
	r.closeObj()
}

func (r *KityMinderJSONRenderer) openObj() {
	r.WriteByte('{')
}

func (r *KityMinderJSONRenderer) closeObj() {
	r.WriteByte('}')
}

func (r *KityMinderJSONRenderer) openChildren(node *ast.Node) {
	if nil != node.FirstChild {
		r.WriteString(",\"children\":[")
	}
}

func (r *KityMinderJSONRenderer) closeChildren(node *ast.Node) {
	if nil != node.FirstChild {
		r.WriteByte(']')
	}
}

func (r *KityMinderJSONRenderer) comma(node *ast.Node) {
	if next := node.Next; nil != next {
		for ; nil != next; next = next.Next {
			if ast.NodeKramdownBlockIAL != next.Type {
				break
			}
		}
		if nil != next && next.IsBlock() {
			r.WriteString(",")
		}
	}
}

func headingChildren(heading *ast.Node) (ret []*ast.Node) {
	start := heading
	if nil != heading.Next && ast.NodeKramdownBlockIAL == heading.Next.Type {
		start = heading.Next.Next
	}
	currentLevel := heading.HeadingLevel
	for n := start; nil != n; n = n.Next {
		if ast.NodeHeading == n.Type {
			if currentLevel >= n.HeadingLevel {
				break
			}
		}
		if ast.NodeKramdownBlockIAL == n.Type {
			if !bytes.Contains(n.Tokens, []byte("type=\"doc\"")) {
				ret = append(ret, n)
			}
		} else {
			ret = append(ret, n)
		}
	}
	return
}
