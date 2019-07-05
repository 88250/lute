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
	"strings"
)

// NodeRendererFunc is a function that renders a given node.
type NodeRendererFunc func(writer strings.Builder, n Node, entering bool) (WalkStatus, error)

// A Renderer interface renders given AST node to given
// writer with given Renderer.
type Renderer struct {
	nodeRendererFuncs map[NodeType]NodeRendererFunc
}

func NewRenderer() *Renderer {
	return &Renderer{nodeRendererFuncs: map[NodeType]NodeRendererFunc{}}
}

func (r *Renderer) Register(nodeType NodeType, v NodeRendererFunc) {
	r.nodeRendererFuncs[nodeType] = v
}

func (r *Renderer) Render(writer strings.Builder, n Node) {
	r.nodeRendererFuncs = map[NodeType]NodeRendererFunc{}
	for kind, nr := range r.nodeRendererFuncs {
		r.nodeRendererFuncs[kind] = nr
	}

	err := Walk(n, func(n Node, entering bool) (WalkStatus, error) {
		s := WalkStatus(WalkContinue)
		var err error
		f := r.nodeRendererFuncs[n.Type()]
		if f != nil {
			s, err = f(writer, n, entering)
		}
		return s, err
	})

	_ = err
}

// HTML renderer

func (r *Renderer) renderText(writer strings.Builder, node Node, entering bool) (WalkStatus, error) {
	if !entering {
		return WalkContinue, nil
	}

	n := node.(*Text)
	writer.WriteString(n.Value)

	return WalkContinue, nil
}
