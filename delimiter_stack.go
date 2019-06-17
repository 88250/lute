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

type delimiterStackElement struct {
	node             *Text  // the text node point to
	typ              string // the type of delimiter ([, ![, *, _)
	num              int    // the number of delimiters
	active           bool   // whether the delimiter is "active" (all are active to start)
	openerCloserBoth string // whether the delimiter is a potential opener, a potential closer, or both (which depends on what sort of characters precede and follow the delimiters)
}

type delimiterStack struct {
	elements []*delimiterStackElement
	count    int
}

func (s *delimiterStack) popMatch(e *delimiterStackElement) (elements []*delimiterStackElement) {
	for i := s.count - 1; 0 <= i; i-- {
		t := s.elements[i]
		if e.typ == t.typ && e.num == t.num {
			s.count = i
			elements = append(s.elements[i:], e)
			s.elements = s.elements[:i]

			return
		}
	}

	return nil
}

func (s *delimiterStack) push(e *delimiterStackElement) {
	s.elements = append(s.elements[:s.count], e)
	s.count++
}

func (s *delimiterStack) pop() *delimiterStackElement {
	if 0 == s.count {
		return nil
	}

	s.count--

	return s.elements[s.count]
}

func (s *delimiterStack) popAll() []*delimiterStackElement {
	ret := s.elements
	s.count = 0
	s.elements = []*delimiterStackElement{}

	return ret
}

func (s *delimiterStack) peek() *delimiterStackElement {
	if 0 == s.count {
		return nil
	}

	return s.elements[s.count-1]
}

func (s *delimiterStack) isEmpty() bool {
	return 0 == len(s.elements)
}
