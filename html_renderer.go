// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma"
	chromahtml "github.com/alecthomas/chroma/formatters/html"
	chromalexers "github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// newHTMLRenderer 创建一个 HTML 渲染器。
func newHTMLRenderer(option options) (ret *Renderer) {
	ret = &Renderer{rendererFuncs: map[int]RendererFunc{}, option: option}

	// 注册 CommonMark 渲染函数

	ret.rendererFuncs[NodeDocument] = ret.renderDocumentHTML
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraphHTML
	ret.rendererFuncs[NodeText] = ret.renderTextHTML
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpanHTML
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlockHTMl
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasisHTML
	ret.rendererFuncs[NodeStrong] = ret.renderStrongHTML
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquoteHTML
	ret.rendererFuncs[NodeHeading] = ret.renderHeadingHTML
	ret.rendererFuncs[NodeList] = ret.renderListHTML
	ret.rendererFuncs[NodeListItem] = ret.renderListItemHTML
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreakHTML
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreakHTML
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreakHTML
	ret.rendererFuncs[NodeHTMLBlock] = ret.renderHTMLHTML
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTMLHTML
	ret.rendererFuncs[NodeLink] = ret.renderLinkHTML
	ret.rendererFuncs[NodeImage] = ret.renderImageHTML

	// 注册 GFM 渲染函数

	ret.rendererFuncs[NodeStrikethrough] = ret.renderStrikethrough
	ret.rendererFuncs[NodeTaskListItemMarker] = ret.renderTaskListItemMarker
	ret.rendererFuncs[NodeTable] = ret.renderTable
	ret.rendererFuncs[NodeTableHead] = ret.renderTableHead
	ret.rendererFuncs[NodeTableRow] = ret.renderTableRow
	ret.rendererFuncs[NodeTableCell] = ret.renderTableCell

	return
}

