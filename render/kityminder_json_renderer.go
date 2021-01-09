// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"strings"

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
	ret.RendererFuncs[ast.NodeHeading] = ret.renderHeading
	ret.RendererFuncs[ast.NodeList] = ret.renderList
	ret.RendererFuncs[ast.NodeListItem] = ret.renderListItem
	ret.RendererFuncs[ast.NodeThematicBreak] = ret.renderThematicBreak
	ret.RendererFuncs[ast.NodeHardBreak] = ret.renderHardBreak
	ret.RendererFuncs[ast.NodeSoftBreak] = ret.renderSoftBreak
	ret.RendererFuncs[ast.NodeHTMLBlock] = ret.renderHTML
	ret.RendererFuncs[ast.NodeTable] = ret.renderTable
	ret.RendererFuncs[ast.NodeToC] = ret.renderToC
	ret.RendererFuncs[ast.NodeYamlFrontMatter] = ret.renderYamlFrontMatter
	ret.RendererFuncs[ast.NodeBlockEmbed] = ret.renderBlockEmbed
	ret.RendererFuncs[ast.NodeBlockQueryEmbed] = ret.renderBlockQueryEmbed
	ret.RendererFuncs[ast.NodeKramdownBlockIAL] = ret.renderKramdownBlockIAL
	ret.DefaultRendererFunc = ret.renderDefault
	return ret
}

func (r *KityMinderJSONRenderer) renderBlockQueryEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderBlockEmbed(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderDefault(n *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderYamlFrontMatter(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderMathBlock(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderTable(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.dataText("table TODO")
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderHTML(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderParagraph(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		md := r.formatNode(node)
		r.dataText(md)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		if next := node.Next; nil != next {
			if ast.NodeKramdownBlockIAL == next.Type {
				next = next.Next
			}
			if nil != next {
				r.comma()
			}
		}
	}
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderBlockquote(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		r.dataText("Blockquote")
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderHeading(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		md := r.formatNode(node)
		r.dataText(md)
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
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderList(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		md := r.formatNode(node)
		r.dataText(md)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderListItem(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		r.openObj()
		md := r.formatNode(node)
		r.dataText(md)
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) renderThematicBreak(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderHardBreak(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderSoftBreak(node *ast.Node, entering bool) ast.WalkStatus {
	return ast.WalkSkipChildren
}

func (r *KityMinderJSONRenderer) renderCodeBlock(node *ast.Node, entering bool) ast.WalkStatus {
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
		r.dataText("文档名 TODO")
		r.openChildren(node)
	} else {
		r.closeChildren(node)
		r.closeObj()
		r.WriteByte(lex.ItemCloseBrace)
	}
	return ast.WalkContinue
}

func (r *KityMinderJSONRenderer) dataText(text string) {
	r.WriteString("\"data\":")
	r.openObj()
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, "\n", "\\n")
	text = strings.ReplaceAll(text, "\"", "")
	text = strings.ReplaceAll(text, "'", "")
	r.WriteString("\"text\":\"" + text + "\"")
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

func (r *KityMinderJSONRenderer) comma() {
	r.WriteString(",")
}

func (r *KityMinderJSONRenderer) formatNode(node *ast.Node) string {
	renderer := NewFormatRenderer(r.Tree, r.Options)
	renderer.Writer = &bytes.Buffer{}
	renderer.NodeWriterStack = append(renderer.NodeWriterStack, renderer.Writer)
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		rendererFunc := renderer.RendererFuncs[n.Type]
		return rendererFunc(n, entering)
	})
	return strings.TrimSpace(renderer.Writer.String())
}

func headingChildren(heading *ast.Node) (ret []*ast.Node) {
	start := heading
	if nil != heading.Next && ast.NodeKramdownBlockIAL == heading.Next.Type {
		start = heading.Next.Next
	}
	currentLevel := heading.HeadingLevel
	for n := start.Next; nil != n; n = n.Next {
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
