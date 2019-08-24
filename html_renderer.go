// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	ret.rendererFuncs[NodeDocument] = ret.renderDocument
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraph
	ret.rendererFuncs[NodeText] = ret.renderText
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpan
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlock
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasis
	ret.rendererFuncs[NodeStrong] = ret.renderStrong
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquote
	ret.rendererFuncs[NodeHeading] = ret.renderHeading
	ret.rendererFuncs[NodeList] = ret.renderList
	ret.rendererFuncs[NodeListItem] = ret.renderListItem
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreak
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreak
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreak
	ret.rendererFuncs[NodeHTMLBlock] = ret.renderHTML
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTML
	ret.rendererFuncs[NodeLink] = ret.renderLink
	ret.rendererFuncs[NodeImage] = ret.renderImage

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

func (r *Renderer) renderImage(node Node, entering bool) (WalkStatus, error) {
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

func (r *Renderer) renderLink(node Node, entering bool) (WalkStatus, error) {
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

func (r *Renderer) renderHTML(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.newline()
	r.write(node.Tokens())
	r.newline()

	return WalkContinue, nil
}

func (r *Renderer) renderInlineHTML(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.write(node.Tokens())

	return WalkContinue, nil
}

func (r *Renderer) renderDocument(node Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *Renderer) renderParagraph(node Node, entering bool) (WalkStatus, error) {
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

func (r *Renderer) renderText(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	text := fromItems(escapeHTML(node.Tokens()))
	if r.option.AutoSpace {
		text = space(text)
	}
	r.writeString(text)

	return WalkContinue, nil
}

func (r *Renderer) renderCodeSpan(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<code>")
		r.write(escapeHTML(node.Tokens()))
		return WalkSkipChildren, nil
	}
	r.writeString("</code>")
	return WalkContinue, nil
}

func (r *Renderer) renderCodeBlock(node Node, entering bool) (WalkStatus, error) {
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

func (r *Renderer) renderEmphasis(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("em", nil, false)
	} else {
		r.tag("/em", nil, false)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrong(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("<strong>")
		r.write(node.Tokens())
	} else {
		r.writeString("</strong>")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderBlockquote(n Node, entering bool) (WalkStatus, error) {
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

func (r *Renderer) renderHeading(node Node, entering bool) (WalkStatus, error) {
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

func (r *Renderer) renderList(node Node, entering bool) (WalkStatus, error) {
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

func (r *Renderer) renderListItem(node Node, entering bool) (WalkStatus, error) {
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

func (r *Renderer) renderThematicBreak(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.tag("hr", nil, true)
		r.newline()
	}

	return WalkContinue, nil
}

func (r *Renderer) renderHardBreak(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.tag("br", nil, true)
		r.newline()
	}

	return WalkContinue, nil
}

func (r *Renderer) renderSoftBreak(node Node, entering bool) (WalkStatus, error) {
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
