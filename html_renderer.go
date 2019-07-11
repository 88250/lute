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
	"fmt"
)

func NewHTMLRenderer() (ret *Renderer) {
	ret = &Renderer{rendererFuncs: map[NodeType]RendererFunc{}}

	ret.rendererFuncs[NodeRoot] = ret.renderRoot
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraph
	ret.rendererFuncs[NodeText] = ret.renderText
	ret.rendererFuncs[NodeInlineCode] = ret.renderInlineCode
	ret.rendererFuncs[NodeCode] = ret.renderCode
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasis
	ret.rendererFuncs[NodeStrong] = ret.renderStrong
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquote
	ret.rendererFuncs[NodeHeading] = ret.renderHeading
	ret.rendererFuncs[NodeList] = ret.renderList
	ret.rendererFuncs[NodeListItem] = ret.renderListItem
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreak
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreak
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreak
	ret.rendererFuncs[NodeHTML] = ret.renderHTML
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTML
	ret.rendererFuncs[NodeLink] = ret.renderLink

	return
}

func (r *Renderer) renderLink(node Node, entering bool) (WalkStatus, error) {
	if entering {
		n := node.(*Link)
		out := "<a href=\"" + n.URL + "\""
		if "" != n.Title {
			out += " title=\"" + n.Title + "\""
		}
		out += ">"
		r.WriteString(out)

		return WalkContinue, nil
	}

	r.WriteString("</a>")

	return WalkContinue, nil
}

func (r *Renderer) renderHTML(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	r.Newline()
	n := node.(*HTML)
	r.WriteString(n.Value)
	r.Newline()

	return WalkContinue, nil
}

func (r *Renderer) renderInlineHTML(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	n := node.(*InlineHTML)
	r.WriteString(n.Value)

	return WalkContinue, nil
}

func (r *Renderer) renderRoot(node Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *Renderer) renderParagraph(node Node, entering bool) (WalkStatus, error) {
	inTightList := false
	if NodeListItem == node.Parent().Type() {
		inTightList = node.Parent().Parent().(*List).Tight
	}
	if entering {
		if !inTightList {
			r.WriteString("<p>")
		}
	} else {
		if !inTightList {
			r.WriteString("</p>")
			r.Newline()
		}
	}

	return WalkContinue, nil
}

func (r *Renderer) renderText(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	n := node.(*Text)
	r.WriteString(n.Value)

	return WalkContinue, nil
}

func (r *Renderer) renderInlineCode(n Node, entering bool) (WalkStatus, error) {
	if entering {
		r.WriteString("<code>" + n.(*InlineCode).Value)

		return WalkSkipChildren, nil
	}
	r.WriteString("</code>")
	return WalkContinue, nil
}

func (r *Renderer) renderCode(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.Newline()
		n := node.(*Code)
		if "" != n.InfoStr {
			r.WriteString("<pre><code class=\"language-" + n.InfoStr + "\">" + n.Value)
		} else {
			r.WriteString("<pre><code>" + n.Value)
		}
		return WalkSkipChildren, nil
	}
	r.WriteString("</code></pre>")
	r.Newline()
	return WalkContinue, nil
}

func (r *Renderer) renderEmphasis(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.WriteString("<em>" + node.(*Emphasis).rawText)
	} else {
		r.WriteString("</em>")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrong(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.WriteString("<strong>" + node.(*Strong).rawText)
	} else {
		r.WriteString("</strong>")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderBlockquote(n Node, entering bool) (WalkStatus, error) {
	if entering {
		r.Newline()
		r.WriteString("<blockquote>")
		r.Newline()
	} else {
		r.Newline()
		r.WriteString("</blockquote>")
		r.Newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderHeading(node Node, entering bool) (WalkStatus, error) {
	n := node.(*Heading)
	if entering {
		r.WriteString("<h" + " 123456"[n.Level:n.Level+1] + ">")
	} else {
		r.WriteString("</h" + " 123456"[n.Level:n.Level+1] + ">")
		r.Newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderList(node Node, entering bool) (WalkStatus, error) {
	n := node.(*List)
	tag := "ul"
	if !n.Bullet {
		tag = "ol"
	}
	if entering {
		r.Newline()
		r.WriteString("<" + tag)
		if !n.Bullet && 1 != n.Start {
			r.WriteString(fmt.Sprintf(" start=\"%d\">", n.Start))
		} else {
			r.WriteString(">")
		}
		r.Newline()
	} else {
		r.Newline()
		r.WriteString("</" + tag + ">")
		r.Newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListItem(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.WriteString("<li>")
		if !node.(*ListItem).Tight {
			r.Newline()
		}
	} else {
		r.WriteString("</li>")
		r.Newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderThematicBreak(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.Newline()
		r.WriteString("<hr />")
		r.Newline()
	}

	return WalkContinue, nil
}

func (r *Renderer) renderHardBreak(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.WriteString("<br />")
	}

	return WalkContinue, nil
}

func (r *Renderer) renderSoftBreak(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.Newline()
	}

	return WalkContinue, nil
}
