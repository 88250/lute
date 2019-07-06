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

	return
}

func (r *Renderer) renderRoot(node Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *Renderer) renderParagraph(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.WriteString("<p>")
	} else {
		r.WriteString("</p>\n")
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

func (r *Renderer) renderCode(n Node, entering bool) (WalkStatus, error) {
	if entering {
		r.WriteString("<pre><code>" + n.(*Code).Value)

		return WalkSkipChildren, nil
	}
	r.WriteString("</code></pre>\n")
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
		r.WriteString("<blockquote>\n")
	} else {
		r.WriteString("</blockquote>\n")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderHeading(node Node, entering bool) (WalkStatus, error) {
	n := node.(*Heading)
	if entering {
		r.WriteString("<h" + " 123456"[n.Depth:n.Depth+1] + ">")
	} else {
		r.WriteString("</h" + " 123456"[n.Depth:n.Depth+1] + ">\n")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderList(node Node, entering bool) (WalkStatus, error) {
	n := node.(*List)
	tag := "ul"
	if ListTypeOrdered == n.ListType {
		tag = "ol"
	}
	if entering {
		r.WriteString("<" + tag)
		if ListTypeOrdered == n.ListType && n.Start != 1 {
			r.WriteString(fmt.Sprintf(" start=\"%d\">\n", n.Start))
		} else {
			r.WriteString(">\n")
		}
	} else {
		r.WriteString("</" + tag + ">\n")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListItem(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.WriteString("<li>")
		fc := node.FirstChild()
		if fc != nil {
			if _, ok := fc.(*Text); !ok {
				r.WriteString("\n")
			}
		}
	} else {
		r.WriteString("</li>\n")
	}
	return WalkContinue, nil
}
