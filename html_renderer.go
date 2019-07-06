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

func NewHTMLRenderer() (ret *Renderer) {
	ret = &Renderer{rendererFuncs: map[NodeType]RendererFunc{}}

	ret.rendererFuncs[NodeRoot] = ret.renderRoot
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraph
	ret.rendererFuncs[NodeText] = ret.renderText
	ret.rendererFuncs[NodeInlineCode] = ret.renderInlineCode
	ret.rendererFuncs[NodeCode] = ret.renderCode
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasis
	ret.rendererFuncs[NodeStrong] = ret.renderStrong

	return
}

func (r *Renderer) renderRoot(node Node, entering bool) (WalkStatus, error) {
	return WalkContinue, nil
}

func (r *Renderer) renderParagraph(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writer.WriteString("<p>")
	} else {
		r.writer.WriteString("</p>\n")
	}

	return WalkContinue, nil
}

func (r *Renderer) renderText(node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	n := node.(*Text)
	r.writer.WriteString(n.Value)

	return WalkContinue, nil
}

func (r *Renderer) renderInlineCode(n Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writer.WriteString("<code>" + n.(*InlineCode).Value)

		return WalkSkipChildren, nil
	}
	r.writer.WriteString("</code>")
	return WalkContinue, nil
}

func (r *Renderer) renderCode(n Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writer.WriteString("<pre><code>" + n.(*Code).Value)

		return WalkSkipChildren, nil
	}
	r.writer.WriteString("</code></pre>\n")
	return WalkContinue, nil
}

func (r *Renderer) renderEmphasis(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writer.WriteString("<em>" + node.(*Emphasis).rawText)
	} else {
		r.writer.WriteString("</em>")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrong(node Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writer.WriteString("<strong>" + node.(*Emphasis).rawText)
	} else {
		r.writer.WriteString("</strong>")
	}
	return WalkContinue, nil
}
