// Lute - A structural markdown engine.
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

type openerBottom struct {
	typ      string
	num      int
	position int
}

type delimiterStackElement struct {
	node             Node   // the text node point to
	typ              string // the type of delimiter ([, ![, *, _)
	num              int    // the number of delimiters
	active           bool   // whether the delimiter is "active" (all are active to start)
	openerCloserBoth string // whether the delimiter is a potential opener, a potential closer, or both (which depends on what sort of characters precede and follow the delimiters)
}

type delimiterStack struct {
	elements      []*delimiterStackElement
	openersBottom []*openerBottom
}

func (s *delimiterStack) matchOpener(e *delimiterStackElement) {
	length := len(s.elements) - 1
	for i := length; 0 <= i; i-- {
		t := s.elements[i]
		if e.typ == t.typ && e.num == t.num {
			e := &delimiterStackElement{}
			if 1 == e.num {
				e.node = &Emphasis{NodeType: NodeEmphasis}
			} else {
				e.node = &Strong{NodeType: NodeStrong}
			}

			for j := i + 1; j < len(s.elements); j++ {
				e.node.Append(s.elements[j].node)
			}
			s.elements = s.elements[:i]
			s.elements = append(s.elements, e)
		}
	}
}

func (s *delimiterStack) popMatch(e *delimiterStackElement) (elements []*delimiterStackElement) {
	for i := len(s.elements) - 1; 0 <= i; i-- {
		t := s.elements[i]
		if e.typ == t.typ && e.num == t.num {
			elements = append(s.elements[i:], e)
			s.elements = s.elements[:i]

			return
		}
	}

	return nil
}

func (s *delimiterStack) insert(position int, e *delimiterStackElement) {
	begin := append(s.elements[0:position], e)
	end := s.elements[position+1:]
	s.elements = append(begin, end...)
}

func (s *delimiterStack) removeDelimiters(begin, end int) {
	for i := begin; i < len(s.elements); i++ {
		if "*" == s.elements[i].typ {
			s.remove(i)
		}
	}
}

func (s *delimiterStack) removeRange(begin, end int) {
	b := s.elements[0:begin]
	e := s.elements[end:]
	s.elements = append(b, e...)

}

func (s *delimiterStack) remove(position int) {
	begin := s.elements[0:position]
	end := s.elements[position:]
	s.elements = append(begin, end...)
}

func (s *delimiterStack) push(e *delimiterStackElement) {
	s.elements = append(s.elements, e)
}

func (s *delimiterStack) pop() *delimiterStackElement {
	if 0 == len(s.elements) {
		return nil
	}

	return s.elements[len(s.elements)-1]
}

func (s *delimiterStack) popAll() []*delimiterStackElement {
	ret := s.elements
	s.elements = []*delimiterStackElement{}

	return ret
}

func (s *delimiterStack) peek() *delimiterStackElement {
	if 0 == len(s.elements) {
		return nil
	}

	return s.elements[len(s.elements)-1]
}

func (s *delimiterStack) isEmpty() bool {
	return 0 == len(s.elements)
}