func (r *Renderer) renderTableCell(node Node, entering bool) (WalkStatus, error) {
	tag := "td"
	if NodeTableHead == node.Parent().Type() {
		tag = "th"
	}
	if entering {
		cell := node.(*TableCell)
		var attrs [][]string
		switch cell.Aligns {
		case 1:
			attrs = append(attrs, []string{"align", "left"})
		case 2:
			attrs = append(attrs, []string{"align", "center"})
		case 3:
			attrs = append(attrs, []string{"align", "right"})
		}
		r.tag(tag, attrs, false)
	} else {
		r.tag("/"+tag, nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableRow(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("tr", nil, false)
		r.newline()
	} else {
		r.tag("/tr", nil, false)
		r.newline()
		if node == node.Parent().LastChild() {
			r.tag("/tbody", nil, false)
		}
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableHead(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("thead", nil, false)
		r.newline()
		r.tag("tr", nil, false)
		r.newline()
	} else {
		r.tag("/tr", nil, false)
		r.newline()
		r.tag("/thead", nil, false)
		r.newline()
		if nil != node.Next() {
			r.tag("tbody", nil, false)
		}
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTable(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("table", nil, false)
		r.newline()
	} else {
		r.tag("/table", nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrikethrough(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("del", nil, false)
	} else {
		r.tag("/del", nil, false)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderImageHTML(node Node, entering bool) (WalkStatus, error) {
	n := node.(*Image)
	if entering {
		if 0 == r.disableTags {
			r.writeString("<img src=\"")
			r.write(escapeHTML(toItems(n.Destination)))
			r.writeString("\" alt=\"")
		}
		r.disableTags++
		return WalkContinue, nil
	}

	r.disableTags--
	if 0 == r.disableTags {
		r.writeString("\"")
		if "" != n.Title {
			r.writeString(" title=\"")
			r.write(escapeHTML(toItems(n.Title)))
			r.writeString("\"")
		}
		r.writeString(" />")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderLinkHTML(node Node, entering bool) (WalkStatus, error) {
	if entering {
		n := node.(*Link)
		attrs := [][]string{{"href", fromItems(escapeHTML(toItems(n.Destination)))}}
		if "" != n.Title {
			attrs = append(attrs, []string{"title", fromItems(escapeHTML(toItems(n.Title)))})
		}
		r.tag("a", attrs, false)

		return WalkContinue, nil
	}

	r.tag("/a", nil, false)

	return WalkContinue, nil
}

func (r *Renderer) renderHTMLHTML(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.newline()
	r.write(node.Tokens())
	r.newline()

	return WalkContinue, nil
}

func (r *Renderer) renderInlineHTMLHTML(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.write(node.Tokens())

	return WalkContinue, nil
}

func (r *Renderer) renderDocumentHTML(node Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *Renderer) renderParagraphHTML(node Node, entering bool) (WalkStatus, error) {
	if grandparent := node.Parent().Parent(); nil != grandparent {
		if list, ok := grandparent.(*List); ok { // List.ListItem.Paragraph
			if list.tight {
				return WalkContinue, nil
			}
		}
	}

	if entering {
		r.newline()
		r.tag("p", nil, false)
	} else {
		r.tag("/p", nil, false)
		r.newline()
	}

	return WalkContinue, nil
}

func (r *Renderer) renderTextHTML(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.write(escapeHTML(node.Tokens()))

	return WalkContinue, nil
}

func (r *Renderer) renderCodeSpanHTML(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<code>")
		r.write(escapeHTML(node.Tokens()))
		return WalkSkipChildren, nil
	}
	r.writeString("</code>")
	return WalkContinue, nil
}

func (r *Renderer) renderCodeBlockHTMl(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		n := node.(*CodeBlock)
		tokens := n.tokens
		if "" != n.info {
			infoWords := strings.Fields(n.info)
			language := infoWords[0]
			r.writeString("<pre><code class=\"language-" + language + "\">")
			rendered := false
			if r.option.CodeSyntaxHighlight {
				codeBlock := fromItems(tokens)
				var lexer chroma.Lexer
				if "" != language {
					lexer = chromalexers.Get(language)
				} else {
					lexer = chromalexers.Analyse(codeBlock)
				}
				if nil == lexer {
					lexer = chromalexers.Fallback
				}
				iterator, err := lexer.Tokenise(nil, codeBlock)
				if nil == err {
					formatter := chromahtml.New(chromahtml.PreventSurroundingPre(), chromahtml.WithClasses(), chromahtml.ClassPrefix("highlight-"))
					var b bytes.Buffer
					if err = formatter.Format(&b, styles.GitHub, iterator); nil == err {
						r.write(b.Bytes())
						rendered = true
						// 生成 CSS 临时调试用：
						//formatter.WriteCSS(os.Stdout, styles.GitHub)
						//os.Stdout.WriteString("\n")
					}
				}
			}

			if !rendered {
				tokens = escapeHTML(tokens)
				r.write(tokens)
			}
		} else {
			if r.option.CodeSyntaxHighlight {
				codeBlock := fromItems(tokens)
				var lexer = chromalexers.Analyse(codeBlock)
				if nil == lexer {
					lexer = chromalexers.Fallback
				}
				language := lexer.Config().Name
				r.writeString("<pre><code class=\"language-" + language + "\">")
				rendered := false

				iterator, err := lexer.Tokenise(nil, codeBlock)
				if nil == err {
					formatter := chromahtml.New(chromahtml.PreventSurroundingPre(), chromahtml.WithClasses(), chromahtml.ClassPrefix("highlight-"))
					var b bytes.Buffer
					if err = formatter.Format(&b, styles.GitHub, iterator); nil == err {
						r.write(b.Bytes())
						rendered = true
					}
				}

				if !rendered {
					tokens = escapeHTML(tokens)
					r.write(tokens)
				}
			} else {
				r.writeString("<pre><code>")
				tokens = escapeHTML(tokens)
				r.write(tokens)
			}
		}
		return WalkSkipChildren, nil
	}
	r.writeString("</code></pre>")
	r.newline()
	return WalkContinue, nil
}

func (r *Renderer) renderEmphasisHTML(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("em", nil, false)
	} else {
		r.tag("/em", nil, false)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrongHTML(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<strong>")
		r.write(node.Tokens())
	} else {
		r.writeString("</strong>")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderBlockquoteHTML(n Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.writeString("<blockquote>")
		r.newline()
	} else {
		r.newline()
		r.writeString("</blockquote>")
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderHeadingHTML(node Node, entering bool) (WalkStatus, error) {
	n := node.(*Heading)
	if entering {
		r.newline()
		r.writeString("<h" + " 123456"[n.Level:n.Level+1] + ">")
	} else {
		r.writeString("</h" + " 123456"[n.Level:n.Level+1] + ">")
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListHTML(node Node, entering bool) (WalkStatus, error) {
	n := node.(*List)
	tag := "ul"
	if 1 == n.listData.typ {
		tag = "ol"
	}
	if entering {
		r.newline()
		attrs := [][]string{{"start", fmt.Sprintf("%d", n.start)}}
		if nil == n.bulletChar && 1 != n.start {
			r.tag(tag, attrs, false)
		} else {
			r.tag(tag, nil, false)
		}
		r.newline()
	} else {
		r.newline()
		r.tag("/"+tag, nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListItemHTML(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("li", nil, false)
	} else {
		r.tag("/li", nil, false)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTaskListItemMarker(node Node, entering bool) (WalkStatus, error) {
	if entering {
		n := node.(*TaskListItemMarker)
		var attrs [][]string
		if n.checked {
			attrs = append(attrs, []string{"checked", ""})
		}
		attrs = append(attrs, []string{"disabled", ""}, []string{"type", "checkbox"})
		r.tag("input", attrs, true)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderThematicBreakHTML(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.tag("hr", nil, true)
		r.newline()
	}

	return WalkContinue, nil
}

func (r *Renderer) renderHardBreakHTML(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("br", nil, true)
		r.newline()
	}

	return WalkContinue, nil
}

func (r *Renderer) renderSoftBreakHTML(node Node, entering bool) (WalkStatus, error) {
	if entering {
		if r.option.SoftBreak2HardBreak {
			r.tag("br", nil, true)
			r.newline()
		} else {
			r.newline()
		}
	}

	return WalkContinue, nil
}

func (r *Renderer) tag(name string, attrs [][]string, selfclosing bool) {
	if r.disableTags > 0 {
		return
	}

	r.writeString("<")
	r.write(toItems(name))
	if 0 < len(attrs) {
		for _, attr := range attrs {
			r.writeString(" " + attr[0] + "=\"" + attr[1] + "\"")
		}
	}
	if selfclosing {
		r.writeString(" /")
	}
	r.writeString(">")
}
