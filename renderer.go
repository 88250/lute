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
	"errors"
	"fmt"
	"strings"
)

type RendererFunc func(n Node, entering bool) (WalkStatus, error)

const Newline = "\n"

type Renderer struct {
	writer        strings.Builder
	lastOut       string
	rendererFuncs map[NodeType]RendererFunc
	disableTags   int
}

func (r *Renderer) Render(n Node) error {
	r.lastOut = Newline
	return Walk(n, func(n Node, entering bool) (WalkStatus, error) {
		f := r.rendererFuncs[n.Type()]
		if nil == f {
			return WalkStop, errors.New(fmt.Sprintf("not found render function for node [type=%d, text=%s]", n.Type(), n.RawText()))
		}

		return f(n, entering)
	})
}

func (r *Renderer) WriteString(content string) {
	r.writer.WriteString(content)
	r.lastOut = content
}

func (r *Renderer) Newline() {
	if Newline != r.lastOut {
		r.WriteString(Newline)
	}
}
